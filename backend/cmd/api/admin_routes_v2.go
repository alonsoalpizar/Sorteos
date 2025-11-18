package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	adminHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/admin"
	"github.com/sorteos-platform/backend/internal/adapters/http/middleware"
	redisAdapter "github.com/sorteos-platform/backend/internal/adapters/redis"
	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"

	// Use cases
	categoryuc "github.com/sorteos-platform/backend/internal/usecase/admin/category"
	configuc "github.com/sorteos-platform/backend/internal/usecase/admin/config"

	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// setupAdminRoutesV2 configura las rutas de administración (versión simplificada)
// Solo incluye los endpoints que tienen use cases 100% completos
func setupAdminRoutesV2(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) {
	// Inicializar middleware
	tokenMgr := redisAdapter.NewTokenManager(rdb, &cfg.JWT)
	blacklistService := redisinfra.NewTokenBlacklistService(rdb)
	authMiddleware := middleware.NewAuthMiddleware(tokenMgr, blacklistService, log)

	// Grupo base de admin - requiere autenticación y rol admin/super_admin
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(authMiddleware.Authenticate())
	adminGroup.Use(authMiddleware.RequireRole("admin", "super_admin"))

	// ==================== CATEGORY MANAGEMENT ====================
	setupCategoryRoutesV2(adminGroup, gormDB, log)

	// ==================== SYSTEM CONFIG ====================
	setupConfigRoutesV2(adminGroup, gormDB, log)

	// ==================== SETTLEMENTS ====================
	setupSettlementRoutesV2(adminGroup, gormDB, log)

	// ==================== USER MANAGEMENT ====================
	setupUserRoutesV2(adminGroup, gormDB, log)

	// ==================== ORGANIZER MANAGEMENT ====================
	setupOrganizerRoutesV2(adminGroup, gormDB, log)

	// ==================== PAYMENT MANAGEMENT ====================
	setupPaymentRoutesV2(adminGroup, gormDB, log)

	// ==================== RAFFLE MANAGEMENT ====================
	setupRaffleRoutesV2(adminGroup, gormDB, log)

	// ==================== NOTIFICATIONS ====================
	setupNotificationRoutesV2(adminGroup, gormDB, log)

	// ==================== REPORTS & DASHBOARD ====================
	setupReportsRoutesV2(adminGroup, gormDB, log)
}

// setupCategoryRoutesV2 configura rutas de gestión de categorías
func setupCategoryRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar use cases
	createCategory := categoryuc.NewCreateCategoryUseCase(db, log)
	updateCategory := categoryuc.NewUpdateCategoryUseCase(db, log)
	deleteCategory := categoryuc.NewDeleteCategoryUseCase(db, log)
	listCategories := categoryuc.NewListCategoriesUseCase(db, log)

	// Inicializar handler
	handler := adminHandler.NewCategoryHandler(
		createCategory,
		updateCategory,
		deleteCategory,
		listCategories,
	)

	// Configurar rutas
	categories := adminGroup.Group("/categories")
	{
		categories.GET("", handler.ListCategories)        // GET /api/v1/admin/categories
		categories.POST("", handler.CreateCategory)       // POST /api/v1/admin/categories
		categories.PUT("/:id", handler.UpdateCategory)    // PUT /api/v1/admin/categories/:id
		categories.DELETE("/:id", handler.DeleteCategory) // DELETE /api/v1/admin/categories/:id
	}

	log.Info("Admin category routes registered",
		logger.Int("endpoints", 4),
		logger.String("base_path", "/api/v1/admin/categories"))
}

