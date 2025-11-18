package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ResetUserPasswordUseCase caso de uso para resetear contraseña de usuario (admin)
type ResetUserPasswordUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewResetUserPasswordUseCase crea una nueva instancia
func NewResetUserPasswordUseCase(db *gorm.DB, log *logger.Logger) *ResetUserPasswordUseCase {
	return &ResetUserPasswordUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ResetUserPasswordUseCase) Execute(ctx context.Context, userID int64, adminID int64) (string, error) {
	// Obtener usuario
	var user domain.User
	if err := uc.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.ErrUserNotFound
		}
		uc.log.Error("Error finding user", logger.Int64("user_id", userID), logger.Error(err))
		return "", errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Generar token de reset seguro
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		uc.log.Error("Error generating reset token", logger.Error(err))
		return "", errors.Wrap(errors.ErrInternalServer, err)
	}
	resetToken := base64.URLEncoding.EncodeToString(tokenBytes)

	// Establecer expiración (24 horas)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Actualizar usuario con token de reset
	updates := map[string]interface{}{
		"password_reset_token":      resetToken,
		"password_reset_expires_at": expiresAt,
	}

	if err := uc.db.Model(&domain.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		uc.log.Error("Error updating password reset token",
			logger.Int64("user_id", userID),
			logger.Error(err))
		return "", errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría
	uc.log.Warn("Admin initiated password reset for user",
		logger.Int64("admin_id", adminID),
		logger.Int64("user_id", userID),
		logger.String("user_email", user.Email),
		logger.String("action", "admin_reset_user_password"))

	// TODO: Enviar email con link de reset
	// resetLink := fmt.Sprintf("https://sorteos.club/reset-password?token=%s", resetToken)
	// Esto se implementará cuando tengamos el servicio de email configurado

	return resetToken, nil
}
