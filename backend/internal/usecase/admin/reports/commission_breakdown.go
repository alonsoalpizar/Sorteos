package reports

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// CommissionTier tier de comisión
type CommissionTier struct {
	CommissionPercent  float64                  `json:"commission_percent"`
	RaffleCount        int64                    `json:"raffle_count"`
	GrossRevenue       float64                  `json:"gross_revenue"`
	FeesCollected      float64                  `json:"fees_collected"`
	OrganizerCount     int64                    `json:"organizer_count"`
	Organizers         []*CommissionOrganizer   `json:"organizers,omitempty"`
}

// CommissionOrganizer organizador en un tier
type CommissionOrganizer struct {
	OrganizerID   int64   `json:"organizer_id"`
	OrganizerName string  `json:"organizer_name"`
	RaffleCount   int64   `json:"raffle_count"`
	GrossRevenue  float64 `json:"gross_revenue"`
	FeesCollected float64 `json:"fees_collected"`
}

// CommissionBreakdownInput datos de entrada
type CommissionBreakdownInput struct {
	DateFrom         string
	DateTo           string
	IncludeOrganizers bool // Si true, incluye lista de organizadores por tier
}

// CommissionBreakdownOutput resultado
type CommissionBreakdownOutput struct {
	Tiers               []*CommissionTier
	TotalRaffles        int64
	TotalGrossRevenue   float64
	TotalFeesCollected  float64
	DefaultTierPercent  float64
	CustomTiersCount    int64
}

