package domain

import (
	"fmt"
	"regexp"
	"time"
)

// UserRole representa el rol del usuario
type UserRole string

const (
	UserRoleUser       UserRole = "user"
	UserRoleAdmin      UserRole = "admin"
	UserRoleSuperAdmin UserRole = "super_admin"
)

// KYCLevel representa el nivel de verificación KYC
type KYCLevel string

const (
	KYCLevelNone            KYCLevel = "none"
	KYCLevelEmailVerified   KYCLevel = "email_verified"
	KYCLevelPhoneVerified   KYCLevel = "phone_verified"
	KYCLevelCedulaVerified  KYCLevel = "cedula_verified"
	KYCLevelFullKYC         KYCLevel = "full_kyc"
)

// UserStatus representa el estado del usuario
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
	UserStatusDeleted   UserStatus = "deleted"
)

// User representa un usuario del sistema
type User struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	UUID string `json:"uuid" gorm:"type:uuid;unique;not null"`

	// Credenciales
	Email              string     `json:"email" gorm:"uniqueIndex;not null"`
	EmailVerified      bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt    *time.Time `json:"email_verified_at,omitempty"`
	Phone              *string    `json:"phone,omitempty" gorm:"uniqueIndex"`
	PhoneVerified      bool       `json:"phone_verified" gorm:"default:false"`
	PhoneVerifiedAt    *time.Time `json:"phone_verified_at,omitempty"`
	PasswordHash       string     `json:"-" gorm:"not null"` // Never serialize password

	// Información personal
	FirstName       *string    `json:"first_name,omitempty"`
	LastName        *string    `json:"last_name,omitempty"`
	Cedula          *string    `json:"cedula,omitempty" gorm:"uniqueIndex"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"` // Fecha de nacimiento
	ProfilePhotoURL *string    `json:"profile_photo_url,omitempty"`

	// Dirección
	AddressLine1 *string `json:"address_line1,omitempty"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Country      string  `json:"country" gorm:"type:char(2);default:'CR'"`

	// Información bancaria (encriptado en app layer)
	IBAN *string `json:"iban,omitempty"`

	// Roles y verificación
	Role     UserRole   `json:"role" gorm:"type:user_role;default:'user';not null"`
	KYCLevel KYCLevel   `json:"kyc_level" gorm:"type:kyc_level;default:'none';not null"`
	Status   UserStatus `json:"status" gorm:"type:user_status;default:'active';not null"`

	// Límites
	MaxActiveRaffles    int     `json:"max_active_raffles" gorm:"default:10"`
	PurchaseLimitDaily  float64 `json:"purchase_limit_daily" gorm:"type:decimal(12,2);default:50000.00"`

	// Tokens (no serializar en JSON)
	RefreshToken          *string    `json:"-"`
	RefreshTokenExpiresAt *time.Time `json:"-"`

	// Códigos de verificación (no serializar en JSON)
	EmailVerificationCode      *string    `json:"-"`
	EmailVerificationExpiresAt *time.Time `json:"-"`
	PhoneVerificationCode      *string    `json:"-"`
	PhoneVerificationExpiresAt *time.Time `json:"-"`
	PasswordResetToken         *string    `json:"-"`
	PasswordResetExpiresAt     *time.Time `json:"-"`

	// Admin fields (for Almighty module)
	SuspensionReason *string    `json:"suspension_reason,omitempty"`
	SuspendedBy      *int64     `json:"suspended_by,omitempty"` // Admin user ID que suspendió
	SuspendedAt      *time.Time `json:"suspended_at,omitempty"`
	LastKYCReview    *time.Time `json:"last_kyc_review,omitempty"`
	KYCReviewer      *int64     `json:"kyc_reviewer,omitempty"` // Admin user ID que revisó KYC

	// Auditoría
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP *string    `json:"-"` // No exponer IP públicamente

	// Soft delete
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// TableName especifica el nombre de la tabla
func (User) TableName() string {
	return "users"
}

// Validaciones

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{7,14}$`) // E.164 format
)

// ValidateEmail valida el formato del email
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if len(email) > 255 {
		return fmt.Errorf("email is too long (max 255 characters)")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePhone valida el formato del teléfono
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // Phone is optional
	}
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone format (use E.164: +50612345678)")
	}
	return nil
}

// ValidatePassword valida la fortaleza de la contraseña
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password is too long (max 128 characters)")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' ||
			 char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character (!@#$%%^&*)")
	}

	return nil
}

// IsActive verifica si el usuario está activo
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && u.DeletedAt == nil
}

// IsAdmin verifica si el usuario es administrador
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleSuperAdmin
}

