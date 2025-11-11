package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ReservationStatus represents the state of a reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusExpired   ReservationStatus = "expired"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

// ReservationExpirationDuration is the time window for completing a reservation (5 minutes)
const ReservationExpirationDuration = 5 * time.Minute

var (
	ErrReservationExpired      = errors.New("reservation has expired")
	ErrReservationAlreadyPaid  = errors.New("reservation already confirmed")
	ErrReservationCancelled    = errors.New("reservation has been cancelled")
	ErrInvalidReservationState = errors.New("invalid reservation state")
	ErrNoNumbersSelected       = errors.New("no numbers selected for reservation")
	ErrInvalidAmount           = errors.New("invalid reservation amount")
)

// Reservation represents a temporary hold on raffle numbers
type Reservation struct {
	ID          uuid.UUID         `json:"id"`
	RaffleID    uuid.UUID         `json:"raffle_id"`
	UserID      uuid.UUID         `json:"user_id"`
	NumberIDs   []string          `json:"number_ids" gorm:"type:text[]"` // PostgreSQL text array
	Status      ReservationStatus `json:"status"`
	SessionID   string            `json:"session_id"`   // For idempotency tracking
	TotalAmount float64           `json:"total_amount"` // Total cost for reserved numbers
	ExpiresAt   time.Time         `json:"expires_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// NewReservation creates a new pending reservation
func NewReservation(raffleID, userID uuid.UUID, numberIDs []string, sessionID string, totalAmount float64) (*Reservation, error) {
	if len(numberIDs) == 0 {
		return nil, ErrNoNumbersSelected
	}

	if totalAmount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()
	return &Reservation{
		ID:          uuid.New(),
		RaffleID:    raffleID,
		UserID:      userID,
		NumberIDs:   numberIDs,
		Status:      ReservationStatusPending,
		SessionID:   sessionID,
		TotalAmount: totalAmount,
		ExpiresAt:   now.Add(ReservationExpirationDuration),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// IsExpired checks if the reservation has expired
func (r *Reservation) IsExpired() bool {
	return time.Now().After(r.ExpiresAt) && r.Status == ReservationStatusPending
}

// CanBePaid checks if the reservation can be paid
func (r *Reservation) CanBePaid() error {
	if r.Status == ReservationStatusConfirmed {
		return ErrReservationAlreadyPaid
	}

	if r.Status == ReservationStatusCancelled {
		return ErrReservationCancelled
	}

	if r.IsExpired() {
		return ErrReservationExpired
	}

	return nil
}

// Confirm marks the reservation as confirmed after successful payment
func (r *Reservation) Confirm() error {
	if err := r.CanBePaid(); err != nil {
		return err
	}

	r.Status = ReservationStatusConfirmed
	r.UpdatedAt = time.Now()
	return nil
}

// Expire marks the reservation as expired
func (r *Reservation) Expire() error {
	if r.Status != ReservationStatusPending {
		return ErrInvalidReservationState
	}

	r.Status = ReservationStatusExpired
	r.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the reservation
func (r *Reservation) Cancel() error {
	if r.Status == ReservationStatusConfirmed {
		return ErrReservationAlreadyPaid
	}

	r.Status = ReservationStatusCancelled
	r.UpdatedAt = time.Now()
	return nil
}

// TimeRemaining returns the time remaining before expiration
func (r *Reservation) TimeRemaining() time.Duration {
	if r.Status != ReservationStatusPending {
		return 0
	}

	remaining := time.Until(r.ExpiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetLockKeys returns the Redis lock keys for this reservation's numbers
func (r *Reservation) GetLockKeys() []string {
	keys := make([]string, len(r.NumberIDs))
	for i, numberID := range r.NumberIDs {
		keys[i] = r.RaffleID.String() + ":" + numberID
	}
	return keys
}
