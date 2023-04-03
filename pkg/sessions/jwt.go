package sessions

import (
	"context"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/utils"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	User      *repo.User `json:"user"`
	SessionId string     `json:"session_id"`
	jwt.RegisteredClaims
}

func (s *Sessions) CreateToken(ctx context.Context, user *repo.User) (string, error) {
	expirationTime := time.Now().Add(168 * time.Hour)
	sesId := utils.GenerateId()
	claims := &JwtClaims{
		User:      user,
		SessionId: sesId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		s.Logger.Z(ctx).Error("cannot create token: ", err)
		return "", utils.ErrCreateToken
	}
	err = s.Repo.AddSession(ctx, user.Id, claims.SessionId)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Sessions) IsValid(ctx context.Context, token string) (*repo.User, error) {
	tok, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		s.Logger.Z(ctx).Error(err)
	}
	if claims, ok := tok.Claims.(*JwtClaims); ok && tok.Valid {
		v, err := s.Repo.CheckSession(ctx, claims.User.Id, claims.SessionId)
		if err != nil {
			return nil, err
		}
		if v {
			return claims.User, nil
		}
	}

	return nil, utils.ErrUnauthorized
}