// CommissionBreakdownUseCase caso de uso para desglose de comisiones
type CommissionBreakdownUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewCommissionBreakdownUseCase crea una nueva instancia
func NewCommissionBreakdownUseCase(db *gorm.DB, log *logger.Logger) *CommissionBreakdownUseCase {
	return &CommissionBreakdownUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *CommissionBreakdownUseCase) Execute(ctx context.Context, input *CommissionBreakdownInput, adminID int64) (*CommissionBreakdownOutput, error) {
	// Validar fechas
	if input.DateFrom == "" || input.DateTo == "" {
		return nil, errors.New("VALIDATION_FAILED", "date_from and date_to are required", 400, nil)
	}

	defaultCommissionPercent := 10.0

	// Obtener rifas completadas en el período con revenue y organizadores
	type RaffleData struct {
		OrganizerID        int64
		OrganizerFirstName *string
		OrganizerLastName  *string
		OrganizerEmail     string
		CustomCommission   *float64
		Revenue            float64
	}

	var raffleData []RaffleData

	// Query: raffles completadas con revenue calculado desde payments
	rows, err := uc.db.Raw(`
		SELECT
			raffles.user_id as organizer_id,
			users.first_name as organizer_first_name,
			users.last_name as organizer_last_name,
			users.email as organizer_email,
			organizer_profiles.custom_commission_rate as custom_commission,
			COALESCE(SUM(payments.amount), 0) as revenue
		FROM raffles
		LEFT JOIN users ON users.id = raffles.user_id
		LEFT JOIN organizer_profiles ON organizer_profiles.user_id = raffles.user_id
		LEFT JOIN payments ON payments.raffle_id = raffles.uuid::text AND payments.status = 'succeeded'
		WHERE raffles.status = 'completed'
			AND raffles.completed_at >= ?
			AND raffles.completed_at <= ?
		GROUP BY raffles.user_id, users.first_name, users.last_name, users.email, organizer_profiles.custom_commission_rate
	`, input.DateFrom, input.DateTo+" 23:59:59").Rows()

	if err != nil {
		uc.log.Error("Error fetching raffle data for commission breakdown", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	for rows.Next() {
		var data RaffleData
		if err := rows.Scan(
			&data.OrganizerID,
			&data.OrganizerFirstName,
			&data.OrganizerLastName,
			&data.OrganizerEmail,
			&data.CustomCommission,
			&data.Revenue,
		); err != nil {
			uc.log.Error("Error scanning raffle data", logger.Error(err))
			continue
		}
		raffleData = append(raffleData, data)
	}

	// Agrupar por commission percent
	tierMap := make(map[float64]*CommissionTier)
	organizerMap := make(map[string]map[int64]*CommissionOrganizer) // tier -> organizer_id -> data

	for _, data := range raffleData {
		commissionPercent := defaultCommissionPercent
		if data.CustomCommission != nil {
			commissionPercent = *data.CustomCommission
		}

		// Crear tier si no existe
		if _, exists := tierMap[commissionPercent]; !exists {
			tierMap[commissionPercent] = &CommissionTier{
				CommissionPercent: commissionPercent,
				Organizers:        make([]*CommissionOrganizer, 0),
			}
			organizerMap[fmt.Sprintf("%.2f", commissionPercent)] = make(map[int64]*CommissionOrganizer)
		}

		tier := tierMap[commissionPercent]

		// Contar raffles para este organizador
		var raffleCount int64
		uc.db.Table("raffles").
			Where("user_id = ?", data.OrganizerID).
			Where("status = ?", "completed").
			Where("completed_at >= ?", input.DateFrom).
			Where("completed_at <= ?", input.DateTo+" 23:59:59").
			Count(&raffleCount)

		tier.RaffleCount += raffleCount
		tier.GrossRevenue += data.Revenue
		tier.FeesCollected += data.Revenue * commissionPercent / 100.0

		// Si se requiere incluir organizadores
		if input.IncludeOrganizers {
			tierKey := fmt.Sprintf("%.2f", commissionPercent)
			if _, exists := organizerMap[tierKey][data.OrganizerID]; !exists {
				organizerName := data.OrganizerEmail
				if data.OrganizerFirstName != nil && data.OrganizerLastName != nil {
					organizerName = *data.OrganizerFirstName + " " + *data.OrganizerLastName
				}

				organizerMap[tierKey][data.OrganizerID] = &CommissionOrganizer{
					OrganizerID:   data.OrganizerID,
					OrganizerName: organizerName,
					RaffleCount:   raffleCount,
					GrossRevenue:  data.Revenue,
					FeesCollected: data.Revenue * commissionPercent / 100.0,
				}
			}
		}
	}

	// Convertir map a slice y calcular organizer count
	tiers := make([]*CommissionTier, 0, len(tierMap))
	for commissionPercent, tier := range tierMap {
		// Agregar organizadores al tier
		if input.IncludeOrganizers {
			tierKey := fmt.Sprintf("%.2f", commissionPercent)
			for _, org := range organizerMap[tierKey] {
				tier.Organizers = append(tier.Organizers, org)
			}
			tier.OrganizerCount = int64(len(tier.Organizers))
		} else {
			// Contar organizadores únicos para este tier
			var count int64
			if commissionPercent == defaultCommissionPercent {
				// Default tier: organizadores sin custom commission
				uc.db.Table("users").
					Joins("LEFT JOIN organizer_profiles ON organizer_profiles.user_id = users.id").
					Where("users.role = ?", "organizer").
					Where("organizer_profiles.custom_commission_rate IS NULL OR organizer_profiles.custom_commission_rate = ?", defaultCommissionPercent).
					Count(&count)
			} else {
				// Custom tier
				uc.db.Table("organizer_profiles").
					Where("custom_commission_rate = ?", commissionPercent).
					Count(&count)
			}
			tier.OrganizerCount = count
		}

		tiers = append(tiers, tier)
	}

	// Ordenar tiers por commission_percent
	for i := 0; i < len(tiers)-1; i++ {
		for j := i + 1; j < len(tiers); j++ {
			if tiers[i].CommissionPercent > tiers[j].CommissionPercent {
				tiers[i], tiers[j] = tiers[j], tiers[i]
			}
		}
	}

	// Calcular totales
	totalRaffles := int64(0)
	totalGrossRevenue := 0.0
	totalFeesCollected := 0.0
	customTiersCount := int64(0)

	for _, tier := range tiers {
		totalRaffles += tier.RaffleCount
		totalGrossRevenue += tier.GrossRevenue
		totalFeesCollected += tier.FeesCollected

		if tier.CommissionPercent != defaultCommissionPercent {
			customTiersCount++
		}
	}

	// Log auditoría
	uc.log.Info("Admin generated commission breakdown",
		logger.Int64("admin_id", adminID),
		logger.String("date_from", input.DateFrom),
		logger.String("date_to", input.DateTo),
		logger.Int("tier_count", len(tiers)),
		logger.Float64("total_fees_collected", totalFeesCollected),
		logger.String("action", "admin_commission_breakdown"))

	return &CommissionBreakdownOutput{
		Tiers:              tiers,
		TotalRaffles:       totalRaffles,
		TotalGrossRevenue:  totalGrossRevenue,
		TotalFeesCollected: totalFeesCollected,
		DefaultTierPercent: defaultCommissionPercent,
		CustomTiersCount:   customTiersCount,
	}, nil
}
