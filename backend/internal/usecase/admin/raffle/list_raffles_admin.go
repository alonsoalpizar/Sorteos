package raffle

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// RaffleAdminMetrics métricas extendidas de una rifa para admin
type RaffleAdminMetrics struct {
	Raffle          *domain.Raffle `json:"raffle"`
	OrganizerName   string         `json:"organizer_name"`
	OrganizerEmail  string         `json:"organizer_email"`
	SoldCount       int            `json:"sold_count"`
	ReservedCount   int            `json:"reserved_count"`
	AvailableCount  int            `json:"available_count"`
	TotalRevenue    float64        `json:"total_revenue"`
	PlatformFee     float64        `json:"platform_fee"`
	NetRevenue      float64        `json:"net_revenue"`
	ConversionRate  float64        `json:"conversion_rate"` // sold / total
}

// ListRafflesAdminInput datos de entrada
type ListRafflesAdminInput struct {
	Page         int
	PageSize     int
	Status       *domain.RaffleStatus // Incluye suspended
	OrganizerID  *int64
	CategoryID   *int64
	Search       string // Buscar en title
	DateFrom     *string
	DateTo       *string
	OrderBy      string
	IncludeAll   bool // Si true, incluye rifas eliminadas (deleted_at)
}

// ListRafflesAdminOutput resultado
type ListRafflesAdminOutput struct {
	Raffles    []*RaffleAdminMetrics
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListRafflesAdminUseCase caso de uso para listar rifas (admin)
type ListRafflesAdminUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListRafflesAdminUseCase crea una nueva instancia
func NewListRafflesAdminUseCase(db *gorm.DB, log *logger.Logger) *ListRafflesAdminUseCase {
	return &ListRafflesAdminUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListRafflesAdminUseCase) Execute(ctx context.Context, input *ListRafflesAdminInput, adminID int64) (*ListRafflesAdminOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Calcular offset
	offset := (input.Page - 1) * input.PageSize

	// Construir query base con JOIN a users para obtener info del organizador
	query := uc.db.Model(&domain.Raffle{}).
		Select(`raffles.*,
			users.name as organizer_name,
			users.email as organizer_email`).
		Joins("LEFT JOIN users ON users.id = raffles.user_id")

	// Por defecto, excluir eliminadas
	if !input.IncludeAll {
		query = query.Where("raffles.deleted_at IS NULL")
	}

	// Aplicar filtros
	if input.Status != nil {
		query = query.Where("raffles.status = ?", *input.Status)
	}

	if input.OrganizerID != nil {
		query = query.Where("raffles.user_id = ?", *input.OrganizerID)
	}

	if input.CategoryID != nil {
		query = query.Where("raffles.category_id = ?", *input.CategoryID)
	}

	if input.Search != "" {
		searchPattern := "%" + input.Search + "%"
		query = query.Where("raffles.title ILIKE ?", searchPattern)
	}

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("raffles.created_at >= ?", *input.DateFrom)
	}

	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("raffles.created_at <= ?", *input.DateTo)
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		uc.log.Error("Error counting raffles", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Aplicar ordenamiento
	orderBy := "raffles.created_at DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}
	query = query.Order(orderBy)

	// Obtener raffles con paginación
	var results []struct {
		domain.Raffle
		OrganizerName  string
		OrganizerEmail string
	}

	if err := query.Offset(offset).Limit(input.PageSize).Scan(&results).Error; err != nil {
		uc.log.Error("Error listing raffles", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Construir métricas
	raffles := make([]*RaffleAdminMetrics, 0, len(results))
	for _, result := range results {
		// Convertir decimal.Decimal a float64
		totalRevenue, _ := result.TotalRevenue.Float64()
		platformFee, _ := result.PlatformFeeAmount.Float64()
		netRevenue, _ := result.NetAmount.Float64()

		// Calcular conversion rate
		conversionRate := 0.0
		if result.TotalNumbers > 0 {
			conversionRate = float64(result.SoldCount) / float64(result.TotalNumbers) * 100
		}

		raffles = append(raffles, &RaffleAdminMetrics{
			Raffle:         &result.Raffle,
			OrganizerName:  result.OrganizerName,
			OrganizerEmail: result.OrganizerEmail,
			SoldCount:      result.SoldCount,
			ReservedCount:  result.ReservedCount,
			AvailableCount: result.AvailableCount(),
			TotalRevenue:   totalRevenue,
			PlatformFee:    platformFee,
			NetRevenue:     netRevenue,
			ConversionRate: conversionRate,
		})
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	// Log auditoría
	uc.log.Info("Admin listed raffles",
		logger.Int64("admin_id", adminID),
		logger.Int("total_results", len(raffles)),
		logger.String("action", "admin_list_raffles"))

	return &ListRafflesAdminOutput{
		Raffles:    raffles,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}
