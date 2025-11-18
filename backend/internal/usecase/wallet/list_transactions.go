package wallet

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// ListTransactionsInput representa los datos de entrada para listar transacciones
type ListTransactionsInput struct {
	UserID int64 `json:"user_id" binding:"required"`
	Limit  int   `json:"limit" binding:"required,min=1,max=100"`
	Offset int   `json:"offset" binding:"min=0"`
}

// ListTransactionsOutput representa los datos de salida
type ListTransactionsOutput struct {
	Transactions []*domain.WalletTransaction `json:"transactions"`
	Total        int64                       `json:"total"`
	Limit        int                         `json:"limit"`
	Offset       int                         `json:"offset"`
}

// ListTransactionsUseCase maneja el listado de transacciones de billetera
type ListTransactionsUseCase struct {
	walletRepo      domain.WalletRepository
	transactionRepo domain.WalletTransactionRepository
	userRepo        domain.UserRepository
	logger          *logger.Logger
}

// NewListTransactionsUseCase crea una nueva instancia del use case
func NewListTransactionsUseCase(
	walletRepo domain.WalletRepository,
	transactionRepo domain.WalletTransactionRepository,
	userRepo domain.UserRepository,
	logger *logger.Logger,
) *ListTransactionsUseCase {
	return &ListTransactionsUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		logger:          logger,
	}
}

// Execute ejecuta el caso de uso de listado de transacciones
func (uc *ListTransactionsUseCase) Execute(ctx context.Context, input *ListTransactionsInput) (*ListTransactionsOutput, error) {
	// Validaciones
	if input.Limit <= 0 {
		input.Limit = 20 // Default
	}
	if input.Limit > 100 {
		input.Limit = 100 // Max
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

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

	// Verificar que la billetera exista
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

	// Obtener transacciones
	transactions, total, err := uc.transactionRepo.FindByWalletID(wallet.ID, input.Limit, input.Offset)
	if err != nil {
		uc.logger.Error("Error listando transacciones",
			logger.Int64("wallet_id", wallet.ID),
			logger.Int("limit", input.Limit),
			logger.Int("offset", input.Offset),
			logger.Error(err))
		return nil, err
	}

	return &ListTransactionsOutput{
		Transactions: transactions,
		Total:        total,
		Limit:        input.Limit,
		Offset:       input.Offset,
	}, nil
}
