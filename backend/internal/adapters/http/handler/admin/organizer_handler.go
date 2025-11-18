package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/usecase/admin/organizer"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// OrganizerHandler maneja las peticiones HTTP relacionadas con administración de organizadores
type OrganizerHandler struct {
	listOrganizersUC           *organizer.ListOrganizersUseCase
	getOrganizerDetailUC       *organizer.GetOrganizerDetailUseCase
	updateOrganizerCommissionUC *organizer.UpdateOrganizerCommissionUseCase
	verifyOrganizerUC          *organizer.VerifyOrganizerUseCase
	log                        *logger.Logger
}

// NewOrganizerHandler crea una nueva instancia del handler
func NewOrganizerHandler(gormDB *gorm.DB, log *logger.Logger) *OrganizerHandler {
	// Inicializar repositorio
	organizerRepo := db.NewOrganizerProfileRepository(gormDB, log)

	return &OrganizerHandler{
		listOrganizersUC:           organizer.NewListOrganizersUseCase(organizerRepo, log),
		getOrganizerDetailUC:       organizer.NewGetOrganizerDetailUseCase(organizerRepo, log),
		updateOrganizerCommissionUC: organizer.NewUpdateOrganizerCommissionUseCase(organizerRepo, log),
		verifyOrganizerUC:          organizer.NewVerifyOrganizerUseCase(organizerRepo, log),
		log:                        log,
	}
}

// List lista organizadores con filtros y paginación
// GET /api/v1/admin/organizers
func (h *OrganizerHandler) List(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &organizer.ListOrganizersInput{
		Page:     1,
		PageSize: 20,
		Search:   c.Query("search"),
		OrderBy:  c.Query("order_by"),
	}

	// Parse page y page_size
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			input.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			input.PageSize = pageSize
		}
	}

	// Parse filtros opcionales
	if verifiedStr := c.Query("verified"); verifiedStr != "" {
		verified := verifiedStr == "true"
		input.Verified = &verified
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = &dateTo
	}

	// Ejecutar use case
	output, err := h.listOrganizersUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// GetByID obtiene detalles completos de un organizador
// GET /api/v1/admin/organizers/:id
func (h *OrganizerHandler) GetByID(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse user ID (el organizador se identifica por su user_id)
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	// Ejecutar use case
	output, err := h.getOrganizerDetailUC.Execute(c.Request.Context(), userID, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// UpdateCommission actualiza la comisión personalizada de un organizador
// PUT /api/v1/admin/organizers/:id/commission
func (h *OrganizerHandler) UpdateCommission(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse user ID
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	// Parse body
	var body struct {
		Commission *float64 `json:"commission"` // NULL = usar default global
		Notes      string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	input := &organizer.UpdateOrganizerCommissionInput{
		UserID:     userID,
		Commission: body.Commission,
		Notes:      body.Notes,
	}

	// Ejecutar use case
	err = h.updateOrganizerCommissionUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Organizer commission updated successfully",
	})
}

// Verify verifica un organizador (aprueba su perfil y datos bancarios)
// PUT /api/v1/admin/organizers/:id/verify
func (h *OrganizerHandler) Verify(c *gin.Context) {
	// Obtener admin ID
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse user ID
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	// Parse body (notas opcionales)
	var body struct {
		Notes string `json:"notes"`
	}
	c.ShouldBindJSON(&body)

	input := &organizer.VerifyOrganizerInput{
		UserID: userID,
		Notes:  body.Notes,
	}

	// Ejecutar use case
	err = h.verifyOrganizerUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Organizer verified successfully",
	})
}
