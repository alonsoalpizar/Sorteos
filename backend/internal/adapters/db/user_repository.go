package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// UserRepositoryImpl implementa domain.UserRepository usando GORM
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create crea un nuevo usuario
func (r *UserRepositoryImpl) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// FindByID busca un usuario por ID
func (r *UserRepositoryImpl) FindByID(id int64) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &user, nil
}

// FindByUUID busca un usuario por UUID
func (r *UserRepositoryImpl) FindByUUID(uuid string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &user, nil
}

// FindByEmail busca un usuario por email
func (r *UserRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &user, nil
}

// FindByPhone busca un usuario por teléfono
func (r *UserRepositoryImpl) FindByPhone(phone string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("phone = ? AND deleted_at IS NULL", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &user, nil
}

// FindByCedula busca un usuario por cédula
func (r *UserRepositoryImpl) FindByCedula(cedula string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("cedula = ? AND deleted_at IS NULL", cedula).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &user, nil
}

// Update actualiza un usuario existente
func (r *UserRepositoryImpl) Update(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// UpdateRefreshToken actualiza el refresh token
func (r *UserRepositoryImpl) UpdateRefreshToken(userID int64, token string, expiresAt time.Time) error {
	result := r.db.Model(&domain.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"refresh_token":           token,
			"refresh_token_expires_at": expiresAt,
		})

	if result.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// UpdateLastLogin actualiza la fecha y IP del último login
func (r *UserRepositoryImpl) UpdateLastLogin(userID int64, ip string) error {
	now := time.Now()
	result := r.db.Model(&domain.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
		})

	if result.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// VerifyEmail marca el email como verificado
func (r *UserRepositoryImpl) VerifyEmail(userID int64) error {
	now := time.Now()
	result := r.db.Model(&domain.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"email_verified":    true,
			"email_verified_at": now,
			"kyc_level":         domain.KYCLevelEmailVerified,
		})

	if result.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// VerifyPhone marca el teléfono como verificado
func (r *UserRepositoryImpl) VerifyPhone(userID int64) error {
	now := time.Now()

	// Primero obtener el usuario para ver su KYC actual
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}

	// Si ya tiene email verificado, subir a phone_verified
	newKYCLevel := domain.KYCLevelPhoneVerified
	if user.KYCLevel < domain.KYCLevelEmailVerified {
		newKYCLevel = user.KYCLevel // Mantener el nivel actual si no tiene email
	}

	result := r.db.Model(&domain.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"phone_verified":    true,
			"phone_verified_at": now,
			"kyc_level":         newKYCLevel,
		})

	if result.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// UpdateKYCLevel actualiza el nivel de KYC
func (r *UserRepositoryImpl) UpdateKYCLevel(userID int64, level domain.KYCLevel) error {
	result := r.db.Model(&domain.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Update("kyc_level", level)

	if result.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// SoftDelete marca el usuario como eliminado (soft delete)
func (r *UserRepositoryImpl) SoftDelete(userID int64) error {
	now := time.Now()
	result := r.db.Model(&domain.User{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"status":     domain.UserStatusDeleted,
		})

	if result.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// List retorna una lista paginada de usuarios
func (r *UserRepositoryImpl) List(offset, limit int, filters map[string]interface{}) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	query := r.db.Model(&domain.User{}).Where("deleted_at IS NULL")

	// Aplicar filtros
	if role, ok := filters["role"].(domain.UserRole); ok {
		query = query.Where("role = ?", role)
	}
	if status, ok := filters["status"].(domain.UserStatus); ok {
		query = query.Where("status = ?", status)
	}
	if kycLevel, ok := filters["kyc_level"].(domain.KYCLevel); ok {
		query = query.Where("kyc_level = ?", kycLevel)
	}
	if emailVerified, ok := filters["email_verified"].(bool); ok {
		query = query.Where("email_verified = ?", emailVerified)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where(
			"email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener página
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return users, total, nil
}
