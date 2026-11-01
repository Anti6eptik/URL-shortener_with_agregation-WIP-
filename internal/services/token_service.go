package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"goshop/internal/config"
	"goshop/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Cfg *config.Config
}

func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{Cfg: cfg}
}

type UserClaims struct {
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	TokenID string `json:"jti"`
	jwt.RegisteredClaims
}

func (s *TokenService) generateJTI() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *TokenService) NewUserClaims(userID int, email string, role string, expiresIn time.Duration) *UserClaims {
	jti := s.generateJTI()
	claims := UserClaims{
		UserID:  userID,
		Email:   email,
		Role:    role,
		TokenID: jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Issuer:    "goshop",
		},
	}
	return &claims
}

func (s *TokenService) GenerateAccessToken(user *domain.User) (string, error) {
	claims := s.NewUserClaims(user.ID, user.Email, user.Role, s.Cfg.AccessTokenExpiry)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.Cfg.JWTSecret)
}

func (s *TokenService) ValidateAccessToken(tokenSting string) (*UserClaims, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenSting, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexoected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuedAt(), jwt.WithIssuer("goshop"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
