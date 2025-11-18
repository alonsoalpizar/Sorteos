-- Migration: 000010_payment_processors
-- Purpose: Gestionar credenciales de procesadores de pago

CREATE TABLE payment_processors (
    id BIGSERIAL PRIMARY KEY,

    -- Provider Info
    provider VARCHAR(50) NOT NULL, -- 'stripe', 'paypal', 'credix', etc.
    name VARCHAR(255) NOT NULL, -- 'Stripe Production', 'PayPal Sandbox'

    -- Status
    is_active BOOLEAN DEFAULT true,
    is_sandbox BOOLEAN DEFAULT false,

    -- Credentials (store encrypted or reference from env vars)
    client_id TEXT, -- PayPal Client ID, Stripe Publishable Key
    secret_key TEXT, -- Secret key (encrypted in app layer)
    webhook_secret TEXT, -- Webhook verification secret (encrypted in app layer)

    -- Configuration
    currency VARCHAR(3) DEFAULT 'CRC', -- ISO 4217
    config JSONB, -- Additional provider-specific config

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_payment_processors_provider ON payment_processors(provider);
CREATE INDEX idx_payment_processors_active ON payment_processors(is_active);

-- Trigger for updated_at
CREATE TRIGGER update_payment_processors_updated_at
    BEFORE UPDATE ON payment_processors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert existing Stripe configuration
INSERT INTO payment_processors (
    provider,
    name,
    is_active,
    is_sandbox,
    currency,
    config
) VALUES (
    'stripe',
    'Stripe Production',
    true,
    false,
    'CRC',
    '{}'::jsonb
);

-- Insert existing PayPal configuration
INSERT INTO payment_processors (
    provider,
    name,
    is_active,
    is_sandbox,
    currency,
    config
) VALUES (
    'paypal',
    'PayPal Production',
    false,
    false,
    'USD',
    '{}'::jsonb
);

COMMENT ON TABLE payment_processors IS 'Configuración de procesadores de pago';
COMMENT ON COLUMN payment_processors.secret_key IS 'Secret key - encrypt in application layer';
COMMENT ON COLUMN payment_processors.webhook_secret IS 'Webhook secret - encrypt in application layer';
