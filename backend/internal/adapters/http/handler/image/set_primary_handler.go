package image

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	imageuc "github.com/sorteos-platform/backend/internal/usecase/image"
)

// SetPrimaryImageHandler maneja establecer imagen primaria
type SetPrimaryImageHandler struct {
	useCase *imageuc.SetPrimaryImageUseCase
}

// NewSetPrimaryImageHandler crea una nueva instancia
func NewSetPrimaryImageHandler(useCase *imageuc.SetPrimaryImageUseCase) *SetPrimaryImageHandler {
	return &SetPrimaryImageHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *SetPrimaryImageHandler) Handle(c *gin.Context) {
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

	// 2. Obtener image_id del path
	imageIDStr := c.Param("image_id")
	imageID, err := strconv.ParseInt(imageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_IMAGE_ID",
			"message": "ID de imagen inválido",
		})
		return
	}

	// 3. Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "UNAUTHORIZED",
			"message": "Usuario no autenticado",
		})
		return
	}

	// 4. Ejecutar use case
	input := &imageuc.SetPrimaryImageInput{
		ImageID:  imageID,
		RaffleID: raffleID,
		UserID:   userID.(int64),
	}

	if err := h.useCase.Execute(c.Request.Context(), input); err != nil {
		handleError(c, err)
		return
	}

	// 5. Response exitoso
	c.JSON(http.StatusOK, gin.H{
		"message": "Imagen primaria establecida exitosamente",
	})
}
