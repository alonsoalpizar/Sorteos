package reports

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// OrganizerPayoutRow fila del reporte de pagos a organizadores
type OrganizerPayoutRow struct {
	OrganizerID        int64   `json:"organizer_id"`
	OrganizerName      string  `json:"organizer_name"`
	OrganizerEmail     string  `json:"organizer_email"`
	KYCLevel           string  `json:"kyc_level"`
	TotalRaffles       int64   `json:"total_raffles"`
	CompletedRaffles   int64   `json:"completed_raffles"`
	TotalRevenue       float64 `json:"total_revenue"`
	TotalPlatformFees  float64 `json:"total_platform_fees"`
	TotalPayouts       float64 `json:"total_payouts"`        // Settlements paid
	PendingPayout      float64 `json:"pending_payout"`       // Settlements pending+approved
	CustomCommission   *float64 `json:"custom_commission,omitempty"`
	AverageRevenuePerRaffle float64 `json:"average_revenue_per_raffle"`
}

// OrganizerPayoutsReportInput datos de entrada
type OrganizerPayoutsReportInput struct {
	DateFrom    string
	DateTo      string
	VerifiedOnly bool   // Solo organizadores verificados
	MinRevenue   *float64
	OrderBy      string // total_revenue DESC, total_raffles DESC, etc.
	Page         int
	PageSize     int
}

// OrganizerPayoutsReportOutput resultado
type OrganizerPayoutsReportOutput struct {
	Rows              []*OrganizerPayoutRow
	Total             int64
	Page              int
	PageSize          int
	TotalPages        int
	// Totales globales
	TotalRevenue      float64
	TotalPlatformFees float64
	TotalPayouts      float64
	TotalPending      float64
}

// OrganizerPayoutsReportUseCase caso de uso para reporte de pagos a organizadores
type OrganizerPayoutsReportUseCase struct {
	systemParamRepo *db.PostgresSystemParameterRepository
	db  *gorm.DB
	log *logger.Logger
}

