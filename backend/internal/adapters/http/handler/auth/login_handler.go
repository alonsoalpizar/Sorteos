package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/auth"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// LoginHandler maneja el endpoint de login
type LoginHandler struct {
	useCase *auth.LoginUseCase
	logger  *logger.Logger
}

// NewLoginHandler crea una nueva instancia del handler
func NewLoginHandler(useCase *auth.LoginUseCase, logger *logger.Logger) *LoginHandler {
	return &LoginHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle maneja la petición de login
func (h *LoginHandler) Handle(c *gin.Context) {
	var input auth.LoginInput

	// Parsear JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	// Obtener IP y User-Agent
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Ejecutar caso de uso
	output, err := h.useCase.Execute(c.Request.Context(), &input, ip, userAgent)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login exitoso",
		"data": gin.H{
			"user": gin.H{
				"id":            output.User.ID,
				"uuid":          output.User.UUID,
				"email":         output.User.Email,
				"email_verified": output.User.EmailVerified,
				"phone_verified": output.User.PhoneVerified,
				"kyc_level":     output.User.KYCLevel,
				"role":          output.User.Role,
				"status":        output.User.Status,
				"first_name":    output.User.FirstName,
				"last_name":     output.User.LastName,
			},
			"access_token":  output.AccessToken,
			"refresh_token": output.RefreshToken,
			"token_type":    output.TokenType,
			"expires_in":    output.ExpiresIn,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *LoginHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in login handler", logger.Error(err))
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
