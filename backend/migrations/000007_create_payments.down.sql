DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_reservation_id;
DROP INDEX IF EXISTS idx_payments_stripe_intent_id;
DROP INDEX IF EXISTS idx_payments_raffle_id;
DROP INDEX IF EXISTS idx_payments_user_id;
DROP TABLE IF EXISTS payments;
