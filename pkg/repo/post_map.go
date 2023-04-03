package repo

import (
	"reddit_backend/pkg/utils"
	"sort"
	"sync"
)

type PostMap struct {
	Mu    *sync.RWMutex
	Posts map[string]*Post
}

func (p *PostMap) GetPostsByCategory(category string) ([]Post, error) {
	var result []Post
	p.Mu.RLock()
	for _, v := range p.Posts {
		if v.Category == category {
			result = append(result, *v)
		}
	}
	p.Mu.RUnlock()
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	return result, nil
}

func (p *PostMap) GetPostsById(id string) (Post, error) {
	p.Mu.RLock()
	post, ok := p.Posts[id]
	p.Mu.RUnlock()
	if !ok {
		return Post{}, utils.ErrInvPostId
	}
	return *post, nil
}

func (p *PostMap) DeletePostById(id string) error {
	p.Mu.Lock()
	delete(p.Posts, id)
	p.Mu.Unlock()
	return nil
}

func (p *PostMap) UpdatePost(post Post) error {
	p.Mu.Lock()
	p.Posts[post.Id] = &post
	p.Mu.Unlock()
	return nil
}

func (p *PostMap) CreatePost(post Post) error {
	p.Mu.Lock()
	p.Posts[post.Id] = &post
	p.Mu.Unlock()
	return nil
}

func (p *PostMap) GetAllPosts() ([]Post, error) {
	v := make([]Post, 0, len(p.Posts))
	p.Mu.RLock()
	for _, value := range p.Posts {
		v = append(v, *value)
	}
	p.Mu.RUnlock()
	sort.SliceStable(v, func(i, j int) bool {
		return v[i].Score > v[j].Score
	})
	return v, nil
}

func (p *PostMap) GetPostsByUser(username string) ([]Post, error) {
	result := []Post{}
	p.Mu.RLock()
	for _, v := range p.Posts {
		if v.Author.Username == username {
			result = append(result, *v)
		}

	}
	p.Mu.RUnlock()
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	return result, nil
}
