package image

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/internal/infrastructure/image"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// UploadImageInput datos de entrada para subir imagen
type UploadImageInput struct {
	RaffleID    int64
	UserID      int64
	File        multipart.File
	Filename    string
	FileSize    int64
	ContentType string
}

// UploadImageOutput resultado de la subida
type UploadImageOutput struct {
	Image *domain.RaffleImage
}

// UploadImageUseCase caso de uso para subir imágenes
type UploadImageUseCase struct {
	raffleRepo      db.RaffleRepository
	imageRepo       db.RaffleImageRepository
	imageProcessor  *image.ImageProcessor
	logger          *logger.Logger
	maxImagesPerRaffle int
	uploadDir       string
	baseURL         string
}

// NewUploadImageUseCase crea una nueva instancia
func NewUploadImageUseCase(
	raffleRepo db.RaffleRepository,
	imageRepo db.RaffleImageRepository,
	logger *logger.Logger,
	uploadDir, baseURL string,
) *UploadImageUseCase {
	return &UploadImageUseCase{
		raffleRepo:         raffleRepo,
		imageRepo:          imageRepo,
		imageProcessor:     image.NewImageProcessor(uploadDir, baseURL),
		logger:             logger,
		maxImagesPerRaffle: 5,
		uploadDir:          uploadDir,
		baseURL:            baseURL,
	}
}

// Execute ejecuta el caso de uso
func (uc *UploadImageUseCase) Execute(ctx context.Context, input *UploadImageInput) (*UploadImageOutput, error) {
	// 1. Validar que el sorteo existe
	raffle, err := uc.raffleRepo.FindByID(input.RaffleID)
	if err != nil {
		uc.logger.Error("Sorteo no encontrado", logger.Int64("raffle_id", input.RaffleID), logger.Error(err))
		return nil, errors.ErrNotFound
	}

	// 2. Validar ownership (solo el dueño puede subir imágenes)
	if raffle.UserID != input.UserID {
		uc.logger.Warn("Usuario no autorizado para subir imágenes",
			logger.Int64("user_id", input.UserID),
			logger.Int64("raffle_owner", raffle.UserID))
		return nil, errors.ErrUnauthorized
	}

	// 3. Validar cantidad de imágenes
	count, err := uc.imageRepo.CountByRaffleID(input.RaffleID)
	if err != nil {
		return nil, err
	}

	if count >= int64(uc.maxImagesPerRaffle) {
		return nil, fmt.Errorf("máximo %d imágenes permitidas por sorteo", uc.maxImagesPerRaffle)
	}

	// 4. Validar tipo de archivo
	if !image.ValidateMimeType(input.ContentType) {
		return nil, fmt.Errorf("tipo de archivo no permitido: %s", input.ContentType)
	}

	// 5. Validar tamaño de archivo (10 MB)
	if input.FileSize > domain.MaxImageFileSize {
		return nil, fmt.Errorf("archivo excede el tamaño máximo (10 MB)")
	}

	// 6. Guardar archivo temporal
	tempFilename := uuid.New().String() + filepath.Ext(input.Filename)
	tempDir := filepath.Join(uc.uploadDir, "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio temporal: %w", err)
	}

	tempPath := filepath.Join(tempDir, tempFilename)
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return nil, fmt.Errorf("error creando archivo temporal: %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempPath) // Limpiar archivo temporal al final

	// Copiar contenido
	if _, err := tempFile.ReadFrom(input.File); err != nil {
		return nil, fmt.Errorf("error guardando archivo: %w", err)
	}

	// 7. Procesar imagen (generar variantes)
	processedImage, err := uc.imageProcessor.ProcessImage(tempPath, input.RaffleID)
	if err != nil {
		uc.logger.Error("Error procesando imagen", logger.Error(err))
		return nil, fmt.Errorf("error procesando imagen: %w", err)
	}

	// 8. Crear registro en base de datos
	raffleImage := &domain.RaffleImage{
		RaffleID:         input.RaffleID,
		Filename:         filepath.Base(processedImage.OriginalPath),
		OriginalFilename: input.Filename,
		FilePath:         processedImage.OriginalPath,
		FileSize:         processedImage.FileSize,
		MimeType:         processedImage.MimeType,
		Width:            &processedImage.Width,
		Height:           &processedImage.Height,
		URLOriginal:      &processedImage.OriginalURL,
		URLLarge:         &processedImage.LargeURL,
		URLMedium:        &processedImage.MediumURL,
		URLThumbnail:     &processedImage.ThumbnailURL,
		DisplayOrder:     int(count), // Siguiente orden
		IsPrimary:        count == 0,  // Primera imagen es primaria por defecto
	}

	if err := uc.imageRepo.Create(raffleImage); err != nil {
		// Limpiar archivos generados si falla la creación
		uc.imageProcessor.DeleteVariants(input.RaffleID, raffleImage.Filename)
		return nil, err
	}

	uc.logger.Info("Imagen subida exitosamente",
		logger.Int64("raffle_id", input.RaffleID),
		logger.Int64("image_id", raffleImage.ID),
		logger.String("filename", raffleImage.Filename))

	return &UploadImageOutput{
		Image: raffleImage,
	}, nil
}
