package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/auth"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// VerifyEmailHandler maneja el endpoint de verificación de email
type VerifyEmailHandler struct {
	useCase *auth.VerifyEmailUseCase
	logger  *logger.Logger
}

// NewVerifyEmailHandler crea una nueva instancia del handler
func NewVerifyEmailHandler(useCase *auth.VerifyEmailUseCase, logger *logger.Logger) *VerifyEmailHandler {
	return &VerifyEmailHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle maneja la petición de verificación de email
func (h *VerifyEmailHandler) Handle(c *gin.Context) {
	var input auth.VerifyEmailInput

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
	response := gin.H{
		"success": output.Success,
		"message": output.Message,
	}

	// Si se generaron tokens (login automático después de verificación)
	if output.AccessToken != "" {
		response["data"] = gin.H{
			"user": gin.H{
				"id":             output.User.ID,
				"uuid":           output.User.UUID,
				"email":          output.User.Email,
				"email_verified": output.User.EmailVerified,
				"phone_verified": output.User.PhoneVerified,
				"kyc_level":      output.User.KYCLevel,
				"role":           output.User.Role,
				"status":         output.User.Status,
				"first_name":     output.User.FirstName,
				"last_name":      output.User.LastName,
			},
			"access_token":  output.AccessToken,
			"refresh_token": output.RefreshToken,
		}
	}

	c.JSON(http.StatusOK, response)
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *VerifyEmailHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in verify email handler", logger.Error(err))
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
