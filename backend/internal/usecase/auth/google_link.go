package auth

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/crypto"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GoogleLinkUseCase maneja la vinculación de cuentas Google con cuentas existentes
type GoogleLinkUseCase struct {
	userRepo       domain.UserRepository
	auditRepo      domain.AuditLogRepository
	tokenMgr       TokenManager
	googleAuthUC   *GoogleAuthUseCase
	logger         *logger.Logger
}

// NewGoogleLinkUseCase crea una nueva instancia del use case
func NewGoogleLinkUseCase(
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	tokenMgr TokenManager,
	googleAuthUC *GoogleAuthUseCase,
	logger *logger.Logger,
) *GoogleLinkUseCase {
	return &GoogleLinkUseCase{
		userRepo:     userRepo,
		auditRepo:    auditRepo,
		tokenMgr:     tokenMgr,
		googleAuthUC: googleAuthUC,
		logger:       logger,
	}
}

// Execute ejecuta el caso de uso de vincular cuenta Google
func (uc *GoogleLinkUseCase) Execute(ctx context.Context, input *GoogleLinkInput, ip, userAgent string) (*GoogleAuthOutput, error) {
	// 1. Verificar el token de Google
	tokenInfo, err := uc.googleAuthUC.verifyGoogleToken(input.IDToken)
	if err != nil {
		return nil, err
	}

	uc.logger.Info("Attempting to link Google account",
		logger.String("google_id", tokenInfo.Sub),
		logger.String("email", tokenInfo.Email),
	)

	// 2. Buscar usuario por email
	user, err := uc.userRepo.FindByEmail(tokenInfo.Email)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return nil, errors.New("USER_NOT_FOUND", "No existe una cuenta con este correo electrónico", 404, nil)
		}
		uc.logger.Error("Error finding user by email", logger.Error(err))
		return nil, err
	}

	// 3. Verificar que el usuario no tenga ya Google vinculado
	if user.GoogleID != nil && *user.GoogleID != "" {
		return nil, errors.New("GOOGLE_ALREADY_LINKED", "Esta cuenta ya tiene Google vinculado", 400, nil)
	}

	// 4. Verificar la contraseña
	if err := crypto.ComparePassword(input.Password, user.PasswordHash); err != nil {
		uc.logger.Warn("Invalid password for Google link attempt",
			logger.Int64("user_id", user.ID),
			logger.String("email", user.Email),
		)

		// Registrar intento fallido en audit log
		auditLog := domain.NewAuditLog(domain.AuditActionUserLoggedIn).
			WithUser(user.ID).
			WithSeverity(domain.AuditSeverityWarning).
			WithDescription("Intento de vincular Google con contraseña incorrecta").
			WithRequest(ip, userAgent, "/auth/google/link", "POST", 401).
			Build()
		_ = uc.auditRepo.Create(auditLog)

		return nil, errors.ErrInvalidCredentials
	}

	// 5. Vincular la cuenta de Google
	if err := uc.userRepo.LinkGoogleAccount(user.ID, tokenInfo.Sub); err != nil {
		uc.logger.Error("Error linking Google account", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Actualizar el objeto user con los nuevos datos
	googleID := tokenInfo.Sub
	user.GoogleID = &googleID
	user.AuthProvider = "google"

	// Si el usuario no tenía foto de perfil, usar la de Google
	if user.ProfilePhotoURL == nil && tokenInfo.Picture != "" {
		user.ProfilePhotoURL = &tokenInfo.Picture
		_ = uc.userRepo.Update(user)
	}

	// Si el usuario no tenía email verificado, marcarlo como verificado
	if !user.EmailVerified {
		_ = uc.userRepo.VerifyEmail(user.ID)
		user.EmailVerified = true
		user.KYCLevel = domain.KYCLevelEmailVerified
	}

	// 6. Generar tokens JWT
	accessToken, refreshToken, err := uc.tokenMgr.GenerateTokenPair(user)
	if err != nil {
		uc.logger.Error("Error generating token pair", logger.Error(err))
		return nil, err
	}

	// Actualizar último login
	_ = uc.userRepo.UpdateLastLogin(user.ID, ip)

	// Registrar en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionUserLoggedIn).
		WithUser(user.ID).
		WithDescription("Cuenta Google vinculada exitosamente").
		WithRequest(ip, userAgent, "/auth/google/link", "POST", 200).
		Build()
	_ = uc.auditRepo.Create(auditLog)

	uc.logger.Info("Google account linked successfully",
		logger.Int64("user_id", user.ID),
		logger.String("email", user.Email),
		logger.String("google_id", tokenInfo.Sub),
	)

	return &GoogleAuthOutput{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}
