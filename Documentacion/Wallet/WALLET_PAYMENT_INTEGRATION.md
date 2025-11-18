# Integración con Procesadores de Pago - Sistema de Billetera

## 🏦 Procesadores Locales de Costa Rica

El sistema de billetera está diseñado para ser **agnóstico del procesador de pagos**, permitiendo integrar con cualquier proveedor local de Costa Rica.

### Procesadores Comunes en CR:
- **BAC Credomatic** - Pasarela de pagos (tarjetas Visa/Mastercard)
- **Banco de Costa Rica (BCR)** - Pagos en línea
- **Davivienda** - Gateway de pagos
- **SINPE Móvil** - Transferencias instantáneas entre cuentas
- **Otros**: Scotiabank, Promerica, Popular, etc.

---

## 💰 Moneda: Colón Costarricense (CRC)

**Moneda por defecto:** `CRC` (₡)

**Conversión aproximada:**
- ₡1,000 CRC ≈ $2 USD
- ₡5,000 CRC ≈ $10 USD
- ₡500,000 CRC ≈ $1,000 USD

**Límites configurados:**
- **Mínimo:** ₡5,000 CRC (~$10 USD)
- **Máximo:** ₡5,000,000 CRC (~$10,000 USD)

---

## 🔌 Arquitectura de Integración

### Patrón: Strategy + Adapter

```
┌─────────────────────────────────────────────────┐
│           AddFundsUseCase (Core)                │
│    (Agnóstico del procesador de pagos)          │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
         ┌─────────────────────┐
         │ PaymentProvider     │ ◄─── Interface (Domain)
         │  (Interface)        │
         └─────────────────────┘
                   △
                   │ Implementaciones (Adapters)
         ┌─────────┴──────────┬──────────────┐
         │                    │              │
    ┌────▼────────┐    ┌─────▼──────┐  ┌───▼──────────┐
    │ BAC Adapter │    │ BCR Adapter│  │ SINPE Adapter│
    └─────────────┘    └────────────┘  └──────────────┘
```

---

## 📋 Interface PaymentProvider (Domain)

**Ubicación:** `internal/domain/payment_provider.go` (a crear)

```go
package domain

import (
	"context"
	"github.com/shopspring/decimal"
)

// PaymentProvider define el contrato para procesadores de pago
type PaymentProvider interface {
	// CreatePayment crea una nueva intención de pago
	CreatePayment(ctx context.Context, input CreatePaymentInput) (*PaymentOutput, error)

	// VerifyPayment verifica el estado de un pago
	VerifyPayment(ctx context.Context, paymentID string) (*PaymentStatus, error)

	// ProcessWebhook procesa notificaciones del procesador
	ProcessWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error)

	// GetName retorna el nombre del procesador
	GetName() string
}

// CreatePaymentInput datos para crear un pago
type CreatePaymentInput struct {
	Amount         decimal.Decimal
	Currency       string // "CRC"
	Description    string
	IdempotencyKey string
	UserID         int64
	UserEmail      string
	CallbackURL    string // URL de retorno después del pago
	WebhookURL     string // URL para notificaciones del procesador
	Metadata       map[string]interface{}
}

// PaymentOutput resultado de crear un pago
type PaymentOutput struct {
	PaymentID   string // ID del procesador
	Status      string // "pending", "processing", "completed", "failed"
	PaymentURL  string // URL para redirigir al usuario
	ExpiresAt   *time.Time
}

// PaymentStatus estado de un pago
type PaymentStatus struct {
	PaymentID   string
	Status      string
	Amount      decimal.Decimal
	PaidAt      *time.Time
	ErrorMessage *string
}

// WebhookEvent evento de webhook
type WebhookEvent struct {
	PaymentID    string
	Status       string
	Amount       decimal.Decimal
	TransactionID string
	PaidAt       *time.Time
}
```

---

## 🔨 Implementación: Adapter BAC Credomatic (Ejemplo)

**Ubicación:** `internal/adapters/payment/bac_payment_provider.go`

