-- Migration: Create email_notifications table
-- Purpose: Store notification history for admin panel
-- Date: 2025-11-18

-- Create enum types for notifications
CREATE TYPE notification_type AS ENUM ('email', 'sms', 'push', 'announcement');
CREATE TYPE notification_status AS ENUM ('queued', 'scheduled', 'sent', 'failed');
CREATE TYPE notification_priority AS ENUM ('low', 'normal', 'high', 'critical');

-- Create email_notifications table
CREATE TABLE email_notifications (
    id BIGSERIAL PRIMARY KEY,

    -- Admin quien envió la notificación
    admin_id BIGINT NOT NULL,

    -- Tipo de notificación
    type notification_type NOT NULL DEFAULT 'email',

    -- Destinatarios (JSON array de objetos {email, name})
    recipients JSONB NOT NULL,

    -- Contenido
    subject TEXT,
    body TEXT NOT NULL,

    -- Template (opcional)
    template_id BIGINT,
    variables JSONB,

    -- Metadatos
    priority notification_priority NOT NULL DEFAULT 'normal',
    status notification_status NOT NULL DEFAULT 'queued',

    -- Timestamps de procesamiento
    sent_at TIMESTAMP,
    scheduled_at TIMESTAMP,

    -- Información del proveedor (SendGrid, Mailgun, etc.)
    provider_id TEXT,
    provider_status TEXT,

    -- Errores
    error TEXT,

    -- Metadatos adicionales (para extensibilidad)
    metadata JSONB,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Foreign keys
    CONSTRAINT fk_admin_id FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes para optimizar queries comunes
CREATE INDEX idx_email_notifications_admin_id ON email_notifications(admin_id);
CREATE INDEX idx_email_notifications_type ON email_notifications(type);
CREATE INDEX idx_email_notifications_status ON email_notifications(status);
CREATE INDEX idx_email_notifications_priority ON email_notifications(priority);
CREATE INDEX idx_email_notifications_created_at ON email_notifications(created_at DESC);
CREATE INDEX idx_email_notifications_sent_at ON email_notifications(sent_at DESC) WHERE sent_at IS NOT NULL;
CREATE INDEX idx_email_notifications_scheduled_at ON email_notifications(scheduled_at) WHERE scheduled_at IS NOT NULL;

-- Index compuesto para filtros comunes
CREATE INDEX idx_email_notifications_type_status ON email_notifications(type, status);
CREATE INDEX idx_email_notifications_admin_type_created ON email_notifications(admin_id, type, created_at DESC);

-- Trigger para updated_at
CREATE TRIGGER update_email_notifications_updated_at
    BEFORE UPDATE ON email_notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentarios para documentación
COMMENT ON TABLE email_notifications IS 'Historial de notificaciones enviadas desde el panel admin';
COMMENT ON COLUMN email_notifications.recipients IS 'Array JSON de destinatarios con formato [{"email":"user@example.com","name":"User Name"}]';
COMMENT ON COLUMN email_notifications.variables IS 'Variables para reemplazar en templates, formato JSON {"variable":"value"}';
COMMENT ON COLUMN email_notifications.provider_id IS 'ID de referencia del proveedor de email (SendGrid message ID, etc.)';
COMMENT ON COLUMN email_notifications.metadata IS 'Metadatos adicionales para extensibilidad futura';
