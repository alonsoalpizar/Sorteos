package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// RaffleImageRepository define los métodos de acceso a datos para RaffleImage
type RaffleImageRepository interface {
	Create(image *domain.RaffleImage) error
	Update(image *domain.RaffleImage) error
	Delete(id int64) error
	SoftDelete(id int64) error
	FindByID(id int64) (*domain.RaffleImage, error)
	FindByRaffleID(raffleID int64) ([]*domain.RaffleImage, error)
	FindPrimaryByRaffleID(raffleID int64) (*domain.RaffleImage, error)
	SetPrimary(id int64) error
	UpdateDisplayOrder(id int64, order int) error
	CountByRaffleID(raffleID int64) (int64, error)
}

// RaffleImageRepositoryImpl implementa RaffleImageRepository
type RaffleImageRepositoryImpl struct {
	db *gorm.DB
}

// NewRaffleImageRepository crea una nueva instancia del repositorio
func NewRaffleImageRepository(db *gorm.DB) RaffleImageRepository {
	return &RaffleImageRepositoryImpl{db: db}
}

// Create crea una nueva imagen
func (r *RaffleImageRepositoryImpl) Create(image *domain.RaffleImage) error {
	if err := r.db.Create(image).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// Update actualiza una imagen existente
func (r *RaffleImageRepositoryImpl) Update(image *domain.RaffleImage) error {
	if err := r.db.Save(image).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// Delete elimina físicamente una imagen
func (r *RaffleImageRepositoryImpl) Delete(id int64) error {
	if err := r.db.Delete(&domain.RaffleImage{}, id).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// SoftDelete elimina lógicamente una imagen
func (r *RaffleImageRepositoryImpl) SoftDelete(id int64) error {
	now := time.Now()
	if err := r.db.Model(&domain.RaffleImage{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// FindByID busca una imagen por ID
func (r *RaffleImageRepositoryImpl) FindByID(id int64) (*domain.RaffleImage, error) {
	var image domain.RaffleImage
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &image, nil
}

// FindByRaffleID busca todas las imágenes de un sorteo
func (r *RaffleImageRepositoryImpl) FindByRaffleID(raffleID int64) ([]*domain.RaffleImage, error) {
	var images []*domain.RaffleImage
	if err := r.db.Where("raffle_id = ? AND deleted_at IS NULL", raffleID).
		Order("display_order ASC, id ASC").
		Find(&images).Error; err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return images, nil
}

// FindPrimaryByRaffleID busca la imagen principal de un sorteo
func (r *RaffleImageRepositoryImpl) FindPrimaryByRaffleID(raffleID int64) (*domain.RaffleImage, error) {
	var image domain.RaffleImage
	if err := r.db.Where("raffle_id = ? AND is_primary = ? AND deleted_at IS NULL", raffleID, true).
		First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &image, nil
}

// SetPrimary establece una imagen como principal (y desmarca las demás)
func (r *RaffleImageRepositoryImpl) SetPrimary(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Obtener la imagen
		var image domain.RaffleImage
		if err := tx.First(&image, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.ErrNotFound
			}
			return errors.Wrap(errors.ErrDatabaseError, err)
		}

		// Desmarcar todas las imágenes primary del mismo raffle
		if err := tx.Model(&domain.RaffleImage{}).
			Where("raffle_id = ? AND id != ?", image.RaffleID, id).
			Update("is_primary", false).Error; err != nil {
			return errors.Wrap(errors.ErrDatabaseError, err)
		}

		// Marcar esta imagen como primary
		if err := tx.Model(&domain.RaffleImage{}).
			Where("id = ?", id).
			Update("is_primary", true).Error; err != nil {
			return errors.Wrap(errors.ErrDatabaseError, err)
		}

		return nil
	})
}

// UpdateDisplayOrder actualiza el orden de visualización
func (r *RaffleImageRepositoryImpl) UpdateDisplayOrder(id int64, order int) error {
	if err := r.db.Model(&domain.RaffleImage{}).
		Where("id = ?", id).
		Update("display_order", order).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// CountByRaffleID cuenta las imágenes de un sorteo
func (r *RaffleImageRepositoryImpl) CountByRaffleID(raffleID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.RaffleImage{}).
		Where("raffle_id = ? AND deleted_at IS NULL", raffleID).
		Count(&count).Error; err != nil {
		return 0, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return count, nil
}
