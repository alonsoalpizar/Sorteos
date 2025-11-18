package settlement

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ApproveSettlementInput datos de entrada
type ApproveSettlementInput struct {
	SettlementID int64
	Notes        string
}

// ApproveSettlementOutput resultado
type ApproveSettlementOutput struct {
	SettlementID  int64
	ApprovedAt    time.Time
	ApprovedBy    int64
	NetAmount     float64
	OrganizerID   int64
	OrganizerName string
	Success       bool
}

// ApproveSettlementUseCase caso de uso para aprobar liquidación
type ApproveSettlementUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewApproveSettlementUseCase crea una nueva instancia
func NewApproveSettlementUseCase(db *gorm.DB, log *logger.Logger) *ApproveSettlementUseCase {
	return &ApproveSettlementUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ApproveSettlementUseCase) Execute(ctx context.Context, input *ApproveSettlementInput, adminID int64) (*ApproveSettlementOutput, error) {
	// Obtener settlement
	var settlement struct {
		ID           int64
		OrganizerID  int64
		RaffleID     int64
		NetAmount    float64
		Status       string
		ApprovedAt   *time.Time
		ApprovedBy   *int64
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

	// Validar que esté en estado pending
	if settlement.Status != "pending" {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("cannot approve settlement with status %s", settlement.Status), 400, nil)
	}

	// Validar que no esté ya aprobado
	if settlement.ApprovedAt != nil {
		return nil, errors.New("VALIDATION_FAILED", "settlement is already approved", 400, nil)
	}

	// Obtener organizador para validaciones
	var organizer struct {
		ID        int64
		Email     string
		FirstName *string
		LastName  *string
		KYCLevel  string
	}

	if err := uc.db.Table("users").
		Where("id = ?", settlement.OrganizerID).
		Scan(&organizer).Error; err != nil {
		uc.log.Error("Error finding organizer", logger.Int64("organizer_id", settlement.OrganizerID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Validar nivel KYC del organizador
	if organizer.KYCLevel != "verified" && organizer.KYCLevel != "enhanced" {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("cannot approve settlement: organizer KYC level is %s, required verified or enhanced", organizer.KYCLevel), 400, nil)
	}

	// Verificar que tenga cuenta bancaria verificada
	var bankAccountCount int64
	uc.db.Table("organizer_bank_accounts").
		Where("user_id = ? AND verified_at IS NOT NULL", settlement.OrganizerID).
		Count(&bankAccountCount)

	if bankAccountCount == 0 {
		return nil, errors.New("VALIDATION_FAILED",
			"cannot approve settlement: organizer has no verified bank account", 400, nil)
	}

	// Actualizar settlement
	now := time.Now()
	updates := map[string]interface{}{
		"status":      "approved",
		"approved_at": now,
		"approved_by": adminID,
		"updated_at":  now,
	}

	if input.Notes != "" {
		// Agregar notas al admin_notes
		var currentNotes string
		uc.db.Table("settlements").
			Select("admin_notes").
			Where("id = ?", input.SettlementID).
			Scan(&currentNotes)

		timestamp := now.Format("2006-01-02 15:04:05")
		newNote := fmt.Sprintf("[%s] Admin ID %d: APPROVED - %s", timestamp, adminID, input.Notes)

		if currentNotes != "" {
			updates["admin_notes"] = currentNotes + "\n---\n" + newNote
		} else {
			updates["admin_notes"] = newNote
		}
	}

	if err := uc.db.Table("settlements").
		Where("id = ?", input.SettlementID).
		Updates(updates).Error; err != nil {
		uc.log.Error("Error approving settlement", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Construir nombre del organizador
	organizerName := organizer.Email
	if organizer.FirstName != nil && organizer.LastName != nil {
		organizerName = *organizer.FirstName + " " + *organizer.LastName
	}

	// Log auditoría crítica
	uc.log.Error("Admin approved settlement",
		logger.Int64("admin_id", adminID),
		logger.Int64("settlement_id", input.SettlementID),
		logger.Int64("organizer_id", settlement.OrganizerID),
		logger.Float64("net_amount", settlement.NetAmount),
		logger.String("action", "admin_approve_settlement"),
		logger.String("severity", "critical"))

	// TODO: Enviar notificación al organizador
	// - Email confirmando aprobación
	// - Detalles del pago pendiente
	// - Próximos pasos

	return &ApproveSettlementOutput{
		SettlementID:  input.SettlementID,
		ApprovedAt:    now,
		ApprovedBy:    adminID,
		NetAmount:     settlement.NetAmount,
		OrganizerID:   settlement.OrganizerID,
		OrganizerName: organizerName,
		Success:       true,
	}, nil
}
