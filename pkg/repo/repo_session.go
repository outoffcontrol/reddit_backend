package repo

import "context"

type SessionDB interface {
	AddSession(ctx context.Context, userId string, token string) error
	CheckSession(ctx context.Context, userId string, token string) (bool, error)
}

type Sessions struct {
	Sessions []string
}
