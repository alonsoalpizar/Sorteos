package system

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// DatabaseHealth salud de la base de datos
type DatabaseHealth struct {
	Status          string  `json:"status"`           // healthy, degraded, down
	ResponseTime    float64 `json:"response_time_ms"` // Tiempo de respuesta en ms
	ConnectionCount int     `json:"connection_count"`
	Error           *string `json:"error,omitempty"`
}

// CacheHealth salud del cache (Redis)
type CacheHealth struct {
	Status       string  `json:"status"` // healthy, degraded, down
	ResponseTime float64 `json:"response_time_ms"`
	Error        *string `json:"error,omitempty"`
}

// SystemMetrics métricas del sistema
type SystemMetrics struct {
	TotalUsers       int64 `json:"total_users"`
	TotalRaffles     int64 `json:"total_raffles"`
	ActiveRaffles    int64 `json:"active_raffles"`
	TotalPayments    int64 `json:"total_payments"`
	TotalSettlements int64 `json:"total_settlements"`
}

// ViewSystemHealthOutput resultado
type ViewSystemHealthOutput struct {
	OverallStatus string          `json:"overall_status"` // healthy, degraded, down
	Database      *DatabaseHealth `json:"database"`
	Cache         *CacheHealth    `json:"cache"`
	Metrics       *SystemMetrics  `json:"metrics"`
	Uptime        float64         `json:"uptime_hours"` // Horas de uptime
	Timestamp     string          `json:"timestamp"`
	Version       string          `json:"version"`
}

// ViewSystemHealthUseCase caso de uso para ver salud del sistema
type ViewSystemHealthUseCase struct {
	db        *gorm.DB
	log       *logger.Logger
	startTime time.Time
}

// NewViewSystemHealthUseCase crea una nueva instancia
func NewViewSystemHealthUseCase(db *gorm.DB, log *logger.Logger) *ViewSystemHealthUseCase {
	return &ViewSystemHealthUseCase{
		db:        db,
		log:       log,
		startTime: time.Now(), // TODO: Obtener de variable global
	}
}

// Execute ejecuta el caso de uso
func (uc *ViewSystemHealthUseCase) Execute(ctx context.Context, adminID int64) (*ViewSystemHealthOutput, error) {
	output := &ViewSystemHealthOutput{
		OverallStatus: "healthy",
		Timestamp:     time.Now().Format(time.RFC3339),
		Version:       "1.0.0", // TODO: Obtener de configuración
	}

	// Calcular uptime
	uptime := time.Since(uc.startTime).Hours()
	output.Uptime = uptime

	// Check Database Health
	dbHealth := uc.checkDatabaseHealth(ctx)
	output.Database = dbHealth

	// Check Cache Health (Redis)
	cacheHealth := uc.checkCacheHealth(ctx)
	output.Cache = cacheHealth

	// Get System Metrics
	metrics := uc.getSystemMetrics(ctx)
	output.Metrics = metrics

	// Determinar overall status
	if dbHealth.Status == "down" {
		output.OverallStatus = "down"
	} else if dbHealth.Status == "degraded" || cacheHealth.Status == "degraded" {
		output.OverallStatus = "degraded"
	}

	// Log auditoría
	uc.log.Info("Admin viewed system health",
		logger.Int64("admin_id", adminID),
		logger.String("overall_status", output.OverallStatus),
		logger.String("action", "admin_view_system_health"))

	return output, nil
}

// checkDatabaseHealth verifica la salud de la base de datos
func (uc *ViewSystemHealthUseCase) checkDatabaseHealth(ctx context.Context) *DatabaseHealth {
	health := &DatabaseHealth{
		Status: "healthy",
	}

	// Medir tiempo de respuesta
	start := time.Now()

	// Simple query para verificar conectividad
	var count int64
	err := uc.db.WithContext(ctx).Raw("SELECT 1").Count(&count).Error

	responseTime := time.Since(start).Milliseconds()
	health.ResponseTime = float64(responseTime)

	if err != nil {
		health.Status = "down"
		errMsg := err.Error()
		health.Error = &errMsg
		uc.log.Error("Database health check failed", logger.Error(err))
		return health
	}

	// Verificar si el tiempo de respuesta es alto
	if responseTime > 1000 {
		health.Status = "degraded"
		msg := "High response time"
		health.Error = &msg
	}

	// Obtener número de conexiones (PostgreSQL específico)
	var connCount int
	uc.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM pg_stat_activity").Scan(&connCount)
	health.ConnectionCount = connCount

	return health
}

// checkCacheHealth verifica la salud del cache (Redis)
func (uc *ViewSystemHealthUseCase) checkCacheHealth(ctx context.Context) *CacheHealth {
	health := &CacheHealth{
		Status: "healthy",
	}

	// TODO: Implementar check de Redis cuando esté configurado
	// Por ahora, marcamos como N/A o skip

	start := time.Now()

	// Simulación de ping a Redis
	// err := redisClient.Ping(ctx).Err()

	responseTime := time.Since(start).Milliseconds()
	health.ResponseTime = float64(responseTime)

	// Si Redis no está configurado, marcar como N/A
	msg := "Redis not configured"
	health.Error = &msg
	health.Status = "healthy" // No afecta overall status si no está configurado

	return health
}

// getSystemMetrics obtiene métricas del sistema
func (uc *ViewSystemHealthUseCase) getSystemMetrics(ctx context.Context) *SystemMetrics {
	metrics := &SystemMetrics{}

	// Total users
	uc.db.WithContext(ctx).Table("users").Count(&metrics.TotalUsers)

	// Total raffles
	uc.db.WithContext(ctx).Table("raffles").Where("deleted_at IS NULL").Count(&metrics.TotalRaffles)

	// Active raffles
	uc.db.WithContext(ctx).Table("raffles").
		Where("status = ?", "active").
		Where("deleted_at IS NULL").
		Count(&metrics.ActiveRaffles)

	// Total payments
	uc.db.WithContext(ctx).Table("payments").Count(&metrics.TotalPayments)

	// Total settlements
	uc.db.WithContext(ctx).Table("settlements").Count(&metrics.TotalSettlements)

	return metrics
}
