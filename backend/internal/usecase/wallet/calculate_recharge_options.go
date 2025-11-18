package wallet

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// CalculateRechargeOptionsInput no requiere inputs, retorna opciones predefinidas
type CalculateRechargeOptionsInput struct{}

// CalculateRechargeOptionsOutput contiene las opciones predefinidas con sus desgloses
type CalculateRechargeOptionsOutput struct {
	Options []*domain.RechargeBreakdown `json:"options"`
}

// CalculateRechargeOptionsUseCase calcula las opciones de recarga predefinidas
type CalculateRechargeOptionsUseCase struct {
	calculator *domain.RechargeCalculator
	logger     *logger.Logger
}

// NewCalculateRechargeOptionsUseCase crea una nueva instancia del use case
func NewCalculateRechargeOptionsUseCase(
	calculator *domain.RechargeCalculator,
	logger *logger.Logger,
) *CalculateRechargeOptionsUseCase {
	return &CalculateRechargeOptionsUseCase{
		calculator: calculator,
		logger:     logger,
	}
}

// Execute ejecuta el caso de uso
func (uc *CalculateRechargeOptionsUseCase) Execute(ctx context.Context, input *CalculateRechargeOptionsInput) (*CalculateRechargeOptionsOutput, error) {
	// Obtener opciones predefinidas del calculator
	options := uc.calculator.GetPredefinedRechargeOptions()

	uc.logger.Info("Recharge options calculated",
		logger.Int("count", len(options)))

	return &CalculateRechargeOptionsOutput{
		Options: options,
	}, nil
}
