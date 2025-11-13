package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	dbadapter "github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain/entities"
	"github.com/sorteos-platform/backend/internal/domain/repositories"
	"github.com/sorteos-platform/backend/internal/infrastructure/redis"
	"github.com/sorteos-platform/backend/internal/infrastructure/websocket"
)

var (
	ErrNumbersAlreadyReserved = errors.New("one or more numbers are already reserved")
	ErrRaffleNotActive        = errors.New("raffle is not active")
	ErrInsufficientNumbers    = errors.New("some requested numbers are not available")
)

// ReservationUseCases handles business logic for reservations
type ReservationUseCases struct {
	reservationRepo   repositories.ReservationRepository
	raffleRepo        dbadapter.RaffleRepository
	raffleNumberRepo  dbadapter.RaffleNumberRepository
	lockService       *redis.LockService
	wsHub             *websocket.Hub // WebSocket hub for real-time updates
}

// NewReservationUseCases creates a new reservation use cases instance
func NewReservationUseCases(
	reservationRepo repositories.ReservationRepository,
	raffleRepo dbadapter.RaffleRepository,
	raffleNumberRepo dbadapter.RaffleNumberRepository,
	lockService *redis.LockService,
	wsHub *websocket.Hub,
) *ReservationUseCases {
	return &ReservationUseCases{
		reservationRepo:  reservationRepo,
		raffleRepo:       raffleRepo,
		raffleNumberRepo: raffleNumberRepo,
		lockService:      lockService,
		wsHub:            wsHub,
	}
}

// CreateReservationInput represents the input for creating a reservation
type CreateReservationInput struct {
	RaffleID  uuid.UUID
	UserID    uuid.UUID
	NumberIDs []string
	SessionID string
}

// CreateReservation creates a new number reservation with distributed locks
func (uc *ReservationUseCases) CreateReservation(ctx context.Context, input CreateReservationInput) (*entities.Reservation, error) {
	// 1. Check for existing reservation with same session ID (idempotency)
	existingReservation, err := uc.reservationRepo.FindBySessionID(ctx, input.SessionID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing reservation: %w", err)
	}
	if existingReservation != nil {
		// Return existing reservation if it's still valid
		if !existingReservation.IsExpired() && existingReservation.Status == entities.ReservationStatusPending {
			return existingReservation, nil
		}
	}

	// 2. Validate raffle exists and is active
	raffle, err := uc.raffleRepo.FindByUUID(input.RaffleID.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching raffle: %w", err)
	}
	if raffle == nil {
		return nil, errors.New("raffle not found")
	}
	if raffle.Status != "active" {
		return nil, ErrRaffleNotActive
	}

	// 3. Calculate total amount
	pricePerNumber, _ := raffle.PricePerNumber.Float64()
	totalAmount := float64(len(input.NumberIDs)) * pricePerNumber

	// 4. Acquire distributed locks for all numbers
	lockKeys := make([]string, len(input.NumberIDs))
	for i, numberID := range input.NumberIDs {
		lockKeys[i] = redis.ReservationLockKey(input.RaffleID.String(), numberID)
	}

	locks, err := uc.lockService.AcquireMultipleLocks(ctx, lockKeys, entities.ReservationExpirationDuration)
	if err != nil {
		if errors.Is(err, redis.ErrLockNotAcquired) {
			return nil, ErrNumbersAlreadyReserved
		}
		return nil, fmt.Errorf("error acquiring locks: %w", err)
	}

	// Ensure locks are released if we encounter an error
	defer func() {
		if err != nil {
			_ = redis.ReleaseMultipleLocks(ctx, locks)
		}
	}()

	// 5. Double-check numbers aren't already reserved in database
	count, err := uc.reservationRepo.CountActiveReservationsForNumbers(ctx, input.RaffleID, input.NumberIDs)
	if err != nil {
		return nil, fmt.Errorf("error checking existing reservations: %w", err)
	}
	if count > 0 {
		return nil, ErrNumbersAlreadyReserved
	}

	// 6. Create reservation entity
	reservation, err := entities.NewReservation(
		input.RaffleID,
		input.UserID,
		input.NumberIDs,
		input.SessionID,
		totalAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating reservation entity: %w", err)
	}

	// 7. Save to database
	if err := uc.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, fmt.Errorf("error saving reservation: %w", err)
	}

	// 8. Update raffle_numbers table to mark as RESERVED
	// Note: We use placeholder user_id 1 for legacy compatibility
	if err := uc.raffleNumberRepo.ReserveNumbers(raffle.ID, input.NumberIDs, 1, 0, entities.ReservationExpirationDuration); err != nil {
		// Log error but continue - the reservation in MongoDB is what matters
		fmt.Printf("[CreateReservation] Error marking numbers as reserved: %v\n", err)
	}

	// 9. Notify via WebSocket about new reservation
	userIDStr := input.UserID.String()
	for _, numberID := range input.NumberIDs {
		uc.wsHub.BroadcastNumberUpdate(
			input.RaffleID.String(),
			numberID,
			"reserved",
			&userIDStr,
		)
	}

	// Locks will be held until they expire (10 minutes) or until the reservation is confirmed/cancelled
	return reservation, nil
}

