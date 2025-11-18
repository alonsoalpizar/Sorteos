package settlement

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ProcessPayoutInput datos de entrada
type ProcessPayoutInput struct {
	SettlementID     int64
	PaymentReference string  // Número de transacción bancaria
	PaymentMethod    string  // wire_transfer, paypal, stripe_connect, etc.
	Notes            string  // Notas adicionales
	PaidAmount       *float64 // Si es diferente al net_amount, especificar
}

// ProcessPayoutOutput resultado
type ProcessPayoutOutput struct {
	SettlementID     int64
	PaidAt           time.Time
	PaymentReference string
	PaymentMethod    string
	PaidAmount       float64
	NetAmount        float64
	OrganizerID      int64
	OrganizerName    string
	Success          bool
	FailureReason    string
}

// ProcessPayoutUseCase caso de uso para procesar pago de liquidación
type ProcessPayoutUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewProcessPayoutUseCase crea una nueva instancia
func NewProcessPayoutUseCase(db *gorm.DB, log *logger.Logger) *ProcessPayoutUseCase {
	return &ProcessPayoutUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ProcessPayoutUseCase) Execute(ctx context.Context, input *ProcessPayoutInput, adminID int64) (*ProcessPayoutOutput, error) {
	// Validar campos obligatorios
	if input.PaymentReference == "" {
		return nil, errors.New("VALIDATION_FAILED", "payment_reference is required", 400, nil)
	}
	if input.PaymentMethod == "" {
		return nil, errors.New("VALIDATION_FAILED", "payment_method is required", 400, nil)
	}

	// Validar payment method
	validMethods := map[string]bool{
		"wire_transfer":   true,
		"ach":             true,
		"paypal":          true,
		"stripe_connect":  true,
		"manual":          true,
	}
	if !validMethods[input.PaymentMethod] {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("invalid payment_method: %s", input.PaymentMethod), 400, nil)
	}

	// Obtener settlement
	var settlement struct {
		ID          int64
		OrganizerID int64
		RaffleID    int64
		NetAmount   float64
		Status      string
		PaidAt      *time.Time
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

	// Validar que esté aprobado
	if settlement.Status != "approved" {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("cannot process payout for settlement with status %s, must be approved", settlement.Status), 400, nil)
	}

	// Validar que no esté ya pagado
	if settlement.PaidAt != nil {
		return nil, errors.New("VALIDATION_FAILED", "settlement is already paid", 400, nil)
	}

	// Determinar paid amount
	paidAmount := settlement.NetAmount
	if input.PaidAmount != nil {
		if *input.PaidAmount <= 0 {
			return nil, errors.New("VALIDATION_FAILED", "paid_amount must be greater than 0", 400, nil)
		}
		paidAmount = *input.PaidAmount

		// Advertencia si difiere del net_amount
		if paidAmount != settlement.NetAmount {
			uc.log.Error("WARNING: Paid amount differs from net_amount",
				logger.Int64("settlement_id", input.SettlementID),
				logger.Float64("net_amount", settlement.NetAmount),
				logger.Float64("paid_amount", paidAmount),
				logger.String("severity", "warning"))
		}
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

	// Verificar que tenga cuenta bancaria (seguridad adicional)
	var bankAccountCount int64
	uc.db.Table("organizer_bank_accounts").
		Where("user_id = ? AND verified_at IS NOT NULL", settlement.OrganizerID).
		Count(&bankAccountCount)

	if bankAccountCount == 0 {
		return nil, errors.New("VALIDATION_FAILED",
			"cannot process payout: organizer has no verified bank account", 400, nil)
	}

	output := &ProcessPayoutOutput{
		SettlementID:     input.SettlementID,
		PaymentReference: input.PaymentReference,
		PaymentMethod:    input.PaymentMethod,
		PaidAmount:       paidAmount,
		NetAmount:        settlement.NetAmount,
		OrganizerID:      settlement.OrganizerID,
	}

	// TODO: Integrar con payment provider real
	// - Si es stripe_connect: hacer transfer a connected account
	// - Si es paypal: hacer mass payment
	// - Si es wire_transfer/ach: validar con banco (o solo registrar)
	var payoutError error

	// Simulación de pago
	// if input.PaymentMethod == "stripe_connect" {
	//     payoutError = uc.stripeService.Transfer(organizerStripeID, paidAmount, input.PaymentReference)
	// } else if input.PaymentMethod == "paypal" {
	//     payoutError = uc.paypalService.MassPayout(organizerPaypalEmail, paidAmount, input.PaymentReference)
	// }

	if payoutError != nil {
		output.Success = false
		output.FailureReason = payoutError.Error()

		uc.log.Error("Payout failed",
			logger.Int64("settlement_id", input.SettlementID),
			logger.Float64("amount", paidAmount),
			logger.String("method", input.PaymentMethod),
			logger.Error(payoutError))

		return output, nil
	}

	// Actualizar settlement
	now := time.Now()
	updates := map[string]interface{}{
		"status":            "paid",
		"paid_at":           now,
		"payment_reference": input.PaymentReference,
		"payment_method":    input.PaymentMethod,
		"updated_at":        now,
	}

	// Agregar notas al admin_notes
	var currentNotes string
	uc.db.Table("settlements").
		Select("admin_notes").
		Where("id = ?", input.SettlementID).
		Scan(&currentNotes)

	timestamp := now.Format("2006-01-02 15:04:05")
	newNote := fmt.Sprintf("[%s] Admin ID %d: PAID - Method: %s, Reference: %s, Amount: $%.2f",
		timestamp, adminID, input.PaymentMethod, input.PaymentReference, paidAmount)

	if paidAmount != settlement.NetAmount {
		newNote += fmt.Sprintf(" (Net amount was $%.2f)", settlement.NetAmount)
	}

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
		uc.log.Error("Error updating settlement as paid", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Construir nombre del organizador
	organizerName := organizer.Email
	if organizer.FirstName != nil && organizer.LastName != nil {
		organizerName = *organizer.FirstName + " " + *organizer.LastName
	}
	output.OrganizerName = organizerName
	output.Success = true

	now = time.Now()
	output.PaidAt = now

	// Log auditoría crítica
	uc.log.Error("Admin processed payout",
		logger.Int64("admin_id", adminID),
		logger.Int64("settlement_id", input.SettlementID),
		logger.Int64("organizer_id", settlement.OrganizerID),
		logger.Float64("amount", paidAmount),
		logger.String("method", input.PaymentMethod),
		logger.String("reference", input.PaymentReference),
		logger.String("action", "admin_process_payout"),
		logger.String("severity", "critical"))

	// TODO: Enviar notificación al organizador
	// - Email confirmando pago procesado
	// - Detalles del pago (reference, método, monto)
	// - Comprobante/recibo

	return output, nil
}
