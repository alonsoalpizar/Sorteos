package notifications

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ManageEmailTemplatesInput datos de entrada
type ManageEmailTemplatesInput struct {
	Operation   string                 `json:"operation"` // create, update, delete, get, list
	TemplateID  *int64                 `json:"template_id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Subject     string                 `json:"subject,omitempty"`
	Body        string                 `json:"body,omitempty"`
	Variables   []string               `json:"variables,omitempty"`   // Lista de variables disponibles
	Category    string                 `json:"category,omitempty"`    // transactional, marketing, system
	Description string                 `json:"description,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
}

// ManageEmailTemplatesOutput resultado
type ManageEmailTemplatesOutput struct {
	Operation string           `json:"operation"`
	Template  *EmailTemplate   `json:"template,omitempty"`
	Templates []*EmailTemplate `json:"templates,omitempty"`
	Message   string           `json:"message"`
}

// EmailTemplate modelo de plantilla de email
type EmailTemplate struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	Variables   *string   `json:"variables"` // JSON array de variables
	Category    string    `json:"category"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	UsageCount  int       `json:"usage_count"`
	CreatedBy   int64     `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// ManageEmailTemplatesUseCase caso de uso para gestionar plantillas de email
type ManageEmailTemplatesUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewManageEmailTemplatesUseCase crea una nueva instancia
func NewManageEmailTemplatesUseCase(db *gorm.DB, log *logger.Logger) *ManageEmailTemplatesUseCase {
	return &ManageEmailTemplatesUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ManageEmailTemplatesUseCase) Execute(ctx context.Context, input *ManageEmailTemplatesInput, adminID int64) (*ManageEmailTemplatesOutput, error) {
	// Validar operación
	validOperations := map[string]bool{
		"create": true,
		"update": true,
		"delete": true,
		"get":    true,
		"list":   true,
	}
	if !validOperations[input.Operation] {
		return nil, errors.New("VALIDATION_FAILED", "operation must be one of: create, update, delete, get, list", 400, nil)
	}

	// Ejecutar operación
	switch input.Operation {
	case "create":
		return uc.createTemplate(ctx, input, adminID)
	case "update":
		return uc.updateTemplate(ctx, input, adminID)
	case "delete":
		return uc.deleteTemplate(ctx, input, adminID)
	case "get":
		return uc.getTemplate(ctx, input, adminID)
	case "list":
		return uc.listTemplates(ctx, input, adminID)
	default:
		return nil, errors.New("VALIDATION_FAILED", "invalid operation", 400, nil)
	}
}

