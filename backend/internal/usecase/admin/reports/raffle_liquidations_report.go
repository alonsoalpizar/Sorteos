package reports

import (
	"context"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// RaffleLiquidationRow fila del reporte de liquidaciones
type RaffleLiquidationRow struct {
	RaffleID           int64   `json:"raffle_id"`
	RaffleTitle        string  `json:"raffle_title"`
	OrganizerID        int64   `json:"organizer_id"`
	OrganizerName      string  `json:"organizer_name"`
	OrganizerEmail     string  `json:"organizer_email"`
	CompletedAt        string  `json:"completed_at"`
	GrossRevenue       float64 `json:"gross_revenue"`
	PlatformFeePercent float64 `json:"platform_fee_percent"`
	PlatformFee        float64 `json:"platform_fee"`
	NetRevenue         float64 `json:"net_revenue"`
	SettlementID       *int64  `json:"settlement_id,omitempty"`
	SettlementStatus   *string `json:"settlement_status,omitempty"`
	PaidAt             *string `json:"paid_at,omitempty"`
}

// RaffleLiquidationsReportInput datos de entrada
type RaffleLiquidationsReportInput struct {
	DateFrom    string
	DateTo      string
	OrganizerID *int64
	CategoryID  *int64
	SettlementStatus *string // pending, approved, paid, rejected, null (sin settlement)
	OrderBy     string
}

// RaffleLiquidationsReportOutput resultado
type RaffleLiquidationsReportOutput struct {
	Rows              []*RaffleLiquidationRow
	Total             int64
	TotalGrossRevenue float64
	TotalPlatformFees float64
	TotalNetRevenue   float64
	// Breakdown por settlement status
	WithSettlement    int64
	WithoutSettlement int64
	PendingCount      int64
	ApprovedCount     int64
	PaidCount         int64
	RejectedCount     int64
}

// RaffleLiquidationsReportUseCase caso de uso para reporte de liquidaciones
type RaffleLiquidationsReportUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewRaffleLiquidationsReportUseCase crea una nueva instancia
func NewRaffleLiquidationsReportUseCase(db *gorm.DB, log *logger.Logger) *RaffleLiquidationsReportUseCase {
	return &RaffleLiquidationsReportUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *RaffleLiquidationsReportUseCase) Execute(ctx context.Context, input *RaffleLiquidationsReportInput, adminID int64) (*RaffleLiquidationsReportOutput, error) {
	// Validar fechas
	if input.DateFrom == "" || input.DateTo == "" {
		return nil, errors.New("VALIDATION_FAILED", "date_from and date_to are required", 400, nil)
	}

	// Query base: rifas completadas con JOIN a users (organizadores) y settlements
	query := uc.db.Table("raffles").
		Select(`
			raffles.id as raffle_id,
			raffles.title as raffle_title,
			raffles.user_id as organizer_id,
			COALESCE(users.first_name || ' ' || users.last_name, users.email) as organizer_name,
			users.email as organizer_email,
			raffles.completed_at,
			settlements.id as settlement_id,
			settlements.status as settlement_status,
			settlements.paid_at
		`).
		Joins("LEFT JOIN users ON users.id = raffles.user_id").
		Joins("LEFT JOIN settlements ON settlements.raffle_id = raffles.id").
		Where("raffles.status = ?", "completed").
		Where("raffles.completed_at >= ?", input.DateFrom).
		Where("raffles.completed_at <= ?", input.DateTo+" 23:59:59")

	// Aplicar filtros opcionales
	if input.OrganizerID != nil {
		query = query.Where("raffles.user_id = ?", *input.OrganizerID)
	}

	if input.CategoryID != nil {
		query = query.Where("raffles.category_id = ?", *input.CategoryID)
	}

	if input.SettlementStatus != nil {
		if *input.SettlementStatus == "null" {
			query = query.Where("settlements.id IS NULL")
		} else {
			query = query.Where("settlements.status = ?", *input.SettlementStatus)
		}
	}

	// Ordenamiento
	orderBy := "raffles.completed_at DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}
	query = query.Order(orderBy)

	// Ejecutar query
	type RaffleRow struct {
		RaffleID         int64
		RaffleTitle      string
		OrganizerID      int64
		OrganizerName    string
		OrganizerEmail   string
		CompletedAt      string
		SettlementID     *int64
		SettlementStatus *string
		PaidAt           *string
	}

	var raffleRows []RaffleRow
	if err := query.Scan(&raffleRows).Error; err != nil {
		uc.log.Error("Error generating liquidations report", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Para cada rifa, calcular revenue desde payments
	rows := make([]*RaffleLiquidationRow, 0, len(raffleRows))
	totalGrossRevenue := 0.0
	totalPlatformFees := 0.0

	withSettlement := int64(0)
	withoutSettlement := int64(0)
	pendingCount := int64(0)
	approvedCount := int64(0)
	paidCount := int64(0)
	rejectedCount := int64(0)

	for _, raffleRow := range raffleRows {
		// Calcular gross revenue de esta rifa
		var grossRevenue float64
		uc.db.Table("payments").
			Select("COALESCE(SUM(amount), 0)").
			Where("raffle_id = (SELECT uuid FROM raffles WHERE id = ?)", raffleRow.RaffleID).
			Where("status = ?", "succeeded").
			Scan(&grossRevenue)

		// Platform fee (TODO: considerar custom commission de organizer)
		platformFeePercent := 10.0
		platformFee := grossRevenue * platformFeePercent / 100.0
		netRevenue := grossRevenue - platformFee

		row := &RaffleLiquidationRow{
			RaffleID:           raffleRow.RaffleID,
			RaffleTitle:        raffleRow.RaffleTitle,
			OrganizerID:        raffleRow.OrganizerID,
			OrganizerName:      raffleRow.OrganizerName,
			OrganizerEmail:     raffleRow.OrganizerEmail,
			CompletedAt:        raffleRow.CompletedAt,
			GrossRevenue:       grossRevenue,
			PlatformFeePercent: platformFeePercent,
			PlatformFee:        platformFee,
			NetRevenue:         netRevenue,
			SettlementID:       raffleRow.SettlementID,
			SettlementStatus:   raffleRow.SettlementStatus,
			PaidAt:             raffleRow.PaidAt,
		}

		rows = append(rows, row)

		totalGrossRevenue += grossRevenue
		totalPlatformFees += platformFee

		// Contar settlements
		if raffleRow.SettlementID != nil {
			withSettlement++
			if raffleRow.SettlementStatus != nil {
				switch *raffleRow.SettlementStatus {
				case "pending":
					pendingCount++
				case "approved":
					approvedCount++
				case "paid":
					paidCount++
				case "rejected":
					rejectedCount++
				}
			}
		} else {
			withoutSettlement++
		}
	}

	totalNetRevenue := totalGrossRevenue - totalPlatformFees

	// Log auditoría
	uc.log.Info("Admin generated liquidations report",
		logger.Int64("admin_id", adminID),
		logger.String("date_from", input.DateFrom),
		logger.String("date_to", input.DateTo),
		logger.Int("total_raffles", len(rows)),
		logger.Float64("total_revenue", totalGrossRevenue),
		logger.String("action", "admin_liquidations_report"))

	return &RaffleLiquidationsReportOutput{
		Rows:              rows,
		Total:             int64(len(rows)),
		TotalGrossRevenue: totalGrossRevenue,
		TotalPlatformFees: totalPlatformFees,
		TotalNetRevenue:   totalNetRevenue,
		WithSettlement:    withSettlement,
		WithoutSettlement: withoutSettlement,
		PendingCount:      pendingCount,
		ApprovedCount:     approvedCount,
		PaidCount:         paidCount,
		RejectedCount:     rejectedCount,
	}, nil
}
