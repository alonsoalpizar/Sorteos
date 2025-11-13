package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/adapters/db"
	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"
	"github.com/sorteos-platform/backend/internal/infrastructure/websocket"
	"github.com/sorteos-platform/backend/internal/usecases"
	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// startBackgroundJobs inicia todos los jobs de fondo
func startBackgroundJobs(gormDB *gorm.DB, rdb *redis.Client, wsHub *websocket.Hub, cfg *config.Config, log *logger.Logger) {
	// Inicializar repositorios
	reservationRepo := db.NewReservationRepository(gormDB)
	raffleRepo := db.NewRaffleRepository(gormDB)
	raffleNumberRepo := db.NewRaffleNumberRepository(gormDB)

	// Inicializar lock service
	lockService := redisinfra.NewLockService(rdb)

	// Inicializar use case
	reservationUseCases := usecases.NewReservationUseCases(
		reservationRepo,
		raffleRepo,
		raffleNumberRepo,
		lockService,
		wsHub,
	)

	// Job de expiración de reservas (ejecutar cada 30 segundos)
	go startReservationExpirationJob(reservationUseCases, log)

	log.Info("Background jobs started")
}

// startReservationExpirationJob inicia el job de expiración de reservas
func startReservationExpirationJob(reservationUC *usecases.ReservationUseCases, log *logger.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Info("Starting expire reservations job", logger.String("interval", "30s"))

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)

		count, err := reservationUC.ExpireOldReservations(ctx)
		if err != nil {
			log.Error("Error expiring reservations", logger.Error(err))
		} else if count > 0 {
			log.Info("Expired reservations", logger.Int("count", count))
		}

		cancel()
	}
}
