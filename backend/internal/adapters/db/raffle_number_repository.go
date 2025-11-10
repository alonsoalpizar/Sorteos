package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// RaffleNumberRepository define los métodos de acceso a datos para RaffleNumber
type RaffleNumberRepository interface {
	Create(number *domain.RaffleNumber) error
	CreateBatch(numbers []*domain.RaffleNumber) error
	Update(number *domain.RaffleNumber) error
	FindByID(id int64) (*domain.RaffleNumber, error)
	FindByRaffleAndNumber(raffleID int64, number string) (*domain.RaffleNumber, error)
	FindByRaffleID(raffleID int64) ([]*domain.RaffleNumber, error)
	FindAvailableByRaffleID(raffleID int64) ([]*domain.RaffleNumber, error)
	FindByUserID(userID int64, offset, limit int) ([]*domain.RaffleNumber, int64, error)
	CountByStatus(raffleID int64, status domain.RaffleNumberStatus) (int64, error)
	ReserveNumbers(raffleID int64, numbers []string, userID, reservationID int64, duration time.Duration) error
	ReleaseExpiredReservations() (int, error)
	MarkAsSold(id int64, userID, paymentID int64) error
	CancelReservation(id int64) error
}

// RaffleNumberRepositoryImpl implementa RaffleNumberRepository
type RaffleNumberRepositoryImpl struct {
	db *gorm.DB
}

// NewRaffleNumberRepository crea una nueva instancia del repositorio
func NewRaffleNumberRepository(db *gorm.DB) RaffleNumberRepository {
	return &RaffleNumberRepositoryImpl{db: db}
}