// IsSuperAdmin verifica si el usuario es super administrador
func (u *User) IsSuperAdmin() bool {
	return u.Role == UserRoleSuperAdmin
}

// HasMinimumKYC verifica si el usuario tiene el nivel mínimo de KYC
func (u *User) HasMinimumKYC(minLevel KYCLevel) bool {
	levels := map[KYCLevel]int{
		KYCLevelNone:           0,
		KYCLevelEmailVerified:  1,
		KYCLevelPhoneVerified:  2,
		KYCLevelCedulaVerified: 3,
		KYCLevelFullKYC:        4,
	}
	return levels[u.KYCLevel] >= levels[minLevel]
}

// CanCreateRaffles verifica si el usuario puede crear sorteos
func (u *User) CanCreateRaffles() bool {
	// Debe tener al menos email verificado
	return u.IsActive() && u.HasMinimumKYC(KYCLevelEmailVerified)
}

// CanPurchase verifica si el usuario puede comprar boletos
func (u *User) CanPurchase() bool {
	return u.IsActive()
}

// GetFullName retorna el nombre completo del usuario
func (u *User) GetFullName() string {
	if u.FirstName != nil && u.LastName != nil {
		return fmt.Sprintf("%s %s", *u.FirstName, *u.LastName)
	}
	if u.FirstName != nil {
		return *u.FirstName
	}
	return u.Email
}

// CanWithdraw verifica si el usuario puede retirar ganancias
// Requisitos: full_kyc + IBAN configurado
func (u *User) CanWithdraw() bool {
	return u.IsActive() &&
		u.KYCLevel == KYCLevelFullKYC &&
		u.IBAN != nil &&
		*u.IBAN != ""
}

// ValidateIBAN valida el formato IBAN costarricense
func ValidateIBAN(iban string) error {
	if iban == "" {
		return fmt.Errorf("IBAN is required")
	}

	// IBAN costarricense: CR + 22 dígitos (total 24 caracteres)
	if len(iban) != 24 {
		return fmt.Errorf("IBAN must be exactly 24 characters (CR + 22 digits)")
	}

	if iban[0:2] != "CR" {
		return fmt.Errorf("IBAN must start with 'CR' for Costa Rica")
	}

	// Verificar que el resto sean dígitos
	for i := 2; i < len(iban); i++ {
		if iban[i] < '0' || iban[i] > '9' {
			return fmt.Errorf("IBAN digits (positions 3-24) must be numeric")
		}
	}

	return nil
}

// ValidateDateOfBirth valida que la fecha de nacimiento sea razonable
func ValidateDateOfBirth(dob time.Time) error {
	now := time.Now()

	// No puede ser fecha futura
	if dob.After(now) {
		return fmt.Errorf("date of birth cannot be in the future")
	}

	// Edad mínima 18 años
	minAge := now.AddDate(-18, 0, 0)
	if dob.After(minAge) {
		return fmt.Errorf("user must be at least 18 years old")
	}

	// Edad máxima 120 años (validación de cordura)
	maxAge := now.AddDate(-120, 0, 0)
	if dob.Before(maxAge) {
		return fmt.Errorf("date of birth is not reasonable (max 120 years old)")
	}

	return nil
}

// UserRepository define el contrato para el repositorio de usuarios
type UserRepository interface {
	// Create crea un nuevo usuario
	Create(user *User) error

	// FindByID busca un usuario por ID
	FindByID(id int64) (*User, error)

	// FindByUUID busca un usuario por UUID
	FindByUUID(uuid string) (*User, error)

	// FindByEmail busca un usuario por email
	FindByEmail(email string) (*User, error)

	// FindByPhone busca un usuario por teléfono
	FindByPhone(phone string) (*User, error)

	// FindByCedula busca un usuario por cédula
	FindByCedula(cedula string) (*User, error)

	// Update actualiza un usuario existente
	Update(user *User) error

	// UpdateRefreshToken actualiza el refresh token
	UpdateRefreshToken(userID int64, token string, expiresAt time.Time) error

	// UpdateLastLogin actualiza la fecha y IP del último login
	UpdateLastLogin(userID int64, ip string) error

	// VerifyEmail marca el email como verificado
	VerifyEmail(userID int64) error

	// VerifyPhone marca el teléfono como verificado
	VerifyPhone(userID int64) error

	// UpdateKYCLevel actualiza el nivel de KYC
	UpdateKYCLevel(userID int64, level KYCLevel) error

	// SoftDelete marca el usuario como eliminado (soft delete)
	SoftDelete(userID int64) error

	// List retorna una lista paginada de usuarios
	List(offset, limit int, filters map[string]interface{}) ([]*User, int64, error)
}
