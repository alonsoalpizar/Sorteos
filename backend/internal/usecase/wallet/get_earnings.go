package wallet

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GetEarningsInput representa el input del use case
type GetEarningsInput struct {
	UserID int64 `json:"user_id" validate:"required"`
	Limit  int   `json:"limit"`  // Para paginación del desglose (0 = sin límite)
	Offset int   `json:"offset"` // Para paginación del desglose
}

// GetEarningsOutput representa el output del use case
type GetEarningsOutput struct {
	*domain.UserEarnings
}

// GetEarningsUseCase obtiene las ganancias del usuario
type GetEarningsUseCase struct {
	raffleRepo db.RaffleRepository
	logger     *logger.Logger
}

// NewGetEarningsUseCase crea una nueva instancia del use case
func NewGetEarningsUseCase(
	raffleRepo db.RaffleRepository,
	logger *logger.Logger,
) *GetEarningsUseCase {
	return &GetEarningsUseCase{
		raffleRepo: raffleRepo,
		logger:     logger,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetEarningsUseCase) Execute(ctx context.Context, input *GetEarningsInput) (*GetEarningsOutput, error) {
	// Obtener resumen (totales)
	summary, err := uc.raffleRepo.GetUserEarningsSummary(input.UserID)
	if err != nil {
		uc.logger.Error("Failed to get earnings summary",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener desglose detallado
	raffles, err := uc.raffleRepo.GetUserCompletedRaffles(input.UserID, input.Limit, input.Offset)
	if err != nil {
		uc.logger.Error("Failed to get completed raffles",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	summary.Raffles = raffles

	uc.logger.Info("User earnings retrieved",
		logger.Int64("user_id", input.UserID),
		logger.Int("completed_raffles", summary.CompletedRaffles))

	return &GetEarningsOutput{
		UserEarnings: summary,
	}, nil
}
