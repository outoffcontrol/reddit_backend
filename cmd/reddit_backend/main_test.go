package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/utils"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {
	//mysql
	dsn := "root@tcp(localhost:3306)/coursera?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	mySql, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	mySql.SetMaxOpenConns(1)
	err = mySql.Ping()
	require.NoError(t, err)
	//mongo
	ctx := context.TODO()
	sess, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	collection := sess.Database("coursera").Collection("posts")
	//redis
	redisAddr := flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")
	redisConn, err := redis.DialURL(*redisAddr)
	require.NoError(t, err)

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiMjY5NSIsInVzZXJuYW1lIjoidGVzdF91c2VyIn0sInNlc3Npb25faWQiOiIyZGJkOTEwNS05N2U0LTRlYmQtODcxYi1hZmJjYTY4MDQzODQiLCJleHAiOjE2ODExMTE3MzEsImlhdCI6MTY4MDUwNjkzMX0.dMP4W_W1Wuhx4N1yf6YkTbaU2s4q4DqP3RVLNg2n2Wo"
	sessionId := "2dbd9105-97e4-4ebd-871b-afbca6804384"
	user := repo.User{Username: "test_user", Id: "2695"}
	post := repo.Post{Id: "a8d6d485-2831-4ad4-854d-7cc63d66439d", Score: 1, Views: 0, Type: "text", Title: "testtitle",
		Author: &user, Category: "music", Text: "texttext", Votes: []repo.Vote{{User: user.Id, Vote: 1}},
		Comments: []repo.Comment{{Id: "1234", Body: "test_comment",
			Author: &user, Created: "2023-03-25T16:48:06+03:00"}}, Created: "2023-03-25T16:48:06+03:00", UpvoteParcentage: 100}

	prepareMysqlDB(mySql)
	prepareMongoDB(collection, &ctx, post)
	prepareRedisDB(redisConn, sessionId)
	defer cleanupDB(mySql, collection, &ctx, redisConn)
	postsCollection := &repo.MongoCollection{
		Coll: collection,
	}

	testSiteMux := ServerMux(mySql, postsCollection, redisConn)

	ts := httptest.NewServer(testSiteMux)

	testCases := []Case{
		{ //ok tests
			path:   "/api/posts/",
			status: http.StatusOK,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := []repo.Post{post}
				posts := []repo.Post{}
				err := json.NewDecoder(r.Body).Decode(&posts)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, posts)
			},
		},
		{
			path:   "/api/posts/music",
			status: http.StatusOK,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := []repo.Post{post}
				posts := []repo.Post{}
				err := json.NewDecoder(r.Body).Decode(&posts)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, posts)
			},
		},
		{
			path:   "/api/post/a8d6d485-2831-4ad4-854d-7cc63d66439d",
			status: http.StatusOK,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := post
				posts := repo.Post{}
				err := json.NewDecoder(r.Body).Decode(&posts)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, posts)
			},
		},
		{
			path:   "/api/post/a8d6d485-2831-4ad4-854d-7cc63d66439d/downvote",
			status: http.StatusOK,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := post
				expectResult.Votes[0].Vote = -1
				expectResult.Score = -1
				expectResult.UpvoteParcentage = 0
				posts := repo.Post{}
				err := json.NewDecoder(r.Body).Decode(&posts)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, posts)
			},
			withAuth: true,
			token:    token,
		},
		{
			path:   "/api/login",
			method: "POST",
			status: http.StatusOK,
			isExpected: func(t *testing.T, r *http.Response) {
				result := &utils.ReturnToken{}
				err := json.NewDecoder(r.Body).Decode(&result)
				require.NoError(t, err)
				require.NotEmpty(t, result)
			},
			body: map[string]interface{}{"username": user.Username, "password": "12345678"},
		},
		{
			path:   "/api/register",
			method: "POST",
			status: http.StatusCreated,
			isExpected: func(t *testing.T, r *http.Response) {
				result := &utils.ReturnToken{}
				err := json.NewDecoder(r.Body).Decode(&result)
				require.NoError(t, err)
				require.NotEmpty(t, result)
			},
			body: map[string]interface{}{"username": "newuser", "password": "12345678"},
		},
		{
			path:   "/api/post/a8d6d485-2831-4ad4-854d-7cc63d66439d",
			method: "DELETE",
			status: http.StatusCreated,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := &utils.ReturnMsg{Msg: "success"}
				result := &utils.ReturnMsg{}
				err := json.NewDecoder(r.Body).Decode(&result)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, result)
			},
			withAuth: true,
			token:    token,
		},
		{
			path:   "/api/posts",
			method: "POST",
			status: http.StatusCreated,
			isExpected: func(t *testing.T, r *http.Response) {
				posts := repo.Post{}
				err := json.NewDecoder(r.Body).Decode(&posts)
				require.NoError(t, err)
				require.NotEmpty(t, posts)
				post = posts
			},
			body:     post,
			withAuth: true,
			token:    token,
		},
		{ //not ok tests
			path:   "/api/posts",
			method: "POST",
			status: http.StatusUnauthorized,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := &utils.ReturnMsg{Msg: "unauthorized"}
				result := &utils.ReturnMsg{}
				err := json.NewDecoder(r.Body).Decode(&result)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, result)
			},
			body:     post,
			withAuth: true,
			token:    token + "1",
		},
		{
			path:   "/api/login",
			method: "POST",
			status: http.StatusUnauthorized,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := &utils.ReturnMsg{Msg: "invalid password"}
				result := &utils.ReturnMsg{}
				err := json.NewDecoder(r.Body).Decode(&result)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, result)
			},
			body: map[string]interface{}{"username": user.Username, "password": "1234567"},
		},
		{
			path:   "/api/post/123/123",
			method: "DELETE",
			status: http.StatusNotFound,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := &utils.ReturnMsg{Msg: "post not found"}
				result := &utils.ReturnMsg{}
				err := json.NewDecoder(r.Body).Decode(&result)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, result)
			},
			withAuth: true,
			token:    token,
		},
		{
			path:   "/api/register",
			method: "POST",
			status: http.StatusUnprocessableEntity,
			isExpected: func(t *testing.T, r *http.Response) {
				expectResult := &utils.Return422Errors{Errors: []utils.Return422Error{
					{
						Location: "body", Param: "username", Value: "newuser", Msg: "username already exists"},
				}}
				result := &utils.Return422Errors{}
				err := json.NewDecoder(r.Body).Decode(&result)
				require.NoError(t, err)
				require.EqualValues(t, expectResult, result)
			},
			body: map[string]interface{}{"username": "newuser", "password": "12345678"},
		},
	}
	runCases(t, ts, mySql, testCases)
}

