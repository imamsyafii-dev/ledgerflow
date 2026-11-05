package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
)

var ErrInvalidOrganization = errors.New("invalid organization")

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type OrganizationRepository interface {
	Create(
		ctx context.Context,
		name string,
		slug string,
	) (model.Organization, error)
}

type OrganizationService struct {
	repository OrganizationRepository
}

func NewOrganizationService(repository OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		repository: repository,
	}
}

func (s *OrganizationService) Create(
	ctx context.Context,
	name string,
	slug string,
) (model.Organization, error) {
	name = strings.TrimSpace(name)
	slug = strings.ToLower(strings.TrimSpace(slug))

	if name == "" {
		return model.Organization{}, fmt.Errorf("%w: name is required", ErrInvalidOrganization)
	}

	if !slugPattern.MatchString(slug) {
		return model.Organization{}, fmt.Errorf(
			"%w: slug must contain lowercase letters, numbers, and hyphens only",
			ErrInvalidOrganization,
		)
	}

	return s.repository.Create(ctx, name, slug)
}
