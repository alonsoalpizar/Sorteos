package repository

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// SystemConfig representa una configuración del sistema
type SystemConfig struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Key       string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Value     string    `gorm:"type:jsonb;not null"` // JSON value
	Category  string    `gorm:"type:varchar(100);index"`
	UpdatedAt time.Time `gorm:"not null"`
	UpdatedBy *int64    `gorm:"index"` // Admin ID que actualizó
}

// TableName especifica el nombre de la tabla
func (SystemConfig) TableName() string {
	return "system_parameters"
}

// SystemConfigRepository interfaz para el repositorio de configuración del sistema
type SystemConfigRepository interface {
	Get(ctx context.Context, key string) (*SystemConfig, error)
	GetByCategory(ctx context.Context, category string) ([]*SystemConfig, error)
	GetAll(ctx context.Context) ([]*SystemConfig, error)
	Set(ctx context.Context, key, value string, category string, updatedBy int64) error
	Delete(ctx context.Context, key string) error
}

// systemConfigRepository implementación del repositorio
type systemConfigRepository struct {
	db *gorm.DB
}

// NewSystemConfigRepository crea una nueva instancia del repositorio
func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &systemConfigRepository{db: db}
}

// Get obtiene una configuración por key
func (r *systemConfigRepository) Get(ctx context.Context, key string) (*SystemConfig, error) {
	var config SystemConfig
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetByCategory obtiene todas las configuraciones de una categoría
func (r *systemConfigRepository) GetByCategory(ctx context.Context, category string) ([]*SystemConfig, error) {
	var configs []*SystemConfig
	if err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("key ASC").
		Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetAll obtiene todas las configuraciones
func (r *systemConfigRepository) GetAll(ctx context.Context) ([]*SystemConfig, error) {
	var configs []*SystemConfig
	if err := r.db.WithContext(ctx).
		Order("category ASC, key ASC").
		Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// Set crea o actualiza una configuración (UPSERT)
func (r *systemConfigRepository) Set(ctx context.Context, key, value string, category string, updatedBy int64) error {
	// Validar que value sea JSON válido
	var js interface{}
	if err := json.Unmarshal([]byte(value), &js); err != nil {
		return err
	}

	now := time.Now()
	config := &SystemConfig{
		Key:       key,
		Value:     value,
		Category:  category,
		UpdatedAt: now,
		UpdatedBy: &updatedBy,
	}

	// UPSERT: Si existe, actualizar; si no, crear
	result := r.db.WithContext(ctx).
		Where("key = ?", key).
		Assign(map[string]interface{}{
			"value":      value,
			"category":   category,
			"updated_at": now,
			"updated_by": updatedBy,
		}).
		FirstOrCreate(config)

	return result.Error
}

// Delete elimina una configuración
func (r *systemConfigRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key = ?", key).Delete(&SystemConfig{}).Error
}
