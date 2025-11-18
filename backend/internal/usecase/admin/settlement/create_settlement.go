package settlement

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// CreateSettlementInput datos de entrada
type CreateSettlementInput struct {
	OrganizerID int64   `json:"organizer_id"`
	RaffleIDs   []int64 `json:"raffle_ids,omitempty"` // IDs de rifas específicas (opcional)
	Mode        string  `json:"mode"`                  // individual, batch
}

// CreateSettlementOutput resultado
type CreateSettlementOutput struct {
	SettlementIDs []int64  `json:"settlement_ids"`
	TotalCreated  int      `json:"total_created"`
	TotalRevenue  float64  `json:"total_revenue"`
	TotalNetAmount float64 `json:"total_net_amount"`
	Message       string   `json:"message"`
}

// CreateSettlementUseCase caso de uso para crear liquidaciones
type CreateSettlementUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCreateSettlementUseCase crea una nueva instancia
func NewCreateSettlementUseCase(db *gorm.DB, log *logger.Logger) *CreateSettlementUseCase {
	return &CreateSettlementUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *CreateSettlementUseCase) Execute(ctx context.Context, input *CreateSettlementInput, adminID int64) (*CreateSettlementOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Verificar que el organizador existe
	var organizerCount int64
	uc.db.WithContext(ctx).Table("users").
		Where("id = ? AND role = ?", input.OrganizerID, "organizer").
		Count(&organizerCount)

	if organizerCount == 0 {
		return nil, errors.New("ORGANIZER_NOT_FOUND", "organizer not found", 404, nil)
	}

	// Obtener platform fee (custom o default)
	platformFeePercent := uc.getPlatformFeePercent(ctx, input.OrganizerID)

	// Obtener rifas elegibles para settlement
	raffles, err := uc.getEligibleRaffles(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(raffles) == 0 {
		return nil, errors.New("NO_ELIGIBLE_RAFFLES", "no eligible raffles found for settlement", 400, nil)
	}

	// Crear settlements
	var settlementIDs []int64
	var totalRevenue, totalNetAmount float64

	for _, raffle := range raffles {
		// Calcular montos
		grossRevenue := raffle.PricePerNumber * float64(raffle.SoldCount)
		platformFee := grossRevenue * (platformFeePercent / 100.0)
		netAmount := grossRevenue - platformFee

		// Crear settlement
		settlement := map[string]interface{}{
			"organizer_id":  input.OrganizerID,
			"raffle_id":     raffle.ID,
			"total_revenue": grossRevenue,
			"platform_fee":  platformFee,
			"net_amount":    netAmount,
			"status":        "pending",
			"calculated_at": time.Now(),
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		}

		var settlementID int64
		result := uc.db.WithContext(ctx).Table("settlements").Create(settlement).Scan(&settlementID)
		if result.Error != nil {
			uc.log.Error("Error creating settlement",
				logger.Int64("raffle_id", raffle.ID),
				logger.Error(result.Error))
			continue
		}

		settlementIDs = append(settlementIDs, settlementID)
		totalRevenue += grossRevenue
		totalNetAmount += netAmount

		uc.log.Info("Settlement created",
			logger.Int64("settlement_id", settlementID),
			logger.Int64("raffle_id", raffle.ID),
			logger.Float64("net_amount", netAmount))
	}

	// Log auditoría
	uc.log.Error("Admin created settlements",
		logger.Int64("admin_id", adminID),
		logger.Int64("organizer_id", input.OrganizerID),
		logger.String("mode", input.Mode),
		logger.Int("total_created", len(settlementIDs)),
		logger.Float64("total_net_amount", totalNetAmount),
		logger.String("action", "admin_create_settlement"),
		logger.String("severity", "info"))

	return &CreateSettlementOutput{
		SettlementIDs:  settlementIDs,
		TotalCreated:   len(settlementIDs),
		TotalRevenue:   totalRevenue,
		TotalNetAmount: totalNetAmount,
		Message:        "Settlements created successfully",
	}, nil
}

// validateInput valida los datos de entrada
func (uc *CreateSettlementUseCase) validateInput(input *CreateSettlementInput) error {
	if input.OrganizerID <= 0 {
		return errors.New("VALIDATION_FAILED", "organizer_id is required", 400, nil)
	}

	validModes := map[string]bool{
		"individual": true,
		"batch":      true,
	}
	if !validModes[input.Mode] {
		return errors.New("VALIDATION_FAILED", "mode must be 'individual' or 'batch'", 400, nil)
	}

	if input.Mode == "individual" && len(input.RaffleIDs) == 0 {
		return errors.New("VALIDATION_FAILED", "raffle_ids required for individual mode", 400, nil)
	}

	return nil
}

// getPlatformFeePercent obtiene el porcentaje de comisión (custom o default)
func (uc *CreateSettlementUseCase) getPlatformFeePercent(ctx context.Context, organizerID int64) float64 {
	var commissionOverride *float64
	uc.db.WithContext(ctx).
		Table("organizer_profiles").
		Select("commission_override").
		Where("user_id = ?", organizerID).
		Scan(&commissionOverride)

	if commissionOverride != nil {
		return *commissionOverride
	}

	// Default: 10%
	return 10.0
}

// getEligibleRaffles obtiene las rifas elegibles para crear settlements
func (uc *CreateSettlementUseCase) getEligibleRaffles(ctx context.Context, input *CreateSettlementInput) ([]RaffleForSettlement, error) {
	var raffles []RaffleForSettlement

	query := uc.db.WithContext(ctx).
		Table("raffles").
		Select("id, title, user_id, price_per_number, sold_count, completed_at").
		Where("user_id = ?", input.OrganizerID).
		Where("status = ?", "completed")

	// Filtro: que no tengan settlement ya creado
	query = query.Where("id NOT IN (SELECT raffle_id FROM settlements WHERE raffle_id IS NOT NULL)")

	// Si es modo individual, filtrar por raffle_ids específicos
	if input.Mode == "individual" && len(input.RaffleIDs) > 0 {
		query = query.Where("id IN ?", input.RaffleIDs)
	}

	result := query.Find(&raffles)
	if result.Error != nil {
		uc.log.Error("Error finding eligible raffles", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	return raffles, nil
}

// RaffleForSettlement datos de rifa para settlement
type RaffleForSettlement struct {
	ID             int64
	Title          string
	UserID         int64
	PricePerNumber float64
	SoldCount      int
	CompletedAt    *time.Time
}