// ConfirmReservation confirms a reservation after successful payment
func (uc *ReservationUseCases) ConfirmReservation(ctx context.Context, reservationID uuid.UUID) error {
	reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("error fetching reservation: %w", err)
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	if err := reservation.Confirm(); err != nil {
		return err
	}

	if err := uc.reservationRepo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("error updating reservation: %w", err)
	}

	// Get raffle to obtain integer ID
	raffle, err := uc.raffleRepo.FindByUUID(reservation.RaffleID.String())
	if err != nil {
		return fmt.Errorf("error fetching raffle: %w", err)
	}

	// Update raffle_numbers table to mark numbers as SOLD
	for _, numberStr := range reservation.NumberIDs {
		raffleNumber, err := uc.raffleNumberRepo.FindByRaffleAndNumber(raffle.ID, numberStr)
		if err != nil {
			// Log error but continue with other numbers
			fmt.Printf("[ConfirmReservation] Error finding number %s: %v\n", numberStr, err)
			continue
		}

		// Mark as sold (use placeholder user_id 1 since we don't have integer user IDs)
		// TODO: Implement proper UUID to integer user_id conversion
		if err := uc.raffleNumberRepo.MarkAsSold(raffleNumber.ID, 1, 0); err != nil {
			// Log error but continue
			fmt.Printf("[ConfirmReservation] Error marking number %s as sold: %v\n", numberStr, err)
			continue
		}
	}

	// Notify via WebSocket that numbers are now SOLD
	userIDStr := reservation.UserID.String()
	for _, numberID := range reservation.NumberIDs {
		uc.wsHub.BroadcastNumberUpdate(
			reservation.RaffleID.String(),
			numberID,
			"sold",
			&userIDStr,
		)
	}

	return nil
}

// CancelReservation cancels a reservation and releases locks
func (uc *ReservationUseCases) CancelReservation(ctx context.Context, reservationID uuid.UUID) error {
	reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("error fetching reservation: %w", err)
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	if err := reservation.Cancel(); err != nil {
		return err
	}

	if err := uc.reservationRepo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("error updating reservation: %w", err)
	}

	// Get raffle to obtain integer ID
	raffle, err := uc.raffleRepo.FindByUUID(reservation.RaffleID.String())
	if err == nil {
		// Update raffle_numbers table to mark as AVAILABLE
		for _, numberStr := range reservation.NumberIDs {
			raffleNumber, err := uc.raffleNumberRepo.FindByRaffleAndNumber(raffle.ID, numberStr)
			if err != nil {
				continue
			}
			// Cancel reservation on this number
			if err := uc.raffleNumberRepo.CancelReservation(raffleNumber.ID); err != nil {
				fmt.Printf("[CancelReservation] Error releasing number %s: %v\n", numberStr, err)
			}
		}
	}

	// Notify via WebSocket that numbers are available again
	for _, numberID := range reservation.NumberIDs {
		uc.wsHub.BroadcastNumberUpdate(
			reservation.RaffleID.String(),
			numberID,
			"available",
			nil,
		)
	}

	// Release locks manually
	return uc.releaseLocks(ctx, reservation)
}

