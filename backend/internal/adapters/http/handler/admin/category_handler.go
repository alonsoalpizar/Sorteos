package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/category"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// CategoryHandler maneja todas las operaciones de administración de categorías
type CategoryHandler struct {
	createCategory *category.CreateCategoryUseCase
	updateCategory *category.UpdateCategoryUseCase
	deleteCategory *category.DeleteCategoryUseCase
	listCategories *category.ListCategoriesUseCase
	reorderCategoriesUC *category.ReorderCategoriesUseCase
}

// NewCategoryHandler crea una nueva instancia
func NewCategoryHandler(
	createCategory *category.CreateCategoryUseCase,
	updateCategory *category.UpdateCategoryUseCase,
	deleteCategory *category.DeleteCategoryUseCase,
	listCategories *category.ListCategoriesUseCase,
	reorderCategories *category.ReorderCategoriesUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		createCategory: createCategory,
		updateCategory: updateCategory,
		deleteCategory: deleteCategory,
		listCategories: listCategories,
		reorderCategoriesUC: reorderCategories,
	}
}

// ListCategories maneja GET /api/v1/admin/categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	input := &category.ListCategoriesInput{
		Page:     page,
		PageSize: pageSize,
		Search:   stringPtr(c.Query("search")),
		OrderBy:  stringPtr(c.Query("order_by")),
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		input.IsActive = &isActive
	}

	output, err := h.listCategories.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// CreateCategory maneja POST /api/v1/admin/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
		IconURL     *string `json:"icon_url"`
		IsActive    bool    `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &category.CreateCategoryInput{
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		IsActive:    req.IsActive,
	}

	output, err := h.createCategory.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// UpdateCategory maneja PUT /api/v1/admin/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_CATEGORY_ID", "invalid category ID format", 400, err))
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		IconURL     *string `json:"icon_url"`
		IsActive    *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &category.UpdateCategoryInput{
		CategoryID:  categoryID,
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		IsActive:    req.IsActive,
	}

	output, err := h.updateCategory.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// DeleteCategory maneja DELETE /api/v1/admin/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_CATEGORY_ID", "invalid category ID format", 400, err))
		return
	}

	input := &category.DeleteCategoryInput{
		CategoryID: categoryID,
	}

	output, err := h.deleteCategory.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ReorderCategories reordena las categorías según el array de IDs proporcionado
// POST /api/v1/admin/categories/reorder
func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body
	var input category.ReorderCategoriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.reorderCategoriesUC.Execute(c.Request.Context(), &input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
