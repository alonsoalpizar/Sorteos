package organizer

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// CalculateOrganizerRevenueInput datos de entrada
type CalculateOrganizerRevenueInput struct {
	OrganizerID int64
	DateFrom    *string // Formato: YYYY-MM-DD
	DateTo      *string // Formato: YYYY-MM-DD
	GroupBy     *string // "month", "year", null para total
}

// CalculateOrganizerRevenueOutput resultado
type CalculateOrganizerRevenueOutput struct {
	OrganizerID    int64                `json:"organizer_id"`
	OrganizerEmail string               `json:"organizer_email"`
	OrganizerName  string               `json:"organizer_name"`
	TotalRevenue   *RevenueBreakdown    `json:"total_revenue"`
	ByPeriod       []*PeriodRevenue     `json:"by_period,omitempty"` // Si GroupBy está presente
	DateFrom       string               `json:"date_from"`
	DateTo         string               `json:"date_to"`
}

// RevenueBreakdown desglose de revenue
type RevenueBreakdown struct {
	GrossRevenue    float64 `json:"gross_revenue"`     // Total vendido
	PlatformFees    float64 `json:"platform_fees"`     // Comisión de plataforma
	NetRevenue      float64 `json:"net_revenue"`       // Lo que le corresponde al organizador
	PendingPayout   float64 `json:"pending_payout"`    // Pendiente de pagar
	PaidOut         float64 `json:"paid_out"`          // Ya pagado
	TotalRaffles    int     `json:"total_raffles"`     // Número de rifas
	CompletedRaffles int    `json:"completed_raffles"` // Rifas completadas
}

// PeriodRevenue revenue por periodo
type PeriodRevenue struct {
	Period   string            `json:"period"`   // "2025-11" o "2025"
	Revenue  *RevenueBreakdown `json:"revenue"`
}

// CalculateOrganizerRevenueUseCase caso de uso para calcular revenue de organizador
type CalculateOrganizerRevenueUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCalculateOrganizerRevenueUseCase crea una nueva instancia
func NewCalculateOrganizerRevenueUseCase(db *gorm.DB, log *logger.Logger) *CalculateOrganizerRevenueUseCase {
	return &CalculateOrganizerRevenueUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *CalculateOrganizerRevenueUseCase) Execute(ctx context.Context, input *CalculateOrganizerRevenueInput, adminID int64) (*CalculateOrganizerRevenueOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Verificar que el organizador existe
	var organizer struct {
		ID    int64
		Email string
		FirstName *string
		LastName *string
	}
	result := uc.db.WithContext(ctx).
		Table("users").
		Select("id, email, first_name, last_name").
		Where("id = ? AND role = ?", input.OrganizerID, "organizer").
		First(&organizer)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("ORGANIZER_NOT_FOUND", "organizer not found", 404, nil)
		}
		uc.log.Error("Error finding organizer", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Construir nombre
	organizerName := organizer.Email
	if organizer.FirstName != nil && organizer.LastName != nil {
		organizerName = *organizer.FirstName + " " + *organizer.LastName
	} else if organizer.FirstName != nil {
		organizerName = *organizer.FirstName
	}

	// Determinar rango de fechas
	dateFrom := input.DateFrom
	dateTo := input.DateTo
	if dateFrom == nil {
		// Por defecto: hace 1 año
		oneYearAgo := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		dateFrom = &oneYearAgo
	}
	if dateTo == nil {
		// Por defecto: hoy
		today := time.Now().Format("2006-01-02")
		dateTo = &today
	}

	// Calcular revenue total
	totalRevenue, err := uc.calculateRevenueForPeriod(ctx, input.OrganizerID, *dateFrom, *dateTo)
	if err != nil {
		return nil, err
	}

	// Construir output base
	output := &CalculateOrganizerRevenueOutput{
		OrganizerID:    input.OrganizerID,
		OrganizerEmail: organizer.Email,
		OrganizerName:  organizerName,
		TotalRevenue:   totalRevenue,
		DateFrom:       *dateFrom,
		DateTo:         *dateTo,
	}

	// Si se requiere agrupación por periodo
	if input.GroupBy != nil && *input.GroupBy != "" {
		byPeriod, err := uc.calculateRevenueByPeriod(ctx, input.OrganizerID, *dateFrom, *dateTo, *input.GroupBy)
		if err != nil {
			return nil, err
		}
		output.ByPeriod = byPeriod
	}

	// Log auditoría
	uc.log.Info("Admin calculated organizer revenue",
		logger.Int64("admin_id", adminID),
		logger.Int64("organizer_id", input.OrganizerID),
		logger.String("date_from", *dateFrom),
		logger.String("date_to", *dateTo),
		logger.Float64("gross_revenue", totalRevenue.GrossRevenue),
		logger.String("action", "admin_calculate_organizer_revenue"))

	return output, nil
}

// validateInput valida los datos de entrada
func (uc *CalculateOrganizerRevenueUseCase) validateInput(input *CalculateOrganizerRevenueInput) error {
	if input.OrganizerID <= 0 {
		return errors.New("VALIDATION_FAILED", "organizer_id is required", 400, nil)
	}

	// Validar formato de fechas si están presentes
	if input.DateFrom != nil {
		if _, err := time.Parse("2006-01-02", *input.DateFrom); err != nil {
			return errors.New("VALIDATION_FAILED", "date_from must be in format YYYY-MM-DD", 400, err)
		}
	}
	if input.DateTo != nil {
		if _, err := time.Parse("2006-01-02", *input.DateTo); err != nil {
			return errors.New("VALIDATION_FAILED", "date_to must be in format YYYY-MM-DD", 400, err)
		}
	}

	// Validar GroupBy
	if input.GroupBy != nil && *input.GroupBy != "" {
		validGroupBy := map[string]bool{
			"month": true,
			"year":  true,
		}
		if !validGroupBy[*input.GroupBy] {
			return errors.New("VALIDATION_FAILED", "group_by must be 'month' or 'year'", 400, nil)
		}
	}

	return nil
}

