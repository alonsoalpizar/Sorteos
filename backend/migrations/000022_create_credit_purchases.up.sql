-- Migración 000022: Crear tabla credit_purchases para integración con Pagadito
-- Esta tabla registra todas las compras de créditos realizadas a través de Pagadito

-- Tipo ENUM para estados de compra
CREATE TYPE credit_purchase_status AS ENUM (
    'pending',      -- Pago iniciado, esperando redirección a Pagadito
    'processing',   -- Usuario en página de Pagadito
    'completed',    -- Pago completado, créditos acreditados
    'failed',       -- Pago fallido o cancelado
    'expired'       -- Expiró sin completarse (30 min timeout)
);

-- Tabla principal de compras de créditos
CREATE TABLE IF NOT EXISTS credit_purchases (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE RESTRICT,

    -- Montos
    desired_credit DECIMAL(12,2) NOT NULL,        -- Crédito que usuario quiere (ej: ₡5,000)
    charge_amount DECIMAL(12,2) NOT NULL,         -- Monto total cobrado (ej: ₡5,500 con fees)
    currency VARCHAR(3) NOT NULL DEFAULT 'CRC',

    -- Desglose de comisiones (para transparencia)
    fixed_fee DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    processor_fee DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    platform_fee DECIMAL(12,2) NOT NULL DEFAULT 0.00,

    -- Integración Pagadito
    ern VARCHAR(255) NOT NULL UNIQUE,             -- External Reference Number (único por transacción)
    pagadito_token VARCHAR(255),                  -- Token de transacción de Pagadito
    pagadito_reference VARCHAR(255),              -- NAP (Número de Aprobación Pagadito)
    pagadito_status VARCHAR(50),                  -- COMPLETED, REGISTERED, VERIFYING, REVOKED, FAILED

    -- Estado
    status credit_purchase_status NOT NULL DEFAULT 'pending',

    -- Idempotencia
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,

    -- Metadata adicional
    metadata JSONB,
    error_message TEXT,

    -- Relación con wallet_transaction (cuando se complete)
    wallet_transaction_id BIGINT REFERENCES wallet_transactions(id),

    -- Auditoría
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,                -- Timeout de 30 minutos
    completed_at TIMESTAMP,
    failed_at TIMESTAMP,

    -- Constraints
    CONSTRAINT chk_credit_purchase_amounts CHECK (
        desired_credit > 0 AND
        charge_amount >= desired_credit AND
        fixed_fee >= 0 AND
        processor_fee >= 0 AND
        platform_fee >= 0
    ),
    CONSTRAINT chk_credit_purchase_currency CHECK (currency IN ('CRC', 'USD'))
);

-- Índices para búsquedas comunes
CREATE INDEX idx_credit_purchases_user_id ON credit_purchases(user_id, created_at DESC);
CREATE INDEX idx_credit_purchases_status ON credit_purchases(status);
CREATE INDEX idx_credit_purchases_ern ON credit_purchases(ern);
CREATE INDEX idx_credit_purchases_token ON credit_purchases(pagadito_token) WHERE pagadito_token IS NOT NULL;
CREATE INDEX idx_credit_purchases_wallet_tx ON credit_purchases(wallet_transaction_id) WHERE wallet_transaction_id IS NOT NULL;
CREATE INDEX idx_credit_purchases_expires_at ON credit_purchases(expires_at)
    WHERE status IN ('pending', 'processing');

-- Trigger para actualizar updated_at automáticamente
CREATE TRIGGER update_credit_purchases_updated_at
    BEFORE UPDATE ON credit_purchases
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentarios descriptivos
COMMENT ON TABLE credit_purchases IS 'Registro de compras de créditos a través de Pagadito';
COMMENT ON COLUMN credit_purchases.ern IS 'External Reference Number - identificador único para Pagadito';
COMMENT ON COLUMN credit_purchases.pagadito_reference IS 'NAP (Número de Aprobación Pagadito) - devuelto al completar pago';
COMMENT ON COLUMN credit_purchases.desired_credit IS 'Cantidad de créditos que el usuario recibirá en su billetera';
COMMENT ON COLUMN credit_purchases.charge_amount IS 'Monto total cobrado incluyendo todas las comisiones';
COMMENT ON COLUMN credit_purchases.idempotency_key IS 'Clave de idempotencia para prevenir compras duplicadas';
COMMENT ON COLUMN credit_purchases.expires_at IS 'Momento en que expira la compra si no se completa (30 minutos)';
