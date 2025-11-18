package user

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// DeleteUserInput datos de entrada
type DeleteUserInput struct {
	UserID int64
	Reason string // Razón de la eliminación
}

// DeleteUserUseCase caso de uso para eliminar usuario (soft delete)
type DeleteUserUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewDeleteUserUseCase crea una nueva instancia
func NewDeleteUserUseCase(db *gorm.DB, log *logger.Logger) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *DeleteUserUseCase) Execute(ctx context.Context, input *DeleteUserInput, adminID int64) error {
	// Validar que admin no se elimine a sí mismo
	if input.UserID == adminID {
		return errors.New("VALIDATION_FAILED", "cannot delete your own account", 400, nil)
	}

	// Validar razón
	if input.Reason == "" {
		return errors.New("VALIDATION_FAILED", "reason is required", 400, nil)
	}

	// Obtener usuario
	var user domain.User
	if err := uc.db.Where("id = ? AND deleted_at IS NULL", input.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		uc.log.Error("Error finding user", logger.Int64("user_id", input.UserID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Verificar que el usuario no tenga rifas activas
	var activeRafflesCount int64
	if err := uc.db.Model(&domain.Raffle{}).
		Where("user_id = ? AND status IN (?, ?) AND deleted_at IS NULL",
			input.UserID, domain.RaffleStatusActive, domain.RaffleStatusSuspended).
		Count(&activeRafflesCount).Error; err != nil {
		uc.log.Error("Error counting active raffles", logger.Int64("user_id", input.UserID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	if activeRafflesCount > 0 {
		return errors.New("VALIDATION_FAILED",
			"cannot delete user with active raffles. Please cancel or complete them first", 400, nil)
	}

	// Iniciar transacción
	tx := uc.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// 1. Cancelar rifas en draft del usuario
	if err := tx.Model(&domain.Raffle{}).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", input.UserID, domain.RaffleStatusDraft).
		Updates(map[string]interface{}{
			"status":     domain.RaffleStatusCancelled,
			"deleted_at": now,
		}).Error; err != nil {
		tx.Rollback()
		uc.log.Error("Error cancelling draft raffles", logger.Int64("user_id", input.UserID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 2. Soft delete del usuario
	if err := tx.Model(&domain.User{}).
		Where("id = ?", input.UserID).
		Updates(map[string]interface{}{
			"deleted_at":        now,
			"is_active":         false,
			"suspension_reason": input.Reason,
			"suspended_by":      adminID,
		}).Error; err != nil {
		tx.Rollback()
		uc.log.Error("Error soft deleting user", logger.Int64("user_id", input.UserID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Commit transacción
	if err := tx.Commit().Error; err != nil {
		uc.log.Error("Error committing transaction", logger.Int64("user_id", input.UserID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría crítica
	uc.log.Error("Admin deleted user",
		logger.Int64("admin_id", adminID),
		logger.Int64("user_id", input.UserID),
		logger.String("user_email", user.Email),
		logger.String("reason", input.Reason),
		logger.String("action", "admin_delete_user"),
		logger.String("severity", "critical"))

	return nil
}
