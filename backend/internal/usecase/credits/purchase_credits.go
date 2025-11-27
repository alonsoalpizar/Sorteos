package credits

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/internal/infrastructure/pagadito"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// PurchaseCreditsInput datos de entrada para compra de créditos
type PurchaseCreditsInput struct {
	UserID         int64           `json:"user_id" binding:"required"`
	DesiredCredit  decimal.Decimal `json:"desired_credit" binding:"required"` // Crédito que el usuario quiere
	Currency       string          `json:"currency" binding:"required"`       // CRC o USD
	IdempotencyKey string          `json:"idempotency_key" binding:"required"`
}

// PurchaseCreditsOutput datos de salida
type PurchaseCreditsOutput struct {
	Purchase   *domain.CreditPurchase `json:"purchase"`
	PaymentURL string                 `json:"payment_url"` // URL para redirigir al usuario
}

// PurchaseCreditsUseCase maneja la compra de créditos vía Pagadito
type PurchaseCreditsUseCase struct {
	purchaseRepo  domain.CreditPurchaseRepository
	walletRepo    domain.WalletRepository
	userRepo      domain.UserRepository
	auditRepo     domain.AuditLogRepository
	pagaditoClient pagadito.Client
	logger        *logger.Logger
}

// NewPurchaseCreditsUseCase crea una nueva instancia
func NewPurchaseCreditsUseCase(
	purchaseRepo domain.CreditPurchaseRepository,
	walletRepo domain.WalletRepository,
	userRepo domain.UserRepository,
	auditRepo domain.AuditLogRepository,
	pagaditoClient pagadito.Client,
	logger *logger.Logger,
) *PurchaseCreditsUseCase {
	return &PurchaseCreditsUseCase{
		purchaseRepo:   purchaseRepo,
		walletRepo:     walletRepo,
		userRepo:       userRepo,
		auditRepo:      auditRepo,
		pagaditoClient: pagaditoClient,
		logger:         logger,
	}
}

