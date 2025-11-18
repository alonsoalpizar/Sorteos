-- Migration: 000013_system_parameters
-- Purpose: Parámetros de negocio configurables dinámicamente

CREATE TABLE system_parameters (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    value_type VARCHAR(20) DEFAULT 'string',
    description TEXT,
    category VARCHAR(50),
    is_sensitive BOOLEAN DEFAULT false,
    updated_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_system_parameters_category ON system_parameters(category);
CREATE INDEX idx_system_parameters_updated_by ON system_parameters(updated_by);

CREATE TRIGGER update_system_parameters_updated_at
    BEFORE UPDATE ON system_parameters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

INSERT INTO system_parameters (key, value, value_type, category, description) VALUES
    ('platform_fee_percentage', '10.0', 'float', 'business', 'Comisión de plataforma por defecto (%)'),
    ('max_active_raffles_per_user', '10', 'int', 'business', 'Máximo de rifas activas simultáneas por organizador'),
    ('reservation_ttl_minutes', '5', 'int', 'business', 'Tiempo de expiración de reservas (minutos)'),
    ('min_raffle_numbers', '10', 'int', 'business', 'Mínimo de números en una rifa'),
    ('max_raffle_numbers', '10000', 'int', 'business', 'Máximo de números en una rifa'),
    ('min_price_per_number', '100.0', 'float', 'business', 'Precio mínimo por número (CRC)'),
    ('require_kyc_for_organizers', 'true', 'bool', 'security', 'Requerir email verificado para crear rifas'),
    ('auto_settlement_creation', 'true', 'bool', 'business', 'Crear settlements automáticamente al completar rifa'),
    ('enable_refunds', 'true', 'bool', 'payment', 'Permitir reembolsos de pagos'),
    ('support_email', 'soporte@sorteos.club', 'string', 'email', 'Email de soporte');

COMMENT ON TABLE system_parameters IS 'Parámetros de negocio configurables';
COMMENT ON COLUMN system_parameters.value_type IS 'Tipo de dato: string, int, float, bool, json';
