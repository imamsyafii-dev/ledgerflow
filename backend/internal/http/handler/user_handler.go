package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/middleware"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type updateUserRoleRequest struct {
	Role string `json:"role"`
}

type updateUserStatusRequest struct {
	IsActive *bool `json:"is_active"`
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	organizationID, exists := middleware.OrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session context",
		})
		return
	}

	users, err := h.service.List(c.Request.Context(), organizationID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserOrganization) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	organizationID, exists := middleware.OrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session context"})
		return
	}

	var request createUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.Create(
		c.Request.Context(),
		organizationID,
		request.Name,
		request.Email,
		request.Password,
		model.UserRole(request.Role),
	)
	if err != nil {
		h.writeUserError(c, err, "failed to create user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	organizationID, actorID, ok := requestUserContext(c)
	if !ok {
		return
	}

	var request updateUserRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.service.UpdateRole(
		c.Request.Context(),
		organizationID,
		actorID,
		c.Param("user_id"),
		model.UserRole(request.Role),
	)
	if err != nil {
		h.writeUserError(c, err, "failed to update user role")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	organizationID, actorID, ok := requestUserContext(c)
	if !ok {
		return
	}

	var request updateUserStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil || request.IsActive == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_active is required"})
		return
	}

	user, err := h.service.UpdateStatus(
		c.Request.Context(),
		organizationID,
		actorID,
		c.Param("user_id"),
		*request.IsActive,
	)
	if err != nil {
		h.writeUserError(c, err, "failed to update user status")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func requestUserContext(c *gin.Context) (string, string, bool) {
	organizationID, organizationExists := middleware.OrganizationID(c)
	actorID, actorExists := c.Get(middleware.UserIDContextKey)
	if !organizationExists || !actorExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session context"})
		return "", "", false
	}

	actorIDValue, ok := actorID.(string)
	if !ok || actorIDValue == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session context"})
		return "", "", false
	}

	return organizationID, actorIDValue, true
}

func (h *UserHandler) writeUserError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrInvalidUser),
		errors.Is(err, service.ErrInvalidUserOrganization):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrEmailAlreadyUsed),
		errors.Is(err, service.ErrCannotModifyOwnUser),
		errors.Is(err, service.ErrLastFinanceAdmin):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fallback})
	}
}
