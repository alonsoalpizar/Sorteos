package db

import (
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain/repositories"
	"github.com/sorteos-platform/backend/internal/infrastructure/database"
)

// NewPaymentRepository crea un nuevo repositorio de pagos
func NewPaymentRepository(db *gorm.DB) repositories.PaymentRepository {
	return database.NewPostgresPaymentRepository(db)
}
