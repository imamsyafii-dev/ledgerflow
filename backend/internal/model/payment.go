package model

import "time"

type PaymentStatus string

const (
	PaymentStatusSucceeded PaymentStatus = "SUCCEEDED"
)

type Payment struct {
	ID                string        `json:"id"`
	OrganizationID    string        `json:"organization_id"`
	InvoiceID         string        `json:"invoice_id"`
	InvoiceNumber     string        `json:"invoice_number"`
	CustomerName      string        `json:"customer_name"`
	Provider          string        `json:"provider"`
	ProviderPaymentID string        `json:"provider_payment_id"`
	Amount            int64         `json:"amount"`
	Currency          string        `json:"currency"`
	Status            PaymentStatus `json:"status"`
	PaidAt            *time.Time    `json:"paid_at"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type PaymentSucceededWebhook struct {
	EventID           string `json:"event_id"`
	EventType         string `json:"event_type"`
	ProviderPaymentID string `json:"payment_id"`
	OrganizationID    string `json:"organization_id"`
	InvoiceID         string `json:"invoice_id"`
	Amount            int64  `json:"amount"`
	Currency          string `json:"currency"`
}
