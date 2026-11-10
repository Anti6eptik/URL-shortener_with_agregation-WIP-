package repository

import (
	"context"
	"errors"
	"strings"
	"url-shortener/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, role
	`

	err := r.pool.QueryRow(ctx, query, user.Email, user.Password_hash).Scan(&user.ID, &user.Created_at, &user.Role)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": //unique violaction
				if strings.Contains(pgErr.Message, "email") {
					return domain.NewCustomError("create", domain.ErrDuplicateEmail, pgErr)
				}
				return domain.NewCustomError("create", domain.ErrInvalidData, pgErr)
			default:
				return domain.NewCustomError("create", domain.ErrDatabaseError, pgErr)
			}
		}
		return domain.NewCustomError("create", domain.ErrUnexpectedError, err)
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at, role FROM users WHERE email = $1`

	var user domain.User

	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password_hash, &user.Created_at, &user.Updated_at, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewCustomError("get", domain.ErrUserNotFound, nil)
		}
		return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at, role FROM users WHERE id = $1`

	var user domain.User

	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password_hash, &user.Created_at, &user.Updated_at, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewCustomError("get", domain.ErrUserNotFound, nil)
		}
		return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateByEmail(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET password_hash = $1, updated_at = NOW()
		WHERE email = $2
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query, user.Password_hash, user.Email).Scan(&user.Updated_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.NewCustomError("put", domain.ErrUserNotFound, nil)
		}
		return domain.NewCustomError("put", domain.ErrDatabaseError, err)
	}
	return nil
}

func (r *UserRepository) GetAll(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	query := `
	SELECT id, email, created_at, updated_at, role
	FROM users
	ORDER BY id
	LIMIT $1
	OFFSET $2`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Email, &user.Created_at, &user.Updated_at, &user.Role)
		if err != nil {
			return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
		}
		users = append(users, user)
		if err = rows.Err(); err != nil {
			return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
		}
	}
	return users, nil
}

func (r *UserRepository) GetTotalCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return -1, domain.NewCustomError("get", domain.ErrDatabaseError, err)
	}
	return count, nil
}

func (r *UserRepository) DeleteById(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": //foreign key violation
				return domain.NewCustomError("delete", domain.ErrDeleteReferenced, pgErr)
			default:
				return domain.NewCustomError("delete", domain.ErrDatabaseError, pgErr)
			}
		}
		return domain.NewCustomError("delete", domain.ErrUnexpectedError, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.NewCustomError("delete", domain.ErrUserNotFound, nil)
	}
	return nil
}
