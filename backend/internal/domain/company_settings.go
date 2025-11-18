package domain

import (
	"fmt"
	"time"
)

// CompanySettings representa la configuración global de la empresa.
// Esta tabla sigue el patrón singleton (solo puede haber un registro).
type CompanySettings struct {
	ID int64 `json:"id" gorm:"primaryKey"`

	// Company Info
	CompanyName string  `json:"company_name" gorm:"not null;default:'Sorteos.club'"`
	TaxID       *string `json:"tax_id,omitempty"` // RUC o Tax ID

	// Address
	AddressLine1 *string `json:"address_line1,omitempty"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Country      string  `json:"country" gorm:"type:char(2);default:'CR'"` // ISO 3166-1 alpha-2

	// Contact
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`
	Website      string  `json:"website" gorm:"default:'https://sorteos.club'"`
	SupportEmail string  `json:"support_email" gorm:"default:'soporte@sorteos.club'"`

	// Branding
	LogoURL *string `json:"logo_url,omitempty"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (CompanySettings) TableName() string {
	return "company_settings"
}

// Validate valida los campos de CompanySettings
func (cs *CompanySettings) Validate() error {
	// Company name es requerido
	if cs.CompanyName == "" {
		return fmt.Errorf("company name is required")
	}

	if len(cs.CompanyName) > 255 {
		return fmt.Errorf("company name is too long (max 255 characters)")
	}

	// Validar email si está presente
	if cs.Email != nil && *cs.Email != "" {
		if err := ValidateEmail(*cs.Email); err != nil {
			return fmt.Errorf("invalid company email: %w", err)
		}
	}

	// Validar support email
	if cs.SupportEmail != "" {
		if err := ValidateEmail(cs.SupportEmail); err != nil {
			return fmt.Errorf("invalid support email: %w", err)
		}
	}

	// Validar country code (debe ser exactamente 2 caracteres)
	if len(cs.Country) != 2 {
		return fmt.Errorf("country code must be exactly 2 characters (ISO 3166-1 alpha-2)")
	}

	return nil
}

// CompanySettingsRepository define los métodos de acceso a datos para CompanySettings
type CompanySettingsRepository interface {
	// Get obtiene la configuración de la empresa (singleton)
	Get() (*CompanySettings, error)

	// Update actualiza la configuración de la empresa
	Update(settings *CompanySettings) error
}
