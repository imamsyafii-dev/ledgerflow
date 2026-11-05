package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/handler"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/repository"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) *gin.Engine {
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		MaxAge: 12 * time.Hour,
	}))

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

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)

	paymentRepository := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepository)
	paymentWebhookHandler := handler.NewPaymentWebhookHandler(paymentService)

	ledgerRepository := repository.NewLedgerRepository(db)
	ledgerService := service.NewLedgerService(ledgerRepository)
	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	api := engine.Group("/api/v1")
	{
		api.POST("/organizations", organizationHandler.Create)
		api.POST("/customers", customerHandler.Create)
		api.GET("/customers", customerHandler.List)
		api.POST("/invoices", invoiceHandler.Create)
		api.GET("/invoices", invoiceHandler.List)
		api.POST("/invoices/:invoice_id/issue", invoiceHandler.Issue)
		api.GET("/ledger/journal-entries", ledgerHandler.ListJournalEntries)
	}

	engine.POST("/webhooks/mock-gateway", paymentWebhookHandler.HandleSucceeded)

	return engine
}
