package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/organizer"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// OrganizerHandler maneja todas las operaciones de administración de organizadores
type OrganizerHandler struct {
	listOrganizers        *organizer.ListOrganizersUseCase
	viewOrganizerDetails  *organizer.ViewOrganizerDetailsUseCase
	updateCommission      *organizer.UpdateCommissionUseCase
	calculateRevenue      *organizer.CalculateOrganizerRevenueUseCase
}

// NewOrganizerHandler crea una nueva instancia
func NewOrganizerHandler(
	listOrganizers *organizer.ListOrganizersUseCase,
	viewOrganizerDetails *organizer.ViewOrganizerDetailsUseCase,
	updateCommission *organizer.UpdateCommissionUseCase,
	calculateRevenue *organizer.CalculateOrganizerRevenueUseCase,
) *OrganizerHandler {
	return &OrganizerHandler{
		listOrganizers:       listOrganizers,
		viewOrganizerDetails: viewOrganizerDetails,
		updateCommission:     updateCommission,
		calculateRevenue:     calculateRevenue,
	}
}

// ListOrganizers maneja GET /api/v1/admin/organizers
func (h *OrganizerHandler) ListOrganizers(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parsear query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Construir input
	input := &organizer.ListOrganizersInput{
		Page:     page,
		PageSize: pageSize,
		Search:   stringPtr(c.Query("search")),
		OrderBy:  stringPtr(c.Query("order_by")),
	}

	// Filtros opcionales
	if statusStr := c.Query("status"); statusStr != "" {
		input.Status = stringPtr(statusStr)
	}
	if kycLevelStr := c.Query("kyc_level"); kycLevelStr != "" {
		input.KYCLevel = stringPtr(kycLevelStr)
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		input.DateFrom = stringPtr(dateFrom)
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		input.DateTo = stringPtr(dateTo)
	}
	if minRevenue := c.Query("min_revenue"); minRevenue != "" {
		if rev, err := strconv.ParseFloat(minRevenue, 64); err == nil {
			input.MinRevenue = &rev
		}
	}

	// Ejecutar use case
	output, err := h.listOrganizers.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetOrganizerByID maneja GET /api/v1/admin/organizers/:id
func (h *OrganizerHandler) GetOrganizerByID(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	organizerIDStr := c.Param("id")
	organizerID, err := strconv.ParseInt(organizerIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_ORGANIZER_ID", "invalid organizer ID format", 400, err))
		return
	}

	input := &organizer.ViewOrganizerDetailsInput{
		OrganizerID: organizerID,
	}

	output, err := h.viewOrganizerDetails.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdateCommission maneja PUT /api/v1/admin/organizers/:id/commission
func (h *OrganizerHandler) UpdateCommission(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	organizerIDStr := c.Param("id")
	organizerID, err := strconv.ParseInt(organizerIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_ORGANIZER_ID", "invalid organizer ID format", 400, err))
		return
	}

	var req struct {
		NewCommission float64 `json:"new_commission" binding:"required"`
		Reason        *string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &organizer.UpdateCommissionInput{
		OrganizerID:   organizerID,
		NewCommission: req.NewCommission,
		Reason:        req.Reason,
	}

	output, err := h.updateCommission.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// CalculateRevenue maneja POST /api/v1/admin/organizers/:id/revenue
func (h *OrganizerHandler) CalculateRevenue(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	organizerIDStr := c.Param("id")
	organizerID, err := strconv.ParseInt(organizerIDStr, 10, 64)
	if err != nil {
		handleError(c, errors.New("INVALID_ORGANIZER_ID", "invalid organizer ID format", 400, err))
		return
	}

	var req struct {
		DateFrom string  `json:"date_from" binding:"required"`
		DateTo   string  `json:"date_to" binding:"required"`
		GroupBy  *string `json:"group_by"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &organizer.CalculateOrganizerRevenueInput{
		OrganizerID: organizerID,
		DateFrom:    req.DateFrom,
		DateTo:      req.DateTo,
		GroupBy:     req.GroupBy,
	}

	output, err := h.calculateRevenue.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