```go
package payment

import (
	"context"
	"fmt"
	"github.com/sorteos-platform/backend/internal/domain"
	"github.com/sorteos-platform/backend/pkg/logger"
	"github.com/shopspring/decimal"
)

// BACPaymentProvider implementa PaymentProvider para BAC Credomatic
type BACPaymentProvider struct {
	apiKey    string
	apiSecret string
	baseURL   string
	logger    *logger.Logger
}

// NewBACPaymentProvider crea una instancia del proveedor BAC
func NewBACPaymentProvider(apiKey, apiSecret, baseURL string, log *logger.Logger) *BACPaymentProvider {
	return &BACPaymentProvider{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
		logger:    log,
	}
}

// CreatePayment crea un pago en BAC
func (p *BACPaymentProvider) CreatePayment(ctx context.Context, input domain.CreatePaymentInput) (*domain.PaymentOutput, error) {
	// 1. Validar monto en CRC
	if input.Currency != "CRC" {
		return nil, fmt.Errorf("BAC solo acepta CRC")
	}

	// 2. Llamar API de BAC para crear sesión de pago
	bacRequest := map[string]interface{}{
		"amount":       input.Amount.String(),
		"currency":     "CRC",
		"description":  input.Description,
		"reference":    input.IdempotencyKey,
		"customer_email": input.UserEmail,
		"callback_url": input.CallbackURL,
		"webhook_url":  input.WebhookURL,
	}

	// HTTP POST a BAC API (ejemplo simplificado)
	// response := httpClient.Post(p.baseURL + "/payments", bacRequest)

	// 3. Retornar URL de pago de BAC
	return &domain.PaymentOutput{
		PaymentID:  "BAC-123456", // ID retornado por BAC
		Status:     "pending",
		PaymentURL: "https://bac.net/pay/xyz123", // URL de BAC para pagar
		ExpiresAt:  expirationTime,
	}, nil
}

// VerifyPayment verifica estado del pago en BAC
func (p *BACPaymentProvider) VerifyPayment(ctx context.Context, paymentID string) (*domain.PaymentStatus, error) {
	// GET a BAC API: /payments/{paymentID}
	// ...
	return &domain.PaymentStatus{
		PaymentID: paymentID,
		Status:    "completed",
		Amount:    decimal.NewFromFloat(5000.00),
		PaidAt:    &now,
	}, nil
}

// ProcessWebhook procesa webhook de BAC
func (p *BACPaymentProvider) ProcessWebhook(ctx context.Context, payload []byte, signature string) (*domain.WebhookEvent, error) {
	// 1. Verificar firma de BAC
	if !p.verifySignature(payload, signature) {
		return nil, fmt.Errorf("firma inválida")
	}

	// 2. Parsear payload de BAC
	var bacEvent BACWebhookPayload
	json.Unmarshal(payload, &bacEvent)

	// 3. Convertir a formato estándar
	return &domain.WebhookEvent{
		PaymentID:     bacEvent.PaymentID,
		Status:        bacEvent.Status,
		Amount:        decimal.NewFromString(bacEvent.Amount),
		TransactionID: bacEvent.TransactionID,
		PaidAt:        bacEvent.PaidAt,
	}, nil
}

func (p *BACPaymentProvider) GetName() string {
	return "BAC Credomatic"
}
```

---

## 🔄 Flujo Completo de Compra de Créditos

