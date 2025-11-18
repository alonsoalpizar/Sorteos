package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/config"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// ConfigHandler maneja todas las operaciones de configuración del sistema
type ConfigHandler struct {
	getConfig    *config.GetSystemConfigUseCase
	updateConfig *config.UpdateSystemConfigUseCase
	listConfigs  *config.ListSystemConfigsUseCase
}

// NewConfigHandler crea una nueva instancia
func NewConfigHandler(
	getConfig *config.GetSystemConfigUseCase,
	updateConfig *config.UpdateSystemConfigUseCase,
	listConfigs *config.ListSystemConfigsUseCase,
) *ConfigHandler {
	return &ConfigHandler{
		getConfig:    getConfig,
		updateConfig: updateConfig,
		listConfigs:  listConfigs,
	}
}

// ListConfigs maneja GET /api/v1/admin/config
func (h *ConfigHandler) ListConfigs(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	input := &config.ListSystemConfigsInput{
		Category: stringPtr(c.Query("category")),
	}

	output, err := h.listConfigs.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetConfig maneja GET /api/v1/admin/config/:key
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	configKey := c.Param("key")
	if configKey == "" {
		handleError(c, errors.New("INVALID_CONFIG_KEY", "config key is required", 400, nil))
		return
	}

	input := &config.GetSystemConfigInput{
		ConfigKey: configKey,
	}

	output, err := h.getConfig.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdateConfig maneja PUT /api/v1/admin/config/:key
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	configKey := c.Param("key")
	if configKey == "" {
		handleError(c, errors.New("INVALID_CONFIG_KEY", "config key is required", 400, nil))
		return
	}

	var req struct {
		ConfigValue string `json:"config_value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errors.New("INVALID_INPUT", "invalid request body", 400, err))
		return
	}

	input := &config.UpdateSystemConfigInput{
		ConfigKey:   configKey,
		ConfigValue: req.ConfigValue,
	}

	output, err := h.updateConfig.Execute(c.Request.Context(), input, adminID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}
