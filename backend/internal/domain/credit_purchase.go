package domain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// CreditPurchaseStatus estados de una compra de créditos
type CreditPurchaseStatus string

const (
	CreditPurchaseStatusPending    CreditPurchaseStatus = "pending"     // Iniciado
	CreditPurchaseStatusProcessing CreditPurchaseStatus = "processing"  // En Pagadito
	CreditPurchaseStatusCompleted  CreditPurchaseStatus = "completed"   // Exitoso
	CreditPurchaseStatusFailed     CreditPurchaseStatus = "failed"      // Fallido
	CreditPurchaseStatusExpired    CreditPurchaseStatus = "expired"     // Expirado
)

// PagaditoStatus estados posibles de Pagadito
type PagaditoStatus string

const (
	PagaditoStatusCompleted  PagaditoStatus = "COMPLETED"   // Pago exitoso
	PagaditoStatusRegistered PagaditoStatus = "REGISTERED"  // Cancelado por usuario
	PagaditoStatusVerifying  PagaditoStatus = "VERIFYING"   // En verificación admin
	PagaditoStatusRevoked    PagaditoStatus = "REVOKED"     // Denegado
	PagaditoStatusFailed     PagaditoStatus = "FAILED"      // Fallido
)

// CreditPurchase representa una compra de créditos vía Pagadito
type CreditPurchase struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	UUID   string `json:"uuid" gorm:"type:uuid;unique;not null;default:uuid_generate_v4()"`
	UserID int64  `json:"user_id" gorm:"index;not null"`
	WalletID int64  `json:"wallet_id" gorm:"index;not null"`

	// Montos
	DesiredCredit decimal.Decimal `json:"desired_credit" gorm:"type:decimal(12,2);not null"`
	ChargeAmount  decimal.Decimal `json:"charge_amount" gorm:"type:decimal(12,2);not null"`
	Currency      string          `json:"currency" gorm:"type:varchar(3);default:'CRC';not null"`

	// Desglose de comisiones
	FixedFee      decimal.Decimal `json:"fixed_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	ProcessorFee  decimal.Decimal `json:"processor_fee" gorm:"type:decimal(12,2);not null;default:0.00"`
	PlatformFee   decimal.Decimal `json:"platform_fee" gorm:"type:decimal(12,2);not null;default:0.00"`

	// Integración Pagadito
	ERN               string  `json:"ern" gorm:"uniqueIndex;not null"`
	PagaditoToken     *string `json:"pagadito_token,omitempty" gorm:"index"`
	PagaditoReference *string `json:"pagadito_reference,omitempty"`
	PagaditoStatus    *string `json:"pagadito_status,omitempty"`

	// Estado
	Status CreditPurchaseStatus `json:"status" gorm:"type:credit_purchase_status;default:'pending';not null;index"`

	// Idempotencia
	IdempotencyKey string `json:"idempotency_key" gorm:"uniqueIndex;not null"`

	// Metadata
	Metadata datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb"`
	ErrorMessage *string `json:"error_message,omitempty"`

	// Transacción de billetera relacionada
	WalletTransactionID *int64 `json:"wallet_transaction_id,omitempty" gorm:"index"`

	// Auditoría
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"index"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	FailedAt    *time.Time `json:"failed_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (CreditPurchase) TableName() string {
	return "credit_purchases"
}

// IsPending verifica si está pendiente
func (cp *CreditPurchase) IsPending() bool {
	return cp.Status == CreditPurchaseStatusPending
}

// IsProcessing verifica si está en proceso
func (cp *CreditPurchase) IsProcessing() bool {
	return cp.Status == CreditPurchaseStatusProcessing
}

// IsCompleted verifica si está completado
func (cp *CreditPurchase) IsCompleted() bool {
	return cp.Status == CreditPurchaseStatusCompleted
}

// IsFailed verifica si falló
func (cp *CreditPurchase) IsFailed() bool {
	return cp.Status == CreditPurchaseStatusFailed
}

// IsExpired verifica si expiró
func (cp *CreditPurchase) IsExpired() bool {
	if cp.Status == CreditPurchaseStatusExpired {
		return true
	}
	// Verificar si pasó el tiempo de expiración
	return time.Now().After(cp.ExpiresAt) &&
		(cp.Status == CreditPurchaseStatusPending || cp.Status == CreditPurchaseStatusProcessing)
}

// MarkAsProcessing marca como en proceso (usuario en Pagadito)
func (cp *CreditPurchase) MarkAsProcessing(pagaditoToken string) error {
	if !cp.IsPending() {
		return fmt.Errorf("solo se puede marcar como 'processing' compras pendientes (estado actual: %s)", cp.Status)
	}
	cp.Status = CreditPurchaseStatusProcessing
	cp.PagaditoToken = &pagaditoToken
	cp.UpdatedAt = time.Now()
	return nil
}

// MarkAsCompleted marca como completado (pago exitoso)
func (cp *CreditPurchase) MarkAsCompleted(pagaditoReference string, walletTransactionID int64) error {
	if !cp.IsProcessing() && !cp.IsPending() {
		return fmt.Errorf("solo se puede completar compras en proceso o pendientes (estado actual: %s)", cp.Status)
	}
	cp.Status = CreditPurchaseStatusCompleted
	cp.PagaditoReference = &pagaditoReference
	statusStr := string(PagaditoStatusCompleted)
	cp.PagaditoStatus = &statusStr
	cp.WalletTransactionID = &walletTransactionID
	now := time.Now()
	cp.CompletedAt = &now
	cp.UpdatedAt = now
	return nil
}

// MarkAsFailed marca como fallido
func (cp *CreditPurchase) MarkAsFailed(reason string, pagaditoStatus PagaditoStatus) error {
	cp.Status = CreditPurchaseStatusFailed
	cp.ErrorMessage = &reason
	if pagaditoStatus != "" {
		statusStr := string(pagaditoStatus)
		cp.PagaditoStatus = &statusStr
	}
	now := time.Now()
	cp.FailedAt = &now
	cp.UpdatedAt = now
	return nil
}

// MarkAsExpired marca como expirado
func (cp *CreditPurchase) MarkAsExpired() error {
	if cp.IsCompleted() {
		return fmt.Errorf("no se puede marcar como expirado una compra completada")
	}
	cp.Status = CreditPurchaseStatusExpired
	cp.UpdatedAt = time.Now()
	return nil
}

// Validate valida la compra de créditos
func (cp *CreditPurchase) Validate() error {
	if cp.UserID <= 0 {
		return fmt.Errorf("user_id es requerido")
	}

	if cp.WalletID <= 0 {
		return fmt.Errorf("wallet_id es requerido")
	}

	if cp.DesiredCredit.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("desired_credit debe ser mayor a cero")
	}

	if cp.ChargeAmount.LessThan(cp.DesiredCredit) {
		return fmt.Errorf("charge_amount debe ser mayor o igual a desired_credit")
	}

	if cp.Currency == "" {
		return fmt.Errorf("currency es requerida")
	}

	if len(cp.Currency) != 3 {
		return fmt.Errorf("currency debe tener 3 caracteres (ISO 4217)")
	}

	if cp.ERN == "" {
		return fmt.Errorf("ERN es requerido")
	}

	if cp.IdempotencyKey == "" {
		return fmt.Errorf("idempotency_key es requerido")
	}

	return nil
}

// GenerateERN genera un External Reference Number único para Pagadito
// Formato: CP_{user_id}_{timestamp}_{random}
func GenerateERN(userID int64) (string, error) {
	// Timestamp unix
	timestamp := time.Now().Unix()

	// Random hex (6 bytes = 12 caracteres hex)
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("error generando bytes aleatorios: %w", err)
	}
	randomHex := hex.EncodeToString(randomBytes)

	// Formato: CP_{user_id}_{timestamp}_{random}
	ern := fmt.Sprintf("CP_%d_%d_%s", userID, timestamp, randomHex)

	// Convertir a mayúsculas para consistencia
	return strings.ToUpper(ern), nil
}

// CreditPurchaseRepository define el contrato para el repositorio
type CreditPurchaseRepository interface {
	// Create crea una nueva compra
	Create(purchase *CreditPurchase) error

	// FindByID busca por ID
	FindByID(id int64) (*CreditPurchase, error)

	// FindByUUID busca por UUID
	FindByUUID(uuid string) (*CreditPurchase, error)

	// FindByERN busca por ERN (External Reference Number)
	FindByERN(ern string) (*CreditPurchase, error)

	// FindByIdempotencyKey busca por clave de idempotencia
	FindByIdempotencyKey(key string) (*CreditPurchase, error)

	// FindByPagaditoToken busca por token de Pagadito
	FindByPagaditoToken(token string) (*CreditPurchase, error)

	// FindByUserID busca compras de un usuario (paginado)
	FindByUserID(userID int64, limit, offset int) ([]*CreditPurchase, int64, error)

	// Update actualiza una compra
	Update(purchase *CreditPurchase) error

	// MarkExpired marca como expiradas las compras que superaron el TTL
	MarkExpired() (int64, error)
}
