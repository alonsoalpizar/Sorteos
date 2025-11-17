package domain

import (
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
)

// Category representa una categoría de sorteo
type Category struct {
	ID           int64
	Name         string
	Slug         string
	Icon         string // emoji
	Description  string
	DisplayOrder int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CategoryRepository interfaz para persistencia de categorías
type CategoryRepository interface {
	// FindAll obtiene todas las categorías activas ordenadas
	FindAll() ([]*Category, error)

	// FindByID obtiene una categoría por ID
	FindByID(id int64) (*Category, error)

	// FindBySlug obtiene una categoría por slug
	FindBySlug(slug string) (*Category, error)

	// Create crea una nueva categoría (admin)
	Create(category *Category) error

	// Update actualiza una categoría (admin)
	Update(category *Category) error

	// Delete elimina una categoría (admin, soft delete)
	Delete(id int64) error
}

// Validate valida los datos de una categoría
func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.ErrValidationFailed
	}
	if len(c.Name) < 3 || len(c.Name) > 100 {
		return errors.ErrValidationFailed
	}
	if c.Slug == "" {
		return errors.ErrValidationFailed
	}
	if c.Icon == "" {
		return errors.ErrValidationFailed
	}
	return nil
}
