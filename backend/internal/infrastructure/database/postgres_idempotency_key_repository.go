package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/sorteos-platform/backend/internal/domain/entities"
	"github.com/sorteos-platform/backend/internal/domain/repositories"
)

// PostgresIdempotencyKeyRepository implements IdempotencyKeyRepository using PostgreSQL
type PostgresIdempotencyKeyRepository struct {
	db *gorm.DB
}

// NewPostgresIdempotencyKeyRepository creates a new PostgreSQL idempotency key repository
func NewPostgresIdempotencyKeyRepository(db *gorm.DB) repositories.IdempotencyKeyRepository {
	return &PostgresIdempotencyKeyRepository{db: db}
}

// Create stores a new idempotency key
func (r *PostgresIdempotencyKeyRepository) Create(ctx context.Context, key *entities.IdempotencyKey) error {
	return r.db.WithContext(ctx).Create(key).Error
}

// FindByKey retrieves an idempotency key by its key value and user ID
func (r *PostgresIdempotencyKeyRepository) FindByKey(ctx context.Context, key string, userID uuid.UUID) (*entities.IdempotencyKey, error) {
	var idempotencyKey entities.IdempotencyKey
	err := r.db.WithContext(ctx).
		Where("idempotency_key = ? AND user_id = ?", key, userID).
		First(&idempotencyKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &idempotencyKey, nil
}

// Update updates an existing idempotency key
func (r *PostgresIdempotencyKeyRepository) Update(ctx context.Context, key *entities.IdempotencyKey) error {
	return r.db.WithContext(ctx).Save(key).Error
}

// DeleteExpired removes expired idempotency keys
func (r *PostgresIdempotencyKeyRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entities.IdempotencyKey{}).Error
}