// NewOrganizerPayoutsReportUseCase crea una nueva instancia
func NewOrganizerPayoutsReportUseCase(gormDB *gorm.DB, log *logger.Logger) *OrganizerPayoutsReportUseCase {
	return &OrganizerPayoutsReportUseCase{
		db:              gormDB,
		systemParamRepo: db.NewSystemParameterRepository(gormDB, log),
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *OrganizerPayoutsReportUseCase) Execute(ctx context.Context, input *OrganizerPayoutsReportInput, adminID int64) (*OrganizerPayoutsReportOutput, error) {
	// Validar fechas
	if input.DateFrom == "" || input.DateTo == "" {
		return nil, errors.New("VALIDATION_FAILED", "date_from and date_to are required", 400, nil)
	}

	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	offset := (input.Page - 1) * input.PageSize

	// Query base: organizadores con sus rifas
	baseQuery := uc.db.Table("users").
		Where("role = ?", "organizer")

	if input.VerifiedOnly {
		baseQuery = baseQuery.Where("kyc_level IN (?)", []string{"verified", "enhanced"})
	}

	// Obtener organizadores
	type OrganizerRow struct {
		ID               int64
		FirstName        *string
		LastName         *string
		Email            string
		KYCLevel         string
		CustomCommission *float64
	}

	var organizers []OrganizerRow
	if err := baseQuery.
		Select("id, first_name, last_name, email, kyc_level").
		Scan(&organizers).Error; err != nil {
		uc.log.Error("Error fetching organizers", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Para cada organizador, calcular métricas
	rows := make([]*OrganizerPayoutRow, 0)

	for _, org := range organizers {
		// Nombre del organizador
		organizerName := org.Email
		if org.FirstName != nil && org.LastName != nil {
			organizerName = *org.FirstName + " " + *org.LastName
		}

		// Contar rifas totales y completadas en el período
		var totalRaffles int64
		var completedRaffles int64

		uc.db.Table("raffles").
			Where("user_id = ?", org.ID).
			Where("created_at >= ?", input.DateFrom).
			Where("created_at <= ?", input.DateTo+" 23:59:59").
			Count(&totalRaffles)

		uc.db.Table("raffles").
			Where("user_id = ?", org.ID).
			Where("status = ?", "completed").
			Where("completed_at >= ?", input.DateFrom).
			Where("completed_at <= ?", input.DateTo+" 23:59:59").
			Count(&completedRaffles)

		// Calcular revenue total de rifas completadas
		var totalRevenue float64
		uc.db.Table("payments").
			Joins("JOIN raffles ON raffles.uuid::text = payments.raffle_id").
			Where("raffles.user_id = ?", org.ID).
			Where("raffles.status = ?", "completed").
			Where("raffles.completed_at >= ?", input.DateFrom).
			Where("raffles.completed_at <= ?", input.DateTo+" 23:59:59").
			Where("payments.status = ?", "succeeded").
			Select("COALESCE(SUM(payments.amount), 0)").
			Scan(&totalRevenue)

		// Si no cumple con min_revenue, skip
		if input.MinRevenue != nil && totalRevenue < *input.MinRevenue {
			continue
		}

		// Platform fees (TODO: considerar custom commission)
		// Obtener platform_fee_percentage desde system_parameters
		platformFeePercent, _ := uc.systemParamRepo.GetFloat("platform_fee_percentage", 10.0)
		if org.CustomCommission != nil {
			platformFeePercent = *org.CustomCommission
		}
		totalPlatformFees := totalRevenue * platformFeePercent / 100.0

		// Settlements pagados (paid)
		var totalPayouts float64
		uc.db.Table("settlements").
			Where("organizer_id = ?", org.ID).
			Where("status = ?", "paid").
			Where("paid_at >= ?", input.DateFrom).
			Where("paid_at <= ?", input.DateTo+" 23:59:59").
			Select("COALESCE(SUM(net_amount), 0)").
			Scan(&totalPayouts)

		// Settlements pendientes (pending + approved)
		var pendingPayout float64
		uc.db.Table("settlements").
			Where("organizer_id = ?", org.ID).
			Where("status IN (?)", []string{"pending", "approved"}).
			Where("calculated_at >= ?", input.DateFrom).
			Where("calculated_at <= ?", input.DateTo+" 23:59:59").
			Select("COALESCE(SUM(net_amount), 0)").
			Scan(&pendingPayout)

		// Average revenue per raffle
		averageRevenuePerRaffle := 0.0
		if completedRaffles > 0 {
			averageRevenuePerRaffle = totalRevenue / float64(completedRaffles)
		}

		row := &OrganizerPayoutRow{
			OrganizerID:             org.ID,
			OrganizerName:           organizerName,
			OrganizerEmail:          org.Email,
			KYCLevel:                org.KYCLevel,
			TotalRaffles:            totalRaffles,
			CompletedRaffles:        completedRaffles,
			TotalRevenue:            totalRevenue,
			TotalPlatformFees:       totalPlatformFees,
			TotalPayouts:            totalPayouts,
			PendingPayout:           pendingPayout,
			CustomCommission:        org.CustomCommission,
			AverageRevenuePerRaffle: averageRevenuePerRaffle,
		}

		rows = append(rows, row)
	}

	// Ordenar rows
	orderBy := "total_revenue DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}

	// Implementar ordenamiento manual (TODO: mejorar con ORDER BY en query)
	// Por ahora ordenamos por revenue descendente
	for i := 0; i < len(rows)-1; i++ {
		for j := i + 1; j < len(rows); j++ {
			if orderBy == "total_revenue DESC" && rows[i].TotalRevenue < rows[j].TotalRevenue {
				rows[i], rows[j] = rows[j], rows[i]
			}
		}
	}

	// Calcular totales globales
	totalRevenue := 0.0
	totalPlatformFees := 0.0
	totalPayouts := 0.0
	totalPending := 0.0

	for _, row := range rows {
		totalRevenue += row.TotalRevenue
		totalPlatformFees += row.TotalPlatformFees
		totalPayouts += row.TotalPayouts
		totalPending += row.PendingPayout
	}

	// Paginación
	total := int64(len(rows))
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	// Aplicar paginación a rows
	start := offset
	end := offset + input.PageSize
	if start > len(rows) {
		start = len(rows)
	}
	if end > len(rows) {
		end = len(rows)
	}
	paginatedRows := rows[start:end]

	// Log auditoría
	uc.log.Info("Admin generated organizer payouts report",
		logger.Int64("admin_id", adminID),
		logger.String("date_from", input.DateFrom),
		logger.String("date_to", input.DateTo),
		logger.Int("total_organizers", len(rows)),
		logger.Float64("total_revenue", totalRevenue),
		logger.String("action", "admin_organizer_payouts_report"))

	return &OrganizerPayoutsReportOutput{
		Rows:              paginatedRows,
		Total:             total,
		Page:              input.Page,
		PageSize:          input.PageSize,
		TotalPages:        totalPages,
		TotalRevenue:      totalRevenue,
		TotalPlatformFees: totalPlatformFees,
		TotalPayouts:      totalPayouts,
		TotalPending:      totalPending,
	}, nil
}
