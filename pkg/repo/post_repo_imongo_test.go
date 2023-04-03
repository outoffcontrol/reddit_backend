package repo

import (
	"context"
	"errors"
	"reddit_backend/pkg/logging"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

type MockMongoStruct struct {
	mockCollection   *MockIMongoCollection
	mockSingleResult *MockIMongoSingleResult
	mockCursor       *MockIMongoCursor
	mockUpdateResult *MockIMongoUpdateResult
	ctx              context.Context
}

func TestGetAllPostsImongo(t *testing.T) {
	testCases := []struct {
		name          string
		testPost      *Post
		prepMongoMock func(mockMongoStruct *MockMongoStruct, testPost *Post)
		repoPost      func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo
		wantErr       string
	}{
		{
			name: "get all posts",
			testPost: &Post{Id: "a8d6d485-2831-4ad4-854d-7cc63d66439d", Score: 0, Views: 0, Type: "text", Title: "testtitle",
				Author: nil, Category: "music", Text: "texttext", Votes: nil, Comments: nil, Created: "2023-03-25T16:48:06+03:00", UpvoteParcentage: 50},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				var results []Post
				results = append(results, *testPost)
				mockMongoStruct.mockCollection.EXPECT().Find(mockMongoStruct.ctx, gomock.Any()).Return(mockMongoStruct.mockCursor, nil)
				mockMongoStruct.mockCursor.EXPECT().All(mockMongoStruct.ctx, gomock.Any()).SetArg(1, results).Return(nil)
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
		},
		{
			name:     "get all posts (cannot decode all posts)",
			testPost: &Post{},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				mockMongoStruct.mockCollection.EXPECT().Find(mockMongoStruct.ctx, gomock.Any()).Return(mockMongoStruct.mockCursor, nil)
				mockMongoStruct.mockCursor.EXPECT().All(mockMongoStruct.ctx, gomock.Any()).
					Return(errors.New("error decoding key id: cannot decode 32-bit integer into a string type"))
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
			wantErr: "error decoding key id: cannot decode 32-bit integer into a string type",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	zapLogger := zap.NewNop()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	ctx := context.TODO()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}
	mockCollection := NewMockIMongoCollection(ctrl)
	mockSingleResult := NewMockIMongoSingleResult(ctrl)
	mockCursor := NewMockIMongoCursor(ctrl)
	mockMongoStruct := &MockMongoStruct{
		mockCollection:   mockCollection,
		mockSingleResult: mockSingleResult,
		mockCursor:       mockCursor,
		ctx:              ctx,
	}
	for _, tc := range testCases {
		tc.prepMongoMock(mockMongoStruct, tc.testPost)
		repoPost := tc.repoPost(loggerCtx, mockMongoStruct.mockCollection, ctx)
		resp, err := repoPost.GetAllPosts(ctx)
		if err != nil {
			if err.Error() != tc.wantErr {
				t.Errorf("expected err: %s, got: %s", tc.wantErr, err)
				return
			}
			return
		}
		require.NotEmpty(t, resp[0])
		require.EqualValues(t, resp[0].Id, tc.testPost.Id)

	}
}

func TestGetPostsByIdImongo(t *testing.T) {
	testCases := []struct {
		name          string
		testPost      *Post
		prepMongoMock func(mockMongoStruct *MockMongoStruct, testPost *Post)
		repoPost      func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo
		wantErr       string
	}{
		{
			name: "get post by ID",
			testPost: &Post{Id: "a8d6d485-2831-4ad4-854d-7cc63d66439d", Score: 0, Views: 0, Type: "text", Title: "testtitle",
				Author: nil, Category: "music", Text: "texttext", Votes: nil, Comments: nil, Created: "2023-03-25T16:48:06+03:00", UpvoteParcentage: 50},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				mockMongoStruct.mockCollection.EXPECT().FindOne(mockMongoStruct.ctx, gomock.Any()).Return(mockMongoStruct.mockSingleResult)
				mockMongoStruct.mockSingleResult.EXPECT().Decode(gomock.Any()).SetArg(0, testPost).Return(nil)
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
		},
		{
			name:     "get posts by ID (error from Mongo)",
			testPost: &Post{},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				mockMongoStruct.mockCollection.EXPECT().FindOne(mockMongoStruct.ctx, gomock.Any()).Return(mockMongoStruct.mockSingleResult)
				mockMongoStruct.mockSingleResult.EXPECT().Decode(gomock.Any()).
					Return(errors.New("write command error: [{write errors: [{error}]}, {<nil>}]"))
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
			wantErr: "write command error: [{write errors: [{error}]}, {<nil>}]",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	zapLogger := zap.NewNop()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	ctx := context.TODO()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}
	mockCollection := NewMockIMongoCollection(ctrl)
	mockSingleResult := NewMockIMongoSingleResult(ctrl)
	mockCursor := NewMockIMongoCursor(ctrl)
	mockMongoStruct := &MockMongoStruct{
		mockCollection:   mockCollection,
		mockSingleResult: mockSingleResult,
		mockCursor:       mockCursor,
		ctx:              ctx,
	}
	for _, tc := range testCases {
		tc.prepMongoMock(mockMongoStruct, tc.testPost)
		repoPost := tc.repoPost(loggerCtx, mockMongoStruct.mockCollection, ctx)
		resp, err := repoPost.GetPostsById(ctx, tc.testPost.Id)
		if err != nil {
			if err.Error() != tc.wantErr {
				t.Errorf("expected err: %s, got: %s", tc.wantErr, err)
				return
			}
			return
		}
		require.NotEmpty(t, resp)
		require.EqualValues(t, resp.Id, tc.testPost.Id)

	}
}

func TestUpdatePostImongo(t *testing.T) {
	testCases := []struct {
		name          string
		comment       string
		testPost      *Post
		testUser      *User
		prepMongoMock func(mockMongoStruct *MockMongoStruct, testPost *Post)
		repoPost      func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo
		wantErr       string
	}{
		{
			name: "ok",
			testPost: &Post{Id: "a8d6d485-2831-4ad4-854d-7cc63d66439d", Score: 0, Views: 0, Type: "text", Title: "testtitle",
				Author: nil, Category: "music", Text: "texttext", Votes: nil, Comments: nil, Created: "2023-03-25T16:48:06+03:00", UpvoteParcentage: 50},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				mockMongoStruct.mockCollection.EXPECT().UpdateOne(mockMongoStruct.ctx, gomock.Any(), gomock.Any()).
					Return(mockMongoStruct.mockUpdateResult, nil)
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
		},
		{
			name:     "not ok)",
			testPost: &Post{},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				mockMongoStruct.mockCollection.EXPECT().UpdateOne(mockMongoStruct.ctx, gomock.Any(), gomock.Any()).
					Return(nil, errors.New("no responses remaining"))
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
			wantErr: "no responses remaining",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	zapLogger := zap.NewNop()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	ctx := context.TODO()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}
	mockCollection := NewMockIMongoCollection(ctrl)
	mockSingleResult := NewMockIMongoSingleResult(ctrl)
	mockCursor := NewMockIMongoCursor(ctrl)
	mockUpdateResult := NewMockIMongoUpdateResult(ctrl)
	mockMongoStruct := &MockMongoStruct{
		mockCollection:   mockCollection,
		mockSingleResult: mockSingleResult,
		mockCursor:       mockCursor,
		mockUpdateResult: mockUpdateResult,
		ctx:              ctx,
	}
	for _, tc := range testCases {
		tc.prepMongoMock(mockMongoStruct, tc.testPost)
		repoPost := tc.repoPost(loggerCtx, mockMongoStruct.mockCollection, ctx)
		err := repoPost.UpdatePost(ctx, *tc.testPost)
		if err != nil {
			if err.Error() != tc.wantErr {
				t.Errorf("expected err: %s, got: %s", tc.wantErr, err)
				return
			}
			return
		}
	}
}

func TestGetPostsByCatImongo(t *testing.T) {
	testCases := []struct {
		name          string
		testPost      *Post
		prepMongoMock func(mockMongoStruct *MockMongoStruct, testPost *Post)
		repoPost      func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo
		wantErr       string
	}{
		{
			name: "get post by Category",
			testPost: &Post{Id: "a8d6d485-2831-4ad4-854d-7cc63d66439d", Score: 0, Views: 0, Type: "text", Title: "testtitle",
				Author: nil, Category: "music", Text: "texttext", Votes: nil, Comments: nil, Created: "2023-03-25T16:48:06+03:00", UpvoteParcentage: 50},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				var results []Post
				results = append(results, *testPost)
				mockMongoStruct.mockCollection.EXPECT().Find(mockMongoStruct.ctx, gomock.Any()).Return(mockMongoStruct.mockCursor, nil)
				mockMongoStruct.mockCursor.EXPECT().All(mockMongoStruct.ctx, gomock.Any()).SetArg(1, results).Return(nil)
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
		},
		{
			name:     "get posts by Category (error from Mongo)",
			testPost: &Post{},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				mockMongoStruct.mockCollection.EXPECT().Find(mockMongoStruct.ctx, gomock.Any()).
					Return(nil, errors.New("write command error: [{write errors: [{error}]}, {<nil>}]"))
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
			wantErr: "write command error: [{write errors: [{error}]}, {<nil>}]",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	zapLogger := zap.NewNop()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	ctx := context.TODO()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}
	mockCollection := NewMockIMongoCollection(ctrl)
	mockSingleResult := NewMockIMongoSingleResult(ctrl)
	mockCursor := NewMockIMongoCursor(ctrl)
	mockMongoStruct := &MockMongoStruct{
		mockCollection:   mockCollection,
		mockSingleResult: mockSingleResult,
		mockCursor:       mockCursor,
		ctx:              ctx,
	}
	for _, tc := range testCases {
		tc.prepMongoMock(mockMongoStruct, tc.testPost)
		repoPost := tc.repoPost(loggerCtx, mockMongoStruct.mockCollection, ctx)
		resp, err := repoPost.GetPostsByCategory(ctx, tc.testPost.Category)
		if err != nil {
			if err.Error() != tc.wantErr {
				t.Errorf("expected err: %s, got: %s", tc.wantErr, err)
				return
			}
			return
		}
		require.NotEmpty(t, resp)
		require.EqualValues(t, tc.testPost.Category, resp[0].Category)
	}
}

func TestGetPostsByUserImongo(t *testing.T) {
	testCases := []struct {
		name          string
		testPost      *Post
		prepMongoMock func(mockMongoStruct *MockMongoStruct, testPost *Post)
		repoPost      func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo
		wantErr       string
	}{
		{
			name: "get post by User",
			testPost: &Post{Id: "a8d6d485-2831-4ad4-854d-7cc63d66439d", Score: 0, Views: 0, Type: "text", Title: "testtitle",
				Author: &User{Username: "test7", Id: "1234"}, Category: "music", Text: "texttext",
				Votes: nil, Comments: nil, Created: "2023-03-25T16:48:06+03:00", UpvoteParcentage: 50},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				var results []Post
				results = append(results, *testPost)
				mockMongoStruct.mockCollection.EXPECT().Find(mockMongoStruct.ctx, gomock.Any()).Return(mockMongoStruct.mockCursor, nil)
				mockMongoStruct.mockCursor.EXPECT().All(mockMongoStruct.ctx, gomock.Any()).SetArg(1, results).Return(nil)
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
		},
		{
			name: "get posts by User (error from Mongo)",
			testPost: &Post{Id: "a8d6d485-2831-4ad4-854d-7cc63d66439d", Score: 0, Views: 0, Type: "text", Title: "testtitle",
				Author: &User{Username: "test7", Id: "1234"}, Category: "music", Text: "texttext",
				Votes: nil, Comments: nil, Created: "2023-03-25T16:48:06+03:00", UpvoteParcentage: 50},
			prepMongoMock: func(mockMongoStruct *MockMongoStruct, testPost *Post) {
				mockMongoStruct.mockCollection.EXPECT().Find(mockMongoStruct.ctx, gomock.Any()).
					Return(nil, errors.New("write command error: [{write errors: [{error}]}, {<nil>}]"))
			},
			repoPost: func(logger *logging.Logger, mongo IMongoCollection, ctx context.Context) *PostMongo {
				return &PostMongo{Logger: logger, Mongo: mongo, Context: &ctx}
			},
			wantErr: "write command error: [{write errors: [{error}]}, {<nil>}]",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	zapLogger := zap.NewNop()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	ctx := context.TODO()
	loggerCtx := &logging.Logger{Zap: logger, RequestIDKey: "requestID", LoggerKey: "logger"}
	mockCollection := NewMockIMongoCollection(ctrl)
	mockSingleResult := NewMockIMongoSingleResult(ctrl)
	mockCursor := NewMockIMongoCursor(ctrl)
	mockMongoStruct := &MockMongoStruct{
		mockCollection:   mockCollection,
		mockSingleResult: mockSingleResult,
		mockCursor:       mockCursor,
		ctx:              ctx,
	}
	for _, tc := range testCases {
		tc.prepMongoMock(mockMongoStruct, tc.testPost)
		repoPost := tc.repoPost(loggerCtx, mockMongoStruct.mockCollection, ctx)
		resp, err := repoPost.GetPostsByUser(ctx, tc.testPost.Author.Username)
		if err != nil {
			if err.Error() != tc.wantErr {
				t.Errorf("expected err: %s, got: %s", tc.wantErr, err)
				return
			}
			return
		}
		require.NotEmpty(t, resp)
		require.EqualValues(t, tc.testPost.Author.Username, resp[0].Author.Username)
	}
}
