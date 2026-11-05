package model

import "time"

type InvoiceStatus string

const (
	InvoiceStatusDraft InvoiceStatus = "DRAFT"
)

type Invoice struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"organization_id"`
	CustomerID     string        `json:"customer_id"`
	Number         string        `json:"number"`
	Amount         int64         `json:"amount"`
	Currency       string        `json:"currency"`
	Status         InvoiceStatus `json:"status"`
	DueDate        time.Time     `json:"due_date"`
	PaidAt         *time.Time    `json:"paid_at"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
