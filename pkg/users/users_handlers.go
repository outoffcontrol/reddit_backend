package users

import (
	"net/http"
	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/sessions"
	"reddit_backend/pkg/utils"
	//"github.com/gorilla/mux"
)

type Users struct {
	Repo    repo.UserDB
	Session *sessions.Sessions
	Logger  *logging.Logger
}

type UserReq struct {
	Password string
	Username string
}

func (u *Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userReq UserReq
	if utils.HandleErrDecodeReq(w, r, &userReq) {
		return
	}
	token, err := u.createUser(r.Context(), &userReq)
	if err != nil {
		if err == utils.ErrUserExist {
			utils.HandleErr(w, r, err, http.StatusUnprocessableEntity, u.Logger, utils.Return422Error{
				Location: "body",
				Param:    "username",
				Value:    userReq.Username,
				Msg:      err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.HandleResult(w, r, &utils.ReturnToken{Token: token}, http.StatusCreated, u.Logger)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var userReq UserReq
	if utils.HandleErrDecodeReq(w, r, &userReq) {
		return
	}
	user, err := u.Repo.CheckPasswordUser(r.Context(), userReq.Username, userReq.Password)
	if utils.HandleErr(w, r, err, http.StatusUnauthorized, u.Logger) {
		return
	}
	token, err := u.Session.CreateToken(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.HandleResult(w, r, &utils.ReturnToken{Token: token}, http.StatusOK, u.Logger)
}
