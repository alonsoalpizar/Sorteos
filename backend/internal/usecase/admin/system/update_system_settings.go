package system

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sorteos-platform/backend/internal/repository"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// UpdateSystemSettingsInput datos de entrada
type UpdateSystemSettingsInput struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"` // Se serializará a JSON
	Category string      `json:"category"`
}

// UpdateSystemSettingsOutput resultado
type UpdateSystemSettingsOutput struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	Category  string      `json:"category"`
	UpdatedAt string      `json:"updated_at"`
	Success   bool        `json:"success"`
}

// UpdateSystemSettingsUseCase caso de uso para actualizar configuración del sistema
type UpdateSystemSettingsUseCase struct {
	configRepo repository.SystemConfigRepository
	log        *logger.Logger
}

// NewUpdateSystemSettingsUseCase crea una nueva instancia
func NewUpdateSystemSettingsUseCase(configRepo repository.SystemConfigRepository, log *logger.Logger) *UpdateSystemSettingsUseCase {
	return &UpdateSystemSettingsUseCase{
		configRepo: configRepo,
		log:        log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateSystemSettingsUseCase) Execute(ctx context.Context, input *UpdateSystemSettingsInput, adminID int64) (*UpdateSystemSettingsOutput, error) {
	// Validar key
	if input.Key == "" {
		return nil, errors.New("VALIDATION_FAILED", "key is required", 400, nil)
	}

	// Validar category
	if input.Category == "" {
		return nil, errors.New("VALIDATION_FAILED", "category is required", 400, nil)
	}

	// Validar value (no puede ser nil)
	if input.Value == nil {
		return nil, errors.New("VALIDATION_FAILED", "value cannot be null", 400, nil)
	}

	// Serializar value a JSON
	valueJSON, err := json.Marshal(input.Value)
	if err != nil {
		uc.log.Error("Error marshaling setting value",
			logger.String("key", input.Key),
			logger.Error(err))
		return nil, errors.New("VALIDATION_FAILED", "value must be JSON-serializable", 400, err)
	}

	// Validaciones específicas por key (ejemplos)
	if err := uc.validateSettingValue(input.Key, input.Value); err != nil {
		return nil, err
	}

	// Actualizar o crear configuración
	if err := uc.configRepo.Set(ctx, input.Key, string(valueJSON), input.Category, adminID); err != nil {
		uc.log.Error("Error updating system setting",
			logger.String("key", input.Key),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener la configuración actualizada para confirmar
	config, err := uc.configRepo.Get(ctx, input.Key)
	if err != nil {
		uc.log.Error("Error getting updated setting", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría crítica (cambio de configuración es crítico)
	uc.log.Error("Admin updated system setting",
		logger.Int64("admin_id", adminID),
		logger.String("key", input.Key),
		logger.String("category", input.Category),
		logger.String("action", "admin_update_system_setting"),
		logger.String("severity", "critical"))

	return &UpdateSystemSettingsOutput{
		Key:       config.Key,
		Value:     input.Value,
		Category:  config.Category,
		UpdatedAt: config.UpdatedAt.Format("2006-01-02 15:04:05"),
		Success:   true,
	}, nil
}

// validateSettingValue valida valores específicos según el key
func (uc *UpdateSystemSettingsUseCase) validateSettingValue(key string, value interface{}) error {
	switch key {
	case "platform_fee_percent":
		// Debe ser un número entre 0 y 100
		if floatVal, ok := value.(float64); ok {
			if floatVal < 0 || floatVal > 100 {
				return errors.New("VALIDATION_FAILED",
					"platform_fee_percent must be between 0 and 100", 400, nil)
			}
		} else {
			return errors.New("VALIDATION_FAILED",
				"platform_fee_percent must be a number", 400, nil)
		}

	case "min_raffle_price":
		// Debe ser un número positivo
		if floatVal, ok := value.(float64); ok {
			if floatVal <= 0 {
				return errors.New("VALIDATION_FAILED",
					"min_raffle_price must be greater than 0", 400, nil)
			}
		} else {
			return errors.New("VALIDATION_FAILED",
				"min_raffle_price must be a number", 400, nil)
		}

	case "max_raffle_numbers":
		// Debe ser un entero positivo
		if floatVal, ok := value.(float64); ok {
			if floatVal <= 0 || floatVal > 1000000 {
				return errors.New("VALIDATION_FAILED",
					"max_raffle_numbers must be between 1 and 1000000", 400, nil)
			}
		} else {
			return errors.New("VALIDATION_FAILED",
				"max_raffle_numbers must be a number", 400, nil)
		}

	case "email_provider":
		// Debe ser smtp, sendgrid, mailgun, ses
		if strVal, ok := value.(string); ok {
			validProviders := map[string]bool{
				"smtp":     true,
				"sendgrid": true,
				"mailgun":  true,
				"ses":      true,
			}
			if !validProviders[strVal] {
				return errors.New("VALIDATION_FAILED",
					fmt.Sprintf("email_provider must be one of: smtp, sendgrid, mailgun, ses"), 400, nil)
			}
		} else {
			return errors.New("VALIDATION_FAILED",
				"email_provider must be a string", 400, nil)
		}

	case "maintenance_mode":
		// Debe ser booleano
		if _, ok := value.(bool); !ok {
			return errors.New("VALIDATION_FAILED",
				"maintenance_mode must be a boolean", 400, nil)
		}

	// Agregar más validaciones según necesidad
	}

	return nil
}
