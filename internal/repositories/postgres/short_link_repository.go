package repository

import (
	"context"
	"errors"
	"strings"
	"url-shortener/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortLinkRepository struct {
	pool *pgxpool.Pool
}

func NewShortLinkRepository(pool *pgxpool.Pool) *ShortLinkRepository {
	return &ShortLinkRepository{pool: pool}
}

func (r *ShortLinkRepository) Create(ctx context.Context, shortlink *domain.Short_Link) error {
	query := `
		INSERT INTO short_links (short_code, original_url, user_id, created_at, expires_at)
		VALUES ($1, $2, $3, NOW(), $4)
		RETURNING id, created_at 
	`

	err := r.pool.QueryRow(ctx, query, shortlink.Short_code, shortlink.Original_URL, shortlink.User_ID, shortlink.Expires_at).Scan(&shortlink.ID, &shortlink.Created_at)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": //unique violation
				if strings.Contains(pgErr.Message, "short_code") {
					return domain.NewCustomError("create", domain.ErrDuplicateShortCode, pgErr)
				}
				return domain.NewCustomError("create", domain.ErrInvalidData, pgErr)
			case "23503": //foreign key violation
				return domain.NewCustomError("create", domain.ErrInvalidData, pgErr)
			default:
				return domain.NewCustomError("create", domain.ErrDatabaseError, pgErr)
			}
		}
		return domain.NewCustomError("create", domain.ErrUnexpectedError, err)
	}
	return nil
}
