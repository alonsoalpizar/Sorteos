package system

import (
	"context"
	"encoding/json"

	"github.com/sorteos-platform/backend/internal/repository"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// SystemSetting configuración individual del sistema
type SystemSetting struct {
	Key       string                 `json:"key"`
	Value     interface{}            `json:"value"` // Deserializado de JSON
	Category  string                 `json:"category"`
	UpdatedAt string                 `json:"updated_at"`
	UpdatedBy *int64                 `json:"updated_by,omitempty"`
}

// GetSystemSettingsInput datos de entrada
type GetSystemSettingsInput struct {
	Category *string // Filtrar por categoría (opcional)
	Key      *string // Obtener un setting específico (opcional)
}

// GetSystemSettingsOutput resultado
type GetSystemSettingsOutput struct {
	Settings       []*SystemSetting
	Categories     []string // Categorías disponibles
	TotalSettings  int
}

// GetSystemSettingsUseCase caso de uso para obtener configuración del sistema
type GetSystemSettingsUseCase struct {
	configRepo repository.SystemConfigRepository
	log        *logger.Logger
}

// NewGetSystemSettingsUseCase crea una nueva instancia
func NewGetSystemSettingsUseCase(configRepo repository.SystemConfigRepository, log *logger.Logger) *GetSystemSettingsUseCase {
	return &GetSystemSettingsUseCase{
		configRepo: configRepo,
		log:        log,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetSystemSettingsUseCase) Execute(ctx context.Context, input *GetSystemSettingsInput, adminID int64) (*GetSystemSettingsOutput, error) {
	var configs []*repository.SystemConfig
	var err error

	// Si se especifica un key, obtener solo ese
	if input.Key != nil && *input.Key != "" {
		config, err := uc.configRepo.Get(ctx, *input.Key)
		if err != nil {
			if err.Error() == "record not found" {
				return nil, errors.New("SETTING_NOT_FOUND",
					"system setting not found", 404, nil)
			}
			uc.log.Error("Error getting system setting", logger.String("key", *input.Key), logger.Error(err))
			return nil, errors.Wrap(errors.ErrDatabaseError, err)
		}
		configs = []*repository.SystemConfig{config}
	} else if input.Category != nil && *input.Category != "" {
		// Si se especifica categoría, filtrar por ella
		configs, err = uc.configRepo.GetByCategory(ctx, *input.Category)
		if err != nil {
			uc.log.Error("Error getting settings by category", logger.String("category", *input.Category), logger.Error(err))
			return nil, errors.Wrap(errors.ErrDatabaseError, err)
		}
	} else {
		// Obtener todas las configuraciones
		configs, err = uc.configRepo.GetAll(ctx)
		if err != nil {
			uc.log.Error("Error getting all system settings", logger.Error(err))
			return nil, errors.Wrap(errors.ErrDatabaseError, err)
		}
	}

	// Convertir a SystemSetting con JSON deserializado
	settings := make([]*SystemSetting, 0, len(configs))
	categoriesMap := make(map[string]bool)

	for _, config := range configs {
		// Deserializar JSON value
		var value interface{}
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			uc.log.Error("Error unmarshaling setting value",
				logger.String("key", config.Key),
				logger.Error(err))
			// Usar valor raw si falla deserialización
			value = config.Value
		}

		setting := &SystemSetting{
			Key:       config.Key,
			Value:     value,
			Category:  config.Category,
			UpdatedAt: config.UpdatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy: config.UpdatedBy,
		}

		settings = append(settings, setting)
		categoriesMap[config.Category] = true
	}

	// Convertir categoriesMap a slice
	categories := make([]string, 0, len(categoriesMap))
	for cat := range categoriesMap {
		categories = append(categories, cat)
	}

	// Log auditoría
	uc.log.Info("Admin viewed system settings",
		logger.Int64("admin_id", adminID),
		logger.Int("total_settings", len(settings)),
		logger.String("action", "admin_view_system_settings"))

	return &GetSystemSettingsOutput{
		Settings:      settings,
		Categories:    categories,
		TotalSettings: len(settings),
	}, nil
}
