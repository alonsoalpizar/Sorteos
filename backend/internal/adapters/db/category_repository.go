package db

import (
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresCategoryRepository implementación de CategoryRepository con PostgreSQL
type PostgresCategoryRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCategoryRepository crea una nueva instancia
func NewCategoryRepository(db *gorm.DB, log *logger.Logger) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{
		db:  db,
		log: log,
	}
}

// FindAll obtiene todas las categorías activas ordenadas
func (r *PostgresCategoryRepository) FindAll() ([]*domain.Category, error) {
	var categories []*domain.Category

	if err := r.db.
		Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Find(&categories).Error; err != nil {
		r.log.Error("Error finding categories", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return categories, nil
}

// FindByID obtiene una categoría por ID
func (r *PostgresCategoryRepository) FindByID(id int64) (*domain.Category, error) {
	var category domain.Category

	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCategoryNotFound
		}
		r.log.Error("Error finding category by ID", logger.Int64("id", id), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &category, nil
}

// FindBySlug obtiene una categoría por slug
func (r *PostgresCategoryRepository) FindBySlug(slug string) (*domain.Category, error) {
	var category domain.Category

	if err := r.db.Where("slug = ?", slug).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCategoryNotFound
		}
		r.log.Error("Error finding category by slug", logger.String("slug", slug), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &category, nil
}

// Create crea una nueva categoría
func (r *PostgresCategoryRepository) Create(category *domain.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		r.log.Error("Error creating category", logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// Update actualiza una categoría
func (r *PostgresCategoryRepository) Update(category *domain.Category) error {
	if err := r.db.Save(category).Error; err != nil {
		r.log.Error("Error updating category", logger.Int64("id", category.ID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// Delete elimina una categoría (soft delete, marca como inactiva)
func (r *PostgresCategoryRepository) Delete(id int64) error {
	if err := r.db.Model(&domain.Category{}).
		Where("id = ?", id).
		Update("is_active", false).Error; err != nil {
		r.log.Error("Error deleting category", logger.Int64("id", id), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}
