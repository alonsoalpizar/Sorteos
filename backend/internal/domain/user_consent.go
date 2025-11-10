package domain

import "time"

// ConsentType representa el tipo de consentimiento GDPR
type ConsentType string

const (
	ConsentTypeTermsOfService ConsentType = "terms_of_service"
	ConsentTypePrivacyPolicy  ConsentType = "privacy_policy"
	ConsentTypeMarketingEmail ConsentType = "marketing_emails"
	ConsentTypeMarketingSMS   ConsentType = "marketing_sms"
	ConsentTypeDataProcessing ConsentType = "data_processing"
)

// UserConsent representa un consentimiento del usuario
type UserConsent struct {
	ID               int64       `json:"id" gorm:"primaryKey"`
	UserID           int64       `json:"user_id" gorm:"not null;index"`
	ConsentType      ConsentType `json:"consent_type" gorm:"type:consent_type;not null"`
	ConsentVersion   string      `json:"consent_version" gorm:"not null"` // e.g., "1.0", "2.1"
	Granted          bool        `json:"granted" gorm:"not null;default:false"`
	GrantedAt        *time.Time  `json:"granted_at,omitempty"`
	RevokedAt        *time.Time  `json:"revoked_at,omitempty"`
	IPAddress        *string     `json:"-"` // No exponer IP públicamente
	UserAgent        *string     `json:"-"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`

	// Relación
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName especifica el nombre de la tabla
func (UserConsent) TableName() string {
	return "user_consents"
}

// IsGranted verifica si el consentimiento está activo
func (uc *UserConsent) IsGranted() bool {
	return uc.Granted && uc.RevokedAt == nil
}

// Grant otorga el consentimiento
func (uc *UserConsent) Grant(ipAddress, userAgent string) {
	now := time.Now()
	uc.Granted = true
	uc.GrantedAt = &now
	uc.RevokedAt = nil
	uc.IPAddress = &ipAddress
	uc.UserAgent = &userAgent
}

// Revoke revoca el consentimiento
func (uc *UserConsent) Revoke() {
	now := time.Now()
	uc.Granted = false
	uc.RevokedAt = &now
}

// UserConsentRepository define el contrato para el repositorio de consentimientos
type UserConsentRepository interface {
	// Create crea un nuevo consentimiento
	Create(consent *UserConsent) error

	// FindByUserAndType busca un consentimiento por usuario y tipo
	FindByUserAndType(userID int64, consentType ConsentType) (*UserConsent, error)

	// FindByUser busca todos los consentimientos de un usuario
	FindByUser(userID int64) ([]*UserConsent, error)

	// Update actualiza un consentimiento existente
	Update(consent *UserConsent) error

	// HasGrantedConsent verifica si el usuario ha otorgado un consentimiento específico
	HasGrantedConsent(userID int64, consentType ConsentType) (bool, error)

	// RevokeConsent revoca un consentimiento
	RevokeConsent(userID int64, consentType ConsentType) error

	// GrantConsent otorga un consentimiento
	GrantConsent(userID int64, consentType ConsentType, version, ipAddress, userAgent string) error
}
