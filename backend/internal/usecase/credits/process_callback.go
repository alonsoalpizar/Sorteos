package credits

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/internal/infrastructure/pagadito"
	walletuc "github.com/sorteos-platform/backend/internal/usecase/wallet"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// ProcessCallbackInput datos de entrada para procesar callback
type ProcessCallbackInput struct {
	Token string `json:"token" binding:"required"` // Token/ERN de la transacción
}

// ProcessCallbackOutput datos de salida
type ProcessCallbackOutput struct {
	Purchase     *domain.CreditPurchase `json:"purchase"`
	Status       string                 `json:"status"`        // COMPLETED, FAILED, PENDING
	Message      string                 `json:"message"`
	RedirectURL  string                 `json:"redirect_url"` // A dónde redirigir al usuario
}

// ProcessPagaditoCallbackUseCase procesa el callback de Pagadito
type ProcessPagaditoCallbackUseCase struct {
	purchaseRepo    domain.CreditPurchaseRepository
	walletRepo      domain.WalletRepository
	transactionRepo domain.WalletTransactionRepository
	auditRepo       domain.AuditLogRepository
	pagaditoClient  pagadito.Client
	addFundsUC      *walletuc.AddFundsUseCase
	logger          *logger.Logger
}

// NewProcessPagaditoCallbackUseCase crea una nueva instancia
func NewProcessPagaditoCallbackUseCase(
	purchaseRepo domain.CreditPurchaseRepository,
	walletRepo domain.WalletRepository,
	transactionRepo domain.WalletTransactionRepository,
	auditRepo domain.AuditLogRepository,
	pagaditoClient pagadito.Client,
	addFundsUC *walletuc.AddFundsUseCase,
	logger *logger.Logger,
) *ProcessPagaditoCallbackUseCase {
	return &ProcessPagaditoCallbackUseCase{
		purchaseRepo:    purchaseRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		auditRepo:       auditRepo,
		pagaditoClient:  pagaditoClient,
		addFundsUC:      addFundsUC,
		logger:          logger,
	}
}

// Execute procesa el callback de Pagadito
func (uc *ProcessPagaditoCallbackUseCase) Execute(ctx context.Context, input *ProcessCallbackInput) (*ProcessCallbackOutput, error) {
	// 1. Buscar compra por ERN (el token es el ERN)
	purchase, err := uc.purchaseRepo.FindByERN(input.Token)
	if err != nil {
		if err == errors.ErrNotFound {
			uc.logger.Warn("Compra no encontrada para token de callback",
				logger.String("token", input.Token))
			return &ProcessCallbackOutput{
				Status:      "ERROR",
				Message:     "Transacción no encontrada",
				RedirectURL: "/credits/error?reason=not_found",
			}, nil
		}
		uc.logger.Error("Error buscando compra",
			logger.String("token", input.Token),
			logger.Error(err))
		return nil, err
	}

	// 2. Verificar si ya fue procesada (idempotencia)
	if purchase.IsCompleted() {
		uc.logger.Info("Compra ya procesada anteriormente",
			logger.Int64("purchase_id", purchase.ID),
			logger.String("status", string(purchase.Status)))

		return &ProcessCallbackOutput{
			Purchase:    purchase,
			Status:      "COMPLETED",
			Message:     "Créditos ya acreditados",
			RedirectURL: fmt.Sprintf("/credits/success?purchase_id=%s", purchase.UUID),
		}, nil
	}

	if purchase.IsFailed() {
		uc.logger.Info("Compra ya marcada como fallida",
			logger.Int64("purchase_id", purchase.ID))

		return &ProcessCallbackOutput{
			Purchase:    purchase,
			Status:      "FAILED",
			Message:     "Pago no completado",
			RedirectURL: fmt.Sprintf("/credits/failed?purchase_id=%s", purchase.UUID),
		}, nil
	}

	// 3. Conectar con Pagadito para verificar estado
	if err := uc.pagaditoClient.Connect(); err != nil {
		uc.logger.Error("Error conectando con Pagadito para verificar estado",
			logger.Error(err))
		return &ProcessCallbackOutput{
			Status:      "PENDING",
			Message:     "Verificando estado del pago...",
			RedirectURL: fmt.Sprintf("/credits/pending?purchase_id=%s", purchase.UUID),
		}, nil
	}

	// 4. Consultar estado en Pagadito
	statusResp, err := uc.pagaditoClient.GetStatus(input.Token)
	if err != nil {
		uc.logger.Error("Error consultando estado en Pagadito",
			logger.String("token", input.Token),
			logger.Error(err))
		return &ProcessCallbackOutput{
			Status:      "PENDING",
			Message:     "Error verificando pago",
			RedirectURL: fmt.Sprintf("/credits/pending?purchase_id=%s", purchase.UUID),
		}, nil
	}

	uc.logger.Info("Estado de Pagadito recibido",
		logger.String("token", input.Token),
		logger.String("status", statusResp.Status),
		logger.String("reference", statusResp.Reference))

	// 5. Procesar según estado de Pagadito
	switch statusResp.Status {
	case "COMPLETED":
		return uc.processCompleted(ctx, purchase, statusResp)

	case "VERIFYING":
		return uc.processVerifying(ctx, purchase, statusResp)

	case "REGISTERED", "FAILED", "REVOKED":
		return uc.processFailed(ctx, purchase, statusResp)

	default:
		uc.logger.Warn("Estado desconocido de Pagadito",
			logger.String("status", statusResp.Status))
		return &ProcessCallbackOutput{
			Purchase:    purchase,
			Status:      "PENDING",
			Message:     "Estado desconocido del pago",
			RedirectURL: fmt.Sprintf("/credits/pending?purchase_id=%s", purchase.UUID),
		}, nil
	}
}

