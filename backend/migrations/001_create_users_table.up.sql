-- Crear extensión para UUIDs si no existe
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum para roles de usuario
CREATE TYPE user_role AS ENUM ('user', 'admin', 'super_admin');

-- Enum para niveles de KYC
CREATE TYPE kyc_level AS ENUM ('none', 'email_verified', 'phone_verified', 'cedula_verified', 'full_kyc');

-- Enum para estados de usuario
CREATE TYPE user_status AS ENUM ('active', 'suspended', 'banned', 'deleted');

-- Tabla de usuarios
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),

    -- Credenciales
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP,
    phone VARCHAR(20),
    phone_verified BOOLEAN DEFAULT FALSE,
    phone_verified_at TIMESTAMP,
    password_hash VARCHAR(255) NOT NULL,

    -- Información personal
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    cedula VARCHAR(20) UNIQUE,

    -- Dirección
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'CR',

    -- Roles y verificación
    role user_role DEFAULT 'user' NOT NULL,
    kyc_level kyc_level DEFAULT 'none' NOT NULL,
    status user_status DEFAULT 'active' NOT NULL,

    -- Límites de usuario
    max_active_raffles INT DEFAULT 10,
    purchase_limit_daily DECIMAL(12,2) DEFAULT 50000.00,

    -- Tokens
    refresh_token TEXT,
    refresh_token_expires_at TIMESTAMP,

    -- Códigos de verificación
    email_verification_code VARCHAR(6),
    email_verification_expires_at TIMESTAMP,
    phone_verification_code VARCHAR(6),
    phone_verification_expires_at TIMESTAMP,
    password_reset_token VARCHAR(64),
    password_reset_expires_at TIMESTAMP,

    -- Auditoría
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip INET,

    -- Soft delete
    deleted_at TIMESTAMP
);

-- Índices para búsquedas frecuentes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_cedula ON users(cedula) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_kyc_level ON users(kyc_level);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentarios
COMMENT ON TABLE users IS 'Tabla principal de usuarios del sistema';
COMMENT ON COLUMN users.kyc_level IS 'Nivel de verificación KYC: none, email_verified, phone_verified, cedula_verified, full_kyc';
COMMENT ON COLUMN users.status IS 'Estado del usuario: active, suspended, banned, deleted';
COMMENT ON COLUMN users.purchase_limit_daily IS 'Límite de compra diaria en colones';
