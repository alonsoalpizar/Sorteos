package raffle

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// ListRafflesInput datos de entrada para listar sorteos
type ListRafflesInput struct {
	Page       int
	PageSize   int
	Status     *domain.RaffleStatus
	UserID     *int64
	CategoryID *int64
	Search     string
	OrderBy    string
	OnlyAvailable bool
}

// ListRafflesOutput resultado del listado
type ListRafflesOutput struct {
	Raffles    []*domain.Raffle
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListRafflesUseCase caso de uso para listar sorteos
type ListRafflesUseCase struct {
	raffleRepo db.RaffleRepository
}

// NewListRafflesUseCase crea una nueva instancia
func NewListRafflesUseCase(raffleRepo db.RaffleRepository) *ListRafflesUseCase {
	return &ListRafflesUseCase{
		raffleRepo: raffleRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListRafflesUseCase) Execute(ctx context.Context, input *ListRafflesInput) (*ListRafflesOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Calcular offset
	offset := (input.Page - 1) * input.PageSize

	// Construir filtros
	filters := make(map[string]interface{})

	if input.Status != nil {
		filters["status"] = *input.Status
	}

	if input.UserID != nil {
		filters["user_id"] = *input.UserID
	}

	if input.CategoryID != nil {
		filters["category_id"] = *input.CategoryID
	}

	if input.Search != "" {
		filters["search"] = input.Search
	}

	if input.OrderBy != "" {
		filters["order_by"] = input.OrderBy
	} else {
		filters["order_by"] = "created_at DESC"
	}

	if input.OnlyAvailable {
		filters["only_available"] = true
	}

	// Obtener sorteos
	raffles, total, err := uc.raffleRepo.List(offset, input.PageSize, filters)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	return &ListRafflesOutput{
		Raffles:    raffles,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}
