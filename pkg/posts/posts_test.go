package post

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/sessions"
	"reddit_backend/pkg/utils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	testCases := []struct {
		testPost       *repo.Post
		testUser       *repo.User
		reqParam       func(post *repo.Post) bytes.Buffer
		req            func(req *http.Request, user *repo.User) *http.Request
		postService    func(postDB *PostDB, session *sessions.Sessions) *Posts
		prepMockMongo  func(postDB *PostDB)
		respStatusCode int
		wantErr        bool
	}{
		{ //ok
			testPost: &repo.Post{Title: "test_title"},
			testUser: &repo.User{Username: "test3"},
			reqParam: func(post *repo.Post) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(post)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				ctx := context.WithValue(req.Context(), sessions.SessionKey, user)
				return req.WithContext(ctx)
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB) {
				postDB.On("CreatePost", mock.Anything, mock.Anything).Return(nil).Once()
			},
			respStatusCode: http.StatusCreated,
		},
		{ //repo return err
			testPost: &repo.Post{Title: "test_title"},
			testUser: &repo.User{Username: "test3"},
			reqParam: func(post *repo.Post) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(post)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				ctx := context.WithValue(req.Context(), sessions.SessionKey, user)
				return req.WithContext(ctx)
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB) {
				postDB.On("CreatePost", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			},
			respStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
		{ //401
			testPost: &repo.Post{Title: "test_title"},
			testUser: &repo.User{Username: "test3"},
			reqParam: func(post *repo.Post) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(post)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				return req
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB) {
				postDB.On("CreatePost", mock.Anything, mock.Anything).Return(nil).Once()
			},
			respStatusCode: http.StatusUnauthorized,
			wantErr:        true,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postDB := &PostDB{}
	sessionDB := sessions.NewMockSessionDB(ctrl)
	session := &sessions.Sessions{Repo: sessionDB, Secret: "test", Logger: nil}

	for _, tc := range testCases {
		reqParam := tc.reqParam(tc.testPost)
		req := httptest.NewRequest("POST", "/api/posts", &reqParam)
		reqWithContext := tc.req(req, tc.testUser)
		w := httptest.NewRecorder()
		postService := tc.postService(postDB, session)
		tc.prepMockMongo(postDB)
		postService.CreatePost(w, reqWithContext)
		resp := w.Result()
		require.EqualValues(t, tc.respStatusCode, resp.StatusCode)
		if tc.wantErr {
			continue
		}
		post := &repo.Post{}
		err := json.NewDecoder(resp.Body).Decode(&post)
		require.NoError(t, err)
		require.EqualValues(t, tc.testPost.Title, post.Title)

	}
}

func TestChangeVote(t *testing.T) {
	testCases := []struct {
		testPost       *repo.Post
		testUser       *repo.User
		url            string
		reqParam       func(post *repo.Post) bytes.Buffer
		req            func(req *http.Request, user *repo.User) *http.Request
		postService    func(postDB *PostDB, session *sessions.Sessions) *Posts
		prepMockMongo  func(postDB *PostDB, testPost *repo.Post)
		vote           int
		respStatusCode int
		wantErr        bool
	}{
		{ //ok upvote
			testPost: &repo.Post{Id: "123", Title: "test_title", Votes: []repo.Vote{{User: "1", Vote: -1}}},
			testUser: &repo.User{Id: "1", Username: "test1"},
			url:      "/api/post/123/upvote",
			reqParam: func(post *repo.Post) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(post)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				ctx := context.WithValue(req.Context(), sessions.SessionKey, user)
				return req.WithContext(ctx)
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB, testPost *repo.Post) {
				postDB.On("GetPostsById", mock.Anything, mock.Anything).Return(*testPost, nil).Once()
				postDB.On("UpdatePost", mock.Anything, mock.Anything).Return(nil).Once()
			},
			vote:           1,
			respStatusCode: http.StatusOK,
		},
		{ //ok upvote
			testPost: &repo.Post{Id: "123", Title: "test_title", Votes: []repo.Vote{{User: "1", Vote: -1}}},
			testUser: &repo.User{Id: "1", Username: "test1"},
			url:      "/api/post/123/unvote",
			reqParam: func(post *repo.Post) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(post)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				ctx := context.WithValue(req.Context(), sessions.SessionKey, user)
				return req.WithContext(ctx)
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB, testPost *repo.Post) {
				postDB.On("GetPostsById", mock.Anything, mock.Anything).Return(*testPost, nil).Once()
				postDB.On("UpdatePost", mock.Anything, mock.Anything).Return(nil).Once()
			},
			vote:           0,
			respStatusCode: http.StatusOK,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postDB := &PostDB{}
	sessionDB := sessions.NewMockSessionDB(ctrl)
	session := &sessions.Sessions{Repo: sessionDB, Secret: "test", Logger: nil}

	for _, tc := range testCases {
		reqParam := tc.reqParam(tc.testPost)
		req := httptest.NewRequest("GET", tc.url, &reqParam)
		reqWithContext := tc.req(req, tc.testUser)
		w := httptest.NewRecorder()
		postService := tc.postService(postDB, session)
		tc.prepMockMongo(postDB, tc.testPost)
		postService.ChangeVote(w, reqWithContext)
		resp := w.Result()
		require.EqualValues(t, tc.respStatusCode, resp.StatusCode)
		if tc.wantErr {
			continue
		}
		post := &repo.Post{}
		err := json.NewDecoder(resp.Body).Decode(&post)
		require.NoError(t, err)
		if tc.vote == 0 {
			require.Empty(t, post.Votes)
			continue
		}
		require.EqualValues(t, tc.vote, post.Votes[0].Vote)
	}
}

