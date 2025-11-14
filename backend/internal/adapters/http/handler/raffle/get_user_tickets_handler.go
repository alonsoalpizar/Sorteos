package raffle

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/domain"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// UserTicketNumberDTO representa un número de sorteo comprado en la respuesta
type UserTicketNumberDTO struct {
	ID        int64  `json:"id"`
	Number    string `json:"number"`
	Price     string `json:"price"`
	SoldAt    string `json:"sold_at"`
	PaymentID *int64 `json:"payment_id,omitempty"`
}

// TicketGroupDTO representa un grupo de tickets por sorteo
type TicketGroupDTO struct {
	Raffle       *RaffleDTO             `json:"raffle"`
	Numbers      []*UserTicketNumberDTO `json:"numbers"`
	TotalNumbers int                    `json:"total_numbers"`
	TotalSpent   string                 `json:"total_spent"`
}

// GetUserTicketsResponse respuesta del listado de tickets
type GetUserTicketsResponse struct {
	Tickets    []*TicketGroupDTO `json:"tickets"`
	Pagination Pagination        `json:"pagination"`
}

// GetUserTicketsHandler maneja la obtención de tickets del usuario
type GetUserTicketsHandler struct {
	useCase *raffleuc.GetUserTicketsUseCase
}

// NewGetUserTicketsHandler crea una nueva instancia
func NewGetUserTicketsHandler(useCase *raffleuc.GetUserTicketsUseCase) *GetUserTicketsHandler {
	return &GetUserTicketsHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *GetUserTicketsHandler) Handle(c *gin.Context) {
	// 1. Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id type"})
		return
	}

	// 2. Parsear query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. Construir input
	input := &raffleuc.GetUserTicketsInput{
		UserID:   userIDInt64,
		Page:     page,
		PageSize: pageSize,
	}

	// 4. Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 5. Construir response
	ticketGroups := make([]*TicketGroupDTO, len(output.Tickets))
	for i, group := range output.Tickets {
		numbers := make([]*UserTicketNumberDTO, len(group.Numbers))
		for j, num := range group.Numbers {
			numbers[j] = toUserTicketNumberDTO(num)
		}

		ticketGroups[i] = &TicketGroupDTO{
			Raffle:       toRaffleDTO(group.Raffle),
			Numbers:      numbers,
			TotalNumbers: group.TotalNumbers,
			TotalSpent:   group.TotalSpent,
		}
	}

	response := &GetUserTicketsResponse{
		Tickets: ticketGroups,
		Pagination: Pagination{
			Page:       output.Page,
			PageSize:   output.PageSize,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// toUserTicketNumberDTO convierte un RaffleNumber a UserTicketNumberDTO
func toUserTicketNumberDTO(num *domain.RaffleNumber) *UserTicketNumberDTO {
	dto := &UserTicketNumberDTO{
		ID:     num.ID,
		Number: num.Number,
	}

	if num.Price != nil {
		dto.Price = num.Price.String()
	} else {
		dto.Price = "0"
	}

	if num.SoldAt != nil && !num.SoldAt.IsZero() {
		dto.SoldAt = num.SoldAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if num.PaymentID != nil && *num.PaymentID > 0 {
		dto.PaymentID = num.PaymentID
	}

	return dto
}