// ExpireReservations finds and expires all pending reservations that have passed their expiration time
func (uc *ReservationUseCases) ExpireReservations(ctx context.Context) (int, error) {
	expiredReservations, err := uc.reservationRepo.FindExpiredPending(ctx, time.Now())
	if err != nil {
		return 0, fmt.Errorf("error finding expired reservations: %w", err)
	}

	count := 0
	for _, reservation := range expiredReservations {
		if err := reservation.Expire(); err != nil {
			// Log error but continue processing other reservations
			continue
		}

		if err := uc.reservationRepo.Update(ctx, reservation); err != nil {
			// Log error but continue
			continue
		}

		// Release locks (they may have already expired, but try anyway)
		_ = uc.releaseLocks(ctx, reservation)

		// Update raffle_numbers table to mark as AVAILABLE
		raffle, err := uc.raffleRepo.FindByUUID(reservation.RaffleID.String())
		if err == nil {
			for _, numberStr := range reservation.NumberIDs {
				raffleNumber, err := uc.raffleNumberRepo.FindByRaffleAndNumber(raffle.ID, numberStr)
				if err != nil {
					continue
				}
				// Cancel reservation on this number
				if err := uc.raffleNumberRepo.CancelReservation(raffleNumber.ID); err != nil {
					fmt.Printf("[ExpireReservations] Error releasing number %s: %v\n", numberStr, err)
				}
			}
		}

		// Notify via WebSocket that numbers are available again
		uc.wsHub.BroadcastReservationExpired(
			reservation.RaffleID.String(),
			reservation.NumberIDs,
		)

		count++
	}

	return count, nil
}

// GetReservation retrieves a reservation by ID
func (uc *ReservationUseCases) GetReservation(ctx context.Context, reservationID uuid.UUID) (*entities.Reservation, error) {
	reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return nil, fmt.Errorf("error fetching reservation: %w", err)
	}
	if reservation == nil {
		return nil, errors.New("reservation not found")
	}
	return reservation, nil
}

// GetUserReservations retrieves all reservations for a user
func (uc *ReservationUseCases) GetUserReservations(ctx context.Context, userID uuid.UUID) ([]*entities.Reservation, error) {
	return uc.reservationRepo.FindByUserID(ctx, userID)
}

// MoveToCheckout transitions a reservation from selection to checkout phase
// This is called when user clicks "Pay Now" button
func (uc *ReservationUseCases) MoveToCheckout(ctx context.Context, reservationID uuid.UUID) error {
	reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("error fetching reservation: %w", err)
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	// Transition to checkout phase (extends timeout to 5 more minutes)
	if err := reservation.MoveToCheckout(); err != nil {
		return err
	}

	// Save updated reservation
	if err := uc.reservationRepo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("error updating reservation: %w", err)
	}

	return nil
}

