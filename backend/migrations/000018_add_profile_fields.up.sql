-- Agregar campos de perfil a la tabla users
ALTER TABLE users
ADD COLUMN IF NOT EXISTS profile_photo_url VARCHAR(255),
ADD COLUMN IF NOT EXISTS date_of_birth DATE,
ADD COLUMN IF NOT EXISTS iban VARCHAR(255);

-- Crear índice en date_of_birth para búsquedas (ej: verificar edad mínima)
CREATE INDEX IF NOT EXISTS idx_users_date_of_birth ON users(date_of_birth);

-- Comentarios para documentación
COMMENT ON COLUMN users.profile_photo_url IS 'URL de la foto de perfil del usuario';
COMMENT ON COLUMN users.date_of_birth IS 'Fecha de nacimiento del usuario';
COMMENT ON COLUMN users.iban IS 'Cuenta IBAN para retiros (encriptado en app layer)';
