-- Rollback: Drop email_notifications table
-- Date: 2025-11-18

-- Drop trigger
DROP TRIGGER IF EXISTS update_email_notifications_updated_at ON email_notifications;

-- Drop indexes
DROP INDEX IF EXISTS idx_email_notifications_admin_type_created;
DROP INDEX IF EXISTS idx_email_notifications_type_status;
DROP INDEX IF EXISTS idx_email_notifications_scheduled_at;
DROP INDEX IF EXISTS idx_email_notifications_sent_at;
DROP INDEX IF EXISTS idx_email_notifications_created_at;
DROP INDEX IF EXISTS idx_email_notifications_priority;
DROP INDEX IF EXISTS idx_email_notifications_status;
DROP INDEX IF EXISTS idx_email_notifications_type;
DROP INDEX IF EXISTS idx_email_notifications_admin_id;

-- Drop table
DROP TABLE IF EXISTS email_notifications;

-- Drop enum types
DROP TYPE IF EXISTS notification_priority;
DROP TYPE IF EXISTS notification_status;
DROP TYPE IF EXISTS notification_type;
