-- Migration: 000012_settlements
-- Purpose: Registrar liquidaciones y pagos a organizadores

CREATE TYPE settlement_status AS ENUM (
    'pending', 'approved', 'paid', 'rejected'
);

CREATE TABLE settlements (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE DEFAULT uuid_generate_v4(),
    raffle_id BIGINT NOT NULL REFERENCES raffles(id) ON DELETE RESTRICT,
    organizer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    gross_revenue DECIMAL(12,2) NOT NULL,
    platform_fee DECIMAL(12,2) NOT NULL,
    platform_fee_percentage DECIMAL(5,2) NOT NULL,
    net_payout DECIMAL(12,2) NOT NULL,
    status settlement_status DEFAULT 'pending',
    payment_method VARCHAR(50),
    payment_reference VARCHAR(255),
    approved_by BIGINT REFERENCES users(id),
    approved_at TIMESTAMP,
    paid_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_settlements_raffle_id ON settlements(raffle_id);
CREATE INDEX idx_settlements_organizer_id ON settlements(organizer_id);
CREATE INDEX idx_settlements_status ON settlements(status, created_at DESC);
CREATE INDEX idx_settlements_approved_by ON settlements(approved_by);

CREATE TRIGGER update_settlements_updated_at
    BEFORE UPDATE ON settlements
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE UNIQUE INDEX idx_settlements_raffle_unique ON settlements(raffle_id);

ALTER TABLE settlements
    ADD CONSTRAINT chk_settlements_net_payout
    CHECK (net_payout = gross_revenue - platform_fee);

ALTER TABLE settlements
    ADD CONSTRAINT chk_settlements_fee_percentage
    CHECK (platform_fee_percentage >= 0 AND platform_fee_percentage <= 50);

COMMENT ON TABLE settlements IS 'Liquidaciones y pagos a organizadores';
COMMENT ON COLUMN settlements.net_payout IS 'gross_revenue - platform_fee';
