package entities

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the state of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSucceeded  PaymentStatus = "succeeded"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

var (
	ErrPaymentAlreadyProcessed = errors.New("payment already processed")
	ErrPaymentFailed           = errors.New("payment failed")
	ErrPaymentCancelled        = errors.New("payment was cancelled")
	ErrInvalidPaymentState     = errors.New("invalid payment state")
	ErrInvalidPaymentAmount    = errors.New("invalid payment amount")
)

// Payment represents a payment transaction
type Payment struct {
	ID                     uuid.UUID     `json:"id"`
	ReservationID          uuid.UUID     `json:"reservation_id"`
	UserID                 uuid.UUID     `json:"user_id"`
	RaffleID               uuid.UUID     `json:"raffle_id"`
	StripePaymentIntentID  string        `json:"stripe_payment_intent_id"`
	StripeClientSecret     string        `json:"stripe_client_secret"`
	Amount                 float64       `json:"amount"`
	Currency               string        `json:"currency"`
	Status                 PaymentStatus `json:"status"`
	PaymentMethod          string        `json:"payment_method,omitempty"`
	ErrorMessage           string        `json:"error_message,omitempty"`
	Metadata               string        `json:"metadata,omitempty"` // JSONB stored as string
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
	PaidAt                 *time.Time    `json:"paid_at,omitempty"`
}

// PaymentMetadata represents additional payment information
type PaymentMetadata struct {
	NumberCount   int      `json:"number_count"`
	NumberIDs     []string `json:"number_ids"`
	RaffleTitle   string   `json:"raffle_title,omitempty"`
	CustomerEmail string   `json:"customer_email,omitempty"`
}

// NewPayment creates a new payment in pending state
func NewPayment(reservationID, userID, raffleID uuid.UUID, stripeIntentID, clientSecret string, amount float64, currency string) (*Payment, error) {
	if amount <= 0 {
		return nil, ErrInvalidPaymentAmount
	}

	if currency == "" {
		currency = "USD"
	}

	now := time.Now()
	return &Payment{
		ID:                    uuid.New(),
		ReservationID:         reservationID,
		UserID:                userID,
		RaffleID:              raffleID,
		StripePaymentIntentID: stripeIntentID,
		StripeClientSecret:    clientSecret,
		Amount:                amount,
		Currency:              currency,
		Status:                PaymentStatusPending,
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

// MarkAsProcessing updates payment status to processing
func (p *Payment) MarkAsProcessing() error {
	if p.Status != PaymentStatusPending {
		return ErrInvalidPaymentState
	}

	p.Status = PaymentStatusProcessing
	p.UpdatedAt = time.Now()
	return nil
}

// MarkAsSucceeded marks the payment as successful
func (p *Payment) MarkAsSucceeded(paymentMethod string) error {
	if p.Status == PaymentStatusSucceeded {
		return ErrPaymentAlreadyProcessed
	}

	if p.Status == PaymentStatusCancelled {
		return ErrPaymentCancelled
	}

	now := time.Now()
	p.Status = PaymentStatusSucceeded
	p.PaymentMethod = paymentMethod
	p.PaidAt = &now
	p.UpdatedAt = now
	return nil
}

// MarkAsFailed marks the payment as failed
func (p *Payment) MarkAsFailed(errorMessage string) error {
	if p.Status == PaymentStatusSucceeded {
		return ErrPaymentAlreadyProcessed
	}

	p.Status = PaymentStatusFailed
	p.ErrorMessage = errorMessage
	p.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the payment
func (p *Payment) Cancel() error {
	if p.Status == PaymentStatusSucceeded {
		return ErrPaymentAlreadyProcessed
	}

	p.Status = PaymentStatusCancelled
	p.UpdatedAt = time.Now()
	return nil
}

// Refund marks the payment as refunded
func (p *Payment) Refund() error {
	if p.Status != PaymentStatusSucceeded {
		return ErrInvalidPaymentState
	}

	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	return nil
}

// IsCompleted checks if payment is in final state
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusSucceeded ||
		p.Status == PaymentStatusFailed ||
		p.Status == PaymentStatusCancelled
}

// SetMetadata sets the payment metadata from a struct
func (p *Payment) SetMetadata(metadata PaymentMetadata) error {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	p.Metadata = string(jsonData)
	return nil
}

// GetMetadata retrieves the payment metadata as a struct
func (p *Payment) GetMetadata() (*PaymentMetadata, error) {
	if p.Metadata == "" {
		return &PaymentMetadata{}, nil
	}

	var metadata PaymentMetadata
	if err := json.Unmarshal([]byte(p.Metadata), &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}