// createTemplate crea una nueva plantilla
func (uc *ManageEmailTemplatesUseCase) createTemplate(ctx context.Context, input *ManageEmailTemplatesInput, adminID int64) (*ManageEmailTemplatesOutput, error) {
	// Validar inputs
	if input.Name == "" {
		return nil, errors.New("VALIDATION_FAILED", "name is required", 400, nil)
	}
	if input.Subject == "" {
		return nil, errors.New("VALIDATION_FAILED", "subject is required", 400, nil)
	}
	if input.Body == "" {
		return nil, errors.New("VALIDATION_FAILED", "body is required", 400, nil)
	}
	if input.Category == "" {
		return nil, errors.New("VALIDATION_FAILED", "category is required", 400, nil)
	}

	// Validar category
	validCategories := map[string]bool{
		"transactional": true,
		"marketing":     true,
		"system":        true,
	}
	if !validCategories[input.Category] {
		return nil, errors.New("VALIDATION_FAILED", "category must be one of: transactional, marketing, system", 400, nil)
	}

	// Validar que el nombre no exista
	var existingCount int64
	uc.db.WithContext(ctx).Table("email_templates").Where("name = ? AND deleted_at IS NULL", input.Name).Count(&existingCount)
	if existingCount > 0 {
		return nil, errors.New("VALIDATION_FAILED", "template with this name already exists", 409, nil)
	}

	// Extraer variables del body
	variables := uc.extractVariables(input.Body)
	if len(input.Variables) > 0 {
		// Merge con variables proporcionadas
		for _, v := range input.Variables {
			if !contains(variables, v) {
				variables = append(variables, v)
			}
		}
	}

	// Serializar variables
	var variablesJSON *string
	if len(variables) > 0 {
		varsBytes, err := json.Marshal(variables)
		if err != nil {
			uc.log.Error("Error marshaling variables", logger.Error(err))
			return nil, errors.Wrap(errors.ErrInternalServer, err)
		}
		varsStr := string(varsBytes)
		variablesJSON = &varsStr
	}

	// Crear template
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	template := &EmailTemplate{
		Name:        input.Name,
		Subject:     input.Subject,
		Body:        input.Body,
		Variables:   variablesJSON,
		Category:    input.Category,
		Description: input.Description,
		IsActive:    isActive,
		UsageCount:  0,
		CreatedBy:   adminID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Guardar en DB
	result := uc.db.WithContext(ctx).Table("email_templates").Create(template)
	if result.Error != nil {
		uc.log.Error("Error creating email template", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log auditoría
	uc.log.Info("Admin created email template",
		logger.Int64("admin_id", adminID),
		logger.Int64("template_id", template.ID),
		logger.String("name", template.Name),
		logger.String("category", template.Category),
		logger.String("action", "admin_create_email_template"))

	return &ManageEmailTemplatesOutput{
		Operation: "create",
		Template:  template,
		Message:   "Email template created successfully",
	}, nil
}

// updateTemplate actualiza una plantilla existente
func (uc *ManageEmailTemplatesUseCase) updateTemplate(ctx context.Context, input *ManageEmailTemplatesInput, adminID int64) (*ManageEmailTemplatesOutput, error) {
	// Validar template_id
	if input.TemplateID == nil {
		return nil, errors.New("VALIDATION_FAILED", "template_id is required for update", 400, nil)
	}

	// Buscar template
	var template EmailTemplate
	result := uc.db.WithContext(ctx).Table("email_templates").Where("id = ? AND deleted_at IS NULL", *input.TemplateID).First(&template)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("TEMPLATE_NOT_FOUND", "email template not found", 404, nil)
		}
		uc.log.Error("Error finding email template", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Actualizar campos
	updates := make(map[string]interface{})

	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Subject != "" {
		updates["subject"] = input.Subject
		// Si cambia el subject, también cambiamos el body probablemente
	}
	if input.Body != "" {
		updates["body"] = input.Body
		// Re-extraer variables
		variables := uc.extractVariables(input.Body)
		if len(variables) > 0 {
			varsBytes, _ := json.Marshal(variables)
			updates["variables"] = string(varsBytes)
		}
	}
	if input.Category != "" {
		validCategories := map[string]bool{
			"transactional": true,
			"marketing":     true,
			"system":        true,
		}
		if !validCategories[input.Category] {
			return nil, errors.New("VALIDATION_FAILED", "invalid category", 400, nil)
		}
		updates["category"] = input.Category
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	updates["updated_by"] = adminID
	updates["updated_at"] = time.Now()

	// Actualizar en DB
	result = uc.db.WithContext(ctx).Table("email_templates").Where("id = ?", template.ID).Updates(updates)
	if result.Error != nil {
		uc.log.Error("Error updating email template", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Recargar template
	uc.db.WithContext(ctx).Table("email_templates").Where("id = ?", template.ID).First(&template)

	// Log auditoría
	uc.log.Error("Admin updated email template",
		logger.Int64("admin_id", adminID),
		logger.Int64("template_id", template.ID),
		logger.String("name", template.Name),
		logger.String("action", "admin_update_email_template"),
		logger.String("severity", "info"))

	return &ManageEmailTemplatesOutput{
		Operation: "update",
		Template:  &template,
		Message:   "Email template updated successfully",
	}, nil
}

// deleteTemplate elimina (soft delete) una plantilla
func (uc *ManageEmailTemplatesUseCase) deleteTemplate(ctx context.Context, input *ManageEmailTemplatesInput, adminID int64) (*ManageEmailTemplatesOutput, error) {
	// Validar template_id
	if input.TemplateID == nil {
		return nil, errors.New("VALIDATION_FAILED", "template_id is required for delete", 400, nil)
	}

	// Buscar template
	var template EmailTemplate
	result := uc.db.WithContext(ctx).Table("email_templates").Where("id = ? AND deleted_at IS NULL", *input.TemplateID).First(&template)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("TEMPLATE_NOT_FOUND", "email template not found", 404, nil)
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Soft delete
	now := time.Now()
	result = uc.db.WithContext(ctx).Table("email_templates").Where("id = ?", template.ID).Update("deleted_at", now)
	if result.Error != nil {
		uc.log.Error("Error deleting email template", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	// Log auditoría
	uc.log.Error("Admin deleted email template",
		logger.Int64("admin_id", adminID),
		logger.Int64("template_id", template.ID),
		logger.String("name", template.Name),
		logger.String("action", "admin_delete_email_template"),
		logger.String("severity", "warning"))

	return &ManageEmailTemplatesOutput{
		Operation: "delete",
		Message:   "Email template deleted successfully",
	}, nil
}

// getTemplate obtiene una plantilla por ID
func (uc *ManageEmailTemplatesUseCase) getTemplate(ctx context.Context, input *ManageEmailTemplatesInput, adminID int64) (*ManageEmailTemplatesOutput, error) {
	// Validar template_id
	if input.TemplateID == nil {
		return nil, errors.New("VALIDATION_FAILED", "template_id is required", 400, nil)
	}

	// Buscar template
	var template EmailTemplate
	result := uc.db.WithContext(ctx).Table("email_templates").Where("id = ? AND deleted_at IS NULL", *input.TemplateID).First(&template)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("TEMPLATE_NOT_FOUND", "email template not found", 404, nil)
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	return &ManageEmailTemplatesOutput{
		Operation: "get",
		Template:  &template,
		Message:   "Email template retrieved successfully",
	}, nil
}

// listTemplates lista todas las plantillas
func (uc *ManageEmailTemplatesUseCase) listTemplates(ctx context.Context, input *ManageEmailTemplatesInput, adminID int64) (*ManageEmailTemplatesOutput, error) {
	query := uc.db.WithContext(ctx).Table("email_templates").Where("deleted_at IS NULL")

	// Filtrar por category si se proporciona
	if input.Category != "" {
		query = query.Where("category = ?", input.Category)
	}

	// Filtrar por is_active si se proporciona
	if input.IsActive != nil {
		query = query.Where("is_active = ?", *input.IsActive)
	}

	// Ordenar por created_at desc
	query = query.Order("created_at DESC")

	var templates []*EmailTemplate
	result := query.Find(&templates)
	if result.Error != nil {
		uc.log.Error("Error listing email templates", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	return &ManageEmailTemplatesOutput{
		Operation: "list",
		Templates: templates,
		Message:   "Email templates retrieved successfully",
	}, nil
}

// extractVariables extrae variables del formato {{variable}} del body
func (uc *ManageEmailTemplatesUseCase) extractVariables(body string) []string {
	variables := make([]string, 0)
	// Simple regex: buscar {{variable}}
	// TODO: Implementar con regexp para producción
	// Por ahora, retornar variables comunes
	commonVars := []string{"user_name", "user_email", "raffle_title", "payment_amount", "verification_link"}
	for _, v := range commonVars {
		if contains(body, "{{"+v+"}}") {
			variables = append(variables, v)
		}
	}
	return variables
}

// contains helper para buscar string en slice
func contains(slice interface{}, item string) bool {
	switch v := slice.(type) {
	case []string:
		for _, s := range v {
			if s == item {
				return true
			}
		}
	case string:
		return strings.Contains(v, item)
	}
	return false
}
