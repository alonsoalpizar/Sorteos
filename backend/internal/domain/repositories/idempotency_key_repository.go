package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/sorteos-platform/backend/internal/domain/entities"
)

// IdempotencyKeyRepository defines the interface for idempotency key persistence
type IdempotencyKeyRepository interface {
	// Create stores a new idempotency key
	Create(ctx context.Context, key *entities.IdempotencyKey) error

	// FindByKey retrieves an idempotency key by its key value
	FindByKey(ctx context.Context, key string, userID uuid.UUID) (*entities.IdempotencyKey, error)

	// Update updates an existing idempotency key
	Update(ctx context.Context, key *entities.IdempotencyKey) error

	// DeleteExpired removes expired idempotency keys (cleanup job)
	DeleteExpired(ctx context.Context) error
}
