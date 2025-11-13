package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain/entities"
	"github.com/sorteos-platform/backend/internal/domain/repositories"
)

// PostgresReservationRepository implements ReservationRepository using PostgreSQL
type PostgresReservationRepository struct {
	db *gorm.DB
}

// NewPostgresReservationRepository creates a new PostgreSQL reservation repository
func NewPostgresReservationRepository(db *gorm.DB) repositories.ReservationRepository {
	return &PostgresReservationRepository{db: db}
}

// Create stores a new reservation
func (r *PostgresReservationRepository) Create(ctx context.Context, reservation *entities.Reservation) error {
	return r.db.WithContext(ctx).Create(reservation).Error
}

// FindByID retrieves a reservation by ID
func (r *PostgresReservationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Reservation, error) {
	var reservation entities.Reservation
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&reservation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reservation, nil
}

// FindBySessionID retrieves a reservation by session ID
func (r *PostgresReservationRepository) FindBySessionID(ctx context.Context, sessionID string) (*entities.Reservation, error) {
	var reservation entities.Reservation
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&reservation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reservation, nil
}

// FindByUserID retrieves all reservations for a user
func (r *PostgresReservationRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Reservation, error) {
	var reservations []*entities.Reservation
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

// FindByRaffleID retrieves all reservations for a raffle
func (r *PostgresReservationRepository) FindByRaffleID(ctx context.Context, raffleID uuid.UUID) ([]*entities.Reservation, error) {
	var reservations []*entities.Reservation
	err := r.db.WithContext(ctx).
		Where("raffle_id = ?", raffleID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

// FindExpiredPending finds all pending reservations that have expired
func (r *PostgresReservationRepository) FindExpiredPending(ctx context.Context, before time.Time) ([]*entities.Reservation, error) {
	var reservations []*entities.Reservation
	err := r.db.WithContext(ctx).
		Model(&entities.Reservation{}).
		Where("status = ? AND expires_at < ?", entities.ReservationStatusPending, before).
		Find(&reservations).Error
	return reservations, err
}

// Update updates an existing reservation
func (r *PostgresReservationRepository) Update(ctx context.Context, reservation *entities.Reservation) error {
	reservation.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(reservation).Error
}

// Delete removes a reservation
func (r *PostgresReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Reservation{}, "id = ?", id).Error
}

// CountActiveReservationsForNumbers counts active reservations for specific numbers
func (r *PostgresReservationRepository) CountActiveReservationsForNumbers(ctx context.Context, raffleID uuid.UUID, numberIDs []string) (int, error) {
	var count int64

	// Query for reservations that:
	// 1. Belong to the raffle
	// 2. Are in pending or confirmed status
	// 3. Have any overlap with the requested numbers
	err := r.db.WithContext(ctx).
		Model(&entities.Reservation{}).
		Where("raffle_id = ?", raffleID).
		Where("status IN ?", []entities.ReservationStatus{
			entities.ReservationStatusPending,
			entities.ReservationStatusConfirmed,
		}).
		Where("number_ids && ?", pq.Array(numberIDs)). // PostgreSQL array overlap operator
		Count(&count).Error

	return int(count), err
}

// FindExpired finds all pending reservations that have expired (alias for FindExpiredPending)
func (r *PostgresReservationRepository) FindExpired(ctx context.Context) ([]*entities.Reservation, error) {
	var reservations []*entities.Reservation
	err := r.db.WithContext(ctx).
		Where("status = ?", entities.ReservationStatusPending).
		Where("expires_at < ?", time.Now()).
		Order("expires_at ASC").
		Limit(100). // Process max 100 per execution
		Find(&reservations).Error

	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// FindActiveByUserAndRaffle finds an active reservation for a user in a specific raffle
func (r *PostgresReservationRepository) FindActiveByUserAndRaffle(ctx context.Context, userID uuid.UUID, raffleID uuid.UUID) (*entities.Reservation, error) {
	var reservation entities.Reservation
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("raffle_id = ?", raffleID).
		Where("status = ?", entities.ReservationStatusPending).
		Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		First(&reservation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &reservation, nil
}

// WithTransaction executes a function within a database transaction
func (r *PostgresReservationRepository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		txCtx := context.WithValue(ctx, "gorm_tx", tx)
		return fn(txCtx)
	})
}
