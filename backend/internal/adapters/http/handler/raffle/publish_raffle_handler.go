package raffle

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
)

// PublishRaffleHandler maneja la publicación de sorteos
type PublishRaffleHandler struct {
	useCase *raffleuc.PublishRaffleUseCase
}

// NewPublishRaffleHandler crea una nueva instancia
func NewPublishRaffleHandler(useCase *raffleuc.PublishRaffleUseCase) *PublishRaffleHandler {
	return &PublishRaffleHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *PublishRaffleHandler) Handle(c *gin.Context) {
	// 1. Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	// 2. Obtener ID del sorteo
	raffleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "ID de sorteo inválido",
		})
		return
	}

	// 3. Ejecutar use case
	input := &raffleuc.PublishRaffleInput{
		RaffleID: raffleID,
		UserID:   userID.(int64),
	}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 4. Response
	c.JSON(http.StatusOK, gin.H{
		"raffle": toRaffleDTO(output.Raffle),
	})
}
