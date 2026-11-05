package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type PaymentWebhookHandler struct {
	service *service.PaymentService
}

func NewPaymentWebhookHandler(
	service *service.PaymentService,
) *PaymentWebhookHandler {
	return &PaymentWebhookHandler{
		service: service,
	}
}

func (h *PaymentWebhookHandler) HandleSucceeded(c *gin.Context) {
	var event model.PaymentSucceededWebhook

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to encode webhook payload",
		})
		return
	}

	payment, duplicate, err := h.service.ProcessSucceededWebhook(
		c.Request.Context(),
		event,
		payload,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPaymentWebhook) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, service.ErrPaymentRejected) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process payment webhook",
		})
		return
	}

	if duplicate {
		c.JSON(http.StatusOK, gin.H{
			"status": "already_processed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": payment,
	})
}
