package utils

import (
	"errors"

	"github.com/google/uuid"
)

func GenerateId() string {
	id := uuid.New()
	return (id.String())
}

type Return422Errors struct {
	Errors []Return422Error `json:"errors"`
}

type Return422Error struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

type ReturnMsg struct {
	Msg string `json:"message"`
}

type ReturnToken struct {
	Token string `json:"token"`
}

var (
	ErrNoUser       = errors.New("user not found")
	ErrUserExist    = errors.New("username already exists")
	ErrBadPass      = errors.New("invalid password")
	ErrNoPost       = errors.New("post not found")
	ErrInvPostId    = errors.New("invalid post id")
	ErrUnauthorized = errors.New("unauthorized")
	ErrCreateToken  = errors.New("error create token")
)
