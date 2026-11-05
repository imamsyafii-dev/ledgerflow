package service

import (
	"context"
	"errors"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
)

var ErrInvalidUserOrganization = errors.New("invalid user organization")

type UserRepository interface {
	ListByOrganization(
		ctx context.Context,
		organizationID string,
	) ([]model.User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) List(
	ctx context.Context,
	organizationID string,
) ([]model.User, error) {
	organizationID = strings.TrimSpace(organizationID)
	if organizationID == "" {
		return nil, ErrInvalidUserOrganization
	}

	return s.repository.ListByOrganization(ctx, organizationID)
}
