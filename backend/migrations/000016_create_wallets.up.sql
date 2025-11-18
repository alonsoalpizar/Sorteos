-- Migration: 000016_create_wallets
-- Purpose: Sistema de billetera/monedero para usuarios

-- Tipos ENUM
CREATE TYPE wallet_status AS ENUM (
    'active', 'frozen', 'closed'
);

CREATE TYPE transaction_type AS ENUM (
    'deposit',           -- Compra de créditos vía procesador de pagos local
    'withdrawal',        -- Retiro a cuenta bancaria (SINPE/transferencia)
    'purchase',          -- Pago de sorteo con saldo
    'refund',            -- Devolución de compra
    'prize_claim',       -- Premio ganado (organizador)
    'settlement_payout', -- Pago de liquidación a organizador
    'adjustment'         -- Ajuste manual (admin)
);

CREATE TYPE transaction_status AS ENUM (
    'pending', 'completed', 'failed', 'reversed'
);

-- Tabla de billeteras (1 por usuario)
CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE DEFAULT uuid_generate_v4(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    -- Saldos
    balance DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    pending_balance DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'CRC',

    -- Estado
    status wallet_status NOT NULL DEFAULT 'active',

    -- Auditoría
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_wallet_balance_positive CHECK (balance >= 0),
    CONSTRAINT chk_wallet_pending_positive CHECK (pending_balance >= 0),
    CONSTRAINT chk_wallet_currency_format CHECK (LENGTH(currency) = 3)
);

-- Índices para wallets
CREATE UNIQUE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_status ON wallets(status);
CREATE INDEX idx_wallets_uuid ON wallets(uuid);

-- Trigger para updated_at
CREATE TRIGGER update_wallets_updated_at
    BEFORE UPDATE ON wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Tabla de transacciones de billetera
CREATE TABLE wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE DEFAULT uuid_generate_v4(),
    wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    -- Detalles de la transacción
    type transaction_type NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',

    -- Snapshots de saldo (para auditoría)
    balance_before DECIMAL(12,2) NOT NULL,
    balance_after DECIMAL(12,2) NOT NULL,

    -- Referencias externas (polimórfico)
    reference_type VARCHAR(50),
    reference_id BIGINT,

    -- Idempotencia
    idempotency_key VARCHAR(255) NOT NULL,

    -- Metadata adicional (JSONB)
    metadata JSONB,

    -- Notas (para ajustes manuales)
    notes TEXT,

    -- Auditoría
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    failed_at TIMESTAMP,
    reversed_at TIMESTAMP,

    -- Constraints
    CONSTRAINT chk_transaction_amount_positive CHECK (amount > 0),
    CONSTRAINT chk_transaction_balance_before_positive CHECK (balance_before >= 0),
    CONSTRAINT chk_transaction_balance_after_positive CHECK (balance_after >= 0)
);

-- Índices para wallet_transactions
CREATE INDEX idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id, created_at DESC);
CREATE INDEX idx_wallet_transactions_user_id ON wallet_transactions(user_id, created_at DESC);
CREATE UNIQUE INDEX idx_wallet_transactions_idempotency_key ON wallet_transactions(idempotency_key);
CREATE INDEX idx_wallet_transactions_status ON wallet_transactions(status);
CREATE INDEX idx_wallet_transactions_type ON wallet_transactions(type);
CREATE INDEX idx_wallet_transactions_reference ON wallet_transactions(reference_type, reference_id);
CREATE INDEX idx_wallet_transactions_created_at ON wallet_transactions(created_at DESC);

-- Comentarios
COMMENT ON TABLE wallets IS 'Billeteras/monederos de usuarios para almacenar créditos pre-comprados';
COMMENT ON COLUMN wallets.balance IS 'Saldo disponible para usar';
COMMENT ON COLUMN wallets.pending_balance IS 'Saldo pendiente de confirmación (ej: pagos procesando)';

COMMENT ON TABLE wallet_transactions IS 'Historial de todas las transacciones de billetera (audit trail completo)';
COMMENT ON COLUMN wallet_transactions.balance_before IS 'Snapshot del saldo antes de la transacción';
COMMENT ON COLUMN wallet_transactions.balance_after IS 'Snapshot del saldo después de la transacción';
COMMENT ON COLUMN wallet_transactions.idempotency_key IS 'Clave única para prevenir transacciones duplicadas';
COMMENT ON COLUMN wallet_transactions.reference_type IS 'Tipo de entidad relacionada (payment, settlement, raffle, admin)';
COMMENT ON COLUMN wallet_transactions.reference_id IS 'ID de la entidad relacionada';
