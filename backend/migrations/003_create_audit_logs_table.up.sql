-- Enum para tipos de acción en auditoría
CREATE TYPE audit_action AS ENUM (
    -- Auth
    'user_registered', 'user_logged_in', 'user_logged_out', 'email_verified', 'phone_verified',
    -- User Management
    'user_updated', 'user_suspended', 'user_banned', 'user_deleted', 'kyc_level_changed',
    -- Raffles
    'raffle_created', 'raffle_published', 'raffle_suspended', 'raffle_completed', 'raffle_deleted',
    -- Reservations
    'numbers_reserved', 'reservation_expired', 'reservation_cancelled',
    -- Payments
    'payment_created', 'payment_confirmed', 'payment_failed', 'payment_refunded',
    -- Settlements
    'settlement_created', 'settlement_approved', 'settlement_paid', 'settlement_rejected',
    -- Admin Actions
    'admin_action_performed', 'system_parameter_changed', 'report_generated'
);

-- Enum para nivel de criticidad
CREATE TYPE audit_severity AS ENUM ('info', 'warning', 'error', 'critical');

-- Tabla de logs de auditoría
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,

    -- Actor (quién realizó la acción)
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    admin_id BIGINT REFERENCES users(id) ON DELETE SET NULL,

    -- Acción
    action audit_action NOT NULL,
    severity audit_severity DEFAULT 'info' NOT NULL,
    description TEXT,

    -- Entidad afectada (polimórfico)
    entity_type VARCHAR(50), -- e.g., 'raffle', 'user', 'payment'
    entity_id BIGINT,

    -- Contexto de la solicitud
    ip_address INET,
    user_agent TEXT,
    endpoint VARCHAR(255),
    http_method VARCHAR(10),
    http_status_code INT,

    -- Datos adicionales (JSON)
    metadata JSONB,

    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para búsquedas frecuentes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_admin_id ON audit_logs(admin_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_severity ON audit_logs(severity);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_ip_address ON audit_logs(ip_address);

-- Índice GIN para búsquedas en metadata JSON
CREATE INDEX idx_audit_logs_metadata ON audit_logs USING GIN(metadata);

-- Particionamiento por fecha (opcional, para mejor performance en grandes volúmenes)
-- Nota: Esto se puede implementar después cuando crezca la tabla
-- ALTER TABLE audit_logs PARTITION BY RANGE (created_at);

-- Comentarios
COMMENT ON TABLE audit_logs IS 'Registro de auditoría de todas las acciones críticas del sistema';
COMMENT ON COLUMN audit_logs.user_id IS 'Usuario que realizó la acción (si aplica)';
COMMENT ON COLUMN audit_logs.admin_id IS 'Admin que realizó la acción (para acciones administrativas)';
COMMENT ON COLUMN audit_logs.entity_type IS 'Tipo de entidad afectada: raffle, user, payment, etc.';
COMMENT ON COLUMN audit_logs.metadata IS 'Datos adicionales en formato JSON (valores anteriores, nuevos, etc.)';
