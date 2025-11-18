package raffle

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ForceStatusChangeInput datos de entrada
type ForceStatusChangeInput struct {
	RaffleID  int64
	NewStatus domain.RaffleStatus
	Reason    string // Razón del cambio forzado
	Notes     string // Notas adicionales
}

// ForceStatusChangeUseCase caso de uso para cambiar forzadamente el estado de una rifa
type ForceStatusChangeUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewForceStatusChangeUseCase crea una nueva instancia
func NewForceStatusChangeUseCase(db *gorm.DB, log *logger.Logger) *ForceStatusChangeUseCase {
	return &ForceStatusChangeUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ForceStatusChangeUseCase) Execute(ctx context.Context, input *ForceStatusChangeInput, adminID int64) error {
	// Validar razón
	if input.Reason == "" {
		return errors.New("VALIDATION_FAILED", "reason is required for forced status change", 400, nil)
	}

	// Obtener rifa
	var raffle domain.Raffle
	if err := uc.db.Where("id = ?", input.RaffleID).First(&raffle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrRaffleNotFound
		}
		uc.log.Error("Error finding raffle", logger.Int64("raffle_id", input.RaffleID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Validar transiciones permitidas
	if !isValidTransition(raffle.Status, input.NewStatus) {
		return errors.New("VALIDATION_FAILED",
			fmt.Sprintf("invalid transition from %s to %s", raffle.Status, input.NewStatus), 400, nil)
	}

	now := time.Now()
	updates := make(map[string]interface{})
	updates["status"] = input.NewStatus
	updates["updated_at"] = now

	// Manejar casos especiales según el nuevo estado
	switch input.NewStatus {
	case domain.RaffleStatusSuspended:
		// Suspender rifa
		updates["suspended_at"] = now
		updates["suspended_by"] = adminID
		updates["suspension_reason"] = input.Reason
		if input.Notes != "" {
			updates["admin_notes"] = input.Notes
		}

		uc.log.Warn("Admin suspended raffle",
			logger.Int64("admin_id", adminID),
			logger.Int64("raffle_id", input.RaffleID),
			logger.String("reason", input.Reason),
			logger.String("action", "admin_suspend_raffle"))

	case domain.RaffleStatusActive:
		// Reactivar rifa (desde suspended)
		updates["suspended_at"] = nil
		updates["suspended_by"] = nil
		updates["suspension_reason"] = nil
		if input.Notes != "" {
			updates["admin_notes"] = input.Notes
		}

		uc.log.Info("Admin activated raffle",
			logger.Int64("admin_id", adminID),
			logger.Int64("raffle_id", input.RaffleID),
			logger.String("reason", input.Reason),
			logger.String("action", "admin_activate_raffle"))

	case domain.RaffleStatusCancelled:
		// Cancelar rifa (debe hacerse con refund en otro caso de uso)
		if raffle.SoldCount > 0 {
			return errors.New("VALIDATION_FAILED",
				"cannot cancel raffle with sold numbers without refund. Use CancelRaffleWithRefundUseCase", 400, nil)
		}

		updates["admin_notes"] = fmt.Sprintf("Cancelled by admin. Reason: %s", input.Reason)

		uc.log.Error("Admin cancelled raffle",
			logger.Int64("admin_id", adminID),
			logger.Int64("raffle_id", input.RaffleID),
			logger.String("reason", input.Reason),
			logger.String("action", "admin_cancel_raffle"))

	default:
		// Otros cambios de estado
		if input.Notes != "" {
			updates["admin_notes"] = input.Notes
		}

		uc.log.Warn("Admin changed raffle status",
			logger.Int64("admin_id", adminID),
			logger.Int64("raffle_id", input.RaffleID),
			logger.String("old_status", string(raffle.Status)),
			logger.String("new_status", string(input.NewStatus)),
			logger.String("reason", input.Reason),
			logger.String("action", "admin_change_raffle_status"))
	}

	// Actualizar rifa
	if err := uc.db.Model(&domain.Raffle{}).Where("id = ?", input.RaffleID).Updates(updates).Error; err != nil {
		uc.log.Error("Error updating raffle status",
			logger.Int64("raffle_id", input.RaffleID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// TODO: Enviar email de notificación al organizador
	// Esto se implementará cuando tengamos el servicio de email configurado

	return nil
}

// isValidTransition valida si una transición de estado es permitida
func isValidTransition(from, to domain.RaffleStatus) bool {
	// Matriz de transiciones permitidas para admin
	validTransitions := map[domain.RaffleStatus][]domain.RaffleStatus{
		domain.RaffleStatusDraft: {
			domain.RaffleStatusActive,
			domain.RaffleStatusCancelled,
		},
		domain.RaffleStatusActive: {
			domain.RaffleStatusSuspended,
			domain.RaffleStatusCancelled,
			domain.RaffleStatusCompleted,
		},
		domain.RaffleStatusSuspended: {
			domain.RaffleStatusActive,
			domain.RaffleStatusCancelled,
		},
		domain.RaffleStatusCompleted: {
			// Completed es estado final, no permite transiciones
		},
		domain.RaffleStatusCancelled: {
			// Cancelled es estado final, no permite transiciones
		},
	}

	allowedStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowed := range allowedStates {
		if allowed == to {
			return true
		}
	}

	return false
}
