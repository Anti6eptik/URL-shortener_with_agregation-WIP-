package postgres

import (
	"context"
	"errors"
	"fmt"
	"goshop/internal/domain"
	"goshop/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.UserRepository = (*userRepo)(nil)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepo {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query, user.Email, user.Password_hash).Scan(user.ID, user.Created_at)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (int, error) {
	query := `SELECT id FROM users WHERE email = $1`

	var id int

	err := r.pool.QueryRow(ctx, query, email).Scan(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, nil
		}
		return -1, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}

	return id, nil
}
