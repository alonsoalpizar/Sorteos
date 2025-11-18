package notifications

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ViewNotificationHistoryInput datos de entrada
type ViewNotificationHistoryInput struct {
	Type       *string    `json:"type,omitempty"`       // email, sms, push, announcement
	Status     *string    `json:"status,omitempty"`     // queued, sent, failed, scheduled
	Priority   *string    `json:"priority,omitempty"`   // low, normal, high, critical
	AdminID    *int64     `json:"admin_id,omitempty"`   // Filtrar por admin que envió
	DateFrom   *string    `json:"date_from,omitempty"`  // Filtro de fecha
	DateTo     *string    `json:"date_to,omitempty"`
	Search     *string    `json:"search,omitempty"`     // Buscar en subject/body
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// ViewNotificationHistoryOutput resultado
type ViewNotificationHistoryOutput struct {
	Notifications []*NotificationHistoryItem `json:"notifications"`
	TotalCount    int                        `json:"total_count"`
	Statistics    *NotificationStatistics    `json:"statistics"`
}

// NotificationHistoryItem item del historial
type NotificationHistoryItem struct {
	ID             int64                  `json:"id"`
	Type           string                 `json:"type"`
	Subject        string                 `json:"subject,omitempty"`
	Recipients     []EmailRecipient       `json:"recipients,omitempty"`
	RecipientCount int                    `json:"recipient_count"`
	Priority       string                 `json:"priority"`
	Status         string                 `json:"status"`
	SentAt         string                 `json:"sent_at,omitempty"`
	ScheduledAt    string                 `json:"scheduled_at,omitempty"`
	ProviderStatus string                 `json:"provider_status,omitempty"`
	Error          string                 `json:"error,omitempty"`
	AdminID        int64                  `json:"admin_id"`
	AdminEmail     string                 `json:"admin_email"`
	CreatedAt      string                 `json:"created_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// NotificationStatistics estadísticas del historial
type NotificationStatistics struct {
	TotalSent       int     `json:"total_sent"`
	TotalFailed     int     `json:"total_failed"`
	TotalQueued     int     `json:"total_queued"`
	TotalScheduled  int     `json:"total_scheduled"`
	SuccessRate     float64 `json:"success_rate"`
	AveragePerDay   float64 `json:"average_per_day"`
	LastSentAt      string  `json:"last_sent_at,omitempty"`
}

// ViewNotificationHistoryUseCase caso de uso para ver historial de notificaciones
type ViewNotificationHistoryUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewViewNotificationHistoryUseCase crea una nueva instancia
func NewViewNotificationHistoryUseCase(db *gorm.DB, log *logger.Logger) *ViewNotificationHistoryUseCase {
	return &ViewNotificationHistoryUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ViewNotificationHistoryUseCase) Execute(ctx context.Context, input *ViewNotificationHistoryInput, adminID int64) (*ViewNotificationHistoryOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Construir query
	query := uc.db.WithContext(ctx).Table("email_notifications")

	// Aplicar filtros
	if input.Type != nil && *input.Type != "" {
		query = query.Where("type = ?", *input.Type)
	}
	if input.Status != nil && *input.Status != "" {
		query = query.Where("status = ?", *input.Status)
	}
	if input.Priority != nil && *input.Priority != "" {
		query = query.Where("priority = ?", *input.Priority)
	}
	if input.AdminID != nil {
		query = query.Where("admin_id = ?", *input.AdminID)
	}
	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("created_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("created_at <= ?", *input.DateTo+" 23:59:59")
	}
	if input.Search != nil && *input.Search != "" {
		searchPattern := "%" + *input.Search + "%"
		query = query.Where("subject ILIKE ? OR body ILIKE ?", searchPattern, searchPattern)
	}

	// Contar total
	var totalCount int64
	countQuery := query
	countQuery.Count(&totalCount)

	// Aplicar paginación
	query = query.Order("created_at DESC").Limit(input.Limit).Offset(input.Offset)

	// Ejecutar query con JOIN a users para obtener email del admin
	rows, err := query.
		Select("email_notifications.*, users.email as admin_email").
		Joins("LEFT JOIN users ON users.id = email_notifications.admin_id").
		Rows()
	if err != nil {
		uc.log.Error("Error querying notification history", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	// Parsear resultados
	notifications := make([]*NotificationHistoryItem, 0)
	for rows.Next() {
		var notification EmailNotification
		var adminEmail string
		var subject, recipients, providerStatus, errorMsg *string

		err := rows.Scan(
			&notification.ID,
			&notification.AdminID,
			&notification.Type,
			&recipients,
			&subject,
			&notification.Body,
			&notification.TemplateID,
			&notification.Variables,
			&notification.Priority,
			&notification.Status,
			&notification.SentAt,
			&notification.ScheduledAt,
			&notification.ProviderID,
			&providerStatus,
			&errorMsg,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&adminEmail,
		)
		if err != nil {
			uc.log.Error("Error scanning notification row", logger.Error(err))
			continue
		}

		// Deserializar recipients
		var recipientsList []EmailRecipient
		if recipients != nil {
			json.Unmarshal([]byte(*recipients), &recipientsList)
		}

		// Construir item
		item := &NotificationHistoryItem{
			ID:             notification.ID,
			Type:           notification.Type,
			Subject:        "",
			Recipients:     recipientsList,
			RecipientCount: len(recipientsList),
			Priority:       notification.Priority,
			Status:         notification.Status,
			AdminID:        notification.AdminID,
			AdminEmail:     adminEmail,
			CreatedAt:      notification.CreatedAt.Format(time.RFC3339),
		}

		if subject != nil {
			item.Subject = *subject
		}
		if notification.SentAt != nil {
			item.SentAt = notification.SentAt.Format(time.RFC3339)
		}
		if notification.ScheduledAt != nil {
			item.ScheduledAt = notification.ScheduledAt.Format(time.RFC3339)
		}
		if providerStatus != nil {
			item.ProviderStatus = *providerStatus
		}
		if errorMsg != nil {
			item.Error = *errorMsg
		}

		notifications = append(notifications, item)
	}

	// Obtener estadísticas
	statistics, err := uc.getStatistics(ctx, input)
	if err != nil {
		uc.log.Error("Error getting statistics", logger.Error(err))
		// No fallar, solo no incluir estadísticas
		statistics = nil
	}

	// Log auditoría
	uc.log.Info("Admin viewed notification history",
		logger.Int64("admin_id", adminID),
		logger.Int("total_count", int(totalCount)),
		logger.Int("limit", input.Limit),
		logger.Int("offset", input.Offset),
		logger.String("action", "admin_view_notification_history"))

	return &ViewNotificationHistoryOutput{
		Notifications: notifications,
		TotalCount:    int(totalCount),
		Statistics:    statistics,
	}, nil
}

// validateInput valida los datos de entrada
func (uc *ViewNotificationHistoryUseCase) validateInput(input *ViewNotificationHistoryInput) error {
	// Validar type
	if input.Type != nil && *input.Type != "" {
		validTypes := map[string]bool{
			"email":        true,
			"sms":          true,
			"push":         true,
			"announcement": true,
		}
		if !validTypes[*input.Type] {
			return errors.New("VALIDATION_FAILED", "invalid notification type", 400, nil)
		}
	}

	// Validar status
	if input.Status != nil && *input.Status != "" {
		validStatuses := map[string]bool{
			"queued":    true,
			"scheduled": true,
			"sent":      true,
			"failed":    true,
		}
		if !validStatuses[*input.Status] {
			return errors.New("VALIDATION_FAILED", "invalid status", 400, nil)
		}
	}

	// Validar priority
	if input.Priority != nil && *input.Priority != "" {
		validPriorities := map[string]bool{
			"low":      true,
			"normal":   true,
			"high":     true,
			"critical": true,
		}
		if !validPriorities[*input.Priority] {
			return errors.New("VALIDATION_FAILED", "invalid priority", 400, nil)
		}
	}

	// Validar paginación
	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	return nil
}

// getStatistics obtiene estadísticas de notificaciones
func (uc *ViewNotificationHistoryUseCase) getStatistics(ctx context.Context, input *ViewNotificationHistoryInput) (*NotificationStatistics, error) {
	query := uc.db.WithContext(ctx).Table("email_notifications")

	// Aplicar mismos filtros que la query principal (excepto paginación)
	if input.Type != nil && *input.Type != "" {
		query = query.Where("type = ?", *input.Type)
	}
	if input.AdminID != nil {
		query = query.Where("admin_id = ?", *input.AdminID)
	}
	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("created_at >= ?", *input.DateFrom)
	}
	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("created_at <= ?", *input.DateTo+" 23:59:59")
	}

	stats := &NotificationStatistics{}

	// Contar por status
	var statusCounts []struct {
		Status string
		Count  int
	}
	query.Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts)

	for _, sc := range statusCounts {
		switch sc.Status {
		case "sent":
			stats.TotalSent = sc.Count
		case "failed":
			stats.TotalFailed = sc.Count
		case "queued":
			stats.TotalQueued = sc.Count
		case "scheduled":
			stats.TotalScheduled = sc.Count
		}
	}

	// Calcular success rate
	total := stats.TotalSent + stats.TotalFailed
	if total > 0 {
		stats.SuccessRate = float64(stats.TotalSent) / float64(total) * 100
	}

	// Obtener última fecha de envío
	var lastSentAt *time.Time
	query.Where("status = ? AND sent_at IS NOT NULL", "sent").
		Order("sent_at DESC").
		Limit(1).
		Pluck("sent_at", &lastSentAt)

	if lastSentAt != nil {
		stats.LastSentAt = lastSentAt.Format(time.RFC3339)
	}

	// Calcular promedio por día (últimos 30 días)
	var countLast30Days int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	uc.db.WithContext(ctx).Table("email_notifications").
		Where("created_at >= ?", thirtyDaysAgo).
		Count(&countLast30Days)

	if countLast30Days > 0 {
		stats.AveragePerDay = float64(countLast30Days) / 30.0
	}

	return stats, nil
}
