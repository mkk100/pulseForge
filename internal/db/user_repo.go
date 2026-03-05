// internal/db/user_repo.go
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

type PostRepo struct {
	pool *pgxpool.Pool
}

type Post struct {
    ID          int64
    Title       string
    Description string
    UserID      int64
    CreatedAt   time.Time
}


func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool} // stores the UserRepo reference to the existing pool
}

func NewPostRepo(pool *pgxpool.Pool) *PostRepo{
	return &PostRepo{pool: pool}
}

func (r *UserRepo) GetUserIDByName(ctx context.Context, name string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		`SELECT userid FROM users WHERE username=$1`,
		name,
	).Scan(&id)
	return id, err
}

func (r *UserRepo) CreateUser(ctx context.Context, name string) (int64, error) {
	var userID int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (userName)
		VALUES ($1)
		RETURNING userId
	`, name).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return userID, nil
}

func (r *PostRepo) CreatePost(ctx context.Context, post Post) (int64, error){
	var postID int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO posts (postTitle, postDescription)
	`)
}