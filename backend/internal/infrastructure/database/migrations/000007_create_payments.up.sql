CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id UUID NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    raffle_id UUID NOT NULL REFERENCES raffles(uuid) ON DELETE CASCADE,

    -- Stripe information
    stripe_payment_intent_id VARCHAR(255) UNIQUE NOT NULL,
    stripe_client_secret VARCHAR(255) NOT NULL,

    -- Payment details
    amount DECIMAL(10,2) NOT NULL, -- Total amount charged
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, succeeded, failed, cancelled, refunded

    -- Metadata
    payment_method VARCHAR(50), -- card, bank_transfer, etc.
    error_message TEXT, -- Stripe error message if failed
    metadata JSONB, -- Additional data from Stripe

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP, -- When payment succeeded

    CONSTRAINT payments_status_check CHECK (
        status IN ('pending', 'processing', 'succeeded', 'failed', 'cancelled', 'refunded')
    )
);

-- Index for finding user payments
CREATE INDEX idx_payments_user_id ON payments(user_id);

-- Index for finding raffle payments
CREATE INDEX idx_payments_raffle_id ON payments(raffle_id);

-- Index for Stripe webhook lookups
CREATE INDEX idx_payments_stripe_intent_id ON payments(stripe_payment_intent_id);

-- Index for reservation lookups
CREATE INDEX idx_payments_reservation_id ON payments(reservation_id);

-- Index for status filtering
CREATE INDEX idx_payments_status ON payments(status);

-- Add trigger for updated_at
CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
