package repo

import (
	"reddit_backend/pkg/utils"
	"sync"
)

type UserMap struct {
	Mu    *sync.RWMutex
	Users map[string]*User
}

func (u *UserMap) CreateUser(id string, username string, password string) (*User, error) {
	u.Mu.RLock()
	_, ok := u.Users[username]
	u.Mu.RUnlock()
	if ok {
		return nil, utils.ErrUserExist
	}
	u.Mu.Lock()
	u.Users[username] = &User{Id: id, Username: username, password: password}
	u.Mu.Unlock()
	return u.Users[username], nil
}

func (u *UserMap) CheckPasswordUser(username string, password string) (*User, error) {
	u.Mu.RLock()
	user, ok := u.Users[username]
	u.Mu.RUnlock()
	if !ok {
		return nil, utils.ErrNoUser
	}
	if user.password != password {
		return nil, utils.ErrBadPass
	}
	return user, nil
}
