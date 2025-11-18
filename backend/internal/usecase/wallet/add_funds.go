package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// AddFundsInput representa los datos de entrada para agregar fondos
type AddFundsInput struct {
	UserID         int64           `json:"user_id" binding:"required"`
	Amount         decimal.Decimal `json:"amount" binding:"required"`
	IdempotencyKey string          `json:"idempotency_key" binding:"required"`
	PaymentMethod  string          `json:"payment_method" binding:"required"` // "stripe", "paypal", etc.
	PaymentIntentID *string        `json:"payment_intent_id,omitempty"`       // ID del procesador externo
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// AddFundsOutput representa los datos de salida
type AddFundsOutput struct {
	Transaction    *domain.WalletTransaction `json:"transaction"`
	NewBalance     decimal.Decimal           `json:"new_balance"`
	PaymentIntentID *string                  `json:"payment_intent_id,omitempty"` // Para completar con procesador
}

// AddFundsUseCase maneja la adición de fondos a la billetera
type AddFundsUseCase struct {
	walletRepo      domain.WalletRepository
	transactionRepo domain.WalletTransactionRepository
	userRepo        domain.UserRepository
	auditRepo       domain.AuditLogRepository
	logger          *logger.Logger
}

// NewAddFundsUseCase crea una nueva instancia del use case
func NewAddFundsUseCase(
	walletRepo domain.WalletRepository,
	transactionRepo domain.WalletTransactionRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	logger *logger.Logger,
) *AddFundsUseCase {
	return &AddFundsUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		auditRepo:       auditRepo,
		logger:          logger,
	}
}

// Execute ejecuta el caso de uso de adición de fondos
// NOTA: Este caso de uso crea una transacción PENDIENTE que se completará
// cuando el webhook del procesador de pagos confirme el pago
func (uc *AddFundsUseCase) Execute(ctx context.Context, input *AddFundsInput) (*AddFundsOutput, error) {
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

		return &AddFundsOutput{
			Transaction:    existingTx,
			NewBalance:     wallet.BalanceAvailable,
			PaymentIntentID: input.PaymentIntentID,
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

	// Obtener o crear billetera
	wallet, err := uc.walletRepo.FindByUserID(input.UserID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "billetera no encontrada - debe crearse primero", err)
		}
		uc.logger.Error("Error buscando billetera",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}

	// Validar que la billetera esté activa
	if !wallet.IsActive() {
		return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "la billetera no está activa", nil)
	}

	// Crear transacción PENDIENTE
	// Esta transacción se completará cuando el webhook confirme el pago
	transaction := &domain.WalletTransaction{
		UUID:           uuid.New().String(),
		WalletID:       wallet.ID,
		UserID:         input.UserID,
		Type:           domain.TransactionTypeDeposit,
		Amount:         input.Amount,
		Status:         domain.TransactionStatusPending,
		BalanceBefore:  wallet.BalanceAvailable,
		BalanceAfter:   wallet.BalanceAvailable, // Aún no se acredita, se hará en el webhook
		IdempotencyKey: input.IdempotencyKey,
	}

	// Agregar referencia al payment intent si existe
	if input.PaymentIntentID != nil {
		refType := "payment_intent"
		transaction.ReferenceType = &refType
		// Nota: ReferenceID es int64, pero PaymentIntentID es string
		// Esto se manejará mejor con metadata
	}

	// Validar transacción
	if err := transaction.Validate(); err != nil {
		return nil, errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Guardar transacción en DB
	if err := uc.transactionRepo.Create(transaction); err != nil {
		uc.logger.Error("Error creando transacción",
			logger.Int64("wallet_id", wallet.ID),
			logger.Error(err))
		return nil, err
	}

	// Registrar en audit log
	entityType := "wallet_transaction"
	entityID := transaction.ID
	auditLog := &domain.AuditLog{
		Action:     "wallet_add_funds_pending",
		EntityType: &entityType,
		EntityID:   &entityID,
		UserID:     &user.ID,
	}
	if err := auditLog.SetMetadata(map[string]interface{}{
		"amount":            input.Amount.String(),
		"payment_method":    input.PaymentMethod,
		"payment_intent_id": input.PaymentIntentID,
	}); err != nil {
		uc.logger.Warn("Error setting audit log metadata", logger.Error(err))
	}
	if err := uc.auditRepo.Create(auditLog); err != nil {
		uc.logger.Warn("Error creando audit log",
			logger.Int64("tx_id", transaction.ID),
			logger.Error(err))
		// No retornar error, audit log no es crítico
	}

	uc.logger.Info("Transacción de depósito creada (pendiente)",
		logger.Int64("tx_id", transaction.ID),
		logger.Int64("user_id", input.UserID),
		logger.String("amount", input.Amount.String()),
		logger.String("payment_method", input.PaymentMethod))

	return &AddFundsOutput{
		Transaction:    transaction,
		NewBalance:     wallet.BalanceAvailable, // No ha cambiado aún
		PaymentIntentID: input.PaymentIntentID,
	}, nil
}

// ConfirmAddFunds completa una transacción de depósito pendiente
// Este método se llamará desde el webhook del procesador de pagos
func (uc *AddFundsUseCase) ConfirmAddFunds(ctx context.Context, transactionID int64) error {
	// Ejecutar dentro de una transacción atómica
	return uc.walletRepo.WithTransaction(func(walletRepo domain.WalletRepository) error {
		// 1. Obtener transacción
		tx, err := uc.transactionRepo.FindByID(transactionID)
		if err != nil {
			return err
		}

		// 2. Verificar que esté pendiente
		if !tx.IsPending() {
			uc.logger.Warn("Intento de confirmar transacción no pendiente",
				logger.Int64("tx_id", transactionID),
				logger.String("status", string(tx.Status)))
			return errors.WrapWithMessage(errors.ErrValidationFailed, "la transacción no está pendiente", nil)
		}

		// 3. Obtener billetera con lock
		wallet, err := walletRepo.FindByID(tx.WalletID)
		if err != nil {
			return err
		}

		if err := walletRepo.Lock(wallet.ID); err != nil {
			return err
		}

		// 4. Acreditar fondos
		if err := wallet.Credit(tx.Amount); err != nil {
			return err
		}

		// 5. Actualizar balance_after en la transacción
		tx.BalanceAfter = wallet.BalanceAvailable

		// 6. Marcar transacción como completada
		if err := tx.MarkAsCompleted(); err != nil {
			return err
		}

		// 7. Guardar cambios
		if err := walletRepo.Update(wallet); err != nil {
			return err
		}

		if err := uc.transactionRepo.Update(tx); err != nil {
			return err
		}

		uc.logger.Info("Depósito confirmado exitosamente",
			logger.Int64("tx_id", tx.ID),
			logger.Int64("wallet_id", wallet.ID),
			logger.String("amount", tx.Amount.String()),
			logger.String("new_balance", wallet.BalanceAvailable.String()))

		return nil
	})
}
