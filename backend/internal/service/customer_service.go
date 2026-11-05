package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
)

var ErrInvalidCustomer = errors.New("invalid customer")

type CustomerRepository interface {
	Create(
		ctx context.Context,
		organizationID string,
		name string,
		email string,
	) (model.Customer, error)
	List(
		ctx context.Context,
		organizationID string,
	) ([]model.Customer, error)
}

type CustomerService struct {
	repository CustomerRepository
}

func NewCustomerService(repository CustomerRepository) *CustomerService {
	return &CustomerService{
		repository: repository,
	}
}

func (s *CustomerService) Create(
	ctx context.Context,
	organizationID string,
	name string,
	email string,
) (model.Customer, error) {
	organizationID = strings.TrimSpace(organizationID)
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if organizationID == "" {
		return model.Customer{}, fmt.Errorf(
			"%w: organization_id is required",
			ErrInvalidCustomer,
		)
	}

	if name == "" {
		return model.Customer{}, fmt.Errorf(
			"%w: name is required",
			ErrInvalidCustomer,
		)
	}

	return s.repository.Create(ctx, organizationID, name, email)
}

func (s *CustomerService) List(
	ctx context.Context,
	organizationID string,
) ([]model.Customer, error) {
	organizationID = strings.TrimSpace(organizationID)

	if organizationID == "" {
		return nil, fmt.Errorf(
			"%w: organization_id is required",
			ErrInvalidCustomer,
		)
	}

	return s.repository.List(ctx, organizationID)
}
