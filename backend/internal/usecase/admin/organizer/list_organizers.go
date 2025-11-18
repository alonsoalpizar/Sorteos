package organizer

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// OrganizerWithMetrics combina perfil de organizador con métricas
type OrganizerWithMetrics struct {
	Profile          *domain.OrganizerProfile `json:"profile"`
	User             *domain.User             `json:"user"`
	TotalRaffles     int                      `json:"total_raffles"`
	ActiveRaffles    int                      `json:"active_raffles"`
	CompletedRaffles int                      `json:"completed_raffles"`
	TotalRevenue     float64                  `json:"total_revenue"`
	PendingPayout    float64                  `json:"pending_payout"`
}

// ListOrganizersInput datos de entrada para listar organizadores
type ListOrganizersInput struct {
	Page      int
	PageSize  int
	Verified  *bool
	DateFrom  *string
	DateTo    *string
	Search    string // Buscar en business_name, user name, email
	OrderBy   string
}

// ListOrganizersOutput resultado del listado
type ListOrganizersOutput struct {
	Organizers []*OrganizerWithMetrics
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// ListOrganizersUseCase caso de uso para listar organizadores (admin)
type ListOrganizersUseCase struct {
	organizerRepo *db.PostgresOrganizerProfileRepository
	log           *logger.Logger
}

// NewListOrganizersUseCase crea una nueva instancia
func NewListOrganizersUseCase(
	organizerRepo *db.PostgresOrganizerProfileRepository,
	log *logger.Logger,
) *ListOrganizersUseCase {
	return &ListOrganizersUseCase{
		organizerRepo: organizerRepo,
		log:           log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListOrganizersUseCase) Execute(ctx context.Context, input *ListOrganizersInput, adminID int64) (*ListOrganizersOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Calcular offset
	offset := (input.Page - 1) * input.PageSize

	// Construir filtros
	filters := make(map[string]interface{})

	if input.Verified != nil {
		filters["verified"] = *input.Verified
	}

	if input.DateFrom != nil && *input.DateFrom != "" {
		filters["date_from"] = *input.DateFrom
	}

	if input.DateTo != nil && *input.DateTo != "" {
		filters["date_to"] = *input.DateTo
	}

	// Obtener perfiles de organizadores
	profiles, total, err := uc.organizerRepo.List(filters, offset, input.PageSize)
	if err != nil {
		uc.log.Error("Error listing organizer profiles", logger.Error(err))
		return nil, err
	}

	// Construir resultado con métricas
	organizers := make([]*OrganizerWithMetrics, 0, len(profiles))

	for _, profile := range profiles {
		// Las métricas ya vienen del perfil
		organizers = append(organizers, &OrganizerWithMetrics{
			Profile:       profile,
			User:          profile.User,
			PendingPayout: profile.PendingPayout,
			// Las demás métricas se pueden obtener posteriormente si es necesario
			// o se pueden incluir en el perfil mediante Preload
		})
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	// Log auditoría
	uc.log.Info("Admin listed organizers",
		logger.Int64("admin_id", adminID),
		logger.Int("total_results", len(organizers)),
		logger.String("action", "admin_list_organizers"))

	return &ListOrganizersOutput{
		Organizers: organizers,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}
