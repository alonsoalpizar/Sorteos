# Pagos e Integraciones - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Criticidad:** MÁXIMA

---

## 1. Arquitectura de Pagos

### 1.1 Diseño Modular (Payment Provider Interface)

**Principio:** Arquitectura hexagonal con providers intercambiables.

```go
// internal/domain/payment_provider.go
type PaymentProvider interface {
    // Autorizar pago (pre-auth)
    Authorize(ctx context.Context, input AuthorizeInput) (*AuthorizeOutput, error)

    // Capturar pago previamente autorizado
    Capture(ctx context.Context, paymentID string) error

    // Reembolso total o parcial
    Refund(ctx context.Context, paymentID string, amount decimal.Decimal) error

    // Verificar firma de webhook
    VerifyWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error)

    // Obtener detalles de pago
    GetPayment(ctx context.Context, paymentID string) (*Payment, error)
}

type AuthorizeInput struct {
    Amount          decimal.Decimal
    Currency        string // "USD", "CRC"
    PaymentMethodID string // tok_xxx (Stripe), ba_xxx (PayPal)
    CustomerID      string
    Description     string
    Metadata        map[string]string
    IdempotencyKey  string
}

type AuthorizeOutput struct {
    PaymentID      string
    Status         PaymentStatus
    Amount         decimal.Decimal
    ExternalID     string // ID del PSP
    RequiresAction bool   // 3D Secure, etc.
    ActionURL      *string
}

type WebhookEvent struct {
    Type      string // "payment.succeeded", "payment.failed"
    PaymentID string
    Status    PaymentStatus
    Metadata  map[string]interface{}
}
```

---

### 1.2 Providers Soportados

| Provider | MVP | Fase 2 | Región | Fees |
|----------|-----|--------|--------|------|
| **Stripe** | ✅ | ✅ | Global | 2.9% + $0.30 |
| **PayPal** | ❌ | ✅ | Global | 3.4% + $0.30 |
| **Local CR** | ❌ | ✅ | Costa Rica | TBD |

**Selección de provider:**
- Feature flag por sorteo: `raffle.payment_provider`
- Fallback automático si provider principal falla
- Admin puede forzar provider en backoffice

---

## 2. Implementación: Stripe Provider

### 2.1 Setup

```go
import "github.com/stripe/stripe-go/v76"

type StripeProvider struct {
    client *stripe.Client
    config StripeConfig
}

type StripeConfig struct {
    SecretKey      string
    WebhookSecret  string
    PublishableKey string
}

func NewStripeProvider(config StripeConfig) *StripeProvider {
    stripe.Key = config.SecretKey
    return &StripeProvider{
        client: &stripe.Client{},
        config: config,
    }
}
```

---

### 2.2 Autorizar Pago

```go
func (p *StripeProvider) Authorize(ctx context.Context, input AuthorizeInput) (*AuthorizeOutput, error) {
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(int64(input.Amount.Mul(decimal.NewFromInt(100)).IntPart())), // cents
        Currency: stripe.String(strings.ToLower(input.Currency)),
        PaymentMethod: stripe.String(input.PaymentMethodID),
        Customer: stripe.String(input.CustomerID),
        Description: stripe.String(input.Description),
        Confirm: stripe.Bool(true), // Confirmar inmediatamente
    }

    // Metadata
    for k, v := range input.Metadata {
        params.AddMetadata(k, v)
    }

    // Idempotencia
    params.SetIdempotencyKey(input.IdempotencyKey)

    pi, err := paymentintent.New(params)
    if err != nil {
        return nil, fmt.Errorf("stripe authorize failed: %w", err)
    }

    output := &AuthorizeOutput{
        PaymentID:  pi.ID,
        ExternalID: pi.ID,
        Amount:     input.Amount,
        Status:     mapStripeStatus(pi.Status),
    }

    // Verificar si requiere autenticación adicional (3D Secure)
    if pi.Status == stripe.PaymentIntentStatusRequiresAction {
        output.RequiresAction = true
        output.ActionURL = &pi.NextAction.RedirectToURL.URL
    }

    return output, nil
}

func mapStripeStatus(status stripe.PaymentIntentStatus) PaymentStatus {
    switch status {
    case stripe.PaymentIntentStatusSucceeded:
        return PaymentSucceeded
    case stripe.PaymentIntentStatusRequiresPaymentMethod, stripe.PaymentIntentStatusRequiresConfirmation:
        return PaymentPending
    case stripe.PaymentIntentStatusCanceled:
        return PaymentFailed
    default:
        return PaymentPending
    }
}
```

