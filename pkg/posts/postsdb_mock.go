// Code generated by mockery v2.20.0. DO NOT EDIT.

package post

import (
	context "context"
	repo "reddit_backend/pkg/repo"

	mock "github.com/stretchr/testify/mock"
)

// PostDB is an autogenerated mock type for the PostDB type
type PostDB struct {
	mock.Mock
}

// AddComment provides a mock function with given fields: ctx, postId, comment
func (_m *PostDB) AddComment(ctx context.Context, postId string, comment repo.Comment) error {
	ret := _m.Called(ctx, postId, comment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, repo.Comment) error); ok {
		r0 = rf(ctx, postId, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePost provides a mock function with given fields: ctx, _a1
func (_m *PostDB) CreatePost(ctx context.Context, _a1 repo.Post) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.Post) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePostById provides a mock function with given fields: ctx, id
func (_m *PostDB) DeletePostById(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllPosts provides a mock function with given fields: ctx
func (_m *PostDB) GetAllPosts(ctx context.Context) ([]repo.Post, error) {
	ret := _m.Called(ctx)

	var r0 []repo.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]repo.Post, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []repo.Post); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repo.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostsByCategory provides a mock function with given fields: ctx, category
func (_m *PostDB) GetPostsByCategory(ctx context.Context, category string) ([]repo.Post, error) {
	ret := _m.Called(ctx, category)

	var r0 []repo.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]repo.Post, error)); ok {
		return rf(ctx, category)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []repo.Post); ok {
		r0 = rf(ctx, category)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repo.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, category)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostsById provides a mock function with given fields: ctx, id
func (_m *PostDB) GetPostsById(ctx context.Context, id string) (repo.Post, error) {
	ret := _m.Called(ctx, id)

	var r0 repo.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (repo.Post, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) repo.Post); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(repo.Post)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostsByUser provides a mock function with given fields: ctx, username
func (_m *PostDB) GetPostsByUser(ctx context.Context, username string) ([]repo.Post, error) {
	ret := _m.Called(ctx, username)

	var r0 []repo.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]repo.Post, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []repo.Post); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repo.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePost provides a mock function with given fields: ctx, _a1
func (_m *PostDB) UpdatePost(ctx context.Context, _a1 repo.Post) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.Post) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPostDB interface {
	mock.TestingT
	Cleanup(func())
}

// NewPostDB creates a new instance of PostDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPostDB(t mockConstructorTestingTNewPostDB) *PostDB {
	mock := &PostDB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
