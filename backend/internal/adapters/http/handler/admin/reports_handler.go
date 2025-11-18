package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/reports"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ReportsHandler maneja las peticiones HTTP relacionadas con reportes y dashboard
type ReportsHandler struct {
	globalDashboardUC          *reports.GlobalDashboardUseCase
	revenueReportUC            *reports.RevenueReportUseCase
	raffleLiquidationsReportUC *reports.RaffleLiquidationsReportUseCase
	exportDataUC               *reports.ExportDataUseCase
	log                        *logger.Logger
}

// NewReportsHandler crea una nueva instancia del handler
func NewReportsHandler(db *gorm.DB, log *logger.Logger) *ReportsHandler {
	return &ReportsHandler{
		globalDashboardUC:          reports.NewGlobalDashboardUseCase(db, log),
		revenueReportUC:            reports.NewRevenueReportUseCase(db, log),
		raffleLiquidationsReportUC: reports.NewRaffleLiquidationsReportUseCase(db, log),
		exportDataUC:               reports.NewExportDataUseCase(db, log),
		log:                        log,
	}
}

// GetDashboard obtiene KPIs del dashboard global
// GET /api/v1/admin/reports/dashboard
func (h *ReportsHandler) GetDashboard(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Ejecutar use case
	kpis, err := h.globalDashboardUC.Execute(c.Request.Context(), adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    kpis,
	})
}

// GetRevenueReport obtiene reporte de ingresos con filtros y agrupación
// GET /api/v1/admin/reports/revenue
func (h *ReportsHandler) GetRevenueReport(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &reports.RevenueReportInput{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		GroupBy:  c.DefaultQuery("group_by", "day"),
	}

	// Parse filtros opcionales
	if organizerIDStr := c.Query("organizer_id"); organizerIDStr != "" {
		if organizerID, err := strconv.ParseInt(organizerIDStr, 10, 64); err == nil {
			input.OrganizerID = &organizerID
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			input.CategoryID = &categoryID
		}
	}

	// Ejecutar use case
	output, err := h.revenueReportUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// GetLiquidationsReport obtiene reporte de liquidaciones de rifas
// GET /api/v1/admin/reports/liquidations
func (h *ReportsHandler) GetLiquidationsReport(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &reports.RaffleLiquidationsReportInput{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		OrderBy:  c.DefaultQuery("order_by", "raffles.completed_at DESC"),
	}

	// Parse filtros opcionales
	if organizerIDStr := c.Query("organizer_id"); organizerIDStr != "" {
		if organizerID, err := strconv.ParseInt(organizerIDStr, 10, 64); err == nil {
			input.OrganizerID = &organizerID
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			input.CategoryID = &categoryID
		}
	}

	if settlementStatus := c.Query("settlement_status"); settlementStatus != "" {
		input.SettlementStatus = &settlementStatus
	}

	// Ejecutar use case
	output, err := h.raffleLiquidationsReportUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ExportData exporta datos a CSV/Excel/PDF
// POST /api/v1/admin/reports/export
func (h *ReportsHandler) ExportData(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body
	var input reports.ExportDataInput
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
	output, err := h.exportDataUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
