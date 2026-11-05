package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type CustomerHandler struct {
	service *service.CustomerService
}

type createCustomerRequest struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: service,
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var request createCustomerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	customer, err := h.service.Create(
		c.Request.Context(),
		request.OrganizationID,
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
	customers, err := h.service.List(
		c.Request.Context(),
		c.Query("organization_id"),
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