// setupConfigRoutesV2 configura rutas de configuración del sistema
func setupConfigRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar use cases
	getConfig := configuc.NewGetSystemConfigUseCase(db, log)
	updateConfig := configuc.NewUpdateSystemConfigUseCase(db, log)
	listConfigs := configuc.NewListSystemConfigsUseCase(db, log)

	// Inicializar handler
	handler := adminHandler.NewConfigHandler(
		getConfig,
		updateConfig,
		listConfigs,
	)

	// Configurar rutas
	config := adminGroup.Group("/config")
	{
		config.GET("", handler.ListConfigs)          // GET /api/v1/admin/config
		config.GET("/:key", handler.GetConfig)       // GET /api/v1/admin/config/:key
		config.PUT("/:key", handler.UpdateConfig)    // PUT /api/v1/admin/config/:key
	}

	log.Info("Admin config routes registered",
		logger.Int("endpoints", 3),
		logger.String("base_path", "/api/v1/admin/config"))
}

// setupSettlementRoutesV2 configura rutas de liquidaciones (settlements)
func setupSettlementRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar handler (el handler ya inicializa todos sus use cases internamente)
	handler := adminHandler.NewSettlementHandler(db, log)

	// Configurar rutas
	settlements := adminGroup.Group("/settlements")
	{
		settlements.GET("", handler.List)                       // GET /api/v1/admin/settlements
		settlements.GET("/:id", handler.GetByID)                // GET /api/v1/admin/settlements/:id
		settlements.POST("", handler.Create)                    // POST /api/v1/admin/settlements
		settlements.PUT("/:id/approve", handler.Approve)        // PUT /api/v1/admin/settlements/:id/approve
		settlements.PUT("/:id/reject", handler.Reject)          // PUT /api/v1/admin/settlements/:id/reject
		settlements.PUT("/:id/payout", handler.MarkPaid)        // PUT /api/v1/admin/settlements/:id/payout
		settlements.POST("/auto-create", handler.AutoCreate)    // POST /api/v1/admin/settlements/auto-create
	}

	log.Info("Admin settlement routes registered",
		logger.Int("endpoints", 7),
		logger.String("base_path", "/api/v1/admin/settlements"))
}

// setupUserRoutesV2 configura rutas de gestión de usuarios
func setupUserRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar handler (el handler ya inicializa todos sus use cases internamente)
	handler := adminHandler.NewUserHandler(db, log)

	// Configurar rutas
	users := adminGroup.Group("/users")
	{
		users.GET("", handler.List)                    // GET /api/v1/admin/users
		users.GET("/:id", handler.GetByID)             // GET /api/v1/admin/users/:id
		users.PUT("/:id/status", handler.UpdateStatus) // PUT /api/v1/admin/users/:id/status
		users.PUT("/:id/kyc", handler.UpdateKYC)       // PUT /api/v1/admin/users/:id/kyc
		users.DELETE("/:id", handler.Delete)           // DELETE /api/v1/admin/users/:id
		users.POST("/:id/reset-password", handler.ResetPassword) // POST /api/v1/admin/users/:id/reset-password
	}

	log.Info("Admin user routes registered",
		logger.Int("endpoints", 6),
		logger.String("base_path", "/api/v1/admin/users"))
}

// setupOrganizerRoutesV2 configura rutas de gestión de organizadores
func setupOrganizerRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar handler (el handler ya inicializa todos sus use cases internamente)
	handler := adminHandler.NewOrganizerHandler(db, log)

	// Configurar rutas
	organizers := adminGroup.Group("/organizers")
	{
		organizers.GET("", handler.List)                           // GET /api/v1/admin/organizers
		organizers.GET("/:id", handler.GetByID)                    // GET /api/v1/admin/organizers/:id
		organizers.PUT("/:id/commission", handler.UpdateCommission) // PUT /api/v1/admin/organizers/:id/commission
		organizers.PUT("/:id/verify", handler.Verify)              // PUT /api/v1/admin/organizers/:id/verify
		organizers.GET("/:id/revenue", handler.GetRevenue)        // GET /api/v1/admin/organizers/:id/revenue
	}

	log.Info("Admin organizer routes registered",
		logger.Int("endpoints", 4),
		logger.String("base_path", "/api/v1/admin/organizers"))
}