---

### 2.3 Webhooks

**Endpoint:** `POST /webhooks/stripe`

```go
func HandleStripeWebhook(c *gin.Context) {
    payload, _ := ioutil.ReadAll(c.Request.Body)
    signature := c.GetHeader("Stripe-Signature")

    // Verificar firma
    event, err := webhook.ConstructEvent(payload, signature, config.StripeWebhookSecret)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid signature"})
        return
    }

    switch event.Type {
    case "payment_intent.succeeded":
        var pi stripe.PaymentIntent
        json.Unmarshal(event.Data.Raw, &pi)
        handlePaymentSucceeded(c.Request.Context(), pi.ID, pi.Metadata)

    case "payment_intent.payment_failed":
        var pi stripe.PaymentIntent
        json.Unmarshal(event.Data.Raw, &pi)
        handlePaymentFailed(c.Request.Context(), pi.ID, pi.LastPaymentError.Message)

    case "charge.dispute.created":
        // Chargeback
        handleChargeback(c.Request.Context(), event)
    }

    c.JSON(200, gin.H{"received": true})
}

func handlePaymentSucceeded(ctx context.Context, paymentID string, metadata map[string]string) {
    reservationID, _ := strconv.ParseInt(metadata["reservation_id"], 10, 64)

    db.Transaction(func(tx *gorm.DB) error {
        // Actualizar payment status
        tx.Model(&Payment{}).Where("external_id = ?", paymentID).
            Update("status", PaymentSucceeded)

        // Actualizar reservation
        tx.Model(&Reservation{}).Where("id = ?", reservationID).
            Update("status", ReservationConfirmed)

        // Marcar números como sold
        reservation := &Reservation{}
        tx.First(reservation, reservationID)
        tx.Model(&RaffleNumber{}).
            Where("raffle_id = ? AND number IN ?", reservation.RaffleID, reservation.Numbers).
            Updates(map[string]interface{}{
                "status": NumberSold,
                "sold_at": time.Now(),
            })

        // Enviar notificación
        notifier.SendEmail(ctx, user.Email, "payment_succeeded", map[string]any{
            "reservation_id": reservationID,
            "numbers": reservation.Numbers,
        })

        return nil
    })
}
```

**Seguridad de Webhooks:**
- Verificar firma HMAC-SHA256
- Retry logic del lado de Stripe (automático)
- Idempotencia: verificar que evento no fue procesado antes

---

### 2.4 Reembolsos

```go
func (p *StripeProvider) Refund(ctx context.Context, paymentID string, amount decimal.Decimal) error {
    params := &stripe.RefundParams{
        PaymentIntent: stripe.String(paymentID),
    }

    if !amount.IsZero() {
        params.Amount = stripe.Int64(int64(amount.Mul(decimal.NewFromInt(100)).IntPart()))
    }

    _, err := refund.New(params)
    return err
}
```

**Casos de uso:**
- Sorteo cancelado por owner
- Disputa ganada por comprador
- Error en pago (doble cargo)

---

## 3. Flujo Completo de Pago

### 3.1 Diagrama de Secuencia

```
Usuario          Frontend        Backend         Stripe          Redis           DB
  |                |                |               |               |             |
  |--Seleccionar-->|                |               |               |             |
  |   números      |                |               |               |             |
  |                |--POST /reservations---------->|               |             |
  |                |                |--Lock Redis------------------>|             |
  |                |                |<-Lock OK---------------------|             |
  |                |                |--Create Reservation---------------------->|
  |                |                |<-Reservation ID---------------------------|
  |                |<-reservation_id, client_secret|               |             |
  |                |                |               |               |             |
  |--Ingresar----->|                |               |               |             |
  |  tarjeta       |                |               |               |             |
  |                |--Stripe.js-------------------->|               |             |
  |                |                |               |--Authorize--->|             |
  |                |                |               |<-PaymentIntent|             |
  |                |<-Success-----------------------|               |             |
  |                |                |               |               |             |
  |                |                |<--Webhook (payment.succeeded)-|             |
  |                |                |--Confirm Reservation--------------------->|
  |                |                |--Mark Numbers Sold------------------------->|
  |                |                |--Release Lock--------------->|             |
  |<-Email confirmación-------------|               |               |             |
```

