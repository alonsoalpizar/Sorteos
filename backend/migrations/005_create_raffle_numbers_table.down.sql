-- Drop function
DROP FUNCTION IF EXISTS release_expired_reservations();

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_raffle_counters ON raffle_numbers;
DROP TRIGGER IF EXISTS trigger_raffle_numbers_updated_at ON raffle_numbers;

-- Drop functions
DROP FUNCTION IF EXISTS update_raffle_counters();
DROP FUNCTION IF EXISTS update_raffle_numbers_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_raffle_numbers_available;
DROP INDEX IF EXISTS idx_raffle_numbers_reserved_until;
DROP INDEX IF EXISTS idx_raffle_numbers_user_id;
DROP INDEX IF EXISTS idx_raffle_numbers_status;
DROP INDEX IF EXISTS idx_raffle_numbers_raffle_id;

-- Drop table
DROP TABLE IF EXISTS raffle_numbers;

-- Drop ENUM type
DROP TYPE IF EXISTS raffle_number_status;
