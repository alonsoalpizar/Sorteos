package raffle

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/domain"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
)

// ListRafflesResponse respuesta del listado
type ListRafflesResponse struct {
	Raffles    []*RaffleDTO `json:"raffles"`
	Pagination Pagination   `json:"pagination"`
}

// Pagination información de paginación
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ListRafflesHandler maneja el listado de sorteos
type ListRafflesHandler struct {
	useCase *raffleuc.ListRafflesUseCase
}

// NewListRafflesHandler crea una nueva instancia
func NewListRafflesHandler(useCase *raffleuc.ListRafflesUseCase) *ListRafflesHandler {
	return &ListRafflesHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *ListRafflesHandler) Handle(c *gin.Context) {
	// 1. Parsear query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 2. Construir input
	input := &raffleuc.ListRafflesInput{
		Page:          page,
		PageSize:      pageSize,
		Search:        c.Query("search"),
		OrderBy:       c.Query("order_by"),
		OnlyAvailable: c.Query("only_available") == "true",
	}

	// Status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := domain.RaffleStatus(statusStr)
		input.Status = &status
	}

	// User ID filter (para "mis sorteos")
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			input.UserID = &userID
		}
	}

	// 3. Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 4. Construir response
	raffles := make([]*RaffleDTO, len(output.Raffles))
	for i, r := range output.Raffles {
		raffles[i] = toRaffleDTO(r)
	}

	response := &ListRafflesResponse{
		Raffles: raffles,
		Pagination: Pagination{
			Page:       output.Page,
			PageSize:   output.PageSize,
			Total:      output.Total,
			TotalPages: output.TotalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}
