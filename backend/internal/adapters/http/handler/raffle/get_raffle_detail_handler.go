package raffle

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/domain"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
)

// GetRaffleDetailResponse respuesta del detalle
type GetRaffleDetailResponse struct {
	Raffle         *RaffleDTO        `json:"raffle"`
	Numbers        []RaffleNumberDTO `json:"numbers,omitempty"`
	Images         []RaffleImageDTO  `json:"images,omitempty"`
	AvailableCount int64             `json:"available_count"`
	ReservedCount  int64             `json:"reserved_count"`
	SoldCount      int64             `json:"sold_count"`
}

// RaffleImageDTO representa una imagen de sorteo
type RaffleImageDTO struct {
	ID           int64  `json:"id"`
	Filename     string `json:"filename"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
	Width        *int   `json:"width,omitempty"`
	Height       *int   `json:"height,omitempty"`
	AltText      string `json:"alt_text,omitempty"`
	DisplayOrder int    `json:"display_order"`
	IsPrimary    bool   `json:"is_primary"`
}

// GetRaffleDetailHandler maneja la obtención de detalle de sorteo
type GetRaffleDetailHandler struct {
	useCase *raffleuc.GetRaffleDetailUseCase
}

// NewGetRaffleDetailHandler crea una nueva instancia
func NewGetRaffleDetailHandler(useCase *raffleuc.GetRaffleDetailUseCase) *GetRaffleDetailHandler {
	return &GetRaffleDetailHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *GetRaffleDetailHandler) Handle(c *gin.Context) {
	// 1. Obtener ID o UUID del path
	idOrUUID := c.Param("id")

	// 2. Construir input
	input := &raffleuc.GetRaffleDetailInput{
		IncludeNumbers: c.Query("include_numbers") == "true",
		IncludeImages:  c.Query("include_images") == "true",
	}

	// Intentar parsear como ID numérico
	if id, err := strconv.ParseInt(idOrUUID, 10, 64); err == nil {
		input.RaffleID = &id
	} else {
		// Es un UUID
		input.RaffleUUID = &idOrUUID
	}

	// 3. Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 4. Construir response
	response := &GetRaffleDetailResponse{
		Raffle:         toRaffleDTO(output.Raffle),
		AvailableCount: output.AvailableCount,
		ReservedCount:  output.ReservedCount,
		SoldCount:      output.SoldCount,
	}

	if output.Numbers != nil {
		response.Numbers = toRaffleNumberDTOs(output.Numbers)
	}

	if output.Images != nil {
		response.Images = toRaffleImageDTOs(output.Images)
	}

	c.JSON(http.StatusOK, response)
}

// toRaffleImageDTOs convierte slice de RaffleImage a DTOs
func toRaffleImageDTOs(images []*domain.RaffleImage) []RaffleImageDTO {
	dtos := make([]RaffleImageDTO, len(images))
	for i, img := range images {
		dtos[i] = RaffleImageDTO{
			ID:           img.ID,
			Filename:     img.Filename,
			FileSize:     img.FileSize,
			MimeType:     img.MimeType,
			Width:        img.Width,
			Height:       img.Height,
			AltText:      img.AltText,
			DisplayOrder: img.DisplayOrder,
			IsPrimary:    img.IsPrimary,
		}
	}
	return dtos
}
