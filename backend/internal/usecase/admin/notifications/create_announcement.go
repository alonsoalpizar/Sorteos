package notifications

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// CreateAnnouncementInput datos de entrada
type CreateAnnouncementInput struct {
	Title       string     `json:"title"`
	Message     string     `json:"message"`
	Type        string     `json:"type"`         // info, warning, maintenance, feature, promotion
	Priority    string     `json:"priority"`     // low, normal, high, critical
	Target      string     `json:"target"`       // all, users, organizers, specific_users
	TargetIDs   []int64    `json:"target_ids,omitempty"` // IDs si target es specific_users
	URL         *string    `json:"url,omitempty"`        // Link a más información
	ActionLabel *string    `json:"action_label,omitempty"` // Texto del botón de acción
	ActionURL   *string    `json:"action_url,omitempty"`   // URL del botón de acción
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`   // Fecha de expiración
	PublishedAt *time.Time `json:"published_at,omitempty"` // Programar publicación
}

// CreateAnnouncementOutput resultado
type CreateAnnouncementOutput struct {
	AnnouncementID int64  `json:"announcement_id"`
	Status         string `json:"status"` // draft, scheduled, published, expired
	PublishedAt    string `json:"published_at,omitempty"`
	ExpiresAt      string `json:"expires_at,omitempty"`
	TargetUsers    int    `json:"target_users"` // Número de usuarios objetivo
	Message        string `json:"message"`
}

// Announcement modelo de anuncio
type Announcement struct {
	ID          int64
	AdminID     int64
	Title       string
	Message     string
	Type        string
	Priority    string
	Target      string
	TargetIDs   *string // JSON array de IDs
	URL         *string
	ActionLabel *string
	ActionURL   *string
	Status      string // draft, scheduled, published, expired
	ViewCount   int
	ClickCount  int
	PublishedAt *time.Time
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// CreateAnnouncementUseCase caso de uso para crear anuncios
type CreateAnnouncementUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCreateAnnouncementUseCase crea una nueva instancia
func NewCreateAnnouncementUseCase(db *gorm.DB, log *logger.Logger) *CreateAnnouncementUseCase {
	return &CreateAnnouncementUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *CreateAnnouncementUseCase) Execute(ctx context.Context, input *CreateAnnouncementInput, adminID int64) (*CreateAnnouncementOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Determinar status inicial
	status := "published"
	if input.PublishedAt != nil && input.PublishedAt.After(time.Now()) {
		status = "scheduled"
	}

	// Si tiene published_at en el pasado, usar now
	publishedAt := time.Now()
	if input.PublishedAt != nil && !input.PublishedAt.After(time.Now()) {
		publishedAt = *input.PublishedAt
	} else if status == "scheduled" {
		publishedAt = *input.PublishedAt
	}

	// Serializar target_ids si existen
	var targetIDsJSON *string
	if len(input.TargetIDs) > 0 {
		// Simple JSON serialization de array de ints
		// TODO: Usar json.Marshal para producción
		targetIDsJSON = new(string)
		*targetIDsJSON = "[]" // Placeholder
	}

	// Crear anuncio
	announcement := &Announcement{
		AdminID:     adminID,
		Title:       input.Title,
		Message:     input.Message,
		Type:        input.Type,
		Priority:    input.Priority,
		Target:      input.Target,
		TargetIDs:   targetIDsJSON,
		URL:         input.URL,
		ActionLabel: input.ActionLabel,
		ActionURL:   input.ActionURL,
		Status:      status,
		ViewCount:   0,
		ClickCount:  0,
		PublishedAt: &publishedAt,
		ExpiresAt:   input.ExpiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Guardar en DB
	result := uc.db.WithContext(ctx).Table("announcements").Create(announcement)
	if result.Error != nil {
		uc.log.Error("Error creating announcement", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Calcular número de usuarios objetivo
	targetUsers, err := uc.calculateTargetUsers(ctx, input.Target, input.TargetIDs)
	if err != nil {
		uc.log.Error("Error calculating target users", logger.Error(err))
		// No fallar, solo loguear
		targetUsers = 0
	}

	// Log auditoría
	uc.log.Error("Admin created announcement",
		logger.Int64("admin_id", adminID),
		logger.Int64("announcement_id", announcement.ID),
		logger.String("title", input.Title),
		logger.String("type", input.Type),
		logger.String("priority", input.Priority),
		logger.String("target", input.Target),
		logger.Int("target_users", targetUsers),
		logger.String("status", status),
		logger.String("action", "admin_create_announcement"),
		logger.String("severity", "info"))

	// TODO: Si es published y no scheduled, notificar usuarios en tiempo real
	// if status == "published" {
	//     go uc.notifyTargetUsers(announcement.ID, targetUsers)
	// }

	// Construir output
	output := &CreateAnnouncementOutput{
		AnnouncementID: announcement.ID,
		Status:         status,
		TargetUsers:    targetUsers,
		Message:        "Announcement created successfully",
	}

	if announcement.PublishedAt != nil {
		output.PublishedAt = announcement.PublishedAt.Format(time.RFC3339)
	}
	if input.ExpiresAt != nil {
		output.ExpiresAt = input.ExpiresAt.Format(time.RFC3339)
	}

	return output, nil
}

// validateInput valida los datos de entrada
func (uc *CreateAnnouncementUseCase) validateInput(input *CreateAnnouncementInput) error {
	// Validar title
	if input.Title == "" {
		return errors.New("VALIDATION_FAILED", "title is required", 400, nil)
	}
	if len(input.Title) > 200 {
		return errors.New("VALIDATION_FAILED", "title cannot exceed 200 characters", 400, nil)
	}

	// Validar message
	if input.Message == "" {
		return errors.New("VALIDATION_FAILED", "message is required", 400, nil)
	}
	if len(input.Message) > 5000 {
		return errors.New("VALIDATION_FAILED", "message cannot exceed 5000 characters", 400, nil)
	}

	// Validar type
	validTypes := map[string]bool{
		"info":        true,
		"warning":     true,
		"maintenance": true,
		"feature":     true,
		"promotion":   true,
	}
	if !validTypes[input.Type] {
		return errors.New("VALIDATION_FAILED", "type must be one of: info, warning, maintenance, feature, promotion", 400, nil)
	}

	// Validar priority
	validPriorities := map[string]bool{
		"low":      true,
		"normal":   true,
		"high":     true,
		"critical": true,
	}
	if !validPriorities[input.Priority] {
		return errors.New("VALIDATION_FAILED", "priority must be one of: low, normal, high, critical", 400, nil)
	}

	// Validar target
	validTargets := map[string]bool{
		"all":            true,
		"users":          true,
		"organizers":     true,
		"specific_users": true,
	}
	if !validTargets[input.Target] {
		return errors.New("VALIDATION_FAILED", "target must be one of: all, users, organizers, specific_users", 400, nil)
	}

	// Si target es specific_users, target_ids es requerido
	if input.Target == "specific_users" && len(input.TargetIDs) == 0 {
		return errors.New("VALIDATION_FAILED", "target_ids is required when target is specific_users", 400, nil)
	}

	// Validar action_url requiere action_label
	if input.ActionURL != nil && input.ActionLabel == nil {
		return errors.New("VALIDATION_FAILED", "action_label is required when action_url is provided", 400, nil)
	}

	// Validar expires_at no puede ser en el pasado
	if input.ExpiresAt != nil && input.ExpiresAt.Before(time.Now()) {
		return errors.New("VALIDATION_FAILED", "expires_at cannot be in the past", 400, nil)
	}

	// Validar expires_at debe ser después de published_at
	if input.ExpiresAt != nil && input.PublishedAt != nil && input.ExpiresAt.Before(*input.PublishedAt) {
		return errors.New("VALIDATION_FAILED", "expires_at must be after published_at", 400, nil)
	}

	return nil
}

// calculateTargetUsers calcula el número de usuarios que recibirán el anuncio
func (uc *CreateAnnouncementUseCase) calculateTargetUsers(ctx context.Context, target string, targetIDs []int64) (int, error) {
	var count int64

	switch target {
	case "all":
		err := uc.db.WithContext(ctx).Table("users").Where("status = ?", "active").Count(&count).Error
		if err != nil {
			return 0, err
		}

	case "users":
		err := uc.db.WithContext(ctx).Table("users").Where("role = ? AND status = ?", "user", "active").Count(&count).Error
		if err != nil {
			return 0, err
		}

	case "organizers":
		err := uc.db.WithContext(ctx).Table("users").Where("role = ? AND status = ?", "organizer", "active").Count(&count).Error
		if err != nil {
			return 0, err
		}

	case "specific_users":
		count = int64(len(targetIDs))
	}

	return int(count), nil
}

// TODO: Implementar métodos auxiliares
// func (uc *CreateAnnouncementUseCase) notifyTargetUsers(announcementID int64, targetCount int)
