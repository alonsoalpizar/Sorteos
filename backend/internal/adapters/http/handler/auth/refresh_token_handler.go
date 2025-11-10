package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/auth"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// RefreshTokenHandler maneja el endpoint de refresh token
type RefreshTokenHandler struct {
	useCase *auth.RefreshTokenUseCase
	logger  *logger.Logger
}

// NewRefreshTokenHandler crea una nueva instancia del handler
func NewRefreshTokenHandler(useCase *auth.RefreshTokenUseCase, logger *logger.Logger) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle maneja la petición de refresh token
func (h *RefreshTokenHandler) Handle(c *gin.Context) {
	var input auth.RefreshTokenInput

	// Parsear JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	// Ejecutar caso de uso
	output, err := h.useCase.Execute(c.Request.Context(), &input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token renovado exitosamente",
		"data": gin.H{
			"access_token":  output.AccessToken,
			"refresh_token": output.RefreshToken,
			"token_type":    output.TokenType,
			"expires_in":    output.ExpiresIn,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *RefreshTokenHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in refresh token handler", logger.Error(err))
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
