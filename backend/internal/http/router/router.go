package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/handler"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/repository"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *gin.Engine {
	engine := gin.Default()

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "ledgerflow-api",
		})
	})

	organizationRepository := repository.NewOrganizationRepository(db)
	organizationService := service.NewOrganizationService(organizationRepository)
	organizationHandler := handler.NewOrganizationHandler(organizationService)

	customerRepository := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepository)
	customerHandler := handler.NewCustomerHandler(customerService)

	api := engine.Group("/api/v1")
	{
		api.POST("/organizations", organizationHandler.Create)
		api.POST("/customers", customerHandler.Create)
	}

	return engine
}
