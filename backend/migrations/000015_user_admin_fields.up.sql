-- Migration: 000015_user_admin_fields
-- Purpose: Agregar campos de administración a tabla users

ALTER TABLE users
    ADD COLUMN suspension_reason TEXT,
    ADD COLUMN suspended_by BIGINT REFERENCES users(id),
    ADD COLUMN suspended_at TIMESTAMP,
    ADD COLUMN last_kyc_review TIMESTAMP,
    ADD COLUMN kyc_reviewer BIGINT REFERENCES users(id);

CREATE INDEX idx_users_suspended_by ON users(suspended_by) WHERE suspended_by IS NOT NULL;
CREATE INDEX idx_users_kyc_reviewer ON users(kyc_reviewer) WHERE kyc_reviewer IS NOT NULL;

COMMENT ON COLUMN users.suspension_reason IS 'Razón de suspensión/ban';
COMMENT ON COLUMN users.suspended_by IS 'Admin que suspendió al usuario';
COMMENT ON COLUMN users.last_kyc_review IS 'Última revisión de KYC por admin';
COMMENT ON COLUMN users.kyc_reviewer IS 'Admin que revisó/cambió KYC';
