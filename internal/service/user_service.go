package service

import "context"

type userRepository interface {
	GetUserIDByName(ctx context.Context, name string) (int64, error)
	CreateUser(ctx context.Context, name string) (int64, error)
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserIDByName(ctx context.Context, name string) (int64, error) {
	return s.repo.GetUserIDByName(ctx, name)
}

func (s *UserService) CreateUser(ctx context.Context, name string) (int64, error) {
	return s.repo.CreateUser(ctx, name)
}
