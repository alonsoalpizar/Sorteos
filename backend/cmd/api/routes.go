package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	authHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/auth"
	raffleHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/raffle"
	websocketHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/websocket"
	"github.com/sorteos-platform/backend/internal/adapters/http/middleware"
	"github.com/sorteos-platform/backend/internal/adapters/db"
	redisAdapter "github.com/sorteos-platform/backend/internal/adapters/redis"
	"github.com/sorteos-platform/backend/internal/adapters/notifier"
	"github.com/sorteos-platform/backend/internal/usecase/auth"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
	"github.com/sorteos-platform/backend/internal/infrastructure/websocket"
	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"
	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// setupAuthRoutes configura las rutas de autenticación y retorna el email notifier para testing
func setupAuthRoutes(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) notifier.Notifier {
	// Inicializar repositorios
	userRepo := db.NewUserRepository(gormDB)
	consentRepo := db.NewUserConsentRepository(gormDB)
	auditRepo := db.NewAuditLogRepository(gormDB)

	// Inicializar token manager
	tokenMgr := redisAdapter.NewTokenManager(rdb, &cfg.JWT)

	// Inicializar blacklist service
	blacklistService := redisinfra.NewTokenBlacklistService(rdb)

	// Inicializar notifier (SMTP o SendGrid según configuración)
	var emailNotifier notifier.Notifier
	if cfg.EmailProvider == "smtp" {
		emailNotifier = notifier.NewSMTPNotifier(&cfg.SMTP, log)
		log.Info("Email provider configured",
			logger.String("provider", "smtp"),
			logger.String("host", cfg.SMTP.Host),
			logger.Int("port", cfg.SMTP.Port),
		)
	} else {
		emailNotifier = notifier.NewSendGridNotifier(&cfg.SendGrid, log)
		log.Info("Email provider configured",
			logger.String("provider", "sendgrid"),
		)
	}

	// Inicializar middlewares
	authMiddleware := middleware.NewAuthMiddleware(tokenMgr, blacklistService, log)
	rateLimiter := middleware.NewRateLimiter(rdb, log)

	// Inicializar use cases
	registerUseCase := auth.NewRegisterUseCase(userRepo, consentRepo, auditRepo, tokenMgr, emailNotifier, log, cfg.SkipEmailVerification)
	loginUseCase := auth.NewLoginUseCase(userRepo, auditRepo, tokenMgr, log)
	refreshTokenUseCase := auth.NewRefreshTokenUseCase(userRepo, tokenMgr, log)
	verifyEmailUseCase := auth.NewVerifyEmailUseCase(userRepo, auditRepo, tokenMgr, log)
	logoutUseCase := auth.NewLogoutUseCase(blacklistService)

	// Inicializar handlers
	registerHandler := authHandler.NewRegisterHandler(registerUseCase, log)
	loginHandler := authHandler.NewLoginHandler(loginUseCase, log)
	refreshHandler := authHandler.NewRefreshTokenHandler(refreshTokenUseCase, log)
	verifyEmailHandler := authHandler.NewVerifyEmailHandler(verifyEmailUseCase, log)
	logoutHandler := authHandler.NewLogoutHandler(logoutUseCase, log)

	// Grupo de rutas de autenticación
	authGroup := router.Group("/api/v1/auth")
	{
		// Rutas públicas con rate limiting
		authGroup.POST("/register",
			rateLimiter.LimitByIP(5, time.Minute),
			registerHandler.Handle,
		)

		authGroup.POST("/login",
			rateLimiter.LimitByIP(5, time.Minute),
			loginHandler.Handle,
		)

		authGroup.POST("/refresh",
			rateLimiter.LimitByIP(10, time.Minute),
			refreshHandler.Handle,
		)

		authGroup.POST("/verify-email",
			rateLimiter.LimitByIP(10, time.Minute),
			verifyEmailHandler.Handle,
		)

		// Rutas protegidas
		protected := authGroup.Group("")
		protected.Use(authMiddleware.Authenticate())
		{
			protected.POST("/logout", logoutHandler.Handle)
		}
	}

	// Grupo de ejemplo de rutas protegidas por rol
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(authMiddleware.Authenticate())
	adminGroup.Use(authMiddleware.RequireRole("admin", "super_admin"))
	{
		// TODO: Implementar endpoints de admin
		adminGroup.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Admin endpoint - lista de usuarios",
			})
		})
	}

	// Ruta de ejemplo protegida por KYC
	router.GET("/api/v1/profile",
		authMiddleware.Authenticate(),
		authMiddleware.RequireMinKYC("email_verified"),
		func(c *gin.Context) {
			userID, _ := middleware.GetUserID(c)
			c.JSON(http.StatusOK, gin.H{
				"message": "Perfil de usuario",
				"user_id": userID,
			})
		},
	)

	// Retornar el email notifier para uso en testing
	return emailNotifier
}

