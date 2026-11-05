package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
)

type OrganizationHandler struct {
	service *service.OrganizationService
}

type createOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func NewOrganizationHandler(
	service *service.OrganizationService,
) *OrganizationHandler {
	return &OrganizationHandler{
		service: service,
	}
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	var request createOrganizationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	organization, err := h.service.Create(
		c.Request.Context(),
		request.Name,
		request.Slug,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrganization) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create organization",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": organization,
	})
}
