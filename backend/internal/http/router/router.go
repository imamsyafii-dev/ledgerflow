package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/config"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/handler"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/middleware"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/model"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/repository"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool, cfg config.Config) *gin.Engine {
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowCredentials: true,
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

	authUserRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(
		authUserRepository,
		cfg.JWTSecret,
	)
	authHandler := handler.NewAuthHandler(
		authService,
		cfg.AppEnv == "production",
	)
	userService := service.NewUserService(authUserRepository)
	userHandler := handler.NewUserHandler(userService)

	engine.POST("/api/v1/organizations", organizationHandler.Create)

	customerRepository := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepository)
	customerHandler := handler.NewCustomerHandler(customerService)

	invoiceRepository := repository.NewInvoiceRepository(db)
	invoiceService := service.NewInvoiceService(invoiceRepository)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)

	paymentRepository := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepository)
	paymentWebhookHandler := handler.NewPaymentWebhookHandler(paymentService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	ledgerRepository := repository.NewLedgerRepository(db)
	ledgerService := service.NewLedgerService(ledgerRepository)
	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	reconciliationRepository := repository.NewReconciliationRepository(db)
	reconciliationService := service.NewReconciliationService(
		reconciliationRepository,
	)
	reconciliationHandler := handler.NewReconciliationHandler(
		reconciliationService,
	)

	api := engine.Group("/api/v1")
	api.Use(middleware.RequireAuth(cfg.JWTSecret, db))
	{
		api.POST("/customers", customerHandler.Create)
		api.GET("/customers", customerHandler.List)
		api.GET(
			"/users",
			middleware.RequireRole(model.UserRoleFinanceAdmin),
			userHandler.List,
		)
		api.POST(
			"/users",
			middleware.RequireRole(model.UserRoleFinanceAdmin),
			userHandler.Create,
		)
		api.PATCH(
			"/users/:user_id/role",
			middleware.RequireRole(model.UserRoleFinanceAdmin),
			userHandler.UpdateRole,
		)
		api.PATCH(
			"/users/:user_id/status",
			middleware.RequireRole(model.UserRoleFinanceAdmin),
			userHandler.UpdateStatus,
		)
		api.POST("/invoices", invoiceHandler.Create)
		api.GET("/invoices", invoiceHandler.List)
		api.POST("/invoices/:invoice_id/issue", invoiceHandler.Issue)
		api.GET("/ledger/journal-entries", ledgerHandler.ListJournalEntries)
		api.GET("/payments", paymentHandler.List)
		api.GET("/reconciliations", reconciliationHandler.List)
	}

	authRoutes := engine.Group("/api/v1/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)
		authRoutes.GET(
			"/me",
			middleware.RequireAuth(cfg.JWTSecret, db),
			authHandler.Me,
		)
	}

	engine.POST("/webhooks/mock-gateway", paymentWebhookHandler.HandleSucceeded)

	return engine
}
