package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type LedgerHandler struct {
	service *service.LedgerService
}

func NewLedgerHandler(service *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{
		service: service,
	}
}

func (h *LedgerHandler) ListJournalEntries(c *gin.Context) {
	entries, err := h.service.ListJournalEntries(
		c.Request.Context(),
		c.Query("organization_id"),
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLedgerRequest) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list journal entries",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": entries,
	})
}
