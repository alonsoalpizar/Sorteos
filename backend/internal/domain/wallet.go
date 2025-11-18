package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// WalletStatus representa el estado de la billetera
type WalletStatus string

const (
	WalletStatusActive WalletStatus = "active"
	WalletStatusFrozen WalletStatus = "frozen"
	WalletStatusClosed WalletStatus = "closed"
)

// Wallet representa la billetera de un usuario
type Wallet struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	UUID   string `json:"uuid" gorm:"type:uuid;unique;not null"`
	UserID int64  `json:"user_id" gorm:"uniqueIndex;not null"`

	// Saldos
	Balance        decimal.Decimal `json:"balance" gorm:"type:decimal(12,2);not null;default:0.00"`
	PendingBalance decimal.Decimal `json:"pending_balance" gorm:"type:decimal(12,2);not null;default:0.00"`
	Currency       string          `json:"currency" gorm:"type:varchar(3);default:'USD';not null"`

	// Estado
	Status WalletStatus `json:"status" gorm:"type:wallet_status;default:'active';not null"`

	// Auditoría
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (Wallet) TableName() string {
	return "wallets"
}

// IsActive verifica si la billetera está activa
func (w *Wallet) IsActive() bool {
	return w.Status == WalletStatusActive
}

// HasSufficientBalance verifica si hay saldo suficiente
func (w *Wallet) HasSufficientBalance(amount decimal.Decimal) bool {
	return w.Balance.GreaterThanOrEqual(amount)
}

// CanDebit verifica si se puede debitar un monto
func (w *Wallet) CanDebit(amount decimal.Decimal) error {
	if !w.IsActive() {
		return fmt.Errorf("billetera no está activa (estado: %s)", w.Status)
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("el monto debe ser mayor a cero")
	}

	if !w.HasSufficientBalance(amount) {
		return fmt.Errorf("saldo insuficiente (disponible: %s, requerido: %s)", w.Balance.String(), amount.String())
	}

	return nil
}

// CanCredit verifica si se puede acreditar un monto
func (w *Wallet) CanCredit(amount decimal.Decimal) error {
	if !w.IsActive() {
		return fmt.Errorf("billetera no está activa (estado: %s)", w.Status)
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("el monto debe ser mayor a cero")
	}

	return nil
}

// Debit debita un monto del saldo
func (w *Wallet) Debit(amount decimal.Decimal) error {
	if err := w.CanDebit(amount); err != nil {
		return err
	}

	w.Balance = w.Balance.Sub(amount)
	w.UpdatedAt = time.Now()
	return nil
}

// Credit acredita un monto al saldo
func (w *Wallet) Credit(amount decimal.Decimal) error {
	if err := w.CanCredit(amount); err != nil {
		return err
	}

	w.Balance = w.Balance.Add(amount)
	w.UpdatedAt = time.Now()
	return nil
}

// CreditPending acredita un monto al saldo pendiente
func (w *Wallet) CreditPending(amount decimal.Decimal) error {
	if !w.IsActive() {
		return fmt.Errorf("billetera no está activa")
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("el monto debe ser mayor a cero")
	}

	w.PendingBalance = w.PendingBalance.Add(amount)
	w.UpdatedAt = time.Now()
	return nil
}

// ConfirmPending mueve saldo pendiente a saldo disponible
func (w *Wallet) ConfirmPending(amount decimal.Decimal) error {
	if amount.GreaterThan(w.PendingBalance) {
		return fmt.Errorf("saldo pendiente insuficiente")
	}

	w.PendingBalance = w.PendingBalance.Sub(amount)
	w.Balance = w.Balance.Add(amount)
	w.UpdatedAt = time.Now()
	return nil
}

// Freeze congela la billetera
func (w *Wallet) Freeze() error {
	if w.Status == WalletStatusClosed {
		return fmt.Errorf("no se puede congelar una billetera cerrada")
	}

	w.Status = WalletStatusFrozen
	w.UpdatedAt = time.Now()
	return nil
}

// Unfreeze descongela la billetera
func (w *Wallet) Unfreeze() error {
	if w.Status != WalletStatusFrozen {
		return fmt.Errorf("la billetera no está congelada")
	}

	w.Status = WalletStatusActive
	w.UpdatedAt = time.Now()
	return nil
}

// Close cierra la billetera
func (w *Wallet) Close() error {
	if !w.Balance.IsZero() || !w.PendingBalance.IsZero() {
		return fmt.Errorf("no se puede cerrar una billetera con saldo")
	}

	w.Status = WalletStatusClosed
	w.UpdatedAt = time.Now()
	return nil
}

// Validate valida la billetera
func (w *Wallet) Validate() error {
	if w.UserID <= 0 {
		return fmt.Errorf("user_id es requerido")
	}

	if w.Balance.LessThan(decimal.Zero) {
		return fmt.Errorf("el saldo no puede ser negativo")
	}

	if w.PendingBalance.LessThan(decimal.Zero) {
		return fmt.Errorf("el saldo pendiente no puede ser negativo")
	}

	if w.Currency == "" {
		return fmt.Errorf("la moneda es requerida")
	}

	if len(w.Currency) != 3 {
		return fmt.Errorf("la moneda debe tener 3 caracteres (ISO 4217)")
	}

	return nil
}

// WalletRepository define el contrato para el repositorio de billeteras
type WalletRepository interface {
	// Create crea una nueva billetera
	Create(wallet *Wallet) error

	// FindByID busca una billetera por ID
	FindByID(id int64) (*Wallet, error)

	// FindByUUID busca una billetera por UUID
	FindByUUID(uuid string) (*Wallet, error)

	// FindByUserID busca una billetera por ID de usuario
	FindByUserID(userID int64) (*Wallet, error)

	// Update actualiza una billetera existente
	Update(wallet *Wallet) error

	// UpdateBalance actualiza solo el saldo (optimización)
	UpdateBalance(walletID int64, balance decimal.Decimal) error

	// Lock adquiere un lock para operaciones concurrentes
	Lock(walletID int64) error

	// Unlock libera un lock
	Unlock(walletID int64) error

	// WithTransaction ejecuta una función dentro de una transacción
	WithTransaction(fn func(repo WalletRepository) error) error
}
