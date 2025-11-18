package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// TransactionType representa el tipo de transacción
type TransactionType string

const (
	TransactionTypeDeposit        TransactionType = "deposit"          // Compra de créditos vía procesador
	TransactionTypeWithdrawal     TransactionType = "withdrawal"       // Retiro a cuenta bancaria
	TransactionTypePurchase       TransactionType = "purchase"         // Pago de sorteo
	TransactionTypeRefund         TransactionType = "refund"           // Devolución de compra
	TransactionTypePrizeClaim     TransactionType = "prize_claim"      // Premio ganado
	TransactionTypeSettlementPayout TransactionType = "settlement_payout" // Pago de liquidación a organizador
	TransactionTypeAdjustment     TransactionType = "adjustment"       // Ajuste manual (admin)
)

// TransactionStatus representa el estado de la transacción
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed" // Transacción revertida
)

// WalletTransaction representa una transacción de billetera
type WalletTransaction struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	UUID     string `json:"uuid" gorm:"type:uuid;unique;not null"`
	WalletID int64  `json:"wallet_id" gorm:"index;not null"`
	UserID   int64  `json:"user_id" gorm:"index;not null"`

	// Detalles de la transacción
	Type   TransactionType   `json:"type" gorm:"type:transaction_type;not null"`
	Amount decimal.Decimal   `json:"amount" gorm:"type:decimal(12,2);not null"`
	Status TransactionStatus `json:"status" gorm:"type:transaction_status;default:'pending';not null"`

	// Snapshots de saldo (para auditoría)
	BalanceBefore decimal.Decimal `json:"balance_before" gorm:"type:decimal(12,2);not null"`
	BalanceAfter  decimal.Decimal `json:"balance_after" gorm:"type:decimal(12,2);not null"`

	// Referencias externas (polimórfico)
	ReferenceType *string `json:"reference_type,omitempty"` // "payment", "settlement", "raffle", "admin"
	ReferenceID   *int64  `json:"reference_id,omitempty"`

	// Idempotencia
	IdempotencyKey string `json:"idempotency_key" gorm:"uniqueIndex;not null"`

	// Metadata adicional (JSONB)
	Metadata datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Notas (para ajustes manuales)
	Notes *string `json:"notes,omitempty"`

	// Auditoría
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	FailedAt    *time.Time `json:"failed_at,omitempty"`
	ReversedAt  *time.Time `json:"reversed_at,omitempty"`
}

// TableName especifica el nombre de la tabla
func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

// IsCompleted verifica si la transacción está completada
func (wt *WalletTransaction) IsCompleted() bool {
	return wt.Status == TransactionStatusCompleted
}

// IsPending verifica si la transacción está pendiente
func (wt *WalletTransaction) IsPending() bool {
	return wt.Status == TransactionStatusPending
}

// IsFailed verifica si la transacción falló
func (wt *WalletTransaction) IsFailed() bool {
	return wt.Status == TransactionStatusFailed
}

// IsReversed verifica si la transacción fue revertida
func (wt *WalletTransaction) IsReversed() bool {
	return wt.Status == TransactionStatusReversed
}

// MarkAsCompleted marca la transacción como completada
func (wt *WalletTransaction) MarkAsCompleted() error {
	if !wt.IsPending() {
		return fmt.Errorf("solo se pueden completar transacciones pendientes")
	}

	wt.Status = TransactionStatusCompleted
	now := time.Now()
	wt.CompletedAt = &now
	return nil
}

// MarkAsFailed marca la transacción como fallida
func (wt *WalletTransaction) MarkAsFailed(reason string) error {
	if !wt.IsPending() {
		return fmt.Errorf("solo se pueden marcar como fallidas transacciones pendientes")
	}

	wt.Status = TransactionStatusFailed
	now := time.Now()
	wt.FailedAt = &now
	if reason != "" {
		wt.Notes = &reason
	}
	return nil
}

// MarkAsReversed marca la transacción como revertida
func (wt *WalletTransaction) MarkAsReversed(reason string) error {
	if !wt.IsCompleted() {
		return fmt.Errorf("solo se pueden revertir transacciones completadas")
	}

	wt.Status = TransactionStatusReversed
	now := time.Now()
	wt.ReversedAt = &now
	if reason != "" {
		wt.Notes = &reason
	}
	return nil
}

