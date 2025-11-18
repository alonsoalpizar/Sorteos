package notifications

import (
	"encoding/json"
	"time"
)

// EmailRecipient destinatario del email
type EmailRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// EmailNotification registro de notificación en DB
type EmailNotification struct {
	ID             int64
	AdminID        int64
	Type           string // email, sms, push
	Recipients     json.RawMessage // JSONB array
	Subject        *string
	Body           string
	TemplateID     *int64
	Variables      *json.RawMessage // JSONB object
	Priority       string
	Status         string // queued, scheduled, sent, failed
	SentAt         *time.Time
	ScheduledAt    *time.Time
	ProviderID     *string // Email ID del proveedor (SendGrid, Mailgun, etc.)
	ProviderStatus *string
	Error          *string
	Metadata       *json.RawMessage // JSONB object
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
