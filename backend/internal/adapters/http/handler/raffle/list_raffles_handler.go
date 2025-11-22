package raffle

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/domain"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
)

// ListRafflesResponse respuesta del listado público
type ListRafflesResponse struct {
	Raffles    []*PublicRaffleDTO `json:"raffles"` // DTO público sin info financiera
	Pagination Pagination         `json:"pagination"`
}

// OwnerListRafflesResponse respuesta del listado para el organizador
type OwnerListRafflesResponse struct {
	Raffles    []*OwnerRaffleDTO `json:"raffles"` // DTO con información financiera
	Pagination Pagination        `json:"pagination"`
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

	// 2. Obtener usuario autenticado (si existe)
	var authenticatedUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(int64); ok {
			authenticatedUserID = &uid
		}
	}

	// 3. Construir input
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
	var requestedUserID *int64
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			input.UserID = &userID
			requestedUserID = &userID
		}
	}

	// Category ID filter (para filtrar por categoría)
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			input.CategoryID = &categoryID
		}
	}

	// Exclude user ID filter (para /explore - excluir sorteos propios)
	if c.Query("exclude_mine") == "true" && authenticatedUserID != nil {
		input.ExcludeUserID = authenticatedUserID
	}

	// 4. Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 5. Determinar si el usuario está consultando sus propios sorteos
	isOwnerQuery := authenticatedUserID != nil && requestedUserID != nil && *authenticatedUserID == *requestedUserID

	// 6. Construir response con DTOs apropiados
	if isOwnerQuery {
		// Usuario consultando SUS propios sorteos - incluir información financiera
		ownerRaffles := make([]*OwnerRaffleDTO, len(output.Raffles))
		for i, r := range output.Raffles {
			ownerRaffles[i] = toOwnerRaffleDTO(r)
		}

		response := &OwnerListRafflesResponse{
			Raffles: ownerRaffles,
			Pagination: Pagination{
				Page:       output.Page,
				PageSize:   output.PageSize,
				Total:      output.Total,
				TotalPages: output.TotalPages,
			},
		}
		c.JSON(http.StatusOK, response)
	} else {
		// Listado público o de otros usuarios - sin información financiera
		raffles := make([]*PublicRaffleDTO, len(output.Raffles))
		for i, r := range output.Raffles {
			raffles[i] = toPublicRaffleDTO(r)
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
}
