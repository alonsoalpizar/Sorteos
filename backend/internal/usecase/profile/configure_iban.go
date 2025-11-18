package profile

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
)

// ConfigureIBANUseCase configura el IBAN del usuario para retiros
type ConfigureIBANUseCase struct {
	userRepo domain.UserRepository
}

// NewConfigureIBANUseCase crea una nueva instancia del caso de uso
func NewConfigureIBANUseCase(userRepo domain.UserRepository) *ConfigureIBANUseCase {
	return &ConfigureIBANUseCase{
		userRepo: userRepo,
	}
}

// ConfigureIBANRequest datos para configurar IBAN
type ConfigureIBANRequest struct {
	IBAN string `json:"iban" binding:"required"`
}

// Execute ejecuta el caso de uso
func (uc *ConfigureIBANUseCase) Execute(ctx context.Context, userID int64, iban string) (*domain.User, error) {
	// Obtener usuario
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verificar que tenga al menos cedula_verified
	if !user.HasMinimumKYC(domain.KYCLevelCedulaVerified) {
		return nil, fmt.Errorf("user must have at least cedula_verified KYC level to configure IBAN")
	}

	// Validar formato IBAN
	if err := domain.ValidateIBAN(iban); err != nil {
		return nil, fmt.Errorf("invalid IBAN: %w", err)
	}

	// TODO: Encriptar IBAN antes de guardar
	// Por ahora guardamos en texto plano - implementar encriptación en siguiente iteración
	user.IBAN = &iban

	// Guardar cambios
	if err := uc.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
