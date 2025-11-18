package category

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// DeleteCategoryInput datos de entrada
type DeleteCategoryInput struct {
	CategoryID int64 `json:"category_id"`
}

// DeleteCategoryOutput resultado
type DeleteCategoryOutput struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	DeletedAt  string `json:"deleted_at"`
	Message    string `json:"message"`
}

// DeleteCategoryUseCase caso de uso para eliminar categoría (soft delete)
type DeleteCategoryUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewDeleteCategoryUseCase crea una nueva instancia
func NewDeleteCategoryUseCase(db *gorm.DB, log *logger.Logger) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, input *DeleteCategoryInput, adminID int64) (*DeleteCategoryOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Buscar categoría
	var category struct {
		ID   int64
		Name string
	}

	result := uc.db.WithContext(ctx).
		Table("categories").
		Select("id, name").
		Where("id = ? AND deleted_at IS NULL", input.CategoryID).
		First(&category)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("CATEGORY_NOT_FOUND", "category not found", 404, nil)
		}
		uc.log.Error("Error finding category", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Verificar si hay rifas usando esta categoría
	var raffleCount int64
	uc.db.WithContext(ctx).
		Table("raffles").
		Where("category_id = ? AND deleted_at IS NULL", input.CategoryID).
		Count(&raffleCount)

	if raffleCount > 0 {
		return nil, errors.New("CATEGORY_IN_USE",
			"cannot delete category: it is being used by active raffles", 409, nil)
	}

	// Soft delete de la categoría
	now := time.Now()
	result = uc.db.WithContext(ctx).
		Table("categories").
		Where("id = ?", input.CategoryID).
		Update("deleted_at", now)

	if result.Error != nil {
		uc.log.Error("Error deleting category", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log auditoría
	uc.log.Error("Admin deleted category",
		logger.Int64("admin_id", adminID),
		logger.Int64("category_id", input.CategoryID),
		logger.String("category_name", category.Name),
		logger.String("action", "admin_delete_category"),
		logger.String("severity", "warning"))

	return &DeleteCategoryOutput{
		CategoryID: input.CategoryID,
		Name:       category.Name,
		DeletedAt:  now.Format(time.RFC3339),
		Message:    "Category deleted successfully",
	}, nil
}

// validateInput valida los datos de entrada
func (uc *DeleteCategoryUseCase) validateInput(input *DeleteCategoryInput) error {
	if input.CategoryID <= 0 {
		return errors.New("VALIDATION_FAILED", "category_id is required", 400, nil)
	}

	return nil
}
