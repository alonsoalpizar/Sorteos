package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// ParameterValueType representa el tipo de valor de un parámetro
type ParameterValueType string

const (
	ParameterTypeString ParameterValueType = "string"
	ParameterTypeInt    ParameterValueType = "int"
	ParameterTypeFloat  ParameterValueType = "float"
	ParameterTypeBool   ParameterValueType = "bool"
	ParameterTypeJSON   ParameterValueType = "json"
)

// ParameterCategory representa la categoría de un parámetro
type ParameterCategory string

const (
	ParameterCategoryBusiness ParameterCategory = "business"
	ParameterCategoryPayment  ParameterCategory = "payment"
	ParameterCategorySecurity ParameterCategory = "security"
	ParameterCategoryEmail    ParameterCategory = "email"
	ParameterCategorySystem   ParameterCategory = "system"
)

// SystemParameter representa un parámetro de configuración del sistema
type SystemParameter struct {
	ID int64 `json:"id" gorm:"primaryKey"`

	// Parameter Key (unique identifier)
	Key string `json:"key" gorm:"uniqueIndex;not null;size:100"`

	// Value (stored as text, parsed according to ValueType)
	Value     string             `json:"value" gorm:"not null"`
	ValueType ParameterValueType `json:"value_type" gorm:"type:varchar(20);default:'string'"`

	// Metadata
	Description *string            `json:"description,omitempty"`
	Category    *ParameterCategory `json:"category,omitempty" gorm:"type:varchar(50)"`
	IsSensitive bool               `json:"is_sensitive" gorm:"default:false"` // Si es true, no mostrar en logs

	// Audit
	UpdatedBy   *int64     `json:"updated_by,omitempty"` // Último admin que modificó
	UpdatedByID *int64     `json:"-" gorm:"column:updated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (SystemParameter) TableName() string {
	return "system_parameters"
}

// Validate valida los campos del SystemParameter
func (sp *SystemParameter) Validate() error {
	// Key es requerido
	if sp.Key == "" {
		return fmt.Errorf("key is required")
	}

	if len(sp.Key) > 100 {
		return fmt.Errorf("key is too long (max 100 characters)")
	}

	// Value es requerido
	if sp.Value == "" {
		return fmt.Errorf("value is required")
	}

	// Validar ValueType
	validTypes := []ParameterValueType{
		ParameterTypeString,
		ParameterTypeInt,
		ParameterTypeFloat,
		ParameterTypeBool,
		ParameterTypeJSON,
	}
	validType := false
	for _, vt := range validTypes {
		if sp.ValueType == vt {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid value_type: %s (valid: string, int, float, bool, json)", sp.ValueType)
	}

	// Validar que el valor coincida con el tipo
	if err := sp.ValidateValueForType(); err != nil {
		return err
	}

	return nil
}

// ValidateValueForType valida que el valor coincida con el tipo declarado
func (sp *SystemParameter) ValidateValueForType() error {
	switch sp.ValueType {
	case ParameterTypeInt:
		if _, err := strconv.ParseInt(sp.Value, 10, 64); err != nil {
			return fmt.Errorf("value is not a valid integer: %w", err)
		}
	case ParameterTypeFloat:
		if _, err := strconv.ParseFloat(sp.Value, 64); err != nil {
			return fmt.Errorf("value is not a valid float: %w", err)
		}
	case ParameterTypeBool:
		if sp.Value != "true" && sp.Value != "false" {
			return fmt.Errorf("value must be 'true' or 'false' for bool type")
		}
	case ParameterTypeJSON:
		var jsonTest interface{}
		if err := json.Unmarshal([]byte(sp.Value), &jsonTest); err != nil {
			return fmt.Errorf("value is not valid JSON: %w", err)
		}
	case ParameterTypeString:
		// String siempre es válido
	default:
		return fmt.Errorf("unknown value_type: %s", sp.ValueType)
	}
	return nil
}

// GetString obtiene el valor como string
func (sp *SystemParameter) GetString() string {
	return sp.Value
}

// GetInt obtiene el valor como int64
func (sp *SystemParameter) GetInt() (int64, error) {
	if sp.ValueType != ParameterTypeInt {
		return 0, fmt.Errorf("parameter %s is not an int (type: %s)", sp.Key, sp.ValueType)
	}
	return strconv.ParseInt(sp.Value, 10, 64)
}

// GetFloat obtiene el valor como float64
func (sp *SystemParameter) GetFloat() (float64, error) {
	if sp.ValueType != ParameterTypeFloat {
		return 0, fmt.Errorf("parameter %s is not a float (type: %s)", sp.Key, sp.ValueType)
	}
	return strconv.ParseFloat(sp.Value, 64)
}

// GetBool obtiene el valor como bool
func (sp *SystemParameter) GetBool() (bool, error) {
	if sp.ValueType != ParameterTypeBool {
		return false, fmt.Errorf("parameter %s is not a bool (type: %s)", sp.Key, sp.ValueType)
	}
	return sp.Value == "true", nil
}

// GetJSON parsea el valor como JSON
func (sp *SystemParameter) GetJSON(target interface{}) error {
	if sp.ValueType != ParameterTypeJSON {
		return fmt.Errorf("parameter %s is not JSON (type: %s)", sp.Key, sp.ValueType)
	}
	return json.Unmarshal([]byte(sp.Value), target)
}

// SetValue establece el valor con validación de tipo
func (sp *SystemParameter) SetValue(value interface{}) error {
	switch sp.ValueType {
	case ParameterTypeString:
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("value must be a string")
		}
		sp.Value = strValue

	case ParameterTypeInt:
		var intValue int64
		switch v := value.(type) {
		case int:
			intValue = int64(v)
		case int64:
			intValue = v
		case float64:
			intValue = int64(v)
		default:
			return fmt.Errorf("value must be an integer")
		}
		sp.Value = strconv.FormatInt(intValue, 10)

	case ParameterTypeFloat:
		var floatValue float64
		switch v := value.(type) {
		case float64:
			floatValue = v
		case float32:
			floatValue = float64(v)
		case int:
			floatValue = float64(v)
		case int64:
			floatValue = float64(v)
		default:
			return fmt.Errorf("value must be a number")
		}
		sp.Value = strconv.FormatFloat(floatValue, 'f', -1, 64)

	case ParameterTypeBool:
		boolValue, ok := value.(bool)
		if !ok {
			return fmt.Errorf("value must be a boolean")
		}
		sp.Value = strconv.FormatBool(boolValue)

	case ParameterTypeJSON:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value to JSON: %w", err)
		}
		sp.Value = string(jsonBytes)

	default:
		return fmt.Errorf("unknown value_type: %s", sp.ValueType)
	}

	return sp.ValidateValueForType()
}

// MaskIfSensitive enmascara el valor si el parámetro es sensible
func (sp *SystemParameter) MaskIfSensitive() *SystemParameter {
	if !sp.IsSensitive {
		return sp
	}

	masked := *sp
	if len(masked.Value) > 4 {
		masked.Value = "****" + masked.Value[len(masked.Value)-4:]
	} else {
		masked.Value = "****"
	}
	return &masked
}

// SystemParameterRepository define los métodos de acceso a datos
type SystemParameterRepository interface {
	// GetByKey obtiene un parámetro por su key
	GetByKey(key string) (*SystemParameter, error)

	// GetString obtiene un valor string con valor por defecto
	GetString(key string, defaultValue string) (string, error)

	// GetInt obtiene un valor int con valor por defecto
	GetInt(key string, defaultValue int64) (int64, error)

	// GetFloat obtiene un valor float con valor por defecto
	GetFloat(key string, defaultValue float64) (float64, error)

	// GetBool obtiene un valor bool con valor por defecto
	GetBool(key string, defaultValue bool) (bool, error)

	// GetJSON parsea un valor JSON
	GetJSON(key string, target interface{}) error

	// List obtiene parámetros con filtros y paginación
	List(category *ParameterCategory, offset, limit int) ([]*SystemParameter, int64, error)

	// ListByCategory obtiene parámetros agrupados por categoría
	ListByCategory() (map[ParameterCategory][]*SystemParameter, error)

	// Update actualiza un parámetro
	Update(param *SystemParameter, updatedBy int64) error

	// Set establece el valor de un parámetro (crea si no existe)
	Set(key string, value interface{}, valueType ParameterValueType, updatedBy int64) error
}
