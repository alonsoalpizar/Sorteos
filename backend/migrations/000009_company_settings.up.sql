-- Migration: 000009_company_settings
-- Purpose: Almacenar datos maestros de la empresa Sorteos.club

CREATE TABLE company_settings (
    id BIGSERIAL PRIMARY KEY,

    -- Company Info
    company_name VARCHAR(255) NOT NULL DEFAULT 'Sorteos.club',
    tax_id VARCHAR(50), -- RUC o Tax ID

    -- Address
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'CR', -- ISO 3166-1 alpha-2

    -- Contact
    phone VARCHAR(20),
    email VARCHAR(255),
    website VARCHAR(255) DEFAULT 'https://sorteos.club',
    support_email VARCHAR(255) DEFAULT 'soporte@sorteos.club',

    -- Branding
    logo_url VARCHAR(500),

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Trigger for updated_at
CREATE TRIGGER update_company_settings_updated_at
    BEFORE UPDATE ON company_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default data
INSERT INTO company_settings (
    company_name,
    country,
    email,
    support_email,
    website
) VALUES (
    'Sorteos.club',
    'CR',
    'info@sorteos.club',
    'soporte@sorteos.club',
    'https://sorteos.club'
);

-- Only one row allowed (singleton pattern)
CREATE UNIQUE INDEX idx_company_settings_singleton ON company_settings ((id IS NOT NULL));

COMMENT ON TABLE company_settings IS 'Configuracion global de la empresa (singleton)';
