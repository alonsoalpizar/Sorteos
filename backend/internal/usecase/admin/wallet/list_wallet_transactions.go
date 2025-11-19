package wallet

import (
	"context"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// ListWalletTransactionsInput entrada
type ListWalletTransactionsInput struct {
	WalletID int64
	Page     int
	PageSize int
	Type     string // Filtrar por tipo de transacción
	Status   string // Filtrar por estado
}

// TransactionSummary resumen de transacción
type TransactionSummary struct {
	ID            int64  `json:"id"`
	UUID          string `json:"uuid"`
	Type          string `json:"type"`
	Amount        string `json:"amount"`
	Status        string `json:"status"`
	BalanceBefore string `json:"balance_before"`
	BalanceAfter  string `json:"balance_after"`
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   int64  `json:"reference_id,omitempty"`
	Notes         string `json:"notes,omitempty"`
	CreatedAt     string `json:"created_at"`
	CompletedAt   string `json:"completed_at,omitempty"`
}

// Pagination info
type TransactionPagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

// ListWalletTransactionsOutput salida
type ListWalletTransactionsOutput struct {
	Transactions []TransactionSummary  `json:"transactions"`
	Pagination   TransactionPagination `json:"pagination"`
}

// ListWalletTransactionsUseCase lista transacciones de una billetera
type ListWalletTransactionsUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListWalletTransactionsUseCase crea una nueva instancia
func NewListWalletTransactionsUseCase(db *gorm.DB, log *logger.Logger) *ListWalletTransactionsUseCase {
	return &ListWalletTransactionsUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListWalletTransactionsUseCase) Execute(ctx context.Context, input *ListWalletTransactionsInput) (*ListWalletTransactionsOutput, error) {
	// Validar que la wallet existe
	var wallet domain.Wallet
	if err := uc.db.WithContext(ctx).Where("id = ?", input.WalletID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("billetera no encontrada")
		}
		return nil, err
	}

	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Construir query
	query := uc.db.WithContext(ctx).Model(&domain.WalletTransaction{}).
		Where("wallet_id = ?", input.WalletID)

	// Filtrar por tipo
	if input.Type != "" {
		query = query.Where("type = ?", input.Type)
	}

	// Filtrar por estado
	if input.Status != "" {
		query = query.Where("status = ?", input.Status)
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		uc.log.Error("Error counting transactions", logger.Error(err))
		return nil, err
	}

	// Ordenar y paginar
	offset := (input.Page - 1) * input.PageSize
	var transactions []domain.WalletTransaction
	if err := query.Order("created_at DESC").
		Limit(input.PageSize).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		uc.log.Error("Error listing transactions", logger.Error(err))
		return nil, err
	}

	// Mapear a DTO
	summaries := make([]TransactionSummary, len(transactions))
	for i, tx := range transactions {
		summary := TransactionSummary{
			ID:            tx.ID,
			UUID:          tx.UUID,
			Type:          string(tx.Type),
			Amount:        tx.Amount.String(),
			Status:        string(tx.Status),
			BalanceBefore: tx.BalanceBefore.String(),
			BalanceAfter:  tx.BalanceAfter.String(),
			CreatedAt:     tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if tx.ReferenceType != nil {
			summary.ReferenceType = *tx.ReferenceType
		}
		if tx.ReferenceID != nil {
			summary.ReferenceID = *tx.ReferenceID
		}
		if tx.Notes != nil {
			summary.Notes = *tx.Notes
		}
		if tx.CompletedAt != nil {
			summary.CompletedAt = tx.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		summaries[i] = summary
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	return &ListWalletTransactionsOutput{
		Transactions: summaries,
		Pagination: TransactionPagination{
			Total:      total,
			Page:       input.Page - 1, // Convertir de 1-indexed a 0-indexed para el frontend
			Limit:      input.PageSize,
			TotalPages: totalPages,
		},
	}, nil
}
