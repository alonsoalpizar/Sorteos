package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// RaffleStatus representa el estado de un sorteo
type RaffleStatus string

const (
	RaffleStatusDraft     RaffleStatus = "draft"
	RaffleStatusActive    RaffleStatus = "active"
	RaffleStatusSuspended RaffleStatus = "suspended"
	RaffleStatusCompleted RaffleStatus = "completed"
	RaffleStatusCancelled RaffleStatus = "cancelled"
)

// DrawMethod representa el método de selección del ganador
type DrawMethod string

const (
	DrawMethodLoteriaCostaRica DrawMethod = "loteria_nacional_cr"
	DrawMethodManual           DrawMethod = "manual"
	DrawMethodRandom           DrawMethod = "random"
)

// SettlementStatus representa el estado de liquidación
type SettlementStatus string

const (
	SettlementStatusPending    SettlementStatus = "pending"
	SettlementStatusProcessing SettlementStatus = "processing"
	SettlementStatusCompleted  SettlementStatus = "completed"
	SettlementStatusFailed     SettlementStatus = "failed"
)

// Raffle representa un sorteo/rifa en el sistema
type Raffle struct {
	ID   int64
	UUID uuid.UUID

	// Owner
	UserID int64

	// Basic info
	Title       string
	Description string
	Status      RaffleStatus
	CategoryID  *int64

	// Pricing
	PricePerNumber      decimal.Decimal
	TotalNumbers        int
	MinNumber           int
	MaxNumber           int

	// Draw info
	DrawDate   time.Time
	DrawMethod DrawMethod

	// Winner info
	WinnerNumber *string
	WinnerUserID *int64

	// Counters
	SoldCount     int
	ReservedCount int

	// Revenue
	TotalRevenue         decimal.Decimal
	PlatformFeePercentage decimal.Decimal
	PlatformFeeAmount    decimal.Decimal
	NetAmount            decimal.Decimal

	// Settlement
	SettledAt        *time.Time
	SettlementStatus SettlementStatus

	// Metadata
	Metadata datatypes.JSON

	// Timestamps
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
	CompletedAt *time.Time
	DeletedAt   *time.Time
}

