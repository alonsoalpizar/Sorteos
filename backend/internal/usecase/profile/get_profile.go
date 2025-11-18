package profile

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
)

// GetProfileUseCase obtiene el perfil completo del usuario
type GetProfileUseCase struct {
	userRepo          domain.UserRepository
	kycDocumentRepo   domain.KYCDocumentRepository
	walletRepo        domain.WalletRepository
}

// NewGetProfileUseCase crea una nueva instancia del caso de uso
func NewGetProfileUseCase(
	userRepo domain.UserRepository,
	kycDocumentRepo domain.KYCDocumentRepository,
	walletRepo domain.WalletRepository,
) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepo:        userRepo,
		kycDocumentRepo: kycDocumentRepo,
		walletRepo:      walletRepo,
	}
}

// ProfileResponse respuesta del perfil completo
type ProfileResponse struct {
	User         *domain.User           `json:"user"`
	KYCDocuments []*domain.KYCDocument  `json:"kyc_documents"`
	Wallet       *domain.Wallet         `json:"wallet"`
	CanWithdraw  bool                   `json:"can_withdraw"`
}

// Execute ejecuta el caso de uso
func (uc *GetProfileUseCase) Execute(ctx context.Context, userID int64) (*ProfileResponse, error) {
	// Obtener usuario
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Obtener documentos KYC
	kycDocs, err := uc.kycDocumentRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get KYC documents: %w", err)
	}

	// Obtener wallet
	wallet, err := uc.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Verificar si puede retirar
	canWithdraw := user.CanWithdraw()

	return &ProfileResponse{
		User:         user,
		KYCDocuments: kycDocs,
		Wallet:       wallet,
		CanWithdraw:  canWithdraw,
	}, nil
}