// calculateRevenueForPeriod calcula el revenue para un periodo específico
func (uc *CalculateOrganizerRevenueUseCase) calculateRevenueForPeriod(ctx context.Context, organizerID int64, dateFrom, dateTo string) (*RevenueBreakdown, error) {
	breakdown := &RevenueBreakdown{}

	// Query para calcular revenue desde raffles completadas
	query := `
		SELECT
			COUNT(*) as total_raffles,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_raffles,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN price_per_number * sold_count ELSE 0 END), 0) as gross_revenue
		FROM raffles
		WHERE user_id = ?
			AND deleted_at IS NULL
			AND created_at >= ?
			AND created_at <= ?
	`

	var totalRaffles, completedRaffles int
	var grossRevenue float64

	err := uc.db.WithContext(ctx).Raw(query, organizerID, dateFrom, dateTo+" 23:59:59").
		Row().
		Scan(&totalRaffles, &completedRaffles, &grossRevenue)

	if err != nil {
		uc.log.Error("Error calculating raffle revenue", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	breakdown.TotalRaffles = totalRaffles
	breakdown.CompletedRaffles = completedRaffles
	breakdown.GrossRevenue = grossRevenue

	// Calcular platform fees
	// Primero intentar obtener custom commission del organizer
	var customCommission *float64
	uc.db.WithContext(ctx).
		Table("organizer_profiles").
		Select("commission_override").
		Where("user_id = ?", organizerID).
		Scan(&customCommission)

	platformFeePercent := 10.0 // Default 10%
	if customCommission != nil {
		platformFeePercent = *customCommission
	}

	breakdown.PlatformFees = grossRevenue * (platformFeePercent / 100.0)
	breakdown.NetRevenue = grossRevenue - breakdown.PlatformFees

	// Calcular pending payout y paid out desde settlements
	var paidOut, pendingPayout float64

	// Paid out: settlements con status 'paid'
	uc.db.WithContext(ctx).
		Table("settlements").
		Select("COALESCE(SUM(net_amount), 0)").
		Where("organizer_id = ? AND status = ? AND calculated_at >= ? AND calculated_at <= ?",
			organizerID, "paid", dateFrom, dateTo+" 23:59:59").
		Scan(&paidOut)

	// Pending payout: settlements con status 'pending' o 'approved'
	uc.db.WithContext(ctx).
		Table("settlements").
		Select("COALESCE(SUM(net_amount), 0)").
		Where("organizer_id = ? AND status IN (?, ?) AND calculated_at >= ? AND calculated_at <= ?",
			organizerID, "pending", "approved", dateFrom, dateTo+" 23:59:59").
		Scan(&pendingPayout)

	breakdown.PaidOut = paidOut
	breakdown.PendingPayout = pendingPayout

	return breakdown, nil
}

// calculateRevenueByPeriod calcula revenue agrupado por periodo
func (uc *CalculateOrganizerRevenueUseCase) calculateRevenueByPeriod(ctx context.Context, organizerID int64, dateFrom, dateTo, groupBy string) ([]*PeriodRevenue, error) {
	var periods []*PeriodRevenue

	// Obtener periodos únicos
	var periodStrings []string
	query := `
		SELECT DISTINCT TO_CHAR(created_at, ?) as period
		FROM raffles
		WHERE user_id = ?
			AND deleted_at IS NULL
			AND created_at >= ?
			AND created_at <= ?
		ORDER BY period
	`

	// Convertir formato Go a PostgreSQL
	pgFormat := "YYYY-MM"
	if groupBy == "year" {
		pgFormat = "YYYY"
	}

	rows, err := uc.db.WithContext(ctx).Raw(query, pgFormat, organizerID, dateFrom, dateTo+" 23:59:59").Rows()
	if err != nil {
		uc.log.Error("Error getting revenue periods", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var period string
		rows.Scan(&period)
		periodStrings = append(periodStrings, period)
	}

	// Calcular revenue para cada periodo
	for _, period := range periodStrings {
		// Determinar rango de fechas para este periodo
		var periodStart, periodEnd string

		if groupBy == "month" {
			// Parse "2025-11"
			t, _ := time.Parse("2006-01", period)
			periodStart = t.Format("2006-01-02")
			periodEnd = t.AddDate(0, 1, -1).Format("2006-01-02") // Último día del mes
		} else {
			// Parse "2025"
			t, _ := time.Parse("2006", period)
			periodStart = t.Format("2006-01-02")
			periodEnd = t.AddDate(1, 0, -1).Format("2006-01-02") // Último día del año
		}

		// Calcular revenue para este periodo
		revenue, err := uc.calculateRevenueForPeriod(ctx, organizerID, periodStart, periodEnd)
		if err != nil {
			return nil, err
		}

		periods = append(periods, &PeriodRevenue{
			Period:  period,
			Revenue: revenue,
		})
	}

	return periods, nil
}