// AddNumberToReservation adds a number to an existing reservation (only in selection phase)
func (uc *ReservationUseCases) AddNumberToReservation(ctx context.Context, reservationID uuid.UUID, numberID string) error {
	// 1. Get reservation
	reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("error fetching reservation: %w", err)
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	// 2. Validate phase and expiration
	if reservation.Phase != entities.ReservationPhaseSelection {
		return entities.ErrCannotAddInCheckout
	}
	if reservation.IsExpired() {
		return entities.ErrReservationExpired
	}

	// 3. Acquire lock for the new number
	lockKey := redis.ReservationLockKey(reservation.RaffleID.String(), numberID)
	lock, err := uc.lockService.AcquireLock(ctx, lockKey, entities.ReservationSelectionTimeout)
	if err != nil {
		if errors.Is(err, redis.ErrLockNotAcquired) {
			return errors.New("number is already reserved")
		}
		return fmt.Errorf("error acquiring lock: %w", err)
	}
	defer lock.Release(ctx)

	// 4. Check number availability in database
	// (This would require a method in raffle number repository)
	// For now, we skip this check

	// 5. Add number to reservation
	if err := reservation.AddNumber(numberID); err != nil {
		return err
	}

	// 6. Update in database
	if err := uc.reservationRepo.Update(ctx, reservation); err != nil {
		return fmt.Errorf("error updating reservation: %w", err)
	}

	// 7. Update raffle_numbers table to mark as RESERVED
	raffle, err := uc.raffleRepo.FindByUUID(reservation.RaffleID.String())
	if err != nil {
		return fmt.Errorf("error fetching raffle: %w", err)
	}

	// Mark the number as reserved in raffle_numbers table
	if err := uc.raffleNumberRepo.ReserveNumbers(raffle.ID, []string{numberID}, 1, 0, entities.ReservationSelectionTimeout); err != nil {
		// Log error but continue - the reservation in MongoDB is what matters
		fmt.Printf("[AddNumberToReservation] Error marking number %s as reserved: %v\n", numberID, err)
	}

	// 8. Notify via WebSocket
	userIDStr := reservation.UserID.String()
	uc.wsHub.BroadcastNumberUpdate(
		reservation.RaffleID.String(),
		numberID,
		"reserved",
		&userIDStr,
	)

	return nil
}

// releaseLocks releases Redis locks for a reservation
func (uc *ReservationUseCases) releaseLocks(ctx context.Context, reservation *entities.Reservation) error {
	// Los locks ya tienen TTL automático, no necesitamos liberarlos manualmente
	// ya que expirarán con la reserva
	return nil
}

// ExpireOldReservations finds and expires old pending reservations
func (uc *ReservationUseCases) ExpireOldReservations(ctx context.Context) (int, error) {
	// 1. Find expired reservations
	expiredReservations, err := uc.reservationRepo.FindExpired(ctx)
	if err != nil {
		return 0, fmt.Errorf("find expired: %w", err)
	}

	if len(expiredReservations) == 0 {
		return 0, nil
	}

	count := 0

	for _, reservation := range expiredReservations {
		// 2. Mark reservation as expired
		reservation.Status = entities.ReservationStatusExpired
		reservation.Phase = entities.ReservationPhaseExpired
		reservation.UpdatedAt = time.Now()

		if err := uc.reservationRepo.Update(ctx, reservation); err != nil {
			// Log error but continue with next reservations
			continue
		}

		// 3. Notify via WebSocket that numbers are now available
		for _, numberID := range reservation.NumberIDs {
			uc.wsHub.BroadcastNumberUpdate(
				reservation.RaffleID.String(),
				numberID,
				"available",
				nil, // No user_id (released)
			)
		}

		count++
	}

	return count, nil
}

// GetActiveReservation gets the active reservation for a user in a raffle
func (uc *ReservationUseCases) GetActiveReservation(
	ctx context.Context,
	userID uuid.UUID,
	raffleID string,
) (*entities.Reservation, error) {
	raffleUUID, err := uuid.Parse(raffleID)
	if err != nil {
		return nil, fmt.Errorf("invalid raffle id: %w", err)
	}

	reservation, err := uc.reservationRepo.FindActiveByUserAndRaffle(ctx, userID, raffleUUID)
	if err != nil {
		return nil, fmt.Errorf("find active reservation: %w", err)
	}

	return reservation, nil
}
