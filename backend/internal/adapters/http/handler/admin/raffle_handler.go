package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/internal/usecase/admin/raffle"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// RaffleHandler maneja las peticiones HTTP relacionadas con administración de rifas
type RaffleHandler struct {
	listRafflesUC          *raffle.ListRafflesAdminUseCase
	viewTransactionsUC     *raffle.ViewRaffleTransactionsUseCase
	forceStatusChangeUC    *raffle.ForceStatusChangeUseCase
	manualDrawWinnerUC     *raffle.ManualDrawWinnerUseCase
	addAdminNotesUC        *raffle.AddAdminNotesUseCase
	cancelWithRefundUC     *raffle.CancelRaffleWithRefundUseCase
	log                    *logger.Logger
}

// NewRaffleHandler crea una nueva instancia del handler
func NewRaffleHandler(db *gorm.DB, log *logger.Logger) *RaffleHandler {
	return &RaffleHandler{
		listRafflesUC:          raffle.NewListRafflesAdminUseCase(db, log),
		viewTransactionsUC:     raffle.NewViewRaffleTransactionsUseCase(db, log),
		forceStatusChangeUC:    raffle.NewForceStatusChangeUseCase(db, log),
		manualDrawWinnerUC:     raffle.NewManualDrawWinnerUseCase(db, log),
		addAdminNotesUC:        raffle.NewAddAdminNotesUseCase(db, log),
		cancelWithRefundUC:     raffle.NewCancelRaffleWithRefundUseCase(db, log),
		log:                    log,
	}
}

// List lista rifas con filtros y paginación (incluye suspended)
// GET /api/v1/admin/raffles
func (h *RaffleHandler) List(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &raffle.ListRafflesAdminInput{
		Page:       1,
		PageSize:   20,
		Search:     c.Query("search"),
		OrderBy:    c.Query("order_by"),
		IncludeAll: false,
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
	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.RaffleStatus(statusStr)
		input.Status = &status
	}

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

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = &dateTo
	}

	if includeAllStr := c.Query("include_all"); includeAllStr == "true" {
		input.IncludeAll = true
	}

	// Ejecutar use case
	output, err := h.listRafflesUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ViewTransactions obtiene timeline de transacciones de una rifa
// GET /api/v1/admin/raffles/:id/transactions
func (h *RaffleHandler) ViewTransactions(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse raffle ID
	raffleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_RAFFLE_ID",
				"message": "invalid raffle ID",
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.viewTransactionsUC.Execute(c.Request.Context(), raffleID, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ForceStatusChange fuerza un cambio de estado en una rifa (suspend, activate, etc.)
// PUT /api/v1/admin/raffles/:id/status
func (h *RaffleHandler) ForceStatusChange(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse raffle ID
	raffleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_RAFFLE_ID",
				"message": "invalid raffle ID",
			},
		})
		return
	}

	// Parse body
	var body struct {
		NewStatus domain.RaffleStatus `json:"new_status" binding:"required"`
		Reason    string              `json:"reason" binding:"required"`
		Notes     string              `json:"notes"`
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

	input := &raffle.ForceStatusChangeInput{
		RaffleID:  raffleID,
		NewStatus: body.NewStatus,
		Reason:    body.Reason,
		Notes:     body.Notes,
	}

	// Ejecutar use case
	err = h.forceStatusChangeUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Raffle status changed successfully",
	})
}

// ManualDraw ejecuta sorteo manual (selecciona ganador)
// POST /api/v1/admin/raffles/:id/draw
func (h *RaffleHandler) ManualDraw(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse raffle ID
	raffleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_RAFFLE_ID",
				"message": "invalid raffle ID",
			},
		})
		return
	}

	// Parse body
	var body struct {
		WinnerNumber *string `json:"winner_number,omitempty"` // Si es nil, selección aleatoria
		Reason       string  `json:"reason" binding:"required"`
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

	input := &raffle.ManualDrawWinnerInput{
		RaffleID:     raffleID,
		WinnerNumber: body.WinnerNumber,
		Reason:       body.Reason,
	}

	// Ejecutar use case
	output, err := h.manualDrawWinnerUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// AddNotes agrega notas administrativas a una rifa
// POST /api/v1/admin/raffles/:id/notes
func (h *RaffleHandler) AddNotes(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse raffle ID
	raffleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_RAFFLE_ID",
				"message": "invalid raffle ID",
			},
		})
		return
	}

	// Parse body
	var body struct {
		Notes  string `json:"notes" binding:"required"`
		Append bool   `json:"append"` // Default false = replace
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

	input := &raffle.AddAdminNotesInput{
		RaffleID: raffleID,
		Notes:    body.Notes,
		Append:   body.Append,
	}

	// Ejecutar use case
	err = h.addAdminNotesUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin notes added successfully",
	})
}

// CancelWithRefund cancela rifa y reembolsa a todos los compradores
// POST /api/v1/admin/raffles/:id/cancel
func (h *RaffleHandler) CancelWithRefund(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse raffle ID
	raffleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_RAFFLE_ID",
				"message": "invalid raffle ID",
			},
		})
		return
	}

	// Parse body
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

	input := &raffle.CancelRaffleWithRefundInput{
		RaffleID: raffleID,
		Reason:   body.Reason,
	}

	// Ejecutar use case
	output, err := h.cancelWithRefundUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
