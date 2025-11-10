package raffle

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// GetRaffleDetailInput datos de entrada
type GetRaffleDetailInput struct {
	RaffleID       *int64
	RaffleUUID     *string
	IncludeNumbers bool
	IncludeImages  bool
}

// GetRaffleDetailOutput resultado con detalles completos
type GetRaffleDetailOutput struct {
	Raffle         *domain.Raffle
	Numbers        []*domain.RaffleNumber
	Images         []*domain.RaffleImage
	AvailableCount int64
	ReservedCount  int64
	SoldCount      int64
}

// GetRaffleDetailUseCase caso de uso para obtener detalle de sorteo
type GetRaffleDetailUseCase struct {
	raffleRepo       db.RaffleRepository
	raffleNumberRepo db.RaffleNumberRepository
	raffleImageRepo  db.RaffleImageRepository
}

// NewGetRaffleDetailUseCase crea una nueva instancia
func NewGetRaffleDetailUseCase(
	raffleRepo db.RaffleRepository,
	raffleNumberRepo db.RaffleNumberRepository,
	raffleImageRepo db.RaffleImageRepository,
) *GetRaffleDetailUseCase {
	return &GetRaffleDetailUseCase{
		raffleRepo:       raffleRepo,
		raffleNumberRepo: raffleNumberRepo,
		raffleImageRepo:  raffleImageRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetRaffleDetailUseCase) Execute(ctx context.Context, input *GetRaffleDetailInput) (*GetRaffleDetailOutput, error) {
	// Validar que se provea ID o UUID
	if input.RaffleID == nil && input.RaffleUUID == nil {
		return nil, errors.ErrBadRequest
	}

	// Buscar el sorteo
	var raffle *domain.Raffle
	var err error

	if input.RaffleID != nil {
		raffle, err = uc.raffleRepo.FindByID(*input.RaffleID)
	} else {
		raffle, err = uc.raffleRepo.FindByUUID(*input.RaffleUUID)
	}

	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	output := &GetRaffleDetailOutput{
		Raffle: raffle,
	}

	// Obtener números si se solicita
	if input.IncludeNumbers {
		numbers, err := uc.raffleNumberRepo.FindByRaffleID(raffle.ID)
		if err != nil {
			return nil, errors.Wrap(errors.ErrDatabaseError, err)
		}
		output.Numbers = numbers
	}

	// Obtener conteos por estado
	availableCount, err := uc.raffleNumberRepo.CountByStatus(raffle.ID, domain.RaffleNumberStatusAvailable)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	output.AvailableCount = availableCount

	reservedCount, err := uc.raffleNumberRepo.CountByStatus(raffle.ID, domain.RaffleNumberStatusReserved)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	output.ReservedCount = reservedCount

	soldCount, err := uc.raffleNumberRepo.CountByStatus(raffle.ID, domain.RaffleNumberStatusSold)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	output.SoldCount = soldCount

	// Obtener imágenes si se solicita
	if input.IncludeImages {
		images, err := uc.raffleImageRepo.FindByRaffleID(raffle.ID)
		if err != nil {
			return nil, errors.Wrap(errors.ErrDatabaseError, err)
		}
		output.Images = images
	}

	return output, nil
}