// setupRaffleRoutes configura las rutas de sorteos
func setupRaffleRoutes(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) {
	// Inicializar repositorios
	raffleRepo := db.NewRaffleRepository(gormDB)
	raffleNumberRepo := db.NewRaffleNumberRepository(gormDB)
	raffleImageRepo := db.NewRaffleImageRepository(gormDB)
	userRepo := db.NewUserRepository(gormDB)
	auditRepo := db.NewAuditLogRepository(gormDB)

	// Inicializar token manager y auth middleware
	tokenMgr := redisAdapter.NewTokenManager(rdb, &cfg.JWT)
	blacklistService := redisinfra.NewTokenBlacklistService(rdb)
	authMiddleware := middleware.NewAuthMiddleware(tokenMgr, blacklistService, log)
	rateLimiter := middleware.NewRateLimiter(rdb, log)

	// Inicializar use cases
	createRaffleUseCase := raffleuc.NewCreateRaffleUseCase(
		raffleRepo,
		raffleNumberRepo,
		userRepo,
		auditRepo,
		log,
	)
	listRafflesUseCase := raffleuc.NewListRafflesUseCase(raffleRepo)
	getRaffleDetailUseCase := raffleuc.NewGetRaffleDetailUseCase(
		raffleRepo,
		raffleNumberRepo,
		raffleImageRepo,
	)
	publishRaffleUseCase := raffleuc.NewPublishRaffleUseCase(
		raffleRepo,
		raffleImageRepo,
		raffleNumberRepo,
		auditRepo,
	)
	updateRaffleUseCase := raffleuc.NewUpdateRaffleUseCase(raffleRepo, auditRepo)
	suspendRaffleUseCase := raffleuc.NewSuspendRaffleUseCase(raffleRepo, auditRepo)
	deleteRaffleUseCase := raffleuc.NewDeleteRaffleUseCase(raffleRepo, auditRepo)
	getUserTicketsUseCase := raffleuc.NewGetUserTicketsUseCase(raffleNumberRepo, raffleRepo)

	// Inicializar handlers
	createRaffleHandler := raffleHandler.NewCreateRaffleHandler(createRaffleUseCase)
	listRafflesHandler := raffleHandler.NewListRafflesHandler(listRafflesUseCase)
	getRaffleDetailHandler := raffleHandler.NewGetRaffleDetailHandler(getRaffleDetailUseCase)
	publishRaffleHandler := raffleHandler.NewPublishRaffleHandler(publishRaffleUseCase)
	updateRaffleHandler := raffleHandler.NewUpdateRaffleHandler(updateRaffleUseCase)
	suspendRaffleHandler := raffleHandler.NewSuspendRaffleHandler(suspendRaffleUseCase)
	deleteRaffleHandler := raffleHandler.NewDeleteRaffleHandler(deleteRaffleUseCase)
	getUserTicketsHandler := raffleHandler.NewGetUserTicketsHandler(getUserTicketsUseCase)

	// Grupo de rutas de sorteos
	rafflesGroup := router.Group("/api/v1/raffles")
	{
		// Rutas públicas
		rafflesGroup.GET("", listRafflesHandler.Handle)                          // Listar sorteos

		// Rutas protegidas (requieren autenticación + email verificado)
		protected := rafflesGroup.Group("")
		protected.Use(authMiddleware.Authenticate())
		protected.Use(authMiddleware.RequireMinKYC("email_verified"))
		{
			// IMPORTANTE: Rutas específicas ANTES de rutas con parámetros :id
			protected.GET("/my-tickets", getUserTicketsHandler.Handle)  // Obtener tickets del usuario

			protected.POST("",
				rateLimiter.LimitByUser(10, time.Hour),  // Max 10 sorteos por hora
				createRaffleHandler.Handle,
			)
			protected.PUT("/:id", updateRaffleHandler.Handle)        // Actualizar sorteo
			protected.POST("/:id/publish", publishRaffleHandler.Handle)  // Publicar sorteo
			protected.DELETE("/:id", deleteRaffleHandler.Handle)      // Eliminar sorteo (soft delete)
		}

		// Detalle de sorteo - DESPUÉS de rutas específicas para evitar conflictos
		rafflesGroup.GET("/:id", getRaffleDetailHandler.Handle)                  // Detalle de sorteo

		// Rutas de admin
		admin := rafflesGroup.Group("")
		admin.Use(authMiddleware.Authenticate())
		admin.Use(authMiddleware.RequireRole("admin", "super_admin"))
		{
			admin.POST("/:id/suspend", suspendRaffleHandler.Handle)  // Suspender sorteo
		}
	}
}

// setupWebSocketRoutes configura las rutas de WebSocket
func setupWebSocketRoutes(router *gin.Engine, wsHub *websocket.Hub, rdb *redis.Client, cfg *config.Config, log *logger.Logger) {
	// Inicializar token manager y auth middleware (opcional para WebSocket)
	tokenMgr := redisAdapter.NewTokenManager(rdb, &cfg.JWT)
	blacklistService := redisinfra.NewTokenBlacklistService(rdb)
	authMiddleware := middleware.NewAuthMiddleware(tokenMgr, blacklistService, log)

	// Inicializar handler
	wsHandler := websocketHandler.NewWebSocketHandler(wsHub)

	// Grupo de rutas WebSocket
	rafflesGroup := router.Group("/api/v1/raffles")
	{
		// WebSocket connection endpoint (público con autenticación opcional)
		rafflesGroup.GET("/:id/ws",
			// Opcional: agregar authMiddleware.Authenticate() para requerir login
			wsHandler.HandleConnection,
		)

		// Stats endpoint (solo para admin)
		rafflesGroup.GET("/:id/ws/stats",
			authMiddleware.Authenticate(),
			authMiddleware.RequireRole("admin", "super_admin"),
			wsHandler.GetConnectionStats,
		)
	}

	// Endpoint global de stats (admin only)
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(authMiddleware.Authenticate())
	adminGroup.Use(authMiddleware.RequireRole("admin", "super_admin"))
	{
		adminGroup.GET("/websocket/stats", wsHandler.GetGlobalStats)
	}
}
