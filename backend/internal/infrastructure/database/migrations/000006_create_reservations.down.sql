DROP TRIGGER IF EXISTS update_reservations_updated_at ON reservations;
DROP INDEX IF EXISTS idx_reservations_status_expires;
DROP INDEX IF EXISTS idx_reservations_session_id;
DROP INDEX IF EXISTS idx_reservations_expires_at;
DROP INDEX IF EXISTS idx_reservations_raffle_id;
DROP INDEX IF EXISTS idx_reservations_user_id;
DROP TABLE IF EXISTS reservations;
