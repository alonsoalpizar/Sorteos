package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ProcessRefundInput datos de entrada
// NOTA: PaymentID es UUID string porque la tabla payments usa UUIDs como PK
type ProcessRefundInput struct {
	PaymentID string   // UUID del payment
	Reason    string  // Razón del reembolso
	Amount    *float64 // Si es parcial, especificar amount. Si nil, es total
	Notes     string  // Notas adicionales
}

// ProcessRefundOutput resultado
type ProcessRefundOutput struct {
	PaymentID     string
	RefundAmount  float64
	RefundType    string // "full" o "partial"
	Success       bool
	FailureReason string
}

// ProcessRefundUseCase caso de uso para procesar reembolsos
type ProcessRefundUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewProcessRefundUseCase crea una nueva instancia
func NewProcessRefundUseCase(db *gorm.DB, log *logger.Logger) *ProcessRefundUseCase {
	return &ProcessRefundUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ProcessRefundUseCase) Execute(ctx context.Context, input *ProcessRefundInput, adminID int64) (*ProcessRefundOutput, error) {
	// Validar razón
	if input.Reason == "" {
		return nil, errors.New("VALIDATION_FAILED", "reason is required for refund", 400, nil)
	}

	// Obtener pago
	var payment Payment
	if err := uc.db.Table("payments").Where("id = ?", input.PaymentID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("PAYMENT_NOT_FOUND", "payment not found", 404, nil)
		}
		uc.log.Error("Error finding payment", logger.String("payment_id", input.PaymentID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Validar que el pago esté succeeded
	if payment.Status != "succeeded" {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("cannot refund payment with status %s", payment.Status), 400, nil)
	}

	// Validar que no esté ya refunded
	if payment.Status == "refunded" {
		return nil, errors.New("VALIDATION_FAILED", "payment is already refunded", 400, nil)
	}

	// Determinar tipo de refund y amount
	refundAmount := payment.Amount
	refundType := "full"

	if input.Amount != nil {
		if *input.Amount <= 0 || *input.Amount > payment.Amount {
			return nil, errors.New("VALIDATION_FAILED",
				fmt.Sprintf("invalid refund amount: must be between 0 and %.2f", payment.Amount), 400, nil)
		}
		refundAmount = *input.Amount
		refundType = "partial"
	}

	output := &ProcessRefundOutput{
		PaymentID:    input.PaymentID,
		RefundAmount: refundAmount,
		RefundType:   refundType,
	}

	// TODO: Integrar con payment provider real
	// Por ahora, simulamos el proceso
	var refundError error

	// Simular llamada a Stripe/PayPal
	// if payment.StripePaymentIntent != nil {
	//     refundError = uc.stripeService.Refund(*payment.StripePaymentIntent, refundAmount)
	// } else if payment.PayPalOrderID != nil {
	//     refundError = uc.paypalService.Refund(*payment.PayPalOrderID, refundAmount)
	// }

	if refundError != nil {
		output.Success = false
		output.FailureReason = refundError.Error()

		uc.log.Error("Refund failed",
			logger.String("payment_id", input.PaymentID),
			logger.Float64("amount", refundAmount),
			logger.Error(refundError))

		return output, nil
	}

	// Iniciar transacción
	tx := uc.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Actualizar pago
	now := time.Now()
	updates := map[string]interface{}{
		"status":      "refunded",
		"refunded_at": now,
		"refunded_by": adminID,
		"updated_at":  now,
	}

	// Guardar razón y notas en metadata o admin_notes
	adminNotes := fmt.Sprintf("Refunded by admin ID %d. Reason: %s", adminID, input.Reason)
	if input.Notes != "" {
		adminNotes += fmt.Sprintf(". Notes: %s", input.Notes)
	}
	if refundType == "partial" {
		adminNotes += fmt.Sprintf(". Partial refund: $%.2f of $%.2f", refundAmount, payment.Amount)
	}
	updates["admin_notes"] = adminNotes

	if err := tx.Table("payments").
		Where("id = ?", input.PaymentID).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		uc.log.Error("Error updating payment status", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Si es refund completo, liberar números
	// NOTA: raffle_numbers usa UUIDs para raffle_id y user_id
	if refundType == "full" {
		// Contar cuántos números vamos a liberar para actualizar contadores
		var numbersCount int64
		tx.Table("raffle_numbers").
			Where("raffle_id = ? AND user_id = ?", payment.RaffleID, payment.UserID).
			Count(&numbersCount)

		// Liberar números asociados a este pago
		if err := tx.Table("raffle_numbers").
			Where("raffle_id = ? AND user_id = ?", payment.RaffleID, payment.UserID).
			Updates(map[string]interface{}{
				"user_id":        nil,
				"reserved_until": nil,
				"updated_at":     now,
			}).Error; err != nil {
			tx.Rollback()
			uc.log.Error("Error releasing numbers", logger.Error(err))
			return nil, errors.Wrap(errors.ErrDatabaseError, err)
		}

		// Actualizar contadores en raffle
		if numbersCount > 0 {
			if err := tx.Table("raffles").
				Where("uuid::text = ?", payment.RaffleID).
				Updates(map[string]interface{}{
					"sold_count": gorm.Expr("sold_count - ?", numbersCount),
					"updated_at": now,
				}).Error; err != nil {
				tx.Rollback()
				uc.log.Error("Error updating raffle counters", logger.Error(err))
				return nil, errors.Wrap(errors.ErrDatabaseError, err)
			}
		}
	}

	// Commit transacción
	if err := tx.Commit().Error; err != nil {
		uc.log.Error("Error committing transaction", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	output.Success = true

	// Log auditoría crítica
	uc.log.Error("Admin processed refund",
		logger.Int64("admin_id", adminID),
		logger.String("payment_id", input.PaymentID),
		logger.String("user_id", payment.UserID),
		logger.String("raffle_id", payment.RaffleID),
		logger.Float64("amount", refundAmount),
		logger.String("type", refundType),
		logger.String("reason", input.Reason),
		logger.String("action", "admin_process_refund"),
		logger.String("severity", "critical"))

	// TODO: Enviar email de notificación al usuario
	// - Confirmación de reembolso
	// - Timeline y detalles

	return output, nil
}
