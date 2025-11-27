package db

import (
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresPaymentProcessorRepository implementación de PaymentProcessorRepository con PostgreSQL
type PostgresPaymentProcessorRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewPaymentProcessorRepository crea una nueva instancia
func NewPaymentProcessorRepository(db *gorm.DB, log *logger.Logger) *PostgresPaymentProcessorRepository {
	return &PostgresPaymentProcessorRepository{
		db:  db,
		log: log,
	}
}

// List obtiene todos los procesadores de pago
func (r *PostgresPaymentProcessorRepository) List() ([]*domain.PaymentProcessor, error) {
	var processors []*domain.PaymentProcessor

	if err := r.db.Order("created_at DESC").Find(&processors).Error; err != nil {
		r.log.Error("Error listing payment processors", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return processors, nil
}

// GetByID obtiene un procesador por ID
func (r *PostgresPaymentProcessorRepository) GetByID(id int64) (*domain.PaymentProcessor, error) {
	var processor domain.PaymentProcessor

	if err := r.db.First(&processor, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting payment processor by ID", logger.Int64("id", id), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &processor, nil
}

// GetByProvider obtiene un procesador por tipo de proveedor
func (r *PostgresPaymentProcessorRepository) GetByProvider(provider domain.ProcessorProvider) (*domain.PaymentProcessor, error) {
	var processor domain.PaymentProcessor

	if err := r.db.Where("provider = ?", provider).First(&processor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting payment processor by provider",
			logger.String("provider", string(provider)),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &processor, nil
}

// FindByProvider obtiene un procesador por proveedor y modo sandbox
func (r *PostgresPaymentProcessorRepository) FindByProvider(provider string, isSandbox bool) (*domain.PaymentProcessor, error) {
	var processor domain.PaymentProcessor

	if err := r.db.Where("provider = ? AND is_sandbox = ?", provider, isSandbox).First(&processor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error finding payment processor",
			logger.String("provider", provider),
			logger.Bool("is_sandbox", isSandbox),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &processor, nil
}

// GetActive obtiene el procesador activo
func (r *PostgresPaymentProcessorRepository) GetActive() (*domain.PaymentProcessor, error) {
	var processor domain.PaymentProcessor

	if err := r.db.Where("is_active = ?", true).First(&processor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting active payment processor", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &processor, nil
}

// Update actualiza un procesador de pago
func (r *PostgresPaymentProcessorRepository) Update(processor *domain.PaymentProcessor) error {
	// Validar antes de actualizar
	if err := processor.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	if err := r.db.Save(processor).Error; err != nil {
		r.log.Error("Error updating payment processor",
			logger.Int64("id", processor.ID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// ToggleActive activa/desactiva un procesador
func (r *PostgresPaymentProcessorRepository) ToggleActive(id int64, active bool) error {
	// Si se está activando un procesador, primero desactivar todos los demás
	if active {
		if err := r.db.Model(&domain.PaymentProcessor{}).
			Where("id != ?", id).
			Update("is_active", false).Error; err != nil {
			r.log.Error("Error deactivating other processors", logger.Error(err))
			return errors.Wrap(errors.ErrDatabaseError, err)
		}
	}

	// Activar/desactivar el procesador especificado
	if err := r.db.Model(&domain.PaymentProcessor{}).
		Where("id = ?", id).
		Update("is_active", active).Error; err != nil {
		r.log.Error("Error toggling payment processor status",
			logger.Int64("id", id),
			logger.Bool("active", active),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}
