package audit

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// AuditLog registro de auditoría
type AuditLog struct {
	ID           int64     `json:"id"`
	AdminID      int64     `json:"admin_id"`
	AdminName    string    `json:"admin_name"`
	AdminEmail   string    `json:"admin_email"`
	Action       string    `json:"action"`
	EntityType   string    `json:"entity_type"`   // user, raffle, payment, settlement, etc.
	EntityID     *int64    `json:"entity_id,omitempty"`
	Description  string    `json:"description"`
	Severity     string    `json:"severity"`      // info, warning, error, critical
	IPAddress    *string   `json:"ip_address,omitempty"`
	UserAgent    *string   `json:"user_agent,omitempty"`
	Metadata     string    `json:"metadata,omitempty"`  // JSON string
	CreatedAt    time.Time `json:"created_at"`
}

// ListAuditLogsInput datos de entrada
type ListAuditLogsInput struct {
	Page         int
	PageSize     int
	AdminID      *int64
	Action       *string  // Filtrar por tipo de acción específica
	EntityType   *string  // user, raffle, payment, settlement
	EntityID     *int64   // Filtrar por entidad específica
	Severity     *string  // info, warning, error, critical
	DateFrom     *string
	DateTo       *string
	Search       string   // Buscar en description, admin name
	OrderBy      string
}

// ListAuditLogsOutput resultado
type ListAuditLogsOutput struct {
	Logs       []*AuditLog
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
	// Estadísticas
	InfoCount     int64
	WarningCount  int64
	ErrorCount    int64
	CriticalCount int64
}

// ListAuditLogsUseCase caso de uso para listar audit logs
type ListAuditLogsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListAuditLogsUseCase crea una nueva instancia
func NewListAuditLogsUseCase(db *gorm.DB, log *logger.Logger) *ListAuditLogsUseCase {
	return &ListAuditLogsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListAuditLogsUseCase) Execute(ctx context.Context, input *ListAuditLogsInput, adminID int64) (*ListAuditLogsOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 50
	}

	offset := (input.Page - 1) * input.PageSize

	// NOTA: Asumimos que existe una tabla audit_logs
	// Si no existe, estos logs vendrán del sistema de logging (zap logs parseados)
	// Por ahora, creamos una implementación que usa la tabla audit_logs

	// Construir query base con JOIN a users (admins)
	query := uc.db.Table("audit_logs").
		Select(`audit_logs.*,
			COALESCE(users.first_name || ' ' || users.last_name, users.email) as admin_name,
			users.email as admin_email`).
		Joins("LEFT JOIN users ON users.id = audit_logs.admin_id")

	// Aplicar filtros
	if input.AdminID != nil {
		query = query.Where("audit_logs.admin_id = ?", *input.AdminID)
	}

	if input.Action != nil {
		query = query.Where("audit_logs.action = ?", *input.Action)
	}

	if input.EntityType != nil {
		query = query.Where("audit_logs.entity_type = ?", *input.EntityType)
	}

	if input.EntityID != nil {
		query = query.Where("audit_logs.entity_id = ?", *input.EntityID)
	}

	if input.Severity != nil {
		query = query.Where("audit_logs.severity = ?", *input.Severity)
	}

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("audit_logs.created_at >= ?", *input.DateFrom)
	}

	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("audit_logs.created_at <= ?", *input.DateTo+" 23:59:59")
	}

	if input.Search != "" {
		searchPattern := "%" + input.Search + "%"
		query = query.Where(
			"audit_logs.description ILIKE ? OR audit_logs.action ILIKE ? OR users.first_name ILIKE ? OR users.last_name ILIKE ? OR users.email ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		uc.log.Error("Error counting audit logs", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Aplicar ordenamiento
	orderBy := "audit_logs.created_at DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}
	query = query.Order(orderBy)

	// Obtener logs con paginación
	var logs []*AuditLog

	if err := query.Offset(offset).Limit(input.PageSize).Scan(&logs).Error; err != nil {
		uc.log.Error("Error listing audit logs", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Calcular estadísticas por severity
	var stats struct {
		Info     int64
		Warning  int64
		Error    int64
		Critical int64
	}

	statsQuery := uc.db.Table("audit_logs").
		Select(`
			COUNT(CASE WHEN severity = 'info' THEN 1 END) as info,
			COUNT(CASE WHEN severity = 'warning' THEN 1 END) as warning,
			COUNT(CASE WHEN severity = 'error' THEN 1 END) as error,
			COUNT(CASE WHEN severity = 'critical' THEN 1 END) as critical
		`)

	// Aplicar mismos filtros (sin paginación)
	if input.AdminID != nil {
		statsQuery = statsQuery.Where("admin_id = ?", *input.AdminID)
	}
	if input.Action != nil {
		statsQuery = statsQuery.Where("action = ?", *input.Action)
	}
	if input.EntityType != nil {
		statsQuery = statsQuery.Where("entity_type = ?", *input.EntityType)
	}
	if input.DateFrom != nil && *input.DateFrom != "" {
		statsQuery = statsQuery.Where("created_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		statsQuery = statsQuery.Where("created_at <= ?", *input.DateTo+" 23:59:59")
	}

	if err := statsQuery.Scan(&stats).Error; err != nil {
		uc.log.Error("Error calculating audit log stats", logger.Error(err))
		// No es crítico, continuamos
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	// Log auditoría (meta-auditoría: registrar que se consultaron los logs)
	uc.log.Info("Admin viewed audit logs",
		logger.Int64("admin_id", adminID),
		logger.Int("total_results", len(logs)),
		logger.String("action", "admin_view_audit_logs"))

	return &ListAuditLogsOutput{
		Logs:          logs,
		Total:         total,
		Page:          input.Page,
		PageSize:      input.PageSize,
		TotalPages:    totalPages,
		InfoCount:     stats.Info,
		WarningCount:  stats.Warning,
		ErrorCount:    stats.Error,
		CriticalCount: stats.Critical,
	}, nil
}
