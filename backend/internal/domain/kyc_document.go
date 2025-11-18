package domain

import (
	"fmt"
	"time"
)

// DocumentType representa el tipo de documento KYC
type DocumentType string

const (
	DocumentTypeCedulaFront DocumentType = "cedula_front"
	DocumentTypeCedulaBack  DocumentType = "cedula_back"
	DocumentTypeSelfie      DocumentType = "selfie"
)

// VerificationStatus representa el estado de verificación de un documento
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusApproved VerificationStatus = "approved"
	VerificationStatusRejected VerificationStatus = "rejected"
)

// KYCDocument representa un documento de verificación KYC
type KYCDocument struct {
	ID                 int64              `json:"id" gorm:"primaryKey"`
	UserID             int64              `json:"user_id" gorm:"not null"`
	DocumentType       DocumentType       `json:"document_type" gorm:"type:kyc_document_type;not null"`
	FileURL            string             `json:"file_url" gorm:"not null"`
	VerificationStatus VerificationStatus `json:"verification_status" gorm:"type:kyc_verification_status;default:'pending';not null"`

	// Información de verificación
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	VerifiedBy     *int64     `json:"verified_by,omitempty"` // Admin user ID
	RejectedReason *string    `json:"rejected_reason,omitempty"`

	// Auditoría
	UploadedAt time.Time `json:"uploaded_at" gorm:"not null;default:now()"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// TableName especifica el nombre de la tabla
func (KYCDocument) TableName() string {
	return "kyc_documents"
}

// ValidateDocumentType valida que el tipo de documento sea válido
func ValidateDocumentType(docType DocumentType) error {
	switch docType {
	case DocumentTypeCedulaFront, DocumentTypeCedulaBack, DocumentTypeSelfie:
		return nil
	default:
		return fmt.Errorf("invalid document type: %s", docType)
	}
}

// IsApproved verifica si el documento está aprobado
func (d *KYCDocument) IsApproved() bool {
	return d.VerificationStatus == VerificationStatusApproved
}

// IsPending verifica si el documento está pendiente
func (d *KYCDocument) IsPending() bool {
	return d.VerificationStatus == VerificationStatusPending
}

// IsRejected verifica si el documento está rechazado
func (d *KYCDocument) IsRejected() bool {
	return d.VerificationStatus == VerificationStatusRejected
}

// KYCDocumentRepository define el contrato para el repositorio de documentos KYC
type KYCDocumentRepository interface {
	// Create crea un nuevo documento KYC
	Create(doc *KYCDocument) error

	// FindByID busca un documento por ID
	FindByID(id int64) (*KYCDocument, error)

	// FindByUserID obtiene todos los documentos de un usuario
	FindByUserID(userID int64) ([]*KYCDocument, error)

	// FindByUserIDAndType busca un documento específico de un usuario
	FindByUserIDAndType(userID int64, docType DocumentType) (*KYCDocument, error)

	// Update actualiza un documento existente
	Update(doc *KYCDocument) error

	// UpdateVerificationStatus actualiza el estado de verificación
	UpdateVerificationStatus(docID int64, status VerificationStatus, verifiedBy int64, reason *string) error

	// Delete elimina un documento
	Delete(docID int64) error

	// HasAllDocuments verifica si el usuario ha subido todos los documentos requeridos
	HasAllDocuments(userID int64) (bool, error)

	// AllDocumentsApproved verifica si todos los documentos están aprobados
	AllDocumentsApproved(userID int64) (bool, error)
}
