package image

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// ImageVariant representa una variante de imagen
type ImageVariant string

const (
	VariantOriginal  ImageVariant = "original"
	VariantLarge     ImageVariant = "large"
	VariantMedium    ImageVariant = "medium"
	VariantThumbnail ImageVariant = "thumbnail"
)

// VariantConfig configuración para cada variante
type VariantConfig struct {
	MaxWidth  int
	MaxHeight int
	Quality   int
	Format    string // "webp" o "original"
}

// DefaultVariants configuración predeterminada de variantes
var DefaultVariants = map[ImageVariant]VariantConfig{
	VariantOriginal: {
		MaxWidth:  1200,
		MaxHeight: 1200,
		Quality:   90,
		Format:    "original",
	},
	VariantLarge: {
		MaxWidth:  800,
		MaxHeight: 800,
		Quality:   85,
		Format:    "webp",
	},
	VariantMedium: {
		MaxWidth:  400,
		MaxHeight: 400,
		Quality:   80,
		Format:    "webp",
	},
	VariantThumbnail: {
		MaxWidth:  150,
		MaxHeight: 150,
		Quality:   75,
		Format:    "webp",
	},
}

// ProcessedImage representa una imagen procesada con sus variantes
type ProcessedImage struct {
	OriginalPath  string
	LargePath     string
	MediumPath    string
	ThumbnailPath string
	OriginalURL   string
	LargeURL      string
	MediumURL     string
	ThumbnailURL  string
	Width         int
	Height        int
	FileSize      int64
	MimeType      string
}

// ImageProcessor servicio para procesar imágenes
type ImageProcessor struct {
	uploadDir string
	baseURL   string
}

// NewImageProcessor crea un nuevo procesador de imágenes
func NewImageProcessor(uploadDir, baseURL string) *ImageProcessor {
	return &ImageProcessor{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// ProcessImage procesa una imagen y genera todas las variantes
func (p *ImageProcessor) ProcessImage(sourcePath string, raffleID int64) (*ProcessedImage, error) {
	// Abrir imagen fuente
	sourceImg, err := imaging.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo imagen: %w", err)
	}

	// Obtener dimensiones originales
	bounds := sourceImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Generar UUID para los archivos
	fileUUID := uuid.New().String()

	// Crear directorio para el sorteo
	raffleDir := filepath.Join(p.uploadDir, fmt.Sprintf("%d", raffleID))
	if err := os.MkdirAll(raffleDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio: %w", err)
	}

	result := &ProcessedImage{
		Width:    width,
		Height:   height,
		MimeType: "image/webp", // Guardamos principalmente en WebP
	}

	// Procesar cada variante
	for variant, config := range DefaultVariants {
		variantDir := filepath.Join(raffleDir, string(variant))
		if err := os.MkdirAll(variantDir, 0755); err != nil {
			return nil, fmt.Errorf("error creando directorio de variante: %w", err)
		}

		// Redimensionar imagen
		resized := p.resizeImage(sourceImg, config.MaxWidth, config.MaxHeight)

		// Determinar extensión y formato
		ext := ".webp"
		if config.Format == "original" {
			// Mantener formato original para la versión "original"
			ext = filepath.Ext(sourcePath)
			if ext == "" {
				ext = ".jpg"
			}
		}

		// Nombre del archivo
		filename := fileUUID + ext
		outputPath := filepath.Join(variantDir, filename)

		// Guardar imagen
		if config.Format == "webp" {
			if err := p.saveAsWebP(resized, outputPath, config.Quality); err != nil {
				return nil, fmt.Errorf("error guardando WebP: %w", err)
			}
		} else {
			// Guardar en formato original (JPEG/PNG)
			if err := imaging.Save(resized, outputPath, imaging.JPEGQuality(config.Quality)); err != nil {
				return nil, fmt.Errorf("error guardando imagen original: %w", err)
			}
		}

		// Generar URL relativa (para que funcione en cualquier dominio)
		url := fmt.Sprintf("/uploads/raffles/%d/%s/%s", raffleID, variant, filename)

		// Asignar rutas y URLs según variante
		switch variant {
		case VariantOriginal:
			result.OriginalPath = outputPath
			result.OriginalURL = url
		case VariantLarge:
			result.LargePath = outputPath
			result.LargeURL = url
		case VariantMedium:
			result.MediumPath = outputPath
			result.MediumURL = url
		case VariantThumbnail:
			result.ThumbnailPath = outputPath
			result.ThumbnailURL = url
		}
	}

	// Obtener tamaño del archivo original
	fileInfo, err := os.Stat(result.OriginalPath)
	if err == nil {
		result.FileSize = fileInfo.Size()
	}

	return result, nil
}

// resizeImage redimensiona una imagen manteniendo el aspect ratio
func (p *ImageProcessor) resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Si la imagen ya es más pequeña, retornarla sin cambios
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	// Calcular nuevas dimensiones manteniendo aspect ratio
	ratio := float64(width) / float64(height)
	newWidth := maxWidth
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > maxHeight {
		newHeight = maxHeight
		newWidth = int(float64(newHeight) * ratio)
	}

	// Redimensionar con filtro Lanczos (mejor calidad)
	return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
}

// saveAsWebP guarda una imagen en formato WebP
func (p *ImageProcessor) saveAsWebP(img image.Image, outputPath string, quality int) error {
	// Crear archivo de salida
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creando archivo: %w", err)
	}
	defer outFile.Close()

	// Configurar opciones de WebP
	options := &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	}

	// Codificar y guardar
	if err := webp.Encode(outFile, img, options); err != nil {
		return fmt.Errorf("error codificando WebP: %w", err)
	}

	return nil
}

// DeleteVariants elimina todas las variantes de una imagen
func (p *ImageProcessor) DeleteVariants(raffleID int64, filename string) error {
	raffleDir := filepath.Join(p.uploadDir, fmt.Sprintf("%d", raffleID))

	// Eliminar cada variante
	for variant := range DefaultVariants {
		variantDir := filepath.Join(raffleDir, string(variant))

		// Buscar y eliminar archivos que coincidan
		files, err := filepath.Glob(filepath.Join(variantDir, "*"))
		if err != nil {
			continue
		}

		for _, file := range files {
			if strings.Contains(filepath.Base(file), strings.TrimSuffix(filename, filepath.Ext(filename))) {
				os.Remove(file)
			}
		}
	}

	return nil
}

// ValidateMimeType valida que el tipo MIME sea permitido
func ValidateMimeType(mimeType string) bool {
	allowed := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"image/gif",
	}

	for _, allowed := range allowed {
		if mimeType == allowed {
			return true
		}
	}

	return false
}
