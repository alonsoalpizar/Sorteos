-- Rollback de migración 000022

-- Eliminar trigger
DROP TRIGGER IF EXISTS update_credit_purchases_updated_at ON credit_purchases;

-- Eliminar índices
DROP INDEX IF EXISTS idx_credit_purchases_expires_at;
DROP INDEX IF EXISTS idx_credit_purchases_wallet_tx;
DROP INDEX IF EXISTS idx_credit_purchases_token;
DROP INDEX IF EXISTS idx_credit_purchases_ern;
DROP INDEX IF EXISTS idx_credit_purchases_status;
DROP INDEX IF EXISTS idx_credit_purchases_user_id;

-- Eliminar tabla
DROP TABLE IF EXISTS credit_purchases;

-- Eliminar tipo ENUM
DROP TYPE IF EXISTS credit_purchase_status;
