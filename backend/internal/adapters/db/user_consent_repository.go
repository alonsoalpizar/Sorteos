package db

import (
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// UserConsentRepositoryImpl implementa domain.UserConsentRepository
type UserConsentRepositoryImpl struct {
	db *gorm.DB
}

// NewUserConsentRepository crea una nueva instancia del repositorio
func NewUserConsentRepository(db *gorm.DB) domain.UserConsentRepository {
	return &UserConsentRepositoryImpl{db: db}
}

// Create crea un nuevo consentimiento
func (r *UserConsentRepositoryImpl) Create(consent *domain.UserConsent) error {
	if err := r.db.Create(consent).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// FindByUserAndType busca un consentimiento por usuario y tipo
func (r *UserConsentRepositoryImpl) FindByUserAndType(userID int64, consentType domain.ConsentType) (*domain.UserConsent, error) {
	var consent domain.UserConsent
	if err := r.db.Where("user_id = ? AND consent_type = ?", userID, consentType).First(&consent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return &consent, nil
}

// FindByUser busca todos los consentimientos de un usuario
func (r *UserConsentRepositoryImpl) FindByUser(userID int64) ([]*domain.UserConsent, error) {
	var consents []*domain.UserConsent
	if err := r.db.Where("user_id = ?", userID).Find(&consents).Error; err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return consents, nil
}

// Update actualiza un consentimiento existente
func (r *UserConsentRepositoryImpl) Update(consent *domain.UserConsent) error {
	if err := r.db.Save(consent).Error; err != nil {
		return errors.Wrap(errors.ErrDatabaseError, err)
	}
	return nil
}

// HasGrantedConsent verifica si el usuario ha otorgado un consentimiento específico
func (r *UserConsentRepositoryImpl) HasGrantedConsent(userID int64, consentType domain.ConsentType) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.UserConsent{}).
		Where("user_id = ? AND consent_type = ? AND granted = true AND revoked_at IS NULL", userID, consentType).
		Count(&count).Error; err != nil {
		return false, errors.Wrap(errors.ErrDatabaseError, err)
	}
	return count > 0, nil
}

// RevokeConsent revoca un consentimiento
func (r *UserConsentRepositoryImpl) RevokeConsent(userID int64, consentType domain.ConsentType) error {
	consent, err := r.FindByUserAndType(userID, consentType)
	if err != nil {
		return err
	}

	consent.Revoke()
	return r.Update(consent)
}

// GrantConsent otorga un consentimiento
func (r *UserConsentRepositoryImpl) GrantConsent(userID int64, consentType domain.ConsentType, version, ipAddress, userAgent string) error {
	// Buscar consentimiento existente
	consent, err := r.FindByUserAndType(userID, consentType)
	if err != nil && err != errors.ErrNotFound {
		return err
	}

	if consent != nil {
		// Actualizar consentimiento existente
		consent.ConsentVersion = version
		consent.Grant(ipAddress, userAgent)
		return r.Update(consent)
	}

	// Crear nuevo consentimiento
	newConsent := &domain.UserConsent{
		UserID:         userID,
		ConsentType:    consentType,
		ConsentVersion: version,
	}
	newConsent.Grant(ipAddress, userAgent)

	return r.Create(newConsent)
}
