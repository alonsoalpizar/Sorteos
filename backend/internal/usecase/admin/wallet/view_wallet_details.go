package wallet

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ViewWalletDetailsInput entrada
type ViewWalletDetailsInput struct {
	WalletID int64
}

// WalletDetails detalles completos de una billetera
type WalletDetails struct {
	ID               int64  `json:"id"`
	UUID             string `json:"uuid"`
	UserID           int64  `json:"user_id"`
	UserEmail        string `json:"user_email"`
	UserName         string `json:"user_name"`
	BalanceAvailable string `json:"balance_available"`
	EarningsBalance  string `json:"earnings_balance"`
	PendingBalance   string `json:"pending_balance"`
	TotalBalance     string `json:"total_balance"`
	Currency         string `json:"currency"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// ViewWalletDetailsOutput salida
type ViewWalletDetailsOutput struct {
	Wallet WalletDetails `json:"wallet"`
}

// ViewWalletDetailsUseCase obtiene detalles de una billetera
type ViewWalletDetailsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewViewWalletDetailsUseCase crea una nueva instancia
func NewViewWalletDetailsUseCase(db *gorm.DB, log *logger.Logger) *ViewWalletDetailsUseCase {
	return &ViewWalletDetailsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ViewWalletDetailsUseCase) Execute(ctx context.Context, input *ViewWalletDetailsInput) (*ViewWalletDetailsOutput, error) {
	var result struct {
		domain.Wallet
		UserEmail string
		UserName  string
	}

	err := uc.db.WithContext(ctx).Table("wallets").
		Select("wallets.*, users.email as user_email, CONCAT(users.first_name, ' ', users.last_name) as user_name").
		Joins("LEFT JOIN users ON users.id = wallets.user_id").
		Where("wallets.id = ?", input.WalletID).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("billetera no encontrada")
		}
		uc.log.Error("Error getting wallet details", logger.Error(err))
		return nil, err
	}

	totalBalance := result.BalanceAvailable.Add(result.EarningsBalance)

	details := WalletDetails{
		ID:               result.ID,
		UUID:             result.UUID,
		UserID:           result.UserID,
		UserEmail:        result.UserEmail,
		UserName:         result.UserName,
		BalanceAvailable: result.BalanceAvailable.String(),
		EarningsBalance:  result.EarningsBalance.String(),
		PendingBalance:   result.PendingBalance.String(),
		TotalBalance:     totalBalance.String(),
		Currency:         result.Currency,
		Status:           string(result.Status),
		CreatedAt:        result.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        result.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return &ViewWalletDetailsOutput{
		Wallet: details,
	}, nil
}
