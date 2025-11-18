package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/payment"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// PaymentHandler maneja todas las operaciones de administración de pagos
type PaymentHandler struct {
	listPayments       *payment.ListPaymentsUseCase
	viewPaymentDetails *payment.ViewPaymentDetailsUseCase
	processRefund      *payment.ProcessRefundUseCase
	manageDispute      *payment.ManageDisputeUseCase
}

// NewPaymentHandler crea una nueva instancia
func NewPaymentHandler(
	listPayments *payment.ListPaymentsUseCase,
	viewPaymentDetails *payment.ViewPaymentDetailsUseCase,
	processRefund *payment.ProcessRefundUseCase,
	manageDispute *payment.ManageDisputeUseCase,
) *PaymentHandler {
	return &PaymentHandler{
		listPayments:       listPayments,
		viewPaymentDetails: viewPaymentDetails,
		processRefund:      processRefund,
		manageDispute:      manageDispute,
	}
}

// ListPayments maneja GET /api/v1/admin/payments
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parsear query params
	var page, pageSize int
	if p := c.Query("page"); p != "" {
		page, _ = strconv.Atoi(p)
	} else {
		page = 1
	}
	if ps := c.Query("page_size"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	} else {
		pageSize = 20
	}

	input := &payment.ListPaymentsInput{
		Page:     page,
		PageSize: pageSize,
		Search:   stringPtr(c.Query("search")),
		OrderBy:  stringPtr(c.Query("order_by")),
	}

	// Filtros opcionales
	if statusStr := c.Query("status"); statusStr != "" {
		input.Status = stringPtr(statusStr)
	}
	if providerStr := c.Query("provider"); providerStr != "" {
		input.Provider = stringPtr(providerStr)
	}
	if hasDispute := c.Query("has_dispute"); hasDispute == "true" {
		boolVal := true
		input.HasDispute = &boolVal
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = stringPtr(dateFrom)
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = stringPtr(dateTo)
	}

	output, err := h.listPayments.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetPaymentByID maneja GET /api/v1/admin/payments/:id
func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		handleError(c, errors.New("INVALID_PAYMENT_ID", "payment ID is required", 400, nil))
		return
	}

	input := &payment.ViewPaymentDetailsInput{
		PaymentID: paymentID,
	}

	output, err := h.viewPaymentDetails.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ProcessRefund maneja POST /api/v1/admin/payments/:id/refund
func (h *PaymentHandler) ProcessRefund(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		handleError(c, errors.New("INVALID_PAYMENT_ID", "payment ID is required", 400, nil))
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
		Reason string  `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &payment.ProcessRefundInput{
		PaymentID: paymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
	}

	output, err := h.processRefund.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ManageDispute maneja POST /api/v1/admin/payments/:id/dispute
func (h *PaymentHandler) ManageDispute(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		handleError(c, errors.New("INVALID_PAYMENT_ID", "payment ID is required", 400, nil))
		return
	}

	var req struct {
		Action          string                 `json:"action" binding:"required"`
		DisputeReason   *string                `json:"dispute_reason"`
		DisputeEvidence *string                `json:"dispute_evidence"`
		Resolution      *string                `json:"resolution"`
		AdminNotes      *string                `json:"admin_notes"`
		Metadata        map[string]interface{} `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &payment.ManageDisputeInput{
		PaymentID:       paymentID,
		Action:          req.Action,
		DisputeReason:   req.DisputeReason,
		DisputeEvidence: req.DisputeEvidence,
		Resolution:      req.Resolution,
		AdminNotes:      req.AdminNotes,
		Metadata:        req.Metadata,
	}

	output, err := h.manageDispute.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