func runCases(t *testing.T, ts *httptest.Server, db *sql.DB, cases []Case) {
	client := &http.Client{Timeout: time.Second}

	for _, item := range cases {
		var req *http.Request
		var err error
		if db.Stats().OpenConnections != 1 {
			t.Fatalf("you have %d open connections, must be 1", db.Stats().OpenConnections)
		}
		if item.method == "" || item.method == http.MethodGet {
			req, err = http.NewRequest(item.method, ts.URL+item.path, nil)
			require.NoError(t, err)
		} else {
			data, err := json.Marshal(item.body)
			require.NoError(t, err)
			reqBody := bytes.NewReader(data)
			req, err = http.NewRequest(item.method, ts.URL+item.path, reqBody)
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/json")
		}
		if item.withAuth {
			req.Header.Add("Authorization", "Bearer "+item.token)
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.EqualValues(t, item.status, resp.StatusCode)
		item.isExpected(t, resp)

	}
}

type Case struct {
	method     string
	path       string
	status     int
	isExpected func(t *testing.T, r *http.Response)
	body       interface{}
	withAuth   bool
	token      string
}

func prepareRedisDB(redisConn redis.Conn, sessionId string) {
	_, err := redisConn.Do("FLUSHDB")
	if err != nil {
		log.Fatalf("cannot flush redis: %s", err)
	}
	sesOpts := []string{time.Now().Add(5 * time.Hour).Format("2006-01-02T15:04:05Z07:00")}
	_, err = redisConn.Do("HSET", "2695", sessionId, sesOpts)
	if err != nil {
		log.Fatalf("cant hset session_id to Redis: %s", err)
	}
}

func prepareMysqlDB(db *sql.DB) {
	sqlReq := []string{
		`DROP TABLE IF EXISTS users;`,

		`CREATE TABLE users (
			id varchar(255) NOT NULL,
			username varchar(255) NOT NULL,
			password varchar(255) NOT NULL,
			PRIMARY KEY (id)
		  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`INSERT INTO users (id, username, password) VALUES
		("2695",	'test_user',	'12345678');`,
	}

	for _, r := range sqlReq {
		_, err := db.Exec(r)
		if err != nil {
			log.Fatalf("cannot exec req to mysql: %s", err)
		}
	}
}

func prepareMongoDB(db *mongo.Collection, ctx *context.Context, post repo.Post) {
	err := db.Drop(*ctx)
	if err != nil {
		log.Fatalf("cannot drop mongo collection: %s", err)
	}

	_, err = db.InsertOne(*ctx, &post)
	if err != nil {
		log.Fatalf("cannot insert post to mongo: %s", err)
	}
}

func cleanupDB(dbSql *sql.DB, dbMongo *mongo.Collection, ctx *context.Context, redisConn redis.Conn) {
	_, err := dbSql.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		log.Fatalf("cannot drop db mysql: %s", err)
	}

	err = dbMongo.Drop(*ctx)
	if err != nil {
		log.Fatalf("cannot drop mongo collection: %s", err)
	}

	_, err = redisConn.Do("FLUSHDB")
	if err != nil {
		log.Fatalf("cannot flush redis: %s", err)
	}

}
