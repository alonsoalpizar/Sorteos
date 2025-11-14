package auth

import (
	"context"
	"fmt"
	"time"

	redisinfra "github.com/sorteos-platform/backend/internal/infrastructure/redis"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// LogoutInput representa la entrada del caso de uso de logout
type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

// LogoutUseCase maneja la lógica de logout
type LogoutUseCase struct {
	blacklistService *redisinfra.TokenBlacklistService
}

// NewLogoutUseCase crea una nueva instancia del caso de uso
func NewLogoutUseCase(blacklistService *redisinfra.TokenBlacklistService) *LogoutUseCase {
	return &LogoutUseCase{
		blacklistService: blacklistService,
	}
}

// Execute ejecuta el logout del usuario
func (uc *LogoutUseCase) Execute(ctx context.Context, input *LogoutInput) error {
	if input.AccessToken == "" && input.RefreshToken == "" {
		return errors.WrapWithMessage(errors.ErrValidationFailed, "at least one token is required", nil)
	}

	// Agregar access token a la blacklist si está presente
	// TTL: tiempo restante hasta la expiración del access token (15 minutos por defecto)
	if input.AccessToken != "" {
		// 15 minutos - tiempo típico de expiración del access token
		accessTTL := 15 * time.Minute
		err := uc.blacklistService.AddToBlacklist(ctx, input.AccessToken, accessTTL)
		if err != nil {
			return fmt.Errorf("error blacklisting access token: %w", err)
		}
	}

	// Agregar refresh token a la blacklist si está presente
	// TTL: tiempo restante hasta la expiración del refresh token (7 días por defecto)
	if input.RefreshToken != "" {
		// 7 días - tiempo típico de expiración del refresh token
		refreshTTL := 7 * 24 * time.Hour
		err := uc.blacklistService.AddToBlacklist(ctx, input.RefreshToken, refreshTTL)
		if err != nil {
			return fmt.Errorf("error blacklisting refresh token: %w", err)
		}
	}

	return nil
}
