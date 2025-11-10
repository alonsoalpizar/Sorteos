package auth

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// RefreshTokenInput representa los datos de entrada para refresh token
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenOutput representa los datos de salida del refresh
type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // en segundos
}

// RefreshTokenManager interface extendida para refresh
type RefreshTokenManager interface {
	TokenManager
	ValidateRefreshToken(tokenString string, userID int64) (*Claims, error)
	RefreshTokenPair(refreshToken string, userID int64, user *domain.User) (newAccessToken, newRefreshToken string, err error)
	IsTokenBlacklisted(userID int64) (bool, error)
}

// RefreshTokenUseCase maneja la renovación de tokens
type RefreshTokenUseCase struct {
	userRepo domain.UserRepository
	tokenMgr RefreshTokenManager
	logger   *logger.Logger
}

// NewRefreshTokenUseCase crea una nueva instancia del use case
func NewRefreshTokenUseCase(
	userRepo domain.UserRepository,
	tokenMgr RefreshTokenManager,
	logger *logger.Logger,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo: userRepo,
		tokenMgr: tokenMgr,
		logger:   logger,
	}
}

// Execute ejecuta el caso de uso de refresh token
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error) {
	// Primero intentar parsear el token para obtener el user_id
	// sin validar completamente (para obtener el userID)
	claims, err := uc.tokenMgr.ValidateRefreshToken(input.RefreshToken, 0)
	if err != nil {
		uc.logger.Warn("Invalid refresh token", logger.Error(err))
		return nil, errors.ErrTokenInvalid
	}

	userID := claims.UserID

	// Verificar que el token no esté en blacklist
	blacklisted, err := uc.tokenMgr.IsTokenBlacklisted(userID)
	if err != nil {
		uc.logger.Error("Error checking token blacklist", logger.Error(err))
		return nil, err
	}
	if blacklisted {
		uc.logger.Warn("Attempt to use blacklisted token",
			logger.Int64("user_id", userID),
		)
		return nil, errors.ErrTokenInvalid
	}

	// Buscar usuario
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return nil, errors.ErrTokenInvalid
		}
		uc.logger.Error("Error finding user", logger.Error(err))
		return nil, err
	}

	// Verificar que el usuario esté activo
	if !user.IsActive() {
		uc.logger.Warn("Refresh token attempt for inactive user",
			logger.Int64("user_id", user.ID),
			logger.String("status", string(user.Status)),
		)
		return nil, errors.ErrForbidden
	}

	// Generar nuevo par de tokens
	newAccessToken, newRefreshToken, err := uc.tokenMgr.RefreshTokenPair(input.RefreshToken, userID, user)
	if err != nil {
		uc.logger.Error("Error refreshing token pair",
			logger.Int64("user_id", userID),
			logger.Error(err),
		)
		return nil, err
	}

	uc.logger.Info("Tokens refreshed successfully",
		logger.Int64("user_id", user.ID),
	)

	return &RefreshTokenOutput{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutos en segundos
	}, nil
}
