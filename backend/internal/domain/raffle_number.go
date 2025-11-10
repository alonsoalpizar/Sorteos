package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// RaffleNumberStatus representa el estado de un número de sorteo
type RaffleNumberStatus string

const (
	RaffleNumberStatusAvailable RaffleNumberStatus = "available"
	RaffleNumberStatusReserved  RaffleNumberStatus = "reserved"
	RaffleNumberStatusSold      RaffleNumberStatus = "sold"
)

// RaffleNumber representa un número individual de un sorteo
type RaffleNumber struct {
	ID int64

	// Raffle reference
	RaffleID int64

	// Number info
	Number string
	Status RaffleNumberStatus

	// Buyer info (only when sold)
	UserID        *int64
	ReservationID *int64
	PaymentID     *int64

	// Reservation tracking
	ReservedAt    *time.Time
	ReservedUntil *time.Time
	ReservedBy    *int64 // User who reserved (may differ from buyer)

	// Sale tracking
	SoldAt *time.Time
	Price  *decimal.Decimal

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewRaffleNumber crea un nuevo número de sorteo
func NewRaffleNumber(raffleID int64, number string) *RaffleNumber {
	return &RaffleNumber{
		RaffleID:  raffleID,
		Number:    number,
		Status:    RaffleNumberStatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate valida el número de sorteo
func (rn *RaffleNumber) Validate() error {
	if rn.RaffleID <= 0 {
		return fmt.Errorf("raffle_id es requerido")
	}

	if rn.Number == "" {
		return fmt.Errorf("el número es requerido")
	}

	// Status validation
	if rn.Status != RaffleNumberStatusAvailable &&
		rn.Status != RaffleNumberStatusReserved &&
		rn.Status != RaffleNumberStatusSold {
		return fmt.Errorf("estado inválido")
	}

	// Validations for sold status
	if rn.Status == RaffleNumberStatusSold {
		if rn.UserID == nil {
			return fmt.Errorf("user_id es requerido para números vendidos")
		}
		if rn.Price == nil {
			return fmt.Errorf("price es requerido para números vendidos")
		}
		if rn.SoldAt == nil {
			return fmt.Errorf("sold_at es requerido para números vendidos")
		}
	}

	// Validations for reserved status
	if rn.Status == RaffleNumberStatusReserved {
		if rn.ReservedUntil == nil {
			return fmt.Errorf("reserved_until es requerido para números reservados")
		}
		if rn.ReservedUntil.Before(time.Now()) {
			return fmt.Errorf("la reserva ha expirado")
		}
	}

	return nil
}

// Reserve reserva el número por un tiempo determinado
func (rn *RaffleNumber) Reserve(userID int64, reservationID int64, duration time.Duration) error {
	if rn.Status != RaffleNumberStatusAvailable {
		return fmt.Errorf("el número no está disponible para reservar (estado: %s)", rn.Status)
	}

	now := time.Now()
	until := now.Add(duration)

	rn.Status = RaffleNumberStatusReserved
	rn.ReservedAt = &now
	rn.ReservedUntil = &until
	rn.ReservedBy = &userID
	rn.ReservationID = &reservationID
	rn.UpdatedAt = now

	return nil
}

// CancelReservation cancela la reserva y libera el número
func (rn *RaffleNumber) CancelReservation() error {
	if rn.Status != RaffleNumberStatusReserved {
		return fmt.Errorf("el número no está reservado")
	}

	rn.Status = RaffleNumberStatusAvailable
	rn.ReservedAt = nil
	rn.ReservedUntil = nil
	rn.ReservedBy = nil
	rn.ReservationID = nil
	rn.UpdatedAt = time.Now()

	return nil
}

// MarkAsSold marca el número como vendido
func (rn *RaffleNumber) MarkAsSold(userID int64, paymentID int64, price decimal.Decimal) error {
	// Se puede vender un número disponible o reservado
	if rn.Status != RaffleNumberStatusAvailable && rn.Status != RaffleNumberStatusReserved {
		return fmt.Errorf("el número no puede ser vendido en su estado actual (estado: %s)", rn.Status)
	}

	now := time.Now()

	rn.Status = RaffleNumberStatusSold
	rn.UserID = &userID
	rn.PaymentID = &paymentID
	rn.Price = &price
	rn.SoldAt = &now
	rn.UpdatedAt = now

	// Limpiar info de reserva si existía
	rn.ReservedAt = nil
	rn.ReservedUntil = nil
	rn.ReservedBy = nil
	rn.ReservationID = nil

	return nil
}

// IsAvailable verifica si el número está disponible
func (rn *RaffleNumber) IsAvailable() bool {
	return rn.Status == RaffleNumberStatusAvailable
}

// IsReserved verifica si el número está reservado
func (rn *RaffleNumber) IsReserved() bool {
	return rn.Status == RaffleNumberStatusReserved
}

// IsSold verifica si el número está vendido
func (rn *RaffleNumber) IsSold() bool {
	return rn.Status == RaffleNumberStatusSold
}

// IsReservationExpired verifica si la reserva ha expirado
func (rn *RaffleNumber) IsReservationExpired() bool {
	if !rn.IsReserved() || rn.ReservedUntil == nil {
		return false
	}

	return rn.ReservedUntil.Before(time.Now())
}

// ReleaseIfExpired libera el número si la reserva ha expirado
func (rn *RaffleNumber) ReleaseIfExpired() bool {
	if rn.IsReservationExpired() {
		rn.CancelReservation()
		return true
	}
	return false
}

// TimeUntilExpiration retorna el tiempo restante de la reserva
func (rn *RaffleNumber) TimeUntilExpiration() time.Duration {
	if !rn.IsReserved() || rn.ReservedUntil == nil {
		return 0
	}

	remaining := time.Until(*rn.ReservedUntil)
	if remaining < 0 {
		return 0
	}

	return remaining
}

// CanBeReservedBy verifica si el número puede ser reservado por un usuario
func (rn *RaffleNumber) CanBeReservedBy(userID int64) bool {
	if rn.IsAvailable() {
		return true
	}

	// Si está reservado, verificar si es por el mismo usuario y la reserva sigue activa
	if rn.IsReserved() && rn.ReservedBy != nil && *rn.ReservedBy == userID {
		return !rn.IsReservationExpired()
	}

	return false
}