// NewRaffle crea una nueva instancia de Raffle con valores por defecto
func NewRaffle(userID int64, title string, pricePerNumber decimal.Decimal, totalNumbers int, drawDate time.Time) *Raffle {
	return &Raffle{
		UUID:                  uuid.New(),
		UserID:                userID,
		Title:                 title,
		Status:                RaffleStatusDraft,
		PricePerNumber:        pricePerNumber,
		TotalNumbers:          totalNumbers,
		MinNumber:             0,
		MaxNumber:             totalNumbers - 1,
		DrawDate:              drawDate,
		DrawMethod:            DrawMethodLoteriaCostaRica,
		SoldCount:             0,
		ReservedCount:         0,
		TotalRevenue:          decimal.Zero,
		PlatformFeePercentage: decimal.NewFromFloat(10.0), // 10% default
		PlatformFeeAmount:     decimal.Zero,
		NetAmount:             decimal.Zero,
		SettlementStatus:      SettlementStatusPending,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// Validate valida la entidad Raffle
func (r *Raffle) Validate() error {
	// Title validation
	if r.Title == "" {
		return fmt.Errorf("el título es requerido")
	}
	if len(r.Title) < 5 {
		return fmt.Errorf("el título debe tener al menos 5 caracteres")
	}
	if len(r.Title) > 255 {
		return fmt.Errorf("el título no puede exceder 255 caracteres")
	}

	// Price validation
	if r.PricePerNumber.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("el precio por número debe ser mayor a 0")
	}
	if r.PricePerNumber.GreaterThan(decimal.NewFromInt(1000000)) {
		return fmt.Errorf("el precio por número no puede exceder 1,000,000")
	}

	// Numbers validation
	if r.TotalNumbers <= 0 {
		return fmt.Errorf("el total de números debe ser mayor a 0")
	}
	if r.TotalNumbers > 10000 {
		return fmt.Errorf("el total de números no puede exceder 10,000")
	}
	if r.MaxNumber <= r.MinNumber {
		return fmt.Errorf("el número máximo debe ser mayor al número mínimo")
	}
	if r.MaxNumber-r.MinNumber+1 != r.TotalNumbers {
		return fmt.Errorf("el rango de números no coincide con el total de números")
	}

	// Draw date validation
	if r.DrawDate.Before(time.Now()) {
		return fmt.Errorf("la fecha de sorteo debe ser futura")
	}

	// Draw method validation
	if r.DrawMethod != DrawMethodLoteriaCostaRica &&
		r.DrawMethod != DrawMethodManual &&
		r.DrawMethod != DrawMethodRandom {
		return fmt.Errorf("método de sorteo inválido")
	}

	// Counters validation
	if r.SoldCount < 0 || r.SoldCount > r.TotalNumbers {
		return fmt.Errorf("el contador de vendidos es inválido")
	}
	if r.ReservedCount < 0 || r.ReservedCount > r.TotalNumbers {
		return fmt.Errorf("el contador de reservados es inválido")
	}

	return nil
}

// CanBePublished verifica si el sorteo puede ser publicado
func (r *Raffle) CanBePublished() bool {
	return r.Status == RaffleStatusDraft &&
		r.Validate() == nil &&
		r.DrawDate.After(time.Now().Add(24*time.Hour)) // Al menos 24h de anticipación
}

// Publish marca el sorteo como activo
func (r *Raffle) Publish() error {
	if !r.CanBePublished() {
		return fmt.Errorf("el sorteo no puede ser publicado en su estado actual")
	}

	now := time.Now()
	r.Status = RaffleStatusActive
	r.PublishedAt = &now
	r.UpdatedAt = now

	return nil
}

// Suspend suspende el sorteo
func (r *Raffle) Suspend() error {
	if r.Status != RaffleStatusActive {
		return fmt.Errorf("solo se pueden suspender sorteos activos")
	}

	r.Status = RaffleStatusSuspended
	r.UpdatedAt = time.Now()

	return nil
}

// Activate reactiva un sorteo suspendido
func (r *Raffle) Activate() error {
	if r.Status != RaffleStatusSuspended {
		return fmt.Errorf("solo se pueden activar sorteos suspendidos")
	}

	// Verificar que la fecha de sorteo siga siendo futura
	if r.DrawDate.Before(time.Now()) {
		return fmt.Errorf("no se puede activar un sorteo con fecha pasada")
	}

	r.Status = RaffleStatusActive
	r.UpdatedAt = time.Now()

	return nil
}

// Complete marca el sorteo como completado
func (r *Raffle) Complete(winnerNumber string, winnerUserID *int64) error {
	if r.Status != RaffleStatusActive && r.Status != RaffleStatusSuspended {
		return fmt.Errorf("solo se pueden completar sorteos activos o suspendidos")
	}

	now := time.Now()
	r.Status = RaffleStatusCompleted
	r.WinnerNumber = &winnerNumber
	r.WinnerUserID = winnerUserID
	r.CompletedAt = &now
	r.UpdatedAt = now

	return nil
}

// Cancel cancela el sorteo
func (r *Raffle) Cancel() error {
	if r.Status == RaffleStatusCompleted {
		return fmt.Errorf("no se puede cancelar un sorteo completado")
	}
	if r.SoldCount > 0 {
		return fmt.Errorf("no se puede cancelar un sorteo con números vendidos (debe reembolsarse primero)")
	}

	r.Status = RaffleStatusCancelled
	r.UpdatedAt = time.Now()

	return nil
}

// IsActive verifica si el sorteo está activo
func (r *Raffle) IsActive() bool {
	return r.Status == RaffleStatusActive && r.DeletedAt == nil
}

// IsDraft verifica si el sorteo está en borrador
func (r *Raffle) IsDraft() bool {
	return r.Status == RaffleStatusDraft
}

// IsCompleted verifica si el sorteo está completado
func (r *Raffle) IsCompleted() bool {
	return r.Status == RaffleStatusCompleted
}

// CanBeSoldOut verifica si el sorteo puede agotarse
func (r *Raffle) CanBeSoldOut() bool {
	return r.SoldCount >= r.TotalNumbers
}

// IsSoldOut verifica si todos los números están vendidos
func (r *Raffle) IsSoldOut() bool {
	return r.SoldCount == r.TotalNumbers
}

// AvailableCount retorna la cantidad de números disponibles
func (r *Raffle) AvailableCount() int {
	return r.TotalNumbers - r.SoldCount - r.ReservedCount
}

// CanBeEdited verifica si el sorteo puede ser editado
func (r *Raffle) CanBeEdited() bool {
	// Solo se pueden editar borradores o sorteos activos sin números vendidos
	return (r.IsDraft() || (r.IsActive() && r.SoldCount == 0)) && r.DeletedAt == nil
}

// CalculateRevenue calcula los ingresos del sorteo (se ejecuta automáticamente en la DB)
func (r *Raffle) CalculateRevenue() {
	r.TotalRevenue = r.PricePerNumber.Mul(decimal.NewFromInt(int64(r.SoldCount)))
	r.PlatformFeeAmount = r.TotalRevenue.Mul(r.PlatformFeePercentage.Div(decimal.NewFromInt(100)))
	r.NetAmount = r.TotalRevenue.Sub(r.PlatformFeeAmount)
}

// CanBeSettled verifica si el sorteo puede ser liquidado
func (r *Raffle) CanBeSettled() bool {
	return r.IsCompleted() &&
		r.SettlementStatus == SettlementStatusPending &&
		r.WinnerNumber != nil
}

// MarkAsSettled marca el sorteo como liquidado
func (r *Raffle) MarkAsSettled() error {
	if !r.CanBeSettled() {
		return fmt.Errorf("el sorteo no puede ser liquidado en su estado actual")
	}

	now := time.Now()
	r.SettledAt = &now
	r.SettlementStatus = SettlementStatusCompleted
	r.UpdatedAt = now

	return nil
}
