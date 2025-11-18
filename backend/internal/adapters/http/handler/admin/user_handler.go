package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/internal/usecase/admin/user"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UserHandler maneja las peticiones HTTP relacionadas con administración de usuarios
type UserHandler struct {
	listUsersUC        *user.ListUsersUseCase
	getUserDetailUC    *user.GetUserDetailUseCase
	updateUserStatusUC *user.UpdateUserStatusUseCase
	updateUserKYCUC    *user.UpdateUserKYCUseCase
	deleteUserUC       *user.DeleteUserUseCase
	log                *logger.Logger
}

// NewUserHandler crea una nueva instancia del handler
func NewUserHandler(db *gorm.DB, log *logger.Logger) *UserHandler {
	return &UserHandler{
		listUsersUC:        user.NewListUsersUseCase(db, log),
		getUserDetailUC:    user.NewGetUserDetailUseCase(db, log),
		updateUserStatusUC: user.NewUpdateUserStatusUseCase(db, log),
		updateUserKYCUC:    user.NewUpdateUserKYCUseCase(db, log),
		deleteUserUC:       user.NewDeleteUserUseCase(db, log),
		log:                log,
	}
}

// List lista usuarios con filtros y paginación
// GET /api/v1/admin/users
func (h *UserHandler) List(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &user.ListUsersInput{
		Page:     1,
		PageSize: 20,
		Search:   c.Query("search"),
		OrderBy:  c.Query("order_by"),
	}

	// Parse page y page_size
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			input.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			input.PageSize = pageSize
		}
	}

	// Parse filtros opcionales
	if roleStr := c.Query("role"); roleStr != "" {
		role := domain.UserRole(roleStr)
		input.Role = &role
	}

	if status := c.Query("status"); status != "" {
		input.Status = &status
	}

	if kycLevelStr := c.Query("kyc_level"); kycLevelStr != "" {
		kycLevel := domain.KYCLevel(kycLevelStr)
		input.KYCLevel = &kycLevel
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = &dateTo
	}

	// Ejecutar use case
	output, err := h.listUsersUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// GetByID obtiene detalles completos de un usuario
// GET /api/v1/admin/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse user ID
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.getUserDetailUC.Execute(c.Request.Context(), userID, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// UpdateStatus actualiza el estado de un usuario (suspend, activate, ban, unban)
// PUT /api/v1/admin/users/:id/status
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse user ID
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	// Parse body
	var body struct {
		Action string `json:"action" binding:"required"` // suspend, activate, ban, unban
		Reason string `json:"reason"`                    // Required for suspend and ban
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	input := &user.UpdateUserStatusInput{
		UserID: userID,
		Action: user.UserStatusAction(body.Action),
		Reason: body.Reason,
	}

	// Ejecutar use case
	err = h.updateUserStatusUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User status updated successfully",
	})
}

// UpdateKYC actualiza el nivel KYC de un usuario
// PUT /api/v1/admin/users/:id/kyc
func (h *UserHandler) UpdateKYC(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse user ID
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	// Parse body
	var body struct {
		KYCLevel string `json:"kyc_level" binding:"required"` // none, email_verified, phone_verified, cedula_verified, full_kyc
		Notes    string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	input := &user.UpdateUserKYCInput{
		UserID:   userID,
		KYCLevel: domain.KYCLevel(body.KYCLevel),
		Notes:    body.Notes,
	}

	// Ejecutar use case
	err = h.updateUserKYCUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User KYC level updated successfully",
	})
}

// Delete elimina un usuario (soft delete)
// DELETE /api/v1/admin/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse user ID
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	// Parse body (razón requerida)
	var body struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": "reason is required",
			},
		})
		return
	}

	input := &user.DeleteUserInput{
		UserID: userID,
		Reason: body.Reason,
	}

	// Ejecutar use case
	err = h.deleteUserUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}
