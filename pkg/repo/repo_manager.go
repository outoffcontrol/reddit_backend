package repo

type RepoManager struct {
	Users    UserDB
	Posts    PostDB
	Sessions SessionDB
}

func NewRepoManager(sessionDB SessionDB, userDB UserDB, postDB PostDB) *RepoManager {
	var r RepoManager
	r.Users = userDB
	r.Posts = postDB
	r.Sessions = sessionDB
	return &r
}