func TestCreateComment(t *testing.T) {
	testCases := []struct {
		testPost       *repo.Post
		testUser       *repo.User
		comment        *Comment
		reqParam       func(comment *Comment) bytes.Buffer
		req            func(req *http.Request, user *repo.User) *http.Request
		postService    func(postDB *PostDB, session *sessions.Sessions) *Posts
		prepMockMongo  func(postDB *PostDB, testPost *repo.Post)
		respStatusCode int
		wantErr        bool
	}{
		{ //ok
			testPost: &repo.Post{Title: "test_title", Comments: []repo.Comment{}},
			testUser: &repo.User{Username: "test6"},
			comment:  &Comment{Comment: "text_comment"},
			reqParam: func(comment *Comment) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(comment)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				ctx := context.WithValue(req.Context(), sessions.SessionKey, user)
				return req.WithContext(ctx)
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB, testPost *repo.Post) {
				postDB.On("GetPostsById", mock.Anything, mock.Anything).Return(*testPost, nil).Once()
				postDB.On("AddComment", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			respStatusCode: http.StatusCreated,
		},
		{ //401
			testPost: &repo.Post{Title: "test_title", Comments: []repo.Comment{}},
			testUser: &repo.User{Username: "test6"},
			comment:  &Comment{Comment: "text_comment"},
			reqParam: func(comment *Comment) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(comment)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				return req
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo:  func(postDB *PostDB, testPost *repo.Post) {},
			respStatusCode: http.StatusUnauthorized,
			wantErr:        true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postDB := &PostDB{}
	sessionDB := sessions.NewMockSessionDB(ctrl)
	session := &sessions.Sessions{Repo: sessionDB, Secret: "test", Logger: nil}

	for _, tc := range testCases {
		reqParam := tc.reqParam(tc.comment)
		req := httptest.NewRequest("POST", "/api/posts", &reqParam)
		reqWithContext := tc.req(req, tc.testUser)
		w := httptest.NewRecorder()
		postService := tc.postService(postDB, session)
		tc.prepMockMongo(postDB, tc.testPost)
		postService.CreateComment(w, reqWithContext)
		resp := w.Result()
		require.EqualValues(t, tc.respStatusCode, resp.StatusCode)
		if tc.wantErr {
			continue
		}
		post := &repo.Post{}
		err := json.NewDecoder(resp.Body).Decode(&post)
		require.NoError(t, err)
		require.EqualValues(t, tc.comment.Comment, post.Comments[0].Body)
	}
}

func TestDeletePost(t *testing.T) {
	testCases := []struct {
		testPost       *repo.Post
		testUser       *repo.User
		reqParam       func(post *repo.Post) bytes.Buffer
		req            func(req *http.Request, user *repo.User) *http.Request
		postService    func(postDB *PostDB, session *sessions.Sessions) *Posts
		prepMockMongo  func(postDB *PostDB, testPost *repo.Post)
		respStatusCode int
		wantErr        bool
	}{
		{ //ok
			testPost: &repo.Post{Id: "123", Title: "test_title"},
			testUser: &repo.User{Username: "test3"},
			reqParam: func(post *repo.Post) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(post)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				ctx := context.WithValue(req.Context(), sessions.SessionKey, user)
				return req.WithContext(ctx)
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB, testPost *repo.Post) {
				postDB.On("GetPostsById", mock.Anything, mock.Anything).Return(*testPost, nil).Once()
				postDB.On("DeletePostById", mock.Anything, mock.Anything).Return(nil).Once()
			},
			respStatusCode: http.StatusCreated,
		},
		{ //401
			testPost: &repo.Post{Id: "123", Title: "test_title"},
			testUser: &repo.User{Username: "test3"},
			reqParam: func(post *repo.Post) bytes.Buffer {
				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(post)
				require.NoError(t, err)
				return buf
			},
			req: func(req *http.Request, user *repo.User) *http.Request {
				return req
			},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo:  func(postDB *PostDB, testPost *repo.Post) {},
			respStatusCode: http.StatusUnauthorized,
			wantErr:        true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postDB := &PostDB{}
	sessionDB := sessions.NewMockSessionDB(ctrl)
	session := &sessions.Sessions{Repo: sessionDB, Secret: "test", Logger: nil}

	for _, tc := range testCases {
		reqParam := tc.reqParam(tc.testPost)
		req := httptest.NewRequest("DELETE", "/api/post/123", &reqParam)
		reqWithContext := tc.req(req, tc.testUser)
		w := httptest.NewRecorder()
		postService := tc.postService(postDB, session)
		tc.prepMockMongo(postDB, tc.testPost)
		postService.DeletePost(w, reqWithContext)
		resp := w.Result()
		require.EqualValues(t, tc.respStatusCode, resp.StatusCode)
		if tc.wantErr {
			continue
		}
		result := &utils.ReturnMsg{}
		err := json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		require.EqualValues(t, "success", result.Msg)
	}
}

func TestGetPostsByCat(t *testing.T) {
	testCases := []struct {
		testPost       *repo.Post
		postService    func(postDB *PostDB, session *sessions.Sessions) *Posts
		prepMockMongo  func(postDB *PostDB, testPost *repo.Post)
		respStatusCode int
		wantErr        bool
	}{
		{ //ok
			testPost: &repo.Post{Category: "music", Title: "test_title", Comments: []repo.Comment{}},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB, testPost *repo.Post) {
				postDB.On("GetPostsByCategory", mock.Anything, mock.Anything).Return([]repo.Post{*testPost}, nil).Once()
			},
			respStatusCode: http.StatusOK,
		},
		{ //not ok
			testPost: &repo.Post{Category: "music", Title: "test_title", Comments: []repo.Comment{}},
			postService: func(postDB *PostDB, session *sessions.Sessions) *Posts {
				return &Posts{Repo: postDB, Session: session, Logger: nil}
			},
			prepMockMongo: func(postDB *PostDB, testPost *repo.Post) {
				postDB.On("GetPostsByCategory", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			respStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postDB := &PostDB{}
	sessionDB := sessions.NewMockSessionDB(ctrl)
	session := &sessions.Sessions{Repo: sessionDB, Secret: "test", Logger: nil}

	for _, tc := range testCases {

		req := httptest.NewRequest("GET", "/api/posts/music", nil)
		w := httptest.NewRecorder()
		postService := tc.postService(postDB, session)
		tc.prepMockMongo(postDB, tc.testPost)
		postService.GetPostsByCategory(w, req)
		resp := w.Result()
		require.EqualValues(t, tc.respStatusCode, resp.StatusCode)
		if tc.wantErr {
			continue
		}
		post := []repo.Post{}
		err := json.NewDecoder(resp.Body).Decode(&post)
		require.NoError(t, err)
		require.EqualValues(t, tc.testPost.Category, post[0].Category)
	}
}
