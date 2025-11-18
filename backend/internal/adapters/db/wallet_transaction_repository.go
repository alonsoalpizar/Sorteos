package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresWalletTransactionRepository implementación de WalletTransactionRepository con PostgreSQL
type PostgresWalletTransactionRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewWalletTransactionRepository crea una nueva instancia
func NewWalletTransactionRepository(db *gorm.DB, log *logger.Logger) *PostgresWalletTransactionRepository {
	return &PostgresWalletTransactionRepository{
		db:  db,
		log: log,
	}
}

// Create crea una nueva transacción
func (r *PostgresWalletTransactionRepository) Create(tx *domain.WalletTransaction) error {
	// Generar UUID si no existe
	if tx.UUID == "" {
		tx.UUID = uuid.New().String()
	}

	// Validar antes de crear
	if err := tx.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Verificar que la idempotency key sea única
	var existing domain.WalletTransaction
	err := r.db.Where("idempotency_key = ?", tx.IdempotencyKey).First(&existing).Error
	if err == nil {
		// Ya existe una transacción con esta idempotency key
		r.log.Warn("Transacción duplicada detectada (idempotency)",
			logger.String("idempotency_key", tx.IdempotencyKey),
			logger.Int64("existing_tx_id", existing.ID))
		return errors.Wrap(errors.ErrConflict, fmt.Errorf("transacción duplicada"))
	}
	if err != gorm.ErrRecordNotFound {
		r.log.Error("Error verificando idempotency key",
			logger.String("idempotency_key", tx.IdempotencyKey),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Crear transacción
	if err := r.db.Create(tx).Error; err != nil {
		r.log.Error("Error creando transacción",
			logger.Int64("wallet_id", tx.WalletID),
			logger.String("type", string(tx.Type)),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	r.log.Info("Transacción creada exitosamente",
		logger.Int64("tx_id", tx.ID),
		logger.Int64("wallet_id", tx.WalletID),
		logger.String("type", string(tx.Type)),
		logger.String("amount", tx.Amount.String()))

	return nil
}

// FindByID busca una transacción por ID
func (r *PostgresWalletTransactionRepository) FindByID(id int64) (*domain.WalletTransaction, error) {
	var tx domain.WalletTransaction

	if err := r.db.First(&tx, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando transacción por ID",
			logger.Int64("id", id),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &tx, nil
}

// FindByUUID busca una transacción por UUID
func (r *PostgresWalletTransactionRepository) FindByUUID(uuid string) (*domain.WalletTransaction, error) {
	var tx domain.WalletTransaction

	if err := r.db.Where("uuid = ?", uuid).First(&tx).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando transacción por UUID",
			logger.String("uuid", uuid),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &tx, nil
}

// FindByIdempotencyKey busca una transacción por clave de idempotencia
func (r *PostgresWalletTransactionRepository) FindByIdempotencyKey(key string) (*domain.WalletTransaction, error) {
	var tx domain.WalletTransaction

	if err := r.db.Where("idempotency_key = ?", key).First(&tx).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando transacción por idempotency_key",
			logger.String("idempotency_key", key),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &tx, nil
}

// FindByWalletID busca transacciones de una billetera (paginado)
func (r *PostgresWalletTransactionRepository) FindByWalletID(walletID int64, limit, offset int) ([]*domain.WalletTransaction, int64, error) {
	var transactions []*domain.WalletTransaction
	var total int64

	// Contar total
	if err := r.db.Model(&domain.WalletTransaction{}).
		Where("wallet_id = ?", walletID).
		Count(&total).Error; err != nil {
		r.log.Error("Error contando transacciones",
			logger.Int64("wallet_id", walletID),
			logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener transacciones paginadas
	if err := r.db.
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		r.log.Error("Error buscando transacciones por wallet_id",
			logger.Int64("wallet_id", walletID),
			logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return transactions, total, nil
}

// FindByUserID busca transacciones de un usuario (paginado)
func (r *PostgresWalletTransactionRepository) FindByUserID(userID int64, limit, offset int) ([]*domain.WalletTransaction, int64, error) {
	var transactions []*domain.WalletTransaction
	var total int64

	// Contar total
	if err := r.db.Model(&domain.WalletTransaction{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		r.log.Error("Error contando transacciones",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener transacciones paginadas
	if err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		r.log.Error("Error buscando transacciones por user_id",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return transactions, total, nil
}

// FindByReference busca transacciones por referencia externa
func (r *PostgresWalletTransactionRepository) FindByReference(referenceType string, referenceID int64) ([]*domain.WalletTransaction, error) {
	var transactions []*domain.WalletTransaction

	if err := r.db.
		Where("reference_type = ? AND reference_id = ?", referenceType, referenceID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		r.log.Error("Error buscando transacciones por referencia",
			logger.String("reference_type", referenceType),
			logger.Int64("reference_id", referenceID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return transactions, nil
}

// Update actualiza una transacción existente
func (r *PostgresWalletTransactionRepository) Update(tx *domain.WalletTransaction) error {
	// Validar antes de actualizar
	if err := tx.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	if err := r.db.Save(tx).Error; err != nil {
		r.log.Error("Error actualizando transacción",
			logger.Int64("tx_id", tx.ID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	r.log.Info("Transacción actualizada",
		logger.Int64("tx_id", tx.ID),
		logger.String("status", string(tx.Status)))

	return nil
}

// WithTransaction ejecuta una función dentro de una transacción
func (r *PostgresWalletTransactionRepository) WithTransaction(fn func(repo domain.WalletTransactionRepository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, tx.Error)
	}

	// Crear repositorio con la transacción
	txRepo := &PostgresWalletTransactionRepository{
		db:  tx,
		log: r.log,
	}

	// Ejecutar función
	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		r.log.Error("Error en commit de transacción", logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}
