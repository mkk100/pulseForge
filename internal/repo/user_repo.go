package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
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
