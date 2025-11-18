package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UpdatePaymentProcessorInput datos de entrada
type UpdatePaymentProcessorInput struct {
	Processor string                 // stripe, paypal, etc.
	Enabled   bool                   // Habilitar/deshabilitar
	Config    map[string]interface{} // Configuración específica del processor
	Priority  *int                   // Prioridad (para cuando hay múltiples processors)
	Notes     string                 // Notas administrativas
}

// PaymentProcessorConfig configuración de un processor
type PaymentProcessorConfig struct {
	Processor  string                 `json:"processor"`
	Enabled    bool                   `json:"enabled"`
	Config     map[string]interface{} `json:"config"`
	Priority   int                    `json:"priority"`
	Notes      string                 `json:"notes"`
	UpdatedAt  time.Time              `json:"updated_at"`
	UpdatedBy  int64                  `json:"updated_by"`
}

// UpdatePaymentProcessorUseCase caso de uso para actualizar configuración de payment processors
type UpdatePaymentProcessorUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUpdatePaymentProcessorUseCase crea una nueva instancia
func NewUpdatePaymentProcessorUseCase(db *gorm.DB, log *logger.Logger) *UpdatePaymentProcessorUseCase {
	return &UpdatePaymentProcessorUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdatePaymentProcessorUseCase) Execute(ctx context.Context, input *UpdatePaymentProcessorInput, adminID int64) (*PaymentProcessorConfig, error) {
	// Validar processor
	validProcessors := []string{"stripe", "paypal", "mercadopago", "pagadito"}
	isValid := false
	for _, p := range validProcessors {
		if input.Processor == p {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("invalid processor: must be one of %v", validProcessors), 400, nil)
	}

	// Validar prioridad
	if input.Priority != nil && (*input.Priority < 1 || *input.Priority > 10) {
		return nil, errors.New("VALIDATION_FAILED", "priority must be between 1 and 10", 400, nil)
	}

	// Por ahora, guardamos la configuración en una tabla system_config
	// En el futuro esto podría ser una tabla dedicada payment_processors

	// NOTA: Asumiendo que existe una tabla system_config con estructura:
	// - key (string, PK)
	// - value (jsonb)
	// - updated_at (timestamp)
	// - updated_by (bigint)

	configKey := fmt.Sprintf("payment_processor.%s", input.Processor)
	now := time.Now()

	priority := 5 // Default priority
	if input.Priority != nil {
		priority = *input.Priority
	}

	configData := &PaymentProcessorConfig{
		Processor:  input.Processor,
		Enabled:    input.Enabled,
		Config:     input.Config,
		Priority:   priority,
		Notes:      input.Notes,
		UpdatedAt:  now,
		UpdatedBy:  adminID,
	}

	// Construir update/insert
	// Usamos Raw SQL para manejar UPSERT (ON CONFLICT)
	query := `
		INSERT INTO system_config (key, value, updated_at, updated_by)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = EXCLUDED.updated_at,
			updated_by = EXCLUDED.updated_by
	`

	if err := uc.db.Exec(query, configKey, configData, now, adminID).Error; err != nil {
		uc.log.Error("Error updating payment processor config",
			logger.String("processor", input.Processor),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría crítica
	uc.log.Error("Admin updated payment processor",
		logger.Int64("admin_id", adminID),
		logger.String("processor", input.Processor),
		logger.Bool("enabled", input.Enabled),
		logger.Int("priority", priority),
		logger.String("action", "admin_update_payment_processor"),
		logger.String("severity", "critical"))

	// Si estamos deshabilitando un processor, verificar que haya al menos uno activo
	if !input.Enabled {
		var activeCount int64
		if err := uc.db.Raw(`
			SELECT COUNT(*)
			FROM system_config
			WHERE key LIKE 'payment_processor.%'
			AND value->>'enabled' = 'true'
		`).Scan(&activeCount).Error; err != nil {
			uc.log.Warn("Could not verify active processors", logger.Error(err))
		} else if activeCount == 0 {
			uc.log.Error("WARNING: All payment processors are now disabled",
				logger.Int64("admin_id", adminID),
				logger.String("severity", "critical"))
		}
	}

	return configData, nil
}

// GetPaymentProcessorConfig obtiene la configuración de un processor
func (uc *UpdatePaymentProcessorUseCase) GetPaymentProcessorConfig(processor string) (*PaymentProcessorConfig, error) {
	configKey := fmt.Sprintf("payment_processor.%s", processor)

	var config PaymentProcessorConfig
	if err := uc.db.Raw(`
		SELECT value
		FROM system_config
		WHERE key = ?
	`, configKey).Scan(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("CONFIG_NOT_FOUND", "payment processor config not found", 404, nil)
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &config, nil
}

// ListPaymentProcessors lista todos los processors configurados
func (uc *UpdatePaymentProcessorUseCase) ListPaymentProcessors() ([]*PaymentProcessorConfig, error) {
	var configs []*PaymentProcessorConfig

	if err := uc.db.Raw(`
		SELECT value
		FROM system_config
		WHERE key LIKE 'payment_processor.%'
		ORDER BY value->>'priority' ASC
	`).Scan(&configs).Error; err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return configs, nil
}
