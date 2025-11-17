package image

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// SetPrimaryImageInput datos de entrada
type SetPrimaryImageInput struct {
	ImageID  int64
	RaffleID int64
	UserID   int64
}

// SetPrimaryImageUseCase caso de uso para establecer imagen primaria
type SetPrimaryImageUseCase struct {
	raffleRepo db.RaffleRepository
	imageRepo  db.RaffleImageRepository
	logger     *logger.Logger
}

// NewSetPrimaryImageUseCase crea una nueva instancia
func NewSetPrimaryImageUseCase(
	raffleRepo db.RaffleRepository,
	imageRepo db.RaffleImageRepository,
	logger *logger.Logger,
) *SetPrimaryImageUseCase {
	return &SetPrimaryImageUseCase{
		raffleRepo: raffleRepo,
		imageRepo:  imageRepo,
		logger:     logger,
	}
}

// Execute ejecuta el caso de uso
func (uc *SetPrimaryImageUseCase) Execute(ctx context.Context, input *SetPrimaryImageInput) error {
	// 1. Buscar la imagen
	img, err := uc.imageRepo.FindByID(input.ImageID)
	if err != nil {
		return err
	}

	// 2. Validar que pertenece al sorteo especificado
	if img.RaffleID != input.RaffleID {
		return errors.ErrNotFound
	}

	// 3. Validar ownership
	raffle, err := uc.raffleRepo.FindByID(input.RaffleID)
	if err != nil {
		return err
	}

	if raffle.UserID != input.UserID {
		uc.logger.Warn("Usuario no autorizado para cambiar imagen primaria",
			logger.Int64("user_id", input.UserID),
			logger.Int64("raffle_owner", raffle.UserID))
		return errors.ErrUnauthorized
	}

	// 4. Si ya es primaria, no hacer nada
	if img.IsPrimary {
		return nil
	}

	// 5. Establecer como primaria (el repository manejará quitar is_primary de las demás)
	if err := uc.imageRepo.SetPrimary(input.ImageID); err != nil {
		return err
	}

	uc.logger.Info("Imagen establecida como primaria",
		logger.Int64("raffle_id", input.RaffleID),
		logger.Int64("image_id", input.ImageID))

	return nil
}
