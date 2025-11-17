package category

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sorteos-platform/backend/internal/domain"
	categoryuc "github.com/sorteos-platform/backend/internal/usecase/category"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// CategoryDTO representa una categoría en el response
type CategoryDTO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Icon         string `json:"icon"`
	Description  string `json:"description,omitempty"`
	DisplayOrder int    `json:"display_order"`
}

// ListCategoriesResponse respuesta del listado
type ListCategoriesResponse struct {
	Categories []*CategoryDTO `json:"categories"`
}

// ListCategoriesHandler maneja el listado de categorías
type ListCategoriesHandler struct {
	useCase *categoryuc.ListCategoriesUseCase
}

// NewListCategoriesHandler crea una nueva instancia
func NewListCategoriesHandler(useCase *categoryuc.ListCategoriesUseCase) *ListCategoriesHandler {
	return &ListCategoriesHandler{
		useCase: useCase,
	}
}

// Handle maneja el request
func (h *ListCategoriesHandler) Handle(c *gin.Context) {
	// Ejecutar use case
	categories, err := h.useCase.Execute(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir response
	dtos := make([]*CategoryDTO, len(categories))
	for i, cat := range categories {
		dtos[i] = toCategoryDTO(cat)
	}

	response := &ListCategoriesResponse{
		Categories: dtos,
	}

	c.JSON(http.StatusOK, response)
}

// toCategoryDTO convierte domain.Category a DTO
func toCategoryDTO(cat *domain.Category) *CategoryDTO {
	return &CategoryDTO{
		ID:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		Icon:         cat.Icon,
		Description:  cat.Description,
		DisplayOrder: cat.DisplayOrder,
	}
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
