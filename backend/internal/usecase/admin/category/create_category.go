package category

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// CreateCategoryInput datos de entrada
type CreateCategoryInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	IsActive    bool    `json:"is_active"`
}

// CreateCategoryOutput resultado
type CreateCategoryOutput struct {
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	IconURL     string  `json:"icon_url,omitempty"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	Message     string  `json:"message"`
}

// CreateCategoryUseCase caso de uso para crear categoría
type CreateCategoryUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCreateCategoryUseCase crea una nueva instancia
func NewCreateCategoryUseCase(db *gorm.DB, log *logger.Logger) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *CreateCategoryUseCase) Execute(ctx context.Context, input *CreateCategoryInput, adminID int64) (*CreateCategoryOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Verificar que no exista una categoría con el mismo nombre
	var count int64
	result := uc.db.WithContext(ctx).
		Table("categories").
		Where("name = ? AND deleted_at IS NULL", input.Name).
		Count(&count)

	if result.Error != nil {
		uc.log.Error("Error checking category name", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	if count > 0 {
		return nil, errors.New("CATEGORY_EXISTS", "category with this name already exists", 409, nil)
	}

	// Crear categoría
	now := time.Now()
	category := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
		"icon_url":    input.IconURL,
		"is_active":   input.IsActive,
		"created_at":  now,
		"updated_at":  now,
	}

	result = uc.db.WithContext(ctx).
		Table("categories").
		Create(category)

	if result.Error != nil {
		uc.log.Error("Error creating category", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Obtener ID de la categoría creada
	var categoryID int64
	uc.db.WithContext(ctx).
		Table("categories").
		Select("id").
		Where("name = ?", input.Name).
		Scan(&categoryID)

	// Log auditoría
	uc.log.Info("Admin created category",
		logger.Int64("admin_id", adminID),
		logger.Int64("category_id", categoryID),
		logger.String("name", input.Name),
		logger.Bool("is_active", input.IsActive),
		logger.String("action", "admin_create_category"))

	// Construir output
	description := ""
	if input.Description != nil {
		description = *input.Description
	}
	iconURL := ""
	if input.IconURL != nil {
		iconURL = *input.IconURL
	}

	return &CreateCategoryOutput{
		CategoryID:  categoryID,
		Name:        input.Name,
		Description: description,
		IconURL:     iconURL,
		IsActive:    input.IsActive,
		CreatedAt:   now.Format(time.RFC3339),
		Message:     "Category created successfully",
	}, nil
}

// validateInput valida los datos de entrada
func (uc *CreateCategoryUseCase) validateInput(input *CreateCategoryInput) error {
	if input.Name == "" {
		return errors.New("VALIDATION_FAILED", "name is required", 400, nil)
	}

	if len(input.Name) > 100 {
		return errors.New("VALIDATION_FAILED", "name must be 100 characters or less", 400, nil)
	}

	if input.Description != nil && len(*input.Description) > 500 {
		return errors.New("VALIDATION_FAILED", "description must be 500 characters or less", 400, nil)
	}

	return nil
}
