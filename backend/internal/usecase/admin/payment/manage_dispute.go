package payment

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ManageDisputeInput datos de entrada
type ManageDisputeInput struct {
	PaymentID       string                 `json:"payment_id"` // UUID o int64 como string
	Action          string                 `json:"action"`     // open, update, close, escalate
	DisputeReason   *string                `json:"dispute_reason,omitempty"`
	DisputeEvidence *string                `json:"dispute_evidence,omitempty"`
	Resolution      *string                `json:"resolution,omitempty"`      // Para cerrar: accepted, rejected, refunded
	AdminNotes      *string                `json:"admin_notes,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`        // Info adicional de Stripe/PayPal
}

// ManageDisputeOutput resultado
type ManageDisputeOutput struct {
	PaymentID        string                 `json:"payment_id"`
	DisputeStatus    string                 `json:"dispute_status"`    // open, under_review, closed
	DisputeReason    string                 `json:"dispute_reason,omitempty"`
	Resolution       string                 `json:"resolution,omitempty"`
	OrganizerID      int64                  `json:"organizer_id"`
	OrganizerEmail   string                 `json:"organizer_email"`
	NotificationSent bool                   `json:"notification_sent"`
	UpdatedAt        string                 `json:"updated_at"`
	Message          string                 `json:"message"`
}

// DisputeMetadata metadata de la disputa
type DisputeMetadata struct {
	Reason          string                 `json:"reason"`
	Evidence        string                 `json:"evidence,omitempty"`
	Resolution      string                 `json:"resolution,omitempty"`
	AdminNotes      string                 `json:"admin_notes,omitempty"`
	OpenedAt        string                 `json:"opened_at"`
	ClosedAt        string                 `json:"closed_at,omitempty"`
	AdditionalData  map[string]interface{} `json:"additional_data,omitempty"`
}

