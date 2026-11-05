package service

import (
	"context"
	"errors"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidRefund       = errors.New("invalid refund request")
	ErrRefundNotFound      = errors.New("refund request not found")
	ErrRefundNotEligible   = errors.New("payment is not eligible for refund")
	ErrPendingRefundExists = errors.New(
		"payment already has a pending refund request",
	)
)

type RefundRepository interface {
	Create(
		ctx context.Context,
		organizationID string,
		paymentID string,
		requestedBy string,
		reason string,
	) (model.RefundRequest, error)

	List(
		ctx context.Context,
		organizationID string,
	) ([]model.RefundRequest, error)

	Decide(
		ctx context.Context,
		organizationID string,
		refundID string,
		actorID string,
		status model.RefundStatus,
		note string,
	) (model.RefundRequest, error)
}

type RefundService struct {
	repository RefundRepository
}

func NewRefundService(repository RefundRepository) *RefundService {
	return &RefundService{
		repository: repository,
	}
}

func (s *RefundService) Create(
	ctx context.Context,
	organizationID string,
	paymentID string,
	requestedBy string,
	reason string,
) (model.RefundRequest, error) {
	reason = strings.TrimSpace(reason)

	if organizationID == "" ||
		paymentID == "" ||
		requestedBy == "" ||
		reason == "" {
		return model.RefundRequest{}, ErrInvalidRefund
	}

	refund, err := s.repository.Create(
		ctx,
		organizationID,
		paymentID,
		requestedBy,
		reason,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.RefundRequest{}, ErrRefundNotEligible
	}

	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return model.RefundRequest{}, ErrPendingRefundExists
	}

	return refund, err
}

func (s *RefundService) List(
	ctx context.Context,
	organizationID string,
) ([]model.RefundRequest, error) {
	if strings.TrimSpace(organizationID) == "" {
		return nil, ErrInvalidRefund
	}

	return s.repository.List(ctx, organizationID)
}

func (s *RefundService) Decide(
	ctx context.Context,
	organizationID string,
	refundID string,
	actorID string,
	status model.RefundStatus,
	note string,
) (model.RefundRequest, error