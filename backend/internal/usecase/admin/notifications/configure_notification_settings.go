package notifications

import (
	"context"
	"encoding/json"

	"github.com/sorteos-platform/backend/internal/repository"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// ConfigureNotificationSettingsInput datos de entrada
type ConfigureNotificationSettingsInput struct {
	Operation string      `json:"operation"` // get, update
	Settings  *NotificationSettingsData `json:"settings,omitempty"`
}

// NotificationSettingsData configuración de notificaciones
type NotificationSettingsData struct {
	EmailProvider       *string                `json:"email_provider,omitempty"`        // smtp, sendgrid, mailgun, ses
	SMTPConfig          *SMTPConfig            `json:"smtp_config,omitempty"`
	SendGridConfig      *SendGridConfig        `json:"sendgrid_config,omitempty"`
	MailgunConfig       *MailgunConfig         `json:"mailgun_config,omitempty"`
	SESConfig           *SESConfig             `json:"ses_config,omitempty"`
	DefaultFromEmail    *string                `json:"default_from_email,omitempty"`
	DefaultFromName     *string                `json:"default_from_name,omitempty"`
	ReplyToEmail        *string                `json:"reply_to_email,omitempty"`
	EnableEmailQueue    *bool                  `json:"enable_email_queue,omitempty"`    // Usar cola o envío directo
	MaxRetries          *int                   `json:"max_retries,omitempty"`           // Reintentos para emails fallidos
	RetryDelay          *int                   `json:"retry_delay_minutes,omitempty"`   // Minutos entre reintentos
	BatchSize           *int                   `json:"batch_size,omitempty"`            // Tamaño de lote para bulk emails
	RateLimitPerHour    *int                   `json:"rate_limit_per_hour,omitempty"`   // Límite de emails por hora
	EnableTracking      *bool                  `json:"enable_tracking,omitempty"`       // Tracking de aperturas/clicks
	EnableSMSNotif      *bool                  `json:"enable_sms_notifications,omitempty"`
	EnablePushNotif     *bool                  `json:"enable_push_notifications,omitempty"`
	MaintenanceModeNotif *bool                 `json:"maintenance_mode_notifications,omitempty"` // Deshabilitar todas las notifs
}

// SMTPConfig configuración SMTP
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"` // TODO: Encriptar en producción
	UseTLS   bool   `json:"use_tls"`
}

// SendGridConfig configuración SendGrid
type SendGridConfig struct {
	APIKey string `json:"api_key"` // TODO: Encriptar
}

// MailgunConfig configuración Mailgun
type MailgunConfig struct {
	Domain string `json:"domain"`
	APIKey string `json:"api_key"` // TODO: Encriptar
}

