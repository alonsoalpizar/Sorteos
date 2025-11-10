package auth

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/crypto"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// LoginInput representa los datos de entrada para login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginOutput representa los datos de salida del login
type LoginOutput struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"` // en segundos
}

// LoginUseCase maneja el login de usuarios
type LoginUseCase struct {
	userRepo  domain.UserRepository
	auditRepo domain.AuditLogRepository
	tokenMgr  TokenManager
	logger    *logger.Logger
}

// NewLoginUseCase crea una nueva instancia del use case
func NewLoginUseCase(
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	tokenMgr TokenManager,
	logger *logger.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:  userRepo,
		auditRepo: auditRepo,
		tokenMgr:  tokenMgr,
		logger:    logger,
	}
}

// Execute ejecuta el caso de uso de login
func (uc *LoginUseCase) Execute(ctx context.Context, input *LoginInput, ip, userAgent string) (*LoginOutput, error) {
	// Buscar usuario por email
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		if err == errors.ErrUserNotFound {
			// No revelar que el usuario no existe (seguridad)
			return nil, errors.ErrInvalidCredentials
		}
		uc.logger.Error("Error finding user by email", logger.Error(err))
		return nil, err
	}

	// Verificar que el usuario esté activo
	if !user.IsActive() {
		uc.logger.Warn("Login attempt for inactive user",
			logger.Int64("user_id", user.ID),
			logger.String("status", string(user.Status)),
		)

		// Registrar intento en audit log
		auditLog := domain.NewAuditLog(domain.AuditActionUserLoggedIn).
			WithUser(user.ID).
			WithSeverity(domain.AuditSeverityWarning).
			WithDescription("Intento de login en cuenta inactiva").
			WithRequest(ip, userAgent, "/auth/login", "POST", 403).
			Build()
		_ = uc.auditRepo.Create(auditLog)

		switch user.Status {
		case domain.UserStatusSuspended:
			return nil, errors.New("ACCOUNT_SUSPENDED", "Tu cuenta ha sido suspendida. Contacta soporte.", 403, nil)
		case domain.UserStatusBanned:
			return nil, errors.New("ACCOUNT_BANNED", "Tu cuenta ha sido bloqueada permanentemente.", 403, nil)
		default:
			return nil, errors.ErrForbidden
		}
	}

	// Verificar contraseña
	if err := crypto.ComparePassword(input.Password, user.PasswordHash); err != nil {
		uc.logger.Warn("Invalid password attempt",
			logger.Int64("user_id", user.ID),
			logger.String("email", user.Email),
		)

		// Registrar intento fallido en audit log
		auditLog := domain.NewAuditLog(domain.AuditActionUserLoggedIn).
			WithUser(user.ID).
			WithSeverity(domain.AuditSeverityWarning).
			WithDescription("Intento de login con contraseña incorrecta").
			WithRequest(ip, userAgent, "/auth/login", "POST", 401).
			Build()
		_ = uc.auditRepo.Create(auditLog)

		return nil, errors.ErrInvalidCredentials
	}

	// Generar tokens JWT
	accessToken, refreshToken, err := uc.tokenMgr.GenerateTokenPair(user)
	if err != nil {
		uc.logger.Error("Error generating token pair", logger.Error(err))
		return nil, err
	}

	// Actualizar último login
	if err := uc.userRepo.UpdateLastLogin(user.ID, ip); err != nil {
		uc.logger.Warn("Error updating last login", logger.Error(err))
		// No fallar el login por esto
	}

	// Registrar login exitoso en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionUserLoggedIn).
		WithUser(user.ID).
		WithDescription("Login exitoso").
		WithRequest(ip, userAgent, "/auth/login", "POST", 200).
		Build()

	if err := uc.auditRepo.Create(auditLog); err != nil {
		uc.logger.Warn("Error creating audit log", logger.Error(err))
	}

	uc.logger.Info("User logged in successfully",
		logger.Int64("user_id", user.ID),
		logger.String("email", user.Email),
	)

	return &LoginOutput{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutos en segundos
	}, nil
}
