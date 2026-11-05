package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
)

var ErrInvalidLedgerRequest = errors.New("invalid ledger request")

type LedgerRepository interface {
	ListJournalEntries(
		ctx context.Context,
		organizationID string,
	) ([]model.JournalEntry, error)
}

type LedgerService struct {
	repository LedgerRepository
}

func NewLedgerService(repository LedgerRepository) *LedgerService {
	return &LedgerService{
		repository: repository,
	}
}

func (s *LedgerService) ListJournalEntries(
	ctx context.Context,
	organizationID string,
) ([]model.JournalEntry, error) {
	organizationID = strings.TrimSpace(organizationID)

	if organizationID == "" {
		return nil, fmt.Errorf(
			"%w: organization_id is required",
			ErrInvalidLedgerRequest,
		)
	}

	return s.repository.ListJournalEntries(ctx, organizationID)
}
