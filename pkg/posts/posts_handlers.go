package post

import (
	"net/http"
	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/sessions"
	"reddit_backend/pkg/utils"
	"strings"

	"github.com/gorilla/mux"
)

type Posts struct {
	Repo    repo.PostDB
	Session *sessions.Sessions
	Logger  *logging.Logger
}

type Post struct {
	Category string
	Type     string
	Title    string
	Url      string
	Text     string
}

type Comment struct {
	Comment string
}

func (p *Posts) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := p.Repo.GetAllPosts(r.Context())
	if utils.HandleErr(w, r, err, http.StatusInternalServerError, p.Logger) {
		return
	}
	utils.HandleResult(w, r, posts, http.StatusOK, p.Logger)
}

func (p *Posts) GetPosts(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	username := param["user"]
	posts, err := p.Repo.GetPostsByUser(r.Context(), username)
	if utils.HandleErr(w, r, err, 500, p.Logger) {
		return
	}

	utils.HandleResult(w, r, posts, http.StatusOK, p.Logger)
}

func (p *Posts) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	category := param["category"]
	posts, err := p.Repo.GetPostsByCategory(r.Context(), category)
	if utils.HandleErr(w, r, err, http.StatusInternalServerError, p.Logger) {
		return
	}
	utils.HandleResult(w, r, posts, http.StatusOK, p.Logger)
}

func (p *Posts) GetPostsById(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id := param["id"]
	posts, err := p.Repo.GetPostsById(r.Context(), id)
	if utils.HandleErr(w, r, err, 400, p.Logger) {
		return
	}
	utils.HandleResult(w, r, posts, http.StatusOK, p.Logger)
}

func (p *Posts) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, err := p.Session.GetUserFromContext(r.Context())
	if utils.HandleErr(w, r, err, 401, p.Logger) {
		return
	}
	var post Post
	if utils.HandleErrDecodeReq(w, r, &post) {
		return
	}

	result, err := p.createPost(r.Context(), post, user)
	if utils.HandleErr(w, r, err, http.StatusInternalServerError, p.Logger) {
		return
	}
	utils.HandleResult(w, r, result, http.StatusCreated, p.Logger)
}

func (p *Posts) DeletePost(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	postId := param["id"]
	_, err := p.Session.GetUserFromContext(r.Context())
	if utils.HandleErr(w, r, err, 401, p.Logger) {
		return
	}
	err = p.Repo.DeletePostById(r.Context(), postId)
	if utils.HandleErr(w, r, err, 500, p.Logger) {
		return
	}
	utils.HandleResult(w, r, &utils.ReturnMsg{Msg: "success"}, http.StatusCreated, p.Logger)
}

func (p *Posts) DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["id"]
	commentId := params["comment_id"]
	post, err := p.deleteComment(r.Context(), postId, commentId)
	if utils.HandleErr(w, r, err, 404, p.Logger) {
		return
	}
	utils.HandleResult(w, r, &post, http.StatusCreated, p.Logger)
}

func (p *Posts) CreateComment(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	postId := param["id"]
	user, err := p.Session.GetUserFromContext(r.Context())
	if utils.HandleErr(w, r, err, 401, p.Logger) {
		return
	}
	var comment Comment
	if utils.HandleErrDecodeReq(w, r, &comment) {
		return
	}
	result, err := p.addCommentToPost(r.Context(), postId, user, comment.Comment)
	if utils.HandleErr(w, r, err, 404, p.Logger) {
		return
	}
	utils.HandleResult(w, r, result, http.StatusCreated, p.Logger)
}

func (p *Posts) ChangeVote(w http.ResponseWriter, r *http.Request) {
	var vote int
	if strings.Contains(r.RequestURI, "/unvote") {
		vote = 0
	} else if strings.Contains(r.RequestURI, "/upvote") {
		vote = 1
	} else {
		vote = -1
	}
	param := mux.Vars(r)
	postId := param["id"]
	user, err := p.Session.GetUserFromContext(r.Context())
	if utils.HandleErr(w, r, err, 401, p.Logger) {
		return
	}
	post, err := p.Repo.GetPostsById(r.Context(), postId)
	if utils.HandleErr(w, r, err, 400, p.Logger) {
		return
	}
	result, err := p.changeVote(r.Context(), post, user, vote)
	if utils.HandleErr(w, r, err, 500, p.Logger) {
		return
	}
	utils.HandleResult(w, r, result, http.StatusOK, p.Logger)
}
