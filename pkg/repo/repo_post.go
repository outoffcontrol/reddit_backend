package repo

import "context"

//go:generate mockery --name PostDB  --output ../posts --outpkg post --filename postsdb_mock.go
type PostDB interface {
	GetPostsByCategory(ctx context.Context, category string) ([]Post, error)
	GetPostsById(ctx context.Context, id string) (Post, error)
	DeletePostById(ctx context.Context, id string) error
	UpdatePost(ctx context.Context, post Post) error
	CreatePost(ctx context.Context, post Post) error
	GetAllPosts(ctx context.Context) ([]Post, error)
	GetPostsByUser(ctx context.Context, username string) ([]Post, error)
	AddComment(ctx context.Context, postId string, comment Comment) error
}

type Post struct {
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Url              string    `json:"url,omitempty"`
	Author           *User     `json:"author"`
	Category         string    `json:"category"`
	Text             string    `json:"text,omitempty"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"`
	Created          string    `json:"created"`
	UpvoteParcentage int       `json:"upvotePercentage"`
	Id               string    `json:"id"`
}

type Vote struct {
	User string `json:"user"` //user_id
	Vote int    `json:"vote"`
}

type Comment struct {
	Created string `json:"created"`
	Author  *User  `json:"author"`
	Body    string `json:"body"`
	Id      string `json:"id"`
}
