package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// ProcessorProvider representa el tipo de procesador de pagos
type ProcessorProvider string

const (
	ProcessorProviderStripe ProcessorProvider = "stripe"
	ProcessorProviderPayPal ProcessorProvider = "paypal"
	ProcessorProviderCredix ProcessorProvider = "credix"
)

// ValidProviders es la lista de proveedores válidos
var ValidProviders = []ProcessorProvider{
	ProcessorProviderStripe,
	ProcessorProviderPayPal,
	ProcessorProviderCredix,
}

// PaymentProcessor representa la configuración de un procesador de pagos
type PaymentProcessor struct {
	ID int64 `json:"id" gorm:"primaryKey"`

	// Provider Info
	Provider ProcessorProvider `json:"provider" gorm:"type:varchar(50);not null"`
	Name     string            `json:"name" gorm:"not null"` // e.g., "Stripe Production", "PayPal Sandbox"

	// Status
	IsActive  bool `json:"is_active" gorm:"default:true"`
	IsSandbox bool `json:"is_sandbox" gorm:"default:false"`

	// Credentials (sensitive data - should be encrypted in app layer)
	// Estos campos se almacenan como texto pero deben ser encriptados/desencriptados
	// en la capa de aplicación antes de guardar/leer
	ClientID      *string `json:"client_id,omitempty"` // PayPal Client ID, Stripe Publishable Key
	SecretKey     *string `json:"-"`                   // Never serialize - Secret key (encrypted)
	WebhookSecret *string `json:"-"`                   // Never serialize - Webhook verification secret (encrypted)

	// Configuration
	Currency string          `json:"currency" gorm:"type:char(3);default:'CRC'"` // ISO 4217
	Config   json.RawMessage `json:"config,omitempty" gorm:"type:jsonb"`         // Additional provider-specific config

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (PaymentProcessor) TableName() string {
	return "payment_processors"
}

// Validate valida los campos del PaymentProcessor
func (pp *PaymentProcessor) Validate() error {
	// Provider es requerido
	if pp.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	// Validar que el provider sea válido
	validProvider := false
	for _, vp := range ValidProviders {
		if pp.Provider == vp {
			validProvider = true
			break
		}
	}
	if !validProvider {
		return fmt.Errorf("invalid provider: %s (valid: stripe, paypal, credix)", pp.Provider)
	}

	// Name es requerido
	if pp.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(pp.Name) > 255 {
		return fmt.Errorf("name is too long (max 255 characters)")
	}

	// Currency debe ser código de 3 caracteres
	if len(pp.Currency) != 3 {
		return fmt.Errorf("currency code must be exactly 3 characters (ISO 4217)")
	}

	// Validar que Config sea JSON válido si está presente
	if pp.Config != nil && len(pp.Config) > 0 {
		var configTest map[string]interface{}
		if err := json.Unmarshal(pp.Config, &configTest); err != nil {
			return fmt.Errorf("invalid config JSON: %w", err)
		}
	}

	return nil
}

// IsStripe verifica si el procesador es Stripe
func (pp *PaymentProcessor) IsStripe() bool {
	return pp.Provider == ProcessorProviderStripe
}

// IsPayPal verifica si el procesador es PayPal
func (pp *PaymentProcessor) IsPayPal() bool {
	return pp.Provider == ProcessorProviderPayPal
}

// MaskSecrets enmascara los secretos para logging seguro
func (pp *PaymentProcessor) MaskSecrets() *PaymentProcessor {
	masked := *pp
	if masked.SecretKey != nil && *masked.SecretKey != "" {
		maskedValue := "***" + (*masked.SecretKey)[len(*masked.SecretKey)-4:]
		masked.SecretKey = &maskedValue
	}
	if masked.WebhookSecret != nil && *masked.WebhookSecret != "" {
		maskedValue := "***" + (*masked.WebhookSecret)[len(*masked.WebhookSecret)-4:]
		masked.WebhookSecret = &maskedValue
	}
	return &masked
}

// PaymentProcessorRepository define los métodos de acceso a datos
type PaymentProcessorRepository interface {
	// List obtiene todos los procesadores de pago
	List() ([]*PaymentProcessor, error)

	// GetByID obtiene un procesador por ID
	GetByID(id int64) (*PaymentProcessor, error)

	// GetByProvider obtiene un procesador por tipo de proveedor
	GetByProvider(provider ProcessorProvider) (*PaymentProcessor, error)

	// GetActive obtiene el procesador activo
	GetActive() (*PaymentProcessor, error)

	// Update actualiza un procesador de pago
	Update(processor *PaymentProcessor) error

	// ToggleActive activa/desactiva un procesador
	ToggleActive(id int64, active bool) error
}
