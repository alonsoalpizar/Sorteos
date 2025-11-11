package redis

import (
	"github.com/redis/go-redis/v9"

	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"
)

// NewLockService crea un nuevo servicio de locks distribuidos
func NewLockService(client *redis.Client) *redisinfra.LockService {
	return redisinfra.NewLockService(client)
}