// ManageDisputeUseCase caso de uso para gestionar disputas de pagos
type ManageDisputeUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewManageDisputeUseCase crea una nueva instancia
func NewManageDisputeUseCase(db *gorm.DB, log *logger.Logger) *ManageDisputeUseCase {
	return &ManageDisputeUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ManageDisputeUseCase) Execute(ctx context.Context, input *ManageDisputeInput, adminID int64) (*ManageDisputeOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Buscar payment
	var payment struct {
		ID             string
		UserID         int64
		RaffleID       int64
		Amount         float64
		Status         string
		HasDispute     bool
		DisputeStatus  *string
		DisputeMetadata *string // JSON
		UpdatedAt      time.Time
	}

	result := uc.db.WithContext(ctx).
		Table("payments").
		Select("id, user_id, raffle_id, amount, status, has_dispute, dispute_status, dispute_metadata, updated_at").
		Where("id = ?", input.PaymentID).
		First(&payment)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("PAYMENT_NOT_FOUND", "payment not found", 404, nil)
		}
		uc.log.Error("Error finding payment", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Obtener organizador de la rifa
	var organizer struct {
		UserID int64
		Email  string
	}

	uc.db.WithContext(ctx).
		Table("raffles").
		Select("user_id, (SELECT email FROM users WHERE users.id = raffles.user_id) as email").
		Where("id = ?", payment.RaffleID).
		First(&organizer)

	// Parsear metadata existente
	var disputeMetadata DisputeMetadata
	if payment.DisputeMetadata != nil {
		json.Unmarshal([]byte(*payment.DisputeMetadata), &disputeMetadata)
	}

	// Ejecutar acción según el tipo
	var newDisputeStatus string
	var notificationSent bool
	var err error

	switch input.Action {
	case "open":
		newDisputeStatus, err = uc.openDispute(ctx, &payment, input, &disputeMetadata)
		if err != nil {
			return nil, err
		}
		notificationSent = uc.notifyOrganizer(organizer.UserID, organizer.Email, "dispute_opened", input.DisputeReason)

	case "update":
		newDisputeStatus = "under_review"
		if input.AdminNotes != nil {
			disputeMetadata.AdminNotes = *input.AdminNotes
		}
		if input.DisputeEvidence != nil {
			disputeMetadata.Evidence = *input.DisputeEvidence
		}
		if input.Metadata != nil {
			disputeMetadata.AdditionalData = input.Metadata
		}

	case "close":
		newDisputeStatus, err = uc.closeDispute(ctx, &payment, input, &disputeMetadata)
		if err != nil {
			return nil, err
		}
		notificationSent = uc.notifyOrganizer(organizer.UserID, organizer.Email, "dispute_closed", input.Resolution)

	case "escalate":
		newDisputeStatus = "escalated"
		if input.AdminNotes != nil {
			disputeMetadata.AdminNotes = *input.AdminNotes
		}
		notificationSent = uc.notifyOrganizer(organizer.UserID, organizer.Email, "dispute_escalated", nil)

	default:
		return nil, errors.New("VALIDATION_FAILED", "invalid action", 400, nil)
	}

	// Serializar metadata
	metadataJSON, err := json.Marshal(disputeMetadata)
	if err != nil {
		uc.log.Error("Error marshaling dispute metadata", logger.Error(err))
		return nil, errors.Wrap(errors.ErrInternalServer, err)
	}
	metadataStr := string(metadataJSON)

	// Actualizar payment en DB
	updates := map[string]interface{}{
		"has_dispute":      true,
		"dispute_status":   newDisputeStatus,
		"dispute_metadata": metadataStr,
		"updated_at":       time.Now(),
	}

	result = uc.db.WithContext(ctx).
		Table("payments").
		Where("id = ?", input.PaymentID).
		Updates(updates)

	if result.Error != nil {
		uc.log.Error("Error updating payment dispute", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log auditoría crítica
	uc.log.Error("Admin managed payment dispute",
		logger.Int64("admin_id", adminID),
		logger.String("payment_id", input.PaymentID),
		logger.String("action", input.Action),
		logger.String("dispute_status", newDisputeStatus),
		logger.Int64("organizer_id", organizer.UserID),
		logger.String("action", "admin_manage_dispute"),
		logger.String("severity", "critical"))

	return &ManageDisputeOutput{
		PaymentID:        input.PaymentID,
		DisputeStatus:    newDisputeStatus,
		DisputeReason:    disputeMetadata.Reason,
		Resolution:       disputeMetadata.Resolution,
		OrganizerID:      organizer.UserID,
		OrganizerEmail:   organizer.Email,
		NotificationSent: notificationSent,
		UpdatedAt:        time.Now().Format(time.RFC3339),
		Message:          "Dispute managed successfully",
	}, nil
}

// validateInput valida los datos de entrada
func (uc *ManageDisputeUseCase) validateInput(input *ManageDisputeInput) error {
	if input.PaymentID == "" {
		return errors.New("VALIDATION_FAILED", "payment_id is required", 400, nil)
	}

	validActions := map[string]bool{
		"open":     true,
		"update":   true,
		"close":    true,
		"escalate": true,
	}
	if !validActions[input.Action] {
		return errors.New("VALIDATION_FAILED", "action must be one of: open, update, close, escalate", 400, nil)
	}

	// Validar campos requeridos según acción
	if input.Action == "open" && input.DisputeReason == nil {
		return errors.New("VALIDATION_FAILED", "dispute_reason is required when opening a dispute", 400, nil)
	}

	if input.Action == "close" && input.Resolution == nil {
		return errors.New("VALIDATION_FAILED", "resolution is required when closing a dispute", 400, nil)
	}

	if input.Action == "close" {
		validResolutions := map[string]bool{
			"accepted":  true,
			"rejected":  true,
			"refunded":  true,
		}
		if !validResolutions[*input.Resolution] {
			return errors.New("VALIDATION_FAILED", "resolution must be one of: accepted, rejected, refunded", 400, nil)
		}
	}

	return nil
}

// openDispute abre una nueva disputa
func (uc *ManageDisputeUseCase) openDispute(ctx context.Context, payment interface{}, input *ManageDisputeInput, metadata *DisputeMetadata) (string, error) {
	metadata.Reason = *input.DisputeReason
	if input.DisputeEvidence != nil {
		metadata.Evidence = *input.DisputeEvidence
	}
	if input.AdminNotes != nil {
		metadata.AdminNotes = *input.AdminNotes
	}
	if input.Metadata != nil {
		metadata.AdditionalData = input.Metadata
	}
	metadata.OpenedAt = time.Now().Format(time.RFC3339)

	return "open", nil
}

// closeDispute cierra una disputa
func (uc *ManageDisputeUseCase) closeDispute(ctx context.Context, payment interface{}, input *ManageDisputeInput, metadata *DisputeMetadata) (string, error) {
	metadata.Resolution = *input.Resolution
	if input.AdminNotes != nil {
		metadata.AdminNotes = *input.AdminNotes
	}
	metadata.ClosedAt = time.Now().Format(time.RFC3339)

	// TODO: Si la resolución es "refunded", procesar el refund automáticamente
	// if *input.Resolution == "refunded" {
	//     refundUseCase := NewProcessRefundUseCase(uc.db, uc.log)
	//     refundInput := &ProcessRefundInput{
	//         PaymentID: payment.ID,
	//         Amount: payment.Amount,
	//         Reason: "Dispute resolved with refund",
	//     }
	//     _, err := refundUseCase.Execute(ctx, refundInput, adminID)
	//     if err != nil {
	//         return "", err
	//     }
	// }

	return "closed", nil
}

// notifyOrganizer notifica al organizador sobre la disputa
func (uc *ManageDisputeUseCase) notifyOrganizer(organizerID int64, email, eventType string, details interface{}) bool {
	// TODO: Integrar con sistema de notificaciones
	// emailNotifier.SendEmail(email, subject, body)

	uc.log.Info("Organizer notified about dispute",
		logger.Int64("organizer_id", organizerID),
		logger.String("email", email),
		logger.String("event_type", eventType))

	return true // Simulado como exitoso
}
