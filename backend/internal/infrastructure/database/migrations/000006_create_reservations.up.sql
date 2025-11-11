CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    raffle_id UUID NOT NULL REFERENCES raffles(uuid) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    number_ids TEXT[] NOT NULL, -- Array of number identifiers (e.g., ["0001", "0042", "0123"])
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, confirmed, expired, cancelled
    session_id VARCHAR(255) NOT NULL, -- For idempotency and tracking
    total_amount DECIMAL(10,2) NOT NULL,
    expires_at TIMESTAMP NOT NULL, -- Auto-expire after 5 minutes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT reservations_status_check CHECK (status IN ('pending', 'confirmed', 'expired', 'cancelled'))
);

-- Index for finding user reservations
CREATE INDEX idx_reservations_user_id ON reservations(user_id);

-- Index for finding raffle reservations
CREATE INDEX idx_reservations_raffle_id ON reservations(raffle_id);

-- Index for finding expired reservations (for cron job)
CREATE INDEX idx_reservations_expires_at ON reservations(expires_at) WHERE status = 'pending';

-- Index for session-based lookups (idempotency)
CREATE INDEX idx_reservations_session_id ON reservations(session_id);

-- Composite index for status + expires_at (cleanup queries)
CREATE INDEX idx_reservations_status_expires ON reservations(status, expires_at);

-- Add trigger for updated_at
CREATE TRIGGER update_reservations_updated_at
    BEFORE UPDATE ON reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
