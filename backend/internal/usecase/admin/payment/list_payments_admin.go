package payment

import (
	"context"
	"time"

	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// Payment estructura simplificada para admin (compatible con tabla payments que usa UUID)
type Payment struct {
	ID                    string     `json:"id"` // UUID
	ReservationID         string     `json:"reservation_id"`
	UserID                string     `json:"user_id"` // UUID reference
	RaffleID              string     `json:"raffle_id"` // UUID reference
	StripePaymentIntentID string     `json:"stripe_payment_intent_id"`
	StripeClientSecret    string     `json:"stripe_client_secret"`
	Amount                float64    `json:"amount"`
	Currency              string     `json:"currency"`
	Status                string     `json:"status"`
	PaymentMethod         *string    `json:"payment_method,omitempty"`
	ErrorMessage          *string    `json:"error_message,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	PaidAt                *time.Time `json:"paid_at,omitempty"`
	// Campos adicionales para admin (pueden no estar en la tabla actual)
	Provider     string     `json:"provider,omitempty"` // stripe, paypal, etc
	RefundedAt   *time.Time `json:"refunded_at,omitempty"`
	RefundedBy   *int64     `json:"refunded_by,omitempty"`
	AdminNotes   string     `json:"admin_notes,omitempty"`
}

// PaymentWithDetails pago con detalles adicionales
type PaymentWithDetails struct {
	Payment       *Payment `json:"payment"`
	UserName      string   `json:"user_name"`
	UserEmail     string   `json:"user_email"`
	RaffleTitle   string   `json:"raffle_title"`
	OrganizerName string   `json:"organizer_name"`
}

// ListPaymentsAdminInput datos de entrada
type ListPaymentsAdminInput struct {
	Page          int
	PageSize      int
	Status        *string // succeeded, pending, failed, refunded, cancelled
	UserID        *int64
	RaffleID      *int64
	OrganizerID   *int64
	Provider      *string // stripe, paypal, etc.
	DateFrom      *string
	DateTo        *string
	MinAmount     *float64
	MaxAmount     *float64
	Search        string // Buscar en payment_intent, order_id
	OrderBy       string
	IncludeRefund bool // Si true, incluye pagos refunded
}

// ListPaymentsAdminOutput resultado
type ListPaymentsAdminOutput struct {
	Payments       []*PaymentWithDetails
	Total          int64
	Page           int
	PageSize       int
	TotalPages     int
	TotalAmount    float64 // Suma total de los pagos filtrados
	SucceededCount int
	RefundedCount  int
	FailedCount    int
}

// ListPaymentsAdminUseCase caso de uso para listar pagos (admin)
type ListPaymentsAdminUseCase struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewListPaymentsAdminUseCase crea una nueva instancia
func NewListPaymentsAdminUseCase(db *gorm.DB, log *logger.Logger) *ListPaymentsAdminUseCase {
	return &ListPaymentsAdminUseCase{
		db:  db,
		log: log,
	}
}

// Execute ejecuta el caso de uso
func (uc *ListPaymentsAdminUseCase) Execute(ctx context.Context, input *ListPaymentsAdminInput, adminID int64) (*ListPaymentsAdminOutput, error) {
	// Validar paginación
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	// Calcular offset
	offset := (input.Page - 1) * input.PageSize

	// Construir query base con JOIN a users, raffles, y organizadores
	// NOTA: payments usa UUIDs para user_id y raffle_id
	query := uc.db.Table("payments").
		Select(`payments.*,
			COALESCE(users.first_name || ' ' || users.last_name, users.email) as user_name,
			users.email as user_email,
			raffles.title as raffle_title,
			COALESCE(organizers.first_name || ' ' || organizers.last_name, organizers.email) as organizer_name`).
		Joins("LEFT JOIN users ON users.uuid::text = payments.user_id").
		Joins("LEFT JOIN raffles ON raffles.uuid::text = payments.raffle_id").
		Joins("LEFT JOIN users AS organizers ON organizers.id = raffles.user_id")

	// Aplicar filtros
	if input.Status != nil {
		query = query.Where("payments.status = ?", *input.Status)
	} else if !input.IncludeRefund {
		// Por defecto, excluir refunded
		query = query.Where("payments.status != ?", "refunded")
	}

	if input.UserID != nil {
		// Convert int64 user ID to UUID
		query = query.Where("users.id = ?", *input.UserID)
	}

	if input.RaffleID != nil {
		// Convert int64 raffle ID to UUID
		query = query.Where("raffles.id = ?", *input.RaffleID)
	}

	if input.OrganizerID != nil {
		query = query.Where("raffles.user_id = ?", *input.OrganizerID)
	}

	if input.Provider != nil {
		query = query.Where("payments.provider = ?", *input.Provider)
	}

	if input.DateFrom != nil && *input.DateFrom != "" {
		query = query.Where("payments.created_at >= ?", *input.DateFrom)
	}

	if input.DateTo != nil && *input.DateTo != "" {
		query = query.Where("payments.created_at <= ?", *input.DateTo)
	}

	if input.MinAmount != nil {
		query = query.Where("payments.amount >= ?", *input.MinAmount)
	}

	if input.MaxAmount != nil {
		query = query.Where("payments.amount <= ?", *input.MaxAmount)
	}

	if input.Search != "" {
		searchPattern := "%" + input.Search + "%"
		query = query.Where(
			"payments.stripe_payment_intent ILIKE ? OR payments.paypal_order_id ILIKE ?",
			searchPattern, searchPattern)
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		uc.log.Error("Error counting payments", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Aplicar ordenamiento
	orderBy := "payments.created_at DESC"
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	}
	query = query.Order(orderBy)

	// Obtener payments con paginación
	var results []struct {
		Payment
		UserName      string
		UserEmail     string
		RaffleTitle   string
		OrganizerName string
	}

	if err := query.Offset(offset).Limit(input.PageSize).Scan(&results).Error; err != nil {
		uc.log.Error("Error listing payments", logger.Error(err))
		return nil, errors.Wrap(errors.ErrDatabaseError, err)
	}

	// Construir resultado
	payments := make([]*PaymentWithDetails, 0, len(results))
	totalAmount := 0.0
	succeededCount := 0
	refundedCount := 0
	failedCount := 0

	for _, result := range results {
		payments = append(payments, &PaymentWithDetails{
			Payment:       &result.Payment,
			UserName:      result.UserName,
			UserEmail:     result.UserEmail,
			RaffleTitle:   result.RaffleTitle,
			OrganizerName: result.OrganizerName,
		})

		// Calcular estadísticas
		if result.Status == "succeeded" {
			succeededCount++
			totalAmount += result.Amount
		} else if result.Status == "refunded" {
			refundedCount++
		} else if result.Status == "failed" {
			failedCount++
		}
	}

	// Calcular total de páginas
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	// Log auditoría
	uc.log.Info("Admin listed payments",
		logger.Int64("admin_id", adminID),
		logger.Int("total_results", len(payments)),
		logger.String("action", "admin_list_payments"))

	return &ListPaymentsAdminOutput{
		Payments:       payments,
		Total:          total,
		Page:           input.Page,
		PageSize:       input.PageSize,
		TotalPages:     totalPages,
		TotalAmount:    totalAmount,
		SucceededCount: succeededCount,
		RefundedCount:  refundedCount,
		FailedCount:    failedCount,
	}, nil
}
