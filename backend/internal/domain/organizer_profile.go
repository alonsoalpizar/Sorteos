package domain

import (
	"fmt"
	"time"
)

// PayoutSchedule representa la frecuencia de pago a organizadores
type PayoutSchedule string

const (
	PayoutScheduleManual  PayoutSchedule = "manual"
	PayoutScheduleWeekly  PayoutSchedule = "weekly"
	PayoutScheduleMonthly PayoutSchedule = "monthly"
)

// BankAccountType representa el tipo de cuenta bancaria
type BankAccountType string

const (
	BankAccountTypeChecking BankAccountType = "checking"
	BankAccountTypeSavings  BankAccountType = "savings"
)

// OrganizerProfile representa el perfil extendido de un organizador de rifas
type OrganizerProfile struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	UserID int64 `json:"user_id" gorm:"not null;uniqueIndex"` // FK a users

	// Business Info
	BusinessName *string `json:"business_name,omitempty"`
	TaxID        *string `json:"tax_id,omitempty"` // RUC o Tax ID de la empresa

	// Bank Info (datos sensibles - deben ser encriptados en app layer)
	BankName          *string          `json:"bank_name,omitempty"`
	BankAccountNumber *string          `json:"-"` // Never serialize - Encrypted
	BankAccountType   *BankAccountType `json:"bank_account_type,omitempty"`
	BankAccountHolder *string          `json:"bank_account_holder,omitempty"`

	// Payout Configuration
	PayoutSchedule     PayoutSchedule `json:"payout_schedule" gorm:"type:varchar(20);default:'manual'"`
	CommissionOverride *float64       `json:"commission_override,omitempty" gorm:"type:decimal(5,2)"` // NULL = use global default

	// Financial Tracking
	TotalPayouts  float64 `json:"total_payouts" gorm:"type:decimal(12,2);default:0.00"`
	PendingPayout float64 `json:"pending_payout" gorm:"type:decimal(12,2);default:0.00"`

	// Verification
	Verified     bool       `json:"verified" gorm:"default:false"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	VerifiedBy   *int64     `json:"verified_by,omitempty"` // Admin user ID que verificó
	VerifiedByID *int64     `json:"-" gorm:"column:verified_by"`

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relación (no se mapea a columna, solo para carga eager)
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName especifica el nombre de la tabla
func (OrganizerProfile) TableName() string {
	return "organizer_profiles"
}

// Validate valida los campos del OrganizerProfile
func (op *OrganizerProfile) Validate() error {
	// UserID es requerido
	if op.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	// Validar business name si está presente
	if op.BusinessName != nil && len(*op.BusinessName) > 255 {
		return fmt.Errorf("business name is too long (max 255 characters)")
	}

	// Validar commission override si está presente
	if op.CommissionOverride != nil {
		if *op.CommissionOverride < 0 || *op.CommissionOverride > 50 {
			return fmt.Errorf("commission override must be between 0 and 50 percent")
		}
	}

	// Validar bank account type si está presente
	if op.BankAccountType != nil {
		if *op.BankAccountType != BankAccountTypeChecking && *op.BankAccountType != BankAccountTypeSavings {
			return fmt.Errorf("invalid bank account type: %s (valid: checking, savings)", *op.BankAccountType)
		}
	}

	// Validar payout schedule
	validSchedules := []PayoutSchedule{PayoutScheduleManual, PayoutScheduleWeekly, PayoutScheduleMonthly}
	validSchedule := false
	for _, vs := range validSchedules {
		if op.PayoutSchedule == vs {
			validSchedule = true
			break
		}
	}
	if !validSchedule {
		return fmt.Errorf("invalid payout schedule: %s (valid: manual, weekly, monthly)", op.PayoutSchedule)
	}

	return nil
}

// HasCustomCommission verifica si el organizador tiene comisión personalizada
func (op *OrganizerProfile) HasCustomCommission() bool {
	return op.CommissionOverride != nil
}

// GetEffectiveCommission obtiene la comisión efectiva (custom o global default)
func (op *OrganizerProfile) GetEffectiveCommission(globalDefault float64) float64 {
	if op.HasCustomCommission() {
		return *op.CommissionOverride
	}
	return globalDefault
}

// HasBankInfo verifica si el organizador tiene información bancaria completa
func (op *OrganizerProfile) HasBankInfo() bool {
	return op.BankName != nil && *op.BankName != "" &&
		op.BankAccountNumber != nil && *op.BankAccountNumber != "" &&
		op.BankAccountHolder != nil && *op.BankAccountHolder != ""
}

// CanReceivePayouts verifica si el organizador puede recibir pagos
func (op *OrganizerProfile) CanReceivePayouts() bool {
	return op.Verified && op.HasBankInfo()
}

// MaskBankInfo enmascara la información bancaria para logging/display
func (op *OrganizerProfile) MaskBankInfo() *OrganizerProfile {
	masked := *op
	if masked.BankAccountNumber != nil && *masked.BankAccountNumber != "" {
		accountLen := len(*masked.BankAccountNumber)
		if accountLen > 4 {
			maskedValue := "****" + (*masked.BankAccountNumber)[accountLen-4:]
			masked.BankAccountNumber = &maskedValue
		}
	}
	return &masked
}

// OrganizerRevenue representa el desglose de ingresos de un organizador
type OrganizerRevenue struct {
	OrganizerID     int64   `json:"organizer_id"`
	TotalRaffles    int     `json:"total_raffles"`
	CompletedRaffles int    `json:"completed_raffles"`
	TotalRevenue    float64 `json:"total_revenue"`
	PlatformFees    float64 `json:"platform_fees"`
	NetRevenue      float64 `json:"net_revenue"`
	TotalPayouts    float64 `json:"total_payouts"`
	PendingPayout   float64 `json:"pending_payout"`
}

// OrganizerProfileRepository define los métodos de acceso a datos
type OrganizerProfileRepository interface {
	// Create crea un nuevo perfil de organizador
	Create(profile *OrganizerProfile) error

	// GetByUserID obtiene un perfil por user_id
	GetByUserID(userID int64) (*OrganizerProfile, error)

	// GetByID obtiene un perfil por ID
	GetByID(id int64) (*OrganizerProfile, error)

	// List obtiene perfiles con filtros y paginación
	List(filters map[string]interface{}, offset, limit int) ([]*OrganizerProfile, int64, error)

	// Update actualiza un perfil de organizador
	Update(profile *OrganizerProfile) error

	// UpdateCommission actualiza solo la comisión
	UpdateCommission(userID int64, commission *float64) error

	// UpdateFinancials actualiza totales financieros
	UpdateFinancials(userID int64, totalPayouts, pendingPayout float64) error

	// GetRevenue obtiene el desglose de ingresos de un organizador
	GetRevenue(userID int64, dateFrom, dateTo *time.Time) (*OrganizerRevenue, error)

	// Verify marca un organizador como verificado
	Verify(userID int64, verifiedBy int64) error
}
