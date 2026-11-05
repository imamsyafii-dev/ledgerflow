package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
)

var ErrInvalidInvoice = errors.New("invalid invoice")

type InvoiceRepository interface {
	Create(
		ctx context.Context,
		organizationID string,
		customerID string,
		number string,
		amount int64,
		currency string,
		dueDate time.Time,
	) (model.Invoice, error)
}

type InvoiceService struct {
	repository InvoiceRepository
}

func NewInvoiceService(repository InvoiceRepository) *InvoiceService {
	return &InvoiceService{
		repository: repository,
	}
}

func (s *InvoiceService) Create(
	ctx context.Context,
	organizationID string,
	customerID string,
	number string,
	amount int64,
	currency string,
	dueDate string,
) (model.Invoice, error) {
	organizationID = strings.TrimSpace(organizationID)
	customerID = strings.TrimSpace(customerID)
	number = strings.TrimSpace(number)
	currency = strings.ToUpper(strings.TrimSpace(currency))

	if organizationID == "" {
		return model.Invoice{}, fmt.Errorf("%w: organization_id is required", ErrInvalidInvoice)
	}

	if customerID == "" {
		return model.Invoice{}, fmt.Errorf("%w: customer_id is required", ErrInvalidInvoice)
	}

	if number == "" {
		return model.Invoice{}, fmt.Errorf("%w: number is required", ErrInvalidInvoice)
	}

	if amount <= 0 {
		return model.Invoice{}, fmt.Errorf("%w: amount must be greater than zero", ErrInvalidInvoice)
	}

	if len(currency) != 3 {
		return model.Invoice{}, fmt.Errorf("%w: currency must have three letters", ErrInvalidInvoice)
	}

	parsedDueDate, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return model.Invoice{}, fmt.Errorf(
			"%w: due_date must use YYYY-MM-DD format",
			ErrInvalidInvoice,
		)
	}

	return s.repository.Create(
		ctx,
		organizationID,
		customerID,
		number,
		amount,
		currency,
		parsedDueDate,
	)
}
