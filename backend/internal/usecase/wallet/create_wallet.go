package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// CreateWalletInput representa los datos de entrada para crear billetera
type CreateWalletInput struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Currency string `json:"currency" binding:"required,len=3"` // "USD", "CRC"
}

// CreateWalletOutput representa los datos de salida
type CreateWalletOutput struct {
	Wallet *domain.Wallet `json:"wallet"`
}

// CreateWalletUseCase maneja la creación de billeteras
type CreateWalletUseCase struct {
	walletRepo domain.WalletRepository
	userRepo   domain.UserRepository
	auditRepo  domain.AuditLogRepository
	logger     *logger.Logger
}

// NewCreateWalletUseCase crea una nueva instancia del use case
func NewCreateWalletUseCase(
	walletRepo domain.WalletRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	logger *logger.Logger,
) *CreateWalletUseCase {
	return &CreateWalletUseCase{
		walletRepo: walletRepo,
		userRepo:   userRepo,
		auditRepo:  auditRepo,
		logger:     logger,
	}
}

// Execute ejecuta el caso de uso de creación de billetera
func (uc *CreateWalletUseCase) Execute(ctx context.Context, input *CreateWalletInput) (*CreateWalletOutput, error) {
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

	// Verificar que el usuario no tenga ya una billetera
	existingWallet, err := uc.walletRepo.FindByUserID(input.UserID)
	if err != nil && err != errors.ErrNotFound {
		uc.logger.Error("Error verificando billetera existente",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}
	if existingWallet != nil {
		return nil, errors.WrapWithMessage(errors.ErrConflict, "el usuario ya tiene una billetera", nil)
	}

	// Crear billetera
	wallet := &domain.Wallet{
		UUID:             uuid.New().String(),
		UserID:           input.UserID,
		BalanceAvailable: decimal.Zero,
		EarningsBalance:  decimal.Zero,
		PendingBalance:   decimal.Zero,
		Currency:         input.Currency,
		Status:           domain.WalletStatusActive,
	}

	// Validar billetera
	if err := wallet.Validate(); err != nil {
		return nil, errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Guardar en base de datos
	if err := uc.walletRepo.Create(wallet); err != nil {
		uc.logger.Error("Error creando billetera",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}

	// Registrar en audit log
	entityType := "wallet"
	entityID := wallet.ID
	auditLog := &domain.AuditLog{
		Action:     "wallet_created",
		EntityType: &entityType,
		EntityID:   &entityID,
		UserID:     &user.ID,
	}
	if err := auditLog.SetMetadata(map[string]interface{}{
		"user_id":  user.ID,
		"currency": wallet.Currency,
	}); err != nil {
		uc.logger.Warn("Error setting audit log metadata", logger.Error(err))
	}
	if err := uc.auditRepo.Create(auditLog); err != nil {
		uc.logger.Warn("Error creando audit log",
			logger.Int64("wallet_id", wallet.ID),
			logger.Error(err))
		// No retornar error, audit log no es crítico
	}

	uc.logger.Info("Billetera creada exitosamente",
		logger.Int64("wallet_id", wallet.ID),
		logger.Int64("user_id", wallet.UserID),
		logger.String("wallet_uuid", wallet.UUID))

	return &CreateWalletOutput{
		Wallet: wallet,
	}, nil
}
