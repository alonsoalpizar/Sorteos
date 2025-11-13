-- Add phase support for double timeout system (selection 10min + checkout 5min)

-- Add ENUM type for reservation phase
CREATE TYPE reservation_phase AS ENUM ('selection', 'checkout', 'completed', 'expired');

-- Add new columns to reservations table
ALTER TABLE reservations
ADD COLUMN phase reservation_phase DEFAULT 'selection',
ADD COLUMN selection_started_at TIMESTAMP DEFAULT NOW(),
ADD COLUMN checkout_started_at TIMESTAMP NULL;

-- Add index for phase queries
CREATE INDEX idx_reservations_phase ON reservations(phase);

-- Add index for active pending reservations in selection phase
CREATE INDEX idx_reservations_active_selection
ON reservations(phase, expires_at)
WHERE status = 'pending' AND phase = 'selection';

-- Add index for checkout phase tracking
CREATE INDEX idx_reservations_checkout
ON reservations(phase, expires_at)
WHERE status = 'pending' AND phase = 'checkout';

-- Update existing records to have selection_started_at = created_at
UPDATE reservations
SET selection_started_at = created_at
WHERE selection_started_at IS NULL;

-- Make selection_started_at NOT NULL after backfilling
ALTER TABLE reservations
ALTER COLUMN selection_started_at SET NOT NULL;

-- Add comment explaining the new fields
COMMENT ON COLUMN reservations.phase IS 'Reservation phase: selection (10min) -> checkout (5min) -> completed/expired';
COMMENT ON COLUMN reservations.selection_started_at IS 'When user started selecting numbers (10 min timer)';
COMMENT ON COLUMN reservations.checkout_started_at IS 'When user clicked "Pay Now" (5 min timer starts)';
