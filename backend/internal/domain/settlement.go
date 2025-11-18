package domain

import (
	"fmt"
	"time"
)

// SettlementStatus representa el estado de una liquidación
type SettlementStatus string

const (
	SettlementStatusPending    SettlementStatus = "pending"    // Pendiente de aprobación
	SettlementStatusApproved   SettlementStatus = "approved"   // Aprobado, listo para pagar
	SettlementStatusPaid       SettlementStatus = "paid"       // Pagado al organizador
	SettlementStatusRejected   SettlementStatus = "rejected"   // Rechazado por admin
)

// Settlement representa una liquidación/pago a un organizador
type Settlement struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	UUID string `json:"uuid" gorm:"type:uuid;unique;not null;default:uuid_generate_v4()"`

	// References
	RaffleID     int64 `json:"raffle_id" gorm:"not null;uniqueIndex"` // FK a raffles - cada rifa solo puede tener un settlement
	OrganizerID  int64 `json:"organizer_id" gorm:"not null;index"`     // FK a users

	// Amounts (calculados automáticamente)
	GrossRevenue           float64 `json:"gross_revenue" gorm:"type:decimal(12,2);not null"`            // Total vendido
	PlatformFee            float64 `json:"platform_fee" gorm:"type:decimal(12,2);not null"`             // Comisión de plataforma
	PlatformFeePercentage  float64 `json:"platform_fee_percentage" gorm:"type:decimal(5,2);not null"`   // % aplicado
	NetPayout              float64 `json:"net_payout" gorm:"type:decimal(12,2);not null"`               // A pagar al organizador

	// Status
	Status SettlementStatus `json:"status" gorm:"type:settlement_status;default:'pending'"`

	// Payment Info
	PaymentMethod    *string `json:"payment_method,omitempty"`    // 'bank_transfer', 'paypal', 'sinpe', etc.
	PaymentReference *string `json:"payment_reference,omitempty"` // Número de transferencia, PayPal transaction ID, etc.

	// Approval
	ApprovedBy   *int64     `json:"approved_by,omitempty"` // Admin user ID que aprobó
	ApprovedByID *int64     `json:"-" gorm:"column:approved_by"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`

	// Payment
	PaidAt *time.Time `json:"paid_at,omitempty"`

	// Notes
	Notes *string `json:"notes,omitempty"` // Notas de admin (ej: razón de rechazo)

	// Audit
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relaciones (no se mapean a columnas)
	Raffle    *Raffle `json:"raffle,omitempty" gorm:"foreignKey:RaffleID"`
	Organizer *User   `json:"organizer,omitempty" gorm:"foreignKey:OrganizerID"`
}

// TableName especifica el nombre de la tabla
func (Settlement) TableName() string {
	return "settlements"
}

// Validate valida los campos del Settlement
func (s *Settlement) Validate() error {
	// RaffleID es requerido
	if s.RaffleID == 0 {
		return fmt.Errorf("raffle_id is required")
	}

	// OrganizerID es requerido
	if s.OrganizerID == 0 {
		return fmt.Errorf("organizer_id is required")
	}

	// Amounts deben ser positivos
	if s.GrossRevenue < 0 {
		return fmt.Errorf("gross_revenue must be non-negative")
	}

	if s.PlatformFee < 0 {
		return fmt.Errorf("platform_fee must be non-negative")
	}

	if s.NetPayout < 0 {
		return fmt.Errorf("net_payout must be non-negative")
	}

	// Validar que net_payout = gross_revenue - platform_fee
	expectedNetPayout := s.GrossRevenue - s.PlatformFee
	// Permitir pequeñas diferencias por redondeo (0.01)
	if abs(s.NetPayout-expectedNetPayout) > 0.01 {
		return fmt.Errorf("net_payout must equal gross_revenue - platform_fee (%.2f != %.2f - %.2f)",
			s.NetPayout, s.GrossRevenue, s.PlatformFee)
	}

	// Validar platform fee percentage
	if s.PlatformFeePercentage < 0 || s.PlatformFeePercentage > 50 {
		return fmt.Errorf("platform_fee_percentage must be between 0 and 50")
	}

	// Validar status
	validStatuses := []SettlementStatus{
		SettlementStatusPending,
		SettlementStatusApproved,
		SettlementStatusPaid,
		SettlementStatusRejected,
	}
	validStatus := false
	for _, vs := range validStatuses {
		if s.Status == vs {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return fmt.Errorf("invalid settlement status: %s", s.Status)
	}

	return nil
}

// IsPending verifica si el settlement está pendiente
func (s *Settlement) IsPending() bool {
	return s.Status == SettlementStatusPending
}

// IsApproved verifica si el settlement está aprobado
func (s *Settlement) IsApproved() bool {
	return s.Status == SettlementStatusApproved
}

// IsPaid verifica si el settlement está pagado
func (s *Settlement) IsPaid() bool {
	return s.Status == SettlementStatusPaid
}

// IsRejected verifica si el settlement está rechazado
func (s *Settlement) IsRejected() bool {
	return s.Status == SettlementStatusRejected
}

// CanApprove verifica si el settlement puede ser aprobado
func (s *Settlement) CanApprove() bool {
	return s.IsPending()
}

// CanReject verifica si el settlement puede ser rechazado
func (s *Settlement) CanReject() bool {
	return s.IsPending()
}

// CanMarkPaid verifica si el settlement puede ser marcado como pagado
func (s *Settlement) CanMarkPaid() bool {
	return s.IsApproved()
}

// CalculateFromRaffle calcula los montos del settlement desde una rifa
func (s *Settlement) CalculateFromRaffle(raffle *Raffle, commissionPercentage float64) {
	// Convert decimal.Decimal to float64
	grossRevenue, _ := raffle.TotalRevenue.Float64()

	s.GrossRevenue = grossRevenue
	s.PlatformFeePercentage = commissionPercentage
	s.PlatformFee = s.GrossRevenue * (commissionPercentage / 100.0)
	s.NetPayout = s.GrossRevenue - s.PlatformFee
}

// abs retorna el valor absoluto de un float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// SettlementRepository define los métodos de acceso a datos
type SettlementRepository interface {
	// Create crea un nuevo settlement
	Create(settlement *Settlement) error

	// GetByID obtiene un settlement por ID
	GetByID(id int64) (*Settlement, error)

	// GetByUUID obtiene un settlement por UUID
	GetByUUID(uuid string) (*Settlement, error)

	// GetByRaffleID obtiene un settlement por raffle_id
	GetByRaffleID(raffleID int64) (*Settlement, error)

	// List obtiene settlements con filtros y paginación
	List(filters map[string]interface{}, offset, limit int) ([]*Settlement, int64, error)

	// UpdateStatus actualiza solo el status
	UpdateStatus(id int64, status SettlementStatus) error

	// Approve aprueba un settlement
	Approve(id int64, adminID int64) error

	// Reject rechaza un settlement
	Reject(id int64, reason string) error

	// MarkPaid marca un settlement como pagado
	MarkPaid(id int64, paymentMethod, paymentReference string) error

	// GetPendingByOrganizer obtiene settlements pendientes de un organizador
	GetPendingByOrganizer(organizerID int64) ([]*Settlement, error)

	// GetTotalsByStatus obtiene totales agrupados por status
	GetTotalsByStatus() (map[SettlementStatus]float64, error)
}
