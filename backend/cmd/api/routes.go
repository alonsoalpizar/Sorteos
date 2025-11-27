package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	authHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/auth"
	categoryHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/category"
	imageHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/image"
	profileHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/profile"
	raffleHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/raffle"
	websocketHandler "github.com/sorteos-platform/backend/internal/adapters/http/handler/websocket"
	"github.com/sorteos-platform/backend/internal/adapters/http/middleware"
	"github.com/sorteos-platform/backend/internal/adapters/db"
	redisAdapter "github.com/sorteos-platform/backend/internal/adapters/redis"
	"github.com/sorteos-platform/backend/internal/adapters/notifier"
	"github.com/sorteos-platform/backend/internal/usecase/auth"
	categoryuc "github.com/sorteos-platform/backend/internal/usecase/category"
	creditsuc "github.com/sorteos-platform/backend/internal/usecase/credits"
	imageuc "github.com/sorteos-platform/backend/internal/usecase/image"
	profileuc "github.com/sorteos-platform/backend/internal/usecase/profile"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
	walletuc "github.com/sorteos-platform/backend/internal/usecase/wallet"
	"github.com/sorteos-platform/backend/internal/infrastructure/websocket"
	"github.com/sorteos-platform/backend/internal/infrastructure/pagadito"
	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"
	"github.com/sorteos-platform/backend/cmd/api/handlers"
	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// setupAuthRoutes configura las rutas de autenticación y retorna el email notifier para testing
func setupAuthRoutes(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) notifier.Notifier {
	// Inicializar repositorios
	userRepo := db.NewUserRepository(gormDB)
	walletRepo := db.NewWalletRepository(gormDB, log)
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
	registerUseCase := auth.NewRegisterUseCase(userRepo, walletRepo, consentRepo, auditRepo, tokenMgr, emailNotifier, log, cfg.SkipEmailVerification)
	loginUseCase := auth.NewLoginUseCase(userRepo, auditRepo, tokenMgr, log)
	refreshTokenUseCase := auth.NewRefreshTokenUseCase(userRepo, tokenMgr, log)
	verifyEmailUseCase := auth.NewVerifyEmailUseCase(userRepo, auditRepo, tokenMgr, log)
	logoutUseCase := auth.NewLogoutUseCase(blacklistService)
	googleAuthUseCase := auth.NewGoogleAuthUseCase(userRepo, walletRepo, auditRepo, tokenMgr, log)
	googleLinkUseCase := auth.NewGoogleLinkUseCase(userRepo, auditRepo, tokenMgr, googleAuthUseCase, log)

	// Inicializar handlers
	registerHandler := authHandler.NewRegisterHandler(registerUseCase, log)
	loginHandler := authHandler.NewLoginHandler(loginUseCase, log)
	refreshHandler := authHandler.NewRefreshTokenHandler(refreshTokenUseCase, log)
	verifyEmailHandler := authHandler.NewVerifyEmailHandler(verifyEmailUseCase, log)
	logoutHandler := authHandler.NewLogoutHandler(logoutUseCase, log)
	googleAuthHandler := authHandler.NewGoogleAuthHandler(googleAuthUseCase, log)
	googleLinkHandler := authHandler.NewGoogleLinkHandler(googleLinkUseCase, log)

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

		// Google OAuth endpoints
		authGroup.POST("/google",
			rateLimiter.LimitByIP(10, time.Minute),
			googleAuthHandler.Handle,
		)

		authGroup.POST("/google/link",
			rateLimiter.LimitByIP(5, time.Minute),
			googleLinkHandler.Handle,
		)

		// Rutas protegidas
		protected := authGroup.Group("")
		protected.Use(authMiddleware.Authenticate())
		{
			protected.POST("/logout", logoutHandler.Handle)
		}
	}

	// Retornar el email notifier para uso en testing
	return emailNotifier
}

