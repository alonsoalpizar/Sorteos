package image

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	imageuc "github.com/sorteos-platform/backend/internal/usecase/image"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// ImageDTO representa una imagen en el response
type ImageDTO struct {
	ID               int64   `json:"id"`
	RaffleID         int64   `json:"raffle_id"`
	Filename         string  `json:"filename"`
	OriginalFilename string  `json:"original_filename"`
	URLOriginal      *string `json:"url_original,omitempty"`
	URLLarge         *string `json:"url_large,omitempty"`
	URLMedium        *string `json:"url_medium,omitempty"`
	URLThumbnail     *string `json:"url_thumbnail,omitempty"`
	Width            *int    `json:"width,omitempty"`
	Height           *int    `json:"height,omitempty"`
	FileSize         int64   `json:"file_size"`
	DisplayOrder     int     `json:"display_order"`
	IsPrimary        bool    `json:"is_primary"`
	CreatedAt        string  `json:"created_at"`
}

// UploadImageResponse respuesta del upload
type UploadImageResponse struct {
	Image *ImageDTO `json:"image"`
}

// UploadImageHandler maneja la subida de imágenes
type UploadImageHandler struct {
	useCase *imageuc.UploadImageUseCase
}

// NewUploadImageHandler crea una nueva instancia
func NewUploadImageHandler(useCase *imageuc.UploadImageUseCase) *UploadImageHandler {
	return &UploadImageHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *UploadImageHandler) Handle(c *gin.Context) {
	// 1. Obtener raffle_id del path
	raffleIDStr := c.Param("id")
	raffleID, err := strconv.ParseInt(raffleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_RAFFLE_ID",
			"message": "ID de sorteo inválido",
		})
		return
	}

	// 2. Obtener user_id del contexto (middleware de autenticación)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "UNAUTHORIZED",
			"message": "Usuario no autenticado",
		})
		return
	}

	// 3. Obtener archivo del form
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "MISSING_FILE",
			"message": "Archivo de imagen requerido",
		})
		return
	}
	defer file.Close()

	// 4. Validar tamaño
	if header.Size > 10*1024*1024 { // 10 MB
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "FILE_TOO_LARGE",
			"message": "El archivo excede el tamaño máximo (10 MB)",
		})
		return
	}

	// 5. Ejecutar use case
	input := &imageuc.UploadImageInput{
		RaffleID:    raffleID,
		UserID:      userID.(int64),
		File:        file,
		Filename:    header.Filename,
		FileSize:    header.Size,
		ContentType: header.Header.Get("Content-Type"),
	}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 6. Construir response
	dto := &ImageDTO{
		ID:               output.Image.ID,
		RaffleID:         output.Image.RaffleID,
		Filename:         output.Image.Filename,
		OriginalFilename: output.Image.OriginalFilename,
		URLOriginal:      output.Image.URLOriginal,
		URLLarge:         output.Image.URLLarge,
		URLMedium:        output.Image.URLMedium,
		URLThumbnail:     output.Image.URLThumbnail,
		Width:            output.Image.Width,
		Height:           output.Image.Height,
		FileSize:         output.Image.FileSize,
		DisplayOrder:     output.Image.DisplayOrder,
		IsPrimary:        output.Image.IsPrimary,
		CreatedAt:        output.Image.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, &UploadImageResponse{
		Image: dto,
	})
}

// handleError maneja errores de forma consistente
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Status, gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	// Error genérico
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    "INTERNAL_SERVER_ERROR",
		"message": err.Error(),
	})
}
