package model

import "time"

type UserRole string

const (
	UserRoleFinanceAdmin UserRole = "finance_admin"
	UserRoleApprover     UserRole = "approver"
	UserRoleViewer       UserRole = "viewer"
)

type User struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           UserRole  `json:"role"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
