package category

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UpdateCategoryInput datos de entrada
type UpdateCategoryInput struct {
	CategoryID  int64   `json:"category_id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// UpdateCategoryOutput resultado
type UpdateCategoryOutput struct {
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	IsActive    bool   `json:"is_active"`
	UpdatedAt   string `json:"updated_at"`
	Message     string `json:"message"`
}

// UpdateCategoryUseCase caso de uso para actualizar categoría
type UpdateCategoryUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUpdateCategoryUseCase crea una nueva instancia
func NewUpdateCategoryUseCase(db *gorm.DB, log *logger.Logger) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateCategoryUseCase) Execute(ctx context.Context, input *UpdateCategoryInput, adminID int64) (*UpdateCategoryOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Buscar categoría
	var category struct {
		ID          int64
		Name        string
		Description *string
		IconURL     *string
		IsActive    bool
	}

	result := uc.db.WithContext(ctx).
		Table("categories").
		Select("id, name, description, icon_url, is_active").
		Where("id = ? AND deleted_at IS NULL", input.CategoryID).
		First(&category)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("CATEGORY_NOT_FOUND", "category not found", 404, nil)
		}
		uc.log.Error("Error finding category", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Si se actualiza el nombre, verificar que no exista otra categoría con ese nombre
	if input.Name != nil && *input.Name != category.Name {
		var count int64
		uc.db.WithContext(ctx).
			Table("categories").
			Where("name = ? AND id != ? AND deleted_at IS NULL", *input.Name, input.CategoryID).
			Count(&count)

		if count > 0 {
			return nil, errors.New("CATEGORY_NAME_EXISTS", "another category with this name already exists", 409, nil)
		}
	}

	// Construir updates
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.Name != nil {
		updates["name"] = *input.Name
		category.Name = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
		category.Description = input.Description
	}
	if input.IconURL != nil {
		updates["icon_url"] = *input.IconURL
		category.IconURL = input.IconURL
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
		category.IsActive = *input.IsActive
	}

	// Actualizar categoría
	result = uc.db.WithContext(ctx).
		Table("categories").
		Where("id = ?", input.CategoryID).
		Updates(updates)

	if result.Error != nil {
		uc.log.Error("Error updating category", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log auditoría
	uc.log.Info("Admin updated category",
		logger.Int64("admin_id", adminID),
		logger.Int64("category_id", input.CategoryID),
		logger.String("action", "admin_update_category"))

	// Construir output
	description := ""
	if category.Description != nil {
		description = *category.Description
	}
	iconURL := ""
	if category.IconURL != nil {
		iconURL = *category.IconURL
	}

	return &UpdateCategoryOutput{
		CategoryID:  input.CategoryID,
		Name:        category.Name,
		Description: description,
		IconURL:     iconURL,
		IsActive:    category.IsActive,
		UpdatedAt:   time.Now().Format(time.RFC3339),
		Message:     "Category updated successfully",
	}, nil
}

// validateInput valida los datos de entrada
func (uc *UpdateCategoryUseCase) validateInput(input *UpdateCategoryInput) error {
	if input.CategoryID <= 0 {
		return errors.New("VALIDATION_FAILED", "category_id is required", 400, nil)
	}

	// Al menos un campo debe ser actualizado
	if input.Name == nil && input.Description == nil && input.IconURL == nil && input.IsActive == nil {
		return errors.New("VALIDATION_FAILED", "at least one field must be provided for update", 400, nil)
	}

	if input.Name != nil && *input.Name == "" {
		return errors.New("VALIDATION_FAILED", "name cannot be empty", 400, nil)
	}

	if input.Name != nil && len(*input.Name) > 100 {
		return errors.New("VALIDATION_FAILED", "name must be 100 characters or less", 400, nil)
	}

	if input.Description != nil && len(*input.Description) > 500 {
		return errors.New("VALIDATION_FAILED", "description must be 500 characters or less", 400, nil)
	}

	return nil
}
