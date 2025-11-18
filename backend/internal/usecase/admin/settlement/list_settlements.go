package settlement

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// SettlementWithDetails liquidación con detalles adicionales
type SettlementWithDetails struct {
	ID                int64      `json:"id"`
	RaffleID          int64      `json:"raffle_id"`
	OrganizerID       int64      `json:"organizer_id"`
	TotalRevenue      float64    `json:"total_revenue"`
	PlatformFee       float64    `json:"platform_fee"`
	NetAmount         float64    `json:"net_amount"`
	Status            string     `json:"status"`
	CalculatedAt      time.Time  `json:"calculated_at"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	ApprovedBy        *int64     `json:"approved_by,omitempty"`
	RejectedAt        *time.Time `json:"rejected_at,omitempty"`
	RejectedBy        *int64     `json:"rejected_by,omitempty"`
	RejectionReason   *string    `json:"rejection_reason,omitempty"`
	PaidAt            *time.Time `json:"paid_at,omitempty"`
	PaymentReference  *string    `json:"payment_reference,omitempty"`
	PaymentMethod     *string    `json:"payment_method,omitempty"`
	AdminNotes        *string    `json:"admin_notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	// Detalles adicionales
	RaffleTitle       string `json:"raffle_title"`
	OrganizerName     string `json:"organizer_name"`
	OrganizerEmail    string `json:"organizer_email"`
	OrganizerKYCLevel string `json:"organizer_kyc_level"`
}

// ListSettlementsInput datos de entrada
type ListSettlementsInput struct {
	Page         int
	PageSize     int
	Status       *string // pending, approved, paid, rejected
	OrganizerID  *int64
	RaffleID     *int64
	DateFrom     *string
	DateTo       *string
	MinAmount    *float64
	MaxAmount    *float64
	Search       string // Buscar en raffle title, organizer name
	OrderBy      string
	KYCLevel     *domain.KYCLevel
	PendingOnly  bool // Si true, solo pending
}

// ListSettlementsOutput resultado
type ListSettlementsOutput struct {
	Settlements       []*SettlementWithDetails
	Total             int64
	Page              int
	PageSize          int
	TotalPages        int
	// Estadísticas por status
	TotalPending      int64
	TotalApproved     int64
	TotalPaid         int64
	TotalRejected     int64
	// Montos totales
	TotalPendingAmount  float64
	TotalApprovedAmount float64
	TotalPaidAmount     float64
}

// ListSettlementsUseCase caso de uso para listar liquidaciones (admin)
type ListSettlementsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListSettlementsUseCase crea una nueva instancia
func NewListSettlementsUseCase(db *gorm.DB, log *logger.Logger) *ListSettlementsUseCase {
	return &ListSettlementsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListSettlementsUseCase) Execute(ctx context.Context, input *ListSettlementsInput, adminID int64) (*ListSettlementsOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Calcular offset
	offset := (input.Page - 1) * input.PageSize

	// Construir query base con JOIN a raffles y users (organizers)
	query := uc.db.Table("settlements").
		Select(`settlements.*,
			raffles.title as raffle_title,
			COALESCE(users.first_name || ' ' || users.last_name, users.email) as organizer_name,
			users.email as organizer_email,
			users.kyc_level as organizer_kyc_level`).
		Joins("LEFT JOIN raffles ON raffles.id = settlements.raffle_id").
		Joins("LEFT JOIN users ON users.id = settlements.organizer_id")

	// Aplicar filtros
	if input.Status != nil {
		query = query.Where("settlements.status = ?", *input.Status)
	} else if input.PendingOnly {
		query = query.Where("settlements.status = ?", "pending")
	}

	if input.OrganizerID != nil {
		query = query.Where("settlements.organizer_id = ?", *input.OrganizerID)
	}

	if input.RaffleID != nil {
		query = query.Where("settlements.raffle_id = ?", *input.RaffleID)
	}

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("settlements.calculated_at >= ?", *input.DateFrom)
	}

	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("settlements.calculated_at <= ?", *input.DateTo)
	}

	if input.MinAmount != nil {
		query = query.Where("settlements.net_amount >= ?", *input.MinAmount)
	}

	if input.MaxAmount != nil {
		query = query.Where("settlements.net_amount <= ?", *input.MaxAmount)
	}

	if input.KYCLevel != nil {
		query = query.Where("users.kyc_level = ?", *input.KYCLevel)
	}

	if input.Search != "" {
		searchPattern := "%" + input.Search + "%"
		query = query.Where(
			"raffles.title ILIKE ? OR users.first_name ILIKE ? OR users.last_name ILIKE ? OR users.email ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		uc.log.Error("Error counting settlements", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Aplicar ordenamiento
	orderBy := "settlements.calculated_at DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}
	query = query.Order(orderBy)

	// Obtener settlements con paginación
	var results []SettlementWithDetails

	if err := query.Offset(offset).Limit(input.PageSize).Scan(&results).Error; err != nil {
		uc.log.Error("Error listing settlements", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Convertir a slice de punteros
	settlements := make([]*SettlementWithDetails, len(results))
	for i := range results {
		settlements[i] = &results[i]
	}

	// Calcular estadísticas por status
	var stats []struct {
		Status string
		Count  int64
		Amount float64
	}

	statsQuery := uc.db.Table("settlements").
		Select("status, COUNT(*) as count, COALESCE(SUM(net_amount), 0) as amount").
		Group("status")

	// Aplicar mismos filtros (sin paginación)
	if input.OrganizerID != nil {
		statsQuery = statsQuery.Where("organizer_id = ?", *input.OrganizerID)
	}
	if input.RaffleID != nil {
		statsQuery = statsQuery.Where("raffle_id = ?", *input.RaffleID)
	}
	if input.DateFrom != nil && *input.DateFrom != "" {
		statsQuery = statsQuery.Where("calculated_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		statsQuery = statsQuery.Where("calculated_at <= ?", *input.DateTo)
	}

	if err := statsQuery.Scan(&stats).Error; err != nil {
		uc.log.Error("Error calculating settlement stats", logger.Error(err))
		// No es crítico, continuamos
	}

	// Organizar estadísticas
	var totalPending, totalApproved, totalPaid, totalRejected int64
	var totalPendingAmount, totalApprovedAmount, totalPaidAmount float64

	for _, stat := range stats {
		switch stat.Status {
		case "pending":
			totalPending = stat.Count
			totalPendingAmount = stat.Amount
		case "approved":
			totalApproved = stat.Count
			totalApprovedAmount = stat.Amount
		case "paid":
			totalPaid = stat.Count
			totalPaidAmount = stat.Amount
		case "rejected":
			totalRejected = stat.Count
		}
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	// Log auditoría
	uc.log.Info("Admin listed settlements",
		logger.Int64("admin_id", adminID),
		logger.Int("total_results", len(settlements)),
		logger.String("action", "admin_list_settlements"))

	return &ListSettlementsOutput{
		Settlements:         settlements,
		Total:               total,
		Page:                input.Page,
		PageSize:            input.PageSize,
		TotalPages:          totalPages,
		TotalPending:        totalPending,
		TotalApproved:       totalApproved,
		TotalPaid:           totalPaid,
		TotalRejected:       totalRejected,
		TotalPendingAmount:  totalPendingAmount,
		TotalApprovedAmount: totalApprovedAmount,
		TotalPaidAmount:     totalPaidAmount,
	}, nil
}
