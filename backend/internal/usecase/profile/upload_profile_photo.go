package profile

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
)

// UploadProfilePhotoUseCase maneja la carga de foto de perfil
type UploadProfilePhotoUseCase struct {
	userRepo domain.UserRepository
}

// NewUploadProfilePhotoUseCase crea una nueva instancia del caso de uso
func NewUploadProfilePhotoUseCase(userRepo domain.UserRepository) *UploadProfilePhotoUseCase {
	return &UploadProfilePhotoUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *UploadProfilePhotoUseCase) Execute(ctx context.Context, userID int64, photoURL string) (*domain.User, error) {
	// Obtener usuario
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Actualizar foto de perfil
	user.ProfilePhotoURL = &photoURL

	// Guardar cambios
	if err := uc.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
