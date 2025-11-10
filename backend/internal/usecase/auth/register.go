package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/crypto"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// RegisterInput representa los datos de entrada para registro
type RegisterInput struct {
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required,min=12"`
	Phone           *string `json:"phone,omitempty"`
	FirstName       *string `json:"first_name,omitempty"`
	LastName        *string `json:"last_name,omitempty"`
	AcceptedTerms   bool    `json:"accepted_terms" binding:"required"`
	AcceptedPrivacy bool    `json:"accepted_privacy" binding:"required"`
}

// RegisterOutput representa los datos de salida del registro
type RegisterOutput struct {
	User                 *domain.User `json:"user"`
	VerificationCodeSent bool         `json:"verification_code_sent"`
	Message              string       `json:"message"`
}

// RegisterUseCase maneja el registro de usuarios
type RegisterUseCase struct {
	userRepo    domain.UserRepository
	consentRepo domain.UserConsentRepository
	auditRepo   domain.AuditLogRepository
	tokenMgr    TokenManager
	notifier    Notifier
	logger      *logger.Logger
}

// TokenManager interface para gestión de tokens
type TokenManager interface {
	StoreVerificationCode(userID int64, codeType, code string, ttl time.Duration) error
	GenerateTokenPair(user *domain.User) (accessToken, refreshToken string, err error)
}

// Notifier interface para envío de notificaciones
type Notifier interface {
	SendVerificationEmail(email, code string) error
}

// NewRegisterUseCase crea una nueva instancia del use case
func NewRegisterUseCase(
	userRepo domain.UserRepository,
	consentRepo domain.UserConsentRepository,
	auditRepo domain.AuditLogRepository,
	tokenMgr TokenManager,
	notifier Notifier,
	logger *logger.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:    userRepo,
		consentRepo: consentRepo,
		auditRepo:   auditRepo,
		tokenMgr:    tokenMgr,
		notifier:    notifier,
		logger:      logger,
	}
}

// Execute ejecuta el caso de uso de registro
func (uc *RegisterUseCase) Execute(ctx context.Context, input *RegisterInput, ip, userAgent string) (*RegisterOutput, error) {
	// Validar email
	if err := domain.ValidateEmail(input.Email); err != nil {
		return nil, errors.WrapWithMessage(errors.ErrValidationFailed, err.Error(), err)
	}

	// Validar contraseña
	if err := domain.ValidatePassword(input.Password); err != nil {
		return nil, errors.WrapWithMessage(errors.ErrValidationFailed, err.Error(), err)
	}

	// Validar teléfono si se proporciona
	if input.Phone != nil {
		if err := domain.ValidatePhone(*input.Phone); err != nil {
			return nil, errors.WrapWithMessage(errors.ErrValidationFailed, err.Error(), err)
		}
	}

	// Verificar que el email no esté registrado
	existingUser, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil && err != errors.ErrUserNotFound {
		uc.logger.Error("Error checking existing email", logger.Error(err))
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	// Verificar que el teléfono no esté registrado
	if input.Phone != nil {
		existingUser, err := uc.userRepo.FindByPhone(*input.Phone)
		if err != nil && err != errors.ErrUserNotFound {
			uc.logger.Error("Error checking existing phone", logger.Error(err))
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.ErrPhoneAlreadyExists
		}
	}

	// Hashear contraseña
	passwordHash, err := crypto.HashPassword(input.Password)
	if err != nil {
		uc.logger.Error("Error hashing password", logger.Error(err))
		return nil, err
	}

	// Crear usuario
	user := &domain.User{
		UUID:         uuid.New().String(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Phone:        input.Phone,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Role:         domain.UserRoleUser,
		KYCLevel:     domain.KYCLevelNone,
		Status:       domain.UserStatusActive,
		Country:      "CR",
	}

	if err := uc.userRepo.Create(user); err != nil {
		uc.logger.Error("Error creating user", logger.Error(err))
		return nil, err
	}

	// Crear consentimientos GDPR
	if input.AcceptedTerms {
		termsConsent := &domain.UserConsent{
			UserID:         user.ID,
			ConsentType:    domain.ConsentTypeTermsOfService,
			ConsentVersion: "1.0",
			IPAddress:      &ip,
			UserAgent:      &userAgent,
		}
		termsConsent.Grant(ip, userAgent)
		if err := uc.consentRepo.Create(termsConsent); err != nil {
			uc.logger.Warn("Error creating terms consent", logger.Error(err))
		}
	}

	if input.AcceptedPrivacy {
		privacyConsent := &domain.UserConsent{
			UserID:         user.ID,
			ConsentType:    domain.ConsentTypePrivacyPolicy,
			ConsentVersion: "1.0",
			IPAddress:      &ip,
			UserAgent:      &userAgent,
		}
		privacyConsent.Grant(ip, userAgent)
		if err := uc.consentRepo.Create(privacyConsent); err != nil {
			uc.logger.Warn("Error creating privacy consent", logger.Error(err))
		}
	}

	// Generar código de verificación
	code, err := crypto.GenerateVerificationCode()
	if err != nil {
		uc.logger.Error("Error generating verification code", logger.Error(err))
		return nil, errors.Wrap(errors.ErrInternalServer, err)
	}

	// Guardar código en Redis (expira en 15 minutos)
	if err := uc.tokenMgr.StoreVerificationCode(user.ID, "email", code, 15*time.Minute); err != nil {
		uc.logger.Error("Error storing verification code", logger.Error(err))
		return nil, err
	}

	// Enviar email de verificación
	verificationSent := false
	if err := uc.notifier.SendVerificationEmail(user.Email, code); err != nil {
		uc.logger.Warn("Error sending verification email", logger.Error(err))
		// No fallar el registro si el email no se envía
	} else {
		verificationSent = true
	}

	// Registrar en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionUserRegistered).
		WithUser(user.ID).
		WithDescription("Usuario registrado exitosamente").
		WithRequest(ip, userAgent, "/auth/register", "POST", 201).
		WithMetadata(map[string]interface{}{
			"email": user.Email,
			"phone": user.Phone,
		}).
		Build()

	if err := uc.auditRepo.Create(auditLog); err != nil {
		uc.logger.Warn("Error creating audit log", logger.Error(err))
	}

	uc.logger.Info("User registered successfully",
		logger.Int64("user_id", user.ID),
		logger.String("email", user.Email),
	)

	message := "Registro exitoso. Por favor verifica tu email."
	if !verificationSent {
		message = "Registro exitoso. No pudimos enviar el email de verificación, por favor contacta soporte."
	}

	return &RegisterOutput{
		User:                 user,
		VerificationCodeSent: verificationSent,
		Message:              message,
	}, nil
}
