-- Drop triggers
DROP TRIGGER IF EXISTS trigger_calculate_raffle_revenue ON raffles;
DROP TRIGGER IF EXISTS trigger_raffles_updated_at ON raffles;

-- Drop functions
DROP FUNCTION IF EXISTS calculate_raffle_revenue();
DROP FUNCTION IF EXISTS update_raffles_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_raffles_active;
DROP INDEX IF EXISTS idx_raffles_created_at;
DROP INDEX IF EXISTS idx_raffles_draw_date;
DROP INDEX IF EXISTS idx_raffles_status;
DROP INDEX IF EXISTS idx_raffles_user_id;

-- Drop table
DROP TABLE IF EXISTS raffles;

-- Drop ENUM types
DROP TYPE IF EXISTS raffle_status;
