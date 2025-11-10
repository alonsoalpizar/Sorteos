-- Enum para tipos de consentimiento
CREATE TYPE consent_type AS ENUM ('terms_of_service', 'privacy_policy', 'marketing_emails', 'marketing_sms', 'data_processing');

-- Tabla de consentimientos GDPR
CREATE TABLE user_consents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Tipo y versión del consentimiento
    consent_type consent_type NOT NULL,
    consent_version VARCHAR(20) NOT NULL, -- e.g., "1.0", "2.1"

    -- Estado del consentimiento
    granted BOOLEAN NOT NULL DEFAULT FALSE,
    granted_at TIMESTAMP,
    revoked_at TIMESTAMP,

    -- Auditoría
    ip_address INET,
    user_agent TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraint: un usuario solo puede tener un consentimiento activo por tipo
    UNIQUE(user_id, consent_type)
);

-- Índices
CREATE INDEX idx_user_consents_user_id ON user_consents(user_id);
CREATE INDEX idx_user_consents_type ON user_consents(consent_type);
CREATE INDEX idx_user_consents_granted ON user_consents(granted) WHERE granted = TRUE;

-- Trigger para updated_at
CREATE TRIGGER update_user_consents_updated_at
    BEFORE UPDATE ON user_consents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentarios
COMMENT ON TABLE user_consents IS 'Registro de consentimientos GDPR por usuario';
COMMENT ON COLUMN user_consents.consent_version IS 'Versión del documento de términos/política aceptado';
COMMENT ON COLUMN user_consents.granted IS 'TRUE si el consentimiento está activo, FALSE si fue revocado';
