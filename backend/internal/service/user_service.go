package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/auth"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidUserOrganization = errors.New("invalid user organization")
	ErrInvalidUser             = errors.New("invalid user")
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailAlreadyUsed        = errors.New("email is already used")
	ErrCannotModifyOwnUser     = errors.New("cannot modify your own role or status")
	ErrLastFinanceAdmin        = errors.New("organization must keep at least one active finance admin")
)

type UserRepository interface {
	Create(
		ctx context.Context,
		organizationID string,
		name string,
		email string,
		passwordHash string,
		role model.UserRole,
	) (model.User, error)
	ListByOrganization(
		ctx context.Context,
		organizationID string,
	) ([]model.User, error)
	FindByIDAndOrganization(
		ctx context.Context,
		userID string,
		organizationID string,
	) (model.User, error)
	UpdateRole(
		ctx context.Context,
		userID string,
		organizationID string,
		role model.UserRole,
	) (model.User, error)
	UpdateStatus(
		ctx context.Context,
		userID string,
		organizationID string,
		isActive bool,
	) (model.User, error)
	CountActiveFinanceAdmins(
		ctx context.Context,
		organizationID string,
	) (int, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) List(
	ctx context.Context,
	organizationID string,
) ([]model.User, error) {
	organizationID = strings.TrimSpace(organizationID)
	if organizationID == "" {
		return nil, ErrInvalidUserOrganization
	}

	return s.repository.ListByOrganization(ctx, organizationID)
}

func (s *UserService) Create(
	ctx context.Context,
	organizationID string,
	name string,
	email string,
	password string,
	role model.UserRole,
) (model.User, error) {
	organizationID = strings.TrimSpace(organizationID)
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if organizationID == "" {
		return model.User{}, ErrInvalidUserOrganization
	}

	if err := validateUserInput(name, email, password, role); err != nil {
		return model.User{}, err
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return model.User{}, fmt.Errorf("hash user password: %w", err)
	}

	user, err := s.repository.Create(
		ctx,
		organizationID,
		name,
		email,
		passwordHash,
		role,
	)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return model.User{}, ErrEmailAlreadyUsed
		}

		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateRole(
	ctx context.Context,
	organizationID string,
	actorID string,
	userID string,
	role model.UserRole,
) (model.User, error) {
	if err := validateUserUpdate(organizationID, actorID, userID); err != nil {
		return model.User{}, err
	}

	if !isValidRole(role) {
		return model.User{}, fmt.Errorf("%w: role is invalid", ErrInvalidUser)
	}

	user, err := s.findUser(ctx, userID, organizationID)
	if err != nil {
		return model.User{}, err
	}

	if user.ID == actorID {
		return model.User{}, ErrCannotModifyOwnUser
	}

	if user.IsActive && user.Role == model.UserRoleFinanceAdmin && role != model.UserRoleFinanceAdmin {
		if err := s.ensureAnotherFinanceAdmin(ctx, organizationID); err != nil {
			return model.User{}, err
		}
	}

	updatedUser, err := s.repository.UpdateRole(ctx, userID, organizationID, role)
	if err != nil {
		return model.User{}, fmt.Errorf("update user role: %w", err)
	}

	return updatedUser, nil
}

func (s *UserService) UpdateStatus(
	ctx context.Context,
	organizationID string,
	actorID string,
	userID string,
	isActive bool,
) (model.User, error) {
	if err := validateUserUpdate(organizationID, actorID, userID); err != nil {
		return model.User{}, err
	}

	user, err := s.findUser(ctx, userID, organizationID)
	if err != nil {
		return model.User{}, err
	}

	if user.ID == actorID {
		return model.User{}, ErrCannotModifyOwnUser
	}

	if user.Role == model.UserRoleFinanceAdmin && user.IsActive && !isActive {
		if err := s.ensureAnotherFinanceAdmin(ctx, organizationID); err != nil {
			return model.User{}, err
		}
	}

	updatedUser, err := s.repository.UpdateStatus(ctx, userID, organizationID, isActive)
	if err != nil {
		return model.User{}, fmt.Errorf("update user status: %w", err)
	}

	return updatedUser, nil
}

func validateUserInput(
	name string,
	email string,
	password string,
	role model.UserRole,
) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidUser)
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%w: email is invalid", ErrInvalidUser)
	}

	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidUser)
	}

	if !isValidRole(role) {
		return fmt.Errorf("%w: role is invalid", ErrInvalidUser)
	}

	return nil
}

func validateUserUpdate(organizationID string, actorID string, userID string) error {
	if strings.TrimSpace(organizationID) == "" {
		return ErrInvalidUserOrganization
	}

	if strings.TrimSpace(actorID) == "" || strings.TrimSpace(userID) == "" {
		return fmt.Errorf("%w: user id is required", ErrInvalidUser)
	}

	return nil
}

func (s *UserService) findUser(
	ctx context.Context,
	userID string,
	organizationID string,
) (model.User, error) {
	user, err := s.repository.FindByIDAndOrganization(ctx, userID, organizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func (s *UserService) ensureAnotherFinanceAdmin(
	ctx context.Context,
	organizationID string,
) error {
	count, err := s.repository.CountActiveFinanceAdmins(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("count finance admins: %w", err)
	}

	if count <= 1 {
		return ErrLastFinanceAdmin
	}

	return nil
}

func isValidRole(role model.UserRole) bool {
	return role == model.UserRoleFinanceAdmin ||
		role == model.UserRoleApprover ||
		role == model.UserRoleViewer
}
