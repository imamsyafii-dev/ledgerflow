package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) List(c *gin.Context) {
	payments, err := h.service.List(
		c.Request.Context(),
		c.Query("organization_id"),
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPaymentRequest) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list payments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": payments,
	})
}
