-- Drop triggers
DROP TRIGGER IF EXISTS trigger_prevent_delete_only_primary ON raffle_images;
DROP TRIGGER IF EXISTS trigger_auto_set_primary_image ON raffle_images;
DROP TRIGGER IF EXISTS trigger_raffle_images_updated_at ON raffle_images;

-- Drop functions
DROP FUNCTION IF EXISTS prevent_delete_only_primary();
DROP FUNCTION IF EXISTS auto_set_primary_image();
DROP FUNCTION IF EXISTS update_raffle_images_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS uq_one_primary_per_raffle;
DROP INDEX IF EXISTS idx_raffle_images_order;
DROP INDEX IF EXISTS idx_raffle_images_primary;
DROP INDEX IF EXISTS idx_raffle_images_raffle_id;

-- Drop table
DROP TABLE IF EXISTS raffle_images;
