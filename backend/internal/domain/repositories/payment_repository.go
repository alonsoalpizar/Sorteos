package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/sorteos-platform/backend/internal/domain/entities"
)

// PaymentRepository defines the interface for payment persistence
type PaymentRepository interface {
	// Create stores a new payment
	Create(ctx context.Context, payment *entities.Payment) error

	// FindByID retrieves a payment by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error)

	// FindByReservationID retrieves a payment by reservation ID
	FindByReservationID(ctx context.Context, reservationID uuid.UUID) (*entities.Payment, error)

	// FindByStripePaymentIntentID retrieves a payment by Stripe Payment Intent ID
	FindByStripePaymentIntentID(ctx context.Context, intentID string) (*entities.Payment, error)

	// FindByUserID retrieves all payments for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Payment, error)

	// FindByRaffleID retrieves all payments for a raffle
	FindByRaffleID(ctx context.Context, raffleID uuid.UUID) ([]*entities.Payment, error)

	// Update updates an existing payment
	Update(ctx context.Context, payment *entities.Payment) error

	// Delete removes a payment
	Delete(ctx context.Context, id uuid.UUID) error
}
