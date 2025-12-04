package raffle

import (
	"github.com/sorteos-platform/backend/internal/domain"
)

// OrganizerInfo información pública del organizador (sin exponer user_id)
type OrganizerInfo struct {
	Name     string `json:"name"`
	Verified bool   `json:"verified"`
}

// PublicRaffleDTO - Para usuarios NO autenticados o compradores sin compras
// Oculta información financiera y sensible
type PublicRaffleDTO struct {
	ID            int64          `json:"id"`
	UUID          string         `json:"uuid"`
	Organizer     *OrganizerInfo `json:"organizer"` // En lugar de user_id
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Status        string         `json:"status"`
	PricePerNumber string        `json:"price_per_number"`
	TotalNumbers  int            `json:"total_numbers"`
	DrawDate      string         `json:"draw_date"`
	DrawMethod    string         `json:"draw_method"`
	SoldCount     int            `json:"sold_count"`
	ReservedCount int            `json:"reserved_count"`
	AvailableCount int           `json:"available_count"`
	CategoryID    *int64         `json:"category_id,omitempty"`
	CreatedAt     string         `json:"created_at"`
	PublishedAt   *string        `json:"published_at,omitempty"`
}

// BuyerRaffleDTO - Para usuarios autenticados que HAN comprado en este sorteo
// Incluye su gasto personal pero NO la información financiera del organizador
type BuyerRaffleDTO struct {
	PublicRaffleDTO            // Hereda campos públicos
	MyTotalSpent    string     `json:"my_total_spent"`     // Cuánto ha gastado este usuario
	MyNumbersCount  int        `json:"my_numbers_count"`   // Cuántos números tiene
}

// OwnerRaffleDTO - Para el organizador del sorteo y admins
// Incluye TODA la información financiera
type OwnerRaffleDTO struct {
	PublicRaffleDTO                        // Hereda campos públicos
	UserID                int64            `json:"user_id"` // Solo visible para owner/admin
	TotalRevenue          string           `json:"total_revenue"`
	PlatformFeePercentage string           `json:"platform_fee_percentage"`
	PlatformFeeAmount     string           `json:"platform_fee_amount"`
	NetAmount             string           `json:"net_amount"`
	SettlementStatus      string           `json:"settlement_status"`
}

// toPublicRaffleDTO convierte un domain.Raffle a DTO público (legacy, usa toPublicRaffleDTOWithOrganizer)
func toPublicRaffleDTO(raffle *domain.Raffle) *PublicRaffleDTO {
	return toPublicRaffleDTOWithOrganizer(raffle, &OrganizerInfo{
		Name:     "Organizador",
		Verified: false,
	})
}

// toPublicRaffleDTOWithOrganizer convierte un domain.Raffle a DTO público con info del organizador
func toPublicRaffleDTOWithOrganizer(raffle *domain.Raffle, organizer *OrganizerInfo) *PublicRaffleDTO {
	dto := &PublicRaffleDTO{
		ID:             raffle.ID,
		UUID:           raffle.UUID.String(),
		Title:          raffle.Title,
		Description:    raffle.Description,
		Status:         string(raffle.Status),
		PricePerNumber: raffle.PricePerNumber.String(),
		TotalNumbers:   raffle.TotalNumbers,
		DrawDate:       raffle.DrawDate.Format("2006-01-02T15:04:05Z07:00"),
		DrawMethod:     string(raffle.DrawMethod),
		SoldCount:      raffle.SoldCount,
		ReservedCount:  raffle.ReservedCount,
		AvailableCount: raffle.TotalNumbers - raffle.SoldCount - raffle.ReservedCount,
		CategoryID:     raffle.CategoryID,
		CreatedAt:      raffle.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Organizer:      organizer,
	}

	if raffle.PublishedAt != nil {
		publishedStr := raffle.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		dto.PublishedAt = &publishedStr
	}

	return dto
}

// toBuyerRaffleDTO convierte a DTO de comprador con su gasto personal (legacy)
func toBuyerRaffleDTO(raffle *domain.Raffle, myTotalSpent string, myNumbersCount int) *BuyerRaffleDTO {
	publicDTO := toPublicRaffleDTO(raffle)

	return &BuyerRaffleDTO{
		PublicRaffleDTO: *publicDTO,
		MyTotalSpent:    myTotalSpent,
		MyNumbersCount:  myNumbersCount,
	}
}

// toBuyerRaffleDTOWithOrganizer convierte a DTO de comprador con info del organizador
func toBuyerRaffleDTOWithOrganizer(raffle *domain.Raffle, organizer *OrganizerInfo, myTotalSpent string, myNumbersCount int) *BuyerRaffleDTO {
	publicDTO := toPublicRaffleDTOWithOrganizer(raffle, organizer)

	return &BuyerRaffleDTO{
		PublicRaffleDTO: *publicDTO,
		MyTotalSpent:    myTotalSpent,
		MyNumbersCount:  myNumbersCount,
	}
}

// toOwnerRaffleDTO convierte a DTO completo (solo para owner/admin) (legacy)
func toOwnerRaffleDTO(raffle *domain.Raffle) *OwnerRaffleDTO {
	publicDTO := toPublicRaffleDTO(raffle)

	return &OwnerRaffleDTO{
		PublicRaffleDTO:       *publicDTO,
		UserID:                raffle.UserID,
		TotalRevenue:          raffle.TotalRevenue.String(),
		PlatformFeePercentage: raffle.PlatformFeePercentage.String(),
		PlatformFeeAmount:     raffle.PlatformFeeAmount.String(),
		NetAmount:             raffle.NetAmount.String(),
		SettlementStatus:      string(raffle.SettlementStatus),
	}
}

// toOwnerRaffleDTOWithOrganizer convierte a DTO completo con info del organizador
func toOwnerRaffleDTOWithOrganizer(raffle *domain.Raffle, organizer *OrganizerInfo) *OwnerRaffleDTO {
	publicDTO := toPublicRaffleDTOWithOrganizer(raffle, organizer)

	return &OwnerRaffleDTO{
		PublicRaffleDTO:       *publicDTO,
		UserID:                raffle.UserID,
		TotalRevenue:          raffle.TotalRevenue.String(),
		PlatformFeePercentage: raffle.PlatformFeePercentage.String(),
		PlatformFeeAmount:     raffle.PlatformFeeAmount.String(),
		NetAmount:             raffle.NetAmount.String(),
		SettlementStatus:      string(raffle.SettlementStatus),
	}
}
