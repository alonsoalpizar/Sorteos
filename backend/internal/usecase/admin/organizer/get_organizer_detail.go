package organizer

import (
	"context"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GetOrganizerDetailOutput resultado del detalle de organizador
type GetOrganizerDetailOutput struct {
	Profile  *domain.OrganizerProfile `json:"profile"`
	User     *domain.User             `json:"user"`
	Revenue  *domain.OrganizerRevenue `json:"revenue"`
}

// GetOrganizerDetailUseCase caso de uso para obtener detalle de organizador (admin)
type GetOrganizerDetailUseCase struct {
	organizerRepo *db.PostgresOrganizerProfileRepository
	log           *logger.Logger
}

// NewGetOrganizerDetailUseCase crea una nueva instancia
func NewGetOrganizerDetailUseCase(
	organizerRepo *db.PostgresOrganizerProfileRepository,
	log *logger.Logger,
) *GetOrganizerDetailUseCase {
	return &GetOrganizerDetailUseCase{
		organizerRepo: organizerRepo,
		log:           log,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetOrganizerDetailUseCase) Execute(ctx context.Context, userID int64, adminID int64) (*GetOrganizerDetailOutput, error) {
	// Obtener perfil de organizador
	profile, err := uc.organizerRepo.GetByUserID(userID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.New("NOT_FOUND", "organizer profile not found", 404, nil)
		}
		uc.log.Error("Error getting organizer profile", logger.Int64("user_id", userID), logger.Error(err))
		return nil, err
	}

	// Obtener revenue (sin filtro de fecha = todos los tiempos)
	revenue, err := uc.organizerRepo.GetRevenue(userID, nil, nil)
	if err != nil {
		uc.log.Error("Error getting organizer revenue", logger.Int64("user_id", userID), logger.Error(err))
		return nil, err
	}

	// Log auditoría
	uc.log.Info("Admin viewed organizer detail",
		logger.Int64("admin_id", adminID),
		logger.Int64("user_id", userID),
		logger.String("action", "admin_view_organizer_detail"))

	return &GetOrganizerDetailOutput{
		Profile:  profile,
		User:     profile.User,
		Revenue:  revenue,
	}, nil
}
