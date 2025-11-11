package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	dbadapter "github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain/entities"
	"github.com/sorteos-platform/backend/internal/domain/repositories"
	"github.com/sorteos-platform/backend/internal/infrastructure/payment"
)

var (
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrReservationNotPending   = errors.New("reservation is not pending")
	ErrPaymentAlreadyExists    = errors.New("payment already exists for this reservation")
	ErrIdempotencyKeyMismatch  = errors.New("idempotency key mismatch: different request with same key")
	ErrPaymentIntentNotFound   = errors.New("payment intent not found")
	ErrPaymentNotFound         = errors.New("payment not found")
)

// PaymentUseCases handles business logic for payments
type PaymentUseCases struct {
	paymentRepo         repositories.PaymentRepository
	reservationRepo     repositories.ReservationRepository
	raffleRepo          dbadapter.RaffleRepository
	idempotencyKeyRepo  repositories.IdempotencyKeyRepository
	paymentProvider     payment.PaymentProvider
	reservationUseCases *ReservationUseCases
}

// NewPaymentUseCases creates a new payment use cases instance
func NewPaymentUseCases(
	paymentRepo repositories.PaymentRepository,
	reservationRepo repositories.ReservationRepository,
	raffleRepo dbadapter.RaffleRepository,
	idempotencyKeyRepo repositories.IdempotencyKeyRepository,
	paymentProvider payment.PaymentProvider,
	reservationUseCases *ReservationUseCases,
) *PaymentUseCases {
	return &PaymentUseCases{
		paymentRepo:         paymentRepo,
		reservationRepo:     reservationRepo,
		raffleRepo:          raffleRepo,
		idempotencyKeyRepo:  idempotencyKeyRepo,
		paymentProvider:     paymentProvider,
		reservationUseCases: reservationUseCases,
	}
}

// CreatePaymentIntentInput represents input for creating a payment intent
type CreatePaymentIntentInput struct {
	ReservationID   uuid.UUID
	UserID          uuid.UUID
	IdempotencyKey  string
}

// CreatePaymentIntentOutput represents the output of creating a payment intent
type CreatePaymentIntentOutput struct {
	PaymentID    uuid.UUID
	ClientSecret string
	Amount       float64
	Currency     string
}

