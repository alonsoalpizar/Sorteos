package raffle

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/domain"
	raffleuc "github.com/sorteos-platform/backend/internal/usecase/raffle"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// UpdateRaffleRequest estructura del request
type UpdateRaffleRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	DrawDate    *string `json:"draw_date,omitempty"` // ISO 8601
	DrawMethod  *string `json:"draw_method,omitempty"`
}

// UpdateRaffleHandler maneja la actualización de sorteos
type UpdateRaffleHandler struct {
	useCase *raffleuc.UpdateRaffleUseCase
}

// NewUpdateRaffleHandler crea una nueva instancia
func NewUpdateRaffleHandler(useCase *raffleuc.UpdateRaffleUseCase) *UpdateRaffleHandler {
	return &UpdateRaffleHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *UpdateRaffleHandler) Handle(c *gin.Context) {
	// 1. Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized})
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

	// 3. Parsear request
	var req UpdateRaffleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "VALIDATION_FAILED",
			"message": err.Error(),
		})
		return
	}

	// 4. Construir input
	input := &raffleuc.UpdateRaffleInput{
		RaffleID:    raffleID,
		UserID:      userID.(int64),
		UserRole:    userRole.(domain.UserRole),
		Title:       req.Title,
		Description: req.Description,
	}

	// Parsear fecha si se provee
	if req.DrawDate != nil {
		drawDate, err := time.Parse(time.RFC3339, *req.DrawDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_DATE_FORMAT",
				"message": "La fecha debe estar en formato ISO 8601",
			})
			return
		}
		input.DrawDate = &drawDate
	}

	// DrawMethod
	if req.DrawMethod != nil {
		method := domain.DrawMethod(*req.DrawMethod)
		input.DrawMethod = &method
	}

	// 5. Ejecutar use case
	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 6. Response
	c.JSON(http.StatusOK, gin.H{
		"raffle": toRaffleDTO(output.Raffle),
	})
}

// SuspendRaffleRequest estructura del request
type SuspendRaffleRequest struct {
	Reason string `json:"reason" binding:"required,min=10"`
}

// SuspendRaffleHandler maneja la suspensión de sorteos (admin only)
type SuspendRaffleHandler struct {
	useCase *raffleuc.SuspendRaffleUseCase
}

// NewSuspendRaffleHandler crea una nueva instancia
func NewSuspendRaffleHandler(useCase *raffleuc.SuspendRaffleUseCase) *SuspendRaffleHandler {
	return &SuspendRaffleHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *SuspendRaffleHandler) Handle(c *gin.Context) {
	// 1. Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized})
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

	// 3. Parsear request
	var req SuspendRaffleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "VALIDATION_FAILED",
			"message": err.Error(),
		})
		return
	}

	// 4. Ejecutar use case
	input := &raffleuc.SuspendRaffleInput{
		RaffleID: raffleID,
		UserID:   userID.(int64),
		UserRole: userRole.(domain.UserRole),
		Reason:   req.Reason,
	}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 5. Response
	c.JSON(http.StatusOK, gin.H{
		"raffle": toRaffleDTO(output.Raffle),
	})
}

// DeleteRaffleHandler maneja la eliminación de sorteos (soft delete)
type DeleteRaffleHandler struct {
	useCase *raffleuc.DeleteRaffleUseCase
}

// NewDeleteRaffleHandler crea una nueva instancia
func NewDeleteRaffleHandler(useCase *raffleuc.DeleteRaffleUseCase) *DeleteRaffleHandler {
	return &DeleteRaffleHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *DeleteRaffleHandler) Handle(c *gin.Context) {
	// 1. Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized})
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
	input := &raffleuc.DeleteRaffleInput{
		RaffleID: raffleID,
		UserID:   userID.(int64),
		UserRole: userRole.(domain.UserRole),
	}

	err = h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	// 4. Response
	c.JSON(http.StatusNoContent, nil)
}
