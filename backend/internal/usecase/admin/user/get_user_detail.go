package user

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UserStats estadísticas del usuario
type UserStats struct {
	TotalRaffles      int     `json:"total_raffles"`
	ActiveRaffles     int     `json:"active_raffles"`
	CompletedRaffles  int     `json:"completed_raffles"`
	TotalRevenue      float64 `json:"total_revenue"`
	TotalTicketsBought int    `json:"total_tickets_bought"`
	TotalSpent        float64 `json:"total_spent"`
}

// GetUserDetailOutput resultado del detalle de usuario
type GetUserDetailOutput struct {
	User         *domain.User  `json:"user"`
	Stats        *UserStats    `json:"stats"`
	RecentRaffles []*domain.Raffle `json:"recent_raffles,omitempty"`
}

// GetUserDetailUseCase caso de uso para obtener detalle de usuario (admin)
type GetUserDetailUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewGetUserDetailUseCase crea una nueva instancia
func NewGetUserDetailUseCase(db *gorm.DB, log *logger.Logger) *GetUserDetailUseCase {
	return &GetUserDetailUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *GetUserDetailUseCase) Execute(ctx context.Context, userID int64, adminID int64) (*GetUserDetailOutput, error) {
	// Obtener usuario
	var user domain.User
	if err := uc.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		uc.log.Error("Error finding user", logger.Int64("user_id", userID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener estadísticas de rifas (como organizador)
	var raffleStats struct {
		TotalRaffles     int
		ActiveRaffles    int
		CompletedRaffles int
		TotalRevenue     float64
	}

	err := uc.db.Model(&domain.Raffle{}).
		Select(`
			COUNT(*) as total_raffles,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_raffles,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_raffles,
			COALESCE(SUM(total_revenue), 0) as total_revenue
		`).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Scan(&raffleStats).Error

	if err != nil {
		uc.log.Error("Error getting raffle stats", logger.Int64("user_id", userID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener estadísticas de tickets comprados (como participante)
	var ticketStats struct {
		TotalTickets int
		TotalSpent   float64
	}

	err = uc.db.Table("raffle_numbers").
		Select(`
			COUNT(*) as total_tickets,
			COALESCE(SUM(price), 0) as total_spent
		`).
		Where("user_id = ?", userID).
		Scan(&ticketStats).Error

	if err != nil {
		uc.log.Error("Error getting ticket stats", logger.Int64("user_id", userID), logger.Error(err))
		// No fallar por esto, continuar con stats vacías
		ticketStats.TotalTickets = 0
		ticketStats.TotalSpent = 0
	}

	// Obtener rifas recientes (últimas 5 como organizador)
	var recentRaffles []*domain.Raffle
	err = uc.db.
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Limit(5).
		Find(&recentRaffles).Error

	if err != nil {
		uc.log.Error("Error getting recent raffles", logger.Int64("user_id", userID), logger.Error(err))
		// No fallar por esto
		recentRaffles = []*domain.Raffle{}
	}

	// Log auditoría
	uc.log.Info("Admin viewed user detail",
		logger.Int64("admin_id", adminID),
		logger.Int64("user_id", userID),
		logger.String("action", "admin_view_user_detail"))

	return &GetUserDetailOutput{
		User: &user,
		Stats: &UserStats{
			TotalRaffles:       raffleStats.TotalRaffles,
			ActiveRaffles:      raffleStats.ActiveRaffles,
			CompletedRaffles:   raffleStats.CompletedRaffles,
			TotalRevenue:       raffleStats.TotalRevenue,
			TotalTicketsBought: ticketStats.TotalTickets,
			TotalSpent:         ticketStats.TotalSpent,
		},
		RecentRaffles: recentRaffles,
	}, nil
}
