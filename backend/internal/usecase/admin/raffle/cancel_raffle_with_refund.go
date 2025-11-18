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

// CancelRaffleWithRefundInput datos de entrada
type CancelRaffleWithRefundInput struct {
	RaffleID int64
	Reason   string // Razón de la cancelación
}

// CancelRaffleWithRefundOutput resultado
type CancelRaffleWithRefundOutput struct {
	RaffleID         int64
	TotalPayments    int
	RefundsInitiated int
	RefundsFailed    int
	TotalRefunded    float64
}

// CancelRaffleWithRefundUseCase caso de uso para cancelar rifa con reembolsos
type CancelRaffleWithRefundUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCancelRaffleWithRefundUseCase crea una nueva instancia
func NewCancelRaffleWithRefundUseCase(db *gorm.DB, log *logger.Logger) *CancelRaffleWithRefundUseCase {
	return &CancelRaffleWithRefundUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *CancelRaffleWithRefundUseCase) Execute(ctx context.Context, input *CancelRaffleWithRefundInput, adminID int64) (*CancelRaffleWithRefundOutput, error) {
	// Validar razón
	if input.Reason == "" {
		return nil, errors.New("VALIDATION_FAILED", "reason is required for cancellation with refund", 400, nil)
	}

	// Obtener rifa
	var raffle domain.Raffle
	if err := uc.db.Where("id = ?", input.RaffleID).First(&raffle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		uc.log.Error("Error finding raffle", logger.Int64("raffle_id", input.RaffleID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Validar que la rifa no esté completed
	if raffle.Status == domain.RaffleStatusCompleted {
		return nil, errors.New("VALIDATION_FAILED", "cannot cancel completed raffle", 400, nil)
	}

	// Validar que la rifa no esté ya cancelled
	if raffle.Status == domain.RaffleStatusCancelled {
		return nil, errors.New("VALIDATION_FAILED", "raffle is already cancelled", 400, nil)
	}

	// Iniciar transacción
	tx := uc.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Obtener todos los pagos succeeded de esta rifa
	var payments []struct {
		ID                  int64
		UserID              int64
		Amount              float64
		StripePaymentIntent *string
		PayPalOrderID       *string
	}

	if err := tx.Table("payments").
		Select("id, user_id, amount, stripe_payment_intent, paypal_order_id").
		Where("raffle_id = ? AND status = ?", input.RaffleID, "succeeded").
		Find(&payments).Error; err != nil {
		tx.Rollback()
		uc.log.Error("Error getting payments for refund", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	output := &CancelRaffleWithRefundOutput{
		RaffleID:      input.RaffleID,
		TotalPayments: len(payments),
	}

	// Procesar reembolsos
	// NOTA: En producción, esto debería integrarse con Stripe/PayPal API
	// Por ahora, solo marcamos los pagos como refunded
	for _, payment := range payments {
		// TODO: Integrar con payment provider real
		// if payment.StripePaymentIntent != nil {
		//     err := stripe.Refund(*payment.StripePaymentIntent, payment.Amount)
		// } else if payment.PayPalOrderID != nil {
		//     err := paypal.Refund(*payment.PayPalOrderID, payment.Amount)
		// }

		// Marcar pago como refunded
		if err := tx.Table("payments").
			Where("id = ?", payment.ID).
			Updates(map[string]interface{}{
				"status":     "refunded",
				"updated_at": time.Now(),
			}).Error; err != nil {
			uc.log.Error("Error marking payment as refunded",
				logger.Int64("payment_id", payment.ID),
				logger.Error(err))
			output.RefundsFailed++
			continue
		}

		output.RefundsInitiated++
		output.TotalRefunded += payment.Amount

		uc.log.Info("Payment marked as refunded",
			logger.Int64("payment_id", payment.ID),
			logger.Int64("user_id", payment.UserID),
			logger.Float64("amount", payment.Amount))
	}

	// Actualizar rifa a cancelled
	now := time.Now()
	updates := map[string]interface{}{
		"status":     domain.RaffleStatusCancelled,
		"updated_at": now,
		"deleted_at": now, // Soft delete
		"admin_notes": fmt.Sprintf(
			"Cancelled by admin ID %d with refunds. Reason: %s. Refunds: %d/%d successful",
			adminID, input.Reason, output.RefundsInitiated, output.TotalPayments,
		),
	}

	if err := tx.Model(&domain.Raffle{}).
		Where("id = ?", input.RaffleID).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		uc.log.Error("Error cancelling raffle", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Liberar números reservados/vendidos
	if err := tx.Exec(`
		UPDATE raffle_numbers
		SET user_id = NULL,
		    reserved_until = NULL,
		    updated_at = ?
		WHERE raffle_id = ?
	`, now, input.RaffleID).Error; err != nil {
		tx.Rollback()
		uc.log.Error("Error releasing raffle numbers", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Commit transacción
	if err := tx.Commit().Error; err != nil {
		uc.log.Error("Error committing transaction", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría crítica
	uc.log.Error("Admin cancelled raffle with refunds",
		logger.Int64("admin_id", adminID),
		logger.Int64("raffle_id", input.RaffleID),
		logger.Int("total_payments", output.TotalPayments),
		logger.Int("refunds_initiated", output.RefundsInitiated),
		logger.Int("refunds_failed", output.RefundsFailed),
		logger.Float64("total_refunded", output.TotalRefunded),
		logger.String("reason", input.Reason),
		logger.String("action", "admin_cancel_raffle_with_refund"),
		logger.String("severity", "critical"))

	// TODO: Enviar emails de notificación
	// - A todos los participantes (confirmación de reembolso)
	// - Al organizador (notificación de cancelación)

	return output, nil
}
