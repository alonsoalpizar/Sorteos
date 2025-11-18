package db

import (
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresSettlementRepository implementación de SettlementRepository con PostgreSQL
type PostgresSettlementRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewSettlementRepository crea una nueva instancia
func NewSettlementRepository(db *gorm.DB, log *logger.Logger) *PostgresSettlementRepository {
	return &PostgresSettlementRepository{
		db:  db,
		log: log,
	}
}

// Create crea un nuevo settlement
func (r *PostgresSettlementRepository) Create(settlement *domain.Settlement) error {
	// Validar antes de crear
	if err := settlement.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	if err := r.db.Create(settlement).Error; err != nil {
		r.log.Error("Error creating settlement",
			logger.Int64("raffle_id", settlement.RaffleID),
			logger.Int64("organizer_id", settlement.OrganizerID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// GetByID obtiene un settlement por ID
func (r *PostgresSettlementRepository) GetByID(id int64) (*domain.Settlement, error) {
	var settlement domain.Settlement

	if err := r.db.
		Preload("Raffle").
		Preload("Organizer").
		First(&settlement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting settlement by ID",
			logger.Int64("id", id),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &settlement, nil
}

// GetByUUID obtiene un settlement por UUID
func (r *PostgresSettlementRepository) GetByUUID(uuid string) (*domain.Settlement, error) {
	var settlement domain.Settlement

	if err := r.db.
		Preload("Raffle").
		Preload("Organizer").
		Where("uuid = ?", uuid).
		First(&settlement).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting settlement by UUID",
			logger.String("uuid", uuid),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &settlement, nil
}

// GetByRaffleID obtiene un settlement por raffle_id
func (r *PostgresSettlementRepository) GetByRaffleID(raffleID int64) (*domain.Settlement, error) {
	var settlement domain.Settlement

	if err := r.db.
		Preload("Raffle").
		Preload("Organizer").
		Where("raffle_id = ?", raffleID).
		First(&settlement).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting settlement by raffle_id",
			logger.Int64("raffle_id", raffleID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &settlement, nil
}

// List obtiene settlements con filtros y paginación
func (r *PostgresSettlementRepository) List(filters map[string]interface{}, offset, limit int) ([]*domain.Settlement, int64, error) {
	var settlements []*domain.Settlement
	var total int64

	query := r.db.Model(&domain.Settlement{})

	// Aplicar filtros
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if organizerID, ok := filters["organizer_id"].(int64); ok && organizerID > 0 {
		query = query.Where("organizer_id = ?", organizerID)
	}

	if dateFrom, ok := filters["date_from"].(time.Time); ok {
		query = query.Where("created_at >= ?", dateFrom)
	}

	if dateTo, ok := filters["date_to"].(time.Time); ok {
		query = query.Where("created_at <= ?", dateTo)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		r.log.Error("Error counting settlements", logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener registros con paginación
	if err := query.
		Preload("Raffle").
		Preload("Organizer").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&settlements).Error; err != nil {
		r.log.Error("Error listing settlements", logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return settlements, total, nil
}

// UpdateStatus actualiza solo el status
func (r *PostgresSettlementRepository) UpdateStatus(id int64, status domain.SettlementStatus) error {
	if err := r.db.Model(&domain.Settlement{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		r.log.Error("Error updating settlement status",
			logger.Int64("id", id),
			logger.String("status", string(status)),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// Approve aprueba un settlement
func (r *PostgresSettlementRepository) Approve(id int64, adminID int64) error {
	// Primero obtener el settlement para validar que puede ser aprobado
	settlement, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if !settlement.CanApprove() {
		return errors.ErrValidationFailed
	}

	now := time.Now()

	if err := r.db.Model(&domain.Settlement{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      domain.SettlementStatusApproved,
			"approved_by": adminID,
			"approved_at": now,
		}).Error; err != nil {
		r.log.Error("Error approving settlement",
			logger.Int64("id", id),
			logger.Int64("admin_id", adminID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// Reject rechaza un settlement
func (r *PostgresSettlementRepository) Reject(id int64, reason string) error {
	// Primero obtener el settlement para validar que puede ser rechazado
	settlement, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if !settlement.CanReject() {
		return errors.ErrValidationFailed
	}

	if err := r.db.Model(&domain.Settlement{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": domain.SettlementStatusRejected,
			"notes":  reason,
		}).Error; err != nil {
		r.log.Error("Error rejecting settlement",
			logger.Int64("id", id),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// MarkPaid marca un settlement como pagado
func (r *PostgresSettlementRepository) MarkPaid(id int64, paymentMethod, paymentReference string) error {
	// Primero obtener el settlement para validar que puede ser marcado como pagado
	settlement, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if !settlement.CanMarkPaid() {
		return errors.ErrValidationFailed
	}

	now := time.Now()

	if err := r.db.Model(&domain.Settlement{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":            domain.SettlementStatusPaid,
			"payment_method":    paymentMethod,
			"payment_reference": paymentReference,
			"paid_at":           now,
		}).Error; err != nil {
		r.log.Error("Error marking settlement as paid",
			logger.Int64("id", id),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// GetPendingByOrganizer obtiene settlements pendientes de un organizador
func (r *PostgresSettlementRepository) GetPendingByOrganizer(organizerID int64) ([]*domain.Settlement, error) {
	var settlements []*domain.Settlement

	if err := r.db.
		Preload("Raffle").
		Where("organizer_id = ? AND status = ?", organizerID, domain.SettlementStatusPending).
		Order("created_at ASC").
		Find(&settlements).Error; err != nil {
		r.log.Error("Error getting pending settlements by organizer",
			logger.Int64("organizer_id", organizerID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return settlements, nil
}

// GetTotalsByStatus obtiene totales agrupados por status
func (r *PostgresSettlementRepository) GetTotalsByStatus() (map[domain.SettlementStatus]float64, error) {
	type statusTotal struct {
		Status domain.SettlementStatus
		Total  float64
	}

	var results []statusTotal

	if err := r.db.Model(&domain.Settlement{}).
		Select("status, COALESCE(SUM(net_payout), 0) as total").
		Group("status").
		Scan(&results).Error; err != nil {
		r.log.Error("Error getting settlement totals by status", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Convertir a map
	totals := make(map[domain.SettlementStatus]float64)
	for _, result := range results {
		totals[result.Status] = result.Total
	}

	return totals, nil
}
