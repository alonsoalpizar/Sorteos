-- Migration: 000011_organizer_profiles
-- Purpose: Perfiles extendidos de organizadores con info bancaria y comisiones

CREATE TABLE organizer_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(255),
    tax_id VARCHAR(50),
    bank_name VARCHAR(100),
    bank_account_number TEXT,
    bank_account_type VARCHAR(20),
    bank_account_holder VARCHAR(255),
    payout_schedule VARCHAR(20) DEFAULT 'manual',
    commission_override DECIMAL(5,2),
    total_payouts DECIMAL(12,2) DEFAULT 0.00,
    pending_payout DECIMAL(12,2) DEFAULT 0.00,
    verified BOOLEAN DEFAULT false,
    verified_at TIMESTAMP,
    verified_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_organizer_profiles_user_id ON organizer_profiles(user_id);
CREATE INDEX idx_organizer_profiles_verified ON organizer_profiles(verified, created_at DESC);
CREATE INDEX idx_organizer_profiles_total_payouts ON organizer_profiles(total_payouts DESC);

CREATE TRIGGER update_organizer_profiles_updated_at
    BEFORE UPDATE ON organizer_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE organizer_profiles IS 'Perfiles extendidos de organizadores';
COMMENT ON COLUMN organizer_profiles.bank_account_number IS 'Encrypted in application layer';
COMMENT ON COLUMN organizer_profiles.commission_override IS 'NULL = use global default from system_parameters';
