package settlement

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// RejectSettlementInput datos de entrada
type RejectSettlementInput struct {
	SettlementID int64
	Reason       string // Razón obligatoria del rechazo
	Notes        string // Notas adicionales
}

// RejectSettlementOutput resultado
type RejectSettlementOutput struct {
	SettlementID  int64
	RejectedAt    time.Time
	RejectedBy    int64
	Reason        string
	OrganizerID   int64
	OrganizerName string
	Success       bool
}

// RejectSettlementUseCase caso de uso para rechazar liquidación
type RejectSettlementUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewRejectSettlementUseCase crea una nueva instancia
func NewRejectSettlementUseCase(db *gorm.DB, log *logger.Logger) *RejectSettlementUseCase {
	return &RejectSettlementUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *RejectSettlementUseCase) Execute(ctx context.Context, input *RejectSettlementInput, adminID int64) (*RejectSettlementOutput, error) {
	// Validar razón
	if input.Reason == "" {
		return nil, errors.New("VALIDATION_FAILED", "reason is required for rejection", 400, nil)
	}

	// Obtener settlement
	var settlement struct {
		ID          int64
		OrganizerID int64
		RaffleID    int64
		NetAmount   float64
		Status      string
		RejectedAt  *time.Time
		RejectedBy  *int64
	}

	if err := uc.db.Table("settlements").
		Where("id = ?", input.SettlementID).
		Scan(&settlement).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("SETTLEMENT_NOT_FOUND", "settlement not found", 404, nil)
		}
		uc.log.Error("Error finding settlement", logger.Int64("settlement_id", input.SettlementID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Validar que esté en estado pending o approved (se puede rechazar después de aprobar)
	if settlement.Status != "pending" && settlement.Status != "approved" {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("cannot reject settlement with status %s", settlement.Status), 400, nil)
	}

	// Validar que no esté ya rechazado
	if settlement.Status == "rejected" {
		return nil, errors.New("VALIDATION_FAILED", "settlement is already rejected", 400, nil)
	}

	// Validar que no esté pagado
	if settlement.Status == "paid" {
		return nil, errors.New("VALIDATION_FAILED", "cannot reject paid settlement", 400, nil)
	}

	// Obtener organizador
	var organizer struct {
		ID        int64
		Email     string
		FirstName *string
		LastName  *string
	}

	if err := uc.db.Table("users").
		Where("id = ?", settlement.OrganizerID).
		Scan(&organizer).Error; err != nil {
		uc.log.Error("Error finding organizer", logger.Int64("organizer_id", settlement.OrganizerID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Actualizar settlement
	now := time.Now()
	updates := map[string]interface{}{
		"status":           "rejected",
		"rejected_at":      now,
		"rejected_by":      adminID,
		"rejection_reason": input.Reason,
		"updated_at":       now,
	}

	// Agregar notas al admin_notes
	var currentNotes string
	uc.db.Table("settlements").
		Select("admin_notes").
		Where("id = ?", input.SettlementID).
		Scan(&currentNotes)

	timestamp := now.Format("2006-01-02 15:04:05")
	newNote := fmt.Sprintf("[%s] Admin ID %d: REJECTED - Reason: %s", timestamp, adminID, input.Reason)
	if input.Notes != "" {
		newNote += fmt.Sprintf(". Notes: %s", input.Notes)
	}

	if currentNotes != "" {
		updates["admin_notes"] = currentNotes + "\n---\n" + newNote
	} else {
		updates["admin_notes"] = newNote
	}

	if err := uc.db.Table("settlements").
		Where("id = ?", input.SettlementID).
		Updates(updates).Error; err != nil {
		uc.log.Error("Error rejecting settlement", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Construir nombre del organizador
	organizerName := organizer.Email
	if organizer.FirstName != nil && organizer.LastName != nil {
		organizerName = *organizer.FirstName + " " + *organizer.LastName
	}

	// Log auditoría crítica
	uc.log.Error("Admin rejected settlement",
		logger.Int64("admin_id", adminID),
		logger.Int64("settlement_id", input.SettlementID),
		logger.Int64("organizer_id", settlement.OrganizerID),
		logger.Float64("net_amount", settlement.NetAmount),
		logger.String("reason", input.Reason),
		logger.String("action", "admin_reject_settlement"),
		logger.String("severity", "critical"))

	// TODO: Enviar notificación al organizador
	// - Email notificando rechazo
	// - Razón del rechazo
	// - Acciones a tomar para corregir

	return &RejectSettlementOutput{
		SettlementID:  input.SettlementID,
		RejectedAt:    now,
		RejectedBy:    adminID,
		Reason:        input.Reason,
		OrganizerID:   settlement.OrganizerID,
		OrganizerName: organizerName,
		Success:       true,
	}, nil
}
