package db

import (
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain/repositories"
	"github.com/sorteos-platform/backend/internal/infrastructure/database"
)

// NewReservationRepository crea un nuevo repositorio de reservas
func NewReservationRepository(db *gorm.DB) repositories.ReservationRepository {
	return database.NewPostgresReservationRepository(db)
}
