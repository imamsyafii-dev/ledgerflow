package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
)

var ErrInvalidReconciliationRequest = errors.New(
	"invalid reconciliation request",
)

type ReconciliationRepository interface {
	List(
		ctx context.Context,
		organizationID string,
	) ([]model.Reconciliation, error)
}

type ReconciliationService struct {
	repository ReconciliationRepository
}

func NewReconciliationService(
	repository ReconciliationRepository,
) *ReconciliationService {
	return &ReconciliationService{
		repository: repository,
	}
}

func (s *ReconciliationService) List(
	ctx context.Context,
	organizationID string,
) ([]model.Reconciliation, error) {
	organizationID = strings.TrimSpace(organizationID)

	if organizationID == "" {
		return nil, fmt.Errorf(
			"%w: organization_id is required",
			ErrInvalidReconciliationRequest,
		)
	}

	return s.repository.List(ctx, organizationID)
}
