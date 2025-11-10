package domain

import (
	"encoding/json"
	"time"
)

// AuditAction representa el tipo de acción auditada
type AuditAction string

const (
	// Auth
	AuditActionUserRegistered AuditAction = "user_registered"
	AuditActionUserLoggedIn   AuditAction = "user_logged_in"
	AuditActionUserLoggedOut  AuditAction = "user_logged_out"
	AuditActionEmailVerified  AuditAction = "email_verified"
	AuditActionPhoneVerified  AuditAction = "phone_verified"

	// User Management
	AuditActionUserUpdated      AuditAction = "user_updated"
	AuditActionUserSuspended    AuditAction = "user_suspended"
	AuditActionUserBanned       AuditAction = "user_banned"
	AuditActionUserDeleted      AuditAction = "user_deleted"
	AuditActionKYCLevelChanged  AuditAction = "kyc_level_changed"

	// Raffles
	AuditActionRaffleCreated    AuditAction = "raffle_created"
	AuditActionRafflePublished  AuditAction = "raffle_published"
	AuditActionRaffleSuspended  AuditAction = "raffle_suspended"
	AuditActionRaffleCompleted  AuditAction = "raffle_completed"
	AuditActionRaffleDeleted    AuditAction = "raffle_deleted"

	// Reservations
	AuditActionNumbersReserved      AuditAction = "numbers_reserved"
	AuditActionReservationExpired   AuditAction = "reservation_expired"
	AuditActionReservationCancelled AuditAction = "reservation_cancelled"

	// Payments
	AuditActionPaymentCreated   AuditAction = "payment_created"
	AuditActionPaymentConfirmed AuditAction = "payment_confirmed"
	AuditActionPaymentFailed    AuditAction = "payment_failed"
	AuditActionPaymentRefunded  AuditAction = "payment_refunded"

	// Settlements
	AuditActionSettlementCreated  AuditAction = "settlement_created"
	AuditActionSettlementApproved AuditAction = "settlement_approved"
	AuditActionSettlementPaid     AuditAction = "settlement_paid"
	AuditActionSettlementRejected AuditAction = "settlement_rejected"

	// Admin Actions
	AuditActionAdminActionPerformed   AuditAction = "admin_action_performed"
	AuditActionSystemParameterChanged AuditAction = "system_parameter_changed"
	AuditActionReportGenerated        AuditAction = "report_generated"
)

// AuditSeverity representa el nivel de criticidad
type AuditSeverity string

const (
	AuditSeverityInfo     AuditSeverity = "info"
	AuditSeverityWarning  AuditSeverity = "warning"
	AuditSeverityError    AuditSeverity = "error"
	AuditSeverityCritical AuditSeverity = "critical"
)

// AuditLog representa un registro de auditoría
type AuditLog struct {
	ID          int64         `json:"id" gorm:"primaryKey"`
	UserID      *int64        `json:"user_id,omitempty" gorm:"index"`
	AdminID     *int64        `json:"admin_id,omitempty" gorm:"index"`
	Action      AuditAction   `json:"action" gorm:"type:audit_action;not null;index"`
	Severity    AuditSeverity `json:"severity" gorm:"type:audit_severity;default:'info';not null;index"`
	Description *string       `json:"description,omitempty"`

	// Entidad afectada (polimórfico)
	EntityType *string `json:"entity_type,omitempty" gorm:"index:idx_audit_entity"` // e.g., "raffle", "user", "payment"
	EntityID   *int64  `json:"entity_id,omitempty" gorm:"index:idx_audit_entity"`

	// Contexto de la solicitud
	IPAddress      *string `json:"ip_address,omitempty" gorm:"index"`
	UserAgent      *string `json:"user_agent,omitempty"`
	Endpoint       *string `json:"endpoint,omitempty"`
	HTTPMethod     *string `json:"http_method,omitempty"`
	HTTPStatusCode *int    `json:"http_status_code,omitempty"`

	// Datos adicionales (JSON)
	Metadata json.RawMessage `json:"metadata,omitempty" gorm:"type:jsonb"`

	// Timestamp
	CreatedAt time.Time `json:"created_at" gorm:"index:idx_audit_created_at"`

	// Relaciones
	User  *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Admin *User `json:"admin,omitempty" gorm:"foreignKey:AdminID"`
}

