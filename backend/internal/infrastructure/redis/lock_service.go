package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrLockNotAcquired = errors.New("failed to acquire lock")
	ErrLockNotHeld     = errors.New("lock not held by this client")
)

// LockService provides distributed locking using Redis
type LockService struct {
	client *redis.Client
}

// NewLockService creates a new lock service
func NewLockService(client *redis.Client) *LockService {
	return &LockService{client: client}
}

// Lock represents a distributed lock
type Lock struct {
	key      string
	value    string
	ttl      time.Duration
	client   *redis.Client
	released bool
}

// AcquireLock attempts to acquire a lock with the given key
// Returns a Lock object if successful, or an error if the lock is already held
func (s *LockService) AcquireLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
	// Generate a unique value for this lock acquisition
	value := fmt.Sprintf("%d", time.Now().UnixNano())

	// Try to set the key with NX (only if it doesn't exist) and expiration
	success, err := s.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return nil, fmt.Errorf("redis set error: %w", err)
	}

	if !success {
		return nil, ErrLockNotAcquired
	}

	return &Lock{
		key:      key,
		value:    value,
		ttl:      ttl,
		client:   s.client,
		released: false,
	}, nil
}

// AcquireMultipleLocks attempts to acquire multiple locks atomically
// If any lock cannot be acquired, all locks are released and an error is returned
func (s *LockService) AcquireMultipleLocks(ctx context.Context, keys []string, ttl time.Duration) ([]*Lock, error) {
	locks := make([]*Lock, 0, len(keys))

	// Try to acquire all locks
	for _, key := range keys {
		lock, err := s.AcquireLock(ctx, key, ttl)
		if err != nil {
			// Failed to acquire this lock - release all previously acquired locks
			for _, acquiredLock := range locks {
				_ = acquiredLock.Release(ctx) // Ignore errors during cleanup
			}
			return nil, fmt.Errorf("failed to acquire lock for key %s: %w", key, err)
		}
		locks = append(locks, lock)
	}

	return locks, nil
}

// Release releases the lock
func (l *Lock) Release(ctx context.Context) error {
	if l.released {
		return nil
	}

	// Use Lua script to ensure we only delete the lock if we own it
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return fmt.Errorf("redis eval error: %w", err)
	}

	if result == int64(0) {
		return ErrLockNotHeld
	}

	l.released = true
	return nil
}

// Extend extends the lock's TTL
func (l *Lock) Extend(ctx context.Context, additionalTTL time.Duration) error {
	if l.released {
		return errors.New("cannot extend released lock")
	}

	// Use Lua script to extend only if we own the lock
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value, additionalTTL.Milliseconds()).Result()
	if err != nil {
		return fmt.Errorf("redis eval error: %w", err)
	}

	if result == int64(0) {
		return ErrLockNotHeld
	}

	l.ttl = additionalTTL
	return nil
}

// IsLocked checks if a key is currently locked
func (s *LockService) IsLocked(ctx context.Context, key string) (bool, error) {
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// ReleaseMultipleLocks releases multiple locks
func ReleaseMultipleLocks(ctx context.Context, locks []*Lock) error {
	var firstError error
	for _, lock := range locks {
		if err := lock.Release(ctx); err != nil && firstError == nil {
			firstError = err
		}
	}
	return firstError
}

// ReservationLockKey generates a Redis lock key for a raffle number
func ReservationLockKey(raffleID, numberID string) string {
	return fmt.Sprintf("lock:reservation:%s:%s", raffleID, numberID)
}

// ForceReleaseLock forcefully releases a lock without verifying ownership
// Use this only for administrative operations like cancellation or expiration
func (s *LockService) ForceReleaseLock(ctx context.Context, key string) error {
	_, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis del error: %w", err)
	}
	return nil
}

// ForceReleaseMultipleLocks forcefully releases multiple locks
func (s *LockService) ForceReleaseMultipleLocks(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	_, err := s.client.Del(ctx, keys...).Result()
	if err != nil {
		return fmt.Errorf("redis del error: %w", err)
	}
	return nil
}
