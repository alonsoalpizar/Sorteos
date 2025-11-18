-- Revertir migración: eliminar campos de perfil
ALTER TABLE users
DROP COLUMN IF EXISTS iban,
DROP COLUMN IF EXISTS date_of_birth,
DROP COLUMN IF EXISTS profile_photo_url;

-- Eliminar índice
DROP INDEX IF EXISTS idx_users_date_of_birth;
