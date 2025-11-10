-- Eliminar trigger
DROP TRIGGER IF EXISTS update_user_consents_updated_at ON user_consents;

-- Eliminar tabla
DROP TABLE IF EXISTS user_consents;

-- Eliminar tipo ENUM
DROP TYPE IF EXISTS consent_type;
