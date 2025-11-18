package wallet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// DebitFundsInput representa los datos de entrada para debitar fondos
type DebitFundsInput struct {
	UserID         int64           `json:"user_id" binding:"required"`
	Amount         decimal.Decimal `json:"amount" binding:"required"`
	IdempotencyKey string          `json:"idempotency_key" binding:"required"`
	ReferenceType  *string         `json:"reference_type,omitempty"` // "payment", "raffle", etc.
	ReferenceID    *int64          `json:"reference_id,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// DebitFundsOutput representa los datos de salida
type DebitFundsOutput struct {
	Transaction *domain.WalletTransaction `json:"transaction"`
	NewBalance  decimal.Decimal           `json:"new_balance"`
}

// DebitFundsUseCase maneja el débito de fondos de la billetera
type DebitFundsUseCase struct {
	walletRepo     domain.WalletRepository
	transactionRepo domain.WalletTransactionRepository
	userRepo       domain.UserRepository
	auditRepo      domain.AuditLogRepository
	logger         *logger.Logger
}

// NewDebitFundsUseCase crea una nueva instancia del use case
func NewDebitFundsUseCase(
	walletRepo domain.WalletRepository,
	transactionRepo domain.WalletTransactionRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	logger *logger.Logger,
) *DebitFundsUseCase {
	return &DebitFundsUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		auditRepo:       auditRepo,
		logger:          logger,
	}
}

// Execute ejecuta el caso de uso de débito de fondos
func (uc *DebitFundsUseCase) Execute(ctx context.Context, input *DebitFundsInput) (*DebitFundsOutput, error) {
	// Validar amount
	if input.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "el monto debe ser mayor a cero", nil)
	}

	// Verificar idempotencia (prevenir transacciones duplicadas)
	existingTx, err := uc.transactionRepo.FindByIdempotencyKey(input.IdempotencyKey)
	if err != nil && err != errors.ErrNotFound {
		uc.logger.Error("Error verificando idempotencia",
			logger.String("idempotency_key", input.IdempotencyKey),
			logger.Error(err))
		return nil, err
	}
	if existingTx != nil {
		// Transacción duplicada detectada - retornar la existente
		uc.logger.Warn("Transacción duplicada detectada (idempotencia)",
			logger.String("idempotency_key", input.IdempotencyKey),
			logger.Int64("existing_tx_id", existingTx.ID))

		wallet, err := uc.walletRepo.FindByUserID(input.UserID)
		if err != nil {
			return nil, err
		}

		return &DebitFundsOutput{
			Transaction: existingTx,
			NewBalance:  wallet.BalanceAvailable,
		}, nil
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

	// Ejecutar débito dentro de una transacción atómica
	var transaction *domain.WalletTransaction
	var newBalance decimal.Decimal

	err = uc.walletRepo.WithTransaction(func(walletRepo domain.WalletRepository) error {
		// 1. Obtener billetera con lock (SELECT ... FOR UPDATE)
		wallet, err := walletRepo.FindByUserID(input.UserID)
		if err != nil {
			if err == errors.ErrNotFound {
				return errors.WrapWithMessage(errors.ErrValidationFailed, "billetera no encontrada", err)
			}
			return err
		}

		// 2. Adquirir lock explícito
		if err := walletRepo.Lock(wallet.ID); err != nil {
			uc.logger.Error("Error adquiriendo lock de billetera",
				logger.Int64("wallet_id", wallet.ID),
				logger.Error(err))
			return err
		}

		// 3. Validar que se pueda debitar
		if err := wallet.CanDebit(input.Amount); err != nil {
			uc.logger.Warn("Débito rechazado - validación fallida",
				logger.Int64("user_id", input.UserID),
				logger.String("amount", input.Amount.String()),
				logger.String("balance", wallet.BalanceAvailable.String()),
				logger.Error(err))
			return errors.Wrap(errors.ErrValidationFailed, err)
		}

		// 4. Crear snapshot de saldos
		balanceBefore := wallet.BalanceAvailable
		balanceAfter := wallet.BalanceAvailable.Sub(input.Amount)

		// 5. Crear transacción
		transaction = &domain.WalletTransaction{
			UUID:          uuid.New().String(),
			WalletID:      wallet.ID,
			UserID:        input.UserID,
			Type:          domain.TransactionTypePurchase,
			Amount:        input.Amount,
			Status:        domain.TransactionStatusCompleted,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			ReferenceType: input.ReferenceType,
			ReferenceID:   input.ReferenceID,
			IdempotencyKey: input.IdempotencyKey,
			Notes:         input.Notes,
		}

		// Marcar como completada inmediatamente
		now := time.Now()
		transaction.CompletedAt = &now

		// 6. Validar transacción
		if err := transaction.Validate(); err != nil {
			return errors.Wrap(errors.ErrValidationFailed, err)
		}

		// 7. Debitar de la billetera
		if err := wallet.Debit(input.Amount); err != nil {
			return errors.Wrap(errors.ErrValidationFailed, err)
		}

		// 8. Actualizar billetera en DB
		if err := walletRepo.Update(wallet); err != nil {
			uc.logger.Error("Error actualizando billetera",
				logger.Int64("wallet_id", wallet.ID),
				logger.Error(err))
			return err
		}

		// 9. Guardar transacción en DB
		if err := uc.transactionRepo.Create(transaction); err != nil {
			uc.logger.Error("Error creando transacción",
				logger.Int64("wallet_id", wallet.ID),
				logger.Error(err))
			return err
		}

		newBalance = wallet.BalanceAvailable
		return nil
	})

	if err != nil {
		uc.logger.Error("Error en transacción de débito",
			logger.Int64("user_id", input.UserID),
			logger.String("amount", input.Amount.String()),
			logger.Error(err))
		return nil, err
	}

	// Registrar en audit log
	entityType := "wallet_transaction"
	entityID := transaction.ID
	auditLog := &domain.AuditLog{
		Action:     "wallet_debit",
		EntityType: &entityType,
		EntityID:   &entityID,
		UserID:     &user.ID,
	}
	if err := auditLog.SetMetadata(map[string]interface{}{
		"amount":         input.Amount.String(),
		"new_balance":    newBalance.String(),
		"reference_type": input.ReferenceType,
		"reference_id":   input.ReferenceID,
	}); err != nil {
		uc.logger.Warn("Error setting audit log metadata", logger.Error(err))
	}
	if err := uc.auditRepo.Create(auditLog); err != nil {
		uc.logger.Warn("Error creando audit log",
			logger.Int64("tx_id", transaction.ID),
			logger.Error(err))
		// No retornar error, audit log no es crítico
	}

	uc.logger.Info("Débito realizado exitosamente",
		logger.Int64("tx_id", transaction.ID),
		logger.Int64("user_id", input.UserID),
		logger.String("amount", input.Amount.String()),
		logger.String("new_balance", newBalance.String()))

	return &DebitFundsOutput{
		Transaction: transaction,
		NewBalance:  newBalance,
	}, nil
}