// Execute ejecuta el caso de uso de compra de créditos
func (uc *PurchaseCreditsUseCase) Execute(ctx context.Context, input *PurchaseCreditsInput) (*PurchaseCreditsOutput, error) {
	// 1. Validar monto
	if input.DesiredCredit.LessThanOrEqual(decimal.Zero) {
		return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "el monto debe ser mayor a cero", nil)
	}

	// Validar mínimo (₡1,000 o $2)
	minAmount := decimal.NewFromInt(1000) // ₡1,000
	if input.Currency == "USD" {
		minAmount = decimal.NewFromInt(2) // $2
	}
	if input.DesiredCredit.LessThan(minAmount) {
		return nil, errors.WrapWithMessage(
			errors.ErrValidationFailed,
			fmt.Sprintf("el monto mínimo es %s %s", minAmount.String(), input.Currency),
			nil,
		)
	}

	// Validar máximo (₡100,000 o $200)
	maxAmount := decimal.NewFromInt(100000) // ₡100,000
	if input.Currency == "USD" {
		maxAmount = decimal.NewFromInt(200) // $200
	}
	if input.DesiredCredit.GreaterThan(maxAmount) {
		return nil, errors.WrapWithMessage(
			errors.ErrValidationFailed,
			fmt.Sprintf("el monto máximo es %s %s", maxAmount.String(), input.Currency),
			nil,
		)
	}

	// 2. Verificar idempotencia
	existingPurchase, err := uc.purchaseRepo.FindByIdempotencyKey(input.IdempotencyKey)
	if err != nil && err != errors.ErrNotFound {
		uc.logger.Error("Error verificando idempotencia",
			logger.String("idempotency_key", input.IdempotencyKey),
			logger.Error(err))
		return nil, err
	}
	if existingPurchase != nil {
		// Compra duplicada detectada - retornar la existente
		uc.logger.Warn("Compra duplicada detectada (idempotencia)",
			logger.String("idempotency_key", input.IdempotencyKey),
			logger.Int64("existing_purchase_id", existingPurchase.ID))

		// Construir payment URL si aún está pendiente/processing
		paymentURL := ""
		if existingPurchase.IsPending() || existingPurchase.IsProcessing() {
			// TODO: Reconstruir URL de Pagadito o devolverla si está guardada
			paymentURL = "https://pagadito.com/..."
		}

		return &PurchaseCreditsOutput{
			Purchase:   existingPurchase,
			PaymentURL: paymentURL,
		}, nil
	}

	// 3. Validar que el usuario exista
	_, err = uc.userRepo.FindByID(input.UserID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "usuario no encontrado", err)
		}
		uc.logger.Error("Error buscando usuario",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}

	// 4. Obtener billetera del usuario
	wallet, err := uc.walletRepo.FindByUserID(input.UserID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.WrapWithMessage(errors.ErrValidationFailed, "billetera no encontrada", err)
		}
		uc.logger.Error("Error obteniendo billetera",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}

	// Validar estado de la billetera
	if wallet.Status != domain.WalletStatusActive {
		return nil, errors.WrapWithMessage(
			errors.ErrValidationFailed,
			fmt.Sprintf("billetera no activa (estado: %s)", wallet.Status),
			nil,
		)
	}

	// 5. Calcular comisiones usando RechargeCalculator existente
	// TODO: Cargar estos valores desde system_parameters
	fixedFee := decimal.NewFromInt(200)        // ₡200 fijo
	processorRate := decimal.NewFromFloat(0.035) // 3.5% procesador
	platformFeeRate := decimal.NewFromFloat(0.02) // 2% plataforma

	calculator := domain.NewRechargeCalculator(fixedFee, processorRate, platformFeeRate)
	breakdown := calculator.CalculateCharge(input.DesiredCredit)

	// 6. Generar ERN (External Reference Number)
	ern, err := domain.GenerateERN(input.UserID)
	if err != nil {
		uc.logger.Error("Error generando ERN",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, errors.Wrap(errors.ErrInternalServer, err)
	}

	// 7. Crear registro de compra en DB (estado: pending)
	purchase := &domain.CreditPurchase{
		UserID:         input.UserID,
		WalletID:       wallet.ID,
		DesiredCredit:  input.DesiredCredit,
		ChargeAmount:   breakdown.ChargeAmount,
		Currency:       input.Currency,
		FixedFee:       breakdown.FixedFee,
		ProcessorFee:   breakdown.ProcessorFee,
		PlatformFee:    breakdown.PlatformFee,
		ERN:            ern,
		Status:         domain.CreditPurchaseStatusPending,
		IdempotencyKey: input.IdempotencyKey,
		ExpiresAt:      time.Now().Add(30 * time.Minute), // TTL de 30 minutos
	}

	if err := uc.purchaseRepo.Create(purchase); err != nil {
		uc.logger.Error("Error creando compra en DB",
			logger.Int64("user_id", input.UserID),
			logger.Error(err))
		return nil, err
	}

	// 8. Conectar con Pagadito
	if err := uc.pagaditoClient.Connect(); err != nil {
		uc.logger.Error("Error conectando con Pagadito",
			logger.Error(err))

		// Marcar compra como fallida
		purchase.MarkAsFailed("Error conectando con procesador de pagos", "")
		uc.purchaseRepo.Update(purchase)

		return nil, errors.WrapWithMessage(
			errors.ErrStripeError, // Reutilizando error de procesador de pagos
			"servicio de pago temporalmente no disponible",
			err,
		)
	}

	// 9. Crear transacción en Pagadito
	// NOTA CRÍTICA: TODO DEBE IR EN USD
	// - price: USD
	// - amount: USD
	// - currency: "USD" (NO "CRC")
	//
	// Pagadito solo acepta USD. Si el usuario paga en CRC,
	// convertimos el monto a USD antes de enviar.

	// Convertir a USD si la moneda es CRC
	priceInUSD := breakdown.ChargeAmount
	amountInUSD := breakdown.ChargeAmount
	if input.Currency == "CRC" {
		// Convertir de CRC a USD usando tasa aproximada
		exchangeRate := decimal.NewFromInt(500) // 1 USD = 500 CRC
		priceInUSD = breakdown.ChargeAmount.Div(exchangeRate)
		amountInUSD = priceInUSD // amount = quantity (1) * price
	}

	transactionReq := &pagadito.TransactionRequest{
		ERN:      ern,
		Amount:   amountInUSD, // SIEMPRE en USD
		Currency: "USD",        // SIEMPRE "USD"
		Details: []pagadito.TransactionDetail{
			{
				Quantity:    1,
				Description: fmt.Sprintf("Recarga de créditos %s %s", input.DesiredCredit.String(), input.Currency),
				Price:       priceInUSD, // SIEMPRE en USD
				URLProduct:  "https://example.com/product",
			},
		},
		CustomParams: map[string]string{
			"user_id":     fmt.Sprintf("%d", input.UserID),
			"wallet_id":   fmt.Sprintf("%d", wallet.ID),
			"purchase_id": fmt.Sprintf("%d", purchase.ID),
		},
		AllowPending: true,
	}

	transactionResp, err := uc.pagaditoClient.CreateTransaction(transactionReq)
	if err != nil {
		uc.logger.Error("Error creando transacción en Pagadito",
			logger.String("ern", ern),
			logger.Error(err))

		// Marcar compra como fallida
		purchase.MarkAsFailed(fmt.Sprintf("Error del procesador: %v", err), "")
		uc.purchaseRepo.Update(purchase)

		return nil, errors.WrapWithMessage(
			errors.ErrStripeError, // Reutilizando error de procesador de pagos
			"error creando transacción de pago",
			err,
		)
	}

	// 10. Actualizar compra con token de Pagadito y cambiar estado a processing
	purchase.MarkAsProcessing(transactionResp.Token)
	if err := uc.purchaseRepo.Update(purchase); err != nil {
		uc.logger.Error("Error actualizando compra con token de Pagadito",
			logger.Int64("purchase_id", purchase.ID),
			logger.Error(err))
		// No retornar error - la transacción ya está creada en Pagadito
	}

	// 11. Log de auditoría
	entityType := "credit_purchase"
	metadataBytes, _ := json.Marshal(map[string]interface{}{
		"ern":            ern,
		"desired_credit": input.DesiredCredit.String(),
		"charge_amount":  breakdown.ChargeAmount.String(),
		"currency":       input.Currency,
	})
	uc.auditRepo.Create(&domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "credit_purchase_initiated",
		EntityType: &entityType,
		EntityID:   &purchase.ID,
		Metadata:   metadataBytes,
	})

	uc.logger.Info("Compra de créditos iniciada exitosamente",
		logger.Int64("purchase_id", purchase.ID),
		logger.Int64("user_id", input.UserID),
		logger.String("ern", ern),
		logger.String("payment_url", transactionResp.PaymentURL))

	// 12. Retornar resultado
	return &PurchaseCreditsOutput{
		Purchase:   purchase,
		PaymentURL: transactionResp.PaymentURL,
	}, nil
}
