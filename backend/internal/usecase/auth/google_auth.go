package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// GoogleTokenInfo representa la respuesta de Google al validar el token
type GoogleTokenInfo struct {
	Sub           string `json:"sub"`            // Google unique user ID
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"` // "true" or "false" as string
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GoogleAuthInput representa los datos de entrada para auth con Google
type GoogleAuthInput struct {
	IDToken string `json:"id_token" binding:"required"` // Token de Google
}

// GoogleAuthOutput representa los datos de salida del auth con Google
type GoogleAuthOutput struct {
	User            *domain.User `json:"user,omitempty"`
	AccessToken     string       `json:"access_token,omitempty"`
	RefreshToken    string       `json:"refresh_token,omitempty"`
	TokenType       string       `json:"token_type,omitempty"`
	ExpiresIn       int          `json:"expires_in,omitempty"`
	RequiresLinking bool         `json:"requires_linking"`      // Si el email ya existe, requiere vincular
	ExistingEmail   string       `json:"existing_email,omitempty"` // Email de la cuenta existente
}

// GoogleLinkInput representa los datos para vincular una cuenta existente
type GoogleLinkInput struct {
	IDToken  string `json:"id_token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// GoogleAuthUseCase maneja la autenticación con Google
type GoogleAuthUseCase struct {
	userRepo   domain.UserRepository
	walletRepo domain.WalletRepository
	auditRepo  domain.AuditLogRepository
	tokenMgr   TokenManager
	logger     *logger.Logger
	httpClient *http.Client
}

// NewGoogleAuthUseCase crea una nueva instancia del use case
func NewGoogleAuthUseCase(
	userRepo domain.UserRepository,
	walletRepo domain.WalletRepository,
	auditRepo domain.AuditLogRepository,
	tokenMgr TokenManager,
	logger *logger.Logger,
) *GoogleAuthUseCase {
	return &GoogleAuthUseCase{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		auditRepo:  auditRepo,
		tokenMgr:   tokenMgr,
		logger:     logger,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// verifyGoogleToken verifica el token de Google con la API de Google
// Soporta tanto access_token (implicit flow) como id_token
func (uc *GoogleAuthUseCase) verifyGoogleToken(token string) (*GoogleTokenInfo, error) {
	// Primero intentar con userinfo (access_token del implicit flow)
	tokenInfo, err := uc.verifyWithUserInfo(token)
	if err == nil {
		return tokenInfo, nil
	}

	// Si falla, intentar con tokeninfo (id_token)
	return uc.verifyWithTokenInfo(token)
}

// verifyWithUserInfo verifica usando el endpoint userinfo (para access_token)
func (uc *GoogleAuthUseCase) verifyWithUserInfo(accessToken string) (*GoogleTokenInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		uc.logger.Debug("Error calling Google userinfo", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		uc.logger.Debug("Google userinfo verification failed",
			logger.Int("status", resp.StatusCode),
			logger.String("body", string(body)),
		)
		return nil, fmt.Errorf("userinfo failed with status %d", resp.StatusCode)
	}

	// La respuesta de userinfo tiene un formato ligeramente diferente
	var userInfo struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"` // bool en userinfo, string en tokeninfo
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		uc.logger.Error("Error decoding Google userinfo response", logger.Error(err))
		return nil, err
	}

	// Verificar que el email esté verificado
	if !userInfo.EmailVerified {
		return nil, errors.New("EMAIL_NOT_VERIFIED", "El email no está verificado en Google", 400, nil)
	}

	// Convertir al formato estándar
	emailVerifiedStr := "false"
	if userInfo.EmailVerified {
		emailVerifiedStr = "true"
	}

	return &GoogleTokenInfo{
		Sub:           userInfo.Sub,
		Email:         userInfo.Email,
		EmailVerified: emailVerifiedStr,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
	}, nil
}

// verifyWithTokenInfo verifica usando el endpoint tokeninfo (para id_token)
func (uc *GoogleAuthUseCase) verifyWithTokenInfo(idToken string) (*GoogleTokenInfo, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	resp, err := uc.httpClient.Get(url)
	if err != nil {
		uc.logger.Error("Error calling Google tokeninfo", logger.Error(err))
		return nil, errors.New("GOOGLE_VERIFY_FAILED", "Error verificando token de Google", 500, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		uc.logger.Warn("Google token verification failed",
			logger.Int("status", resp.StatusCode),
			logger.String("body", string(body)),
		)
		return nil, errors.New("INVALID_GOOGLE_TOKEN", "Token de Google inválido o expirado", 401, nil)
	}

	var tokenInfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		uc.logger.Error("Error decoding Google response", logger.Error(err))
		return nil, errors.New("GOOGLE_DECODE_FAILED", "Error procesando respuesta de Google", 500, err)
	}

	// Verificar que el email esté verificado en Google
	if tokenInfo.EmailVerified != "true" {
		return nil, errors.New("EMAIL_NOT_VERIFIED", "El email no está verificado en Google", 400, nil)
	}

	return &tokenInfo, nil
}