// TableName especifica el nombre de la tabla
func (AuditLog) TableName() string {
	return "audit_logs"
}

// SetMetadata establece el metadata como JSON
func (al *AuditLog) SetMetadata(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	al.Metadata = jsonData
	return nil
}

// GetMetadata obtiene el metadata y lo deserializa
func (al *AuditLog) GetMetadata(dest interface{}) error {
	if al.Metadata == nil {
		return nil
	}
	return json.Unmarshal(al.Metadata, dest)
}

// AuditLogRepository define el contrato para el repositorio de audit logs
type AuditLogRepository interface {
	// Create crea un nuevo registro de auditoría
	Create(log *AuditLog) error

	// FindByID busca un log por ID
	FindByID(id int64) (*AuditLog, error)

	// FindByUser busca logs de un usuario específico
	FindByUser(userID int64, offset, limit int) ([]*AuditLog, int64, error)

	// FindByAdmin busca logs de acciones de admin
	FindByAdmin(adminID int64, offset, limit int) ([]*AuditLog, int64, error)

	// FindByEntity busca logs de una entidad específica
	FindByEntity(entityType string, entityID int64, offset, limit int) ([]*AuditLog, int64, error)

	// FindByAction busca logs por tipo de acción
	FindByAction(action AuditAction, offset, limit int) ([]*AuditLog, int64, error)

	// FindBySeverity busca logs por severidad
	FindBySeverity(severity AuditSeverity, offset, limit int) ([]*AuditLog, int64, error)

	// FindByDateRange busca logs en un rango de fechas
	FindByDateRange(start, end time.Time, offset, limit int) ([]*AuditLog, int64, error)

	// List retorna una lista paginada de logs con filtros
	List(offset, limit int, filters map[string]interface{}) ([]*AuditLog, int64, error)
}

// AuditLogBuilder facilita la creación de audit logs
type AuditLogBuilder struct {
	log *AuditLog
}

// NewAuditLog crea un nuevo builder de audit log
func NewAuditLog(action AuditAction) *AuditLogBuilder {
	return &AuditLogBuilder{
		log: &AuditLog{
			Action:    action,
			Severity:  AuditSeverityInfo,
			CreatedAt: time.Now(),
		},
	}
}

// WithUser establece el usuario que realizó la acción
func (b *AuditLogBuilder) WithUser(userID int64) *AuditLogBuilder {
	b.log.UserID = &userID
	return b
}

// WithAdmin establece el admin que realizó la acción
func (b *AuditLogBuilder) WithAdmin(adminID int64) *AuditLogBuilder {
	b.log.AdminID = &adminID
	return b
}

// WithSeverity establece la severidad
func (b *AuditLogBuilder) WithSeverity(severity AuditSeverity) *AuditLogBuilder {
	b.log.Severity = severity
	return b
}

// WithDescription establece la descripción
func (b *AuditLogBuilder) WithDescription(description string) *AuditLogBuilder {
	b.log.Description = &description
	return b
}

// WithEntity establece la entidad afectada
func (b *AuditLogBuilder) WithEntity(entityType string, entityID int64) *AuditLogBuilder {
	b.log.EntityType = &entityType
	b.log.EntityID = &entityID
	return b
}

// WithRequest establece el contexto de la solicitud HTTP
func (b *AuditLogBuilder) WithRequest(ip, userAgent, endpoint, method string, statusCode int) *AuditLogBuilder {
	b.log.IPAddress = &ip
	b.log.UserAgent = &userAgent
	b.log.Endpoint = &endpoint
	b.log.HTTPMethod = &method
	b.log.HTTPStatusCode = &statusCode
	return b
}

// WithMetadata establece metadata adicional
func (b *AuditLogBuilder) WithMetadata(data interface{}) *AuditLogBuilder {
	_ = b.log.SetMetadata(data)
	return b
}

// Build retorna el audit log construido
func (b *AuditLogBuilder) Build() *AuditLog {
	return b.log
}
