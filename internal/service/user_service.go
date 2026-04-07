package service

import (
	"context"
	"errors"
	"goshop/internal/domain"
	"goshop/internal/repository/postgres"
)

type UserService struct {
	UserRepository  *postgres.UserRepository
	PasswordService *PasswordService
}

func NewUserService(
	UserRepository *postgres.UserRepository,
	PasswordService *PasswordService,
) *UserService {
	return &UserService{
		UserRepository:  UserRepository,
		PasswordService: PasswordService,
	}
}

func (s *UserService) Login(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := s.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		var repoErr *domain.UserError
		if errors.As(err, &repoErr) && errors.Is(repoErr.Err, domain.ErrUserNotFound) {
			return nil, NewUnauthorizedError("invalid credentials", nil)
		}
		return nil, NewInternalError("failed to find user", err)
	}
	if !s.PasswordService.CheckPassword(password, user.Password_hash) {
		return nil, NewUnauthorizedError("invalid credentials", nil)
	}

	//создание access и refresh токенов

	return user, nil
}

func (s *UserService) Register(ctx context.Context, email string, password string) (*domain.User, error) {
	if email == "" {
		return nil, NewBadRequestError("email is required", errors.New("no email"))
	}

	err := s.PasswordService.ValidatePassword(password)
	if err != nil {
		return nil, NewBadRequestError("password is too weak", err)
	}

	hashedPassword, err := s.PasswordService.HashPassword(password)
	if err != nil {
		return nil, NewInternalError("", err)
	}

	user := &domain.User{
		Email:         email,
		Password_hash: hashedPassword,
	}

	err = s.UserRepository.Create(ctx, user)
	if err != nil {
		var repoErr *domain.UserError
		if errors.As(err, &repoErr) {
			switch repoErr.Err {
			case domain.ErrDuplicateEmail:
				return nil, NewConflictError("email already exists", repoErr.Details)
			case domain.ErrInvalidUserData:
				return nil, NewBadRequestError("invalid user data", repoErr.Details)
			case domain.ErrDatabaseError:
				return nil, NewInternalError("database error", repoErr.Details)
			case domain.ErrUnexpectedError:
				return nil, NewInternalError("unexpected error", repoErr.Details)
			}
		}
	}

	//создание access и refresh токенов

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, email string, old_password string, new_password string) error {
	user, err := s.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		var repoErr *domain.UserError
		if errors.As(err, &repoErr) && errors.Is(repoErr.Err, domain.ErrUserNotFound) {
			return NewUnauthorizedError("invalid credentials", nil)
		}
		return NewInternalError("failed to find user", err)
	}

	if !s.PasswordService.CheckPassword(old_password, user.Password_hash) {
		return NewUnauthorizedError("invalid credentials", nil)
	}

	err = s.PasswordService.ValidatePassword(new_password)
	if err != nil {
		return NewBadRequestError("new password is too weak", err)
	}

	hashedPassword, err := s.PasswordService.HashPassword(new_password)
	if err != nil {
		return NewInternalError("", err)
	}
	user.Password_hash = hashedPassword

	err = s.UserRepository.UpdateByEmail(ctx, user)
	if err != nil {
		var repoErr *domain.UserError
		if errors.As(err, &repoErr) {
			return NewInternalError("database error", repoErr.Details)
		}
		return NewInternalError("unexpected error", repoErr.Details)
	}

	//создание access и refresh токенов

	return nil
}

func (s *UserService) GetAll(ctx context.Context, page int, limit int) ([]domain.User, int, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, err := s.UserRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, 0, NewInternalError("database error", err)
	}

	total, err := s.UserRepository.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, 0, NewInternalError("database error", err)
	}

	totalPages := (total + limit - 1) / limit

	return users, total, totalPages, nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	err := s.UserRepository.DeleteById(ctx, id)
	if err != nil {
		var repoErr *domain.UserError
		if errors.As(err, &repoErr) {
			switch repoErr.Err {
			case domain.ErrDeleteReferenced:
				return NewConflictError("attempt to delete referenced", repoErr.Details)
			case domain.ErrDatabaseError:
				return NewInternalError("database error", repoErr.Details)
			case domain.ErrUserNotFound:
				return NewNotFoundError("user not found", repoErr.Details)
			case domain.ErrUnexpectedError:
				return NewInternalError("unexpected error", repoErr.Details)
			}
		}
	}
	return nil
}
