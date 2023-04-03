package repo

import "context"

type UserDB interface {
	CreateUser(ctx context.Context, id string, username string, password string) (*User, error)
	CheckPasswordUser(ctx context.Context, username string, password string) (*User, error)
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	password string
}
