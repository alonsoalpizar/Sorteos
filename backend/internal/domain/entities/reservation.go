package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ReservationStatus represents the state of a reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusExpired   ReservationStatus = "expired"
	ReservationStatusCancelled ReservationStatus = "cancelled"
)

// ReservationPhase represents the current phase of the reservation
type ReservationPhase string

const (
	ReservationPhaseSelection ReservationPhase = "selection" // User selecting numbers (10 min)
	ReservationPhaseCheckout  ReservationPhase = "checkout"  // User in payment flow (5 min)
	ReservationPhaseCompleted ReservationPhase = "completed" // Payment successful
	ReservationPhaseExpired   ReservationPhase = "expired"   // Reservation expired
)

// Timeout durations for each phase
const (
	MaxNumbersPerReservation         = 10
	ReservationSelectionTimeout      = 10 * time.Minute // Phase 1: Selection
	ReservationCheckoutTimeout       = 5 * time.Minute  // Phase 2: Checkout/Payment
	ReservationExpirationDuration    = ReservationSelectionTimeout // Legacy compatibility
)

var (
	ErrReservationExpired      = errors.New("reservation has expired")
	ErrReservationAlreadyPaid  = errors.New("reservation already confirmed")
	ErrReservationCancelled    = errors.New("reservation has been cancelled")
	ErrInvalidReservationState = errors.New("invalid reservation state")
	ErrNoNumbersSelected       = errors.New("no numbers selected for reservation")
	ErrInvalidAmount           = errors.New("invalid reservation amount")
	ErrMaxNumbersExceeded      = errors.New("maximum 10 numbers per reservation")
	ErrCannotAddInCheckout     = errors.New("cannot add numbers during checkout phase")
	ErrNotInSelectionPhase     = errors.New("reservation not in selection phase")
)

// Reservation represents a temporary hold on raffle numbers
type Reservation struct {
	ID          uuid.UUID         `json:"id"`
	RaffleID    uuid.UUID         `json:"raffle_id"`
	UserID      uuid.UUID         `json:"user_id"`
	NumberIDs   pq.StringArray    `json:"number_ids" gorm:"type:text[]"` // PostgreSQL text array
	Status      ReservationStatus `json:"status"`
	SessionID   string            `json:"session_id"`   // For idempotency tracking
	TotalAmount float64           `json:"total_amount"` // Total cost for reserved numbers

	// Double timeout system
	Phase               ReservationPhase `json:"phase" gorm:"type:reservation_phase"`
	SelectionStartedAt  time.Time        `json:"selection_started_at"`
	CheckoutStartedAt   *time.Time       `json:"checkout_started_at,omitempty"`

	ExpiresAt   time.Time         `json:"expires_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// NewReservation creates a new pending reservation in selection phase
func NewReservation(raffleID, userID uuid.UUID, numberIDs []string, sessionID string, totalAmount float64) (*Reservation, error) {
	if len(numberIDs) == 0 {
		return nil, ErrNoNumbersSelected
	}

	if len(numberIDs) > MaxNumbersPerReservation {
		return nil, ErrMaxNumbersExceeded
	}

	if totalAmount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()
	return &Reservation{
		ID:                 uuid.New(),
		RaffleID:           raffleID,
		UserID:             userID,
		NumberIDs:          pq.StringArray(numberIDs),
		Status:             ReservationStatusPending,
		SessionID:          sessionID,
		TotalAmount:        totalAmount,
		Phase:              ReservationPhaseSelection,
		SelectionStartedAt: now,
		ExpiresAt:          now.Add(ReservationSelectionTimeout), // 10 minutes for selection
		CreatedAt:          now,
		UpdatedAt:          now,
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

// AddNumber adds a number to an existing reservation (only in selection phase)
func (r *Reservation) AddNumber(numberID string) error {
	if r.Phase != ReservationPhaseSelection {
		return ErrCannotAddInCheckout
	}

	if len(r.NumberIDs) >= MaxNumbersPerReservation {
		return ErrMaxNumbersExceeded
	}

	// Check for duplicates
	for _, existing := range r.NumberIDs {
		if existing == numberID {
			return errors.New("number already in reservation")
		}
	}

	r.NumberIDs = append(r.NumberIDs, numberID)
	r.UpdatedAt = time.Now()
	return nil
}

// MoveToCheckout transitions the reservation from selection to checkout phase
// This extends the timeout by an additional 5 minutes
func (r *Reservation) MoveToCheckout() error {
	if r.Phase != ReservationPhaseSelection {
		return ErrNotInSelectionPhase
	}

	if r.IsExpired() {
		return ErrReservationExpired
	}

	now := time.Now()
	r.Phase = ReservationPhaseCheckout
	r.CheckoutStartedAt = &now
	r.ExpiresAt = now.Add(ReservationCheckoutTimeout) // 5 minutes for checkout
	r.UpdatedAt = now
	return nil
}

// Confirm marks the reservation as confirmed after successful payment
func (r *Reservation) Confirm() error {
	if err := r.CanBePaid(); err != nil {
		return err
	}

	now := time.Now()
	r.Phase = ReservationPhaseCompleted
	r.Status = ReservationStatusConfirmed
	r.UpdatedAt = now
	return nil
}

// Expire marks the reservation as expired
func (r *Reservation) Expire() error {
	if r.Status != ReservationStatusPending {
		return ErrInvalidReservationState
	}

	r.Phase = ReservationPhaseExpired
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
