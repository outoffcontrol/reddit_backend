package repo

import (
	"context"
	"database/sql"
	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/utils"
)

type UserMySql struct {
	Logger *logging.Logger
	MySql  *sql.DB
}

func (u *UserMySql) CreateUser(ctx context.Context, id string, username string, password string) (*User, error) {
	user := &User{}
	userRaw := u.MySql.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)
	err := userRaw.Scan(&user.Id, &user.Username, &user.password)
	if err == nil {
		return nil, utils.ErrUserExist
	}
	if err != nil && err != sql.ErrNoRows {
		u.Logger.Z(ctx).Error("cannot get user from mysql: ", err)
		return nil, err
	}
	user = &User{Id: id, Username: username, password: password}
	_, err = u.MySql.Exec("INSERT INTO users (`id`, `username`, `password`) VALUES (?, ?, ?)", user.Id, user.Username, user.password)
	if err != nil {
		u.Logger.Z(ctx).Error("cannot create user into mysql: ", err)
		return nil, err
	}
	return user, nil
}

func (u *UserMySql) CheckPasswordUser(ctx context.Context, username string, password string) (*User, error) {
	user := &User{}
	userRaw := u.MySql.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)
	err := userRaw.Scan(&user.Id, &user.Username, &user.password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNoUser
		}
		u.Logger.Z(ctx).Error("cannot get user from mysql: ", err)
		return nil, err
	}
	if user.password != password {
		return nil, utils.ErrBadPass
	}
	return user, nil
}
