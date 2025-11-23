package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/auth"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// RegisterHandler maneja el endpoint de registro
type RegisterHandler struct {
	useCase *auth.RegisterUseCase
	logger  *logger.Logger
}

// NewRegisterHandler crea una nueva instancia del handler
func NewRegisterHandler(useCase *auth.RegisterUseCase, logger *logger.Logger) *RegisterHandler {
	return &RegisterHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Handle maneja la petición de registro
func (h *RegisterHandler) Handle(c *gin.Context) {
	var input auth.RegisterInput

	// Parsear JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	// DEBUG: Log del input recibido
	h.logger.Info("Register input received",
		logger.String("email", input.Email),
		logger.Bool("accepted_terms", input.AcceptedTerms),
		logger.Bool("accepted_privacy", input.AcceptedPrivacy),
	)

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
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": output.Message,
		"data": gin.H{
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
				"created_at":     output.User.CreatedAt,
			},
			"verification_code_sent": output.VerificationCodeSent,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *RegisterHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in register handler", logger.Error(err))
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