// Execute ejecuta el caso de uso de autenticación con Google
func (uc *GoogleAuthUseCase) Execute(ctx context.Context, input *GoogleAuthInput, ip, userAgent string) (*GoogleAuthOutput, error) {
	// 1. Verificar el token de Google
	tokenInfo, err := uc.verifyGoogleToken(input.IDToken)
	if err != nil {
		return nil, err
	}

	uc.logger.Info("Google token verified",
		logger.String("google_id", tokenInfo.Sub),
		logger.String("email", tokenInfo.Email),
	)

	// 2. Buscar si ya existe un usuario con este Google ID
	existingByGoogleID, err := uc.userRepo.FindByGoogleID(tokenInfo.Sub)
	if err == nil && existingByGoogleID != nil {
		// Usuario ya existe con Google vinculado - hacer login directo
		return uc.loginExistingUser(ctx, existingByGoogleID, ip, userAgent)
	}

	// 3. Buscar si existe un usuario con el mismo email
	existingByEmail, err := uc.userRepo.FindByEmail(tokenInfo.Email)
	if err == nil && existingByEmail != nil {
		// Email ya existe pero sin Google vinculado - requiere vincular con contraseña
		uc.logger.Info("Email exists, requires linking",
			logger.String("email", tokenInfo.Email),
		)
		return &GoogleAuthOutput{
			RequiresLinking: true,
			ExistingEmail:   tokenInfo.Email,
		}, nil
	}

	// 4. Usuario nuevo - crear cuenta automáticamente
	return uc.createNewGoogleUser(ctx, tokenInfo, ip, userAgent)
}

// loginExistingUser hace login para un usuario existente con Google vinculado
func (uc *GoogleAuthUseCase) loginExistingUser(ctx context.Context, user *domain.User, ip, userAgent string) (*GoogleAuthOutput, error) {
	// Verificar que el usuario esté activo
	if !user.IsActive() {
		switch user.Status {
		case domain.UserStatusSuspended:
			return nil, errors.New("ACCOUNT_SUSPENDED", "Tu cuenta ha sido suspendida. Contacta soporte.", 403, nil)
		case domain.UserStatusBanned:
			return nil, errors.New("ACCOUNT_BANNED", "Tu cuenta ha sido bloqueada permanentemente.", 403, nil)
		default:
			return nil, errors.ErrForbidden
		}
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
	}

	// Registrar login en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionUserLoggedIn).
		WithUser(user.ID).
		WithDescription("Login exitoso via Google OAuth").
		WithRequest(ip, userAgent, "/auth/google", "POST", 200).
		Build()
	_ = uc.auditRepo.Create(auditLog)

	uc.logger.Info("User logged in via Google",
		logger.Int64("user_id", user.ID),
		logger.String("email", user.Email),
	)

	return &GoogleAuthOutput{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}

// createNewGoogleUser crea un nuevo usuario desde Google OAuth
func (uc *GoogleAuthUseCase) createNewGoogleUser(ctx context.Context, tokenInfo *GoogleTokenInfo, ip, userAgent string) (*GoogleAuthOutput, error) {
	now := time.Now()
	googleID := tokenInfo.Sub

	// Crear usuario con datos de Google
	user := &domain.User{
		UUID:            uuid.New().String(),
		Email:           tokenInfo.Email,
		EmailVerified:   true, // Google ya verificó el email
		EmailVerifiedAt: &now,
		PasswordHash:    "", // No tiene contraseña (solo Google)
		GoogleID:        &googleID,
		AuthProvider:    "google",
		FirstName:       &tokenInfo.GivenName,
		LastName:        &tokenInfo.FamilyName,
		ProfilePhotoURL: &tokenInfo.Picture,
		Role:            domain.UserRoleUser,
		KYCLevel:        domain.KYCLevelEmailVerified, // Auto-verificado por Google
		Status:          domain.UserStatusActive,
		Country:         "CR",
	}

	// Crear usuario en BD
	if err := uc.userRepo.Create(user); err != nil {
		uc.logger.Error("Error creating Google user", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Crear wallet para el usuario
	wallet := &domain.Wallet{
		UUID:     uuid.New().String(),
		UserID:   user.ID,
		Currency: "CRC",
		Status:   domain.WalletStatusActive,
	}
	if err := uc.walletRepo.Create(wallet); err != nil {
		uc.logger.Warn("Error creating wallet for Google user", logger.Error(err))
		// No fallar el registro por esto
	}

	// Generar tokens JWT
	accessToken, refreshToken, err := uc.tokenMgr.GenerateTokenPair(user)
	if err != nil {
		uc.logger.Error("Error generating token pair", logger.Error(err))
		return nil, err
	}

	// Actualizar último login
	_ = uc.userRepo.UpdateLastLogin(user.ID, ip)

	// Registrar en audit log
	auditLog := domain.NewAuditLog(domain.AuditActionUserRegistered).
		WithUser(user.ID).
		WithDescription("Usuario registrado via Google OAuth").
		WithRequest(ip, userAgent, "/auth/google", "POST", 201).
		Build()
	_ = uc.auditRepo.Create(auditLog)

	uc.logger.Info("New user created via Google",
		logger.Int64("user_id", user.ID),
		logger.String("email", user.Email),
	)

	return &GoogleAuthOutput{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}