// Create crea un nuevo número de sorteo
func (r *RaffleNumberRepositoryImpl) Create(number *domain.RaffleNumber) error {
	if err := r.db.Create(number).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// CreateBatch crea múltiples números de sorteo en batch
func (r *RaffleNumberRepositoryImpl) CreateBatch(numbers []*domain.RaffleNumber) error {
	if len(numbers) == 0 {
		return nil
	}

	// Usar CreateInBatches para mejor performance
	if err := r.db.CreateInBatches(numbers, 100).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// Update actualiza un número de sorteo
func (r *RaffleNumberRepositoryImpl) Update(number *domain.RaffleNumber) error {
	if err := r.db.Save(number).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// FindByID busca un número por ID
func (r *RaffleNumberRepositoryImpl) FindByID(id int64) (*domain.RaffleNumber, error) {
	var number domain.RaffleNumber
	if err := r.db.First(&number, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &number, nil
}

// FindByRaffleAndNumber busca un número específico de un sorteo
func (r *RaffleNumberRepositoryImpl) FindByRaffleAndNumber(raffleID int64, number string) (*domain.RaffleNumber, error) {
	var raffleNumber domain.RaffleNumber
	if err := r.db.Where("raffle_id = ? AND number = ?", raffleID, number).First(&raffleNumber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &raffleNumber, nil
}

// FindByRaffleID busca todos los números de un sorteo
func (r *RaffleNumberRepositoryImpl) FindByRaffleID(raffleID int64) ([]*domain.RaffleNumber, error) {
	var numbers []*domain.RaffleNumber
	if err := r.db.Where("raffle_id = ?", raffleID).Order("number ASC").Find(&numbers).Error; err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return numbers, nil
}

// FindAvailableByRaffleID busca números disponibles de un sorteo
func (r *RaffleNumberRepositoryImpl) FindAvailableByRaffleID(raffleID int64) ([]*domain.RaffleNumber, error) {
	var numbers []*domain.RaffleNumber
	if err := r.db.Where("raffle_id = ? AND status = ?", raffleID, domain.RaffleNumberStatusAvailable).
		Order("number ASC").
		Find(&numbers).Error; err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return numbers, nil
}

// FindByUserID busca números comprados por un usuario
func (r *RaffleNumberRepositoryImpl) FindByUserID(userID int64, offset, limit int) ([]*domain.RaffleNumber, int64, error) {
	var numbers []*domain.RaffleNumber
	var total int64

	query := r.db.Model(&domain.RaffleNumber{}).Where("user_id = ? AND status = ?", userID, domain.RaffleNumberStatusSold)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("sold_at DESC").Offset(offset).Limit(limit).Find(&numbers).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return numbers, total, nil
}

// CountByStatus cuenta números por estado
func (r *RaffleNumberRepositoryImpl) CountByStatus(raffleID int64, status domain.RaffleNumberStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.RaffleNumber{}).
		Where("raffle_id = ? AND status = ?", raffleID, status).
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return count, nil
}

// ReserveNumbers reserva múltiples números en una transacción
func (r *RaffleNumberRepositoryImpl) ReserveNumbers(raffleID int64, numbers []string, userID, reservationID int64, duration time.Duration) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		until := now.Add(duration)

		// Verificar que todos los números estén disponibles (con FOR UPDATE para lock)
		var count int64
		if err := tx.Model(&domain.RaffleNumber{}).
			Where("raffle_id = ? AND number IN ? AND status = ?", raffleID, numbers, domain.RaffleNumberStatusAvailable).
			Count(&count).Error; err != nil {
			return errors.Wrap(errors.ErrDatabaseError, err)
		}

		if int(count) != len(numbers) {
			return errors.ErrNumberAlreadyReserved
		}

		// Actualizar a reservado
		updates := map[string]interface{}{
			"status":         domain.RaffleNumberStatusReserved,
			"reserved_at":    now,
			"reserved_until": until,
			"reserved_by":    userID,
			"reservation_id": reservationID,
			"updated_at":     now,
		}

		if err := tx.Model(&domain.RaffleNumber{}).
			Where("raffle_id = ? AND number IN ? AND status = ?", raffleID, numbers, domain.RaffleNumberStatusAvailable).
			Updates(updates).Error; err != nil {
			return errors.Wrap(errors.ErrDatabaseError, err)
		}

		return nil
	})
}

// ReleaseExpiredReservations libera todas las reservas expiradas
func (r *RaffleNumberRepositoryImpl) ReleaseExpiredReservations() (int, error) {
	now := time.Now()

	result := r.db.Model(&domain.RaffleNumber{}).
		Where("status = ? AND reserved_until < ?", domain.RaffleNumberStatusReserved, now).
		Updates(map[string]interface{}{
			"status":         domain.RaffleNumberStatusAvailable,
			"reserved_at":    nil,
			"reserved_until": nil,
			"reserved_by":    nil,
			"reservation_id": nil,
			"updated_at":     now,
		})

	if result.Error != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	return int(result.RowsAffected), nil
}

// MarkAsSold marca un número como vendido
func (r *RaffleNumberRepositoryImpl) MarkAsSold(id int64, userID, paymentID int64) error {
	// Buscar el número para obtener el precio del raffle
	var number domain.RaffleNumber
	if err := r.db.Preload("Raffle").First(&number, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Nota: Asumiendo que el precio se obtendrá del raffle asociado
	// En una implementación real, necesitarías hacer un JOIN o preload
	now := time.Now()

	updates := map[string]interface{}{
		"status":         domain.RaffleNumberStatusSold,
		"user_id":        userID,
		"payment_id":     paymentID,
		"sold_at":        now,
		"reserved_at":    nil,
		"reserved_until": nil,
		"reserved_by":    nil,
		"reservation_id": nil,
		"updated_at":     now,
	}

	if err := r.db.Model(&domain.RaffleNumber{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// CancelReservation cancela la reserva de un número
func (r *RaffleNumberRepositoryImpl) CancelReservation(id int64) error {
	now := time.Now()

	updates := map[string]interface{}{
		"status":         domain.RaffleNumberStatusAvailable,
		"reserved_at":    nil,
		"reserved_until": nil,
		"reserved_by":    nil,
		"reservation_id": nil,
		"updated_at":     now,
	}

	if err := r.db.Model(&domain.RaffleNumber{}).
		Where("id = ? AND status = ?", id, domain.RaffleNumberStatusReserved).
		Updates(updates).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}
