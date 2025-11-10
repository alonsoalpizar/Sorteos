package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/sorteos-platform/backend/pkg/logger"
)

// RateLimiter maneja el rate limiting con Redis
type RateLimiter struct {
	rdb    *redis.Client
	logger *logger.Logger
}

// NewRateLimiter crea una nueva instancia del rate limiter
func NewRateLimiter(rdb *redis.Client, logger *logger.Logger) *RateLimiter {
	return &RateLimiter{
		rdb:    rdb,
		logger: logger,
	}
}

// LimitByIP limita las peticiones por dirección IP
func (rl *RateLimiter) LimitByIP(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:ip:%s", ip)

		allowed, err := rl.checkRateLimit(c.Request.Context(), key, maxRequests, window)
		if err != nil {
			rl.logger.Error("Error checking rate limit", logger.Error(err))
			// En caso de error de Redis, permitir la petición
			c.Next()
			return
		}

		if !allowed {
			rl.logger.Warn("Rate limit exceeded",
				logger.String("ip", ip),
				logger.String("endpoint", c.Request.URL.Path),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "TOO_MANY_REQUESTS",
				"message": "Demasiadas solicitudes. Por favor intente más tarde.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LimitByUser limita las peticiones por usuario autenticado
func (rl *RateLimiter) LimitByUser(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Intentar obtener user_id del contexto (requiere AuthMiddleware antes)
		userID, exists := GetUserID(c)
		if !exists {
			// Si no hay usuario autenticado, usar IP
			ip := c.ClientIP()
			key := fmt.Sprintf("rate_limit:ip:%s", ip)

			allowed, err := rl.checkRateLimit(c.Request.Context(), key, maxRequests, window)
			if err != nil {
				rl.logger.Error("Error checking rate limit", logger.Error(err))
				c.Next()
				return
			}

			if !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code":    "TOO_MANY_REQUESTS",
					"message": "Demasiadas solicitudes. Por favor intente más tarde.",
				})
				c.Abort()
				return
			}

			c.Next()
			return
		}

		// Rate limit por usuario
		key := fmt.Sprintf("rate_limit:user:%d", userID)

		allowed, err := rl.checkRateLimit(c.Request.Context(), key, maxRequests, window)
		if err != nil {
			rl.logger.Error("Error checking rate limit", logger.Error(err))
			c.Next()
			return
		}

		if !allowed {
			rl.logger.Warn("Rate limit exceeded",
				logger.Int64("user_id", userID),
				logger.String("endpoint", c.Request.URL.Path),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "TOO_MANY_REQUESTS",
				"message": "Demasiadas solicitudes. Por favor intente más tarde.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LimitByEndpoint limita las peticiones por endpoint específico y usuario
func (rl *RateLimiter) LimitByEndpoint(endpoint string, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		// Intentar obtener user_id del contexto
		userID, exists := GetUserID(c)
		if exists {
			key = fmt.Sprintf("rate_limit:endpoint:%s:user:%d", endpoint, userID)
		} else {
			ip := c.ClientIP()
			key = fmt.Sprintf("rate_limit:endpoint:%s:ip:%s", endpoint, ip)
		}

		allowed, err := rl.checkRateLimit(c.Request.Context(), key, maxRequests, window)
		if err != nil {
			rl.logger.Error("Error checking rate limit", logger.Error(err))
			c.Next()
			return
		}

		if !allowed {
			rl.logger.Warn("Rate limit exceeded for endpoint",
				logger.String("endpoint", endpoint),
				logger.String("key", key),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "TOO_MANY_REQUESTS",
				"message": fmt.Sprintf("Demasiadas solicitudes para %s. Por favor intente más tarde.", endpoint),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit verifica el rate limit usando sliding window en Redis
func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string, maxRequests int, window time.Duration) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-window)

	pipe := rl.rdb.Pipeline()

	// Remover entradas antiguas fuera de la ventana
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))

	// Contar peticiones en la ventana actual
	pipe.ZCard(ctx, key)

	// Agregar nueva petición
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixNano()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})

	// Establecer expiración de la key
	pipe.Expire(ctx, key, window)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Obtener el conteo (segundo comando)
	count := cmds[1].(*redis.IntCmd).Val()

	// Permitir si está por debajo del límite
	return count < int64(maxRequests), nil
}

// GetRemainingRequests retorna las peticiones restantes para una key
func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, key string, maxRequests int, window time.Duration) (int, error) {
	windowStart := time.Now().Add(-window)

	// Remover entradas antiguas
	rl.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))

	// Contar peticiones actuales
	count, err := rl.rdb.ZCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	remaining := maxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// ResetRateLimit reinicia el contador de rate limit para una key
func (rl *RateLimiter) ResetRateLimit(ctx context.Context, key string) error {
	return rl.rdb.Del(ctx, key).Err()
}
