package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/notifications"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// NotificationHandler maneja las peticiones HTTP relacionadas con notificaciones
type NotificationHandler struct {
	sendEmailUC            *notifications.SendEmailUseCase
	sendBulkEmailUC        *notifications.SendBulkEmailUseCase
	manageTemplatesUC      *notifications.ManageEmailTemplatesUseCase
	createAnnouncementUC   *notifications.CreateAnnouncementUseCase
	viewHistoryUC          *notifications.ViewNotificationHistoryUseCase
	log                    *logger.Logger
}

// NewNotificationHandler crea una nueva instancia del handler
func NewNotificationHandler(db *gorm.DB, log *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		sendEmailUC:            notifications.NewSendEmailUseCase(db, log),
		sendBulkEmailUC:        notifications.NewSendBulkEmailUseCase(db, log),
		manageTemplatesUC:      notifications.NewManageEmailTemplatesUseCase(db, log),
		createAnnouncementUC:   notifications.NewCreateAnnouncementUseCase(db, log),
		viewHistoryUC:          notifications.NewViewNotificationHistoryUseCase(db, log),
		log:                    log,
	}
}

// SendEmail envía un email individual o a múltiples destinatarios
// POST /api/v1/admin/notifications/email
func (h *NotificationHandler) SendEmail(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body
	var input notifications.SendEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.sendEmailUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// SendBulkEmail envía emails masivos (newsletters, announcements, campaigns)
// POST /api/v1/admin/notifications/bulk
func (h *NotificationHandler) SendBulkEmail(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body
	var input notifications.SendBulkEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.sendBulkEmailUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ManageTemplates gestiona plantillas de email (list, create, update, delete)
// GET/POST/PUT/DELETE /api/v1/admin/notifications/templates
func (h *NotificationHandler) ManageTemplates(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body
	var input notifications.ManageEmailTemplatesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.manageTemplatesUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// CreateAnnouncement crea anuncios masivos para todos los usuarios o segmentos
// POST /api/v1/admin/notifications/announcements
func (h *NotificationHandler) CreateAnnouncement(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body
	var input notifications.CreateAnnouncementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.createAnnouncementUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    output,
	})
}

// ViewHistory ver historial de notificaciones enviadas con filtros
// GET /api/v1/admin/notifications/history
func (h *NotificationHandler) ViewHistory(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &notifications.ViewNotificationHistoryInput{
		Limit:  20,
		Offset: 0,
	}

	// Parse limit y offset
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			input.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			input.Offset = offset
		}
	}

	// Parse search
	if search := c.Query("search"); search != "" {
		input.Search = &search
	}

	// Parse filtros opcionales
	if notifType := c.Query("type"); notifType != "" {
		input.Type = &notifType
	}

	if status := c.Query("status"); status != "" {
		input.Status = &status
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = &dateTo
	}

	// Ejecutar use case
	output, err := h.viewHistoryUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
