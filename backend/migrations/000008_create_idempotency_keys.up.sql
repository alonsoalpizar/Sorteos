CREATE TABLE IF NOT EXISTS idempotency_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL, -- Client-provided key
    user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,

    -- Request fingerprint
    request_path VARCHAR(255) NOT NULL, -- e.g., "/api/v1/payments"
    request_params JSONB, -- Serialized request body for verification

    -- Response cache
    response_status_code INTEGER,
    response_body JSONB,

    -- State tracking
    status VARCHAR(20) NOT NULL DEFAULT 'processing', -- processing, completed, failed

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,

    -- Auto-expire after 24 hours (for cleanup)
    expires_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours'),

    CONSTRAINT idempotency_keys_status_check CHECK (status IN ('processing', 'completed', 'failed'))
);

-- Index for key lookups (primary use case)
CREATE INDEX idx_idempotency_keys_key ON idempotency_keys(idempotency_key);

-- Index for user lookups
CREATE INDEX idx_idempotency_keys_user_id ON idempotency_keys(user_id);

-- Index for cleanup of expired keys
CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);

-- Composite index for key + user verification
CREATE INDEX idx_idempotency_keys_key_user ON idempotency_keys(idempotency_key, user_id);
