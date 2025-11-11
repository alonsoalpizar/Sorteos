package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain/entities"
	"github.com/sorteos-platform/backend/internal/domain/repositories"
)

// PostgresPaymentRepository implements PaymentRepository using PostgreSQL
type PostgresPaymentRepository struct {
	db *gorm.DB
}

// NewPostgresPaymentRepository creates a new PostgreSQL payment repository
func NewPostgresPaymentRepository(db *gorm.DB) repositories.PaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

// Create stores a new payment
func (r *PostgresPaymentRepository) Create(ctx context.Context, payment *entities.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// FindByID retrieves a payment by ID
func (r *PostgresPaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// FindByReservationID retrieves a payment by reservation ID
func (r *PostgresPaymentRepository) FindByReservationID(ctx context.Context, reservationID uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("reservation_id = ?", reservationID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// FindByStripePaymentIntentID retrieves a payment by Stripe Payment Intent ID
func (r *PostgresPaymentRepository) FindByStripePaymentIntentID(ctx context.Context, intentID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.WithContext(ctx).Where("stripe_payment_intent_id = ?", intentID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// FindByUserID retrieves all payments for a user
func (r *PostgresPaymentRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// FindByRaffleID retrieves all payments for a raffle
func (r *PostgresPaymentRepository) FindByRaffleID(ctx context.Context, raffleID uuid.UUID) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	err := r.db.WithContext(ctx).
		Where("raffle_id = ?", raffleID).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// Update updates an existing payment
func (r *PostgresPaymentRepository) Update(ctx context.Context, payment *entities.Payment) error {
	payment.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(payment).Error
}

// Delete removes a payment
func (r *PostgresPaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Payment{}, "id = ?", id).Error
}