// processCompleted procesa un pago completado
func (uc *ProcessPagaditoCallbackUseCase) processCompleted(
	ctx context.Context,
	purchase *domain.CreditPurchase,
	statusResp *pagadito.StatusResponse,
) (*ProcessCallbackOutput, error) {
	// Acreditar créditos a la billetera usando AddFundsUseCase
	addFundsInput := &walletuc.AddFundsInput{
		UserID:         purchase.UserID,
		Amount:         purchase.DesiredCredit,
		IdempotencyKey: fmt.Sprintf("cp_%d_%s", purchase.ID, purchase.ERN),
		PaymentMethod:  "pagadito",
		PaymentIntentID: &statusResp.Reference, // NAP de Pagadito
		Metadata: map[string]interface{}{
			"credit_purchase_id": purchase.ID,
			"ern":                purchase.ERN,
			"pagadito_reference": statusResp.Reference,
			"charge_amount":      purchase.ChargeAmount.String(),
		},
	}

	addFundsOutput, err := uc.addFundsUC.Execute(ctx, addFundsInput)
	if err != nil {
		uc.logger.Error("Error acreditando fondos a billetera",
			logger.Int64("purchase_id", purchase.ID),
			logger.Int64("user_id", purchase.UserID),
			logger.Error(err))

		// NO marcar compra como fallida - intentar de nuevo más tarde
		return &ProcessCallbackOutput{
			Purchase:    purchase,
			Status:      "PENDING",
			Message:     "Error procesando créditos, intentando nuevamente...",
			RedirectURL: fmt.Sprintf("/credits/pending?purchase_id=%s", purchase.UUID),
		}, nil
	}

	// Actualizar compra como completada
	purchase.MarkAsCompleted(statusResp.Reference, addFundsOutput.Transaction.ID)
	if err := uc.purchaseRepo.Update(purchase); err != nil {
		uc.logger.Error("Error actualizando compra como completada",
			logger.Int64("purchase_id", purchase.ID),
			logger.Error(err))
		// No retornar error - los créditos ya fueron acreditados
	}

	// Log de auditoría
	entityType := "credit_purchase"
	metadataBytes, _ := json.Marshal(map[string]interface{}{
		"ern":                   purchase.ERN,
		"pagadito_reference":    statusResp.Reference,
		"desired_credit":        purchase.DesiredCredit.String(),
		"wallet_transaction_id": addFundsOutput.Transaction.ID,
		"new_balance":           addFundsOutput.NewBalance.String(),
	})
	uc.auditRepo.Create(&domain.AuditLog{
		UserID:     &purchase.UserID,
		Action:     "credit_purchase_completed",
		EntityType: &entityType,
		EntityID:   &purchase.ID,
		Metadata:   metadataBytes,
	})

	uc.logger.Info("Compra de créditos completada exitosamente",
		logger.Int64("purchase_id", purchase.ID),
		logger.Int64("user_id", purchase.UserID),
		logger.String("reference", statusResp.Reference),
		logger.String("credits_added", purchase.DesiredCredit.String()))

	return &ProcessCallbackOutput{
		Purchase:    purchase,
		Status:      "COMPLETED",
		Message:     fmt.Sprintf("¡Créditos acreditados exitosamente! Nuevo saldo: %s", addFundsOutput.NewBalance.String()),
		RedirectURL: fmt.Sprintf("/credits/success?purchase_id=%s&amount=%s", purchase.UUID, purchase.DesiredCredit.String()),
	}, nil
}

