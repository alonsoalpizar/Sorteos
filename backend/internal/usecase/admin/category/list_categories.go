package category

import (
	"context"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ListCategoriesInput datos de entrada
type ListCategoriesInput struct {
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Search   *string `json:"search,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	OrderBy  *string `json:"order_by,omitempty"` // created_at, name, raffle_count
}

// CategoryListItem item de categoría en la lista
type CategoryListItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	IsActive    bool   `json:"is_active"`
	RaffleCount int    `json:"raffle_count"`
	CreatedAt   string `json:"created_at"`
}

// ListCategoriesOutput resultado
type ListCategoriesOutput struct {
	Categories []*CategoryListItem `json:"categories"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalCount int64               `json:"total_count"`
	TotalPages int                 `json:"total_pages"`
}

// ListCategoriesUseCase caso de uso para listar categorías
type ListCategoriesUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListCategoriesUseCase crea una nueva instancia
func NewListCategoriesUseCase(db *gorm.DB, log *logger.Logger) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListCategoriesUseCase) Execute(ctx context.Context, input *ListCategoriesInput, adminID int64) (*ListCategoriesOutput, error) {
	// Validar y sanitizar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Construir query base
	query := uc.db.WithContext(ctx).
		Table("categories").
		Where("deleted_at IS NULL")

	// Aplicar filtros
	if input.Search != nil && *input.Search != "" {
		searchPattern := "%" + *input.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if input.IsActive != nil {
		query = query.Where("is_active = ?", *input.IsActive)
	}

	// Contar total
	var totalCount int64
	countQuery := query
	result := countQuery.Count(&totalCount)
	if result.Error != nil {
		uc.log.Error("Error counting categories", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Aplicar ordenamiento
	orderBy := "created_at DESC"
	if input.OrderBy != nil && *input.OrderBy != "" {
		switch *input.OrderBy {
		case "name":
			orderBy = "name ASC"
		case "name_desc":
			orderBy = "name DESC"
		case "created_at":
			orderBy = "created_at ASC"
		case "created_at_desc":
			orderBy = "created_at DESC"
		}
	}
	query = query.Order(orderBy)

	// Aplicar paginación
	offset := (input.Page - 1) * input.PageSize
	query = query.Limit(input.PageSize).Offset(offset)

	// Ejecutar query
	var categories []struct {
		ID          int64
		Name        string
		Description *string
		IconURL     *string
		IsActive    bool
		CreatedAt   string
	}

	result = query.
		Select("id, name, description, icon_url, is_active, created_at").
		Find(&categories)

	if result.Error != nil {
		uc.log.Error("Error fetching categories", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Obtener count de raffles por categoría
	raffleCounts := make(map[int64]int)
	var raffleCountResults []struct {
		CategoryID int64
		Count      int
	}

	uc.db.WithContext(ctx).
		Table("raffles").
		Select("category_id, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("category_id").
		Find(&raffleCountResults)

	for _, rc := range raffleCountResults {
		raffleCounts[rc.CategoryID] = rc.Count
	}

	// Construir output
	categoryItems := make([]*CategoryListItem, 0, len(categories))
	for _, cat := range categories {
		description := ""
		if cat.Description != nil {
			description = *cat.Description
		}
		iconURL := ""
		if cat.IconURL != nil {
			iconURL = *cat.IconURL
		}

		categoryItems = append(categoryItems, &CategoryListItem{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: description,
			IconURL:     iconURL,
			IsActive:    cat.IsActive,
			RaffleCount: raffleCounts[cat.ID],
			CreatedAt:   cat.CreatedAt,
		})
	}

	totalPages := int(totalCount) / input.PageSize
	if int(totalCount)%input.PageSize > 0 {
		totalPages++
	}

	return &ListCategoriesOutput{
		Categories: categoryItems,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

// validateInput valida los datos de entrada
func (uc *ListCategoriesUseCase) validateInput(input *ListCategoriesInput) error {
	if input.Page <= 0 {
		input.Page = 1
	}

	if input.PageSize <= 0 || input.PageSize > 100 {
		input.PageSize = 50
	}

	return nil
}