// setupPaymentRoutesV2 configura rutas de gestión de pagos
func setupPaymentRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar handler (el handler ya inicializa todos sus use cases internamente)
	handler := adminHandler.NewPaymentHandler(db, log)

	// Configurar rutas
	payments := adminGroup.Group("/payments")
	{
		payments.GET("", handler.List)                      // GET /api/v1/admin/payments
		payments.GET("/:id", handler.GetByID)               // GET /api/v1/admin/payments/:id
		payments.POST("/:id/refund", handler.ProcessRefund) // POST /api/v1/admin/payments/:id/refund
		payments.POST("/:id/dispute", handler.ManageDispute) // POST /api/v1/admin/payments/:id/dispute
	}

	log.Info("Admin payment routes registered",
		logger.Int("endpoints", 4),
		logger.String("base_path", "/api/v1/admin/payments"))
}

// setupRaffleRoutesV2 configura rutas de gestión de rifas
func setupRaffleRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar handler (el handler ya inicializa todos sus use cases internamente)
	handler := adminHandler.NewRaffleHandler(db, log)

	// Configurar rutas
	raffles := adminGroup.Group("/raffles")
	{
		raffles.GET("", handler.List)                             // GET /api/v1/admin/raffles
		raffles.GET("/:id/transactions", handler.ViewTransactions) // GET /api/v1/admin/raffles/:id/transactions
		raffles.PUT("/:id/status", handler.ForceStatusChange)     // PUT /api/v1/admin/raffles/:id/status
		raffles.POST("/:id/draw", handler.ManualDraw)             // POST /api/v1/admin/raffles/:id/draw
		raffles.POST("/:id/notes", handler.AddNotes)              // POST /api/v1/admin/raffles/:id/notes
		raffles.POST("/:id/cancel", handler.CancelWithRefund)     // POST /api/v1/admin/raffles/:id/cancel
	}

	log.Info("Admin raffle routes registered",
		logger.Int("endpoints", 6),
		logger.String("base_path", "/api/v1/admin/raffles"))
}

// setupNotificationRoutesV2 configura rutas de notificaciones
func setupNotificationRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar handler (el handler ya inicializa todos sus use cases internamente)
	handler := adminHandler.NewNotificationHandler(db, log)

	// Configurar rutas
	notifications := adminGroup.Group("/notifications")
	{
		notifications.POST("/email", handler.SendEmail)                  // POST /api/v1/admin/notifications/email
		notifications.POST("/bulk", handler.SendBulkEmail)               // POST /api/v1/admin/notifications/bulk
		notifications.POST("/templates", handler.ManageTemplates)        // POST /api/v1/admin/notifications/templates
		notifications.POST("/announcements", handler.CreateAnnouncement) // POST /api/v1/admin/notifications/announcements
		notifications.GET("/history", handler.ViewHistory)               // GET /api/v1/admin/notifications/history
	}

	log.Info("Admin notification routes registered",
		logger.Int("endpoints", 6),
		logger.String("base_path", "/api/v1/admin/notifications"))
}

// setupReportsRoutesV2 configura rutas de reportes y dashboard
func setupReportsRoutesV2(adminGroup *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// Inicializar handler (el handler ya inicializa todos sus use cases internamente)
	handler := adminHandler.NewReportsHandler(db, log)

	// Configurar rutas
	reports := adminGroup.Group("/reports")
	{
		reports.GET("/dashboard", handler.GetDashboard)                 // GET /api/v1/admin/reports/dashboard
		reports.GET("/revenue", handler.GetRevenueReport)               // GET /api/v1/admin/reports/revenue
		reports.GET("/liquidations", handler.GetLiquidationsReport)     // GET /api/v1/admin/reports/liquidations
		reports.POST("/export", handler.ExportData)                     // POST /api/v1/admin/reports/export
	}

	log.Info("Admin reports routes registered",
		logger.Int("endpoints", 4),
		logger.String("base_path", "/api/v1/admin/reports"))
}
