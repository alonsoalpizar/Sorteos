package db

import (
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresCompanySettingsRepository implementación de CompanySettingsRepository con PostgreSQL
type PostgresCompanySettingsRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCompanySettingsRepository crea una nueva instancia
func NewCompanySettingsRepository(db *gorm.DB, log *logger.Logger) *PostgresCompanySettingsRepository {
	return &PostgresCompanySettingsRepository{
		db:  db,
		log: log,
	}
}

// Get obtiene la configuración de la empresa (singleton)
func (r *PostgresCompanySettingsRepository) Get() (*domain.CompanySettings, error) {
	var settings domain.CompanySettings

	// Buscar el primer registro (debería ser el único)
	if err := r.db.First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting company settings", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &settings, nil
}

// Update actualiza la configuración de la empresa
func (r *PostgresCompanySettingsRepository) Update(settings *domain.CompanySettings) error {
	// Validar antes de actualizar
	if err := settings.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Actualizar usando Save (mantiene el ID)
	if err := r.db.Save(settings).Error; err != nil {
		r.log.Error("Error updating company settings", logger.Int64("id", settings.ID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// GetOrCreate obtiene la configuración existente o crea una nueva con valores por defecto
func (r *PostgresCompanySettingsRepository) GetOrCreate() (*domain.CompanySettings, error) {
	settings, err := r.Get()
	if err == nil {
		return settings, nil
	}

	// Si no existe, crear con valores por defecto
	if err == errors.ErrNotFound {
		defaultSettings := &domain.CompanySettings{
			CompanyName:  "Sorteos.club",
			Country:      "CR",
			Website:      "https://sorteos.club",
			SupportEmail: "soporte@sorteos.club",
		}

		if err := r.db.Create(defaultSettings).Error; err != nil {
			r.log.Error("Error creating default company settings", logger.Error(err))
			return nil, errors.Wrap(errors.ErrDatabaseError, err)
		}

		return defaultSettings, nil
	}

	return nil, err
}

// strPtr helper para crear punteros a string
func strPtr(s string) *string {
	return &s
}
