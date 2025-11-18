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
