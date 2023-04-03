package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/sessions"
	"reddit_backend/pkg/utils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		testUser       *repo.User
		password       string
		userReq        func(username string, password string) bytes.Buffer
		userService    func(repo *MockUserDB, session *sessions.Sessions) *Users
		prepMockSql    func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, username string, password string)
		respStatusCode int
		wantErr        bool
	}{
		{ //ok
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				req := &UserReq{Username: username, Password: password}
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(req)
				require.NoError(t, err)
				return buf
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, username string, password string) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), username, password).Return(respUser, nil)
				sessionDB.EXPECT().AddSession(gomock.Any(), respUser.Id, gomock.Any()).Return(nil)
			},
			respStatusCode: http.StatusCreated,
		},
		{ //user exist
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				req := &UserReq{Username: username, Password: password}
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(req)
				require.NoError(t, err)
				return buf
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, username string, password string) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), username, password).Return(nil, utils.ErrUserExist)
			},
			respStatusCode: http.StatusUnprocessableEntity,
			wantErr:        true,
		},
		{ //empty req param
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				var badReq bytes.Buffer
				return badReq
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, username string, password string) {
			},
			respStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},

		{ //internal error(mysql)
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				req := &UserReq{Username: username, Password: password}
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(req)
				require.NoError(t, err)
				return buf
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, username string, password string) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), username, password).Return(nil, utils.ErrUserExist)
			},
			respStatusCode: http.StatusUnprocessableEntity,
			wantErr:        true,
		},
		{ //internal err (error create token)
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				req := &UserReq{Username: username, Password: password}
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(req)
				require.NoError(t, err)
				return buf
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, username string, password string) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), username, password).Return(respUser, nil)
				sessionDB.EXPECT().AddSession(gomock.Any(), respUser.Id, gomock.Any()).Return(utils.ErrCreateToken)
			},
			respStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userDB := NewMockUserDB(ctrl)
	sessionDB := sessions.NewMockSessionDB(ctrl)
	session := &sessions.Sessions{Repo: sessionDB, Secret: "test", Logger: nil}

	for _, tc := range testCases {
		testUser := tc.testUser
		tc.prepMockSql(userDB, sessionDB, testUser, tc.testUser.Username, tc.password)
		userService := tc.userService(userDB, session)
		userReq := tc.userReq(testUser.Username, tc.password)
		req := httptest.NewRequest("POST", "/api/register", &userReq)
		w := httptest.NewRecorder()
		userService.CreateUser(w, req)
		resp := w.Result()
		require.EqualValues(t, tc.respStatusCode, resp.StatusCode)
		if tc.wantErr {
			continue
		}
		token := &utils.ReturnToken{}
		err := json.NewDecoder(resp.Body).Decode(&token)
		require.NoError(t, err)
		require.NotEmpty(t, token.Token)

	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		testUser       *repo.User
		password       string
		userReq        func(username string, password string) bytes.Buffer
		userService    func(repo *MockUserDB, session *sessions.Sessions) *Users
		prepMockSql    func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, password string)
		respStatusCode int
		wantErr        bool
	}{
		{ //ok
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				req := &UserReq{Username: username, Password: password}
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(req)
				require.NoError(t, err)
				return buf
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, password string) {
				repo.EXPECT().CheckPasswordUser(gomock.Any(), respUser.Username, password).Return(respUser, nil)
				sessionDB.EXPECT().AddSession(gomock.Any(), respUser.Id, gomock.Any()).Return(nil)
			},
			respStatusCode: http.StatusOK,
		},
		{ //bad password
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				req := &UserReq{Username: username, Password: password}
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(req)
				require.NoError(t, err)
				return buf
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, password string) {
				repo.EXPECT().CheckPasswordUser(gomock.Any(), respUser.Username, password).Return(nil, utils.ErrBadPass)
			},
			respStatusCode: http.StatusUnauthorized,
			wantErr:        true,
		},
		{ //empty req param
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				var badReq bytes.Buffer
				return badReq
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql:    func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, password string) {},
			respStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{ //internal err (error create token)
			testUser: &repo.User{Id: utils.GenerateId(), Username: "test2"},
			password: "123123123",
			userReq: func(username string, password string) bytes.Buffer {
				req := &UserReq{Username: username, Password: password}
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(req)
				require.NoError(t, err)
				return buf
			},
			userService: func(repo *MockUserDB, session *sessions.Sessions) *Users {
				return &Users{repo, session, nil}
			},
			prepMockSql: func(repo *MockUserDB, sessionDB *sessions.MockSessionDB, respUser *repo.User, password string) {
				repo.EXPECT().CheckPasswordUser(gomock.Any(), respUser.Username, password).Return(respUser, nil)
				sessionDB.EXPECT().AddSession(gomock.Any(), respUser.Id, gomock.Any()).Return(utils.ErrCreateToken)
			},
			respStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userDB := NewMockUserDB(ctrl)
	sessionDB := sessions.NewMockSessionDB(ctrl)
	session := &sessions.Sessions{Repo: sessionDB, Secret: "test", Logger: nil}
	for _, tc := range testCases {
		testUser := tc.testUser
		tc.prepMockSql(userDB, sessionDB, testUser, tc.password)
		userService := tc.userService(userDB, session)
		userReq := tc.userReq(testUser.Username, tc.password)
		req := httptest.NewRequest("POST", "/api/register", &userReq)
		w := httptest.NewRecorder()
		userService.Login(w, req)
		resp := w.Result()
		require.EqualValues(t, tc.respStatusCode, resp.StatusCode)
		if tc.wantErr {
			continue
		}
		token := &utils.ReturnToken{}
		err := json.NewDecoder(resp.Body).Decode(&token)
		require.NoError(t, err)
		require.NotEmpty(t, token.Token)

	}
}
