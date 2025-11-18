package user

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// UpdateUserKYCInput datos de entrada
type UpdateUserKYCInput struct {
	UserID   int64
	KYCLevel domain.KYCLevel
	Notes    string // Notas del revisor
}

// UpdateUserKYCUseCase caso de uso para actualizar KYC de usuario
type UpdateUserKYCUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUpdateUserKYCUseCase crea una nueva instancia
func NewUpdateUserKYCUseCase(db *gorm.DB, log *logger.Logger) *UpdateUserKYCUseCase {
	return &UpdateUserKYCUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *UpdateUserKYCUseCase) Execute(ctx context.Context, input *UpdateUserKYCInput, adminID int64) error {
	// Validar KYC level
	validLevels := []domain.KYCLevel{
		domain.KYCLevelNone,
		domain.KYCLevelEmailVerified,
		domain.KYCLevelPhoneVerified,
		domain.KYCLevelCedulaVerified,
		domain.KYCLevelFullKYC,
	}

	validLevel := false
	for _, level := range validLevels {
		if input.KYCLevel == level {
			validLevel = true
			break
		}
	}

	if !validLevel {
		return errors.New("VALIDATION_FAILED", "invalid KYC level", 400, nil)
	}

	// Obtener usuario
	var user domain.User
	if err := uc.db.Where("id = ?", input.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		uc.log.Error("Error finding user", logger.Int64("user_id", input.UserID), logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	now := time.Now()

	// Actualizar KYC
	updates := map[string]interface{}{
		"kyc_level":       input.KYCLevel,
		"last_kyc_review": now,
		"kyc_reviewer":    adminID,
	}

	if err := uc.db.Model(&domain.User{}).Where("id = ?", input.UserID).Updates(updates).Error; err != nil {
		uc.log.Error("Error updating user KYC",
			logger.Int64("user_id", input.UserID),
			logger.String("new_level", string(input.KYCLevel)),
			logger.Error(err))
		return errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría
	uc.log.Info("Admin updated user KYC",
		logger.Int64("admin_id", adminID),
		logger.Int64("user_id", input.UserID),
		logger.String("old_level", string(user.KYCLevel)),
		logger.String("new_level", string(input.KYCLevel)),
		logger.String("notes", input.Notes),
		logger.String("action", "admin_update_user_kyc"))

	// TODO: Enviar email de notificación al usuario
	// Esto se implementará cuando tengamos el servicio de email configurado

	return nil
}
