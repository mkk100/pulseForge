package service

import (
	"context"
	"time"

	"pulseforge/internal/repo"
)

type postRepository interface {
	CreatePost(ctx context.Context, post repo.Post) (int64, error)
	ListRecentPosts(ctx context.Context, limit int) ([]repo.Post, error)
}

type CreatePostInput struct {
	Title       string
	Description string
	UserID      int64
}

type Post struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int64     `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PostService struct {
	repo postRepository
}

func NewPostService(repo postRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, input CreatePostInput) (int64, error) {
	return s.repo.CreatePost(ctx, repo.Post{
		Title:       input.Title,
		Description: input.Description,
		UserID:      input.UserID,
	})
}

func (s *PostService) ListRecentPosts(ctx context.Context, limit int) ([]Post, error) {
	posts, err := s.repo.ListRecentPosts(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, Post{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			UserID:      post.UserID,
			CreatedAt:   post.CreatedAt,
		})
	}
	return result, nil
}
