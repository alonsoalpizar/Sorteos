package jobs

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/sorteos-platform/backend/internal/usecases"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// ExpireReservationsJob job para expirar reservas pendientes
type ExpireReservationsJob struct {
	reservationUseCases *usecases.ReservationUseCases
	logger              *logger.Logger
	interval            time.Duration
	stopChan            chan struct{}
}

// NewExpireReservationsJob crea un nuevo job de expiración
func NewExpireReservationsJob(
	reservationUseCases *usecases.ReservationUseCases,
	logger *logger.Logger,
	interval time.Duration,
) *ExpireReservationsJob {
	return &ExpireReservationsJob{
		reservationUseCases: reservationUseCases,
		logger:              logger,
		interval:            interval,
		stopChan:            make(chan struct{}),
	}
}

// Start inicia el job en background
func (j *ExpireReservationsJob) Start() {
	j.logger.Info("Starting expire reservations job", zap.Duration("interval", j.interval))

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	// Ejecutar inmediatamente al iniciar
	j.run()

	// Ejecutar periódicamente
	for {
		select {
		case <-ticker.C:
			j.run()
		case <-j.stopChan:
			j.logger.Info("Expire reservations job stopped")
			return
		}
	}
}

// Stop detiene el job
func (j *ExpireReservationsJob) Stop() {
	close(j.stopChan)
}

// run ejecuta el proceso de expiración
func (j *ExpireReservationsJob) run() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	count, err := j.reservationUseCases.ExpireReservations(ctx)
	duration := time.Since(start)

	if err != nil {
		j.logger.Error("Failed to expire reservations",
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return
	}

	if count > 0 {
		j.logger.Info("Expired reservations",
			zap.Int("count", count),
			zap.Duration("duration", duration),
		)
	}
}
