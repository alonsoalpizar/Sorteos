package reports

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// DashboardKPIs métricas principales del dashboard
type DashboardKPIs struct {
	// Usuarios
	TotalUsers      int64 `json:"total_users"`
	ActiveUsers     int64 `json:"active_users"`
	SuspendedUsers  int64 `json:"suspended_users"`
	BannedUsers     int64 `json:"banned_users"`
	NewUsersToday   int64 `json:"new_users_today"`
	NewUsersWeek    int64 `json:"new_users_week"`
	NewUsersMonth   int64 `json:"new_users_month"`

	// Organizadores
	TotalOrganizers    int64 `json:"total_organizers"`
	VerifiedOrganizers int64 `json:"verified_organizers"`
	PendingOrganizers  int64 `json:"pending_organizers"`

	// Rifas
	TotalRaffles     int64 `json:"total_raffles"`
	ActiveRaffles    int64 `json:"active_raffles"`
	CompletedRaffles int64 `json:"completed_raffles"`
	SuspendedRaffles int64 `json:"suspended_raffles"`
	DraftRaffles     int64 `json:"draft_raffles"`

	// Revenue
	RevenueToday      float64 `json:"revenue_today"`
	RevenueWeek       float64 `json:"revenue_week"`
	RevenueMonth      float64 `json:"revenue_month"`
	RevenueYear       float64 `json:"revenue_year"`
	RevenueAllTime    float64 `json:"revenue_all_time"`

	// Platform Fees
	PlatformFeesToday    float64 `json:"platform_fees_today"`
	PlatformFeesMonth    float64 `json:"platform_fees_month"`
	PlatformFeesAllTime  float64 `json:"platform_fees_all_time"`

	// Settlements
	PendingSettlementsCount  int64   `json:"pending_settlements_count"`
	PendingSettlementsAmount float64 `json:"pending_settlements_amount"`
	ApprovedSettlementsCount int64   `json:"approved_settlements_count"`
	ApprovedSettlementsAmount float64 `json:"approved_settlements_amount"`

	// Payments
	TotalPayments       int64   `json:"total_payments"`
	SucceededPayments   int64   `json:"succeeded_payments"`
	PendingPayments     int64   `json:"pending_payments"`
	FailedPayments      int64   `json:"failed_payments"`
	RefundedPayments    int64   `json:"refunded_payments"`
	TotalPaymentsAmount float64 `json:"total_payments_amount"`

	// Activity (últimas 24h)
	RecentUsers       int64 `json:"recent_users"`        // Usuarios creados en 24h
	RecentRaffles     int64 `json:"recent_raffles"`      // Rifas creadas en 24h
	RecentPayments    int64 `json:"recent_payments"`     // Pagos en 24h
	RecentSettlements int64 `json:"recent_settlements"`  // Settlements en 24h
}

