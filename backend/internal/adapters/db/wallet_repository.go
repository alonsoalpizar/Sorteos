package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresWalletRepository implementación de WalletRepository con PostgreSQL
type PostgresWalletRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewWalletRepository crea una nueva instancia
func NewWalletRepository(db *gorm.DB, log *logger.Logger) *PostgresWalletRepository {
	return &PostgresWalletRepository{
		db:  db,
		log: log,
	}
}

// Create crea una nueva billetera
func (r *PostgresWalletRepository) Create(wallet *domain.Wallet) error {
	// Generar UUID si no existe
	if wallet.UUID == "" {
		wallet.UUID = uuid.New().String()
	}

	// Validar antes de crear
	if err := wallet.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Verificar que no exista otra billetera para el usuario
	var existing domain.Wallet
	err := r.db.Where("user_id = ?", wallet.UserID).First(&existing).Error
	if err == nil {
		return errors.Wrap(errors.ErrConflict, fmt.Errorf("el usuario ya tiene una billetera"))
	}
	if err != gorm.ErrRecordNotFound {
		r.log.Error("Error verificando billetera existente",
			logger.Int64("user_id", wallet.UserID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Crear billetera
	if err := r.db.Create(wallet).Error; err != nil {
		r.log.Error("Error creando billetera",
			logger.Int64("user_id", wallet.UserID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	r.log.Info("Billetera creada exitosamente",
		logger.Int64("wallet_id", wallet.ID),
		logger.Int64("user_id", wallet.UserID))

	return nil
}

// FindByID busca una billetera por ID
func (r *PostgresWalletRepository) FindByID(id int64) (*domain.Wallet, error) {
	var wallet domain.Wallet

	if err := r.db.First(&wallet, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando billetera por ID",
			logger.Int64("id", id),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &wallet, nil
}

// FindByUUID busca una billetera por UUID
func (r *PostgresWalletRepository) FindByUUID(uuid string) (*domain.Wallet, error) {
	var wallet domain.Wallet

	if err := r.db.Where("uuid = ?", uuid).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando billetera por UUID",
			logger.String("uuid", uuid),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &wallet, nil
}

// FindByUserID busca una billetera por ID de usuario
func (r *PostgresWalletRepository) FindByUserID(userID int64) (*domain.Wallet, error) {
	var wallet domain.Wallet

	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error buscando billetera por user_id",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &wallet, nil
}

// Update actualiza una billetera existente
func (r *PostgresWalletRepository) Update(wallet *domain.Wallet) error {
	// Validar antes de actualizar
	if err := wallet.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	if err := r.db.Save(wallet).Error; err != nil {
		r.log.Error("Error actualizando billetera",
			logger.Int64("wallet_id", wallet.ID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	r.log.Info("Billetera actualizada",
		logger.Int64("wallet_id", wallet.ID),
		logger.String("balance", wallet.Balance.String()))

	return nil
}

// UpdateBalance actualiza solo el saldo (optimización)
func (r *PostgresWalletRepository) UpdateBalance(walletID int64, balance decimal.Decimal) error {
	if balance.LessThan(decimal.Zero) {
		return errors.Wrap(errors.ErrValidationFailed, fmt.Errorf("el saldo no puede ser negativo"))
	}

	result := r.db.Model(&domain.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", balance)

	if result.Error != nil {
		r.log.Error("Error actualizando saldo",
			logger.Int64("wallet_id", walletID),
			logger.Error(result.Error))
		return errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// Lock adquiere un lock para operaciones concurrentes
func (r *PostgresWalletRepository) Lock(walletID int64) error {
	// Usar SELECT ... FOR UPDATE para lock pesimista
	var wallet domain.Wallet
	if err := r.db.Clauses(gorm.Locking{Strength: "UPDATE"}).
		First(&wallet, walletID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		r.log.Error("Error adquiriendo lock de billetera",
			logger.Int64("wallet_id", walletID),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// Unlock libera un lock (manejado por transacción)
func (r *PostgresWalletRepository) Unlock(walletID int64) error {
	// El lock se libera automáticamente al finalizar la transacción
	return nil
}

// WithTransaction ejecuta una función dentro de una transacción
func (r *PostgresWalletRepository) WithTransaction(fn func(repo domain.WalletRepository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return errors.Wrap(errors.ErrDatabaseError, tx.Error)
	}

	// Crear repositorio con la transacción
	txRepo := &PostgresWalletRepository{
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
