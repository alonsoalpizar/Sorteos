package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// NewDecimalFromFloat es un helper para convertir float64 a decimal
func NewDecimalFromFloat(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

// UserEarnings representa el resumen de ganancias de un usuario
type UserEarnings struct {
	TotalCollected     decimal.Decimal   `json:"total_collected"`      // Total recolectado de todos los sorteos
	PlatformCommission decimal.Decimal   `json:"platform_commission"`  // Total de comisión (10%)
	NetEarnings        decimal.Decimal   `json:"net_earnings"`         // Ganancias netas (total - comisión)
	CompletedRaffles   int               `json:"completed_raffles"`    // Cantidad de sorteos completados
	Raffles            []RaffleEarning   `json:"raffles"`              // Desglose por sorteo
}

// RaffleEarning representa las ganancias de un sorteo específico
type RaffleEarning struct {
	RaffleID           int64           `json:"raffle_id"`
	RaffleUUID         string          `json:"raffle_uuid"`
	Title              string          `json:"title"`
	DrawDate           time.Time       `json:"draw_date"`
	CompletedAt        *time.Time      `json:"completed_at"`
	TotalRevenue       decimal.Decimal `json:"total_revenue"`        // Total recolectado
	PlatformFeePercent decimal.Decimal `json:"platform_fee_percent"` // Porcentaje (10.00)
	PlatformFeeAmount  decimal.Decimal `json:"platform_fee_amount"`  // Comisión (₡)
	NetAmount          decimal.Decimal `json:"net_amount"`           // Ganancia neta
	SettlementStatus   string          `json:"settlement_status"`    // pending | completed
	SettledAt          *time.Time      `json:"settled_at"`           // Cuando se depositó
}

// RaffleRepository necesita este método para earnings
// (se agregará a la interfaz existente)
type EarningsRepository interface {
	// GetUserCompletedRaffles obtiene los sorteos completados de un organizador
	// con status='completed' y settlement_status='completed'
	GetUserCompletedRaffles(userID int64, limit, offset int) ([]RaffleEarning, error)

	// GetUserEarningsSummary obtiene el resumen total de ganancias
	GetUserEarningsSummary(userID int64) (*UserEarnings, error)
}
