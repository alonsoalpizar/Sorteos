package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresCreditPurchaseRepository implementación de CreditPurchaseRepository con PostgreSQL
type PostgresCreditPurchaseRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCreditPurchaseRepository crea una nueva instancia
func NewCreditPurchaseRepository(db *gorm.DB, log *logger.Logger) *PostgresCreditPurchaseRepository {
	return &PostgresCreditPurchaseRepository{
		db:  db,
		log: log,
	}
}

// Create crea una nueva compra
func (r *PostgresCreditPurchaseRepository) Create(purchase *domain.CreditPurchase) error {
	// Generar UUID si no existe
	if purchase.UUID == "" {
		purchase.UUID = uuid.New().String()
	}

	// Validar antes de crear
	if err := purchase.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Verificar que la idempotency key sea única
	var existing domain.CreditPurchase
	err := r.db.Where("idempotency_key = ?", purchase.IdempotencyKey).First(&existing).Error
	if err == nil {
		// Ya existe una compra con esta idempotency key
		r.log.Warn("Compra duplicada detectada (idempotency)",
			logger.String("idempotency_key", purchase.IdempotencyKey),
			logger.Int64("existing_purchase_id", existing.ID))
		return errors.Wrap(errors.ErrConflict, fmt.Errorf("compra duplicada"))
	}
	if err != gorm.ErrRecordNotFound {
		r.log.Error("Error verificando idempotency key",
			logger.String("idempotency_key", purchase.IdempotencyKey),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Verificar que el ERN sea único
	err = r.db.Where("ern = ?", purchase.ERN).First(&existing).Error
	if err == nil {
		r.log.Warn("ERN duplicado detectado",
			logger.String("ern", purchase.ERN),
			logger.Int64("existing_purchase_id", existing.ID))
		return errors.Wrap(errors.ErrConflict, fmt.Errorf("ERN duplicado"))
	}
	if err != gorm.ErrRecordNotFound {
		r.log.Error("Error verificando ERN",
			logger.String("ern", purchase.ERN),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Crear compra
	if err := r.db.Create(purchase).Error; err != nil {
		r.log.Error("Error creando compra de créditos",
			logger.Int64("user_id", purchase.UserID),
			logger.String("ern", purchase.ERN),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	r.log.Info("Compra de créditos creada exitosamente",
		logger.Int64("purchase_id", purchase.ID),
		logger.Int64("user_id", purchase.UserID),
		logger.String("ern", purchase.ERN),
		logger.String("desired_credit", purchase.DesiredCredit.String()),
		logger.String("charge_amount", purchase.ChargeAmount.String()))

	return nil
}

// FindByID busca una compra por ID
func (r *PostgresCreditPurchaseRepository) FindByID(id int64) (*domain.CreditPurchase, error) {
	var purchase domain.CreditPurchase

	if err := r.db.First(&purchase, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando compra por ID",
			logger.Int64("id", id),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &purchase, nil
}

// FindByUUID busca una compra por UUID
func (r *PostgresCreditPurchaseRepository) FindByUUID(uuid string) (*domain.CreditPurchase, error) {
	var purchase domain.CreditPurchase

	if err := r.db.Where("uuid = ?", uuid).First(&purchase).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando compra por UUID",
			logger.String("uuid", uuid),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &purchase, nil
}

// FindByERN busca una compra por ERN (External Reference Number)
func (r *PostgresCreditPurchaseRepository) FindByERN(ern string) (*domain.CreditPurchase, error) {
	var purchase domain.CreditPurchase

	if err := r.db.Where("ern = ?", ern).First(&purchase).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando compra por ERN",
			logger.String("ern", ern),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &purchase, nil
}

// FindByIdempotencyKey busca una compra por clave de idempotencia
func (r *PostgresCreditPurchaseRepository) FindByIdempotencyKey(key string) (*domain.CreditPurchase, error) {
	var purchase domain.CreditPurchase

	if err := r.db.Where("idempotency_key = ?", key).First(&purchase).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando compra por idempotency key",
			logger.String("idempotency_key", key),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &purchase, nil
}

// FindByPagaditoToken busca una compra por token de Pagadito
func (r *PostgresCreditPurchaseRepository) FindByPagaditoToken(token string) (*domain.CreditPurchase, error) {
	var purchase domain.CreditPurchase

	if err := r.db.Where("pagadito_token = ?", token).First(&purchase).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando compra por pagadito_token",
			logger.String("token", token),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &purchase, nil
}

// FindByUserID busca compras de un usuario (paginado)
func (r *PostgresCreditPurchaseRepository) FindByUserID(userID int64, limit, offset int) ([]*domain.CreditPurchase, int64, error) {
	var purchases []*domain.CreditPurchase
	var total int64

	// Contar total
	if err := r.db.Model(&domain.CreditPurchase{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		r.log.Error("Error contando compras del usuario",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener compras paginadas
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&purchases).Error; err != nil {
		r.log.Error("Error buscando compras del usuario",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return purchases, total, nil
}

// Update actualiza una compra
func (r *PostgresCreditPurchaseRepository) Update(purchase *domain.CreditPurchase) error {
	// Validar antes de actualizar
	if err := purchase.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Actualizar
	if err := r.db.Save(purchase).Error; err != nil {
		r.log.Error("Error actualizando compra",
			logger.Int64("purchase_id", purchase.ID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	r.log.Info("Compra actualizada exitosamente",
		logger.Int64("purchase_id", purchase.ID),
		logger.String("status", string(purchase.Status)))

	return nil
}

// MarkExpired marca como expiradas las compras que superaron el TTL
func (r *PostgresCreditPurchaseRepository) MarkExpired() (int64, error) {
	// Actualizar compras que están en pending/processing y ya expiraron
	result := r.db.
		Model(&domain.CreditPurchase{}).
		Where("status IN (?, ?)", domain.CreditPurchaseStatusPending, domain.CreditPurchaseStatusProcessing).
		Where("expires_at < ?", time.Now()).
		Updates(map[string]interface{}{
			"status":     domain.CreditPurchaseStatusExpired,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		r.log.Error("Error marcando compras como expiradas",
			logger.Error(result.Error))
		return 0, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	if result.RowsAffected > 0 {
		r.log.Info("Compras marcadas como expiradas",
			logger.Int64("count", result.RowsAffected))
	}

	return result.RowsAffected, nil
}
