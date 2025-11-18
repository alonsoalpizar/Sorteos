package wallet

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/usecase/wallet"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// CalculateRechargeOptionsHandler maneja el endpoint de calcular opciones de recarga
type CalculateRechargeOptionsHandler struct {
	useCase *wallet.CalculateRechargeOptionsUseCase
	logger  *logger.Logger
}

// NewCalculateRechargeOptionsHandler crea una nueva instancia del handler
func NewCalculateRechargeOptionsHandler(
	useCase *wallet.CalculateRechargeOptionsUseCase,
	logger *logger.Logger,
) *CalculateRechargeOptionsHandler {
	return &CalculateRechargeOptionsHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Handle maneja la petición de calcular opciones de recarga
// GET /api/v1/wallet/recharge-options
func (h *CalculateRechargeOptionsHandler) Handle(c *gin.Context) {
	// Este endpoint no requiere autenticación ya que solo calcula opciones predefinidas
	// Puede ser usado antes del login para mostrar precios

	input := &wallet.CalculateRechargeOptionsInput{}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("Error calculating recharge options", logger.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Error al calcular opciones de recarga",
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"options": output.Options,
			"currency": "CRC",
			"note": "Los montos mostrados incluyen todas las comisiones. El crédito deseado es lo que recibirás en tu billetera.",
		},
	})
}
