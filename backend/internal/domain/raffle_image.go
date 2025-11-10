package domain

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// RaffleImage representa una imagen de un sorteo
type RaffleImage struct {
	ID int64

	// Raffle reference
	RaffleID int64

	// Image info
	Filename         string
	OriginalFilename string
	FilePath         string
	FileSize         int64
	MimeType         string

	// Image metadata
	Width   *int
	Height  *int
	AltText string

	// Ordering
	DisplayOrder int
	IsPrimary    bool

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Valid MIME types for raffle images
var ValidImageMimeTypes = []string{
	"image/jpeg",
	"image/jpg",
	"image/png",
	"image/webp",
	"image/gif",
}

// Max file size: 10 MB
const MaxImageFileSize = 10 * 1024 * 1024

// NewRaffleImage crea una nueva imagen de sorteo
func NewRaffleImage(raffleID int64, filename, originalFilename, filePath string, fileSize int64, mimeType string) *RaffleImage {
	return &RaffleImage{
		RaffleID:         raffleID,
		Filename:         filename,
		OriginalFilename: originalFilename,
		FilePath:         filePath,
		FileSize:         fileSize,
		MimeType:         mimeType,
		DisplayOrder:     0,
		IsPrimary:        false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Validate valida la imagen
func (ri *RaffleImage) Validate() error {
	if ri.RaffleID <= 0 {
		return fmt.Errorf("raffle_id es requerido")
	}

	if ri.Filename == "" {
		return fmt.Errorf("el nombre del archivo es requerido")
	}

	if ri.FilePath == "" {
		return fmt.Errorf("la ruta del archivo es requerida")
	}

	if ri.FileSize <= 0 {
		return fmt.Errorf("el tamaño del archivo debe ser mayor a 0")
	}

	if ri.FileSize > MaxImageFileSize {
		return fmt.Errorf("el archivo excede el tamaño máximo permitido (10 MB)")
	}

	if !ri.IsValidMimeType() {
		return fmt.Errorf("tipo de archivo no permitido: %s", ri.MimeType)
	}

	if ri.DisplayOrder < 0 {
		return fmt.Errorf("el orden de visualización no puede ser negativo")
	}

	return nil
}

// IsValidMimeType verifica si el MIME type es válido
func (ri *RaffleImage) IsValidMimeType() bool {
	for _, validType := range ValidImageMimeTypes {
		if ri.MimeType == validType {
			return true
		}
	}
	return false
}

// GetFileExtension retorna la extensión del archivo
func (ri *RaffleImage) GetFileExtension() string {
	return filepath.Ext(ri.Filename)
}

// IsImage verifica si es un archivo de imagen válido
func (ri *RaffleImage) IsImage() bool {
	return ri.IsValidMimeType()
}

// SetPrimary marca la imagen como principal
func (ri *RaffleImage) SetPrimary() {
	ri.IsPrimary = true
	ri.UpdatedAt = time.Now()
}

// UnsetPrimary desmarca la imagen como principal
func (ri *RaffleImage) UnsetPrimary() {
	ri.IsPrimary = false
	ri.UpdatedAt = time.Now()
}

// SetDimensions establece las dimensiones de la imagen
func (ri *RaffleImage) SetDimensions(width, height int) {
	ri.Width = &width
	ri.Height = &height
	ri.UpdatedAt = time.Now()
}

// SetAltText establece el texto alternativo
func (ri *RaffleImage) SetAltText(altText string) {
	ri.AltText = altText
	ri.UpdatedAt = time.Now()
}

// SetDisplayOrder establece el orden de visualización
func (ri *RaffleImage) SetDisplayOrder(order int) error {
	if order < 0 {
		return fmt.Errorf("el orden no puede ser negativo")
	}

	ri.DisplayOrder = order
	ri.UpdatedAt = time.Now()

	return nil
}

// SoftDelete marca la imagen como eliminada (soft delete)
func (ri *RaffleImage) SoftDelete() {
	now := time.Now()
	ri.DeletedAt = &now
	ri.UpdatedAt = now
}

// IsDeleted verifica si la imagen está eliminada
func (ri *RaffleImage) IsDeleted() bool {
	return ri.DeletedAt != nil
}

// GetURL genera una URL pública para la imagen
// Esto puede ser personalizado según la implementación (local, S3, CDN)
func (ri *RaffleImage) GetURL(baseURL string) string {
	// Por ahora retorna una URL simple, puede ser mejorado con CDN
	return fmt.Sprintf("%s/uploads/raffles/%d/%s", baseURL, ri.RaffleID, ri.Filename)
}

// GetThumbnailURL genera una URL para la miniatura
func (ri *RaffleImage) GetThumbnailURL(baseURL string) string {
	// Asume que las miniaturas tienen el prefijo "thumb_"
	ext := ri.GetFileExtension()
	nameWithoutExt := strings.TrimSuffix(ri.Filename, ext)
	thumbnailName := fmt.Sprintf("thumb_%s%s", nameWithoutExt, ext)

	return fmt.Sprintf("%s/uploads/raffles/%d/%s", baseURL, ri.RaffleID, thumbnailName)
}

// GetImageInfo retorna información básica de la imagen
func (ri *RaffleImage) GetImageInfo() map[string]interface{} {
	info := map[string]interface{}{
		"id":                ri.ID,
		"filename":          ri.Filename,
		"original_filename": ri.OriginalFilename,
		"file_size":         ri.FileSize,
		"mime_type":         ri.MimeType,
		"display_order":     ri.DisplayOrder,
		"is_primary":        ri.IsPrimary,
	}

	if ri.Width != nil && ri.Height != nil {
		info["width"] = *ri.Width
		info["height"] = *ri.Height
	}

	if ri.AltText != "" {
		info["alt_text"] = ri.AltText
	}

	return info
}
