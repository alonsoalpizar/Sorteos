package payment

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PaymentDetailEvent evento en el timeline del pago
type PaymentDetailEvent struct {
	Type      string    `json:"type"` // created, webhook, status_change, refund, note
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentFullDetails detalles completos del pago
type PaymentFullDetails struct {
	Payment         *Payment              `json:"payment"`
	User            *domain.User          `json:"user"`
	Raffle          *domain.Raffle        `json:"raffle"`
	Organizer       *domain.User          `json:"organizer"`
	Numbers         []string              `json:"numbers"` // Números comprados
	Timeline        []*PaymentDetailEvent `json:"timeline"`
	RefundHistory   []*RefundRecord       `json:"refund_history,omitempty"`
	WebhookEvents   []*WebhookEvent       `json:"webhook_events,omitempty"`
}

// RefundRecord registro de reembolso
type RefundRecord struct {
	RefundedAt time.Time `json:"refunded_at"`
	RefundedBy int64     `json:"refunded_by"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"` // full, partial
	Reason     string    `json:"reason"`
	Notes      string    `json:"notes"`
}

// WebhookEvent evento de webhook
type WebhookEvent struct {
	ReceivedAt time.Time              `json:"received_at"`
	Provider   string                 `json:"provider"`
	EventType  string                 `json:"event_type"`
	Status     string                 `json:"status"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// ViewPaymentDetailsUseCase caso de uso para ver detalles completos de un pago
type ViewPaymentDetailsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewViewPaymentDetailsUseCase crea una nueva instancia
func NewViewPaymentDetailsUseCase(db *gorm.DB, log *logger.Logger) *ViewPaymentDetailsUseCase {
	return &ViewPaymentDetailsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
// NOTA: paymentID es UUID string porque la tabla payments usa UUIDs
func (uc *ViewPaymentDetailsUseCase) Execute(ctx context.Context, paymentID string, adminID int64) (*PaymentFullDetails, error) {
	// Obtener pago
	var payment Payment
	if err := uc.db.Table("payments").Where("id = ?", paymentID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("PAYMENT_NOT_FOUND", "payment not found", 404, nil)
		}
		uc.log.Error("Error finding payment", logger.String("payment_id", paymentID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	details := &PaymentFullDetails{
		Payment:  &payment,
		Timeline: make([]*PaymentDetailEvent, 0),
	}

	// Obtener usuario (payment.UserID es UUID que referencia users.uuid)
	var user domain.User
	if err := uc.db.Where("uuid = ?", payment.UserID).First(&user).Error; err == nil {
		details.User = &user
	}

	// Obtener raffle y organizador (payment.RaffleID es UUID que referencia raffles.uuid)
	var raffle domain.Raffle
	if err := uc.db.Where("uuid = ?", payment.RaffleID).First(&raffle).Error; err == nil {
		details.Raffle = &raffle

		// Obtener organizador
		var organizer domain.User
		if err := uc.db.Where("id = ?", raffle.UserID).First(&organizer).Error; err == nil {
			details.Organizer = &organizer
		}
	}

	// Obtener números comprados
	type RaffleNumberResult struct {
		Number string
	}
	var raffleNumbers []RaffleNumberResult
	if err := uc.db.Table("raffle_numbers").
		Select("number").
		Where("raffle_id = ? AND user_id = ?", payment.RaffleID, payment.UserID).
		Find(&raffleNumbers).Error; err == nil {
		numbers := make([]string, 0, len(raffleNumbers))
		for _, rn := range raffleNumbers {
			numbers = append(numbers, rn.Number)
		}
		details.Numbers = numbers
	}

	// Construir timeline

	// 1. Evento de creación
	details.Timeline = append(details.Timeline, &PaymentDetailEvent{
		Type:      "created",
		Timestamp: payment.CreatedAt,
		Details:   "Payment created",
		Metadata: map[string]interface{}{
			"amount":   payment.Amount,
			"provider": payment.Provider,
		},
	})

	// 2. Cambios de estado (si hay audit logs)
	var auditLogs []struct {
		Action    string
		Details   string
		CreatedAt time.Time
		UserID    *int64
	}

	if err := uc.db.Table("audit_logs").
		Select("action, details, created_at, user_id").
		Where("entity_type = ? AND entity_id = ?", "payment", paymentID).
		Order("created_at ASC").
		Find(&auditLogs).Error; err == nil {
		for _, log := range auditLogs {
			details.Timeline = append(details.Timeline, &PaymentDetailEvent{
				Type:      "status_change",
				Timestamp: log.CreatedAt,
				Details:   log.Details,
			})
		}
	}

	// 3. Si hay webhook events (asumiendo tabla payment_webhook_events)
	var webhookEvents []struct {
		ReceivedAt time.Time
		Provider   string
		EventType  string
		Status     string
		EventData  map[string]interface{}
	}

	if err := uc.db.Table("payment_webhook_events").
		Select("received_at, provider, event_type, status, event_data").
		Where("payment_id = ?", paymentID).
		Order("received_at ASC").
		Find(&webhookEvents).Error; err == nil {

		details.WebhookEvents = make([]*WebhookEvent, 0, len(webhookEvents))
		for _, we := range webhookEvents {
			details.WebhookEvents = append(details.WebhookEvents, &WebhookEvent{
				ReceivedAt: we.ReceivedAt,
				Provider:   we.Provider,
				EventType:  we.EventType,
				Status:     we.Status,
				Data:       we.EventData,
			})

			// Agregar a timeline
			details.Timeline = append(details.Timeline, &PaymentDetailEvent{
				Type:      "webhook",
				Timestamp: we.ReceivedAt,
				Details:   we.EventType,
				Metadata: map[string]interface{}{
					"provider": we.Provider,
					"status":   we.Status,
				},
			})
		}
	}

	// 4. Si fue refunded
	if payment.Status == "refunded" && payment.RefundedAt != nil {
		refundRecord := &RefundRecord{
			RefundedAt: *payment.RefundedAt,
			Amount:     payment.Amount,
			Type:       "full",
		}

		// Extraer razón de admin_notes si existe
		if payment.AdminNotes != "" {
			refundRecord.Notes = payment.AdminNotes
		}

		details.RefundHistory = []*RefundRecord{refundRecord}

		// Agregar a timeline
		details.Timeline = append(details.Timeline, &PaymentDetailEvent{
			Type:      "refund",
			Timestamp: *payment.RefundedAt,
			Details:   "Payment refunded",
			Metadata: map[string]interface{}{
				"amount": payment.Amount,
				"type":   "full",
			},
		})
	}

	// Log auditoría
	uc.log.Info("Admin viewed payment details",
		logger.Int64("admin_id", adminID),
		logger.String("payment_id", paymentID),
		logger.String("action", "admin_view_payment_details"))

	return details, nil
}
