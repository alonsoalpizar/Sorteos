package organizer

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// VerifyOrganizerInput datos de entrada
type VerifyOrganizerInput struct {
	UserID int64
	Notes  string // Notas de verificación
}

// VerifyOrganizerUseCase caso de uso para verificar organizador
type VerifyOrganizerUseCase struct {
	organizerRepo *db.PostgresOrganizerProfileRepository
	log           *logger.Logger
}

// NewVerifyOrganizerUseCase crea una nueva instancia
func NewVerifyOrganizerUseCase(
	organizerRepo *db.PostgresOrganizerProfileRepository,
	log *logger.Logger,
) *VerifyOrganizerUseCase {
	return &VerifyOrganizerUseCase{
		organizerRepo: organizerRepo,
		log:           log,
	}
}

// Execute ejecuta el caso de uso
func (uc *VerifyOrganizerUseCase) Execute(ctx context.Context, input *VerifyOrganizerInput, adminID int64) error {
	// Obtener perfil de organizador
	profile, err := uc.organizerRepo.GetByUserID(input.UserID)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.New("NOT_FOUND", "organizer profile not found", 404, nil)
		}
		uc.log.Error("Error getting organizer profile", logger.Int64("user_id", input.UserID), logger.Error(err))
		return err
	}

	// Validar que el organizador tenga información bancaria completa
	if !profile.HasBankInfo() {
		return errors.New("VALIDATION_FAILED",
			"cannot verify organizer without complete bank information", 400, nil)
	}

	// Verificar organizador
	if err := uc.organizerRepo.Verify(input.UserID, adminID); err != nil {
		uc.log.Error("Error verifying organizer",
			logger.Int64("user_id", input.UserID),
			logger.Int64("admin_id", adminID),
			logger.Error(err))
		return err
	}

	// Log auditoría
	uc.log.Info("Admin verified organizer",
		logger.Int64("admin_id", adminID),
		logger.Int64("user_id", input.UserID),
		logger.String("notes", input.Notes),
		logger.String("action", "admin_verify_organizer"))

	// TODO: Enviar email de notificación al organizador
	// Esto se implementará cuando tengamos el servicio de email configurado

	return nil
}
