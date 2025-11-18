package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// AuditLog representa un registro de auditoría
type AuditLog struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	AdminID     int64     `gorm:"not null;index"`
	Action      string    `gorm:"type:varchar(100);not null;index"`
	EntityType  string    `gorm:"type:varchar(50);index"`
	EntityID    *int64    `gorm:"index"`
	Description string    `gorm:"type:text"`
	Severity    string    `gorm:"type:varchar(20);not null;index;default:'info'"` // info, warning, error, critical
	IPAddress   *string   `gorm:"type:varchar(45)"`
	UserAgent   *string   `gorm:"type:text"`
	Metadata    *string   `gorm:"type:jsonb"` // JSON metadata adicional
	CreatedAt   time.Time `gorm:"not null;index"`
}

// TableName especifica el nombre de la tabla
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogRepository interfaz para el repositorio de audit logs
type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	FindByFilters(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*AuditLog, int64, error)
}

// auditLogRepository implementación del repositorio
type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository crea una nueva instancia del repositorio
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create crea un nuevo registro de auditoría
func (r *auditLogRepository) Create(ctx context.Context, log *AuditLog) error {
	log.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(log).Error
}

// FindByFilters busca logs con filtros dinámicos
func (r *auditLogRepository) FindByFilters(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*AuditLog, int64, error) {
	var logs []*AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&AuditLog{})

	// Aplicar filtros dinámicos
	for key, value := range filters {
		if value != nil {
			switch key {
			case "admin_id":
				query = query.Where("admin_id = ?", value)
			case "action":
				query = query.Where("action = ?", value)
			case "entity_type":
				query = query.Where("entity_type = ?", value)
			case "entity_id":
				query = query.Where("entity_id = ?", value)
			case "severity":
				query = query.Where("severity = ?", value)
			case "date_from":
				query = query.Where("created_at >= ?", value)
			case "date_to":
				query = query.Where("created_at <= ?", value)
			case "search":
				searchPattern := "%" + value.(string) + "%"
				query = query.Where("description ILIKE ? OR action ILIKE ?", searchPattern, searchPattern)
			}
		}
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación y obtener resultados
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
