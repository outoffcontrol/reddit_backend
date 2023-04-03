package users

import (
	"context"
	"reddit_backend/pkg/utils"
)

func (u *Users) createUser(ctx context.Context, userReq *UserReq) (string, error) {
	user, err := u.Repo.CreateUser(ctx, utils.GenerateId(), userReq.Username, userReq.Password)
	if err != nil {
		return "", err
	}
	token, err := u.Session.CreateToken(ctx, user)
	if err != nil {
		return "", err
	}
	return token, nil
}
