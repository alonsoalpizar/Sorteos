package settlement

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// SettlementFullDetails detalles completos de liquidación
type SettlementFullDetails struct {
	Settlement      *SettlementWithDetails
	Raffle          *domain.Raffle
	Organizer       *domain.User
	PaymentsSummary *PaymentsSummary
	Timeline        []*SettlementEvent
	BankAccount     *OrganizerBankAccount
}

// PaymentsSummary resumen de pagos de la rifa
type PaymentsSummary struct {
	TotalPayments      int
	SucceededPayments  int
	RefundedPayments   int
	TotalRevenue       float64
	TotalRefunded      float64
	NetRevenue         float64
	PlatformFeePercent float64
	PlatformFeeAmount  float64
}

// SettlementEvent evento en timeline de liquidación
type SettlementEvent struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Actor     *string                `json:"actor,omitempty"`
	Details   string                 `json:"details"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// OrganizerBankAccount información bancaria del organizador
type OrganizerBankAccount struct {
	AccountHolder string `json:"account_holder"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
	IBAN          string `json:"iban,omitempty"`
	SWIFT         string `json:"swift,omitempty"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
}

// ViewSettlementDetailsUseCase caso de uso para ver detalles completos de liquidación
type ViewSettlementDetailsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewViewSettlementDetailsUseCase crea una nueva instancia
func NewViewSettlementDetailsUseCase(db *gorm.DB, log *logger.Logger) *ViewSettlementDetailsUseCase {
	return &ViewSettlementDetailsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ViewSettlementDetailsUseCase) Execute(ctx context.Context, settlementID int64, adminID int64) (*SettlementFullDetails, error) {
	// Obtener settlement con detalles
	var settlement SettlementWithDetails
	if err := uc.db.Table("settlements").
		Select(`settlements.*,
			raffles.title as raffle_title,
			COALESCE(users.first_name || ' ' || users.last_name, users.email) as organizer_name,
			users.email as organizer_email,
			users.kyc_level as organizer_kyc_level`).
		Joins("LEFT JOIN raffles ON raffles.id = settlements.raffle_id").
		Joins("LEFT JOIN users ON users.id = settlements.organizer_id").
		Where("settlements.id = ?", settlementID).
		Scan(&settlement).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("SETTLEMENT_NOT_FOUND", "settlement not found", 404, nil)
		}
		uc.log.Error("Error finding settlement", logger.Int64("settlement_id", settlementID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	details := &SettlementFullDetails{
		Settlement: &settlement,
		Timeline:   make([]*SettlementEvent, 0),
	}

	// Obtener raffle completa
	var raffle domain.Raffle
	if err := uc.db.Where("id = ?", settlement.RaffleID).First(&raffle).Error; err == nil {
		details.Raffle = &raffle
	}

	// Obtener organizador completo
	var organizer domain.User
	if err := uc.db.Where("id = ?", settlement.OrganizerID).First(&organizer).Error; err == nil {
		details.Organizer = &organizer
	}

	// Obtener resumen de pagos
	var paymentStats struct {
		TotalCount     int
		SucceededCount int
		RefundedCount  int
		TotalRevenue   float64
		TotalRefunded  float64
	}

	uc.db.Table("payments").
		Select(`
			COUNT(*) as total_count,
			COUNT(CASE WHEN status = 'succeeded' THEN 1 END) as succeeded_count,
			COUNT(CASE WHEN status = 'refunded' THEN 1 END) as refunded_count,
			COALESCE(SUM(CASE WHEN status = 'succeeded' THEN amount ELSE 0 END), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN status = 'refunded' THEN amount ELSE 0 END), 0) as total_refunded
		`).
		Where("raffle_id = (SELECT uuid FROM raffles WHERE id = ?)", settlement.RaffleID).
		Scan(&paymentStats)

	platformFeePercent := 10.0 // TODO: Obtener de configuración
	platformFeeAmount := paymentStats.TotalRevenue * platformFeePercent / 100.0
	netRevenue := paymentStats.TotalRevenue - platformFeeAmount - paymentStats.TotalRefunded

	details.PaymentsSummary = &PaymentsSummary{
		TotalPayments:      paymentStats.TotalCount,
		SucceededPayments:  paymentStats.SucceededCount,
		RefundedPayments:   paymentStats.RefundedCount,
		TotalRevenue:       paymentStats.TotalRevenue,
		TotalRefunded:      paymentStats.TotalRefunded,
		NetRevenue:         netRevenue,
		PlatformFeePercent: platformFeePercent,
		PlatformFeeAmount:  platformFeeAmount,
	}

	// Obtener cuenta bancaria del organizador
	var bankAccount OrganizerBankAccount
	if err := uc.db.Table("organizer_bank_accounts").
		Where("user_id = ? AND is_primary = true", settlement.OrganizerID).
		Scan(&bankAccount).Error; err == nil {
		details.BankAccount = &bankAccount
	}

	// Construir timeline
	// 1. Evento de creación/cálculo
	details.Timeline = append(details.Timeline, &SettlementEvent{
		Type:      "calculated",
		Timestamp: settlement.CalculatedAt,
		Details:   "Settlement calculated automatically",
		Metadata: map[string]interface{}{
			"total_revenue": settlement.TotalRevenue,
			"platform_fee":  settlement.PlatformFee,
			"net_amount":    settlement.NetAmount,
		},
	})

	// 2. Evento de aprobación
	if settlement.ApprovedAt != nil {
		var approverName string
		if settlement.ApprovedBy != nil {
			var approver domain.User
			if err := uc.db.Where("id = ?", *settlement.ApprovedBy).First(&approver).Error; err == nil {
				approverName = approver.GetFullName()
			}
		}

		details.Timeline = append(details.Timeline, &SettlementEvent{
			Type:      "approved",
			Timestamp: *settlement.ApprovedAt,
			Actor:     &approverName,
			Details:   "Settlement approved by admin",
			Metadata: map[string]interface{}{
				"admin_id": settlement.ApprovedBy,
			},
		})
	}

	// 3. Evento de rechazo
	if settlement.RejectedAt != nil {
		var rejecterName string
		if settlement.RejectedBy != nil {
			var rejecter domain.User
			if err := uc.db.Where("id = ?", *settlement.RejectedBy).First(&rejecter).Error; err == nil {
				rejecterName = rejecter.GetFullName()
			}
		}

		details.Timeline = append(details.Timeline, &SettlementEvent{
			Type:      "rejected",
			Timestamp: *settlement.RejectedAt,
			Actor:     &rejecterName,
			Details:   "Settlement rejected by admin",
			Metadata: map[string]interface{}{
				"admin_id": settlement.RejectedBy,
				"reason":   settlement.RejectionReason,
			},
		})
	}

	// 4. Evento de pago
	if settlement.PaidAt != nil {
		details.Timeline = append(details.Timeline, &SettlementEvent{
			Type:      "paid",
			Timestamp: *settlement.PaidAt,
			Details:   "Settlement marked as paid",
			Metadata: map[string]interface{}{
				"payment_reference": settlement.PaymentReference,
				"payment_method":    settlement.PaymentMethod,
			},
		})
	}

	// Ordenar timeline cronológicamente (más antiguo primero)
	// Ya está ordenado por construcción, pero podríamos usar sort.Slice si agregamos más eventos

	// Log auditoría
	uc.log.Info("Admin viewed settlement details",
		logger.Int64("admin_id", adminID),
		logger.Int64("settlement_id", settlementID),
		logger.Int64("organizer_id", settlement.OrganizerID),
		logger.String("action", "admin_view_settlement_details"))

	return details, nil
}
