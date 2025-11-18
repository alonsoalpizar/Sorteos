package config

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UpdateSystemConfigInput datos de entrada
type UpdateSystemConfigInput struct {
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
}

// UpdateSystemConfigOutput resultado
type UpdateSystemConfigOutput struct {
	ConfigKey      string `json:"config_key"`
	ConfigValue    string `json:"config_value"`
	PreviousValue  string `json:"previous_value"`
	UpdatedAt      string `json:"updated_at"`
	Message        string `json:"message"`
}

// UpdateSystemConfigUseCase caso de uso para actualizar configuración del sistema
type UpdateSystemConfigUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUpdateSystemConfigUseCase crea una nueva instancia
func NewUpdateSystemConfigUseCase(db *gorm.DB, log *logger.Logger) *UpdateSystemConfigUseCase {
	return &UpdateSystemConfigUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateSystemConfigUseCase) Execute(ctx context.Context, input *UpdateSystemConfigInput, adminID int64) (*UpdateSystemConfigOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Buscar configuración actual
	var currentConfig struct {
		ConfigKey   string
		ConfigValue string
	}

	result := uc.db.WithContext(ctx).
		Table("system_config").
		Select("config_key, config_value").
		Where("config_key = ?", input.ConfigKey).
		First(&currentConfig)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("CONFIG_NOT_FOUND", "configuration key not found", 404, nil)
		}
		uc.log.Error("Error finding config", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	previousValue := currentConfig.ConfigValue

	// Validar que el nuevo valor sea diferente
	if input.ConfigValue == previousValue {
		return nil, errors.New("VALIDATION_FAILED", "new value must be different from current value", 400, nil)
	}

	// Actualizar configuración
	now := time.Now()
	result = uc.db.WithContext(ctx).
		Table("system_config").
		Where("config_key = ?", input.ConfigKey).
		Updates(map[string]interface{}{
			"config_value": input.ConfigValue,
			"updated_at":   now,
		})

	if result.Error != nil {
		uc.log.Error("Error updating config", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log auditoría crítica (cambios en configuración son sensibles)
	uc.log.Error("Admin updated system config",
		logger.Int64("admin_id", adminID),
		logger.String("config_key", input.ConfigKey),
		logger.String("previous_value", previousValue),
		logger.String("new_value", input.ConfigValue),
		logger.String("action", "admin_update_config"),
		logger.String("severity", "critical"))

	return &UpdateSystemConfigOutput{
		ConfigKey:     input.ConfigKey,
		ConfigValue:   input.ConfigValue,
		PreviousValue: previousValue,
		UpdatedAt:     now.Format(time.RFC3339),
		Message:       "System configuration updated successfully",
	}, nil
}

// validateInput valida los datos de entrada
func (uc *UpdateSystemConfigUseCase) validateInput(input *UpdateSystemConfigInput) error {
	if input.ConfigKey == "" {
		return errors.New("VALIDATION_FAILED", "config_key is required", 400, nil)
	}

	if input.ConfigValue == "" {
		return errors.New("VALIDATION_FAILED", "config_value is required", 400, nil)
	}

	// Validar longitud
	if len(input.ConfigValue) > 1000 {
		return errors.New("VALIDATION_FAILED", "config_value must be 1000 characters or less", 400, nil)
	}

	return nil
}
