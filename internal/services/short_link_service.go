package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
	"time"
	"url-shortener/internal/domain"
	repository "url-shortener/internal/repositories/postgres"
)

type ShortLinkService struct {
	ShortLinkRepository *repository.ShortLinkRepository
}

func NewShortLinkService(ShortLinkRepository *repository.ShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{ShortLinkRepository: ShortLinkRepository}
}

const charsetForShortCodes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(n int) (string, error) {
	buf := make([]byte, n)
	max := big.NewInt(int64(len(charsetForShortCodes)))
	for i := range buf {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		buf[i] = charsetForShortCodes[num.Int64()]
	}
	return string(buf), nil
}

const maxAttempts = 10

func (s *ShortLinkService) CreateLink(ctx context.Context, user_id int, original_url string, expires_at *time.Time) (*domain.Short_Link, error) {
	u, err := url.ParseRequestURI(original_url)
	if err != nil && u.Scheme != "" && u.Host != "" {
		return nil, NewBadRequestError("not a valid URL", errors.New("not a valid URL"))
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		shortCode, err := generateShortCode(7)
		if err != nil {
			return nil, NewInternalError("failed to generate short code", err)
		}
		link := &domain.Short_Link{
			Short_code:   shortCode,
			Original_URL: original_url,
			User_ID:      user_id,
			Expires_at:   expires_at,
		}

		err = s.ShortLinkRepository.Create(ctx, link)
		if err != nil {
			var repoErr *domain.CustomError
			if errors.As(err, &repoErr) {
				switch repoErr.Err {
				case domain.ErrDuplicateShortCode:
					continue
				case domain.ErrInvalidData:
					return nil, NewBadRequestError("invalid data", repoErr.Details)
				case domain.ErrDatabaseError:
					return nil, NewInternalError("database error", repoErr.Details)
				case domain.ErrUnexpectedError:
					return nil, NewInternalError("unexpected error", repoErr.Details)
				}
			}
		}
		return link, nil
	}
	return nil, NewNotFoundError("max attempts exceeded", nil)
}

func (s *ShortLinkService) GetByCode(ctx context.Context, short_code string) (string, error) {
	link, err := s.ShortLinkRepository.GetByCode(ctx, short_code)
	if err != nil {
		var repoErr *domain.CustomError
		if errors.As(err, &repoErr) {
			switch repoErr.Err {
			case domain.ErrLinkNotFound:
				return "", NewNotFoundError("link not found", repoErr.Details)
			case domain.ErrDatabaseError:
				return "", NewInternalError("database error", repoErr.Details)
			}
		}
	}
	return link.Original_URL, nil
}

func (s *ShortLinkService) ListByUser(ctx context.Context, userID int) ([]domain.Short_Link, error) {
	links, err := s.ShortLinkRepository.ListByUser(ctx, userID)
	if err != nil {
		return nil, NewInternalError("database error", err)
	}
	return links, nil
}
