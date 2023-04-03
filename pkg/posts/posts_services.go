package post

import (
	"context"
	"reddit_backend/pkg/repo"
	"reddit_backend/pkg/utils"
	"time"
)

func (p *Posts) createPost(ctx context.Context, post Post, user *repo.User) (repo.Post, error) {
	result := repo.Post{
		Score:    1,
		Views:    0,
		Type:     post.Type,
		Title:    post.Title,
		Url:      post.Url,
		Author:   user,
		Category: post.Category,
		Text:     post.Text,
		Votes: []repo.Vote{
			{
				User: user.Id,
				Vote: 1,
			},
		},
		Comments:         []repo.Comment{},
		Created:          time.Now().Format("2006-01-02T15:04:05Z07:00"),
		UpvoteParcentage: 100,
		Id:               utils.GenerateId(),
	}
	err := p.Repo.CreatePost(ctx, result)
	if err != nil {
		return repo.Post{}, err
	}

	return result, nil
}

func (p *Posts) addCommentToPost(ctx context.Context, postId string, user *repo.User, commentMsg string) (repo.Post, error) {
	post, err := p.Repo.GetPostsById(ctx, postId)
	if err != nil {
		return repo.Post{}, utils.ErrNoPost
	}
	comments := make([]repo.Comment, len(post.Comments))
	copy(comments, post.Comments)
	comment := repo.Comment{
		Created: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Author:  user,
		Body:    commentMsg,
		Id:      utils.GenerateId(),
	}
	comments = append(comments, comment)
	post.Comments = comments
	err = p.Repo.AddComment(ctx, post.Id, comment)
	if err != nil {
		return repo.Post{}, err
	}
	return post, nil
}

func (p *Posts) deleteComment(ctx context.Context, postId string, commentId string) (repo.Post, error) {
	post, err := p.Repo.GetPostsById(ctx, postId)
	if err != nil {
		return repo.Post{}, utils.ErrNoPost
	}
	comments := make([]repo.Comment, len(post.Comments))
	copy(comments, post.Comments)
	for i, v := range comments {
		if v.Id == commentId {
			comments[i] = comments[len(comments)-1]
			comments[len(comments)-1] = repo.Comment{}
			comments = comments[:len(comments)-1]
		}
	}
	post.Comments = comments
	err = p.Repo.UpdatePost(ctx, post)
	if err != nil {
		return repo.Post{}, err
	}
	return post, nil
}

func (p *Posts) changeVote(ctx context.Context, post repo.Post, user *repo.User, way int) (repo.Post, error) {
	voteExist := false
	countUpvote := 0
	votes := make([]repo.Vote, len(post.Votes))
	copy(votes, post.Votes)
	if way == 0 {
		for i, v := range votes {
			if v.User == user.Id {
				votes[i] = votes[len(votes)-1]
				votes[len(votes)-1] = repo.Vote{}
				votes = votes[:len(votes)-1]
			}
		}
	}
	score := 0
	for i, v := range votes {
		if v.User == user.Id {
			votes[i].Vote = way
			voteExist = true
		}
		if votes[i].Vote == 1 {
			countUpvote = countUpvote + 1
		}
		score = score + votes[i].Vote
	}
	if !voteExist && way != 0 {
		votes = append(votes, repo.Vote{User: user.Id, Vote: way})
		if way == 1 {
			countUpvote = countUpvote + 1
		}
		score = score + way
	}
	post.Score = score
	if len(votes) > 0 {
		post.UpvoteParcentage = 100 / len(votes) * countUpvote
	} else {
		post.UpvoteParcentage = 0
	}
	post.Votes = votes
	err := p.Repo.UpdatePost(ctx, post)
	if err != nil {
		return repo.Post{}, err
	}
	return post, nil
}
