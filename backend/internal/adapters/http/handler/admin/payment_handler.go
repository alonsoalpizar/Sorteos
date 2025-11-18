package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/payment"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PaymentHandler maneja las peticiones HTTP relacionadas con administración de pagos
type PaymentHandler struct {
	listPaymentsUC     *payment.ListPaymentsAdminUseCase
	viewPaymentDetailsUC *payment.ViewPaymentDetailsUseCase
	processRefundUC    *payment.ProcessRefundUseCase
	manageDisputeUC    *payment.ManageDisputeUseCase
	log                *logger.Logger
}

// NewPaymentHandler crea una nueva instancia del handler
func NewPaymentHandler(db *gorm.DB, log *logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		listPaymentsUC:     payment.NewListPaymentsAdminUseCase(db, log),
		viewPaymentDetailsUC: payment.NewViewPaymentDetailsUseCase(db, log),
		processRefundUC:    payment.NewProcessRefundUseCase(db, log),
		manageDisputeUC:    payment.NewManageDisputeUseCase(db, log),
		log:                log,
	}
}

// List lista pagos con filtros y paginación
// GET /api/v1/admin/payments
func (h *PaymentHandler) List(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &payment.ListPaymentsAdminInput{
		Page:          1,
		PageSize:      20,
		Search:        c.Query("search"),
		OrderBy:       c.Query("order_by"),
		IncludeRefund: false,
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

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			input.UserID = &userID
		}
	}

	if raffleIDStr := c.Query("raffle_id"); raffleIDStr != "" {
		if raffleID, err := strconv.ParseInt(raffleIDStr, 10, 64); err == nil {
			input.RaffleID = &raffleID
		}
	}

	if organizerIDStr := c.Query("organizer_id"); organizerIDStr != "" {
		if organizerID, err := strconv.ParseInt(organizerIDStr, 10, 64); err == nil {
			input.OrganizerID = &organizerID
		}
	}

	if provider := c.Query("provider"); provider != "" {
		input.Provider = &provider
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

	if includeRefundStr := c.Query("include_refund"); includeRefundStr == "true" {
		input.IncludeRefund = true
	}

	// Ejecutar use case
	output, err := h.listPaymentsUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// GetByID obtiene detalles completos de un pago
// GET /api/v1/admin/payments/:id
func (h *PaymentHandler) GetByID(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse payment ID (UUID string)
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_PAYMENT_ID",
				"message": "invalid payment ID",
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.viewPaymentDetailsUC.Execute(c.Request.Context(), paymentID, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ProcessRefund procesa un reembolso (full o partial)
// POST /api/v1/admin/payments/:id/refund
func (h *PaymentHandler) ProcessRefund(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse payment ID
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_PAYMENT_ID",
				"message": "invalid payment ID",
			},
		})
		return
	}

	// Parse body
	var body struct {
		Reason string   `json:"reason" binding:"required"`
		Amount *float64 `json:"amount,omitempty"` // Si es nil, refund total
		Notes  string   `json:"notes"`
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

	input := &payment.ProcessRefundInput{
		PaymentID: paymentID,
		Reason:    body.Reason,
		Amount:    body.Amount,
		Notes:     body.Notes,
	}

	// Ejecutar use case
	output, err := h.processRefundUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ManageDispute gestiona disputas de pagos (open, update, close, escalate)
// POST /api/v1/admin/payments/:id/dispute
func (h *PaymentHandler) ManageDispute(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse payment ID
	paymentID := c.Param("id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_PAYMENT_ID",
				"message": "invalid payment ID",
			},
		})
		return
	}

	// Parse body
	var input payment.ManageDisputeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Asignar payment ID desde params
	input.PaymentID = paymentID

	// Ejecutar use case
	output, err := h.manageDisputeUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
