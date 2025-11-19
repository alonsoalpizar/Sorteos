package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UnfreezeWalletInput entrada
type UnfreezeWalletInput struct {
	WalletID int64
	AdminID  int64
}

// UnfreezeWalletOutput salida
type UnfreezeWalletOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UnfreezeWalletUseCase descongela una billetera
type UnfreezeWalletUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUnfreezeWalletUseCase crea una nueva instancia
func NewUnfreezeWalletUseCase(db *gorm.DB, log *logger.Logger) *UnfreezeWalletUseCase {
	return &UnfreezeWalletUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UnfreezeWalletUseCase) Execute(ctx context.Context, input *UnfreezeWalletInput) (*UnfreezeWalletOutput, error) {
	// Obtener wallet
	var wallet domain.Wallet
	if err := uc.db.WithContext(ctx).Where("id = ?", input.WalletID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("billetera no encontrada")
		}
		return nil, err
	}

	// Verificar que esté congelada
	if wallet.Status != domain.WalletStatusFrozen {
		return nil, fmt.Errorf("la billetera no está congelada")
	}

	// Descongelar
	now := time.Now()
	updates := map[string]interface{}{
		"status":     domain.WalletStatusActive,
		"updated_at": now,
	}

	if err := uc.db.WithContext(ctx).Model(&wallet).Updates(updates).Error; err != nil {
		uc.log.Error("Error unfreezing wallet", logger.Error(err))
		return nil, err
	}

	uc.log.Info("Wallet unfrozen by admin",
		logger.Int64("wallet_id", wallet.ID),
		logger.Int64("admin_id", input.AdminID))

	return &UnfreezeWalletOutput{
		Success: true,
		Message: "Billetera descongelada exitosamente",
	}, nil
}