// CreatePaymentIntent creates a Stripe payment intent for a reservation
func (uc *PaymentUseCases) CreatePaymentIntent(ctx context.Context, input CreatePaymentIntentInput) (*CreatePaymentIntentOutput, error) {
	// 1. Check idempotency key
	if input.IdempotencyKey != "" {
		existingKey, err := uc.idempotencyKeyRepo.FindByKey(ctx, input.IdempotencyKey, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("error checking idempotency key: %w", err)
		}

		if existingKey != nil {
			// Request already processed
			if existingKey.Status == entities.IdempotencyKeyStatusCompleted {
				// Return cached response
				var cachedOutput CreatePaymentIntentOutput
				if err := json.Unmarshal([]byte(existingKey.ResponseBody), &cachedOutput); err == nil {
					return &cachedOutput, nil
				}
			}

			// Verify request matches
			if err := existingKey.VerifyRequestMatch("/payments/intent", input); err != nil {
				return nil, ErrIdempotencyKeyMismatch
			}
		}
	}

	// 2. Get reservation
	reservation, err := uc.reservationRepo.FindByID(ctx, input.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("error fetching reservation: %w", err)
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}

	// 3. Verify reservation belongs to user
	if reservation.UserID != input.UserID {
		return nil, errors.New("reservation does not belong to user")
	}

	// 4. Verify reservation can be paid
	if err := reservation.CanBePaid(); err != nil {
		return nil, err
	}

	// 5. Check if payment already exists
	existingPayment, err := uc.paymentRepo.FindByReservationID(ctx, input.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing payment: %w", err)
	}
	if existingPayment != nil {
		// Payment already exists, return client secret
		return &CreatePaymentIntentOutput{
			PaymentID:    existingPayment.ID,
			ClientSecret: existingPayment.StripeClientSecret,
			Amount:       existingPayment.Amount,
			Currency:     existingPayment.Currency,
		}, nil
	}

	// 6. Get raffle details for metadata
	raffle, err := uc.raffleRepo.FindByUUID(reservation.RaffleID.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching raffle: %w", err)
	}
	if raffle == nil {
		return nil, errors.New("raffle not found")
	}

	// 7. Create payment intent with Stripe
	amountInCents := int64(reservation.TotalAmount * 100) // Convert to cents
	metadata := map[string]string{
		"reservation_id": reservation.ID.String(),
		"raffle_id":      reservation.RaffleID.String(),
		"user_id":        reservation.UserID.String(),
		"number_count":   fmt.Sprintf("%d", len(reservation.NumberIDs)),
	}

	stripeIntent, err := uc.paymentProvider.CreatePaymentIntent(ctx, payment.CreatePaymentIntentInput{
		Amount:      amountInCents,
		Currency:    "usd",
		Description: fmt.Sprintf("Raffle: %s - %d numbers", raffle.Title, len(reservation.NumberIDs)),
		Metadata:    metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating payment intent: %w", err)
	}

	// 8. Create payment record in database
	paymentEntity, err := entities.NewPayment(
		reservation.ID,
		reservation.UserID,
		reservation.RaffleID,
		stripeIntent.ID,
		stripeIntent.ClientSecret,
		reservation.TotalAmount,
		stripeIntent.Currency,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating payment entity: %w", err)
	}

	// Set metadata
	paymentMetadata := entities.PaymentMetadata{
		NumberCount: len(reservation.NumberIDs),
		NumberIDs:   reservation.NumberIDs,
		RaffleTitle: raffle.Title,
	}
	if err := paymentEntity.SetMetadata(paymentMetadata); err != nil {
		return nil, fmt.Errorf("error setting payment metadata: %w", err)
	}

	if err := uc.paymentRepo.Create(ctx, paymentEntity); err != nil {
		return nil, fmt.Errorf("error saving payment: %w", err)
	}

	// 9. Create response
	output := &CreatePaymentIntentOutput{
		PaymentID:    paymentEntity.ID,
		ClientSecret: stripeIntent.ClientSecret,
		Amount:       reservation.TotalAmount,
		Currency:     stripeIntent.Currency,
	}

	// 10. Store idempotency key if provided
	if input.IdempotencyKey != "" {
		idempKey, err := entities.NewIdempotencyKey(input.IdempotencyKey, input.UserID, "/payments/intent", input)
		if err == nil {
			if err := idempKey.MarkAsCompleted(200, output); err == nil {
				_ = uc.idempotencyKeyRepo.Create(ctx, idempKey)
			}
		}
	}

	return output, nil
}

// ProcessPaymentWebhook processes a payment webhook from Stripe
func (uc *PaymentUseCases) ProcessPaymentWebhook(ctx context.Context, eventType string, paymentIntentID string) error {
	// 1. Find payment by Stripe Payment Intent ID
	paymentEntity, err := uc.paymentRepo.FindByStripePaymentIntentID(ctx, paymentIntentID)
	if err != nil {
		return fmt.Errorf("error fetching payment: %w", err)
	}
	if paymentEntity == nil {
		return ErrPaymentNotFound
	}

	// 2. Get reservation
	reservation, err := uc.reservationRepo.FindByID(ctx, paymentEntity.ReservationID)
	if err != nil {
		return fmt.Errorf("error fetching reservation: %w", err)
	}
	if reservation == nil {
		return ErrReservationNotFound
	}

	// 3. Handle event based on type
	switch eventType {
	case "payment_intent.succeeded":
		// Mark payment as succeeded
		if err := paymentEntity.MarkAsSucceeded("card"); err != nil {
			return fmt.Errorf("error marking payment as succeeded: %w", err)
		}

		if err := uc.paymentRepo.Update(ctx, paymentEntity); err != nil {
			return fmt.Errorf("error updating payment: %w", err)
		}

		// Confirm reservation
		if err := uc.reservationUseCases.ConfirmReservation(ctx, reservation.ID); err != nil {
			return fmt.Errorf("error confirming reservation: %w", err)
		}

		// TODO: Mark numbers as sold in raffle (future implementation)

	case "payment_intent.payment_failed":
		// Mark payment as failed
		if err := paymentEntity.MarkAsFailed("Payment failed"); err != nil {
			return fmt.Errorf("error marking payment as failed: %w", err)
		}

		if err := uc.paymentRepo.Update(ctx, paymentEntity); err != nil {
			return fmt.Errorf("error updating payment: %w", err)
		}

		// Reservation remains pending, user can retry

	case "payment_intent.canceled":
		// Cancel payment
		if err := paymentEntity.Cancel(); err != nil {
			return fmt.Errorf("error canceling payment: %w", err)
		}

		if err := uc.paymentRepo.Update(ctx, paymentEntity); err != nil {
			return fmt.Errorf("error updating payment: %w", err)
		}

		// Cancel reservation
		if err := uc.reservationUseCases.CancelReservation(ctx, reservation.ID); err != nil {
			return fmt.Errorf("error canceling reservation: %w", err)
		}

	default:
		// Ignore other event types
		return nil
	}

	return nil
}

// GetPayment retrieves a payment by ID
func (uc *PaymentUseCases) GetPayment(ctx context.Context, paymentID uuid.UUID) (*entities.Payment, error) {
	payment, err := uc.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("error fetching payment: %w", err)
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

// GetUserPayments retrieves all payments for a user
func (uc *PaymentUseCases) GetUserPayments(ctx context.Context, userID uuid.UUID) ([]*entities.Payment, error) {
	return uc.paymentRepo.FindByUserID(ctx, userID)
}