---

### 3.2 Código Frontend (React)

```tsx
import { loadStripe } from '@stripe/stripe-js'
import { Elements, CardElement, useStripe, useElements } from '@stripe/react-stripe-js'

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY)

function CheckoutForm({ reservationId, amount }: CheckoutFormProps) {
  const stripe = useStripe()
  const elements = useElements()
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!stripe || !elements) return

    setLoading(true)

    // 1. Crear Payment Intent en backend
    const { data } = await api.post('/payments', {
      reservation_id: reservationId,
      amount: amount,
    })

    // 2. Confirmar pago con Stripe.js
    const { error, paymentIntent } = await stripe.confirmCardPayment(
      data.client_secret,
      {
        payment_method: {
          card: elements.getElement(CardElement)!,
        },
      }
    )

    if (error) {
      toast.error(error.message)
      setLoading(false)
      return
    }

    if (paymentIntent.status === 'succeeded') {
      toast.success('Pago exitoso')
      navigate(`/confirmacion/${reservationId}`)
    }

    setLoading(false)
  }

  return (
    <form onSubmit={handleSubmit}>
      <CardElement options={{ style: { base: { fontSize: '16px' } } }} />
      <Button type="submit" disabled={!stripe || loading}>
        {loading ? 'Procesando...' : `Pagar $${amount}`}
      </Button>
    </form>
  )
}

function CheckoutPage() {
  return (
    <Elements stripe={stripePromise}>
      <CheckoutForm reservationId={123} amount={15} />
    </Elements>
  )
}
```

---

## 4. Idempotencia

### 4.1 ¿Por qué es crítica?

**Escenario sin idempotencia:**
1. Usuario hace click en "Pagar"
2. Request llega al servidor → cargo exitoso
3. Response se pierde (timeout de red)
4. Frontend reintenta → **doble cargo**

**Solución: Idempotency-Key**

---

### 4.2 Implementación

**Frontend:**
```tsx
import { v4 as uuidv4 } from 'uuid'

const idempotencyKey = useMemo(() => uuidv4(), [reservationId])

await api.post('/payments', {
  reservation_id: reservationId,
  amount: 15,
}, {
  headers: {
    'Idempotency-Key': idempotencyKey,
  },
})
```

**Backend:**
```go
func CreatePayment(c *gin.Context) {
    idempotencyKey := c.GetHeader("Idempotency-Key")
    if idempotencyKey == "" {
        c.JSON(400, gin.H{"error": "Idempotency-Key required"})
        return
    }

    // Verificar si ya existe en Redis
    cacheKey := fmt.Sprintf("payment:idempotency:%s", idempotencyKey)
    existingPaymentID, err := rdb.Get(ctx, cacheKey).Result()
    if err == nil {
        // Ya procesado, retornar pago existente
        payment := paymentRepo.FindByID(ctx, existingPaymentID)
        c.JSON(200, payment)
        return
    }

    // Crear nuevo pago
    payment := createPaymentWithProvider(...)

    // Guardar idempotency key en Redis (24h)
    rdb.Set(ctx, cacheKey, payment.ID, 24*time.Hour)

    c.JSON(201, payment)
}
```

---

## 5. Manejo de Errores y Reintentos

### 5.1 Estrategia de Reintentos

**Exponential Backoff:**
```go
func retryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }

        if !isRetryable(err) {
            return err
        }

        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    return errors.New("max retries exceeded")
}

func isRetryable(err error) bool {
    // Reintentar solo en errores de red/timeout, no en errores de validación
    stripeErr, ok := err.(*stripe.Error)
    if !ok {
        return false
    }

    return stripeErr.Type == stripe.ErrorTypeAPIConnection ||
           stripeErr.Code == stripe.ErrorCodeRateLimitExceeded
}
```

