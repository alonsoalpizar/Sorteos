-- Migration: 000017_add_earnings_balance
-- Purpose: Separar saldo de recargas (no retirable) de saldo de ganancias (retirable via IBAN)

-- Renombrar balance a balance_available (saldo de recargas)
ALTER TABLE wallets RENAME COLUMN balance TO balance_available;

-- Agregar nuevo campo para ganancias de sorteos (retirable)
ALTER TABLE wallets ADD COLUMN earnings_balance DECIMAL(12,2) NOT NULL DEFAULT 0.00;

-- Constraint para earnings_balance
ALTER TABLE wallets ADD CONSTRAINT chk_wallet_earnings_positive CHECK (earnings_balance >= 0);

-- Actualizar comentarios
COMMENT ON COLUMN wallets.balance_available IS 'Saldo disponible de recargas (NO retirable, solo para comprar boletos)';
COMMENT ON COLUMN wallets.earnings_balance IS 'Saldo de ganancias de sorteos propios (retirable via IBAN)';

-- Actualizar constraint anterior (cambió el nombre de la columna)
ALTER TABLE wallets DROP CONSTRAINT IF EXISTS chk_wallet_balance_positive;
ALTER TABLE wallets ADD CONSTRAINT chk_wallet_balance_available_positive CHECK (balance_available >= 0);
