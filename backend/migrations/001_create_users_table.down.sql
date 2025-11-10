-- Eliminar trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Eliminar función
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Eliminar tabla
DROP TABLE IF EXISTS users;

-- Eliminar tipos ENUM
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS kyc_level;
DROP TYPE IF EXISTS user_role;

-- Eliminar extensión (solo si no la usan otras tablas)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
