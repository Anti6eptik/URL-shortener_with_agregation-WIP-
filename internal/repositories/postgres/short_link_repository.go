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
		RETURNING id, created_at`

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

func (r *ShortLinkRepository) GetByCode(ctx context.Context, shortcode string) (*domain.Short_Link, error) {
	query := `SELECT id, short_code, original_url, user_id, created_at, expires_at FROM short_links WHERE short_code = $1`

	var link domain.Short_Link

	err := r.pool.QueryRow(ctx, query, shortcode).Scan(&link.ID, &link.Short_code, &link.Original_URL, &link.User_ID, &link.Created_at, &link.Expires_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewCustomError("get", domain.ErrLinkNotFound, nil)
		}
		return nil, domain.NewCustomError("get", domain.ErrUnexpectedError, err)
	}

	return &link, nil
}

func (r *ShortLinkRepository) ListByUser(ctx context.Context, userID int) ([]domain.Short_Link, error) {
	query := `
	SELECT id, short_code, original_url, user_id, created_at, expires_at
	FROM short_links
	WHERE user_id = $1
	ORDER BY created_at`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
	}
	defer rows.Close()

	links := make([]domain.Short_Link, 0)
	for rows.Next() {
		var link domain.Short_Link
		err := rows.Scan(&link.ID, &link.Short_code, &link.Original_URL, &link.User_ID, &link.Created_at, &link.Expires_at)
		if err != nil {
			return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
		}
		links = append(links, link)
		if err = rows.Err(); err != nil {
			return nil, domain.NewCustomError("get", domain.ErrDatabaseError, err)
		}
	}
	return links, nil
}
