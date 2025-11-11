package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// IdempotencyKeyStatus represents the state of an idempotency key
type IdempotencyKeyStatus string

const (
	IdempotencyKeyStatusProcessing IdempotencyKeyStatus = "processing"
	IdempotencyKeyStatusCompleted  IdempotencyKeyStatus = "completed"
	IdempotencyKeyStatusFailed     IdempotencyKeyStatus = "failed"
)

// IdempotencyKeyExpirationDuration is how long to keep idempotency keys (24 hours)
const IdempotencyKeyExpirationDuration = 24 * time.Hour

var (
	ErrIdempotencyKeyConflict = errors.New("idempotency key conflict: different request with same key")
	ErrIdempotencyKeyExpired  = errors.New("idempotency key has expired")
)

// IdempotencyKey ensures that duplicate requests are not processed twice
type IdempotencyKey struct {
	ID                 uuid.UUID            `json:"id"`
	IdempotencyKey     string               `json:"idempotency_key"`
	UserID             uuid.UUID            `json:"user_id"`
	RequestPath        string               `json:"request_path"`
	RequestParams      string               `json:"request_params,omitempty"` // JSONB as string
	ResponseStatusCode int                  `json:"response_status_code,omitempty"`
	ResponseBody       string               `json:"response_body,omitempty"` // JSONB as string
	Status             IdempotencyKeyStatus `json:"status"`
	CreatedAt          time.Time            `json:"created_at"`
	CompletedAt        *time.Time           `json:"completed_at,omitempty"`
	ExpiresAt          time.Time            `json:"expires_at"`
}

// NewIdempotencyKey creates a new idempotency key record
func NewIdempotencyKey(key string, userID uuid.UUID, requestPath string, requestBody interface{}) (*IdempotencyKey, error) {
	now := time.Now()

	var requestParamsJSON string
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		requestParamsJSON = string(jsonData)
	}

	return &IdempotencyKey{
		ID:             uuid.New(),
		IdempotencyKey: key,
		UserID:         userID,
		RequestPath:    requestPath,
		RequestParams:  requestParamsJSON,
		Status:         IdempotencyKeyStatusProcessing,
		CreatedAt:      now,
		ExpiresAt:      now.Add(IdempotencyKeyExpirationDuration),
	}, nil
}

// MarkAsCompleted marks the idempotency key as completed with response
func (ik *IdempotencyKey) MarkAsCompleted(statusCode int, responseBody interface{}) error {
	now := time.Now()

	var responseBodyJSON string
	if responseBody != nil {
		jsonData, err := json.Marshal(responseBody)
		if err != nil {
			return err
		}
		responseBodyJSON = string(jsonData)
	}

	ik.Status = IdempotencyKeyStatusCompleted
	ik.ResponseStatusCode = statusCode
	ik.ResponseBody = responseBodyJSON
	ik.CompletedAt = &now
	return nil
}

// MarkAsFailed marks the idempotency key as failed
func (ik *IdempotencyKey) MarkAsFailed(statusCode int, errorResponse interface{}) error {
	now := time.Now()

	var errorBodyJSON string
	if errorResponse != nil {
		jsonData, err := json.Marshal(errorResponse)
		if err != nil {
			return err
		}
		errorBodyJSON = string(jsonData)
	}

	ik.Status = IdempotencyKeyStatusFailed
	ik.ResponseStatusCode = statusCode
	ik.ResponseBody = errorBodyJSON
	ik.CompletedAt = &now
	return nil
}

// VerifyRequestMatch checks if the current request matches the stored one
func (ik *IdempotencyKey) VerifyRequestMatch(requestPath string, requestBody interface{}) error {
	// Check path match
	if ik.RequestPath != requestPath {
		return ErrIdempotencyKeyConflict
	}

	// Check request body match
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		currentHash := hashJSON(string(jsonData))
		storedHash := hashJSON(ik.RequestParams)

		if currentHash != storedHash {
			return ErrIdempotencyKeyConflict
		}
	}

	return nil
}

// IsExpired checks if the idempotency key has expired
func (ik *IdempotencyKey) IsExpired() bool {
	return time.Now().After(ik.ExpiresAt)
}

// GetResponseBody unmarshals the stored response body
func (ik *IdempotencyKey) GetResponseBody(target interface{}) error {
	if ik.ResponseBody == "" {
		return nil
	}
	return json.Unmarshal([]byte(ik.ResponseBody), target)
}

// hashJSON creates a SHA-256 hash of JSON data for comparison
func hashJSON(jsonStr string) string {
	hash := sha256.Sum256([]byte(jsonStr))
	return hex.EncodeToString(hash[:])
}
