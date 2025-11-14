package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/auth"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// LogoutHandler maneja el endpoint de logout
type LogoutHandler struct {
	useCase *auth.LogoutUseCase
	logger  *logger.Logger
}

// NewLogoutHandler crea una nueva instancia del handler
func NewLogoutHandler(useCase *auth.LogoutUseCase, logger *logger.Logger) *LogoutHandler {
	return &LogoutHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// LogoutRequest representa la petición de logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Handle maneja la petición de logout
func (h *LogoutHandler) Handle(c *gin.Context) {
	var req LogoutRequest

	// Parsear JSON (refresh token es opcional en el body)
	if err := c.ShouldBindJSON(&req); err != nil {
		// Es válido no enviar body, continuamos
		req.RefreshToken = ""
	}

	// Extraer access token del header Authorization
	accessToken := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			accessToken = parts[1]
		}
	}

	// Preparar input para el caso de uso
	input := &auth.LogoutInput{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	}

	// Ejecutar caso de uso
	err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout exitoso",
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *LogoutHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in logout handler", logger.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Error interno del servidor",
		})
		return
	}

	c.JSON(appErr.Status, ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
	})
}
