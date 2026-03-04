// internal/db/user_repo.go
package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool} // stores the UserRepo reference to the existing pool
}

func (r *UserRepo) GetUserIDByName(ctx context.Context, name string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		`SELECT userid FROM users WHERE username=$1`,
		name,
	).Scan(&id)
	return id, err
}
