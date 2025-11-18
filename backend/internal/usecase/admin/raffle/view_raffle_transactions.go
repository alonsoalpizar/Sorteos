package raffle

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// TransactionEvent evento en el timeline
type TransactionEvent struct {
	Type      string    `json:"type"` // reservation, payment, refund, status_change, note
	Timestamp time.Time `json:"timestamp"`
	UserID    *int64    `json:"user_id,omitempty"`
	UserName  *string   `json:"user_name,omitempty"`
	Amount    *float64  `json:"amount,omitempty"`
	Status    *string   `json:"status,omitempty"`
	Details   string    `json:"details"`
}

// RaffleTransactionMetrics métricas calculadas
type RaffleTransactionMetrics struct {
	TotalReservations  int     `json:"total_reservations"`
	TotalPayments      int     `json:"total_payments"`
	TotalRefunds       int     `json:"total_refunds"`
	ConversionRate     float64 `json:"conversion_rate"` // payments / reservations
	RefundRate         float64 `json:"refund_rate"`     // refunds / payments
	TotalRevenue       float64 `json:"total_revenue"`
	TotalRefunded      float64 `json:"total_refunded"`
	NetRevenue         float64 `json:"net_revenue"`
}

// ViewRaffleTransactionsOutput resultado
type ViewRaffleTransactionsOutput struct {
	Raffle   *domain.Raffle            `json:"raffle"`
	Timeline []*TransactionEvent       `json:"timeline"`
	Metrics  *RaffleTransactionMetrics `json:"metrics"`
}

// ViewRaffleTransactionsUseCase caso de uso para ver timeline de transacciones
type ViewRaffleTransactionsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewViewRaffleTransactionsUseCase crea una nueva instancia
func NewViewRaffleTransactionsUseCase(db *gorm.DB, log *logger.Logger) *ViewRaffleTransactionsUseCase {
	return &ViewRaffleTransactionsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ViewRaffleTransactionsUseCase) Execute(ctx context.Context, raffleID int64, adminID int64) (*ViewRaffleTransactionsOutput, error) {
	// Obtener rifa
	var raffle domain.Raffle
	if err := uc.db.Where("id = ?", raffleID).First(&raffle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		uc.log.Error("Error finding raffle", logger.Int64("raffle_id", raffleID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	timeline := make([]*TransactionEvent, 0)

	// 1. Obtener reservations
	var reservations []struct {
		UserID    *int64
		UserName  string
		CreatedAt time.Time
		Count     int
	}

	if err := uc.db.Table("reservations").
		Select("reservations.user_id, CONCAT(users.first_name, ' ', users.last_name) as user_name, reservations.created_at, COUNT(*) as count").
		Joins("LEFT JOIN users ON users.id = reservations.user_id").
		Where("reservations.raffle_id = ?", raffleID).
		Group("reservations.user_id, CONCAT(users.first_name, ' ', users.last_name), reservations.created_at").
		Find(&reservations).Error; err != nil {
		uc.log.Error("Error getting reservations", logger.Error(err))
		// No fallar por esto
	} else {
		for _, res := range reservations {
			timeline = append(timeline, &TransactionEvent{
				Type:      "reservation",
				Timestamp: res.CreatedAt,
				UserID:    res.UserID,
				UserName:  &res.UserName,
				Details:   fmt.Sprintf("Reserved %d numbers", res.Count),
			})
		}
	}

	// 2. Obtener payments
	var payments []struct {
		ID        int64
		UserID    int64
		UserName  string
		Amount    float64
		Status    string
		CreatedAt time.Time
	}

	if err := uc.db.Table("payments").
		Select("payments.id, payments.user_id, CONCAT(users.first_name, ' ', users.last_name) as user_name, payments.amount, payments.status, payments.created_at").
		Joins("LEFT JOIN users ON users.id = payments.user_id").
		Where("payments.raffle_id = ?", raffleID).
		Find(&payments).Error; err != nil {
		uc.log.Error("Error getting payments", logger.Error(err))
		// No fallar por esto
	} else {
		for _, pay := range payments {
			eventType := "payment"
			if pay.Status == "refunded" {
				eventType = "refund"
			}

			timeline = append(timeline, &TransactionEvent{
				Type:      eventType,
				Timestamp: pay.CreatedAt,
				UserID:    &pay.UserID,
				UserName:  &pay.UserName,
				Amount:    &pay.Amount,
				Status:    &pay.Status,
				Details:   fmt.Sprintf("Payment %s: $%.2f", pay.Status, pay.Amount),
			})
		}
	}

	// 3. Obtener cambios de estado de audit logs (si existe la tabla)
	var auditLogs []struct {
		UserID    *int64
		UserName  *string
		Action    string
		Details   string
		CreatedAt time.Time
	}

	if err := uc.db.Table("audit_logs").
		Select("audit_logs.user_id, CONCAT(users.first_name, ' ', users.last_name) as user_name, audit_logs.action, audit_logs.details, audit_logs.created_at").
		Joins("LEFT JOIN users ON users.id = audit_logs.user_id").
		Where("audit_logs.entity_type = ? AND audit_logs.entity_id = ?", "raffle", raffleID).
		Find(&auditLogs).Error; err == nil {
		for _, log := range auditLogs {
			timeline = append(timeline, &TransactionEvent{
				Type:      "status_change",
				Timestamp: log.CreatedAt,
				UserID:    log.UserID,
				UserName:  log.UserName,
				Details:   fmt.Sprintf("%s: %s", log.Action, log.Details),
			})
		}
	}

	// Ordenar timeline por timestamp (más reciente primero)
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Timestamp.After(timeline[j].Timestamp)
	})

	// Calcular métricas
	metrics := &RaffleTransactionMetrics{}

	// Contar reservations únicas
	metrics.TotalReservations = len(reservations)

	// Contar payments y refunds
	totalPayments := 0
	totalRefunds := 0
	totalRevenue := 0.0
	totalRefunded := 0.0

	for _, pay := range payments {
		if pay.Status == "succeeded" || pay.Status == "refunded" {
			totalPayments++
			totalRevenue += pay.Amount
		}
		if pay.Status == "refunded" {
			totalRefunds++
			totalRefunded += pay.Amount
		}
	}

	metrics.TotalPayments = totalPayments
	metrics.TotalRefunds = totalRefunds
	metrics.TotalRevenue = totalRevenue
	metrics.TotalRefunded = totalRefunded
	metrics.NetRevenue = totalRevenue - totalRefunded

	// Calcular conversion rate
	if metrics.TotalReservations > 0 {
		metrics.ConversionRate = float64(metrics.TotalPayments) / float64(metrics.TotalReservations) * 100
	}

	// Calcular refund rate
	if metrics.TotalPayments > 0 {
		metrics.RefundRate = float64(metrics.TotalRefunds) / float64(metrics.TotalPayments) * 100
	}

	// Log auditoría
	uc.log.Info("Admin viewed raffle transactions",
		logger.Int64("admin_id", adminID),
		logger.Int64("raffle_id", raffleID),
		logger.String("action", "admin_view_raffle_transactions"))

	return &ViewRaffleTransactionsOutput{
		Raffle:   &raffle,
		Timeline: timeline,
		Metrics:  metrics,
	}, nil
}