// processVerifying procesa un pago en verificación
func (uc *ProcessPagaditoCallbackUseCase) processVerifying(
	ctx context.Context,
	purchase *domain.CreditPurchase,
	statusResp *pagadito.StatusResponse,
) (*ProcessCallbackOutput, error) {
	// Mantener estado processing
	purchase.Status = domain.CreditPurchaseStatusProcessing
	statusStr := string(domain.PagaditoStatusVerifying)
	purchase.PagaditoStatus = &statusStr
	uc.purchaseRepo.Update(purchase)

	uc.logger.Info("Pago en verificación manual",
		logger.Int64("purchase_id", purchase.ID),
		logger.String("reference", statusResp.Reference))

	return &ProcessCallbackOutput{
		Purchase:    purchase,
		Status:      "VERIFYING",
		Message:     "Su pago está en verificación. Será notificado cuando se complete.",
		RedirectURL: fmt.Sprintf("/credits/verifying?purchase_id=%s", purchase.UUID),
	}, nil
}

// processFailed procesa un pago fallido/cancelado
func (uc *ProcessPagaditoCallbackUseCase) processFailed(
	ctx context.Context,
	purchase *domain.CreditPurchase,
	statusResp *pagadito.StatusResponse,
) (*ProcessCallbackOutput, error) {
	// Marcar como fallida
	var reason string
	switch statusResp.Status {
	case "REGISTERED":
		reason = "Pago cancelado por el usuario"
	case "REVOKED":
		reason = "Pago rechazado por el procesador"
	case "FAILED":
		reason = "Error procesando el pago"
	default:
		reason = fmt.Sprintf("Pago no completado (estado: %s)", statusResp.Status)
	}

	purchase.MarkAsFailed(reason, domain.PagaditoStatus(statusResp.Status))
	uc.purchaseRepo.Update(purchase)

	// Log de auditoría
	entityType := "credit_purchase"
	metadataBytes, _ := json.Marshal(map[string]interface{}{
		"ern":             purchase.ERN,
		"pagadito_status": statusResp.Status,
		"reason":          reason,
	})
	uc.auditRepo.Create(&domain.AuditLog{
		UserID:     &purchase.UserID,
		Action:     "credit_purchase_failed",
		EntityType: &entityType,
		EntityID:   &purchase.ID,
		Metadata:   metadataBytes,
	})

	uc.logger.Info("Compra de créditos fallida",
		logger.Int64("purchase_id", purchase.ID),
		logger.String("status", statusResp.Status),
		logger.String("reason", reason))

	return &ProcessCallbackOutput{
		Purchase:    purchase,
		Status:      "FAILED",
		Message:     reason,
		RedirectURL: fmt.Sprintf("/credits/failed?purchase_id=%s&reason=%s", purchase.UUID, statusResp.Status),
	}, nil
}
