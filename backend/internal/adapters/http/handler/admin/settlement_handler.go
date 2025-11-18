package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/internal/usecase/admin/settlement"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// SettlementHandler maneja las peticiones HTTP relacionadas con liquidaciones de organizadores
type SettlementHandler struct {
	listSettlementsUC        *settlement.ListSettlementsUseCase
	viewSettlementDetailsUC  *settlement.ViewSettlementDetailsUseCase
	createSettlementUC       *settlement.CreateSettlementUseCase
	approveSettlementUC      *settlement.ApproveSettlementUseCase
	rejectSettlementUC       *settlement.RejectSettlementUseCase
	markSettlementPaidUC     *settlement.MarkSettlementPaidUseCase
	autoCreateSettlementsUC  *settlement.AutoCreateSettlementsUseCase
	log                      *logger.Logger
}

// NewSettlementHandler crea una nueva instancia del handler
func NewSettlementHandler(db *gorm.DB, log *logger.Logger) *SettlementHandler {
	return &SettlementHandler{
		listSettlementsUC:        settlement.NewListSettlementsUseCase(db, log),
		viewSettlementDetailsUC:  settlement.NewViewSettlementDetailsUseCase(db, log),
		createSettlementUC:       settlement.NewCreateSettlementUseCase(db, log),
		approveSettlementUC:      settlement.NewApproveSettlementUseCase(db, log),
		rejectSettlementUC:       settlement.NewRejectSettlementUseCase(db, log),
		markSettlementPaidUC:     settlement.NewMarkSettlementPaidUseCase(db, log),
		autoCreateSettlementsUC:  settlement.NewAutoCreateSettlementsUseCase(db, log),
		log:                      log,
	}
}

// List lista liquidaciones con filtros y paginación
// GET /api/v1/admin/settlements
func (h *SettlementHandler) List(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &settlement.ListSettlementsInput{
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
	if status := c.Query("status"); status != "" {
		input.Status = &status
	}

	if organizerIDStr := c.Query("organizer_id"); organizerIDStr != "" {
		if organizerID, err := strconv.ParseInt(organizerIDStr, 10, 64); err == nil {
			input.OrganizerID = &organizerID
		}
	}

	if raffleIDStr := c.Query("raffle_id"); raffleIDStr != "" {
		if raffleID, err := strconv.ParseInt(raffleIDStr, 10, 64); err == nil {
			input.RaffleID = &raffleID
		}
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = &dateTo
	}

	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		if minAmount, err := strconv.ParseFloat(minAmountStr, 64); err == nil {
			input.MinAmount = &minAmount
		}
	}

	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		if maxAmount, err := strconv.ParseFloat(maxAmountStr, 64); err == nil {
			input.MaxAmount = &maxAmount
		}
	}

	if kycLevelStr := c.Query("kyc_level"); kycLevelStr != "" {
		kycLevel := domain.KYCLevel(kycLevelStr)
		input.KYCLevel = &kycLevel
	}

	if pendingOnlyStr := c.Query("pending_only"); pendingOnlyStr == "true" {
		input.PendingOnly = true
	}

	// Ejecutar use case
	output, err := h.listSettlementsUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// GetByID obtiene detalles completos de una liquidación
// GET /api/v1/admin/settlements/:id
func (h *SettlementHandler) GetByID(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse settlement ID
	settlementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_SETTLEMENT_ID",
				"message": "invalid settlement ID",
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.viewSettlementDetailsUC.Execute(c.Request.Context(), settlementID, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// Create crea nuevas liquidaciones (individual o batch)
// POST /api/v1/admin/settlements
func (h *SettlementHandler) Create(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body
	var input settlement.CreateSettlementInput
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
	output, err := h.createSettlementUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    output,
	})
}

// Approve aprueba una liquidación
// PUT /api/v1/admin/settlements/:id/approve
func (h *SettlementHandler) Approve(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse settlement ID
	settlementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_SETTLEMENT_ID",
				"message": "invalid settlement ID",
			},
		})
		return
	}

	// Parse body (notas opcionales)
	var body struct {
		Notes string `json:"notes"`
	}
	c.ShouldBindJSON(&body)

	input := &settlement.ApproveSettlementInput{
		SettlementID: settlementID,
		Notes:        body.Notes,
	}

	// Ejecutar use case
	output, err := h.approveSettlementUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// Reject rechaza una liquidación
// PUT /api/v1/admin/settlements/:id/reject
func (h *SettlementHandler) Reject(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse settlement ID
	settlementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_SETTLEMENT_ID",
				"message": "invalid settlement ID",
			},
		})
		return
	}

	// Parse body (razón obligatoria)
	var body struct {
		Reason string `json:"reason" binding:"required"`
		Notes  string `json:"notes"`
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

	input := &settlement.RejectSettlementInput{
		SettlementID: settlementID,
		Reason:       body.Reason,
		Notes:        body.Notes,
	}

	// Ejecutar use case
	output, err := h.rejectSettlementUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// MarkPaid marca una liquidación como pagada
// PUT /api/v1/admin/settlements/:id/payout
func (h *SettlementHandler) MarkPaid(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse settlement ID
	settlementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_SETTLEMENT_ID",
				"message": "invalid settlement ID",
			},
		})
		return
	}

	// Parse body
	var body settlement.MarkSettlementPaidInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Asignar settlement ID desde params
	body.SettlementID = settlementID

	// Ejecutar use case
	output, err := h.markSettlementPaidUC.Execute(c.Request.Context(), &body, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// AutoCreate crea liquidaciones automáticamente para rifas completadas
// POST /api/v1/admin/settlements/auto-create
func (h *SettlementHandler) AutoCreate(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body (parametros opcionales)
	var input settlement.AutoCreateSettlementsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Valores por defecto si no se envía body
		input.DaysAfterCompletion = 3
		input.DryRun = false
	}

	// Ejecutar use case
	output, err := h.autoCreateSettlementsUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
