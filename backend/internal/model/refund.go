package model

import "time"

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "PENDING"
	RefundStatusApproved  RefundStatus = "APPROVED"
	RefundStatusRejected  RefundStatus = "REJECTED"
	RefundStatusCompleted RefundStatus = "COMPLETED"
	RefundStatusFailed    RefundStatus = "FAILED"
)

type RefundRequest struct {
	ID               string       `json:"id"`
	OrganizationID   string       `json:"organization_id"`
	PaymentID        string       `json:"payment_id"`
	InvoiceNumber    string       `json:"invoice_number"`
	CustomerName     string       `json:"customer_name"`
	Amount           int64        `json:"amount"`
	Currency         string       `json:"currency"`
	RequestedBy      string       `json:"requested_by"`
	RequestedByName  string       `json:"requested_by_name"`
	Reason           string       `json:"reason"`
	Status           RefundStatus `json:"status"`
	DecidedBy        *string      `json:"decided_by"`
	DecisionNote     *string      `json:"decision_note"`
	DecidedAt        *time.Time   `json:"decided_at"`
	ExecutedBy       *string      `json:"executed_by"`
	ExecutedAt       *time.Time   `json:"executed_at"`
	ProviderRefundID *string      `json:"provider_refund_id"`
	CompletedAt      *time.Time   `json:"completed_at"`
	FailureReason    *string      `json:"failure_reason"`
	CreatedAt        time.Time    `json:"created_at"`
}
