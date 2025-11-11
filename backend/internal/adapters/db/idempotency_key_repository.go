package db

import (
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain/repositories"
	"github.com/sorteos-platform/backend/internal/infrastructure/database"
)

// NewIdempotencyKeyRepository crea un nuevo repositorio de claves de idempotencia
func NewIdempotencyKeyRepository(db *gorm.DB) repositories.IdempotencyKeyRepository {
	return database.NewPostgresIdempotencyKeyRepository(db)
}