```
┌─────────┐                 ┌──────────┐              ┌───────────────┐
│ Usuario │                 │ Backend  │              │ BAC/Procesador│
└────┬────┘                 └────┬─────┘              └───────┬───────┘
     │                           │                            │
     │ 1. POST /wallet/add-funds │                            │
     ├──────────────────────────►│                            │
     │   (amount: 50000 CRC)      │                            │
     │                           │                            │
     │                           │ 2. AddFundsUseCase         │
     │                           │    crea tx PENDING         │
     │                           │                            │
     │                           │ 3. CreatePayment()         │
     │                           ├───────────────────────────►│
     │                           │                            │
     │                           │ 4. payment_url             │
     │                           │◄───────────────────────────┤
     │                           │                            │
     │ 5. payment_url            │                            │
     │◄──────────────────────────┤                            │
     │                           │                            │
     │ 6. Redirigir a BAC        │                            │
     ├────────────────────────────────────────────────────────►│
     │                           │                            │
     │                           │                       7. Usuario paga
     │                           │                          con tarjeta
     │                           │                            │
     │                           │ 8. POST /webhook (async)   │
     │                           │◄───────────────────────────┤
     │                           │                            │
     │                           │ 9. ConfirmAddFunds()       │
     │                           │    - Valida firma          │
     │                           │    - Acredita ₡50,000      │
     │                           │    - tx → COMPLETED        │
     │                           │                            │
     │ 10. Redirect callback_url │                            │
     │◄──────────────────────────────────────────────────────┤
     │                           │                            │
     │ 11. GET /wallet/balance   │                            │
     ├──────────────────────────►│                            │
     │                           │                            │
     │ 12. balance: ₡50,000      │                            │
     │◄──────────────────────────┤                            │
     │                           │                            │
```

---

## 🛠️ Configuración (Environment Variables)

```bash
# Procesador de pagos activo
PAYMENT_PROVIDER=bac  # "bac", "bcr", "sinpe", etc.

# BAC Credomatic
BAC_API_KEY=pk_live_xxxxxxxxxxxxx
BAC_API_SECRET=sk_live_xxxxxxxxxxxxx
BAC_BASE_URL=https://api.bac.net/v1
BAC_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxx

# URLs de callback
PAYMENT_CALLBACK_URL=https://sorteos.club/wallet/payment/success
PAYMENT_WEBHOOK_URL=https://api.sorteos.club/api/v1/wallet/webhook/bac

# Límites (CRC)
WALLET_MIN_DEPOSIT=5000     # ₡5,000 mínimo
WALLET_MAX_DEPOSIT=5000000  # ₡5,000,000 máximo
```

---

## 📝 Modificaciones en AddFundsUseCase

**Ubicación:** `internal/usecase/wallet/add_funds.go`

```go
type AddFundsUseCase struct {
	walletRepo      domain.WalletRepository
	transactionRepo domain.WalletTransactionRepository
	paymentProvider domain.PaymentProvider  // ← Inyectar provider
	userRepo        domain.UserRepository
	auditRepo       domain.AuditLogRepository
	logger          *logger.Logger
}

func (uc *AddFundsUseCase) Execute(ctx context.Context, input *AddFundsInput) (*AddFundsOutput, error) {
	// ... validaciones ...

	// Crear transacción PENDING en DB
	transaction := &domain.WalletTransaction{...}
	uc.transactionRepo.Create(transaction)

	// Llamar al procesador de pagos
	paymentOutput, err := uc.paymentProvider.CreatePayment(ctx, domain.CreatePaymentInput{
		Amount:         input.Amount,
		Currency:       "CRC",
		Description:    fmt.Sprintf("Recarga de billetera - Tx #%d", transaction.ID),
		IdempotencyKey: input.IdempotencyKey,
		UserEmail:      user.Email,
		CallbackURL:    cfg.PaymentCallbackURL,
		WebhookURL:     cfg.PaymentWebhookURL,
	})

	// Retornar URL de pago al frontend
	return &AddFundsOutput{
		Transaction:    transaction,
		NewBalance:     wallet.Balance,
		PaymentID:      &paymentOutput.PaymentID,
		PaymentURL:     &paymentOutput.PaymentURL,  // ← Frontend redirige aquí
	}, nil
}
```

---

## 🎣 Webhook Handler

**Ubicación:** `internal/adapters/http/handler/wallet/webhook_handler.go` (a crear)

