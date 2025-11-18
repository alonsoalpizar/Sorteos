package db

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"

	"gorm.io/gorm"
)

// KYCDocumentRepository implementa domain.KYCDocumentRepository
type KYCDocumentRepository struct {
	db *gorm.DB
}

// NewKYCDocumentRepository crea un nuevo repositorio de documentos KYC
func NewKYCDocumentRepository(db *gorm.DB) domain.KYCDocumentRepository {
	return &KYCDocumentRepository{db: db}
}

// Create crea un nuevo documento KYC
func (r *KYCDocumentRepository) Create(doc *domain.KYCDocument) error {
	return r.db.WithContext(context.Background()).Create(doc).Error
}

// FindByID busca un documento por ID
func (r *KYCDocumentRepository) FindByID(id int64) (*domain.KYCDocument, error) {
	var doc domain.KYCDocument
	err := r.db.WithContext(context.Background()).
		Where("id = ?", id).
		First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// FindByUserID obtiene todos los documentos de un usuario
func (r *KYCDocumentRepository) FindByUserID(userID int64) ([]*domain.KYCDocument, error) {
	var docs []*domain.KYCDocument
	err := r.db.WithContext(context.Background()).
		Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Find(&docs).Error
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// FindByUserIDAndType busca un documento específico de un usuario
func (r *KYCDocumentRepository) FindByUserIDAndType(userID int64, docType domain.DocumentType) (*domain.KYCDocument, error) {
	var doc domain.KYCDocument
	err := r.db.WithContext(context.Background()).
		Where("user_id = ? AND document_type = ?", userID, docType).
		First(&doc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// Update actualiza un documento existente
func (r *KYCDocumentRepository) Update(doc *domain.KYCDocument) error {
	return r.db.WithContext(context.Background()).Save(doc).Error
}

// UpdateVerificationStatus actualiza el estado de verificación
func (r *KYCDocumentRepository) UpdateVerificationStatus(
	docID int64,
	status domain.VerificationStatus,
	verifiedBy int64,
	reason *string,
) error {
	updates := map[string]interface{}{
		"verification_status": status,
		"verified_by":         verifiedBy,
	}

	if status == domain.VerificationStatusApproved {
		updates["verified_at"] = gorm.Expr("NOW()")
		updates["rejected_reason"] = nil
	} else if status == domain.VerificationStatusRejected {
		updates["verified_at"] = nil
		updates["rejected_reason"] = reason
	}

	return r.db.WithContext(context.Background()).
		Model(&domain.KYCDocument{}).
		Where("id = ?", docID).
		Updates(updates).Error
}

// Delete elimina un documento
func (r *KYCDocumentRepository) Delete(docID int64) error {
	return r.db.WithContext(context.Background()).
		Where("id = ?", docID).
		Delete(&domain.KYCDocument{}).Error
}

// HasAllDocuments verifica si el usuario ha subido todos los documentos requeridos
func (r *KYCDocumentRepository) HasAllDocuments(userID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(context.Background()).
		Model(&domain.KYCDocument{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	// Debe tener exactamente 3 documentos: cedula_front, cedula_back, selfie
	return count == 3, nil
}

// AllDocumentsApproved verifica si todos los documentos están aprobados
func (r *KYCDocumentRepository) AllDocumentsApproved(userID int64) (bool, error) {
	var approvedCount int64
	err := r.db.WithContext(context.Background()).
		Model(&domain.KYCDocument{}).
		Where("user_id = ? AND verification_status = ?", userID, domain.VerificationStatusApproved).
		Count(&approvedCount).Error

	if err != nil {
		return false, err
	}

	// Si tiene exactamente 3 documentos aprobados, tiene full KYC
	return approvedCount == 3, nil
}
