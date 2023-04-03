package repo

import (
	"sync"
)

type SessionsMap struct {
	Mu       *sync.RWMutex
	Sessions map[string]*Sessions
}

func (s *SessionsMap) AddSession(userId string, token string) error {
	s.Mu.RLock()
	_, ok := s.Sessions[userId]
	s.Mu.RUnlock()
	if !ok {
		s.Mu.Lock()
		s.Sessions[userId] = &Sessions{Sessions: []string{token}}
		s.Mu.Unlock()
		return nil
	}
	s.Mu.Lock()
	s.Sessions[userId].Sessions = append(s.Sessions[userId].Sessions, token)
	s.Mu.Unlock()
	return nil
}

func (s *SessionsMap) CheckSession(userId string, token string) (bool, error) {
	s.Mu.RLock()
	usersSessions, ok := s.Sessions[userId]
	s.Mu.RUnlock()
	if ok {
		for _, t := range usersSessions.Sessions {
			if t == token {
				return true, nil
			}
		}
	}
	return false, nil
}
