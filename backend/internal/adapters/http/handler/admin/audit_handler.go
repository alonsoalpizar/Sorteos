package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/audit"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// AuditHandler maneja las peticiones HTTP relacionadas con audit logs
type AuditHandler struct {
	listAuditLogsUC *audit.ListAuditLogsUseCase
	log             *logger.Logger
}

// NewAuditHandler crea una nueva instancia del handler
func NewAuditHandler(db *gorm.DB, log *logger.Logger) *AuditHandler {
	return &AuditHandler{
		listAuditLogsUC: audit.NewListAuditLogsUseCase(db, log),
		log:             log,
	}
}

// ListAuditLogs lista audit logs con filtros y paginación
// GET /api/v1/admin/audit
func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &audit.ListAuditLogsInput{
		Page:     1,
		PageSize: 50,
		Search:   c.Query("search"),
		OrderBy:  c.DefaultQuery("order_by", "created_at DESC"),
	}

	// Parse page y page_size
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			input.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			input.PageSize = pageSize
		}
	}

	// Parse filtros opcionales
	if adminIDStr := c.Query("admin_id"); adminIDStr != "" {
		if adminIDFilter, err := strconv.ParseInt(adminIDStr, 10, 64); err == nil {
			input.AdminID = &adminIDFilter
		}
	}

	if action := c.Query("action"); action != "" {
		input.Action = &action
	}

	if entityType := c.Query("entity_type"); entityType != "" {
		input.EntityType = &entityType
	}

	if entityIDStr := c.Query("entity_id"); entityIDStr != "" {
		if entityID, err := strconv.ParseInt(entityIDStr, 10, 64); err == nil {
			input.EntityID = &entityID
		}
	}

	if severity := c.Query("severity"); severity != "" {
		input.Severity = &severity
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = &dateTo
	}

	// Ejecutar use case
	output, err := h.listAuditLogsUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
