package middleware

import (
	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/sessions"
)

type Middleware struct {
	Session *sessions.Sessions
	Logger  *logging.Logger
	Repo    *repo.RepoManager
}
