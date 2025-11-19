package wallet

import (
	"context"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ListWalletsInput entrada para listar billeteras
type ListWalletsInput struct {
	Page     int
	PageSize int
	Search   string // Buscar por email o nombre de usuario
	Status   string // active, frozen, closed
	OrderBy  string
}

// WalletSummary resumen de billetera para listado
type WalletSummary struct {
	ID                int64  `json:"id"`
	UUID              string `json:"uuid"`
	UserID            int64  `json:"user_id"`
	UserEmail         string `json:"user_email"`
	UserName          string `json:"user_name"`
	BalanceAvailable  string `json:"balance_available"`
	EarningsBalance   string `json:"earnings_balance"`
	PendingBalance    string `json:"pending_balance"`
	TotalBalance      string `json:"total_balance"` // Suma de available + earnings
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// Pagination info
type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// ListWalletsOutput salida con paginación
type ListWalletsOutput struct {
	Wallets    []WalletSummary `json:"wallets"`
	Pagination Pagination      `json:"pagination"`
}

// ListWalletsUseCase lista billeteras con filtros y paginación
type ListWalletsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListWalletsUseCase crea una nueva instancia
func NewListWalletsUseCase(db *gorm.DB, log *logger.Logger) *ListWalletsUseCase {
	return &ListWalletsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListWalletsUseCase) Execute(ctx context.Context, input *ListWalletsInput) (*ListWalletsOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Construir query
	query := uc.db.WithContext(ctx).Table("wallets").
		Select("wallets.*, users.email as user_email, CONCAT(users.first_name, ' ', users.last_name) as user_name").
		Joins("LEFT JOIN users ON users.id = wallets.user_id")

	// Filtrar por búsqueda
	if input.Search != "" {
		searchPattern := "%" + input.Search + "%"
		query = query.Where("users.email ILIKE ? OR CONCAT(users.first_name, ' ', users.last_name) ILIKE ?", searchPattern, searchPattern)
	}

	// Filtrar por estado
	if input.Status != "" {
		query = query.Where("wallets.status = ?", input.Status)
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		uc.log.Error("Error counting wallets", logger.Error(err))
		return nil, err
	}

	// Ordenar
	orderBy := "wallets.created_at DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}
	query = query.Order(orderBy)

	// Paginar
	offset := (input.Page - 1) * input.PageSize
	query = query.Limit(input.PageSize).Offset(offset)

	// Ejecutar query
	var wallets []struct {
		domain.Wallet
		UserEmail string
		UserName  string
	}

	if err := query.Scan(&wallets).Error; err != nil {
		uc.log.Error("Error listing wallets", logger.Error(err))
		return nil, err
	}

	// Mapear a DTO
	summaries := make([]WalletSummary, len(wallets))
	for i, w := range wallets {
		totalBalance := w.Wallet.BalanceAvailable.Add(w.Wallet.EarningsBalance)
		summaries[i] = WalletSummary{
			ID:               w.ID,
			UUID:             w.UUID,
			UserID:           w.UserID,
			UserEmail:        w.UserEmail,
			UserName:         w.UserName,
			BalanceAvailable: w.Wallet.BalanceAvailable.String(),
			EarningsBalance:  w.Wallet.EarningsBalance.String(),
			PendingBalance:   w.Wallet.PendingBalance.String(),
			TotalBalance:     totalBalance.String(),
			Currency:         w.Currency,
			Status:           string(w.Status),
			CreatedAt:        w.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:        w.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	return &ListWalletsOutput{
		Wallets: summaries,
		Pagination: Pagination{
			Total:      total,
			Page:       input.Page - 1, // Convertir de 1-indexed a 0-indexed para el frontend
			Limit:      input.PageSize,
			TotalPages: totalPages,
		},
	}, nil
}
