package profile

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
)

// UpdateProfileUseCase actualiza la información personal del usuario
type UpdateProfileUseCase struct {
	userRepo domain.UserRepository
}

// NewUpdateProfileUseCase crea una nueva instancia del caso de uso
func NewUpdateProfileUseCase(userRepo domain.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		userRepo: userRepo,
	}
}

// UpdateProfileRequest datos para actualizar perfil
type UpdateProfileRequest struct {
	FirstName    *string   `json:"first_name,omitempty"`
	LastName     *string   `json:"last_name,omitempty"`
	DateOfBirth  *DateOnly `json:"date_of_birth,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	Cedula       *string   `json:"cedula,omitempty"`
	AddressLine1 *string   `json:"address_line1,omitempty"`
	AddressLine2 *string   `json:"address_line2,omitempty"`
	City         *string   `json:"city,omitempty"`
	State        *string   `json:"state,omitempty"`
	PostalCode   *string   `json:"postal_code,omitempty"`
}

// Execute ejecuta el caso de uso
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID int64, req *UpdateProfileRequest) (*domain.User, error) {
	// Obtener usuario actual
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Validar y actualizar campos
	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}

	if req.LastName != nil {
		user.LastName = req.LastName
	}

	if req.DateOfBirth != nil {
		// Validar fecha de nacimiento
		if err := domain.ValidateDateOfBirth(req.DateOfBirth.Time); err != nil {
			return nil, fmt.Errorf("invalid date of birth: %w", err)
		}
		user.DateOfBirth = &req.DateOfBirth.Time
	}

	if req.Phone != nil {
		// Validar formato de teléfono
		if err := domain.ValidatePhone(*req.Phone); err != nil {
			return nil, fmt.Errorf("invalid phone: %w", err)
		}

		// Verificar que el teléfono no esté en uso por otro usuario
		existingUser, err := uc.userRepo.FindByPhone(*req.Phone)
		if err == nil && existingUser != nil && existingUser.ID != userID {
			return nil, fmt.Errorf("phone number already in use")
		}

		user.Phone = req.Phone
		// Al cambiar el teléfono, marcarlo como no verificado
		user.PhoneVerified = false
		user.PhoneVerifiedAt = nil
	}

	if req.Cedula != nil {
		// Verificar que la cédula no esté vacía
		if *req.Cedula != "" {
			// Verificar que la cédula no esté en uso por otro usuario
			existingUser, err := uc.userRepo.FindByCedula(*req.Cedula)
			if err == nil && existingUser != nil && existingUser.ID != userID {
				return nil, fmt.Errorf("cedula already in use")
			}
		}

		user.Cedula = req.Cedula
	}

	if req.AddressLine1 != nil {
		user.AddressLine1 = req.AddressLine1
	}

	if req.AddressLine2 != nil {
		user.AddressLine2 = req.AddressLine2
	}

	if req.City != nil {
		user.City = req.City
	}

	if req.State != nil {
		user.State = req.State
	}

	if req.PostalCode != nil {
		user.PostalCode = req.PostalCode
	}

	// Guardar cambios
	if err := uc.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
