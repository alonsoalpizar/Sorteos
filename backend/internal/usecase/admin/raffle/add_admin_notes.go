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

// AddAdminNotesInput datos de entrada
type AddAdminNotesInput struct {
	RaffleID int64
	Notes    string
	Append   bool // Si true, agrega a las notas existentes; si false, reemplaza
}

// AddAdminNotesUseCase caso de uso para agregar notas administrativas a una rifa
type AddAdminNotesUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewAddAdminNotesUseCase crea una nueva instancia
func NewAddAdminNotesUseCase(db *gorm.DB, log *logger.Logger) *AddAdminNotesUseCase {
	return &AddAdminNotesUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *AddAdminNotesUseCase) Execute(ctx context.Context, input *AddAdminNotesInput, adminID int64) error {
	// Validar notas
	if input.Notes == "" {
		return errors.New("VALIDATION_FAILED", "notes cannot be empty", 400, nil)
	}

	// Obtener rifa para verificar existencia y obtener notas actuales
	var raffle domain.Raffle
	if err := uc.db.Where("id = ?", input.RaffleID).First(&raffle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrRaffleNotFound
		}
		uc.log.Error("Error finding raffle", logger.Int64("raffle_id", input.RaffleID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Construir notas finales
	finalNotes := input.Notes
	if input.Append && raffle.AdminNotes != nil && *raffle.AdminNotes != "" {
		// Agregar timestamp y admin ID a la nueva nota
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		newNote := fmt.Sprintf("\n---\n[%s - Admin ID: %d]\n%s", timestamp, adminID, input.Notes)
		finalNotes = *raffle.AdminNotes + newNote
	} else {
		// Nueva nota con timestamp
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		finalNotes = fmt.Sprintf("[%s - Admin ID: %d]\n%s", timestamp, adminID, input.Notes)
	}

	// Actualizar notas
	if err := uc.db.Model(&domain.Raffle{}).
		Where("id = ?", input.RaffleID).
		Update("admin_notes", finalNotes).Error; err != nil {
		uc.log.Error("Error updating admin notes",
			logger.Int64("raffle_id", input.RaffleID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría
	uc.log.Info("Admin added notes to raffle",
		logger.Int64("admin_id", adminID),
		logger.Int64("raffle_id", input.RaffleID),
		logger.Bool("append", input.Append),
		logger.String("action", "admin_add_raffle_notes"))

	return nil
}
