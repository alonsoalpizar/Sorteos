package raffle

import (
	"net/http"

	"github.com/gin-gonic/gin"

	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// ListRaffleBuyersHandler maneja la obtención de compradores de un sorteo
type ListRaffleBuyersHandler struct {
	useCase *raffleuc.ListRaffleBuyersUseCase
}

// NewListRaffleBuyersHandler crea una nueva instancia
func NewListRaffleBuyersHandler(useCase *raffleuc.ListRaffleBuyersUseCase) *ListRaffleBuyersHandler {
	return &ListRaffleBuyersHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
// GET /api/v1/raffles/:id/buyers?include_sold=true&include_reserved=true
func (h *ListRaffleBuyersHandler) Handle(c *gin.Context) {
	// 1. Obtener UUID del sorteo
	raffleUUID := c.Param("id")
	if raffleUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de sorteo requerido",
		})
		return
	}

	// 2. Obtener user_id del contexto (autenticado)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "No autenticado",
		})
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error interno",
		})
		return
	}

	// 3. Obtener query params
	includeSold := c.DefaultQuery("include_sold", "true") == "true"
	includeReserved := c.DefaultQuery("include_reserved", "true") == "true"

	// 4. Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), &raffleuc.ListRaffleBuyersInput{
		RaffleUUID:      raffleUUID,
		OwnerUserID:     userIDInt64,
		IncludeSold:     includeSold,
		IncludeReserved: includeReserved,
	})

	if err != nil {
		switch err {
		case errors.ErrRaffleNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Sorteo no encontrado",
			})
		case errors.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Solo el organizador puede ver la lista de compradores",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Error interno del servidor",
			})
		}
		return
	}

	// 5. Retornar resultado
	c.JSON(http.StatusOK, output)
}
