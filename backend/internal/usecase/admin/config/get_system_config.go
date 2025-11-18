package config

import (
	"context"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// GetSystemConfigInput datos de entrada
type GetSystemConfigInput struct {
	ConfigKey string `json:"config_key"`
}

// GetSystemConfigOutput resultado
type GetSystemConfigOutput struct {
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	UpdatedAt   string `json:"updated_at"`
}

// GetSystemConfigUseCase caso de uso para obtener configuración del sistema
type GetSystemConfigUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewGetSystemConfigUseCase crea una nueva instancia
func NewGetSystemConfigUseCase(db *gorm.DB, log *logger.Logger) *GetSystemConfigUseCase {
	return &GetSystemConfigUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetSystemConfigUseCase) Execute(ctx context.Context, input *GetSystemConfigInput, adminID int64) (*GetSystemConfigOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Buscar configuración
	var config struct {
		ConfigKey   string
		ConfigValue string
		Category    *string
		Description *string
		UpdatedAt   string
	}

	result := uc.db.WithContext(ctx).
		Table("system_parameters").
		Select("key as config_key, value as config_value, category, description, updated_at").
		Where("key = ?", input.ConfigKey).
		First(&config)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("CONFIG_NOT_FOUND", "configuration key not found", 404, nil)
		}
		uc.log.Error("Error finding config", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log acceso a configuración sensible
	uc.log.Info("Admin accessed system config",
		logger.Int64("admin_id", adminID),
		logger.String("config_key", input.ConfigKey),
		logger.String("action", "admin_get_config"))

	// Construir output
	category := ""
	if config.Category != nil {
		category = *config.Category
	}
	description := ""
	if config.Description != nil {
		description = *config.Description
	}

	return &GetSystemConfigOutput{
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		Category:    category,
		Description: description,
		UpdatedAt:   config.UpdatedAt,
	}, nil
}

// validateInput valida los datos de entrada
func (uc *GetSystemConfigUseCase) validateInput(input *GetSystemConfigInput) error {
	if input.ConfigKey == "" {
		return errors.New("VALIDATION_FAILED", "config_key is required", 400, nil)
	}

	return nil
}