// Validate valida la transacción
func (wt *WalletTransaction) Validate() error {
	if wt.WalletID <= 0 {
		return fmt.Errorf("wallet_id es requerido")
	}

	if wt.UserID <= 0 {
		return fmt.Errorf("user_id es requerido")
	}

	if wt.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("el monto debe ser mayor a cero")
	}

	if wt.IdempotencyKey == "" {
		return fmt.Errorf("idempotency_key es requerido")
	}

	if wt.Type == "" {
		return fmt.Errorf("el tipo de transacción es requerido")
	}

	// Validar tipos válidos
	validTypes := map[TransactionType]bool{
		TransactionTypeDeposit:          true,
		TransactionTypeWithdrawal:       true,
		TransactionTypePurchase:         true,
		TransactionTypeRefund:           true,
		TransactionTypePrizeClaim:       true,
		TransactionTypeSettlementPayout: true,
		TransactionTypeAdjustment:       true,
	}

	if !validTypes[wt.Type] {
		return fmt.Errorf("tipo de transacción inválido: %s", wt.Type)
	}

	// Validar saldos
	if wt.BalanceBefore.LessThan(decimal.Zero) {
		return fmt.Errorf("balance_before no puede ser negativo")
	}

	if wt.BalanceAfter.LessThan(decimal.Zero) {
		return fmt.Errorf("balance_after no puede ser negativo")
	}

	// Validar coherencia de saldos según tipo
	switch wt.Type {
	case TransactionTypeDeposit, TransactionTypeRefund, TransactionTypePrizeClaim, TransactionTypeSettlementPayout:
		// Debe incrementar el saldo
		expected := wt.BalanceBefore.Add(wt.Amount)
		if !wt.BalanceAfter.Equal(expected) {
			return fmt.Errorf("balance_after inválido para %s (esperado: %s, obtenido: %s)",
				wt.Type, expected.String(), wt.BalanceAfter.String())
		}
	case TransactionTypePurchase, TransactionTypeWithdrawal, TransactionTypeAdjustment:
		// Debe decrementar el saldo
		expected := wt.BalanceBefore.Sub(wt.Amount)
		if !wt.BalanceAfter.Equal(expected) {
			return fmt.Errorf("balance_after inválido para %s (esperado: %s, obtenido: %s)",
				wt.Type, expected.String(), wt.BalanceAfter.String())
		}
	}

	return nil
}

// IsDebit verifica si la transacción es un débito
func (wt *WalletTransaction) IsDebit() bool {
	return wt.Type == TransactionTypePurchase ||
		wt.Type == TransactionTypeWithdrawal ||
		(wt.Type == TransactionTypeAdjustment && wt.BalanceAfter.LessThan(wt.BalanceBefore))
}

// IsCredit verifica si la transacción es un crédito
func (wt *WalletTransaction) IsCredit() bool {
	return wt.Type == TransactionTypeDeposit ||
		wt.Type == TransactionTypeRefund ||
		wt.Type == TransactionTypePrizeClaim ||
		wt.Type == TransactionTypeSettlementPayout ||
		(wt.Type == TransactionTypeAdjustment && wt.BalanceAfter.GreaterThan(wt.BalanceBefore))
}

// WalletTransactionRepository define el contrato para el repositorio de transacciones
type WalletTransactionRepository interface {
	// Create crea una nueva transacción
	Create(tx *WalletTransaction) error

	// FindByID busca una transacción por ID
	FindByID(id int64) (*WalletTransaction, error)

	// FindByUUID busca una transacción por UUID
	FindByUUID(uuid string) (*WalletTransaction, error)

	// FindByIdempotencyKey busca una transacción por clave de idempotencia
	FindByIdempotencyKey(key string) (*WalletTransaction, error)

	// FindByWalletID busca transacciones de una billetera (paginado)
	FindByWalletID(walletID int64, limit, offset int) ([]*WalletTransaction, int64, error)

	// FindByUserID busca transacciones de un usuario (paginado)
	FindByUserID(userID int64, limit, offset int) ([]*WalletTransaction, int64, error)

	// FindByReference busca transacciones por referencia externa
	FindByReference(referenceType string, referenceID int64) ([]*WalletTransaction, error)

	// Update actualiza una transacción existente
	Update(tx *WalletTransaction) error

	// WithTransaction ejecuta una función dentro de una transacción
	WithTransaction(fn func(repo WalletTransactionRepository) error) error
}
