package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/config"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// TokenManager gestiona tokens JWT usando Redis
type TokenManager struct {
	rdb    *redis.Client
	config *config.JWTConfig
}

// NewTokenManager crea una nueva instancia del token manager
func NewTokenManager(rdb *redis.Client, cfg *config.JWTConfig) *TokenManager {
	return &TokenManager{
		rdb:    rdb,
		config: cfg,
	}
}

// Claims representa los claims del JWT
type Claims struct {
	UserID   int64           `json:"user_id"`
	Email    string          `json:"email"`
	Role     domain.UserRole `json:"role"`
	KYCLevel domain.KYCLevel `json:"kyc_level"`
	jwt.RegisteredClaims
}

// GenerateAccessToken genera un access token
func (tm *TokenManager) GenerateAccessToken(user *domain.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(tm.config.AccessTokenExpiry)

	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		KYCLevel: user.KYCLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    tm.config.Issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.config.Secret))
	if err != nil {
		return "", errors.Wrap(errors.ErrInternalServer, err)
	}

	return tokenString, nil
}

// GenerateRefreshToken genera un refresh token
func (tm *TokenManager) GenerateRefreshToken(user *domain.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(tm.config.RefreshTokenExpiry)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    tm.config.Issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.config.Secret))
	if err != nil {
		return "", errors.Wrap(errors.ErrInternalServer, err)
	}

	// Guardar en Redis con TTL
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", user.ID)
	if err := tm.rdb.Set(ctx, key, tokenString, tm.config.RefreshTokenExpiry).Err(); err != nil {
		return "", errors.Wrap(errors.ErrRedisError, err)
	}

	return tokenString, nil
}

// ValidateAccessToken valida un access token y retorna los claims
func (tm *TokenManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.config.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrTokenInvalid
	}

	return claims, nil
}

// ValidateRefreshToken valida un refresh token
func (tm *TokenManager) ValidateRefreshToken(tokenString string, userID int64) (*Claims, error) {
	// Validar formato del token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.config.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrTokenInvalid
	}

	// Verificar que el token existe en Redis
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)
	storedToken, err := tm.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.ErrTokenInvalid
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrRedisError, err)
	}

	// Verificar que el token coincide
	if storedToken != tokenString {
		return nil, errors.ErrTokenInvalid
	}

	return claims, nil
}

// RevokeRefreshToken revoca un refresh token
func (tm *TokenManager) RevokeRefreshToken(userID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)
	if err := tm.rdb.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(errors.ErrRedisError, err)
	}
	return nil
}

// RevokeAllUserTokens revoca todos los tokens de un usuario
func (tm *TokenManager) RevokeAllUserTokens(userID int64) error {
	ctx := context.Background()

	// Revocar refresh token
	refreshKey := fmt.Sprintf("refresh_token:%d", userID)
	if err := tm.rdb.Del(ctx, refreshKey).Err(); err != nil {
		return errors.Wrap(errors.ErrRedisError, err)
	}

	// Blacklist de access tokens (guardar user_id en blacklist por el tiempo de expiración)
	blacklistKey := fmt.Sprintf("token_blacklist:%d", userID)
	if err := tm.rdb.Set(ctx, blacklistKey, "revoked", tm.config.AccessTokenExpiry).Err(); err != nil {
		return errors.Wrap(errors.ErrRedisError, err)
	}

	return nil
}

// IsTokenBlacklisted verifica si un token está en la blacklist
func (tm *TokenManager) IsTokenBlacklisted(userID int64) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("token_blacklist:%d", userID)

	exists, err := tm.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.Wrap(errors.ErrRedisError, err)
	}

	return exists > 0, nil
}

// StoreVerificationCode almacena un código de verificación en Redis
func (tm *TokenManager) StoreVerificationCode(userID int64, codeType, code string, ttl time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("verification:%s:%d", codeType, userID)

	if err := tm.rdb.Set(ctx, key, code, ttl).Err(); err != nil {
		return errors.Wrap(errors.ErrRedisError, err)
	}

	return nil
}

// ValidateVerificationCode valida un código de verificación
func (tm *TokenManager) ValidateVerificationCode(userID int64, codeType, code string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("verification:%s:%d", codeType, userID)

	storedCode, err := tm.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // Código no existe o expiró
	}
	if err != nil {
		return false, errors.Wrap(errors.ErrRedisError, err)
	}

	return storedCode == code, nil
}

// DeleteVerificationCode elimina un código de verificación
func (tm *TokenManager) DeleteVerificationCode(userID int64, codeType string) error {
	ctx := context.Background()
	key := fmt.Sprintf("verification:%s:%d", codeType, userID)

	if err := tm.rdb.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(errors.ErrRedisError, err)
	}

	return nil
}

// GenerateTokenPair genera un par de access y refresh tokens
func (tm *TokenManager) GenerateTokenPair(user *domain.User) (accessToken, refreshToken string, err error) {
	accessToken, err = tm.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = tm.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshTokenPair genera un nuevo par de tokens usando un refresh token
func (tm *TokenManager) RefreshTokenPair(refreshToken string, userID int64, user *domain.User) (newAccessToken, newRefreshToken string, err error) {
	// Validar refresh token
	_, err = tm.ValidateRefreshToken(refreshToken, userID)
	if err != nil {
		return "", "", err
	}

	// Generar nuevos tokens
	newAccessToken, err = tm.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	// Si está configurado para rotar refresh tokens
	if tm.config.RefreshTokenRotate {
		// Revocar el refresh token anterior
		if err := tm.RevokeRefreshToken(userID); err != nil {
			return "", "", err
		}

		// Generar nuevo refresh token
		newRefreshToken, err = tm.GenerateRefreshToken(user)
		if err != nil {
			return "", "", err
		}
	} else {
		newRefreshToken = refreshToken
	}

	return newAccessToken, newRefreshToken, nil
}
