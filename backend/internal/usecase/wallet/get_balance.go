package wallet

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GetBalanceInput representa los datos de entrada para consultar saldo
type GetBalanceInput struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// GetBalanceOutput representa los datos de salida
type GetBalanceOutput struct {
	Wallet         *domain.Wallet  `json:"wallet"`
	Balance        decimal.Decimal `json:"balance"`
	PendingBalance decimal.Decimal `json:"pending_balance"`
	Currency       string          `json:"currency"`
	Status         string          `json:"status"`
}

// GetBalanceUseCase maneja la consulta de saldo de billetera
type GetBalanceUseCase struct {
	walletRepo domain.WalletRepository
	userRepo   domain.UserRepository
	logger     *logger.Logger
}

// NewGetBalanceUseCase crea una nueva instancia del use case
func NewGetBalanceUseCase(
	walletRepo domain.WalletRepository,
	userRepo domain.UserRepository,
	logger *logger.Logger,
) *GetBalanceUseCase {
	return &GetBalanceUseCase{
		walletRepo: walletRepo,
		userRepo:   userRepo,
		logger:     logger,
	}
}

// Execute ejecuta el caso de uso de consulta de saldo
func (uc *GetBalanceUseCase) Execute(ctx context.Context, input *GetBalanceInput) (*GetBalanceOutput, error) {
	// Validar que el usuario exista
	user, err := uc.userRepo.FindByID(input.UserID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "usuario no encontrado", err)
		}
		uc.logger.Error("Error buscando usuario",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}

	// Verificar que el usuario esté activo
	if !user.IsActive() {
		return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "el usuario no está activo", nil)
	}

	// Obtener billetera
	wallet, err := uc.walletRepo.FindByUserID(input.UserID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "billetera no encontrada", err)
		}
		uc.logger.Error("Error buscando billetera",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}

	return &GetBalanceOutput{
		Wallet:         wallet,
		Balance:        wallet.BalanceAvailable,
		PendingBalance: wallet.PendingBalance,
		Currency:       wallet.Currency,
		Status:         string(wallet.Status),
	}, nil
}
