package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/adapters/redis"
	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// AuthMiddleware maneja la autenticación JWT
type AuthMiddleware struct {
	tokenMgr         *redis.TokenManager
	blacklistService *redisinfra.TokenBlacklistService
	logger           *logger.Logger
}

// NewAuthMiddleware crea una nueva instancia del middleware
func NewAuthMiddleware(tokenMgr *redis.TokenManager, blacklistService *redisinfra.TokenBlacklistService, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		tokenMgr:         tokenMgr,
		blacklistService: blacklistService,
		logger:           logger,
	}
}

// Authenticate verifica el JWT token
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.respondUnauthorized(c, "Token de autorización no proporcionado")
			return
		}

		// Verificar formato "Bearer {token}"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.respondUnauthorized(c, "Formato de token inválido")
			return
		}

		tokenString := parts[1]

		// Verificar si el token está en blacklist ANTES de validarlo
		// Esto es más eficiente y evita procesamiento innecesario
		blacklisted, err := m.blacklistService.IsBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			m.logger.Error("Error checking token blacklist", logger.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Error verificando token",
			})
			c.Abort()
			return
		}

		if blacklisted {
			m.logger.Warn("Blacklisted token used")
			m.respondUnauthorized(c, "Token revocado")
			return
		}

		// Validar token
		claims, err := m.tokenMgr.ValidateAccessToken(tokenString)
		if err != nil {
			m.logger.Warn("Invalid token", logger.Error(err))
			m.respondUnauthorized(c, "Token inválido o expirado")
			return
		}

		// Guardar claims en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_kyc_level", claims.KYCLevel)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole verifica que el usuario tenga un rol específico
func (m *AuthMiddleware) RequireRole(allowedRoles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// El middleware Authenticate debe ejecutarse primero
		role, exists := c.Get("user_role")
		if !exists {
			m.respondForbidden(c, "Rol de usuario no encontrado")
			return
		}

		userRole := role.(domain.UserRole)

		// Verificar si el rol está permitido
		allowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			m.logger.Warn("Forbidden access attempt",
				logger.String("user_role", string(userRole)),
				logger.String("required_roles", formatRoles(allowedRoles)),
			)
			m.respondForbidden(c, "No tienes permisos para esta operación")
			return
		}

		c.Next()
	}
}

// RequireMinKYC verifica que el usuario tenga un nivel mínimo de KYC
func (m *AuthMiddleware) RequireMinKYC(minLevel domain.KYCLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		// El middleware Authenticate debe ejecutarse primero
		kycLevel, exists := c.Get("user_kyc_level")
		if !exists {
			m.respondForbidden(c, "Nivel KYC no encontrado")
			return
		}

		userKYCLevel := kycLevel.(domain.KYCLevel)

		// Niveles de KYC (orden ascendente)
		levels := map[domain.KYCLevel]int{
			domain.KYCLevelNone:            0,
			domain.KYCLevelEmailVerified:   1,
			domain.KYCLevelPhoneVerified:   2,
			domain.KYCLevelCedulaVerified:  3,
			domain.KYCLevelFullKYC:         4,
		}

		if levels[userKYCLevel] < levels[minLevel] {
			m.logger.Warn("Insufficient KYC level",
				logger.String("user_kyc", string(userKYCLevel)),
				logger.String("required_kyc", string(minLevel)),
			)
			m.respondForbidden(c, "Nivel de verificación insuficiente")
			return
		}

		c.Next()
	}
}

// OptionalAuth intenta autenticar pero no falla si no hay token
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := m.tokenMgr.ValidateAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Guardar claims en el contexto si el token es válido
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_kyc_level", claims.KYCLevel)
		c.Set("claims", claims)

		c.Next()
	}
}

// respondUnauthorized envía una respuesta 401
func (m *AuthMiddleware) respondUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":    errors.ErrUnauthorized.Code,
		"message": message,
	})
	c.Abort()
}

// respondForbidden envía una respuesta 403
func (m *AuthMiddleware) respondForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, gin.H{
		"code":    errors.ErrForbidden.Code,
		"message": message,
	})
	c.Abort()
}

// formatRoles formatea un slice de roles como string
func formatRoles(roles []domain.UserRole) string {
	strs := make([]string, len(roles))
	for i, role := range roles {
		strs[i] = string(role)
	}
	return strings.Join(strs, ", ")
}

// GetUserID obtiene el user_id del contexto
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int64), true
}

// GetUserRole obtiene el role del contexto
func GetUserRole(c *gin.Context) (domain.UserRole, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(domain.UserRole), true
}

// GetClaims obtiene los claims completos del contexto
func GetClaims(c *gin.Context) (*redis.Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	return claims.(*redis.Claims), true
}
