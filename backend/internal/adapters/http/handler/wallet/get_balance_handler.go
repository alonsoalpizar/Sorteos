package wallet

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/wallet"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GetBalanceHandler maneja el endpoint de consulta de saldo
type GetBalanceHandler struct {
	useCase *wallet.GetBalanceUseCase
	logger  *logger.Logger
}

// NewGetBalanceHandler crea una nueva instancia del handler
func NewGetBalanceHandler(useCase *wallet.GetBalanceUseCase, logger *logger.Logger) *GetBalanceHandler {
	return &GetBalanceHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle maneja la petición de consulta de saldo
// GET /api/v1/wallet/balance
func (h *GetBalanceHandler) Handle(c *gin.Context) {
	// Obtener user_id del contexto (middleware de autenticación lo inyecta)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario no autenticado",
		})
		return
	}

	// Ejecutar caso de uso
	input := &wallet.GetBalanceInput{
		UserID: userID.(int64),
	}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"wallet_id":       output.Wallet.ID,
			"wallet_uuid":     output.Wallet.UUID,
			"balance":         output.Balance.String(),
			"pending_balance": output.PendingBalance.String(),
			"currency":        output.Currency,
			"status":          output.Status,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *GetBalanceHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in get balance handler", logger.Error(err))
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
