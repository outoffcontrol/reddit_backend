package sessions

import (
	"context"
	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/utils"
)

type Sessions struct {
	Secret string
	Logger *logging.Logger
	Repo   repo.SessionDB
}

type sessKey string

var SessionKey sessKey = "user"

func (s *Sessions) GetUserFromContext(ctx context.Context) (*repo.User, error) {
	user, ok := ctx.Value(SessionKey).(*repo.User)
	if !ok || user == nil {
		return nil, utils.ErrUnauthorized
	}

	return ctx.Value(SessionKey).(*repo.User), nil
}
