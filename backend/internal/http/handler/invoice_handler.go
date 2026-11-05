package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type InvoiceHandler struct {
	service *service.InvoiceService
}

type createInvoiceRequest struct {
	OrganizationID string `json:"organization_id"`
	CustomerID     string `json:"customer_id"`
	Number         string `json:"number"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	DueDate        string `json:"due_date"`
}

type issueInvoiceRequest struct {
	OrganizationID string `json:"organization_id"`
}

func NewInvoiceHandler(service *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		service: service,
	}
}

func (h *InvoiceHandler) Create(c *gin.Context) {
	var request createInvoiceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	invoice, err := h.service.Create(
		c.Request.Context(),
		request.OrganizationID,
		request.CustomerID,
		request.Number,
		request.Amount,
		request.Currency,
		request.DueDate,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInvoice) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create invoice",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": invoice,
	})
}

func (h *InvoiceHandler) List(c *gin.Context) {
	organizationID := c.Query("organization_id")

	invoices, err := h.service.List(
		c.Request.Context(),
		organizationID,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInvoice) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list invoices",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": invoices,
	})
}

func (h *InvoiceHandler) Issue(c *gin.Context) {
	var request issueInvoiceRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	invoice, err := h.service.Issue(
		c.Request.Context(),
		request.OrganizationID,
		c.Param("invoice_id"),
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInvoice) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to issue invoice",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": invoice,
	})
}
