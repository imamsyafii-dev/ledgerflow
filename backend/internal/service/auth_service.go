package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/auth"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthUserRepository interface {
	FindByEmail(
		ctx context.Context,
		email string,
	) (model.User, error)
	FindByID(
		ctx context.Context,
		userID string,
	) (model.User, error)
}

type AuthService struct {
	userRepository AuthUserRepository
	jwtSecret      string
}

func NewAuthService(
	userRepository AuthUserRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (model.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return model.User{}, "", ErrInvalidCredentials
	}

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, "", ErrInvalidCredentials
		}

		return model.User{}, "", fmt.Errorf("find login user: %w", err)
	}

	if !user.IsActive {
		return model.User{}, "", ErrInvalidCredentials
	}

	if err := auth.ComparePassword(user.PasswordHash, password); err != nil {
		return model.User{}, "", ErrInvalidCredentials
	}

	token, err := auth.GenerateAccessToken(
		user.ID,
		user.OrganizationID,
		string(user.Role),
		s.jwtSecret,
	)
	if err != nil {
		return model.User{}, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetCurrentUser(
	ctx context.Context,
	userID string,
) (model.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("get current user: %w", err)
	}

	if !user.IsActive {
		return model.User{}, ErrInvalidCredentials
	}

	return user, nil
}
