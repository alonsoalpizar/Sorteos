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

// SendEmailInput datos de entrada
type SendEmailInput struct {
	To          []EmailRecipient       `json:"to"`
	CC          []EmailRecipient       `json:"cc,omitempty"`
	BCC         []EmailRecipient       `json:"bcc,omitempty"`
	Subject     string                 `json:"subject"`
	Body        string                 `json:"body"`
	TemplateID  *int64                 `json:"template_id,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	Priority    string                 `json:"priority"`    // low, normal, high
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

// SendEmailOutput resultado
type SendEmailOutput struct {
	NotificationID int64  `json:"notification_id"`
	Status         string `json:"status"` // queued, scheduled, sent, failed
	SentAt         string `json:"sent_at,omitempty"`
	ScheduledAt    string `json:"scheduled_at,omitempty"`
	Recipients     int    `json:"recipients"`
	Message        string `json:"message"`
}

// SendEmailUseCase caso de uso para enviar emails transaccionales
type SendEmailUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewSendEmailUseCase crea una nueva instancia
func NewSendEmailUseCase(db *gorm.DB, log *logger.Logger) *SendEmailUseCase {
	return &SendEmailUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *SendEmailUseCase) Execute(ctx context.Context, input *SendEmailInput, adminID int64) (*SendEmailOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Si tiene template_id, cargar plantilla
	var finalBody string
	var finalSubject string

	if input.TemplateID != nil {
		// TODO: Cargar template desde email_templates table
		// template, err := uc.loadTemplate(ctx, *input.TemplateID)
		// finalBody = uc.renderTemplate(template.Body, input.Variables)
		// finalSubject = uc.renderTemplate(template.Subject, input.Variables)

		// Por ahora usar body/subject proporcionados
		finalBody = input.Body
		finalSubject = input.Subject
	} else {
		finalBody = input.Body
		finalSubject = input.Subject
	}

	// Serializar recipients
	allRecipients := append(input.To, input.CC...)
	allRecipients = append(allRecipients, input.BCC...)

	recipientsJSON, err := json.Marshal(allRecipients)
	if err != nil {
		uc.log.Error("Error marshaling recipients", logger.Error(err))
		return nil, errors.Wrap(errors.ErrInternalServer, err)
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

	// Determinar status inicial
	status := "queued"
	if input.ScheduledAt != nil && input.ScheduledAt.After(time.Now()) {
		status = "scheduled"
	}

	// Crear registro de notificación
	notification := &EmailNotification{
		AdminID:     adminID,
		Type:        "email",
		Recipients:  string(recipientsJSON),
		Subject:     &finalSubject,
		Body:        finalBody,
		TemplateID:  input.TemplateID,
		Variables:   variablesJSON,
		Priority:    input.Priority,
		Status:      status,
		ScheduledAt: input.ScheduledAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Guardar en DB
	result := uc.db.WithContext(ctx).Table("email_notifications").Create(notification)
	if result.Error != nil {
		uc.log.Error("Error creating email notification", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Si es envío inmediato (no scheduled), intentar enviar
	var sentAt *time.Time
	if status == "queued" {
		// TODO: Integrar con email provider (SendGrid, Mailgun, SES, SMTP)
		// err := uc.sendViaProvider(ctx, notification)
		// if err != nil {
		//     notification.Status = "failed"
		//     notification.Error = err.Error()
		// } else {
		//     notification.Status = "sent"
		//     now := time.Now()
		//     notification.SentAt = &now
		//     sentAt = &now
		// }
		//
		// uc.db.WithContext(ctx).Table("email_notifications").Where("id = ?", notification.ID).Updates(notification)

		// Simulación: marcar como sent inmediatamente
		now := time.Now()
		notification.Status = "sent"
		notification.SentAt = &now
		sentAt = &now

		uc.db.WithContext(ctx).Table("email_notifications").
			Where("id = ?", notification.ID).
			Updates(map[string]interface{}{
				"status":  "sent",
				"sent_at": now,
			})

		uc.log.Info("Email notification queued for immediate delivery",
			logger.Int64("notification_id", notification.ID),
			logger.Int64("admin_id", adminID),
			logger.String("subject", finalSubject),
			logger.Int("recipients", len(allRecipients)))
	} else {
		uc.log.Info("Email notification scheduled",
			logger.Int64("notification_id", notification.ID),
			logger.Int64("admin_id", adminID),
			logger.String("subject", finalSubject),
			logger.String("scheduled_at", input.ScheduledAt.Format(time.RFC3339)))
	}

	// Log auditoría
	uc.log.Error("Admin sent email notification",
		logger.Int64("admin_id", adminID),
		logger.Int64("notification_id", notification.ID),
		logger.String("subject", finalSubject),
		logger.Int("recipients", len(allRecipients)),
		logger.String("priority", input.Priority),
		logger.String("status", notification.Status),
		logger.String("action", "admin_send_email"),
		logger.String("severity", "info"))

	// Construir output
	output := &SendEmailOutput{
		NotificationID: notification.ID,
		Status:         notification.Status,
		Recipients:     len(allRecipients),
		Message:        fmt.Sprintf("Email %s successfully", notification.Status),
	}

	if sentAt != nil {
		output.SentAt = sentAt.Format(time.RFC3339)
	}
	if input.ScheduledAt != nil {
		output.ScheduledAt = input.ScheduledAt.Format(time.RFC3339)
	}

	return output, nil
}

// validateInput valida los datos de entrada
func (uc *SendEmailUseCase) validateInput(input *SendEmailInput) error {
	// Validar que haya al menos un destinatario
	if len(input.To) == 0 {
		return errors.New("VALIDATION_FAILED", "at least one recipient is required", 400, nil)
	}

	// Validar emails
	for _, recipient := range input.To {
		if recipient.Email == "" {
			return errors.New("VALIDATION_FAILED", "recipient email cannot be empty", 400, nil)
		}
		// TODO: Validar formato de email con regex
	}

	// Validar subject
	if input.Subject == "" && input.TemplateID == nil {
		return errors.New("VALIDATION_FAILED", "subject is required when not using a template", 400, nil)
	}

	// Validar body
	if input.Body == "" && input.TemplateID == nil {
		return errors.New("VALIDATION_FAILED", "body is required when not using a template", 400, nil)
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

	// Validar scheduled_at (no puede ser en el pasado)
	if input.ScheduledAt != nil && input.ScheduledAt.Before(time.Now()) {
		return errors.New("VALIDATION_FAILED", "scheduled_at cannot be in the past", 400, nil)
	}

	return nil
}

// TODO: Implementar métodos auxiliares
// func (uc *SendEmailUseCase) loadTemplate(ctx context.Context, templateID int64) (*EmailTemplate, error)
// func (uc *SendEmailUseCase) renderTemplate(template string, variables map[string]interface{}) string
// func (uc *SendEmailUseCase) sendViaProvider(ctx context.Context, notification *EmailNotification) error
