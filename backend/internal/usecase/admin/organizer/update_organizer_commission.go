package organizer

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// UpdateOrganizerCommissionInput datos de entrada
type UpdateOrganizerCommissionInput struct {
	UserID     int64
	Commission *float64 // NULL = usar default global
	Notes      string   // Razón del cambio
}

// UpdateOrganizerCommissionUseCase caso de uso para actualizar comisión de organizador
type UpdateOrganizerCommissionUseCase struct {
	organizerRepo *db.PostgresOrganizerProfileRepository
	log           *logger.Logger
}

// NewUpdateOrganizerCommissionUseCase crea una nueva instancia
func NewUpdateOrganizerCommissionUseCase(
	organizerRepo *db.PostgresOrganizerProfileRepository,
	log *logger.Logger,
) *UpdateOrganizerCommissionUseCase {
	return &UpdateOrganizerCommissionUseCase{
		organizerRepo: organizerRepo,
		log:           log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateOrganizerCommissionUseCase) Execute(ctx context.Context, input *UpdateOrganizerCommissionInput, adminID int64) error {
	// Validar comisión si está presente
	if input.Commission != nil {
		if *input.Commission < 0 || *input.Commission > 50 {
			return errors.New("VALIDATION_FAILED", "commission must be between 0 and 50 percent", 400, nil)
		}
	}

	// Obtener perfil actual para logging
	profile, err := uc.organizerRepo.GetByUserID(input.UserID)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.New("NOT_FOUND", "organizer profile not found", 404, nil)
		}
		uc.log.Error("Error getting organizer profile", logger.Int64("user_id", input.UserID), logger.Error(err))
		return err
	}

	// Actualizar comisión
	if err := uc.organizerRepo.UpdateCommission(input.UserID, input.Commission); err != nil {
		uc.log.Error("Error updating organizer commission",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return err
	}

	// Log auditoría
	oldCommission := "global_default"
	if profile.CommissionOverride != nil {
		oldCommission = fmt.Sprintf("%.2f", *profile.CommissionOverride)
	}

	newCommission := "global_default"
	if input.Commission != nil {
		newCommission = fmt.Sprintf("%.2f", *input.Commission)
	}

	uc.log.Info("Admin updated organizer commission",
		logger.Int64("admin_id", adminID),
		logger.Int64("user_id", input.UserID),
		logger.String("old_value", oldCommission),
		logger.String("new_value", newCommission),
		logger.String("notes", input.Notes),
		logger.String("action", "admin_update_organizer_commission"))

	return nil
}
