package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/settlement"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// SettlementHandler maneja todas las operaciones de administración de settlements
type SettlementHandler struct {
	listSettlements    *settlement.ListSettlementsUseCase
	approveSettlement  *settlement.ApproveSettlementUseCase
	rejectSettlement   *settlement.RejectSettlementUseCase
	createSettlement   *settlement.CreateSettlementUseCase
	markPaid           *settlement.MarkSettlementPaidUseCase
	autoCreate         *settlement.AutoCreateSettlementsUseCase
}

// NewSettlementHandler crea una nueva instancia
func NewSettlementHandler(
	listSettlements *settlement.ListSettlementsUseCase,
	approveSettlement *settlement.ApproveSettlementUseCase,
	rejectSettlement *settlement.RejectSettlementUseCase,
	createSettlement *settlement.CreateSettlementUseCase,
	markPaid *settlement.MarkSettlementPaidUseCase,
	autoCreate *settlement.AutoCreateSettlementsUseCase,
) *SettlementHandler {
	return &SettlementHandler{
		listSettlements:   listSettlements,
		approveSettlement: approveSettlement,
		rejectSettlement:  rejectSettlement,
		createSettlement:  createSettlement,
		markPaid:          markPaid,
		autoCreate:        autoCreate,
	}
}

// ListSettlements maneja GET /api/v1/admin/settlements
func (h *SettlementHandler) ListSettlements(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parsear query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Construir input
	input := &settlement.ListSettlementsInput{
		Page:     page,
		PageSize: pageSize,
		Search:   stringPtr(c.Query("search")),
		OrderBy:  stringPtr(c.Query("order_by")),
	}

	// Filtros opcionales
	if statusStr := c.Query("status"); statusStr != "" {
		input.Status = stringPtr(statusStr)
	}
	if organizerIDStr := c.Query("organizer_id"); organizerIDStr != "" {
		orgID, err := strconv.ParseInt(organizerIDStr, 10, 64)
		if err == nil {
			input.OrganizerID = &orgID
		}
	}
	if kycLevelStr := c.Query("kyc_level"); kycLevelStr != "" {
		input.KYCLevel = stringPtr(kycLevelStr)
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = stringPtr(dateFrom)
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = stringPtr(dateTo)
	}

	// Ejecutar use case
	output, err := h.listSettlements.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ApproveSettlement maneja POST /api/v1/admin/settlements/:id/approve
func (h *SettlementHandler) ApproveSettlement(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	settlementIDStr := c.Param("id")
	settlementID, err := strconv.ParseInt(settlementIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_SETTLEMENT_ID", "invalid settlement ID format", 400, err))
		return
	}

	var req struct {
		Notes *string `json:"notes"`
	}
	c.ShouldBindJSON(&req)

	input := &settlement.ApproveSettlementInput{
		SettlementID: settlementID,
		Notes:        req.Notes,
	}

	output, err := h.approveSettlement.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// RejectSettlement maneja POST /api/v1/admin/settlements/:id/reject
func (h *SettlementHandler) RejectSettlement(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	settlementIDStr := c.Param("id")
	settlementID, err := strconv.ParseInt(settlementIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_SETTLEMENT_ID", "invalid settlement ID format", 400, err))
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "reason is required", 400, err))
		return
	}

	input := &settlement.RejectSettlementInput{
		SettlementID: settlementID,
		Reason:       req.Reason,
	}

	output, err := h.rejectSettlement.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// CreateSettlement maneja POST /api/v1/admin/settlements
func (h *SettlementHandler) CreateSettlement(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var req struct {
		OrganizerID int64   `json:"organizer_id" binding:"required"`
		RaffleIDs   []int64 `json:"raffle_ids"`
		Mode        string  `json:"mode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &settlement.CreateSettlementInput{
		OrganizerID: req.OrganizerID,
		RaffleIDs:   req.RaffleIDs,
		Mode:        req.Mode,
	}

	output, err := h.createSettlement.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// MarkSettlementPaid maneja POST /api/v1/admin/settlements/:id/mark-paid
func (h *SettlementHandler) MarkSettlementPaid(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	settlementIDStr := c.Param("id")
	settlementID, err := strconv.ParseInt(settlementIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_SETTLEMENT_ID", "invalid settlement ID format", 400, err))
		return
	}

	var req struct {
		PaymentMethod    string  `json:"payment_method" binding:"required"`
		PaymentReference *string `json:"payment_reference"`
		Notes            *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &settlement.MarkSettlementPaidInput{
		SettlementID:     settlementID,
		PaymentMethod:    req.PaymentMethod,
		PaymentReference: req.PaymentReference,
		Notes:            req.Notes,
	}

	output, err := h.markPaid.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// AutoCreateSettlements maneja POST /api/v1/admin/settlements/auto-create
func (h *SettlementHandler) AutoCreateSettlements(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var req struct {
		DaysAfterCompletion int  `json:"days_after_completion" binding:"required"`
		DryRun              bool `json:"dry_run"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &settlement.AutoCreateSettlementsInput{
		DaysAfterCompletion: req.DaysAfterCompletion,
		DryRun:              req.DryRun,
	}

	output, err := h.autoCreate.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
