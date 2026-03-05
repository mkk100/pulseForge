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

func (r *PostRepo) ListRecentPosts(ctx context.Context, limit string) ([]Post, error){
	var posts []Post // used pool.Query for multiple rows, important learning lesson function here
	rows, err := r.pool.Query(ctx, `
		SELECT postId, postTitle, postDescription, userId, created_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)

	for rows.Next(){
		var post Post
		if err:= rows.Scan(
			&post.ID,
			&post.Title,
			&post.Description,
			&post.UserID,
			&post.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}
	return posts, err
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