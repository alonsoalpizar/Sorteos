package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Inicializar logger
	log, err := logger.New(cfg.Server.Environment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Sorteos Platform API",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	// Conectar a PostgreSQL
	db, err := initDatabase(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	log.Info("Connected to PostgreSQL",
		zap.String("host", cfg.Database.Host),
		zap.String("database", cfg.Database.DBName),
	)

	// Conectar a Redis
	rdb := initRedis(cfg, log)
	log.Info("Connected to Redis",
		zap.String("host", cfg.Redis.Host),
		zap.Int("db", cfg.Redis.DB),
	)

	// Configurar modo Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar router
	router := gin.New()
	setupMiddleware(router, log, cfg)
	setupRoutes(router, db, rdb, cfg, log)

	// Crear servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Info("Server listening", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	// Cerrar conexiones
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Error("Error closing database", zap.Error(err))
	}

	if err := rdb.Close(); err != nil {
		log.Error("Error closing Redis", zap.Error(err))
	}

	log.Info("Server exited")
}

// initDatabase inicializa la conexión a PostgreSQL
func initDatabase(cfg *config.Config, log *logger.Logger) (*gorm.DB, error) {
	// Configurar logger de GORM
	var gormLogLevel gormlogger.LogLevel
	if cfg.IsDevelopment() {
		gormLogLevel = gormlogger.Info
	} else {
		gormLogLevel = gormlogger.Error
	}

	gormLog := gormlogger.New(
		&gormLogWriter{log: log},
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  cfg.IsDevelopment(),
		},
	)

	// Conectar a PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormLog,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configurar connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Verificar conexión
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// gormLogWriter adapta logger.Logger para GORM
type gormLogWriter struct {
	log *logger.Logger
}

func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	w.log.Info(fmt.Sprintf(format, args...))
}

// initRedis inicializa la conexión a Redis
func initRedis(cfg *config.Config, log *logger.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// Verificar conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	return rdb
}

// setupMiddleware configura los middlewares globales
func setupMiddleware(router *gin.Engine, log *logger.Logger, cfg *config.Config) {
	// Recovery middleware
	router.Use(gin.Recovery())

	// Logger middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		if len(c.Errors) > 0 {
			log.Error("Request failed", fields...)
		} else if statusCode >= 500 {
			log.Error("Server error", fields...)
		} else if statusCode >= 400 {
			log.Warn("Client error", fields...)
		} else {
			log.Info("Request completed", fields...)
		}
	})

	// CORS middleware
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Verificar si el origen está permitido
		allowed := false
		for _, allowedOrigin := range cfg.Server.AllowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, Idempotency-Key")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Request ID middleware
	router.Use(func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	})
}

// setupRoutes configura las rutas de la API
func setupRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client, cfg *config.Config, log *logger.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// Readiness check (verifica dependencias)
	router.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Verificar PostgreSQL
		sqlDB, _ := db.DB()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"error":  "database connection failed",
			})
			return
		}

		// Verificar Redis
		if err := rdb.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"error":  "redis connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ready",
			"database": "connected",
			"redis":    "connected",
		})
	})

	// Setup auth routes
	setupAuthRoutes(router, db, rdb, cfg, log)

	// API v1 - Ruta de prueba
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
				"timestamp": time.Now().UTC(),
			})
		})
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "NOT_FOUND",
			"message": "Endpoint not found",
		})
	})
}
