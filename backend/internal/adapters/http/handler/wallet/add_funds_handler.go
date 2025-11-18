package wallet

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/sorteos-platform/backend/internal/usecase/wallet"
	"github.com/sorteos-platform/backend/pkg/errors"
	"github.com/sorteos-platform/backend/pkg/logger"
)

// AddFundsHandler maneja el endpoint de agregar fondos
type AddFundsHandler struct {
	useCase *wallet.AddFundsUseCase
	logger  *logger.Logger
}

// NewAddFundsHandler crea una nueva instancia del handler
func NewAddFundsHandler(useCase *wallet.AddFundsUseCase, logger *logger.Logger) *AddFundsHandler {
	return &AddFundsHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// AddFundsRequest representa la petición de agregar fondos
type AddFundsRequest struct {
	Amount        string `json:"amount" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"` // "card", "sinpe", "transfer"
}

// Handle maneja la petición de agregar fondos
// POST /api/v1/wallet/add-funds
func (h *AddFundsHandler) Handle(c *gin.Context) {
	// Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Usuario no autenticado",
		})
		return
	}

	var req AddFundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Datos de entrada inválidos: " + err.Error(),
		})
		return
	}

	// Parsear amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_AMOUNT",
			Message: "Monto inválido",
		})
		return
	}

	// Validar amount mínimo/máximo
	if amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_AMOUNT",
			Message: "El monto debe ser mayor a cero",
		})
		return
	}

	minAmount := decimal.NewFromFloat(1000.00) // Mínimo ₡1,000 CRC
	if amount.LessThan(minAmount) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "AMOUNT_TOO_LOW",
			Message: "El monto mínimo es ₡1,000",
		})
		return
	}

	maxAmount := decimal.NewFromFloat(5000000.00) // Máximo ₡5,000,000 CRC (~$10,000 USD)
	if amount.GreaterThan(maxAmount) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "AMOUNT_TOO_HIGH",
			Message: "El monto máximo es ₡5,000,000",
		})
		return
	}

	// Obtener o generar Idempotency-Key
	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		// Generar uno automáticamente si no se proporciona
		idempotencyKey = uuid.New().String()
		h.logger.Warn("Idempotency-Key not provided, generated automatically",
			logger.String("idempotency_key", idempotencyKey))
	}

	// Ejecutar caso de uso
	input := &wallet.AddFundsInput{
		UserID:         userID.(int64),
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
		PaymentMethod:  req.PaymentMethod,
	}

	output, err := h.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Respuesta exitosa
	// Nota: En producción, aquí se retornaría el payment_url o tokens
	// del procesador de pagos local (ej: BAC, BCR, Davivienda, etc.)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Transacción de depósito creada. Complete el pago con el procesador.",
		"data": gin.H{
			"transaction_id":   output.Transaction.ID,
			"transaction_uuid": output.Transaction.UUID,
			"amount":           output.Transaction.Amount.String(),
			"status":           output.Transaction.Status,
			"payment_method":   req.PaymentMethod,
			// TODO: Agregar payment_url o token del procesador local
			// Ejemplo: "payment_url": "https://procesador.local/pay/xyz123"
			"idempotency_key": idempotencyKey,
		},
	})
}

// handleError maneja los errores y retorna la respuesta apropiada
func (h *AddFundsHandler) handleError(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.logger.Error("Unexpected error in add funds handler", logger.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Error interno del servidor",
		})
		return
	}

	c.JSON(appErr.Status, ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
	})
}
