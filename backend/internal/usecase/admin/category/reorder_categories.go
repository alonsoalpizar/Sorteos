package category

import (
	"context"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ReorderCategoriesInput datos de entrada
type ReorderCategoriesInput struct {
	CategoryIDs []int64 `json:"category_ids" binding:"required"` // Array de IDs en el nuevo orden
}

// ReorderCategoriesOutput resultado
type ReorderCategoriesOutput struct {
	UpdatedCount int    `json:"updated_count"`
	Success      bool   `json:"success"`
	Message      string `json:"message"`
}

// ReorderCategoriesUseCase caso de uso para reordenar categorías
type ReorderCategoriesUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewReorderCategoriesUseCase crea una nueva instancia
func NewReorderCategoriesUseCase(db *gorm.DB, log *logger.Logger) *ReorderCategoriesUseCase {
	return &ReorderCategoriesUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ReorderCategoriesUseCase) Execute(ctx context.Context, input *ReorderCategoriesInput, adminID int64) (*ReorderCategoriesOutput, error) {
	// Validar que se proporcionen IDs
	if len(input.CategoryIDs) == 0 {
		return nil, errors.New("VALIDATION_FAILED", "category_ids array cannot be empty", 400, nil)
	}

	// Verificar que todas las categorías existen
	var count int64
	if err := uc.db.WithContext(ctx).
		Table("categories").
		Where("id IN ?", input.CategoryIDs).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		uc.log.Error("Error counting categories", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if int(count) != len(input.CategoryIDs) {
		return nil, errors.New("VALIDATION_FAILED",
			"some category IDs do not exist or are deleted", 400, nil)
	}

	// Actualizar display_order de cada categoría en una transacción
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, categoryID := range input.CategoryIDs {
			displayOrder := i + 1 // display_order empieza en 1

			if err := tx.Table("categories").
				Where("id = ?", categoryID).
				Update("display_order", displayOrder).Error; err != nil {
				uc.log.Error("Error updating category display_order",
					logger.Int64("category_id", categoryID),
					logger.Int("new_order", displayOrder),
					logger.Error(err))
				return err
			}
		}
		return nil
	})

	if err != nil {
		uc.log.Error("Error reordering categories", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría
	uc.log.Info("Admin reordered categories",
		logger.Int64("admin_id", adminID),
		logger.Int("total_categories", len(input.CategoryIDs)),
		logger.String("action", "admin_reorder_categories"))

	return &ReorderCategoriesOutput{
		UpdatedCount: len(input.CategoryIDs),
		Success:      true,
		Message:      "Categories reordered successfully",
	}, nil
}
