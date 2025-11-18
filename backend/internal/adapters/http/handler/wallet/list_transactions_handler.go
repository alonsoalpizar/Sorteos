package wallet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/wallet"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// ListTransactionsHandler maneja el endpoint de listado de transacciones
type ListTransactionsHandler struct {
	useCase *wallet.ListTransactionsUseCase
	logger  *logger.Logger
}

// NewListTransactionsHandler crea una nueva instancia del handler
func NewListTransactionsHandler(useCase *wallet.ListTransactionsUseCase, logger *logger.Logger) *ListTransactionsHandler {
	return &ListTransactionsHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle maneja la petición de listado de transacciones
// GET /api/v1/wallet/transactions?limit=20&offset=0
func (h *ListTransactionsHandler) Handle(c *gin.Context) {
	// Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario no autenticado",
		})
		return
	}

	// Parsear query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validar límites
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Ejecutar caso de uso
	input := &wallet.ListTransactionsInput{
		UserID: userID.(int64),
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Convertir transacciones a formato de respuesta
	transactions := make([]gin.H, len(output.Transactions))
	for i, tx := range output.Transactions {
		transactions[i] = gin.H{
			"id":             tx.ID,
			"uuid":           tx.UUID,
			"type":           tx.Type,
			"amount":         tx.Amount.String(),
			"status":         tx.Status,
			"balance_before": tx.BalanceBefore.String(),
			"balance_after":  tx.BalanceAfter.String(),
			"reference_type": tx.ReferenceType,
			"reference_id":   tx.ReferenceID,
			"notes":          tx.Notes,
			"created_at":     tx.CreatedAt,
			"completed_at":   tx.CompletedAt,
		}
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"transactions": transactions,
			"pagination": gin.H{
				"total":  output.Total,
				"limit":  output.Limit,
				"offset": output.Offset,
			},
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *ListTransactionsHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in list transactions handler", logger.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Error interno del servidor",
		})
		return
	}

	c.JSON(appErr.Status, ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
	})
}