---

### 5.2 Códigos de Error

| Código | Descripción | Acción Frontend |
|--------|-------------|-----------------|
| `card_declined` | Tarjeta rechazada | Mostrar mensaje, pedir otra tarjeta |
| `insufficient_funds` | Fondos insuficientes | Pedir otra tarjeta |
| `expired_card` | Tarjeta expirada | Pedir otra tarjeta |
| `processing_error` | Error temporal | Reintentar automáticamente |
| `rate_limit_error` | Demasiadas requests | Esperar 30s, reintentar |

---

## 6. Conciliación de Pagos

### 6.1 Cron Job Diario

**Objetivo:** Verificar que pagos en DB coincidan con Stripe.

```go
func ReconcilePayments(ctx context.Context) error {
    // Obtener pagos de últimas 24h en estado pending
    payments := paymentRepo.FindPending(ctx, 24*time.Hour)

    for _, payment := range payments {
        // Consultar estado real en Stripe
        stripePayment, err := stripeProvider.GetPayment(ctx, payment.ExternalID)
        if err != nil {
            logger.Error("reconciliation failed", zap.Int64("payment_id", payment.ID))
            continue
        }

        // Actualizar si difiere
        if payment.Status != stripePayment.Status {
            logger.Warn("payment status mismatch",
                zap.Int64("payment_id", payment.ID),
                zap.String("db_status", string(payment.Status)),
                zap.String("stripe_status", string(stripePayment.Status)),
            )

            payment.Status = stripePayment.Status
            paymentRepo.Update(ctx, payment)

            // Si pasó a succeeded, confirmar reserva
            if stripePayment.Status == PaymentSucceeded {
                confirmReservation(ctx, payment.ReservationID)
            }
        }
    }

    return nil
}
```

**Ejecución:** Todos los días a las 2am (cron)

---

## 7. Chargebacks y Disputas

### 7.1 Webhook: `charge.dispute.created`

```go
func handleChargeback(ctx context.Context, event stripe.Event) {
    var dispute stripe.Dispute
    json.Unmarshal(event.Data.Raw, &dispute)

    // Marcar pago como disputed
    payment := paymentRepo.FindByExternalID(ctx, dispute.Charge.ID)
    payment.Status = PaymentDisputed
    paymentRepo.Update(ctx, payment)

    // Notificar al owner del sorteo
    raffle := raffleRepo.FindByID(ctx, payment.Reservation.RaffleID)
    notifier.SendEmail(ctx, raffle.Owner.Email, "chargeback_notification", map[string]any{
        "amount": payment.Amount,
        "reason": dispute.Reason,
    })

    // Bloquear fondos de liquidación
    settlement := settlementRepo.FindByRaffleID(ctx, raffle.ID)
    settlement.Status = SettlementOnHold
    settlementRepo.Update(ctx, settlement)
}
```

**Política:**
- Fondos se retienen hasta resolución de disputa
- Owner puede responder con evidencia (tracking, fotos)
- Si dispute se pierde → reembolso automático + números liberados

---

## 8. Modo "Sin Cobro en Plataforma"

### 8.1 Modelo de Negocio

**Opción 1: Suscripción Mensual**
- Owner paga $20/mes (plan Basic) o $50/mes (plan Pro)
- Puede publicar sorteos ilimitados
- No cobra comisión por venta de boletos
- Owner coordina pagos fuera de plataforma (efectivo, Sinpe, etc.)

**Opción 2: Gratuito con Limitaciones**
- Máximo 3 sorteos activos
- Sin comisión de plataforma
- Watermark en imágenes (opcional)

---

### 8.2 Implementación

```go
type Raffle struct {
    // ...
    ChargeOnPlatform bool // true = cobrar en plataforma, false = sin cobro
    SubscriptionTier string // "free", "basic", "pro"
}

// Validación al crear sorteo
func (uc *CreateRaffleUseCase) Execute(ctx context.Context, input CreateRaffleInput) error {
    user := userRepo.FindByID(ctx, input.UserID)

    if !input.ChargeOnPlatform {
        // Verificar límites según suscripción
        activeRaffles := raffleRepo.CountActive(ctx, user.ID)
        if user.SubscriptionTier == "free" && activeRaffles >= 3 {
            return errors.New("límite de sorteos alcanzado, actualiza a plan Basic")
        }
    }

    // ...
}
```

