package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklistService gestiona tokens revocados en Redis
type TokenBlacklistService struct {
	client *redis.Client
}

// NewTokenBlacklistService crea una nueva instancia del servicio
func NewTokenBlacklistService(client *redis.Client) *TokenBlacklistService {
	return &TokenBlacklistService{
		client: client,
	}
}

// AddToBlacklist agrega un token a la lista negra
// El token expirará automáticamente en Redis después del TTL especificado
func (s *TokenBlacklistService) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:token:%s", token)

	// Guardar en Redis con expiración automática
	// El valor no importa, solo necesitamos verificar la existencia de la key
	err := s.client.Set(ctx, key, "revoked", ttl).Err()
	if err != nil {
		return fmt.Errorf("error adding token to blacklist: %w", err)
	}

	return nil
}

// IsBlacklisted verifica si un token está en la lista negra
func (s *TokenBlacklistService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", token)

	// Verificar si la key existe
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error checking token blacklist: %w", err)
	}

	return exists > 0, nil
}

// RemoveFromBlacklist elimina un token de la lista negra (raramente usado)
func (s *TokenBlacklistService) RemoveFromBlacklist(ctx context.Context, token string) error {
	key := fmt.Sprintf("blacklist:token:%s", token)

	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error removing token from blacklist: %w", err)
	}

	return nil
}

// CountBlacklistedTokens retorna el número de tokens en la blacklist
// Útil para métricas y monitoreo
func (s *TokenBlacklistService) CountBlacklistedTokens(ctx context.Context) (int64, error) {
	// Buscar todas las keys que coincidan con el patrón
	keys, err := s.client.Keys(ctx, "blacklist:token:*").Result()
	if err != nil {
		return 0, fmt.Errorf("error counting blacklisted tokens: %w", err)
	}

	return int64(len(keys)), nil
}

// CleanupExpiredTokens limpia tokens expirados de la blacklist
// Nota: Redis ya hace esto automáticamente con el TTL,
// este método es solo por si necesitamos una limpieza manual
func (s *TokenBlacklistService) CleanupExpiredTokens(ctx context.Context) (int64, error) {
	// Este método no es necesario con Redis ya que el TTL lo maneja automáticamente
	// Lo dejamos como stub por si se necesita en el futuro
	return 0, nil
}
