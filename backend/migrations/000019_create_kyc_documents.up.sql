-- Crear enum para tipo de documento KYC
CREATE TYPE kyc_document_type AS ENUM (
    'cedula_front',
    'cedula_back',
    'selfie'
);

-- Crear enum para estado de verificación de documento
CREATE TYPE kyc_verification_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);

-- Crear tabla kyc_documents
CREATE TABLE IF NOT EXISTS kyc_documents (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type kyc_document_type NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    verification_status kyc_verification_status NOT NULL DEFAULT 'pending',

    -- Información de verificación
    verified_at TIMESTAMP,
    verified_by INTEGER REFERENCES users(id), -- Admin que verificó
    rejected_reason TEXT,

    -- Auditoría
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Constraint: Un usuario solo puede tener un documento de cada tipo activo
    UNIQUE(user_id, document_type)
);

-- Índices para búsquedas eficientes
CREATE INDEX idx_kyc_documents_user_id ON kyc_documents(user_id);
CREATE INDEX idx_kyc_documents_verification_status ON kyc_documents(verification_status);
CREATE INDEX idx_kyc_documents_user_type ON kyc_documents(user_id, document_type);

-- Comentarios para documentación
COMMENT ON TABLE kyc_documents IS 'Almacena documentos de verificación KYC de usuarios';
COMMENT ON COLUMN kyc_documents.document_type IS 'Tipo de documento: cedula_front, cedula_back, selfie';
COMMENT ON COLUMN kyc_documents.verification_status IS 'Estado de verificación: pending, approved, rejected';
COMMENT ON COLUMN kyc_documents.file_url IS 'URL del archivo subido (almacenado en /uploads/kyc/)';
COMMENT ON COLUMN kyc_documents.verified_by IS 'ID del admin que verificó el documento';
COMMENT ON COLUMN kyc_documents.rejected_reason IS 'Razón del rechazo si verification_status=rejected';
