package category

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// ListCategoriesUseCase caso de uso para listar categorías
type ListCategoriesUseCase struct {
	categoryRepo domain.CategoryRepository
	logger       *logger.Logger
}

// NewListCategoriesUseCase crea una nueva instancia
func NewListCategoriesUseCase(
	categoryRepo domain.CategoryRepository,
	logger *logger.Logger,
) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListCategoriesUseCase) Execute(ctx context.Context) ([]*domain.Category, error) {
	categories, err := uc.categoryRepo.FindAll()
	if err != nil {
		uc.logger.Error("Error listing categories", logger.Error(err))
		return nil, err
	}

	uc.logger.Info("Categories listed successfully", logger.Int("count", len(categories)))
	return categories, nil
}
