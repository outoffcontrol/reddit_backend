package main

import (
	"context"
	"database/sql"
	"flag"
	"net/http"

	//"sync"
	"time"

	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/middleware"
	posts "reddit_backend/pkg/posts"
	repo "reddit_backend/pkg/repo"
	"reddit_backend/pkg/sessions"
	users "reddit_backend/pkg/users"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {

	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	//mutex for map data
	//mu := &sync.RWMutex{}
	//sessionMap := repo.SessionsMap{Mu: mu, Sessions: make(map[string]*repo.Sessions)}
	//userMap := repo.UserMap{Mu: mu, Users: make(map[string]*repo.User)}
	//postMap := repo.PostMap{Mu: mu, Posts: make(map[string]*repo.Post)}

	//session redis
	redisAddr := flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")
	redisConn, err := redis.DialURL(*redisAddr)
	if err != nil {
		logger.Fatalf("cant connect to redis")
	}

	//users mysql
	dsn := "root@tcp(localhost:3306)/coursera?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	mySql, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("cant open mysql: ", err)
	}
	mySql.SetMaxOpenConns(10)
	err = mySql.Ping()
	if err != nil {
		logger.Fatalf("cant connect to mysql: ", err)
	}

	//posts mongodb
	ctx := context.TODO()
	sess, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logger.Fatalf("cant connect to mongodb: ", err)
	}
	collection := sess.Database("coursera").Collection("posts")
	postsCollection := &repo.MongoCollection{
		Coll: collection,
	}

	newSiteMux := ServerMux(mySql, postsCollection, redisConn)

	srv := &http.Server{
		Handler:      newSiteMux,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Infow("starting server",
		"addr", srv.Addr,
	)
	srv.ListenAndServe()
}

func ServerMux(dbSql *sql.DB, dbMongo repo.IMongoCollection, redisConn redis.Conn) *http.ServeMux {
	ctx := context.TODO()
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}

	sessionRedis := repo.SessionsRedis{Logger: loggerCtx, RedisConn: redisConn}
	userMySql := repo.UserMySql{Logger: loggerCtx, MySql: dbSql}
	postMongo := repo.PostMongo{Logger: loggerCtx, Mongo: dbMongo, Context: &ctx}
	repoManger := repo.NewRepoManager(&sessionRedis, &userMySql, &postMongo)

	session := &sessions.Sessions{Repo: repoManger.Sessions, Secret: "test", Logger: loggerCtx}
	m := middleware.Middleware{Repo: repoManger, Session: session, Logger: loggerCtx}
	posts := &posts.Posts{Repo: repoManger.Posts, Session: session, Logger: loggerCtx}
	users := &users.Users{Repo: repoManger.Users, Session: session, Logger: loggerCtx}

	r1 := mux.NewRouter()
	r1.HandleFunc("/api/posts/", posts.GetAll).Methods("GET")
	r1.HandleFunc("/api/posts/{category}", posts.GetPostsByCategory).Methods("GET")
	r1.HandleFunc("/api/post/{id}", posts.GetPostsById).Methods("GET")
	r1.HandleFunc("/api/post/{id}/upvote", m.Auth(posts.ChangeVote)).Methods("GET")
	r1.HandleFunc("/api/post/{id}/downvote", m.Auth(posts.ChangeVote)).Methods("GET")
	r1.HandleFunc("/api/post/{id}/unvote", m.Auth(posts.ChangeVote)).Methods("GET")
	r1.HandleFunc("/api/post/{id}", m.Auth(posts.CreateComment)).Methods("POST")
	r1.HandleFunc("/api/posts", m.Auth(posts.CreatePost)).Methods("POST")
	r1.HandleFunc("/api/post/{id}", m.Auth(posts.DeletePost)).Methods("DELETE")
	r1.HandleFunc("/api/post/{id}/{comment_id}", m.Auth(posts.DeleteComment)).Methods("DELETE")
	r1.HandleFunc("/api/user/{user}", posts.GetPosts).Methods("GET")
	r1.HandleFunc("/api/register", users.CreateUser).Methods("POST")
	r1.HandleFunc("/api/login", users.Login).Methods("POST")
	r1.Use(loggerCtx.SetupReqID)
	r1.Use(loggerCtx.SetupLogger)
	r1.Use(m.LoggingMiddleware)
	r1.Use(m.Panic)

	//static
	r2 := mux.NewRouter()
	dir := "./web/redditclone/template/"
	fileServer := http.FileServer(http.Dir(dir))
	r2.PathPrefix("/").Handler(http.StripPrefix("/", fileServer)).Methods("GET")
	//r2.Use(m.LoggingMiddleware)

	r2.Use(middleware.StripUrlMiddleware)
	r2.Use(m.Panic)

	siteMux := http.NewServeMux()
	siteMux.Handle("/api/", r1)
	siteMux.Handle("/", r2)
	return siteMux
}
