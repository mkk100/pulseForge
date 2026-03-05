package db
import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


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

func NewPostRepo(pool *pgxpool.Pool) *PostRepo{
	return &PostRepo{pool: pool}
}

func (r *PostRepo) CreatePost(ctx context.Context, post Post) (int64, error){
	var postID int64
	err := r.pool.QueryRow(ctx, `
    	INSERT INTO posts (postTitle, postDescription, userId)
    	VALUES ($1, $2, $3)
    	RETURNING postId
	`, post.Title, post.Description, post.UserID).Scan(&postID)

	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}
	return postID, nil
}