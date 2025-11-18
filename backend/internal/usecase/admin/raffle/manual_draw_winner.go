package raffle

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ManualDrawWinnerInput datos de entrada
type ManualDrawWinnerInput struct {
	RaffleID      int64
	WinnerNumber  *string // Si es nil, se selecciona aleatoriamente
	Reason        string  // Razón del sorteo manual
}

// ManualDrawWinnerOutput resultado
type ManualDrawWinnerOutput struct {
	WinnerNumber string
	WinnerUserID *int64
	WinnerName   *string
	WinnerEmail  *string
}

// ManualDrawWinnerUseCase caso de uso para ejecutar sorteo manual
type ManualDrawWinnerUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewManualDrawWinnerUseCase crea una nueva instancia
func NewManualDrawWinnerUseCase(db *gorm.DB, log *logger.Logger) *ManualDrawWinnerUseCase {
	return &ManualDrawWinnerUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ManualDrawWinnerUseCase) Execute(ctx context.Context, input *ManualDrawWinnerInput, adminID int64) (*ManualDrawWinnerOutput, error) {
	// Validar razón
	if input.Reason == "" {
		return nil, errors.New("VALIDATION_FAILED", "reason is required for manual draw", 400, nil)
	}

	// Obtener rifa
	var raffle domain.Raffle
	if err := uc.db.Where("id = ?", input.RaffleID).First(&raffle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		uc.log.Error("Error finding raffle", logger.Int64("raffle_id", input.RaffleID), logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Validar que la rifa esté en estado active o suspended (permitimos sortear suspended)
	if raffle.Status != domain.RaffleStatusActive && raffle.Status != domain.RaffleStatusSuspended {
		return nil, errors.New("VALIDATION_FAILED",
			fmt.Sprintf("raffle must be active or suspended to draw winner, current status: %s", raffle.Status), 400, nil)
	}

	// Validar que la rifa no tenga ganador ya
	if raffle.WinnerNumber != nil {
		return nil, errors.New("VALIDATION_FAILED", "raffle already has a winner", 400, nil)
	}

	// Determinar número ganador
	var winnerNumber string
	if input.WinnerNumber != nil && *input.WinnerNumber != "" {
		// Usar número especificado por admin
		winnerNumber = *input.WinnerNumber

		// Validar que el número esté en el rango válido
		// (Aquí asumimos que los números son strings, podría ser necesario convertir)
		// TODO: Validar que el número esté vendido
	} else {
		// Seleccionar aleatoriamente de los números vendidos
		selectedNumber, err := uc.selectRandomSoldNumber(input.RaffleID)
		if err != nil {
			return nil, err
		}
		winnerNumber = selectedNumber
	}

	// Obtener información del ganador desde raffle_numbers
	var raffleNumber struct {
		UserID *int64
	}

	if err := uc.db.Table("raffle_numbers").
		Select("user_id").
		Where("raffle_id = ? AND number = ?", input.RaffleID, winnerNumber).
		First(&raffleNumber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("VALIDATION_FAILED",
				fmt.Sprintf("number %s not found or not sold", winnerNumber), 400, nil)
		}
		uc.log.Error("Error finding winner number", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Obtener info del usuario ganador si existe
	var winnerName, winnerEmail *string
	if raffleNumber.UserID != nil {
		var user domain.User
		if err := uc.db.Where("id = ?", *raffleNumber.UserID).First(&user).Error; err == nil {
			name := user.GetFullName()
			winnerName = &name
			winnerEmail = &user.Email
		}
	}

	now := time.Now()

	// Actualizar rifa con ganador y marcar como completed
	updates := map[string]interface{}{
		"winner_number": winnerNumber,
		"winner_user_id": raffleNumber.UserID,
		"status": domain.RaffleStatusCompleted,
		"completed_at": now,
		"updated_at": now,
		"admin_notes": fmt.Sprintf("Manual draw by admin ID %d. Reason: %s", adminID, input.Reason),
	}

	if err := uc.db.Model(&domain.Raffle{}).
		Where("id = ?", input.RaffleID).
		Updates(updates).Error; err != nil {
		uc.log.Error("Error updating raffle with winner",
			logger.Int64("raffle_id", input.RaffleID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Log auditoría crítica
	uc.log.Error("Admin manually drew winner for raffle",
		logger.Int64("admin_id", adminID),
		logger.Int64("raffle_id", input.RaffleID),
		logger.String("winner_number", winnerNumber),
		logger.String("reason", input.Reason),
		logger.String("action", "admin_manual_draw_winner"),
		logger.String("severity", "critical"))

	// TODO: Enviar emails
	// - Al ganador
	// - Al organizador
	// Esto se implementará cuando tengamos el servicio de email configurado

	return &ManualDrawWinnerOutput{
		WinnerNumber: winnerNumber,
		WinnerUserID: raffleNumber.UserID,
		WinnerName:   winnerName,
		WinnerEmail:  winnerEmail,
	}, nil
}

// selectRandomSoldNumber selecciona aleatoriamente un número vendido
func (uc *ManualDrawWinnerUseCase) selectRandomSoldNumber(raffleID int64) (string, error) {
	// Obtener todos los números vendidos
	var soldNumbers []string
	if err := uc.db.Table("raffle_numbers").
		Select("number").
		Where("raffle_id = ? AND user_id IS NOT NULL", raffleID).
		Pluck("number", &soldNumbers).Error; err != nil {
		uc.log.Error("Error getting sold numbers", logger.Error(err))
		return "", errors.Wrap(errors.ErrDatabaseError, err)
	}

	if len(soldNumbers) == 0 {
		return "", errors.New("VALIDATION_FAILED", "no sold numbers available for draw", 400, nil)
	}

	// Seleccionar aleatoriamente
	maxBig := big.NewInt(int64(len(soldNumbers)))
	randomIndex, err := rand.Int(rand.Reader, maxBig)
	if err != nil {
		uc.log.Error("Error generating random number", logger.Error(err))
		return "", errors.Wrap(errors.ErrInternalServer, err)
	}

	return soldNumbers[randomIndex.Int64()], nil
}
