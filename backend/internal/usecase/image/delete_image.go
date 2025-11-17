package image

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/infrastructure/image"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// DeleteImageInput datos de entrada para eliminar imagen
type DeleteImageInput struct {
	ImageID  int64
	RaffleID int64
	UserID   int64
}

// DeleteImageUseCase caso de uso para eliminar imágenes
type DeleteImageUseCase struct {
	raffleRepo     db.RaffleRepository
	imageRepo      db.RaffleImageRepository
	imageProcessor *image.ImageProcessor
	logger         *logger.Logger
	uploadDir      string
}

// NewDeleteImageUseCase crea una nueva instancia
func NewDeleteImageUseCase(
	raffleRepo db.RaffleRepository,
	imageRepo db.RaffleImageRepository,
	logger *logger.Logger,
	uploadDir string,
) *DeleteImageUseCase {
	return &DeleteImageUseCase{
		raffleRepo:     raffleRepo,
		imageRepo:      imageRepo,
		imageProcessor: image.NewImageProcessor(uploadDir, ""),
		logger:         logger,
		uploadDir:      uploadDir,
	}
}

// Execute ejecuta el caso de uso
func (uc *DeleteImageUseCase) Execute(ctx context.Context, input *DeleteImageInput) error {
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
		uc.logger.Warn("Usuario no autorizado para eliminar imagen",
			logger.Int64("user_id", input.UserID),
			logger.Int64("raffle_owner", raffle.UserID))
		return errors.ErrUnauthorized
	}

	// 4. Si es la única imagen primaria, no permitir eliminarla
	// (La base de datos tiene un trigger para esto, pero validamos también aquí)
	if img.IsPrimary {
		count, err := uc.imageRepo.CountByRaffleID(input.RaffleID)
		if err != nil {
			return err
		}

		if count == 1 {
			// Está bien eliminar la única imagen
		} else {
			// Si hay más imágenes, primero se debe establecer otra como primaria
			return errors.New("CANNOT_DELETE_PRIMARY", "No se puede eliminar la imagen primaria. Establece otra imagen como primaria primero.", 400, nil)
		}
	}

	// 5. Soft delete en base de datos
	if err := uc.imageRepo.SoftDelete(input.ImageID); err != nil {
		return err
	}

	// 6. Eliminar archivos físicos
	if err := uc.imageProcessor.DeleteVariants(img.RaffleID, img.Filename); err != nil {
		// Log error pero no falla la operación
		uc.logger.Error("Error eliminando archivos físicos",
			logger.Int64("image_id", input.ImageID),
			logger.Error(err))
	}

	uc.logger.Info("Imagen eliminada exitosamente",
		logger.Int64("raffle_id", input.RaffleID),
		logger.Int64("image_id", input.ImageID))

	return nil
}