```go
func (h *WebhookHandler) HandleBAC(c *gin.Context) {
	// 1. Leer payload y firma
	payload, _ := c.GetRawData()
	signature := c.GetHeader("X-BAC-Signature")

	// 2. Procesar webhook con el provider
	event, err := h.paymentProvider.ProcessWebhook(ctx, payload, signature)
	if err != nil {
		h.logger.Error("Invalid webhook signature", logger.Error(err))
		c.JSON(400, gin.H{"error": "invalid signature"})
		return
	}

	// 3. Buscar transacción por payment_id (en metadata)
	tx, err := h.transactionRepo.FindByPaymentID(event.PaymentID)

	// 4. Confirmar depósito
	if event.Status == "completed" {
		err := h.addFundsUseCase.ConfirmAddFunds(ctx, tx.ID)
		if err != nil {
			h.logger.Error("Failed to confirm funds", logger.Error(err))
			c.JSON(500, gin.H{"error": "failed to confirm"})
			return
		}
	}

	// 5. Responder OK al procesador
	c.JSON(200, gin.H{"status": "ok"})
}
```

---

## 🔐 Seguridad del Webhook

### Verificación de Firma (Ejemplo BAC)

```go
func (p *BACPaymentProvider) verifySignature(payload []byte, signature string) bool {
	// 1. Calcular HMAC-SHA256 del payload
	mac := hmac.New(sha256.New, []byte(p.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// 2. Comparación constante-time (prevenir timing attacks)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

### IP Whitelisting (Opcional)

```go
var bacAllowedIPs = []string{
	"192.0.2.1",    // IP de BAC (ejemplo)
	"192.0.2.2",
}

func (m *WebhookMiddleware) ValidateIP(c *gin.Context) {
	clientIP := c.ClientIP()

	allowed := false
	for _, ip := range bacAllowedIPs {
		if ip == clientIP {
			allowed = true
			break
		}
	}

	if !allowed {
		c.AbortWithStatus(403)
		return
	}

	c.Next()
}
```

---

## 🧪 Testing con Procesador Mock

**Ubicación:** `internal/adapters/payment/mock_payment_provider.go`

```go
type MockPaymentProvider struct {
	logger *logger.Logger
}

func (p *MockPaymentProvider) CreatePayment(ctx context.Context, input domain.CreatePaymentInput) (*domain.PaymentOutput, error) {
	// Simular pago exitoso automáticamente (para testing)
	return &domain.PaymentOutput{
		PaymentID:  fmt.Sprintf("MOCK-%s", input.IdempotencyKey),
		Status:     "pending",
		PaymentURL: "http://localhost:3000/mock-payment-success",
		ExpiresAt:  expiresAt,
	}, nil
}
```

**Configuración en desarrollo:**
```bash
PAYMENT_PROVIDER=mock  # Usar mock en desarrollo
```

---

## 📋 Checklist de Integración

- [ ] Elegir procesador de pagos local (BAC, BCR, etc.)
- [ ] Obtener credenciales API (api_key, api_secret)
- [ ] Implementar adapter del procesador (`bac_payment_provider.go`)
- [ ] Configurar variables de entorno
- [ ] Crear webhook handler (`/api/v1/wallet/webhook/bac`)
- [ ] Implementar verificación de firma del webhook
- [ ] Testing en sandbox/ambiente de pruebas del procesador
- [ ] Configurar IPs permitidas (whitelist)
- [ ] Monitoreo de webhooks (alertas si fallan)
- [ ] Manejo de reintentos de webhook
- [ ] Testing en producción con monto mínimo

---

## 💡 Recomendaciones

1. **Comenzar con BAC Credomatic**: Es el procesador más popular en CR
2. **Usar SINPE Móvil como alternativa**: Para transferencias instantáneas
3. **Implementar múltiples procesadores**: Permitir al usuario elegir
4. **Rate limiting agresivo en webhooks**: Prevenir spam
5. **Logs detallados**: Todos los webhooks deben loguearse
6. **Retry mechanism**: Si webhook falla, reintentar (exponential backoff)
7. **Monitoring**: Alertar si webhooks no llegan en X minutos

---

**Versión**: 1.0
**Última actualización**: 2025-11-18
**Moneda**: CRC (₡)
**Procesador recomendado**: BAC Credomatic, BCR, o SINPE Móvil
