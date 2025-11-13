-- Rollback double timeout enhancement

-- Drop indexes
DROP INDEX IF EXISTS idx_reservations_checkout;
DROP INDEX IF EXISTS idx_reservations_active_selection;
DROP INDEX IF EXISTS idx_reservations_phase;

-- Drop columns
ALTER TABLE reservations
DROP COLUMN IF EXISTS checkout_started_at,
DROP COLUMN IF EXISTS selection_started_at,
DROP COLUMN IF EXISTS phase;

-- Drop ENUM type
DROP TYPE IF EXISTS reservation_phase;
