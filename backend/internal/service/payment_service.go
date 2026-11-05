package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5"
)

var ErrInvalidPaymentRequest = errors.New("invalid payment request")
var ErrInvalidPaymentWebhook = errors.New("invalid payment webhook")
var ErrPaymentRejected = errors.New("payment rejected")

type PaymentRepository interface {
	ProcessSucceededWebhook(
		ctx context.Context,
		event model.PaymentSucceededWebhook,
		payload []byte,
	) (model.Payment, bool, error)
	List(
		ctx context.Context,
		organizationID string,
	) ([]model.Payment, error)
}

type PaymentService struct {
	repository PaymentRepository
}

func NewPaymentService(repository PaymentRepository) *PaymentService {
	return &PaymentService{
		repository: repository,
	}
}

func (s *PaymentService) ProcessSucceededWebhook(
	ctx context.Context,
	event model.PaymentSucceededWebhook,
	payload []byte,
) (model.Payment, bool, error) {
	event.EventID = strings.TrimSpace(event.EventID)
	event.EventType = strings.TrimSpace(event.EventType)
	event.ProviderPaymentID = strings.TrimSpace(event.ProviderPaymentID)
	event.OrganizationID = strings.TrimSpace(event.OrganizationID)
	event.InvoiceID = strings.TrimSpace(event.InvoiceID)
	event.Currency = strings.ToUpper(strings.TrimSpace(event.Currency))

	if event.EventID == "" {
		return model.Payment{}, false, fmt.Errorf(
			"%w: event_id is required",
			ErrInvalidPaymentWebhook,
		)
	}

	if event.EventType != "payment.succeeded" {
		return model.Payment{}, false, fmt.Errorf(
			"%w: event_type must be payment.succeeded",
			ErrInvalidPaymentWebhook,
		)
	}

	if event.ProviderPaymentID == "" {
		return model.Payment{}, false, fmt.Errorf(
			"%w: payment_id is required",
			ErrInvalidPaymentWebhook,
		)
	}

	if event.OrganizationID == "" || event.InvoiceID == "" {
		return model.Payment{}, false, fmt.Errorf(
			"%w: organization_id and invoice_id are required",
			ErrInvalidPaymentWebhook,
		)
	}

	if event.Amount <= 0 {
		return model.Payment{}, false, fmt.Errorf(
			"%w: amount must be greater than zero",
			ErrInvalidPaymentWebhook,
		)
	}

	if len(event.Currency) != 3 {
		return model.Payment{}, false, fmt.Errorf(
			"%w: currency must have three letters",
			ErrInvalidPaymentWebhook,
		)
	}

	payment, duplicate, err := s.repository.ProcessSucceededWebhook(
		ctx,
		event,
		payload,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Payment{}, false, fmt.Errorf(
				"%w: invoice must be ISSUED and match the webhook amount",
				ErrPaymentRejected,
			)
		}

		return model.Payment{}, false, err
	}

	return payment, duplicate, nil
}

func (s *PaymentService) List(
	ctx context.Context,
	organizationID string,
) ([]model.Payment, error) {
	organizationID = strings.TrimSpace(organizationID)

	if organizationID == "" {
		return nil, fmt.Errorf(
			"%w: organization_id is required",
			ErrInvalidPaymentRequest,
		)
	}

	return s.repository.List(ctx, organizationID)
}
