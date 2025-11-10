package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// AuditLogRepositoryImpl implementa domain.AuditLogRepository
type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

// NewAuditLogRepository crea una nueva instancia del repositorio
func NewAuditLogRepository(db *gorm.DB) domain.AuditLogRepository {
	return &AuditLogRepositoryImpl{db: db}
}

// Create crea un nuevo registro de auditoría
func (r *AuditLogRepositoryImpl) Create(log *domain.AuditLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// FindByID busca un log por ID
func (r *AuditLogRepositoryImpl) FindByID(id int64) (*domain.AuditLog, error) {
	var log domain.AuditLog
	if err := r.db.First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &log, nil
}

// FindByUser busca logs de un usuario específico
func (r *AuditLogRepositoryImpl) FindByUser(userID int64, offset, limit int) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	query := r.db.Model(&domain.AuditLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return logs, total, nil
}

// FindByAdmin busca logs de acciones de admin
func (r *AuditLogRepositoryImpl) FindByAdmin(adminID int64, offset, limit int) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	query := r.db.Model(&domain.AuditLog{}).Where("admin_id = ?", adminID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return logs, total, nil
}

// FindByEntity busca logs de una entidad específica
func (r *AuditLogRepositoryImpl) FindByEntity(entityType string, entityID int64, offset, limit int) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	query := r.db.Model(&domain.AuditLog{}).Where("entity_type = ? AND entity_id = ?", entityType, entityID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return logs, total, nil
}

// FindByAction busca logs por tipo de acción
func (r *AuditLogRepositoryImpl) FindByAction(action domain.AuditAction, offset, limit int) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	query := r.db.Model(&domain.AuditLog{}).Where("action = ?", action)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return logs, total, nil
}

// FindBySeverity busca logs por severidad
func (r *AuditLogRepositoryImpl) FindBySeverity(severity domain.AuditSeverity, offset, limit int) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	query := r.db.Model(&domain.AuditLog{}).Where("severity = ?", severity)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return logs, total, nil
}

// FindByDateRange busca logs en un rango de fechas
func (r *AuditLogRepositoryImpl) FindByDateRange(start, end time.Time, offset, limit int) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	query := r.db.Model(&domain.AuditLog{}).Where("created_at BETWEEN ? AND ?", start, end)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return logs, total, nil
}

// List retorna una lista paginada de logs con filtros
func (r *AuditLogRepositoryImpl) List(offset, limit int, filters map[string]interface{}) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	query := r.db.Model(&domain.AuditLog{})

	// Aplicar filtros
	if userID, ok := filters["user_id"].(int64); ok {
		query = query.Where("user_id = ?", userID)
	}
	if adminID, ok := filters["admin_id"].(int64); ok {
		query = query.Where("admin_id = ?", adminID)
	}
	if action, ok := filters["action"].(domain.AuditAction); ok {
		query = query.Where("action = ?", action)
	}
	if severity, ok := filters["severity"].(domain.AuditSeverity); ok {
		query = query.Where("severity = ?", severity)
	}
	if entityType, ok := filters["entity_type"].(string); ok {
		query = query.Where("entity_type = ?", entityType)
	}
	if entityID, ok := filters["entity_id"].(int64); ok {
		query = query.Where("entity_id = ?", entityID)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener página
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return logs, total, nil
}