// GlobalDashboardUseCase caso de uso para obtener KPIs del dashboard
type GlobalDashboardUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewGlobalDashboardUseCase crea una nueva instancia
func NewGlobalDashboardUseCase(db *gorm.DB, log *logger.Logger) *GlobalDashboardUseCase {
	return &GlobalDashboardUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *GlobalDashboardUseCase) Execute(ctx context.Context, adminID int64) (*DashboardKPIs, error) {
	kpis := &DashboardKPIs{}

	// Timestamps para filtros
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := now.AddDate(0, 0, -7)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	last24h := now.Add(-24 * time.Hour)

	// === USUARIOS ===
	var userStats struct {
		Total     int64
		Active    int64
		Suspended int64
		Banned    int64
		Today     int64
		Week      int64
		Month     int64
	}

	uc.db.Table("users").
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active,
			COUNT(CASE WHEN status = 'suspended' THEN 1 END) as suspended,
			COUNT(CASE WHEN status = 'banned' THEN 1 END) as banned,
			COUNT(CASE WHEN created_at >= ? THEN 1 END) as today,
			COUNT(CASE WHEN created_at >= ? THEN 1 END) as week,
			COUNT(CASE WHEN created_at >= ? THEN 1 END) as month
		`, startOfToday, startOfWeek, startOfMonth).
		Scan(&userStats)

	kpis.TotalUsers = userStats.Total
	kpis.ActiveUsers = userStats.Active
	kpis.SuspendedUsers = userStats.Suspended
	kpis.BannedUsers = userStats.Banned
	kpis.NewUsersToday = userStats.Today
	kpis.NewUsersWeek = userStats.Week
	kpis.NewUsersMonth = userStats.Month

	// === ORGANIZADORES ===
	var organizerStats struct {
		Total    int64
		Verified int64
		Pending  int64
	}

	uc.db.Table("users").
		Where("role = ?", "organizer").
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN kyc_level IN ('verified', 'enhanced') THEN 1 END) as verified,
			COUNT(CASE WHEN kyc_level IN ('none', 'basic', 'pending') THEN 1 END) as pending
		`).
		Scan(&organizerStats)

	kpis.TotalOrganizers = organizerStats.Total
	kpis.VerifiedOrganizers = organizerStats.Verified
	kpis.PendingOrganizers = organizerStats.Pending

	// === RIFAS ===
	var raffleStats struct {
		Total     int64
		Active    int64
		Completed int64
		Suspended int64
		Draft     int64
	}

	uc.db.Table("raffles").
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'suspended' THEN 1 END) as suspended,
			COUNT(CASE WHEN status = 'draft' THEN 1 END) as draft
		`).
		Scan(&raffleStats)

	kpis.TotalRaffles = raffleStats.Total
	kpis.ActiveRaffles = raffleStats.Active
	kpis.CompletedRaffles = raffleStats.Completed
	kpis.SuspendedRaffles = raffleStats.Suspended
	kpis.DraftRaffles = raffleStats.Draft

	// === REVENUE (de pagos exitosos) ===
	var revenueStats struct {
		Today   float64
		Week    float64
		Month   float64
		Year    float64
		AllTime float64
	}

	uc.db.Table("payments").
		Where("status = ?", "succeeded").
		Select(`
			COALESCE(SUM(CASE WHEN paid_at >= ? THEN amount ELSE 0 END), 0) as today,
			COALESCE(SUM(CASE WHEN paid_at >= ? THEN amount ELSE 0 END), 0) as week,
			COALESCE(SUM(CASE WHEN paid_at >= ? THEN amount ELSE 0 END), 0) as month,
			COALESCE(SUM(CASE WHEN paid_at >= ? THEN amount ELSE 0 END), 0) as year,
			COALESCE(SUM(amount), 0) as all_time
		`, startOfToday, startOfWeek, startOfMonth, startOfYear).
		Scan(&revenueStats)

	kpis.RevenueToday = revenueStats.Today
	kpis.RevenueWeek = revenueStats.Week
	kpis.RevenueMonth = revenueStats.Month
	kpis.RevenueYear = revenueStats.Year
	kpis.RevenueAllTime = revenueStats.AllTime

	// === PLATFORM FEES (10% del revenue) ===
	platformFeePercent := 0.10 // TODO: Obtener de configuración
	kpis.PlatformFeesToday = revenueStats.Today * platformFeePercent
	kpis.PlatformFeesMonth = revenueStats.Month * platformFeePercent
	kpis.PlatformFeesAllTime = revenueStats.AllTime * platformFeePercent

	// === SETTLEMENTS ===
	var settlementStats struct {
		PendingCount    int64
		PendingAmount   float64
		ApprovedCount   int64
		ApprovedAmount  float64
	}

	uc.db.Table("settlements").
		Select(`
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN net_amount ELSE 0 END), 0) as pending_amount,
			COUNT(CASE WHEN status = 'approved' THEN 1 END) as approved_count,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN net_amount ELSE 0 END), 0) as approved_amount
		`).
		Scan(&settlementStats)

	kpis.PendingSettlementsCount = settlementStats.PendingCount
	kpis.PendingSettlementsAmount = settlementStats.PendingAmount
	kpis.ApprovedSettlementsCount = settlementStats.ApprovedCount
	kpis.ApprovedSettlementsAmount = settlementStats.ApprovedAmount

	// === PAYMENTS ===
	var paymentStats struct {
		Total     int64
		Succeeded int64
		Pending   int64
		Failed    int64
		Refunded  int64
		Amount    float64
	}

	uc.db.Table("payments").
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'succeeded' THEN 1 END) as succeeded,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
			COUNT(CASE WHEN status = 'refunded' THEN 1 END) as refunded,
			COALESCE(SUM(CASE WHEN status = 'succeeded' THEN amount ELSE 0 END), 0) as amount
		`).
		Scan(&paymentStats)

	kpis.TotalPayments = paymentStats.Total
	kpis.SucceededPayments = paymentStats.Succeeded
	kpis.PendingPayments = paymentStats.Pending
	kpis.FailedPayments = paymentStats.Failed
	kpis.RefundedPayments = paymentStats.Refunded
	kpis.TotalPaymentsAmount = paymentStats.Amount

	// === ACTIVIDAD RECIENTE (últimas 24h) ===
	var recentActivity struct {
		Users       int64
		Raffles     int64
		Payments    int64
		Settlements int64
	}

	// Usuarios recientes
	uc.db.Table("users").Where("created_at >= ?", last24h).Count(&recentActivity.Users)

	// Rifas recientes
	uc.db.Table("raffles").Where("created_at >= ? AND deleted_at IS NULL", last24h).Count(&recentActivity.Raffles)

	// Pagos recientes
	uc.db.Table("payments").Where("created_at >= ?", last24h).Count(&recentActivity.Payments)

	// Settlements recientes
	uc.db.Table("settlements").Where("created_at >= ?", last24h).Count(&recentActivity.Settlements)

	kpis.RecentUsers = recentActivity.Users
	kpis.RecentRaffles = recentActivity.Raffles
	kpis.RecentPayments = recentActivity.Payments
	kpis.RecentSettlements = recentActivity.Settlements

	// Log auditoría
	uc.log.Info("Admin accessed global dashboard",
		logger.Int64("admin_id", adminID),
		logger.String("action", "admin_view_dashboard"))

	return kpis, nil
}