**Frontend:**
```tsx
<Checkbox
  checked={chargeOnPlatform}
  onCheckedChange={setChargeOnPlatform}
>
  Cobrar boletos en la plataforma
</Checkbox>

{!chargeOnPlatform && (
  <Alert>
    <InfoIcon className="w-4 h-4" />
    <AlertDescription>
      Coordinarás los pagos directamente con los compradores.
      La plataforma solo gestiona los sorteos.
    </AlertDescription>
  </Alert>
)}
```

---

## 9. Liquidaciones (Settlements)

### 9.1 Flujo de Liquidación

**Cuando un sorteo finaliza:**
1. Calcular total recaudado
2. Descontar comisiones:
   - Stripe: 2.9% + $0.30 por transacción
   - Plataforma: 5% del total
3. Calcular neto a depositar al owner
4. Esperar confirmación de entrega de premio
5. Depositar a cuenta del owner

```go
type Settlement struct {
    ID            int64
    RaffleID      int64
    UserID        int64 // owner
    GrossAmount   decimal.Decimal // total recaudado
    StripeFees    decimal.Decimal
    PlatformFee   decimal.Decimal
    NetAmount     decimal.Decimal // a depositar
    Status        SettlementStatus
    PaidAt        *time.Time
}

func CalculateSettlement(raffle *Raffle) *Settlement {
    grossAmount := raffle.SoldCount * raffle.PricePerNumber

    stripeFees := grossAmount.Mul(decimal.NewFromFloat(0.029)).
        Add(decimal.NewFromFloat(0.30).Mul(decimal.NewFromInt(raffle.SoldCount)))

    platformFee := grossAmount.Mul(decimal.NewFromFloat(0.05))

    netAmount := grossAmount.Sub(stripeFees).Sub(platformFee)

    return &Settlement{
        RaffleID:    raffle.ID,
        UserID:      raffle.UserID,
        GrossAmount: grossAmount,
        StripeFees:  stripeFees,
        PlatformFee: platformFee,
        NetAmount:   netAmount,
        Status:      SettlementPending,
    }
}
```

---

### 9.2 Payout a Owner (Stripe Connect)

**Futuro (Fase 2):**
- Integrar Stripe Connect
- Owner vincula su cuenta bancaria
- Payout automático tras confirmación de entrega

```go
func PayoutToOwner(ctx context.Context, settlementID int64) error {
    settlement := settlementRepo.FindByID(ctx, settlementID)
    user := userRepo.FindByID(ctx, settlement.UserID)

    // Crear payout en Stripe Connect
    params := &stripe.PayoutParams{
        Amount:   stripe.Int64(int64(settlement.NetAmount.Mul(decimal.NewFromInt(100)).IntPart())),
        Currency: stripe.String("usd"),
        Destination: stripe.String(user.StripeConnectAccountID),
    }

    payout, err := payout.New(params)
    if err != nil {
        return err
    }

    settlement.Status = SettlementPaid
    settlement.PaidAt = &time.Time{}
    settlementRepo.Update(ctx, settlement)

    return nil
}
```

---

## 10. Tests de Pagos

### 10.1 Tests con Stripe Test Mode

**Tarjetas de prueba:**
- `4242 4242 4242 4242` - Aprobada
- `4000 0000 0000 9995` - Rechazada (insufficient_funds)
- `4000 0025 0000 3155` - Requiere 3D Secure

**Tests críticos:**
- [ ] Pago exitoso actualiza reserva y números
- [ ] Mismo Idempotency-Key no duplica cargo
- [ ] Webhook tardío no confirma reserva expirada
- [ ] Reembolso libera números para reventa
- [ ] Conciliación detecta discrepancias

---

## 11. Próximos Pasos

1. Implementar StripeProvider completo
2. Setup webhook endpoint con verificación de firma
3. Crear flujo de checkout en frontend
4. Tests e2e con Stripe test mode
5. Integrar PayPal (Fase 2)
6. Stripe Connect para payouts automáticos

---

**Actualizado:** 2025-11-10
