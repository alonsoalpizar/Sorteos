package user

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UserStatusAction acción de cambio de estado
type UserStatusAction string

const (
	UserStatusActionSuspend   UserStatusAction = "suspend"
	UserStatusActionActivate  UserStatusAction = "activate"
	UserStatusActionBan       UserStatusAction = "ban"
	UserStatusActionUnban     UserStatusAction = "unban"
)

// UpdateUserStatusInput datos de entrada
type UpdateUserStatusInput struct {
	UserID int64
	Action UserStatusAction
	Reason string // Requerido para suspend y ban
}

// UpdateUserStatusUseCase caso de uso para actualizar estado de usuario
type UpdateUserStatusUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUpdateUserStatusUseCase crea una nueva instancia
func NewUpdateUserStatusUseCase(db *gorm.DB, log *logger.Logger) *UpdateUserStatusUseCase {
	return &UpdateUserStatusUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateUserStatusUseCase) Execute(ctx context.Context, input *UpdateUserStatusInput, adminID int64) error {
	// Validar que admin no se suspenda/banee a sí mismo
	if input.UserID == adminID {
		return errors.New("VALIDATION_FAILED", "cannot modify your own status", 400, nil)
	}

	// Validar que reason sea requerido para suspend y ban
	if (input.Action == UserStatusActionSuspend || input.Action == UserStatusActionBan) && input.Reason == "" {
		return errors.New("VALIDATION_FAILED", "reason is required for suspend/ban actions", 400, nil)
	}

	// Obtener usuario
	var user domain.User
	if err := uc.db.Where("id = ?", input.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		uc.log.Error("Error finding user", logger.Int64("user_id", input.UserID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	now := time.Now()
	updates := make(map[string]interface{})

	switch input.Action {
	case UserStatusActionSuspend:
		// Suspender usuario
		updates["suspended_at"] = now
		updates["suspended_by"] = adminID
		updates["suspension_reason"] = input.Reason

		uc.log.Warn("Admin suspended user",
			logger.Int64("admin_id", adminID),
			logger.Int64("user_id", input.UserID),
			logger.String("reason", input.Reason),
			logger.String("action", "admin_suspend_user"))

	case UserStatusActionActivate:
		// Reactivar usuario suspendido
		updates["suspended_at"] = nil
		updates["suspended_by"] = nil
		updates["suspension_reason"] = nil

		uc.log.Info("Admin activated user",
			logger.Int64("admin_id", adminID),
			logger.Int64("user_id", input.UserID),
			logger.String("action", "admin_activate_user"))

	case UserStatusActionBan:
		// Banear usuario (desactivar cuenta permanentemente)
		updates["is_active"] = false
		updates["suspended_at"] = now
		updates["suspended_by"] = adminID
		updates["suspension_reason"] = input.Reason

		uc.log.Error("Admin banned user",
			logger.Int64("admin_id", adminID),
			logger.Int64("user_id", input.UserID),
			logger.String("reason", input.Reason),
			logger.String("action", "admin_ban_user"))

	case UserStatusActionUnban:
		// Desbanear usuario
		updates["is_active"] = true
		updates["suspended_at"] = nil
		updates["suspended_by"] = nil
		updates["suspension_reason"] = nil

		uc.log.Info("Admin unbanned user",
			logger.Int64("admin_id", adminID),
			logger.Int64("user_id", input.UserID),
			logger.String("action", "admin_unban_user"))

	default:
		return errors.New("VALIDATION_FAILED", fmt.Sprintf("invalid action: %s", input.Action), 400, nil)
	}

	// Actualizar usuario
	if err := uc.db.Model(&domain.User{}).Where("id = ?", input.UserID).Updates(updates).Error; err != nil {
		uc.log.Error("Error updating user status",
			logger.Int64("user_id", input.UserID),
			logger.String("action", string(input.Action)),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// TODO: Enviar email de notificación al usuario
	// Esto se implementará cuando tengamos el servicio de email configurado

	return nil
}
