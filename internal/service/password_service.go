package service

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	cost int
}

func NewPasswordService(cost int) *PasswordService {
	return &PasswordService{cost: cost}
}

func (s *PasswordService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain uppercase, lowercase and number")
	}
	return nil
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	if err := s.ValidatePassword(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	return string(hashedBytes), nil
}

func (s *PasswordService) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
