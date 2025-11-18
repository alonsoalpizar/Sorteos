package wallet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/wallet"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GetEarningsHandler maneja el endpoint de consulta de ganancias
type GetEarningsHandler struct {
	useCase *wallet.GetEarningsUseCase
	logger  *logger.Logger
}

// NewGetEarningsHandler crea una nueva instancia del handler
func NewGetEarningsHandler(useCase *wallet.GetEarningsUseCase, logger *logger.Logger) *GetEarningsHandler {
	return &GetEarningsHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle maneja la petición de consulta de ganancias
// GET /api/v1/wallet/earnings?limit=10&offset=0
func (h *GetEarningsHandler) Handle(c *gin.Context) {
	// Obtener user_id del contexto (middleware de autenticación lo inyecta)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario no autenticado",
		})
		return
	}

	// Parsear query params para paginación
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))   // 0 = sin límite
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Ejecutar caso de uso
	input := &wallet.GetEarningsInput{
		UserID: userID.(int64),
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Preparar respuesta
	raffles := make([]gin.H, 0, len(output.Raffles))
	for _, r := range output.Raffles {
		raffles = append(raffles, gin.H{
			"raffle_id":            r.RaffleID,
			"raffle_uuid":          r.RaffleUUID,
			"title":                r.Title,
			"draw_date":            r.DrawDate,
			"completed_at":         r.CompletedAt,
			"total_revenue":        r.TotalRevenue.String(),
			"platform_fee_percent": r.PlatformFeePercent.String(),
			"platform_fee_amount":  r.PlatformFeeAmount.String(),
			"net_amount":           r.NetAmount.String(),
			"settlement_status":    r.SettlementStatus,
			"settled_at":           r.SettledAt,
		})
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_collected":     output.TotalCollected.String(),
			"platform_commission": output.PlatformCommission.String(),
			"net_earnings":        output.NetEarnings.String(),
			"completed_raffles":   output.CompletedRaffles,
			"raffles":             raffles,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *GetEarningsHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in get earnings handler", logger.Error(err))
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
