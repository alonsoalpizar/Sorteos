-- Revertir migración: eliminar tabla kyc_documents
DROP TABLE IF EXISTS kyc_documents;

-- Eliminar enums
DROP TYPE IF EXISTS kyc_verification_status;
DROP TYPE IF EXISTS kyc_document_type;
