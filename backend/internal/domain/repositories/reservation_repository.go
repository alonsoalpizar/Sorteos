package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sorteos-platform/backend/internal/domain/entities"
)

// ReservationRepository defines the interface for reservation persistence
type ReservationRepository interface {
	// Create stores a new reservation
	Create(ctx context.Context, reservation *entities.Reservation) error

	// FindByID retrieves a reservation by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Reservation, error)

	// FindBySessionID retrieves a reservation by session ID (for idempotency)
	FindBySessionID(ctx context.Context, sessionID string) (*entities.Reservation, error)

	// FindByUserID retrieves all reservations for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reservation, error)

	// FindByRaffleID retrieves all reservations for a raffle
	FindByRaffleID(ctx context.Context, raffleID uuid.UUID) ([]*entities.Reservation, error)

	// FindExpiredPending finds all pending reservations that have expired
	FindExpiredPending(ctx context.Context, before time.Time) ([]*entities.Reservation, error)

	// Update updates an existing reservation
	Update(ctx context.Context, reservation *entities.Reservation) error

	// Delete removes a reservation
	Delete(ctx context.Context, id uuid.UUID) error

	// CountActiveReservationsForNumbers counts active reservations for specific numbers
	CountActiveReservationsForNumbers(ctx context.Context, raffleID uuid.UUID, numberIDs []string) (int, error)

	// FindExpired finds all pending reservations that have expired
	FindExpired(ctx context.Context) ([]*entities.Reservation, error)

	// FindActiveByUserAndRaffle finds an active reservation for a user in a specific raffle
	FindActiveByUserAndRaffle(ctx context.Context, userID uuid.UUID, raffleID uuid.UUID) (*entities.Reservation, error)

	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
