package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/middleware"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type CustomerHandler struct {
	service *service.CustomerService
}

type createCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: service,
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var request createCustomerRequest
	organizationID, exists := middleware.OrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session context",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	customer, err := h.service.Create(
		c.Request.Context(),
		organizationID,
		request.Name,
		request.Email,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCustomer) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create customer",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": customer,
	})
}

func (h *CustomerHandler) List(c *gin.Context) {
	organizationID, exists := middleware.OrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session context",
		})
		return
	}

	customers, err := h.service.List(
		c.Request.Context(),
		organizationID,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCustomer) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list customers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": customers,
	})
}
