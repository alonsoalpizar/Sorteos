-- Migration: 000017_add_earnings_balance (rollback)

-- Eliminar constraint de earnings_balance
ALTER TABLE wallets DROP CONSTRAINT IF EXISTS chk_wallet_earnings_positive;

-- Eliminar columna earnings_balance
ALTER TABLE wallets DROP COLUMN IF EXISTS earnings_balance;

-- Restaurar constraint con nombre original
ALTER TABLE wallets DROP CONSTRAINT IF EXISTS chk_wallet_balance_available_positive;
ALTER TABLE wallets ADD CONSTRAINT chk_wallet_balance_positive CHECK (balance_available >= 0);

-- Renombrar de vuelta a balance
ALTER TABLE wallets RENAME COLUMN balance_available TO balance;

-- Restaurar comentario original
COMMENT ON COLUMN wallets.balance IS 'Saldo disponible para usar';
