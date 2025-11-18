package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/repository"
	"github.com/sorteos-platform/backend/internal/usecase/admin/payment"
	"github.com/sorteos-platform/backend/internal/usecase/admin/system"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// SystemHandler maneja las peticiones HTTP relacionadas con configuración del sistema
type SystemHandler struct {
	getSystemSettingsUC    *system.GetSystemSettingsUseCase
	updateSystemSettingsUC *system.UpdateSystemSettingsUseCase
	updatePaymentProcessorUC *payment.UpdatePaymentProcessorUseCase
	log                    *logger.Logger
}

// NewSystemHandler crea una nueva instancia del handler
func NewSystemHandler(db *gorm.DB, log *logger.Logger) *SystemHandler {
	// Crear repository de configuración
	configRepo := repository.NewSystemConfigRepository(db)

	return &SystemHandler{
		getSystemSettingsUC:    system.NewGetSystemSettingsUseCase(configRepo, log),
		updateSystemSettingsUC: system.NewUpdateSystemSettingsUseCase(configRepo, log),
		updatePaymentProcessorUC: payment.NewUpdatePaymentProcessorUseCase(db, log),
		log:                    log,
	}
}

// ListParameters lista parámetros del sistema (similar a GetSystemSettings)
// GET /api/v1/admin/system/parameters
func (h *SystemHandler) ListParameters(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Construir input desde query params
	input := &system.GetSystemSettingsInput{}

	// Parse filtros opcionales
	if category := c.Query("category"); category != "" {
		input.Category = &category
	}

	// Ejecutar use case
	output, err := h.getSystemSettingsUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// UpdateParameter actualiza un parámetro específico del sistema
// PUT /api/v1/admin/system/parameters/:key
func (h *SystemHandler) UpdateParameter(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Obtener key desde URL
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_KEY",
				"message": "parameter key is required",
			},
		})
		return
	}

	// Parse body
	var body struct {
		Value    interface{} `json:"value" binding:"required"`
		Category string      `json:"category" binding:"required"`
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

	input := &system.UpdateSystemSettingsInput{
		Key:      key,
		Value:    body.Value,
		Category: body.Category,
	}

	// Ejecutar use case
	output, err := h.updateSystemSettingsUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// GetCompanySettings obtiene configuración de la empresa
// GET /api/v1/admin/system/company
func (h *SystemHandler) GetCompanySettings(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Filtrar solo configuraciones de categoría "company"
	category := "company"
	input := &system.GetSystemSettingsInput{
		Category: &category,
	}

	// Ejecutar use case
	output, err := h.getSystemSettingsUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// UpdateCompanySettings actualiza configuración de la empresa
// PUT /api/v1/admin/system/company
func (h *SystemHandler) UpdateCompanySettings(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Parse body - puede incluir múltiples settings
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INPUT",
				"message": err.Error(),
			},
		})
		return
	}

	// Actualizar cada configuración de empresa
	results := make(map[string]interface{})
	for key, value := range body {
		input := &system.UpdateSystemSettingsInput{
			Key:      key,
			Value:    value,
			Category: "company",
		}

		output, err := h.updateSystemSettingsUC.Execute(c.Request.Context(), input, adminID)
		if err != nil {
			handleError(c, err)
			return
		}
		results[key] = output
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
		"message": "Company settings updated successfully",
	})
}

// ListPaymentProcessors lista configuraciones de payment processors
// GET /api/v1/admin/system/payment-processors
func (h *SystemHandler) ListPaymentProcessors(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Listar payment processors configurados
	processors, err := h.updatePaymentProcessorUC.ListPaymentProcessors()
	if err != nil {
		handleError(c, err)
		return
	}

	// Log auditoría
	h.log.Info("Admin listed payment processors",
		logger.Int64("admin_id", adminID),
		logger.Int("count", len(processors)),
		logger.String("action", "admin_list_payment_processors"))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"processors": processors,
			"total":      len(processors),
		},
	})
}

// UpdatePaymentProcessor actualiza configuración de un payment processor
// PUT /api/v1/admin/system/payment-processors/:processor
func (h *SystemHandler) UpdatePaymentProcessor(c *gin.Context) {
	// Obtener admin ID del contexto
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	// Obtener processor desde URL
	processor := c.Param("processor")
	if processor == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_PROCESSOR",
				"message": "processor name is required",
			},
		})
		return
	}

	// Parse body
	var body struct {
		Enabled  bool                   `json:"enabled"`
		Config   map[string]interface{} `json:"config"`
		Priority *int                   `json:"priority,omitempty"`
		Notes    string                 `json:"notes,omitempty"`
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

	input := &payment.UpdatePaymentProcessorInput{
		Processor: processor,
		Enabled:   body.Enabled,
		Config:    body.Config,
		Priority:  body.Priority,
		Notes:     body.Notes,
	}

	// Ejecutar use case
	output, err := h.updatePaymentProcessorUC.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
		"message": "Payment processor configuration updated successfully",
	})
}
