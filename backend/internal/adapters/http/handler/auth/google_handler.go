package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/auth"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GoogleAuthHandler maneja el endpoint de autenticación con Google
type GoogleAuthHandler struct {
	useCase *auth.GoogleAuthUseCase
	logger  *logger.Logger
}

// NewGoogleAuthHandler crea una nueva instancia del handler
func NewGoogleAuthHandler(useCase *auth.GoogleAuthUseCase, logger *logger.Logger) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// GoogleAuthRequest representa la solicitud de autenticación con Google
type GoogleAuthRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// Handle maneja la petición de autenticación con Google
func (h *GoogleAuthHandler) Handle(c *gin.Context) {
	var req GoogleAuthRequest

	// Parsear JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "El token de Google es requerido",
		})
		return
	}

	// Obtener IP y User-Agent
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Ejecutar caso de uso
	input := &auth.GoogleAuthInput{
		IDToken: req.IDToken,
	}

	output, err := h.useCase.Execute(c.Request.Context(), input, ip, userAgent)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Si requiere vinculación (email ya existe)
	if output.RequiresLinking {
		c.JSON(http.StatusConflict, gin.H{
			"success":          false,
			"requires_linking": true,
			"message":          "Ya existe una cuenta con este email. Ingresa tu contraseña para vincular tu cuenta de Google.",
			"email":            output.ExistingEmail,
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Autenticación con Google exitosa",
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
				"auth_provider":  output.User.AuthProvider,
			},
			"access_token":  output.AccessToken,
			"refresh_token": output.RefreshToken,
			"token_type":    output.TokenType,
			"expires_in":    output.ExpiresIn,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *GoogleAuthHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in google auth handler", logger.Error(err))
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

// GoogleLinkHandler maneja el endpoint de vinculación de cuenta Google
type GoogleLinkHandler struct {
	useCase *auth.GoogleLinkUseCase
	logger  *logger.Logger
}

// NewGoogleLinkHandler crea una nueva instancia del handler
func NewGoogleLinkHandler(useCase *auth.GoogleLinkUseCase, logger *logger.Logger) *GoogleLinkHandler {
	return &GoogleLinkHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// GoogleLinkRequest representa la solicitud de vinculación
type GoogleLinkRequest struct {
	IDToken  string `json:"id_token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Handle maneja la petición de vinculación de cuenta Google
func (h *GoogleLinkHandler) Handle(c *gin.Context) {
	var req GoogleLinkRequest

	// Parsear JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "El token de Google y la contraseña son requeridos",
		})
		return
	}

	// Obtener IP y User-Agent
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Ejecutar caso de uso
	input := &auth.GoogleLinkInput{
		IDToken:  req.IDToken,
		Password: req.Password,
	}

	output, err := h.useCase.Execute(c.Request.Context(), input, ip, userAgent)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cuenta de Google vinculada exitosamente",
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
				"auth_provider":  output.User.AuthProvider,
			},
			"access_token":  output.AccessToken,
			"refresh_token": output.RefreshToken,
			"token_type":    output.TokenType,
			"expires_in":    output.ExpiresIn,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *GoogleLinkHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in google link handler", logger.Error(err))
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
