package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// FreezeWalletInput entrada
type FreezeWalletInput struct {
	WalletID int64
	AdminID  int64
	Reason   string
}

// FreezeWalletOutput salida
type FreezeWalletOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// FreezeWalletUseCase congela una billetera
type FreezeWalletUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewFreezeWalletUseCase crea una nueva instancia
func NewFreezeWalletUseCase(db *gorm.DB, log *logger.Logger) *FreezeWalletUseCase {
	return &FreezeWalletUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *FreezeWalletUseCase) Execute(ctx context.Context, input *FreezeWalletInput) (*FreezeWalletOutput, error) {
	// Validar razón
	if input.Reason == "" {
		return nil, fmt.Errorf("la razón es requerida")
	}

	// Obtener wallet
	var wallet domain.Wallet
	if err := uc.db.WithContext(ctx).Where("id = ?", input.WalletID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("billetera no encontrada")
		}
		return nil, err
	}

	// Verificar que no esté ya congelada
	if wallet.Status == domain.WalletStatusFrozen {
		return nil, fmt.Errorf("la billetera ya está congelada")
	}

	// Congelar (solo cambiamos el status, la razón se loguea en auditoría)
	now := time.Now()
	updates := map[string]interface{}{
		"status":     domain.WalletStatusFrozen,
		"updated_at": now,
	}

	if err := uc.db.WithContext(ctx).Model(&wallet).Updates(updates).Error; err != nil {
		uc.log.Error("Error freezing wallet", logger.Error(err))
		return nil, err
	}

	uc.log.Info("Wallet frozen by admin",
		logger.Int64("wallet_id", wallet.ID),
		logger.Int64("admin_id", input.AdminID),
		logger.String("reason", input.Reason))

	return &FreezeWalletOutput{
		Success: true,
		Message: "Billetera congelada exitosamente",
	}, nil
}
