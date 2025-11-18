package db

import (
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// PostgresSystemParameterRepository implementación de SystemParameterRepository con PostgreSQL
type PostgresSystemParameterRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewSystemParameterRepository crea una nueva instancia
func NewSystemParameterRepository(db *gorm.DB, log *logger.Logger) *PostgresSystemParameterRepository {
	return &PostgresSystemParameterRepository{
		db:  db,
		log: log,
	}
}

// GetByKey obtiene un parámetro por su key
func (r *PostgresSystemParameterRepository) GetByKey(key string) (*domain.SystemParameter, error) {
	var param domain.SystemParameter

	if err := r.db.Where("key = ?", key).First(&param).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		r.log.Error("Error getting system parameter by key",
			logger.String("key", key),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return &param, nil
}

// GetString obtiene un valor string con valor por defecto
func (r *PostgresSystemParameterRepository) GetString(key string, defaultValue string) (string, error) {
	param, err := r.GetByKey(key)
	if err != nil {
		if err == errors.ErrNotFound {
			return defaultValue, nil
		}
		return "", err
	}

	return param.GetString(), nil
}

// GetInt obtiene un valor int con valor por defecto
func (r *PostgresSystemParameterRepository) GetInt(key string, defaultValue int64) (int64, error) {
	param, err := r.GetByKey(key)
	if err != nil {
		if err == errors.ErrNotFound {
			return defaultValue, nil
		}
		return 0, err
	}

	value, err := param.GetInt()
	if err != nil {
		r.log.Error("Error parsing int parameter",
			logger.String("key", key),
			logger.Error(err))
		return defaultValue, nil
	}

	return value, nil
}

// GetFloat obtiene un valor float con valor por defecto
func (r *PostgresSystemParameterRepository) GetFloat(key string, defaultValue float64) (float64, error) {
	param, err := r.GetByKey(key)
	if err != nil {
		if err == errors.ErrNotFound {
			return defaultValue, nil
		}
		return 0, err
	}

	value, err := param.GetFloat()
	if err != nil {
		r.log.Error("Error parsing float parameter",
			logger.String("key", key),
			logger.Error(err))
		return defaultValue, nil
	}

	return value, nil
}

// GetBool obtiene un valor bool con valor por defecto
func (r *PostgresSystemParameterRepository) GetBool(key string, defaultValue bool) (bool, error) {
	param, err := r.GetByKey(key)
	if err != nil {
		if err == errors.ErrNotFound {
			return defaultValue, nil
		}
		return false, err
	}

	value, err := param.GetBool()
	if err != nil {
		r.log.Error("Error parsing bool parameter",
			logger.String("key", key),
			logger.Error(err))
		return defaultValue, nil
	}

	return value, nil
}

// GetJSON parsea un valor JSON
func (r *PostgresSystemParameterRepository) GetJSON(key string, target interface{}) error {
	param, err := r.GetByKey(key)
	if err != nil {
		return err
	}

	if err := param.GetJSON(target); err != nil {
		r.log.Error("Error parsing JSON parameter",
			logger.String("key", key),
			logger.Error(err))
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	return nil
}

// List obtiene parámetros con filtros y paginación
func (r *PostgresSystemParameterRepository) List(category *domain.ParameterCategory, offset, limit int) ([]*domain.SystemParameter, int64, error) {
	var params []*domain.SystemParameter
	var total int64

	query := r.db.Model(&domain.SystemParameter{})

	// Aplicar filtro de categoría si está presente
	if category != nil {
		query = query.Where("category = ?", *category)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		r.log.Error("Error counting system parameters", logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener registros con paginación
	if err := query.
		Order("category ASC, key ASC").
		Offset(offset).
		Limit(limit).
		Find(&params).Error; err != nil {
		r.log.Error("Error listing system parameters", logger.Error(err))
		return nil, 0, errors.Wrap(errors.ErrDatabaseError, err)
	}

	return params, total, nil
}

// ListByCategory obtiene parámetros agrupados por categoría
func (r *PostgresSystemParameterRepository) ListByCategory() (map[domain.ParameterCategory][]*domain.SystemParameter, error) {
	var params []*domain.SystemParameter

	if err := r.db.Order("category ASC, key ASC").Find(&params).Error; err != nil {
		r.log.Error("Error listing system parameters by category", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Agrupar por categoría
	result := make(map[domain.ParameterCategory][]*domain.SystemParameter)
	for _, param := range params {
		if param.Category != nil {
			result[*param.Category] = append(result[*param.Category], param)
		}
	}

	return result, nil
}

// Update actualiza un parámetro
func (r *PostgresSystemParameterRepository) Update(param *domain.SystemParameter, updatedBy int64) error {
	// Validar antes de actualizar
	if err := param.Validate(); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Actualizar el parámetro y el campo updated_by
	param.UpdatedByID = &updatedBy

	if err := r.db.Save(param).Error; err != nil {
		r.log.Error("Error updating system parameter",
			logger.String("key", param.Key),
			logger.Int64("updated_by", updatedBy),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	return nil
}

// Set establece el valor de un parámetro (crea si no existe)
func (r *PostgresSystemParameterRepository) Set(key string, value interface{}, valueType domain.ParameterValueType, updatedBy int64) error {
	// Buscar si el parámetro existe
	param, err := r.GetByKey(key)

	if err != nil && err != errors.ErrNotFound {
		return err
	}

	// Si no existe, crear nuevo
	if err == errors.ErrNotFound {
		param = &domain.SystemParameter{
			Key:       key,
			ValueType: valueType,
		}
	}

	// Establecer el valor
	if err := param.SetValue(value); err != nil {
		return errors.Wrap(errors.ErrValidationFailed, err)
	}

	// Guardar (crear o actualizar)
	param.UpdatedByID = &updatedBy

	if param.ID == 0 {
		// Crear nuevo
		if err := r.db.Create(param).Error; err != nil {
			r.log.Error("Error creating system parameter",
				logger.String("key", key),
				logger.Error(err))
			return errors.Wrap(errors.ErrDatabaseError, err)
		}
	} else {
		// Actualizar existente
		if err := r.db.Save(param).Error; err != nil {
			r.log.Error("Error updating system parameter",
				logger.String("key", key),
				logger.Error(err))
			return errors.Wrap(errors.ErrDatabaseError, err)
		}
	}

	return nil
}
