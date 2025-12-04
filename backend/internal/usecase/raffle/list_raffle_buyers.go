package raffle

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/adapters/db"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/errors"
)

// BuyerInfo información de un comprador/reservador
type BuyerInfo struct {
	UserID      int64      `json:"user_id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone,omitempty"`
	Numbers     []string   `json:"numbers"`
	TotalAmount string     `json:"total_amount"`
	Status      string     `json:"status"` // "sold" o "reserved"
	ReservedAt  *time.Time `json:"reserved_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	SoldAt      *time.Time `json:"sold_at,omitempty"`
}

// ListRaffleBuyersInput datos de entrada
type ListRaffleBuyersInput struct {
	RaffleUUID   string
	OwnerUserID  int64 // El usuario que hace la solicitud (debe ser owner)
	IncludeSold  bool  // Incluir números vendidos
	IncludeReserved bool // Incluir números reservados
}

// ListRaffleBuyersOutput resultado
type ListRaffleBuyersOutput struct {
	Buyers      []BuyerInfo `json:"buyers"`
	TotalSold   int         `json:"total_sold"`
	TotalReserved int       `json:"total_reserved"`
}

// ListRaffleBuyersUseCase caso de uso para listar compradores de un sorteo
type ListRaffleBuyersUseCase struct {
	raffleRepo       db.RaffleRepository
	raffleNumberRepo db.RaffleNumberRepository
	userRepo         domain.UserRepository
}

// NewListRaffleBuyersUseCase crea una nueva instancia
func NewListRaffleBuyersUseCase(
	raffleRepo db.RaffleRepository,
	raffleNumberRepo db.RaffleNumberRepository,
	userRepo domain.UserRepository,
) *ListRaffleBuyersUseCase {
	return &ListRaffleBuyersUseCase{
		raffleRepo:       raffleRepo,
		raffleNumberRepo: raffleNumberRepo,
		userRepo:         userRepo,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListRaffleBuyersUseCase) Execute(ctx context.Context, input *ListRaffleBuyersInput) (*ListRaffleBuyersOutput, error) {
	// 1. Buscar el sorteo por UUID
	raffle, err := uc.raffleRepo.FindByUUID(input.RaffleUUID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrRaffleNotFound
		}
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 2. Verificar que el usuario sea el owner del sorteo
	if raffle.UserID != input.OwnerUserID {
		return nil, errors.ErrForbidden
	}

	// 3. Obtener todos los números del sorteo
	numbers, err := uc.raffleNumberRepo.FindByRaffleID(raffle.ID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// 4. Agrupar por usuario
	buyerMap := make(map[int64]*BuyerInfo)
	totalSold := 0
	totalReserved := 0

	for _, num := range numbers {
		var userID *int64
		var status string
		var reservedAt, expiresAt, soldAt *time.Time

		if num.Status == domain.RaffleNumberStatusSold && num.UserID != nil {
			if !input.IncludeSold {
				continue
			}
			userID = num.UserID
			status = "sold"
			soldAt = num.SoldAt
			totalSold++
		} else if num.Status == domain.RaffleNumberStatusReserved && num.ReservedBy != nil {
			if !input.IncludeReserved {
				continue
			}
			userID = num.ReservedBy
			status = "reserved"
			reservedAt = num.ReservedAt
			expiresAt = num.ReservedUntil
			totalReserved++
		} else {
			continue
		}

		if userID == nil {
			continue
		}

		// Agregar al mapa o crear nueva entrada
		if buyer, exists := buyerMap[*userID]; exists {
			buyer.Numbers = append(buyer.Numbers, num.Number)
			// Si hay mezcla de estados, priorizar el más reciente
			if status == "reserved" && buyer.Status == "sold" {
				// Mantener como sold + reserved
				buyer.Status = "mixed"
			}
		} else {
			buyerMap[*userID] = &BuyerInfo{
				UserID:     *userID,
				Numbers:    []string{num.Number},
				Status:     status,
				ReservedAt: reservedAt,
				ExpiresAt:  expiresAt,
				SoldAt:     soldAt,
			}
		}
	}

	// 5. Obtener info de usuarios
	buyers := make([]BuyerInfo, 0, len(buyerMap))
	pricePerNumber := raffle.PricePerNumber

	for userID, buyer := range buyerMap {
		user, err := uc.userRepo.FindByID(userID)
		if err != nil {
			// Si no encontramos el usuario, usar placeholders
			buyer.Name = "Usuario eliminado"
			buyer.Email = "N/A"
		} else {
			// Construir nombre
			if user.FirstName != nil && user.LastName != nil {
				buyer.Name = *user.FirstName + " " + *user.LastName
			} else if user.FirstName != nil {
				buyer.Name = *user.FirstName
			} else {
				buyer.Name = user.Email
			}
			buyer.Email = user.Email
			buyer.Phone = user.Phone
		}

		// Calcular total
		totalAmount := pricePerNumber.Mul(decimal.NewFromInt(int64(len(buyer.Numbers))))
		buyer.TotalAmount = totalAmount.String()

		buyers = append(buyers, *buyer)
	}

	return &ListRaffleBuyersOutput{
		Buyers:        buyers,
		TotalSold:     totalSold,
		TotalReserved: totalReserved,
	}, nil
}
