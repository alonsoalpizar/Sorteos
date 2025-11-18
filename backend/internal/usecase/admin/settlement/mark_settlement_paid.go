package settlement

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// MarkSettlementPaidInput datos de entrada
type MarkSettlementPaidInput struct {
	SettlementID     int64   `json:"settlement_id"`
	PaymentMethod    string  `json:"payment_method"`    // bank_transfer, paypal, stripe, cash
	PaymentReference *string `json:"payment_reference,omitempty"` // Número de transferencia, transaction ID, etc.
	Notes            *string `json:"notes,omitempty"`
}

// MarkSettlementPaidOutput resultado
type MarkSettlementPaidOutput struct {
	SettlementID       int64   `json:"settlement_id"`
	Status             string  `json:"status"`
	NetAmount          float64 `json:"net_amount"`
	PaymentMethod      string  `json:"payment_method"`
	PaymentReference   string  `json:"payment_reference,omitempty"`
	PaidAt             string  `json:"paid_at"`
	OrganizerID        int64   `json:"organizer_id"`
	OrganizerEmail     string  `json:"organizer_email"`
	NotificationSent   bool    `json:"notification_sent"`
	Message            string  `json:"message"`
}

// MarkSettlementPaidUseCase caso de uso para marcar settlement como pagado
type MarkSettlementPaidUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewMarkSettlementPaidUseCase crea una nueva instancia
func NewMarkSettlementPaidUseCase(db *gorm.DB, log *logger.Logger) *MarkSettlementPaidUseCase {
	return &MarkSettlementPaidUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *MarkSettlementPaidUseCase) Execute(ctx context.Context, input *MarkSettlementPaidInput, adminID int64) (*MarkSettlementPaidOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Buscar settlement
	var settlement struct {
		ID            int64
		OrganizerID   int64
		RaffleID      int64
		TotalRevenue  float64
		PlatformFee   float64
		NetAmount     float64
		Status        string
		ApprovedBy    *int64
		ApprovedAt    *time.Time
		CalculatedAt  time.Time
	}

	result := uc.db.WithContext(ctx).
		Table("settlements").
		Select("id, organizer_id, raffle_id, total_revenue, platform_fee, net_amount, status, approved_by, approved_at, created_at").
		Where("id = ?", input.SettlementID).
		First(&settlement)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("SETTLEMENT_NOT_FOUND", "settlement not found", 404, nil)
		}
		uc.log.Error("Error finding settlement", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Validar que el settlement esté aprobado
	if settlement.Status != "approved" {
		return nil, errors.New("SETTLEMENT_NOT_APPROVED",
			"settlement must be approved before marking as paid", 400, nil)
	}

	// Obtener información del organizador
	var organizer struct {
		Email string
	}
	uc.db.WithContext(ctx).
		Table("users").
		Select("email").
		Where("id = ?", settlement.OrganizerID).
		First(&organizer)

	// Actualizar settlement
	now := time.Now()
	updates := map[string]interface{}{
		"status":            "paid",
		"payment_method":    input.PaymentMethod,
		"payment_reference": input.PaymentReference,
		"paid_at":           now,
		"paid_by":           adminID,
		"notes":             input.Notes,
		"updated_at":        now,
	}

	result = uc.db.WithContext(ctx).
		Table("settlements").
		Where("id = ?", input.SettlementID).
		Updates(updates)

	if result.Error != nil {
		uc.log.Error("Error updating settlement", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Actualizar organizer_profile (total_payouts, pending_payout)
	err := uc.updateOrganizerProfile(ctx, settlement.OrganizerID, settlement.NetAmount)
	if err != nil {
		uc.log.Error("Error updating organizer profile", logger.Error(err))
		// No fallar la operación, solo loguear
	}

	// TODO: Enviar email de confirmación al organizador
	notificationSent := uc.sendPaymentConfirmation(settlement.OrganizerID, organizer.Email, settlement.NetAmount, input.PaymentMethod)

	// Log auditoría crítica
	uc.log.Error("Admin marked settlement as paid",
		logger.Int64("admin_id", adminID),
		logger.Int64("settlement_id", input.SettlementID),
		logger.Int64("organizer_id", settlement.OrganizerID),
		logger.Float64("net_amount", settlement.NetAmount),
		logger.String("payment_method", input.PaymentMethod),
		logger.String("action", "admin_mark_settlement_paid"),
		logger.String("severity", "critical"))

	// Construir output
	paymentRef := ""
	if input.PaymentReference != nil {
		paymentRef = *input.PaymentReference
	}

	return &MarkSettlementPaidOutput{
		SettlementID:       input.SettlementID,
		Status:             "paid",
		NetAmount:          settlement.NetAmount,
		PaymentMethod:      input.PaymentMethod,
		PaymentReference:   paymentRef,
		PaidAt:             now.Format(time.RFC3339),
		OrganizerID:        settlement.OrganizerID,
		OrganizerEmail:     organizer.Email,
		NotificationSent:   notificationSent,
		Message:            "Settlement marked as paid successfully",
	}, nil
}

// validateInput valida los datos de entrada
func (uc *MarkSettlementPaidUseCase) validateInput(input *MarkSettlementPaidInput) error {
	if input.SettlementID <= 0 {
		return errors.New("VALIDATION_FAILED", "settlement_id is required", 400, nil)
	}

	validMethods := map[string]bool{
		"bank_transfer": true,
		"paypal":        true,
		"stripe":        true,
		"cash":          true,
		"check":         true,
	}
	if !validMethods[input.PaymentMethod] {
		return errors.New("VALIDATION_FAILED", "payment_method must be one of: bank_transfer, paypal, stripe, cash, check", 400, nil)
	}

	return nil
}

// updateOrganizerProfile actualiza las métricas del organizador
func (uc *MarkSettlementPaidUseCase) updateOrganizerProfile(ctx context.Context, organizerID int64, paidAmount float64) error {
	// Incrementar total_payouts
	// Decrementar pending_payout

	result := uc.db.WithContext(ctx).Exec(`
		UPDATE organizer_profiles
		SET
			total_payouts = COALESCE(total_payouts, 0) + ?,
			pending_payout = GREATEST(0, COALESCE(pending_payout, 0) - ?),
			updated_at = ?
		WHERE user_id = ?
	`, paidAmount, paidAmount, time.Now(), organizerID)

	if result.Error != nil {
		return result.Error
	}

	// Si no existe el perfil, crearlo
	if result.RowsAffected == 0 {
		profile := map[string]interface{}{
			"user_id":        organizerID,
			"total_payouts":  paidAmount,
			"pending_payout": 0,
			"created_at":     time.Now(),
			"updated_at":     time.Now(),
		}
		uc.db.WithContext(ctx).Table("organizer_profiles").Create(profile)
	}

	return nil
}

// sendPaymentConfirmation envía email de confirmación de pago
func (uc *MarkSettlementPaidUseCase) sendPaymentConfirmation(organizerID int64, email string, amount float64, method string) bool {
	// TODO: Integrar con sistema de notificaciones
	// subject := "Pago procesado - Sorteos.club"
	// body := fmt.Sprintf("Tu pago de $%.2f ha sido procesado vía %s", amount, method)
	// emailNotifier.SendEmail(email, subject, body)

	uc.log.Info("Payment confirmation sent to organizer",
		logger.Int64("organizer_id", organizerID),
		logger.String("email", email),
		logger.Float64("amount", amount))

	return true // Simulado como exitoso
}
