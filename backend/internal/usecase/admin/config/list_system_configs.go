package config

import (
	"context"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ListSystemConfigsInput datos de entrada
type ListSystemConfigsInput struct {
	Category *string `json:"category,omitempty"` // Filtrar por categoría (email, payment, general, etc.)
}

// ConfigListItem item de configuración en la lista
type ConfigListItem struct {
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	UpdatedAt   string `json:"updated_at"`
}

// ListSystemConfigsOutput resultado
type ListSystemConfigsOutput struct {
	Configs    []*ConfigListItem `json:"configs"`
	TotalCount int               `json:"total_count"`
}

// ListSystemConfigsUseCase caso de uso para listar configuraciones del sistema
type ListSystemConfigsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListSystemConfigsUseCase crea una nueva instancia
func NewListSystemConfigsUseCase(db *gorm.DB, log *logger.Logger) *ListSystemConfigsUseCase {
	return &ListSystemConfigsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListSystemConfigsUseCase) Execute(ctx context.Context, input *ListSystemConfigsInput, adminID int64) (*ListSystemConfigsOutput, error) {
	// Construir query base
	query := uc.db.WithContext(ctx).
		Table("system_config")

	// Aplicar filtros
	if input.Category != nil && *input.Category != "" {
		query = query.Where("category = ?", *input.Category)
	}

	// Ejecutar query
	var configs []struct {
		ConfigKey   string
		ConfigValue string
		Category    *string
		Description *string
		UpdatedAt   string
	}

	result := query.
		Select("config_key, config_value, category, description, updated_at").
		Order("category ASC, config_key ASC").
		Find(&configs)

	if result.Error != nil {
		uc.log.Error("Error fetching configs", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log acceso a configuración
	uc.log.Info("Admin listed system configs",
		logger.Int64("admin_id", adminID),
		logger.Int("count", len(configs)),
		logger.String("action", "admin_list_configs"))

	// Construir output
	configItems := make([]*ConfigListItem, 0, len(configs))
	for _, cfg := range configs {
		category := ""
		if cfg.Category != nil {
			category = *cfg.Category
		}
		description := ""
		if cfg.Description != nil {
			description = *cfg.Description
		}

		configItems = append(configItems, &ConfigListItem{
			ConfigKey:   cfg.ConfigKey,
			ConfigValue: cfg.ConfigValue,
			Category:    category,
			Description: description,
			UpdatedAt:   cfg.UpdatedAt,
		})
	}

	return &ListSystemConfigsOutput{
		Configs:    configItems,
		TotalCount: len(configItems),
	}, nil
}
