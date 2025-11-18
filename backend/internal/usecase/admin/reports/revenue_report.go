package reports

import (
	"context"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// RevenueDataPoint punto de datos para gráficos de tiempo
type RevenueDataPoint struct {
	Date         string  `json:"date"`           // YYYY-MM-DD o YYYY-MM o YYYY
	GrossRevenue float64 `json:"gross_revenue"`  // Total revenue de pagos
	PlatformFees float64 `json:"platform_fees"`  // Comisiones de plataforma
	NetRevenue   float64 `json:"net_revenue"`    // Revenue neto a organizadores
	PaymentCount int64   `json:"payment_count"`  // Número de pagos
	RaffleCount  int64   `json:"raffle_count"`   // Número de rifas completadas
}

// RevenueReportInput datos de entrada
type RevenueReportInput struct {
	DateFrom    string  // YYYY-MM-DD
	DateTo      string  // YYYY-MM-DD
	OrganizerID *int64  // Filtrar por organizador específico
	CategoryID  *int64  // Filtrar por categoría
	GroupBy     string  // day, week, month (default: day)
}

// RevenueReportOutput resultado
type RevenueReportOutput struct {
	DataPoints       []*RevenueDataPoint
	TotalGrossRevenue float64
	TotalPlatformFees float64
	TotalNetRevenue   float64
	TotalPayments     int64
	TotalRaffles      int64
	AverageRevenuePerDay float64
	AverageRevenuePerRaffle float64
}

// RevenueReportUseCase caso de uso para generar reporte de ingresos
type RevenueReportUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewRevenueReportUseCase crea una nueva instancia
func NewRevenueReportUseCase(db *gorm.DB, log *logger.Logger) *RevenueReportUseCase {
	return &RevenueReportUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *RevenueReportUseCase) Execute(ctx context.Context, input *RevenueReportInput, adminID int64) (*RevenueReportOutput, error) {
	// Validar fechas
	if input.DateFrom == "" || input.DateTo == "" {
		return nil, errors.New("VALIDATION_FAILED", "date_from and date_to are required", 400, nil)
	}

	// Validar group_by
	if input.GroupBy == "" {
		input.GroupBy = "day"
	}
	if input.GroupBy != "day" && input.GroupBy != "week" && input.GroupBy != "month" {
		return nil, errors.New("VALIDATION_FAILED", "group_by must be day, week, or month", 400, nil)
	}

	// Determinar formato de agrupación SQL
	var dateFormat string
	switch input.GroupBy {
	case "day":
		dateFormat = "TO_CHAR(paid_at, 'YYYY-MM-DD')"
	case "week":
		dateFormat = "TO_CHAR(DATE_TRUNC('week', paid_at), 'YYYY-MM-DD')"
	case "month":
		dateFormat = "TO_CHAR(paid_at, 'YYYY-MM')"
	}

	// Construir query base
	query := uc.db.Table("payments").
		Select(`
			`+dateFormat+` as date,
			COALESCE(SUM(amount), 0) as gross_revenue,
			COUNT(*) as payment_count
		`).
		Where("status = ?", "succeeded").
		Where("paid_at >= ?", input.DateFrom).
		Where("paid_at <= ?", input.DateTo+" 23:59:59").
		Group("date").
		Order("date ASC")

	// Aplicar filtros opcionales
	if input.OrganizerID != nil {
		// JOIN con raffles para filtrar por organizador
		query = query.
			Joins("JOIN raffles ON raffles.uuid::text = payments.raffle_id").
			Where("raffles.user_id = ?", *input.OrganizerID)
	}

	if input.CategoryID != nil {
		// JOIN con raffles para filtrar por categoría
		if input.OrganizerID == nil {
			query = query.Joins("JOIN raffles ON raffles.uuid::text = payments.raffle_id")
		}
		query = query.Where("raffles.category_id = ?", *input.CategoryID)
	}

	// Ejecutar query
	var rawDataPoints []struct {
		Date         string
		GrossRevenue float64
		PaymentCount int64
	}

	if err := query.Scan(&rawDataPoints).Error; err != nil {
		uc.log.Error("Error generating revenue report", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Platform fee percent (TODO: obtener de configuración)
	platformFeePercent := 0.10

	// Convertir a data points con cálculos
	dataPoints := make([]*RevenueDataPoint, len(rawDataPoints))
	totalGrossRevenue := 0.0
	totalPlatformFees := 0.0
	totalPayments := int64(0)

	for i, raw := range rawDataPoints {
		platformFee := raw.GrossRevenue * platformFeePercent
		netRevenue := raw.GrossRevenue - platformFee

		dataPoints[i] = &RevenueDataPoint{
			Date:         raw.Date,
			GrossRevenue: raw.GrossRevenue,
			PlatformFees: platformFee,
			NetRevenue:   netRevenue,
			PaymentCount: raw.PaymentCount,
		}

		totalGrossRevenue += raw.GrossRevenue
		totalPlatformFees += platformFee
		totalPayments += raw.PaymentCount
	}

	totalNetRevenue := totalGrossRevenue - totalPlatformFees

	// Contar rifas completadas en el período
	var totalRaffles int64
	raffleQuery := uc.db.Table("raffles").
		Where("status = ?", "completed").
		Where("completed_at >= ?", input.DateFrom).
		Where("completed_at <= ?", input.DateTo+" 23:59:59")

	if input.OrganizerID != nil {
		raffleQuery = raffleQuery.Where("user_id = ?", *input.OrganizerID)
	}
	if input.CategoryID != nil {
		raffleQuery = raffleQuery.Where("category_id = ?", *input.CategoryID)
	}

	raffleQuery.Count(&totalRaffles)

	// Agregar raffle_count a cada data point
	for _, dp := range dataPoints {
		// Contar rifas completadas en esta fecha/período
		var raffleCount int64
		dateQuery := uc.db.Table("raffles").
			Where("status = ?", "completed")

		// Ajustar filtro de fecha según group_by
		if input.GroupBy == "day" {
			dateQuery = dateQuery.
				Where("DATE(completed_at) = ?", dp.Date)
		} else if input.GroupBy == "month" {
			dateQuery = dateQuery.
				Where("TO_CHAR(completed_at, 'YYYY-MM') = ?", dp.Date)
		} else if input.GroupBy == "week" {
			dateQuery = dateQuery.
				Where("DATE_TRUNC('week', completed_at)::date = ?", dp.Date)
		}

		if input.OrganizerID != nil {
			dateQuery = dateQuery.Where("user_id = ?", *input.OrganizerID)
		}
		if input.CategoryID != nil {
			dateQuery = dateQuery.Where("category_id = ?", *input.CategoryID)
		}

		dateQuery.Count(&raffleCount)
		dp.RaffleCount = raffleCount
	}

	// Calcular promedios
	averageRevenuePerDay := 0.0
	if len(dataPoints) > 0 {
		averageRevenuePerDay = totalGrossRevenue / float64(len(dataPoints))
	}

	averageRevenuePerRaffle := 0.0
	if totalRaffles > 0 {
		averageRevenuePerRaffle = totalGrossRevenue / float64(totalRaffles)
	}

	// Log auditoría
	uc.log.Info("Admin generated revenue report",
		logger.Int64("admin_id", adminID),
		logger.String("date_from", input.DateFrom),
		logger.String("date_to", input.DateTo),
		logger.String("group_by", input.GroupBy),
		logger.Float64("total_revenue", totalGrossRevenue),
		logger.String("action", "admin_revenue_report"))

	return &RevenueReportOutput{
		DataPoints:              dataPoints,
		TotalGrossRevenue:       totalGrossRevenue,
		TotalPlatformFees:       totalPlatformFees,
		TotalNetRevenue:         totalNetRevenue,
		TotalPayments:           totalPayments,
		TotalRaffles:            totalRaffles,
		AverageRevenuePerDay:    averageRevenuePerDay,
		AverageRevenuePerRaffle: averageRevenuePerRaffle,
	}, nil
}
