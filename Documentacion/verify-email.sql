-- Script para verificar emails manualmente en desarrollo
-- SOLO USAR EN DESARROLLO/TESTING

-- Verificar un usuario específico por email
UPDATE users
SET kyc_level = 'email_verified'
WHERE email = 'AlonsoAlpizar@gmail.com';

-- Ver estado de verificación de todos los usuarios
SELECT id, email, kyc_level, created_at
FROM users
ORDER BY created_at DESC;

-- Verificar todos los usuarios (SOLO EN DESARROLLO)
-- UPDATE users SET kyc_level = 'email_verified';
