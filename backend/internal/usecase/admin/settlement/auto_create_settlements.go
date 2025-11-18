package settlement

import (
	"context"
	"fmt"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// AutoCreateSettlementsInput datos de entrada
type AutoCreateSettlementsInput struct {
	DaysAfterCompletion int  `json:"days_after_completion"` // Esperar X días después de que la rifa se complete
	DryRun              bool `json:"dry_run"`                // Si es true, solo simula sin crear settlements
}

// AutoCreateSettlementsOutput resultado
type AutoCreateSettlementsOutput struct {
	EligibleRaffles    int                 `json:"eligible_raffles"`
	SettlementsCreated int                 `json:"settlements_created"`
	TotalNetAmount     float64             `json:"total_net_amount"`
	TotalPlatformFees  float64             `json:"total_platform_fees"`
	DryRun             bool                `json:"dry_run"`
	ProcessedAt        string              `json:"processed_at"`
	Settlements        []*SettlementSummary `json:"settlements,omitempty"`
	Errors             []string            `json:"errors,omitempty"`
	Message            string              `json:"message"`
}

// SettlementSummary resumen de un settlement creado
type SettlementSummary struct {
	SettlementID  int64   `json:"settlement_id"`
	OrganizerID   int64   `json:"organizer_id"`
	RaffleID      int64   `json:"raffle_id"`
	RaffleTitle   string  `json:"raffle_title"`
	TotalRevenue  float64 `json:"total_revenue"`
	PlatformFee   float64 `json:"platform_fee"`
	NetAmount     float64 `json:"net_amount"`
	Status        string  `json:"status"`
}

// AutoCreateSettlementsUseCase caso de uso para crear settlements automáticamente (batch job)
type AutoCreateSettlementsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewAutoCreateSettlementsUseCase crea una nueva instancia
func NewAutoCreateSettlementsUseCase(db *gorm.DB, log *logger.Logger) *AutoCreateSettlementsUseCase {
	return &AutoCreateSettlementsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *AutoCreateSettlementsUseCase) Execute(ctx context.Context, input *AutoCreateSettlementsInput, adminID int64) (*AutoCreateSettlementsOutput, error) {
	// Validar inputs
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	// Construir output base
	output := &AutoCreateSettlementsOutput{
		EligibleRaffles:    0,
		SettlementsCreated: 0,
		TotalNetAmount:     0,
		TotalPlatformFees:  0,
		DryRun:             input.DryRun,
		ProcessedAt:        time.Now().Format(time.RFC3339),
		Settlements:        make([]*SettlementSummary, 0),
		Errors:             make([]string, 0),
	}

	// Calcular fecha límite
	cutoffDate := time.Now().AddDate(0, 0, -input.DaysAfterCompletion)

	// Buscar raffles elegibles para settlement
	// Criterios:
	// 1. status = 'completed'
	// 2. completed_at <= cutoffDate (han pasado X días desde que se completó)
	// 3. No tienen settlement creado
	// 4. sold_count > 0 (al menos se vendió un número)
	var eligibleRaffles []struct {
		ID             int64
		UserID         int64
		Title          string
		PricePerNumber float64
		SoldCount      int
		CompletedAt    *time.Time
	}

	result := uc.db.WithContext(ctx).
		Table("raffles").
		Select("id, user_id, title, price_per_number, sold_count, completed_at").
		Where("status = ?", "completed").
		Where("completed_at IS NOT NULL").
		Where("completed_at <= ?", cutoffDate).
		Where("sold_count > 0").
		Where("id NOT IN (SELECT raffle_id FROM settlements WHERE raffle_id IS NOT NULL)").
		Find(&eligibleRaffles)

	if result.Error != nil {
		uc.log.Error("Error finding eligible raffles", logger.Error(result.Error))
		return nil, errors.Wrap(errors.ErrDatabaseError, result.Error)
	}

	output.EligibleRaffles = len(eligibleRaffles)

	// Si no hay raffles elegibles, retornar
	if len(eligibleRaffles) == 0 {
		output.Message = "No eligible raffles found for settlement creation"
		uc.log.Info("Auto create settlements - no eligible raffles",
			logger.Int64("admin_id", adminID),
			logger.Int("days_after_completion", input.DaysAfterCompletion),
			logger.Bool("dry_run", input.DryRun))
		return output, nil
	}

	// Agrupar por organizador para batch processing
	rafflesByOrganizer := make(map[int64][]struct {
		ID             int64
		UserID         int64
		Title          string
		PricePerNumber float64
		SoldCount      int
		CompletedAt    *time.Time
	})

	for _, raffle := range eligibleRaffles {
		rafflesByOrganizer[raffle.UserID] = append(rafflesByOrganizer[raffle.UserID], raffle)
	}

	// Procesar cada organizador
	for organizerID, raffles := range rafflesByOrganizer {
		// Obtener comisión del organizador
		platformFeePercent := uc.getPlatformFeePercent(ctx, organizerID)

		// Procesar cada rifa
		for _, raffle := range raffles {
			totalRevenue := raffle.PricePerNumber * float64(raffle.SoldCount)
			platformFee := totalRevenue * (platformFeePercent / 100)
			netAmount := totalRevenue - platformFee

			// Si es dry run, solo simular
			if input.DryRun {
				summary := &SettlementSummary{
					SettlementID:  0, // No se crea en dry run
					OrganizerID:   organizerID,
					RaffleID:      raffle.ID,
					RaffleTitle:   raffle.Title,
					TotalRevenue:  totalRevenue,
					PlatformFee:   platformFee,
					NetAmount:     netAmount,
					Status:        "pending", // Status que tendría
				}
				output.Settlements = append(output.Settlements, summary)
				output.TotalNetAmount += netAmount
				output.TotalPlatformFees += platformFee
				output.SettlementsCreated++
				continue
			}

			// Crear settlement real
			now := time.Now()
			settlement := map[string]interface{}{
				"organizer_id":  organizerID,
				"raffle_id":     raffle.ID,
				"total_revenue": totalRevenue,
				"platform_fee":  platformFee,
				"net_amount":    netAmount,
				"status":        "pending",
				"calculated_at": now,
				"created_at":    now,
				"updated_at":    now,
			}

			result := uc.db.WithContext(ctx).
				Table("settlements").
				Create(settlement)

			if result.Error != nil {
				errMsg := fmt.Sprintf("Failed to create settlement for raffle %d: %v", raffle.ID, result.Error)
				output.Errors = append(output.Errors, errMsg)
				uc.log.Error("Error creating settlement",
					logger.Int64("raffle_id", raffle.ID),
					logger.Int64("organizer_id", organizerID),
					logger.Error(result.Error))
				continue
			}

			// Obtener el ID del settlement creado
			var settlementID int64
			uc.db.WithContext(ctx).
				Table("settlements").
				Select("id").
				Where("raffle_id = ?", raffle.ID).
				Scan(&settlementID)

			// Actualizar organizer_profile (incrementar pending_payout)
			err := uc.updateOrganizerProfile(ctx, organizerID, netAmount)
			if err != nil {
				uc.log.Error("Error updating organizer profile",
					logger.Int64("organizer_id", organizerID),
					logger.Error(err))
				// No fallar la operación, solo loguear
			}

			// Agregar al resumen
			summary := &SettlementSummary{
				SettlementID:  settlementID,
				OrganizerID:   organizerID,
				RaffleID:      raffle.ID,
				RaffleTitle:   raffle.Title,
				TotalRevenue:  totalRevenue,
				PlatformFee:   platformFee,
				NetAmount:     netAmount,
				Status:        "pending",
			}
			output.Settlements = append(output.Settlements, summary)
			output.TotalNetAmount += netAmount
			output.TotalPlatformFees += platformFee
			output.SettlementsCreated++

			uc.log.Info("Settlement created automatically",
				logger.Int64("settlement_id", settlementID),
				logger.Int64("raffle_id", raffle.ID),
				logger.Int64("organizer_id", organizerID),
				logger.Float64("net_amount", netAmount))
		}
	}

	// Log auditoría crítica
	uc.log.Error("Admin executed auto create settlements",
		logger.Int64("admin_id", adminID),
		logger.Int("eligible_raffles", output.EligibleRaffles),
		logger.Int("settlements_created", output.SettlementsCreated),
		logger.Float64("total_net_amount", output.TotalNetAmount),
		logger.Float64("total_platform_fees", output.TotalPlatformFees),
		logger.Bool("dry_run", input.DryRun),
		logger.Int("errors", len(output.Errors)),
		logger.String("action", "admin_auto_create_settlements"),
		logger.String("severity", "critical"))

	// Construir mensaje
	if input.DryRun {
		output.Message = fmt.Sprintf("Dry run completed: %d settlements would be created", output.SettlementsCreated)
	} else {
		output.Message = fmt.Sprintf("Successfully created %d settlements from %d eligible raffles", output.SettlementsCreated, output.EligibleRaffles)
	}

	return output, nil
}

// validateInput valida los datos de entrada
func (uc *AutoCreateSettlementsUseCase) validateInput(input *AutoCreateSettlementsInput) error {
	if input.DaysAfterCompletion < 0 {
		return errors.New("VALIDATION_FAILED", "days_after_completion must be >= 0", 400, nil)
	}

	// Validar que no sea excesivo (máximo 365 días)
	if input.DaysAfterCompletion > 365 {
		return errors.New("VALIDATION_FAILED", "days_after_completion must be <= 365", 400, nil)
	}

	return nil
}

// getPlatformFeePercent obtiene el porcentaje de comisión del organizador
func (uc *AutoCreateSettlementsUseCase) getPlatformFeePercent(ctx context.Context, organizerID int64) float64 {
	var commissionOverride *float64

	uc.db.WithContext(ctx).
		Table("organizer_profiles").
		Select("commission_override").
		Where("user_id = ?", organizerID).
		Scan(&commissionOverride)

	if commissionOverride != nil && *commissionOverride > 0 {
		return *commissionOverride
	}

	// Default commission: 10%
	return 10.0
}

// updateOrganizerProfile actualiza las métricas del organizador
func (uc *AutoCreateSettlementsUseCase) updateOrganizerProfile(ctx context.Context, organizerID int64, pendingAmount float64) error {
	// Incrementar pending_payout

	result := uc.db.WithContext(ctx).Exec(`
		UPDATE organizer_profiles
		SET
			pending_payout = COALESCE(pending_payout, 0) + ?,
			updated_at = ?
		WHERE user_id = ?
	`, pendingAmount, time.Now(), organizerID)

	if result.Error != nil {
		return result.Error
	}

	// Si no existe el perfil, crearlo
	if result.RowsAffected == 0 {
		profile := map[string]interface{}{
			"user_id":        organizerID,
			"pending_payout": pendingAmount,
			"total_payouts":  0,
			"created_at":     time.Now(),
			"updated_at":     time.Now(),
		}
		uc.db.WithContext(ctx).Table("organizer_profiles").Create(profile)
	}

	return nil
}
