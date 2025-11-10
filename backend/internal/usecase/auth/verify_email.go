package auth

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// VerifyEmailInput representa los datos de entrada para verificación de email
type VerifyEmailInput struct {
	UserID int64  `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required,len=6"`
}

// VerifyEmailOutput representa los datos de salida de la verificación
type VerifyEmailOutput struct {
	Success      bool         `json:"success"`
	Message      string       `json:"message"`
	User         *domain.User `json:"user,omitempty"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
}

// VerificationCodeValidator interface para validar códigos
type VerificationCodeValidator interface {
	ValidateVerificationCode(userID int64, codeType, code string) (bool, error)
	DeleteVerificationCode(userID int64, codeType string) error
}

// VerifyEmailUseCase maneja la verificación de email
type VerifyEmailUseCase struct {
	userRepo  domain.UserRepository
	auditRepo domain.AuditLogRepository
	tokenMgr  interface {
		TokenManager
		VerificationCodeValidator
	}
	logger *logger.Logger
}

// NewVerifyEmailUseCase crea una nueva instancia del use case
func NewVerifyEmailUseCase(
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	tokenMgr interface {
		TokenManager
		VerificationCodeValidator
	},
	logger *logger.Logger,
) *VerifyEmailUseCase {
	return &VerifyEmailUseCase{
		userRepo:  userRepo,
		auditRepo: auditRepo,
		tokenMgr:  tokenMgr,
		logger:    logger,
	}
}

// Execute ejecuta el caso de uso de verificación de email
func (uc *VerifyEmailUseCase) Execute(ctx context.Context, input *VerifyEmailInput, ip, userAgent string) (*VerifyEmailOutput, error) {
	// Buscar usuario
	user, err := uc.userRepo.FindByID(input.UserID)
	if err != nil {
		uc.logger.Error("Error finding user", logger.Error(err))
		return nil, err
	}

	// Verificar que el email no esté ya verificado
	if user.EmailVerified {
		return &VerifyEmailOutput{
			Success: true,
			Message: "Email ya verificado anteriormente",
			User:    user,
		}, nil
	}

	// Validar código de verificación
	valid, err := uc.tokenMgr.ValidateVerificationCode(user.ID, "email", input.Code)
	if err != nil {
		uc.logger.Error("Error validating verification code", logger.Error(err))
		return nil, err
	}

	if !valid {
		uc.logger.Warn("Invalid verification code",
			logger.Int64("user_id", user.ID),
			logger.String("email", user.Email),
		)

		// Registrar intento fallido en audit log
		auditLog := domain.NewAuditLog(domain.AuditActionEmailVerified).
			WithUser(user.ID).
			WithSeverity(domain.AuditSeverityWarning).
			WithDescription("Código de verificación incorrecto").
			WithRequest(ip, userAgent, "/auth/verify-email", "POST", 400).
			Build()
		_ = uc.auditRepo.Create(auditLog)

		return nil, errors.New("INVALID_VERIFICATION_CODE", "Código de verificación incorrecto o expirado", 400, nil)
	}

	// Marcar email como verificado
	if err := uc.userRepo.VerifyEmail(user.ID); err != nil {
		uc.logger.Error("Error verifying email", logger.Error(err))
		return nil, err
	}

	// Eliminar código de verificación de Redis
	if err := uc.tokenMgr.DeleteVerificationCode(user.ID, "email"); err != nil {
		uc.logger.Warn("Error deleting verification code", logger.Error(err))
		// No fallar por esto
	}

	// Actualizar usuario en memoria
	user.EmailVerified = true
	user.KYCLevel = domain.KYCLevelEmailVerified

	// Generar tokens JWT (para login automático después de verificación)
	accessToken, refreshToken, err := uc.tokenMgr.GenerateTokenPair(user)
	if err != nil {
		uc.logger.Error("Error generating token pair", logger.Error(err))
		return nil, err
	}

	// Registrar verificación exitosa en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionEmailVerified).
		WithUser(user.ID).
		WithDescription("Email verificado exitosamente").
		WithRequest(ip, userAgent, "/auth/verify-email", "POST", 200).
		Build()

	if err := uc.auditRepo.Create(auditLog); err != nil {
		uc.logger.Warn("Error creating audit log", logger.Error(err))
	}

	uc.logger.Info("Email verified successfully",
		logger.Int64("user_id", user.ID),
		logger.String("email", user.Email),
	)

	return &VerifyEmailOutput{
		Success:      true,
		Message:      "Email verificado exitosamente",
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
