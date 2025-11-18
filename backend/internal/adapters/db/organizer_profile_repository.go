package db

import (
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresOrganizerProfileRepository implementación de OrganizerProfileRepository con PostgreSQL
type PostgresOrganizerProfileRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewOrganizerProfileRepository crea una nueva instancia
func NewOrganizerProfileRepository(db *gorm.DB, log *logger.Logger) *PostgresOrganizerProfileRepository {
	return &PostgresOrganizerProfileRepository{
		db:  db,
		log: log,
	}
}

// Create crea un nuevo perfil de organizador
func (r *PostgresOrganizerProfileRepository) Create(profile *domain.OrganizerProfile) error {
	// Validar antes de crear
	if err := profile.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	if err := r.db.Create(profile).Error; err != nil {
		r.log.Error("Error creating organizer profile",
			logger.Int64("user_id", profile.UserID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// GetByUserID obtiene un perfil por user_id
func (r *PostgresOrganizerProfileRepository) GetByUserID(userID int64) (*domain.OrganizerProfile, error) {
	var profile domain.OrganizerProfile

	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting organizer profile by user_id",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &profile, nil
}

// GetByID obtiene un perfil por ID
func (r *PostgresOrganizerProfileRepository) GetByID(id int64) (*domain.OrganizerProfile, error) {
	var profile domain.OrganizerProfile

	if err := r.db.First(&profile, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting organizer profile by ID",
			logger.Int64("id", id),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &profile, nil
}

// List obtiene perfiles con filtros y paginación
func (r *PostgresOrganizerProfileRepository) List(filters map[string]interface{}, offset, limit int) ([]*domain.OrganizerProfile, int64, error) {
	var profiles []*domain.OrganizerProfile
	var total int64

	query := r.db.Model(&domain.OrganizerProfile{})

	// Aplicar filtros
	if verified, ok := filters["verified"].(bool); ok {
		query = query.Where("verified = ?", verified)
	}

	if payoutSchedule, ok := filters["payout_schedule"].(string); ok && payoutSchedule != "" {
		query = query.Where("payout_schedule = ?", payoutSchedule)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		r.log.Error("Error counting organizer profiles", logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener registros con paginación
	if err := query.
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&profiles).Error; err != nil {
		r.log.Error("Error listing organizer profiles", logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return profiles, total, nil
}

// Update actualiza un perfil de organizador
func (r *PostgresOrganizerProfileRepository) Update(profile *domain.OrganizerProfile) error {
	// Validar antes de actualizar
	if err := profile.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	if err := r.db.Save(profile).Error; err != nil {
		r.log.Error("Error updating organizer profile",
			logger.Int64("id", profile.ID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// UpdateCommission actualiza solo la comisión
func (r *PostgresOrganizerProfileRepository) UpdateCommission(userID int64, commission *float64) error {
	if err := r.db.Model(&domain.OrganizerProfile{}).
		Where("user_id = ?", userID).
		Update("commission_override", commission).Error; err != nil {
		r.log.Error("Error updating organizer commission",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// UpdateFinancials actualiza totales financieros
func (r *PostgresOrganizerProfileRepository) UpdateFinancials(userID int64, totalPayouts, pendingPayout float64) error {
	if err := r.db.Model(&domain.OrganizerProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"total_payouts":  totalPayouts,
			"pending_payout": pendingPayout,
		}).Error; err != nil {
		r.log.Error("Error updating organizer financials",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// GetRevenue obtiene el desglose de ingresos de un organizador
func (r *PostgresOrganizerProfileRepository) GetRevenue(userID int64, dateFrom, dateTo *time.Time) (*domain.OrganizerRevenue, error) {
	var revenue domain.OrganizerRevenue

	query := `
		SELECT
			$1 as organizer_id,
			COUNT(*) as total_raffles,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_raffles,
			COALESCE(SUM(total_revenue), 0) as total_revenue,
			COALESCE(SUM(platform_fee_amount), 0) as platform_fees,
			COALESCE(SUM(net_amount), 0) as net_revenue,
			(SELECT COALESCE(SUM(total_payouts), 0) FROM organizer_profiles WHERE user_id = $1) as total_payouts,
			(SELECT COALESCE(SUM(pending_payout), 0) FROM organizer_profiles WHERE user_id = $1) as pending_payout
		FROM raffles
		WHERE user_id = $1
	`

	args := []interface{}{userID}

	// Aplicar filtros de fecha si están presentes
	if dateFrom != nil {
		query += " AND created_at >= $" + string(rune(len(args)+1))
		args = append(args, *dateFrom)
	}

	if dateTo != nil {
		query += " AND created_at <= $" + string(rune(len(args)+1))
		args = append(args, *dateTo)
	}

	if err := r.db.Raw(query, args...).Scan(&revenue).Error; err != nil {
		r.log.Error("Error getting organizer revenue",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &revenue, nil
}

// Verify marca un organizador como verificado
func (r *PostgresOrganizerProfileRepository) Verify(userID int64, verifiedBy int64) error {
	now := time.Now()

	if err := r.db.Model(&domain.OrganizerProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"verified":    true,
			"verified_at": now,
			"verified_by": verifiedBy,
		}).Error; err != nil {
		r.log.Error("Error verifying organizer",
			logger.Int64("user_id", userID),
			logger.Int64("verified_by", verifiedBy),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}
