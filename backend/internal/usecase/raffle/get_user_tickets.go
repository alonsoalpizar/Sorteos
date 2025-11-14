package raffle

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// TicketGroup representa un grupo de tickets por sorteo
type TicketGroup struct {
	Raffle       *domain.Raffle        `json:"raffle"`
	Numbers      []*domain.RaffleNumber `json:"numbers"`
	TotalNumbers int                    `json:"total_numbers"`
	TotalSpent   string                 `json:"total_spent"` // Decimal como string
}

// GetUserTicketsInput datos de entrada
type GetUserTicketsInput struct {
	UserID   int64
	Page     int
	PageSize int
}

// GetUserTicketsOutput resultado del listado
type GetUserTicketsOutput struct {
	Tickets    []*TicketGroup `json:"tickets"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// GetUserTicketsUseCase caso de uso para obtener tickets del usuario
type GetUserTicketsUseCase struct {
	raffleNumberRepo db.RaffleNumberRepository
	raffleRepo       db.RaffleRepository
}

// NewGetUserTicketsUseCase crea una nueva instancia
func NewGetUserTicketsUseCase(
	raffleNumberRepo db.RaffleNumberRepository,
	raffleRepo db.RaffleRepository,
) *GetUserTicketsUseCase {
	return &GetUserTicketsUseCase{
		raffleNumberRepo: raffleNumberRepo,
		raffleRepo:       raffleRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetUserTicketsUseCase) Execute(ctx context.Context, input *GetUserTicketsInput) (*GetUserTicketsOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Calcular offset
	offset := (input.Page - 1) * input.PageSize

	// Obtener números comprados por el usuario
	numbers, total, err := uc.raffleNumberRepo.FindByUserID(input.UserID, offset, input.PageSize)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Agrupar por raffle_id
	raffleMap := make(map[int64][]*domain.RaffleNumber)
	raffleIDs := []int64{}

	for _, number := range numbers {
		if _, exists := raffleMap[number.RaffleID]; !exists {
			raffleIDs = append(raffleIDs, number.RaffleID)
			raffleMap[number.RaffleID] = []*domain.RaffleNumber{}
		}
		raffleMap[number.RaffleID] = append(raffleMap[number.RaffleID], number)
	}

	// Obtener información de los sorteos
	ticketGroups := []*TicketGroup{}
	for _, raffleID := range raffleIDs {
		raffle, err := uc.raffleRepo.FindByID(raffleID)
		if err != nil {
			// Si no se encuentra el sorteo, continuar con el siguiente
			continue
		}

		raffleNumbers := raffleMap[raffleID]

		// Calcular total gastado usando decimal.Decimal
		totalSpent := decimal.Zero
		for _, num := range raffleNumbers {
			if num.Price != nil {
				totalSpent = totalSpent.Add(*num.Price)
			}
		}

		ticketGroups = append(ticketGroups, &TicketGroup{
			Raffle:       raffle,
			Numbers:      raffleNumbers,
			TotalNumbers: len(raffleNumbers),
			TotalSpent:   totalSpent.String(),
		})
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	return &GetUserTicketsOutput{
		Tickets:    ticketGroups,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}
