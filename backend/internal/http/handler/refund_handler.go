package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/middleware"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type RefundHandler struct{ service *service.RefundService }

func NewRefundHandler(s *service.RefundService) *RefundHandler { return &RefundHandler{s} }

func (h *RefundHandler) List(c *gin.Context) {
	o, _ := middleware.OrganizationID(c)
	x, e := h.service.List(c, o)
	if e != nil {
		c.JSON(400, gin.H{"error": e.Error()})
		return
	}
	c.JSON(200, gin.H{"data": x})
}

func (h *RefundHandler) Create(c *gin.Context) {
	o, _ := middleware.OrganizationID(c)
	u, _ := c.Get(middleware.UserIDContextKey)
	var r struct {
		Reason string `json:"reason"`
	}
	if c.ShouldBindJSON(&r) != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	x, e := h.service.Create(c, o, c.Param("payment_id"), u.(string), r.Reason)
	if e != nil {
		c.JSON(400, gin.H{"error": e.Error()})
		return
	}
	c.JSON(201, gin.H{"data": x})
}

func (h *RefundHandler) Decide(c *gin.Context) {
	o, _ := middleware.OrganizationID(c)
	u, _ := c.Get(middleware.UserIDContextKey)
	var r struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if c.ShouldBindJSON(&r) != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	x, e := h.service.Decide(c, o, c.Param("refund_id"), u.(string), model.RefundStatus(r.Status), r.Note)
	if e != nil {
		c.JSON(400, gin.H{"error": e.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": x})
}

func (h *RefundHandler) Execute(c *gin.Context) {
	organizationID, exists := middleware.OrganizationID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session context",
		})
		return
	}

	userID, exists := c.Get(middleware.UserIDContextKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session context",
		})
		return
	}

	executedBy, ok := userID.(string)
	if !ok || executedBy == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session context",
		})
		return
	}

	refund, err := h.service.Execute(
		c.Request.Context(),
		organizationID,
		c.Param("refund_id"),
		executedBy,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefund):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, service.ErrRefundNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "approved refund request not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to execute refund",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": refund,
	})
}
