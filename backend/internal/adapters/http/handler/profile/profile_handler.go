package profile

import (
	"net/http"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/internal/usecase/profile"

	"github.com/gin-gonic/gin"
)

// ProfileHandler maneja los endpoints de perfil de usuario
type ProfileHandler struct {
	getProfileUC         *profile.GetProfileUseCase
	updateProfileUC      *profile.UpdateProfileUseCase
	uploadPhotoUC        *profile.UploadProfilePhotoUseCase
	configureIBANUC      *profile.ConfigureIBANUseCase
	uploadKYCDocumentUC  *profile.UploadKYCDocumentUseCase
}

// NewProfileHandler crea una nueva instancia del handler
func NewProfileHandler(
	getProfileUC *profile.GetProfileUseCase,
	updateProfileUC *profile.UpdateProfileUseCase,
	uploadPhotoUC *profile.UploadProfilePhotoUseCase,
	configureIBANUC *profile.ConfigureIBANUseCase,
	uploadKYCDocumentUC *profile.UploadKYCDocumentUseCase,
) *ProfileHandler {
	return &ProfileHandler{
		getProfileUC:        getProfileUC,
		updateProfileUC:     updateProfileUC,
		uploadPhotoUC:       uploadPhotoUC,
		configureIBANUC:     configureIBANUC,
		uploadKYCDocumentUC: uploadKYCDocumentUC,
	}
}

// GetProfile obtiene el perfil completo del usuario
// GET /api/v1/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	result, err := h.getProfileUC.Execute(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateProfile actualiza la información personal del usuario
// PUT /api/v1/profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	var req profile.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request: " + err.Error(),
		})
		return
	}

	user, err := h.updateProfileUC.Execute(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UploadProfilePhoto sube la foto de perfil del usuario
// POST /api/v1/profile/photo
func (h *ProfileHandler) UploadProfilePhoto(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Por ahora, esperamos que el frontend envíe el file_url ya procesado
	// TODO: Implementar upload real de archivo usando multipart/form-data
	var req struct {
		PhotoURL string `json:"photo_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request: " + err.Error(),
		})
		return
	}

	user, err := h.uploadPhotoUC.Execute(c.Request.Context(), userID.(int64), req.PhotoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// ConfigureIBAN configura el IBAN para retiros
// POST /api/v1/profile/iban
func (h *ProfileHandler) ConfigureIBAN(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	var req struct {
		IBAN string `json:"iban" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request: " + err.Error(),
		})
		return
	}

	user, err := h.configureIBANUC.Execute(c.Request.Context(), userID.(int64), req.IBAN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UploadKYCDocument sube un documento KYC
// POST /api/v1/profile/kyc/:document_type
func (h *ProfileHandler) UploadKYCDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	docType := c.Param("document_type")
	if docType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "document_type is required",
		})
		return
	}

	// Convertir string a DocumentType
	var documentType domain.DocumentType
	switch docType {
	case "cedula_front":
		documentType = domain.DocumentTypeCedulaFront
	case "cedula_back":
		documentType = domain.DocumentTypeCedulaBack
	case "selfie":
		documentType = domain.DocumentTypeSelfie
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid document_type. Valid values: cedula_front, cedula_back, selfie",
		})
		return
	}

	// Por ahora, esperamos que el frontend envíe el file_url ya procesado
	// TODO: Implementar upload real de archivo usando multipart/form-data
	var req struct {
		FileURL string `json:"file_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request: " + err.Error(),
		})
		return
	}

	doc, err := h.uploadKYCDocumentUC.Execute(c.Request.Context(), userID.(int64), documentType, req.FileURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    doc,
	})
}