// setupRaffleRoutes configura las rutas de sorteos
func setupRaffleRoutes(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) {
	// Inicializar repositorios
	raffleRepo := db.NewRaffleRepository(gormDB)
	raffleNumberRepo := db.NewRaffleNumberRepository(gormDB)
	raffleImageRepo := db.NewRaffleImageRepository(gormDB)
	categoryRepo := db.NewCategoryRepository(gormDB, log)
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
		gormDB,
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

	// Use case de categorías
	listCategoriesUseCase := categoryuc.NewListCategoriesUseCase(categoryRepo, log)

	// Use cases de imágenes
	uploadDir := "/var/www/sorteos.club/uploads/raffles"
	baseURL := "https://sorteos.club"
	uploadImageUseCase := imageuc.NewUploadImageUseCase(raffleRepo, raffleImageRepo, log, uploadDir, baseURL)
	deleteImageUseCase := imageuc.NewDeleteImageUseCase(raffleRepo, raffleImageRepo, log, uploadDir)
	setPrimaryImageUseCase := imageuc.NewSetPrimaryImageUseCase(raffleRepo, raffleImageRepo, log)

	// Inicializar handlers
	createRaffleHandler := raffleHandler.NewCreateRaffleHandler(createRaffleUseCase)
	listRafflesHandler := raffleHandler.NewListRafflesHandler(listRafflesUseCase)
	getRaffleDetailHandler := raffleHandler.NewGetRaffleDetailHandler(getRaffleDetailUseCase, raffleNumberRepo)
	publishRaffleHandler := raffleHandler.NewPublishRaffleHandler(publishRaffleUseCase)
	updateRaffleHandler := raffleHandler.NewUpdateRaffleHandler(updateRaffleUseCase)
	suspendRaffleHandler := raffleHandler.NewSuspendRaffleHandler(suspendRaffleUseCase)
	deleteRaffleHandler := raffleHandler.NewDeleteRaffleHandler(deleteRaffleUseCase)
	getUserTicketsHandler := raffleHandler.NewGetUserTicketsHandler(getUserTicketsUseCase)

	// Handler de categorías
	listCategoriesHandler := categoryHandler.NewListCategoriesHandler(listCategoriesUseCase)

	// Handlers de imágenes
	uploadImageHandler := imageHandler.NewUploadImageHandler(uploadImageUseCase)
	deleteImageHandler := imageHandler.NewDeleteImageHandler(deleteImageUseCase)
	setPrimaryImageHandler := imageHandler.NewSetPrimaryImageHandler(setPrimaryImageUseCase)

	// Ruta pública de categorías
	router.GET("/api/v1/categories", listCategoriesHandler.Handle)

	// Grupo de rutas de sorteos
	rafflesGroup := router.Group("/api/v1/raffles")
	{
		// Rutas con autenticación opcional (personaliza respuesta si está autenticado)
		rafflesGroup.GET("", authMiddleware.OptionalAuth(), listRafflesHandler.Handle) // Listar sorteos

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

			// Rutas de imágenes
			protected.POST("/:id/images", uploadImageHandler.Handle)                    // Subir imagen
			protected.DELETE("/:id/images/:image_id", deleteImageHandler.Handle)        // Eliminar imagen
			protected.PUT("/:id/images/:image_id/primary", setPrimaryImageHandler.Handle) // Establecer primaria
		}

		// Detalle de sorteo - DESPUÉS de rutas específicas para evitar conflictos
		// Usa OptionalAuth para personalizar según si el usuario está logueado
		rafflesGroup.GET("/:id", authMiddleware.OptionalAuth(), getRaffleDetailHandler.Handle)

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

// setupProfileRoutes configura las rutas de perfil de usuario
func setupProfileRoutes(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) {
	// Inicializar repositorios
	userRepo := db.NewUserRepository(gormDB)
	kycDocumentRepo := db.NewKYCDocumentRepository(gormDB)
	walletRepo := db.NewWalletRepository(gormDB, log)

	// Inicializar token manager y middlewares
	tokenMgr := redisAdapter.NewTokenManager(rdb, &cfg.JWT)
	blacklistService := redisinfra.NewTokenBlacklistService(rdb)
	authMiddleware := middleware.NewAuthMiddleware(tokenMgr, blacklistService, log)

	// Inicializar use cases
	getProfileUC := profileuc.NewGetProfileUseCase(userRepo, kycDocumentRepo, walletRepo)
	updateProfileUC := profileuc.NewUpdateProfileUseCase(userRepo)
	uploadPhotoUC := profileuc.NewUploadProfilePhotoUseCase(userRepo)
	configureIBANUC := profileuc.NewConfigureIBANUseCase(userRepo)
	uploadKYCDocumentUC := profileuc.NewUploadKYCDocumentUseCase(userRepo, kycDocumentRepo)

	// Inicializar handler
	profileHdlr := profileHandler.NewProfileHandler(
		getProfileUC,
		updateProfileUC,
		uploadPhotoUC,
		configureIBANUC,
		uploadKYCDocumentUC,
	)

	// Grupo de rutas de perfil (todas requieren autenticación)
	profileGroup := router.Group("/api/v1/profile")
	profileGroup.Use(authMiddleware.Authenticate())
	{
		// GET /api/v1/profile - Obtener perfil completo
		profileGroup.GET("", profileHdlr.GetProfile)

		// PUT /api/v1/profile - Actualizar información personal
		profileGroup.PUT("", profileHdlr.UpdateProfile)

		// POST /api/v1/profile/photo - Subir foto de perfil
		profileGroup.POST("/photo", profileHdlr.UploadProfilePhoto)

		// POST /api/v1/profile/iban - Configurar IBAN (requiere cedula_verified)
		profileGroup.POST("/iban",
			authMiddleware.RequireMinKYC("cedula_verified"),
			profileHdlr.ConfigureIBAN,
		)

		// POST /api/v1/profile/kyc/:document_type - Subir documento KYC
		// Parámetros: cedula_front, cedula_back, selfie
		profileGroup.POST("/kyc/:document_type", profileHdlr.UploadKYCDocument)
	}
}

// loadPagaditoConfig carga la configuración de Pagadito desde la base de datos
func loadPagaditoConfig(repo *db.PostgresPaymentProcessorRepository, log *logger.Logger) (*pagadito.Config, error) {
	// Buscar procesador Pagadito activo (sandbox o producción)
	processor, err := repo.FindByProvider("pagadito", true) // true = sandbox
	if err != nil {
		return nil, fmt.Errorf("error cargando configuración de Pagadito: %w", err)
	}

	if !processor.IsActive {
		return nil, fmt.Errorf("Pagadito está deshabilitado en la configuración")
	}

	// Deserializar configuración del JSONB
	var configMap map[string]interface{}
	if err := json.Unmarshal(processor.Config, &configMap); err != nil {
		return nil, fmt.Errorf("error deserializando configuración: %w", err)
	}

	uid, ok := configMap["uid"].(string)
	if !ok {
		return nil, fmt.Errorf("UID no encontrado en configuración")
	}

	wsk, ok := configMap["wsk"].(string)
	if !ok {
		return nil, fmt.Errorf("WSK no encontrado en configuración")
	}

	apiURL, ok := configMap["api_url"].(string)
	if !ok {
		return nil, fmt.Errorf("API URL no encontrado en configuración")
	}

	callbackURL, ok := configMap["callback_url"].(string)
	if !ok {
		return nil, fmt.Errorf("Callback URL no encontrado en configuración")
	}

	sandboxMode, _ := configMap["sandbox_mode"].(bool)

	return &pagadito.Config{
		UID:         uid,
		WSK:         wsk,
		SandboxMode: sandboxMode,
		APIURL:      apiURL,
		ReturnURL:   callbackURL,
	}, nil
}

// setupCreditsRoutes configura las rutas de compra de créditos con Pagadito
func setupCreditsRoutes(router *gin.Engine, gormDB *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) {
	// Inicializar repositorios
	creditPurchaseRepo := db.NewCreditPurchaseRepository(gormDB, log)
	userRepo := db.NewUserRepository(gormDB)
	walletRepo := db.NewWalletRepository(gormDB, log)
	walletTransactionRepo := db.NewWalletTransactionRepository(gormDB, log)
	auditRepo := db.NewAuditLogRepository(gormDB)
	paymentProcessorRepo := db.NewPaymentProcessorRepository(gormDB, log)

	// Cargar configuración de Pagadito desde base de datos
	pagaditoConfig, err := loadPagaditoConfig(paymentProcessorRepo, log)
	if err != nil {
		log.Warn("No se pudo cargar configuración de Pagadito, rutas de créditos deshabilitadas",
			logger.Error(err))
		return
	}

	// Inicializar cliente Pagadito
	pagaditoClient := pagadito.NewHTTPClient(pagaditoConfig)

	// Inicializar token manager y middlewares
	tokenMgr := redisAdapter.NewTokenManager(rdb, &cfg.JWT)
	blacklistService := redisinfra.NewTokenBlacklistService(rdb)
	authMiddleware := middleware.NewAuthMiddleware(tokenMgr, blacklistService, log)
	rateLimiter := middleware.NewRateLimiter(rdb, log)

	// Inicializar use cases de wallet (necesarios para acreditar fondos)
	addFundsUC := walletuc.NewAddFundsUseCase(
		walletRepo,
		walletTransactionRepo,
		userRepo,
		auditRepo,
		log,
	)

	// Inicializar use cases de créditos
	purchaseCreditsUC := creditsuc.NewPurchaseCreditsUseCase(
		creditPurchaseRepo,
		walletRepo,
		userRepo,
		auditRepo,
		pagaditoClient,
		log,
	)

	processCallbackUC := creditsuc.NewProcessPagaditoCallbackUseCase(
		creditPurchaseRepo,
		walletRepo,
		walletTransactionRepo,
		auditRepo,
		pagaditoClient,
		addFundsUC,
		log,
	)

	// Inicializar handlers
	purchaseCreditsHandler := handlers.NewPurchaseCreditsHandler(purchaseCreditsUC, log)
	pagaditoCallbackHandler := handlers.NewPagaditoCallbackHandler(processCallbackUC, log)
	getPurchaseStatusHandler := handlers.NewGetPurchaseStatusHandler(log)

	// Grupo de rutas de créditos
	creditsGroup := router.Group("/api/v1/credits")
	{
		// POST /api/v1/credits/purchase - Comprar créditos (requiere autenticación)
		creditsGroup.POST("/purchase",
			authMiddleware.Authenticate(),
			authMiddleware.RequireMinKYC("email_verified"),
			rateLimiter.LimitByUser(20, time.Hour), // Max 20 compras por hora
			purchaseCreditsHandler.Handle,
		)

		// GET /api/v1/credits/callback - Callback de Pagadito (PÚBLICO, sin auth)
		creditsGroup.GET("/callback", pagaditoCallbackHandler.Handle)

		// GET /api/v1/credits/purchase/:id - Obtener estado de compra (requiere auth)
		creditsGroup.GET("/purchase/:id",
			authMiddleware.Authenticate(),
			getPurchaseStatusHandler.Handle,
		)
	}

	log.Info("Rutas de créditos con Pagadito configuradas correctamente",
		logger.String("callback_url", pagaditoConfig.ReturnURL))
}
