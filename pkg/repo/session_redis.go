package repo

import (
	"context"
	"reddit_backend/pkg/logging"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

type SessionsRedis struct {
	Logger    *logging.Logger
	RedisConn redis.Conn
}

func (s *SessionsRedis) AddSession(ctx context.Context, userId string, session_id string) error {
	err := s.deleteExpSessions(ctx, userId)
	if err != nil {
		s.Logger.Z(ctx).Error("cant delete exp sessins from Redis: ", err)
		//return  в целом это не должно влиять на возможность залогиниться
	}
	sesOpts := []string{time.Now().Add(72 * time.Hour).Format("2006-01-02T15:04:05Z07:00")}
	_, err = s.RedisConn.Do("HSET", userId, session_id, sesOpts)
	if err != nil {
		s.Logger.Z(ctx).Error("cant hset session_id to Redis: ", err)
		return err
	}
	return nil
}

func (s *SessionsRedis) CheckSession(ctx context.Context, userId string, session_id string) (bool, error) {
	_, err := redis.String(s.RedisConn.Do("HGET", userId, session_id))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			return false, nil
		}
		s.Logger.Z(ctx).Error("cant get data from Redis: ", err)
		return false, err
	}

	return true, nil
}

func (s *SessionsRedis) deleteExpSessions(ctx context.Context, userId string) error {
	v, err := redis.StringMap(s.RedisConn.Do("HGETALL", userId))
	if err != nil {
		s.Logger.Z(ctx).Error("cant get data from Redis: ", err)
		return err
	}

	for i, v := range v {
		arr := strings.Split(v[1:len(v)-1], " ")
		t, err := time.Parse("2006-01-02T15:04:05Z07:00", arr[0])
		if err != nil {
			s.Logger.Z(ctx).Error("cant convert time: ", err)
			continue
		}

		if t.Before(time.Now()) {
			_, err := s.RedisConn.Do("HDEL", userId, i)
			if err != nil {
				s.Logger.Z(ctx).Error("cant delete session from Redis: ", err)
				continue
			}
		}
	}
	return nil
}
