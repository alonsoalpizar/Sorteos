package raffle

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/sorteos-platform/backend/internal/domain"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// CreateRaffleRequest estructura del request
type CreateRaffleRequest struct {
	Title                 string  `json:"title" binding:"required,min=5,max=255"`
	Description           string  `json:"description" binding:"required,min=20"`
	PricePerNumber        float64 `json:"price_per_number" binding:"required,gt=0"`
	TotalNumbers          int     `json:"total_numbers" binding:"required,min=10,max=10000"`
	DrawDate              string  `json:"draw_date" binding:"required"` // ISO 8601
	DrawMethod            string  `json:"draw_method" binding:"required,oneof=loteria_nacional_cr manual random"`
	PlatformFeePercentage *float64 `json:"platform_fee_percentage,omitempty"`
}

// CreateRaffleResponse estructura de la respuesta
type CreateRaffleResponse struct {
	Raffle  *RaffleDTO        `json:"raffle"`
	Numbers []RaffleNumberDTO `json:"numbers,omitempty"`
}

// RaffleDTO representa un sorteo en el response
type RaffleDTO struct {
	ID                    int64   `json:"id"`
	UUID                  string  `json:"uuid"`
	UserID                int64   `json:"user_id"`
	Title                 string  `json:"title"`
	Description           string  `json:"description"`
	Status                string  `json:"status"`
	PricePerNumber        string  `json:"price_per_number"`
	TotalNumbers          int     `json:"total_numbers"`
	DrawDate              string  `json:"draw_date"`
	DrawMethod            string  `json:"draw_method"`
	SoldCount             int     `json:"sold_count"`
	ReservedCount         int     `json:"reserved_count"`
	TotalRevenue          string  `json:"total_revenue"`
	PlatformFeePercentage string  `json:"platform_fee_percentage"`
	PlatformFeeAmount     string  `json:"platform_fee_amount"`
	NetAmount             string  `json:"net_amount"`
	SettlementStatus      string  `json:"settlement_status"`
	CreatedAt             string  `json:"created_at"`
	PublishedAt           *string `json:"published_at,omitempty"`
}

// RaffleNumberDTO representa un número de sorteo
type RaffleNumberDTO struct {
	ID     int64  `json:"id"`
	Number string `json:"number"`
	Status string `json:"status"`
}

// CreateRaffleHandler maneja la creación de sorteos
type CreateRaffleHandler struct {
	useCase *raffleuc.CreateRaffleUseCase
}

// NewCreateRaffleHandler crea una nueva instancia
func NewCreateRaffleHandler(useCase *raffleuc.CreateRaffleUseCase) *CreateRaffleHandler {
	return &CreateRaffleHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *CreateRaffleHandler) Handle(c *gin.Context) {
	// 1. Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized})
		return
	}

	// 2. Parsear request
	var req CreateRaffleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "VALIDATION_FAILED",
			"message": err.Error(),
		})
		return
	}

	// 3. Parsear fecha
	drawDate, err := time.Parse(time.RFC3339, req.DrawDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_DATE_FORMAT",
			"message": "La fecha debe estar en formato ISO 8601",
		})
		return
	}

	// 4. Construir input
	input := &raffleuc.CreateRaffleInput{
		UserID:         userID.(int64),
		Title:          req.Title,
		Description:    req.Description,
		PricePerNumber: decimal.NewFromFloat(req.PricePerNumber),
		TotalNumbers:   req.TotalNumbers,
		DrawDate:       drawDate,
		DrawMethod:     domain.DrawMethod(req.DrawMethod),
	}

	if req.PlatformFeePercentage != nil {
		fee := decimal.NewFromFloat(*req.PlatformFeePercentage)
		input.PlatformFeePercentage = &fee
	}

	// 5. Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 6. Construir response
	response := &CreateRaffleResponse{
		Raffle:  toRaffleDTO(output.Raffle),
		Numbers: toRaffleNumberDTOs(output.Numbers),
	}

	c.JSON(http.StatusCreated, response)
}

// toRaffleDTO convierte domain.Raffle a DTO
func toRaffleDTO(r *domain.Raffle) *RaffleDTO {
	dto := &RaffleDTO{
		ID:                    r.ID,
		UUID:                  r.UUID.String(),
		UserID:                r.UserID,
		Title:                 r.Title,
		Description:           r.Description,
		Status:                string(r.Status),
		PricePerNumber:        r.PricePerNumber.String(),
		TotalNumbers:          r.TotalNumbers,
		DrawDate:              r.DrawDate.Format(time.RFC3339),
		DrawMethod:            string(r.DrawMethod),
		SoldCount:             r.SoldCount,
		ReservedCount:         r.ReservedCount,
		TotalRevenue:          r.TotalRevenue.String(),
		PlatformFeePercentage: r.PlatformFeePercentage.String(),
		PlatformFeeAmount:     r.PlatformFeeAmount.String(),
		NetAmount:             r.NetAmount.String(),
		SettlementStatus:      string(r.SettlementStatus),
		CreatedAt:             r.CreatedAt.Format(time.RFC3339),
	}

	if r.PublishedAt != nil {
		publishedAt := r.PublishedAt.Format(time.RFC3339)
		dto.PublishedAt = &publishedAt
	}

	return dto
}

// toRaffleNumberDTOs convierte slice de RaffleNumber a DTOs
func toRaffleNumberDTOs(numbers []*domain.RaffleNumber) []RaffleNumberDTO {
	dtos := make([]RaffleNumberDTO, len(numbers))
	for i, n := range numbers {
		dtos[i] = RaffleNumberDTO{
			ID:     n.ID,
			Number: n.Number,
			Status: string(n.Status),
		}
	}
	return dtos
}

// handleError maneja los errores de forma consistente
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Status, gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    "INTERNAL_SERVER_ERROR",
		"message": "Error interno del servidor",
	})
}