// SESConfig configuración AWS SES
type SESConfig struct {
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`     // TODO: Encriptar
	SecretAccessKey string `json:"secret_access_key"` // TODO: Encriptar
}

// ConfigureNotificationSettingsOutput resultado
type ConfigureNotificationSettingsOutput struct {
	Operation string                    `json:"operation"`
	Settings  *NotificationSettingsData `json:"settings,omitempty"`
	Message   string                    `json:"message"`
}

// ConfigureNotificationSettingsUseCase caso de uso para configurar ajustes de notificaciones
type ConfigureNotificationSettingsUseCase struct {
	configRepo repository.SystemConfigRepository
	log        *logger.Logger
}

// NewConfigureNotificationSettingsUseCase crea una nueva instancia
func NewConfigureNotificationSettingsUseCase(configRepo repository.SystemConfigRepository, log *logger.Logger) *ConfigureNotificationSettingsUseCase {
	return &ConfigureNotificationSettingsUseCase{
		configRepo: configRepo,
		log:        log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ConfigureNotificationSettingsUseCase) Execute(ctx context.Context, input *ConfigureNotificationSettingsInput, adminID int64) (*ConfigureNotificationSettingsOutput, error) {
	// Validar operación
	if input.Operation != "get" && input.Operation != "update" {
		return nil, errors.New("VALIDATION_FAILED", "operation must be 'get' or 'update'", 400, nil)
	}

	if input.Operation == "get" {
		return uc.getSettings(ctx, adminID)
	} else {
		return uc.updateSettings(ctx, input, adminID)
	}
}

// getSettings obtiene la configuración actual
func (uc *ConfigureNotificationSettingsUseCase) getSettings(ctx context.Context, adminID int64) (*ConfigureNotificationSettingsOutput, error) {
	// Obtener todas las configuraciones de categoría "notification"
	configs, err := uc.configRepo.GetByCategory(ctx, "notification")
	if err != nil {
		uc.log.Error("Error getting notification settings", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Parsear configuraciones a estructura
	settings := &NotificationSettingsData{}

	for _, config := range configs {
		var value interface{}
		json.Unmarshal([]byte(config.Value), &value)

		switch config.Key {
		case "email_provider":
			if strVal, ok := value.(string); ok {
				settings.EmailProvider = &strVal
			}
		case "default_from_email":
			if strVal, ok := value.(string); ok {
				settings.DefaultFromEmail = &strVal
			}
		case "default_from_name":
			if strVal, ok := value.(string); ok {
				settings.DefaultFromName = &strVal
			}
		case "reply_to_email":
			if strVal, ok := value.(string); ok {
				settings.ReplyToEmail = &strVal
			}
		case "enable_email_queue":
			if boolVal, ok := value.(bool); ok {
				settings.EnableEmailQueue = &boolVal
			}
		case "max_retries":
			if floatVal, ok := value.(float64); ok {
				intVal := int(floatVal)
				settings.MaxRetries = &intVal
			}
		case "batch_size":
			if floatVal, ok := value.(float64); ok {
				intVal := int(floatVal)
				settings.BatchSize = &intVal
			}
		case "rate_limit_per_hour":
			if floatVal, ok := value.(float64); ok {
				intVal := int(floatVal)
				settings.RateLimitPerHour = &intVal
			}
		case "enable_tracking":
			if boolVal, ok := value.(bool); ok {
				settings.EnableTracking = &boolVal
			}
		case "enable_sms_notifications":
			if boolVal, ok := value.(bool); ok {
				settings.EnableSMSNotif = &boolVal
			}
		case "enable_push_notifications":
			if boolVal, ok := value.(bool); ok {
				settings.EnablePushNotif = &boolVal
			}
		case "maintenance_mode_notifications":
			if boolVal, ok := value.(bool); ok {
				settings.MaintenanceModeNotif = &boolVal
			}
		// TODO: Parsear configs de SMTP, SendGrid, Mailgun, SES
		}
	}

	return &ConfigureNotificationSettingsOutput{
		Operation: "get",
		Settings:  settings,
		Message:   "Notification settings retrieved successfully",
	}, nil
}

// updateSettings actualiza la configuración
func (uc *ConfigureNotificationSettingsUseCase) updateSettings(ctx context.Context, input *ConfigureNotificationSettingsInput, adminID int64) (*ConfigureNotificationSettingsOutput, error) {
	if input.Settings == nil {
		return nil, errors.New("VALIDATION_FAILED", "settings are required for update", 400, nil)
	}

	settings := input.Settings

	// Validar y actualizar cada configuración
	if settings.EmailProvider != nil {
		if err := uc.validateEmailProvider(*settings.EmailProvider); err != nil {
			return nil, err
		}
		valueJSON, _ := json.Marshal(*settings.EmailProvider)
		uc.configRepo.Set(ctx, "email_provider", string(valueJSON), "notification", adminID)
	}

	if settings.DefaultFromEmail != nil {
		// TODO: Validar formato de email
		valueJSON, _ := json.Marshal(*settings.DefaultFromEmail)
		uc.configRepo.Set(ctx, "default_from_email", string(valueJSON), "notification", adminID)
	}

	if settings.DefaultFromName != nil {
		valueJSON, _ := json.Marshal(*settings.DefaultFromName)
		uc.configRepo.Set(ctx, "default_from_name", string(valueJSON), "notification", adminID)
	}

	if settings.ReplyToEmail != nil {
		// TODO: Validar formato de email
		valueJSON, _ := json.Marshal(*settings.ReplyToEmail)
		uc.configRepo.Set(ctx, "reply_to_email", string(valueJSON), "notification", adminID)
	}

	if settings.EnableEmailQueue != nil {
		valueJSON, _ := json.Marshal(*settings.EnableEmailQueue)
		uc.configRepo.Set(ctx, "enable_email_queue", string(valueJSON), "notification", adminID)
	}

	if settings.MaxRetries != nil {
		if *settings.MaxRetries < 0 || *settings.MaxRetries > 10 {
			return nil, errors.New("VALIDATION_FAILED", "max_retries must be between 0 and 10", 400, nil)
		}
		valueJSON, _ := json.Marshal(*settings.MaxRetries)
		uc.configRepo.Set(ctx, "max_retries", string(valueJSON), "notification", adminID)
	}

	if settings.RetryDelay != nil {
		if *settings.RetryDelay < 1 || *settings.RetryDelay > 1440 {
			return nil, errors.New("VALIDATION_FAILED", "retry_delay must be between 1 and 1440 minutes", 400, nil)
		}
		valueJSON, _ := json.Marshal(*settings.RetryDelay)
		uc.configRepo.Set(ctx, "retry_delay_minutes", string(valueJSON), "notification", adminID)
	}

	if settings.BatchSize != nil {
		if *settings.BatchSize < 1 || *settings.BatchSize > 1000 {
			return nil, errors.New("VALIDATION_FAILED", "batch_size must be between 1 and 1000", 400, nil)
		}
		valueJSON, _ := json.Marshal(*settings.BatchSize)
		uc.configRepo.Set(ctx, "batch_size", string(valueJSON), "notification", adminID)
	}

	if settings.RateLimitPerHour != nil {
		if *settings.RateLimitPerHour < 1 || *settings.RateLimitPerHour > 100000 {
			return nil, errors.New("VALIDATION_FAILED", "rate_limit_per_hour must be between 1 and 100000", 400, nil)
		}
		valueJSON, _ := json.Marshal(*settings.RateLimitPerHour)
		uc.configRepo.Set(ctx, "rate_limit_per_hour", string(valueJSON), "notification", adminID)
	}

	if settings.EnableTracking != nil {
		valueJSON, _ := json.Marshal(*settings.EnableTracking)
		uc.configRepo.Set(ctx, "enable_tracking", string(valueJSON), "notification", adminID)
	}

	if settings.EnableSMSNotif != nil {
		valueJSON, _ := json.Marshal(*settings.EnableSMSNotif)
		uc.configRepo.Set(ctx, "enable_sms_notifications", string(valueJSON), "notification", adminID)
	}

	if settings.EnablePushNotif != nil {
		valueJSON, _ := json.Marshal(*settings.EnablePushNotif)
		uc.configRepo.Set(ctx, "enable_push_notifications", string(valueJSON), "notification", adminID)
	}

	if settings.MaintenanceModeNotif != nil {
		valueJSON, _ := json.Marshal(*settings.MaintenanceModeNotif)
		uc.configRepo.Set(ctx, "maintenance_mode_notifications", string(valueJSON), "notification", adminID)
	}

	// TODO: Guardar configs de SMTP, SendGrid, Mailgun, SES
	// Estas deben ser encriptadas antes de guardar en producción

	// Log auditoría crítica
	uc.log.Error("Admin updated notification settings",
		logger.Int64("admin_id", adminID),
		logger.String("action", "admin_configure_notification_settings"),
		logger.String("severity", "critical"))

	// Recargar settings
	return uc.getSettings(ctx, adminID)
}

// validateEmailProvider valida el proveedor de email
func (uc *ConfigureNotificationSettingsUseCase) validateEmailProvider(provider string) error {
	validProviders := map[string]bool{
		"smtp":     true,
		"sendgrid": true,
		"mailgun":  true,
		"ses":      true,
	}

	if !validProviders[provider] {
		return errors.New("VALIDATION_FAILED", "email_provider must be one of: smtp, sendgrid, mailgun, ses", 400, nil)
	}

	return nil
}
