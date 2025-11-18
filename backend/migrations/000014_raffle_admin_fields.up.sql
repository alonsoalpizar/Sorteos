-- Migration: 000014_raffle_admin_fields
-- Purpose: Agregar campos de administración a tabla raffles

ALTER TABLE raffles
    ADD COLUMN suspension_reason TEXT,
    ADD COLUMN suspended_by BIGINT REFERENCES users(id),
    ADD COLUMN suspended_at TIMESTAMP,
    ADD COLUMN admin_notes TEXT;

CREATE INDEX idx_raffles_suspended_by ON raffles(suspended_by) WHERE suspended_by IS NOT NULL;

COMMENT ON COLUMN raffles.suspension_reason IS 'Razón de suspensión (visible para organizador)';
COMMENT ON COLUMN raffles.suspended_by IS 'Admin que suspendió la rifa';
COMMENT ON COLUMN raffles.admin_notes IS 'Notas privadas de admin (no visibles para organizador)';
