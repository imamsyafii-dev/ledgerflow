package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type ReconciliationHandler struct {
	service *service.ReconciliationService
}

func NewReconciliationHandler(
	service *service.ReconciliationService,
) *ReconciliationHandler {
	return &ReconciliationHandler{
		service: service,
	}
}

func (h *ReconciliationHandler) List(c *gin.Context) {
	reconciliations, err := h.service.List(
		c.Request.Context(),
		c.Query("organization_id"),
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidReconciliationRequest) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list reconciliations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reconciliations,
	})
}
