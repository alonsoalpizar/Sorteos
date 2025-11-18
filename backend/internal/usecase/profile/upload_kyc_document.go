package profile

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
)

// UploadKYCDocumentUseCase maneja la carga de documentos KYC
type UploadKYCDocumentUseCase struct {
	userRepo        domain.UserRepository
	kycDocumentRepo domain.KYCDocumentRepository
}

// NewUploadKYCDocumentUseCase crea una nueva instancia del caso de uso
func NewUploadKYCDocumentUseCase(
	userRepo domain.UserRepository,
	kycDocumentRepo domain.KYCDocumentRepository,
) *UploadKYCDocumentUseCase {
	return &UploadKYCDocumentUseCase{
		userRepo:        userRepo,
		kycDocumentRepo: kycDocumentRepo,
	}
}

// UploadKYCDocumentRequest datos para subir documento
type UploadKYCDocumentRequest struct {
	DocumentType domain.DocumentType `json:"document_type" binding:"required"`
	FileURL      string              `json:"file_url" binding:"required"`
}

// Execute ejecuta el caso de uso
func (uc *UploadKYCDocumentUseCase) Execute(
	ctx context.Context,
	userID int64,
	docType domain.DocumentType,
	fileURL string,
) (*domain.KYCDocument, error) {
	// Validar tipo de documento
	if err := domain.ValidateDocumentType(docType); err != nil {
		return nil, err
	}

	// Verificar si ya existe un documento de este tipo para el usuario
	existingDoc, err := uc.kycDocumentRepo.FindByUserIDAndType(userID, docType)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing document: %w", err)
	}

	// Si existe, actualizarlo; si no, crear uno nuevo
	if existingDoc != nil {
		// Actualizar documento existente
		existingDoc.FileURL = fileURL
		existingDoc.VerificationStatus = domain.VerificationStatusPending
		existingDoc.VerifiedAt = nil
		existingDoc.VerifiedBy = nil
		existingDoc.RejectedReason = nil

		if err := uc.kycDocumentRepo.Update(existingDoc); err != nil {
			return nil, fmt.Errorf("failed to update document: %w", err)
		}

		return existingDoc, nil
	}

	// Crear nuevo documento
	doc := &domain.KYCDocument{
		UserID:             userID,
		DocumentType:       docType,
		FileURL:            fileURL,
		VerificationStatus: domain.VerificationStatusPending,
	}

	if err := uc.kycDocumentRepo.Create(doc); err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Verificar si el usuario ahora tiene todos los documentos aprobados
	allApproved, err := uc.kycDocumentRepo.AllDocumentsApproved(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check documents status: %w", err)
	}

	// Si todos los documentos están aprobados, actualizar KYC level a full_kyc
	if allApproved {
		if err := uc.userRepo.UpdateKYCLevel(userID, domain.KYCLevelFullKYC); err != nil {
			return nil, fmt.Errorf("failed to update KYC level: %w", err)
		}
	}

	return doc, nil
}
