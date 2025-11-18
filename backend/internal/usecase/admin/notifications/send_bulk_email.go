package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// SendBulkEmailInput datos de entrada
type SendBulkEmailInput struct {
	Subject     string                 `json:"subject"`
	Body        string                 `json:"body"`
	TemplateID  *int64                 `json:"template_id,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	Segment     string                 `json:"segment"`      // all_users, all_organizers, custom
	Filters     *BulkEmailFilters      `json:"filters,omitempty"`
	Priority    string                 `json:"priority"`     // low, normal, high
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	BatchSize   int                    `json:"batch_size"`   // Tamaño de lote para envío
}

// BulkEmailFilters filtros para segmentación personalizada
type BulkEmailFilters struct {
	Roles          []string   `json:"roles,omitempty"`           // user, organizer, super_admin
	Status         []string   `json:"status,omitempty"`          // active, suspended
	KYCLevels      []string   `json:"kyc_levels,omitempty"`      // unverified, basic, full
	RegisteredFrom *time.Time `json:"registered_from,omitempty"`
	RegisteredTo   *time.Time `json:"registered_to,omitempty"`
	LastLoginFrom  *time.Time `json:"last_login_from,omitempty"`
	LastLoginTo    *time.Time `json:"last_login_to,omitempty"`
	MinRaffles     *int       `json:"min_raffles,omitempty"`     // Para organizadores
	MinRevenue     *float64   `json:"min_revenue,omitempty"`     // Para organizadores
}

// SendBulkEmailOutput resultado
type SendBulkEmailOutput struct {
	BulkNotificationID int64  `json:"bulk_notification_id"`
	Status             string `json:"status"` // queued, scheduled, processing, completed, failed
	TotalRecipients    int    `json:"total_recipients"`
	BatchesCreated     int    `json:"batches_created"`
	EstimatedDuration  int    `json:"estimated_duration_minutes"`
	ScheduledAt        string `json:"scheduled_at,omitempty"`
	Message            string `json:"message"`
}

// BulkEmailNotification registro de notificación masiva
type BulkEmailNotification struct {
	ID                int64
	AdminID           int64
	Subject           string
	Body              string
	TemplateID        *int64
	Variables         *string // JSON
	Segment           string
	Filters           *string // JSON
	Priority          string
	BatchSize         int
	Status            string // queued, scheduled, processing, completed, failed
	TotalRecipients   int
	SuccessfulSent    int
	FailedSent        int
	ScheduledAt       *time.Time
	StartedAt         *time.Time
	CompletedAt       *time.Time
	Error             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// SendBulkEmailUseCase caso de uso para envío masivo de emails
type SendBulkEmailUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewSendBulkEmailUseCase crea una nueva instancia
func NewSendBulkEmailUseCase(db *gorm.DB, log *logger.Logger) *SendBulkEmailUseCase {
	return &SendBulkEmailUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *SendBulkEmailUseCase) Execute(ctx context.Context, input *SendBulkEmailInput, adminID int64) (*SendBulkEmailOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Obtener lista de destinatarios según segmento y filtros
	recipients, err := uc.getRecipients(ctx, input.Segment, input.Filters)
	if err != nil {
		uc.log.Error("Error getting recipients", logger.Error(err))
		return nil, err
	}

	if len(recipients) == 0 {
		return nil, errors.New("VALIDATION_FAILED", "no recipients found for the specified segment and filters", 400, nil)
	}

	// Serializar variables
	var variablesJSON *string
	if input.Variables != nil {
		varsBytes, err := json.Marshal(input.Variables)
		if err != nil {
			uc.log.Error("Error marshaling variables", logger.Error(err))
			return nil, errors.Wrap(errors.ErrInternalServer, err)
		}
		varsStr := string(varsBytes)
		variablesJSON = &varsStr
	}

	// Serializar filters
	var filtersJSON *string
	if input.Filters != nil {
		filtersBytes, err := json.Marshal(input.Filters)
		if err != nil {
			uc.log.Error("Error marshaling filters", logger.Error(err))
			return nil, errors.Wrap(errors.ErrInternalServer, err)
		}
		filtersStr := string(filtersBytes)
		filtersJSON = &filtersStr
	}

	// Determinar status inicial
	status := "queued"
	if input.ScheduledAt != nil && input.ScheduledAt.After(time.Now()) {
		status = "scheduled"
	}

	// Calcular número de batches
	batchSize := input.BatchSize
	if batchSize == 0 {
		batchSize = 100 // Default
	}
	batchesCount := (len(recipients) + batchSize - 1) / batchSize

	// Crear registro de bulk notification
	bulkNotification := &BulkEmailNotification{
		AdminID:         adminID,
		Subject:         input.Subject,
		Body:            input.Body,
		TemplateID:      input.TemplateID,
		Variables:       variablesJSON,
		Segment:         input.Segment,
		Filters:         filtersJSON,
		Priority:        input.Priority,
		BatchSize:       batchSize,
		Status:          status,
		TotalRecipients: len(recipients),
		SuccessfulSent:  0,
		FailedSent:      0,
		ScheduledAt:     input.ScheduledAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Guardar en DB
	result := uc.db.WithContext(ctx).Table("bulk_email_notifications").Create(bulkNotification)
	if result.Error != nil {
		uc.log.Error("Error creating bulk email notification", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// TODO: Crear batches individuales en email_notifications table
	// for i := 0; i < len(recipients); i += batchSize {
	//     end := i + batchSize
	//     if end > len(recipients) {
	//         end = len(recipients)
	//     }
	//     batch := recipients[i:end]
	//
	//     // Crear EmailNotification para este batch
	//     uc.createBatchNotification(ctx, bulkNotification.ID, batch, input)
	// }

	// TODO: Si es envío inmediato, iniciar procesamiento en background
	// if status == "queued" {
	//     go uc.processBulkEmail(bulkNotification.ID)
	// }

	// Estimación de duración (100 emails por minuto)
	estimatedMinutes := (len(recipients) + 99) / 100

	// Log auditoría crítica (envío masivo)
	uc.log.Error("Admin created bulk email notification",
		logger.Int64("admin_id", adminID),
		logger.Int64("bulk_notification_id", bulkNotification.ID),
		logger.String("subject", input.Subject),
		logger.String("segment", input.Segment),
		logger.Int("total_recipients", len(recipients)),
		logger.Int("batches", batchesCount),
		logger.String("status", status),
		logger.String("action", "admin_send_bulk_email"),
		logger.String("severity", "warning"))

	// Construir output
	output := &SendBulkEmailOutput{
		BulkNotificationID: bulkNotification.ID,
		Status:             status,
		TotalRecipients:    len(recipients),
		BatchesCreated:     batchesCount,
		EstimatedDuration:  estimatedMinutes,
		Message:            fmt.Sprintf("Bulk email %s with %d recipients in %d batches", status, len(recipients), batchesCount),
	}

	if input.ScheduledAt != nil {
		output.ScheduledAt = input.ScheduledAt.Format(time.RFC3339)
	}

	return output, nil
}

// validateInput valida los datos de entrada
func (uc *SendBulkEmailUseCase) validateInput(input *SendBulkEmailInput) error {
	// Validar subject
	if input.Subject == "" && input.TemplateID == nil {
		return errors.New("VALIDATION_FAILED", "subject is required when not using a template", 400, nil)
	}

	// Validar body
	if input.Body == "" && input.TemplateID == nil {
		return errors.New("VALIDATION_FAILED", "body is required when not using a template", 400, nil)
	}

	// Validar segment
	validSegments := map[string]bool{
		"all_users":      true,
		"all_organizers": true,
		"custom":         true,
	}
	if !validSegments[input.Segment] {
		return errors.New("VALIDATION_FAILED", "segment must be one of: all_users, all_organizers, custom", 400, nil)
	}

	// Si es custom, filters es requerido
	if input.Segment == "custom" && input.Filters == nil {
		return errors.New("VALIDATION_FAILED", "filters are required when segment is custom", 400, nil)
	}

	// Validar priority
	if input.Priority == "" {
		input.Priority = "normal"
	}
	validPriorities := map[string]bool{
		"low":    true,
		"normal": true,
		"high":   true,
	}
	if !validPriorities[input.Priority] {
		return errors.New("VALIDATION_FAILED", "priority must be one of: low, normal, high", 400, nil)
	}

	// Validar batch_size
	if input.BatchSize < 0 {
		return errors.New("VALIDATION_FAILED", "batch_size must be positive", 400, nil)
	}
	if input.BatchSize > 1000 {
		return errors.New("VALIDATION_FAILED", "batch_size cannot exceed 1000", 400, nil)
	}

	// Validar scheduled_at
	if input.ScheduledAt != nil && input.ScheduledAt.Before(time.Now()) {
		return errors.New("VALIDATION_FAILED", "scheduled_at cannot be in the past", 400, nil)
	}

	return nil
}

// getRecipients obtiene la lista de destinatarios según segmento y filtros
func (uc *SendBulkEmailUseCase) getRecipients(ctx context.Context, segment string, filters *BulkEmailFilters) ([]EmailRecipient, error) {
	var recipients []EmailRecipient
	query := uc.db.WithContext(ctx).Table("users")

	switch segment {
	case "all_users":
		query = query.Where("role = ? AND status = ?", "user", "active")

	case "all_organizers":
		query = query.Where("role = ? AND status = ?", "organizer", "active")

	case "custom":
		if filters == nil {
			return nil, errors.New("VALIDATION_FAILED", "filters required for custom segment", 400, nil)
		}

		// Aplicar filtros
		if len(filters.Roles) > 0 {
			query = query.Where("role IN ?", filters.Roles)
		}
		if len(filters.Status) > 0 {
			query = query.Where("status IN ?", filters.Status)
		}
		if len(filters.KYCLevels) > 0 {
			query = query.Where("kyc_level IN ?", filters.KYCLevels)
		}
		if filters.RegisteredFrom != nil {
			query = query.Where("created_at >= ?", filters.RegisteredFrom)
		}
		if filters.RegisteredTo != nil {
			query = query.Where("created_at <= ?", filters.RegisteredTo)
		}
		if filters.LastLoginFrom != nil {
			query = query.Where("last_login_at >= ?", filters.LastLoginFrom)
		}
		if filters.LastLoginTo != nil {
			query = query.Where("last_login_at <= ?", filters.LastLoginTo)
		}

		// TODO: Filtros para organizadores (min_raffles, min_revenue) requieren JOIN
		// if filters.MinRaffles != nil {
		//     query = query.Joins("LEFT JOIN organizer_profiles ON users.id = organizer_profiles.user_id")
		//     query = query.Where("organizer_profiles.total_raffles >= ?", *filters.MinRaffles)
		// }
	}

	// Seleccionar email y nombre
	rows, err := query.Select("email, first_name, last_name").Rows()
	if err != nil {
		uc.log.Error("Error querying recipients", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		var firstName, lastName *string
		if err := rows.Scan(&email, &firstName, &lastName); err != nil {
			continue
		}

		name := ""
		if firstName != nil && lastName != nil {
			name = fmt.Sprintf("%s %s", *firstName, *lastName)
		} else if firstName != nil {
			name = *firstName
		}

		recipients = append(recipients, EmailRecipient{
			Email: email,
			Name:  name,
		})
	}

	return recipients, nil
}

// TODO: Implementar métodos auxiliares
// func (uc *SendBulkEmailUseCase) createBatchNotification(ctx context.Context, bulkID int64, batch []EmailRecipient, input *SendBulkEmailInput) error
// func (uc *SendBulkEmailUseCase) processBulkEmail(bulkID int64)
