package repo

import (
	"context"
	"reddit_backend/pkg/logging"
	"reddit_backend/pkg/utils"
	"sort"

	"gopkg.in/mgo.v2/bson"
)

type PostMongo struct {
	Logger  *logging.Logger
	Mongo   IMongoCollection
	Context *context.Context
}

func (p *PostMongo) GetPostsByCategory(ctx context.Context, category string) ([]Post, error) {
	var result []Post
	cur, err := p.Mongo.Find(*p.Context, bson.M{"category": category})
	if err != nil {
		p.Logger.Z(ctx).Error("cannot get posts by category from mongo: ", err)
		return nil, err
	}

	err = cur.All(*p.Context, &result)
	if err != nil {
		p.Logger.Z(ctx).Error("cannot decode posts by category: ", err)
		return nil, err
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	return result, nil
}

func (p *PostMongo) GetPostsById(ctx context.Context, id string) (Post, error) {
	post := &Post{}
	err := p.Mongo.FindOne(*p.Context, bson.M{"id": id}).Decode(&post)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return Post{}, utils.ErrInvPostId
		}
		p.Logger.Z(ctx).Error("cannot get post by id: ", err)
		return Post{}, err
	}
	return *post, nil
}

func (p *PostMongo) DeletePostById(ctx context.Context, id string) error {
	_, err := p.Mongo.DeleteOne(*p.Context, bson.M{"id": id})
	if err != nil {
		p.Logger.Z(ctx).Error("cannot delete post in mongo: ", err)
		return err
	}
	return nil
}

func (p *PostMongo) UpdatePost(ctx context.Context, post Post) error {
	_, err := p.Mongo.UpdateOne(*p.Context, bson.M{"id": post.Id}, bson.M{
		"$set": post,
	})
	if err != nil {
		p.Logger.Z(ctx).Error("cannot update post in mongo: ", err)
		return err
	}
	return nil
}

func (p *PostMongo) AddComment(ctx context.Context, postId string, comment Comment) error {
	_, err := p.Mongo.UpdateOne(*p.Context, bson.M{"id": postId}, bson.M{
		"$push": bson.M{
			"comments": comment,
		},
	})
	if err != nil {
		p.Logger.Z(ctx).Error("cannot update post in mongo: ", err)
		return err
	}
	return nil
}

func (p *PostMongo) CreatePost(ctx context.Context, post Post) error {
	_, err := p.Mongo.InsertOne(*p.Context, &post)
	if err != nil {
		p.Logger.Z(ctx).Error("cannot add post to mongo: ", err)
		return err
	}
	return nil
}

func (p *PostMongo) GetAllPosts(ctx context.Context) ([]Post, error) {
	posts := []Post{}
	cur, err := p.Mongo.Find(*p.Context, bson.M{})
	if err != nil {
		p.Logger.Z(ctx).Error("cannot get all posts from mongo: ", err)
		return nil, err
	}
	err = cur.All(*p.Context, &posts)
	if err != nil {
		p.Logger.Z(ctx).Error("cannot decode all posts: ", err)
		return nil, err
	}
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].Score > posts[j].Score
	})
	return posts, nil
}

func (p *PostMongo) GetPostsByUser(ctx context.Context, username string) ([]Post, error) {
	result := []Post{}
	cur, err := p.Mongo.Find(*p.Context, bson.M{"author.username": username})
	if err != nil {
		p.Logger.Z(ctx).Error("cannot get users posts from mongo: ", err)
		return nil, err
	}

	err = cur.All(*p.Context, &result)
	if err != nil {
		p.Logger.Z(ctx).Error("cannot decode users posts: ", err)
		return nil, err
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	return result, nil
}
