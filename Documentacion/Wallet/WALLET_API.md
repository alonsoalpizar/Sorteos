# API de Billetera/Monedero - Documentación

## Base URL
```
https://api.sorteos.club/api/v1/wallet
```

Todas las rutas requieren autenticación excepto donde se indique lo contrario.

---

## 🔐 Autenticación

Todas las peticiones deben incluir el token JWT en el header:

```http
Authorization: Bearer <access_token>
```

El `user_id` se extrae automáticamente del token JWT.

---

## 📊 Endpoints

### 1. Calcular Opciones de Recarga

Obtiene las opciones predefinidas de recarga con sus desgloses completos (crédito deseado vs monto a cobrar).

**Endpoint:** `GET /api/v1/wallet/recharge-options`

**Autenticación:** No requerida (público)

**Rate Limit:** 60 requests / minuto por IP

**Respuesta Exitosa (200 OK):**
```json
{
  "success": true,
  "data": {
    "options": [
      {
        "desired_credit": "1000.00",
        "fixed_fee": "100.00",
        "processor_rate": "0.03",
        "processor_fee": "133.68",
        "platform_fee_rate": "0.02",
        "platform_fee": "22.11",
        "total_fees": "155.79",
        "charge_amount": "1155.79"
      },
      {
        "desired_credit": "5000.00",
        "fixed_fee": "100.00",
        "processor_rate": "0.03",
        "processor_fee": "268.42",
        "platform_fee_rate": "0.02",
        "platform_fee": "110.53",
        "total_fees": "378.95",
        "charge_amount": "5378.95"
      },
      {
        "desired_credit": "10000.00",
        "fixed_fee": "100.00",
        "processor_rate": "0.03",
        "processor_fee": "418.95",
        "platform_fee_rate": "0.02",
        "platform_fee": "221.05",
        "total_fees": "640.00",
        "charge_amount": "10640.00"
      },
      {
        "desired_credit": "15000.00",
        "fixed_fee": "100.00",
        "processor_rate": "0.03",
        "processor_fee": "569.47",
        "platform_fee_rate": "0.02",
        "platform_fee": "331.58",
        "total_fees": "901.05",
        "charge_amount": "15901.05"
      },
      {
        "desired_credit": "20000.00",
        "fixed_fee": "100.00",
        "processor_rate": "0.03",
        "processor_fee": "720.00",
        "platform_fee_rate": "0.02",
        "platform_fee": "442.11",
        "total_fees": "1162.11",
        "charge_amount": "21162.11"
      }
    ],
    "currency": "CRC",
    "note": "Los montos mostrados incluyen todas las comisiones. El crédito deseado es lo que recibirás en tu billetera."
  }
}
```

**Campos de cada opción:**
- `desired_credit`: Crédito que recibirá el usuario en su billetera (₡1,000, ₡5,000, ₡10,000, ₡15,000, ₡20,000)
- `fixed_fee`: Tarifa fija del procesador de pagos (₡100)
- `processor_rate`: Tasa porcentual del procesador (3% = 0.03)
- `processor_fee`: Comisión calculada del procesador
- `platform_fee_rate`: Tasa de comisión de la plataforma (2% = 0.02)
- `platform_fee`: Comisión calculada de la plataforma
- `total_fees`: Total de comisiones (processor_fee + platform_fee)
- `charge_amount`: **Monto total a cobrar al usuario** (desired_credit + total_fees)

**Fórmula utilizada:** `C = (D + f) / (1 - r)`
- C = Charge amount (monto a cobrar)
- D = Desired credit (crédito deseado)
- f = Fixed fee (tarifa fija)
- r = Total rate (processor_rate + platform_fee_rate)

**Uso recomendado:**
Este endpoint debe llamarse al mostrar la pantalla de recarga para que el usuario vea exactamente cuánto se le cobrará por cada opción de crédito.

---

### 2. Consultar Saldo

Obtiene el saldo actual de la billetera del usuario autenticado.

**Endpoint:** `GET /api/v1/wallet/balance`

**Headers:**
```http
Authorization: Bearer <access_token>
```

**Rate Limit:** 30 requests / minuto

**Respuesta Exitosa (200 OK):**
```json
{
  "success": true,
  "data": {
    "wallet_id": 123,
    "wallet_uuid": "550e8400-e29b-41d4-a716-446655440000",
    "balance": "150.50",
    "pending_balance": "25.00",
    "currency": "USD",
    "status": "active"
  }
}
```

**Campos de Respuesta:**
- `wallet_id`: ID interno de la billetera
- `wallet_uuid`: UUID público de la billetera
- `balance`: Saldo disponible para usar (string decimal)
- `pending_balance`: Saldo pendiente de confirmación (string decimal)
- `currency`: Moneda (ISO 4217: "USD", "CRC")
- `status`: Estado de la billetera ("active", "frozen", "closed")

**Errores Comunes:**
```json
// 401 Unauthorized
{
  "code": "UNAUTHORIZED",
  "message": "Usuario no autenticado"
}

// 404 Not Found
{
  "code": "NOT_FOUND",
  "message": "Billetera no encontrada"
}
```

---

### 2. Listar Transacciones

Obtiene el historial de transacciones de la billetera del usuario autenticado.

**Endpoint:** `GET /api/v1/wallet/transactions`

**Headers:**
```http
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `limit` (opcional): Número de transacciones por página (1-100, default: 20)
- `offset` (opcional): Número de transacciones a saltar (default: 0)

**Ejemplo:**
```
GET /api/v1/wallet/transactions?limit=20&offset=0
```

**Rate Limit:** 30 requests / minuto

**Respuesta Exitosa (200 OK):**
```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "id": 1,
        "uuid": "550e8400-e29b-41d4-a716-446655440001",
        "type": "deposit",
        "amount": "100.00",
        "status": "completed",
        "balance_before": "50.50",
        "balance_after": "150.50",
        "reference_type": "payment_intent",
        "reference_id": null,
        "notes": null,
        "created_at": "2025-11-18T10:30:00Z",
        "completed_at": "2025-11-18T10:30:15Z"
      },
      {
        "id": 2,
        "uuid": "550e8400-e29b-41d4-a716-446655440002",
        "type": "purchase",
        "amount": "25.00",
        "status": "completed",
        "balance_before": "150.50",
        "balance_after": "125.50",
        "reference_type": "raffle",
        "reference_id": 456,
        "notes": null,
        "created_at": "2025-11-18T11:00:00Z",
        "completed_at": "2025-11-18T11:00:01Z"
      }
    ],
    "pagination": {
      "total": 45,
      "limit": 20,
      "offset": 0
    }
  }
}
```

**Tipos de Transacción (`type`):**
- `deposit`: Compra de créditos vía procesador (Stripe)
- `withdrawal`: Retiro a cuenta bancaria
- `purchase`: Pago de sorteo con saldo
- `refund`: Devolución de compra
- `prize_claim`: Premio ganado (organizador)
- `settlement_payout`: Pago de liquidación a organizador
- `adjustment`: Ajuste manual por admin

**Estados de Transacción (`status`):**
- `pending`: Pendiente de confirmación
- `completed`: Completada exitosamente
- `failed`: Fallida
- `reversed`: Revertida

---

### 3. Agregar Fondos

Inicia el proceso de compra de créditos para la billetera.

**Endpoint:** `POST /api/v1/wallet/add-funds`

**Headers:**
```http
Authorization: Bearer <access_token>
Content-Type: application/json
Idempotency-Key: <uuid> (opcional pero recomendado)
```

**Body:**
```json
{
  "amount": "100.00",
  "payment_method": "stripe"
}
```

**Campos del Body:**
- `amount` (requerido): Monto de crédito deseado (string decimal)
  - Mínimo: ₡1,000 CRC
  - Máximo: ₡5,000,000 CRC
  - **Recomendado:** Usar las opciones predefinidas del endpoint `/recharge-options`
- `payment_method` (requerido): Método de pago ("card", "sinpe", "transfer")

**Idempotency-Key:**
- Si se proporciona, previene transacciones duplicadas
- Debe ser un UUID único por intento de pago
- Si no se proporciona, se genera automáticamente

**Rate Limit:** 5 requests / hora (más restrictivo para prevenir fraude)

**Respuesta Exitosa (201 Created):**
```json
{
  "success": true,
  "message": "Transacción de depósito creada. Complete el pago con el procesador.",
  "data": {
    "transaction_id": 123,
    "transaction_uuid": "550e8400-e29b-41d4-a716-446655440003",
    "amount": "100.00",
    "status": "pending",
    "payment_method": "stripe",
    "idempotency_key": "550e8400-e29b-41d4-a716-446655440004"
  }
}
```

**Flujo Completo:**
1. Cliente llama a `POST /add-funds`
2. Backend crea transacción PENDING
3. Cliente redirige a Stripe Checkout (TODO: agregar client_secret)
4. Usuario paga en Stripe
5. Webhook de Stripe llama a `/webhook/stripe`
6. Backend confirma transacción y acredita fondos
7. Transacción pasa a status COMPLETED

**Errores Comunes:**
```json
// 400 Bad Request - Monto inválido
{
  "code": "AMOUNT_TOO_LOW",
  "message": "El monto mínimo es ₡1,000"
}

{
  "code": "AMOUNT_TOO_HIGH",
  "message": "El monto máximo es ₡5,000,000"
}

// 409 Conflict - Transacción duplicada (idempotencia)
{
  "code": "CONFLICT",
  "message": "Transacción duplicada"
}

// 429 Too Many Requests - Rate limit excedido
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Demasiadas peticiones. Intente más tarde."
}
```

---

## 🔄 Flujos de Uso

### Flujo 1: Mostrar Opciones de Recarga

```javascript
// 1. Obtener opciones de recarga predefinidas (sin autenticación)
const optionsRes = await fetch('/api/v1/wallet/recharge-options')
const { data } = await optionsRes.json()

// 2. Mostrar opciones al usuario
data.options.forEach(option => {
  console.log(`Crédito: ₡${option.desired_credit}`)
  console.log(`Total a pagar: ₡${option.charge_amount}`)
  console.log(`Comisiones: ₡${option.total_fees}`)
  console.log('---')
})

// Ejemplo de salida:
// Crédito: ₡1,000
// Total a pagar: ₡1,155.79
// Comisiones: ₡155.79
// ---
// Crédito: ₡5,000
// Total a pagar: ₡5,378.95
// Comisiones: ₡378.95
```

### Flujo 2: Comprar Créditos

```javascript
// 1. Usuario selecciona una opción (ej: ₡5,000)
const selectedOption = data.options[1] // ₡5,000

// 2. Generar idempotency key (UNA sola vez)
const idempotencyKey = crypto.randomUUID()

// 3. Solicitar agregar fondos (con el crédito deseado, no el charge_amount)
const response = await fetch('/api/v1/wallet/add-funds', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json',
    'Idempotency-Key': idempotencyKey
  },
  body: JSON.stringify({
    amount: selectedOption.desired_credit, // '5000.00' - crédito deseado
    payment_method: 'card' // o 'sinpe', 'transfer'
  })
})

const { data } = await response.json()

// 4. Redirigir al procesador de pagos (BAC/BCR/SINPE)
// window.location.href = data.payment_url

// 5. Webhook del procesador confirmará automáticamente
// El usuario recibirá exactamente ₡5,000 en su billetera
```

### Flujo 3: Consultar Saldo antes de Comprar

```javascript
// 1. Consultar saldo actual
const balanceRes = await fetch('/api/v1/wallet/balance', {
  headers: {
    'Authorization': `Bearer ${accessToken}`
  }
})

const { data } = await balanceRes.json()

// 2. Verificar si tiene saldo suficiente
if (parseFloat(data.balance) >= rafflePrice) {
  // Puede pagar con saldo
  await payWithWallet(raffleId)
} else {
  // Necesita agregar fondos
  showAddFundsModal()
}
```

### Flujo 4: Mostrar Historial de Transacciones

```javascript
// 1. Cargar primera página
let offset = 0
const limit = 20

const response = await fetch(
  `/api/v1/wallet/transactions?limit=${limit}&offset=${offset}`,
  {
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  }
)

const { data } = await response.json()

// 2. Renderizar transacciones
data.transactions.forEach(tx => {
  console.log(`${tx.type}: ${tx.amount} (${tx.status})`)
})

// 3. Paginación
const hasNextPage = (offset + limit) < data.pagination.total
if (hasNextPage) {
  // Cargar siguiente página
  offset += limit
  // ... repetir fetch
}
```

---

## 🔒 Seguridad

### Idempotencia
Todas las operaciones de dinero soportan idempotencia mediante el header `Idempotency-Key`:

```http
Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000
```

**Reglas:**
- Debe ser un UUID v4
- Se debe generar UNA sola vez en el cliente
- NO regenerar en retries (usar el mismo UUID)
- Válido por 24 horas
- Previene transacciones duplicadas por doble click o retry

### Rate Limiting

| Endpoint | Límite | Por |
|----------|--------|-----|
| GET /balance | 30 req | minuto |
| GET /transactions | 30 req | minuto |
| POST /add-funds | 5 req | hora |

Cuando se excede el límite:
```json
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Demasiadas peticiones. Intente en X segundos."
}
```

### HTTPS Obligatorio
Todas las peticiones deben usar HTTPS en producción. HTTP será rechazado.

---

## 📝 Códigos de Estado HTTP

| Código | Significado |
|--------|-------------|
| 200 OK | Petición exitosa |
| 201 Created | Recurso creado (add-funds) |
| 400 Bad Request | Datos de entrada inválidos |
| 401 Unauthorized | No autenticado o token inválido |
| 403 Forbidden | Autenticado pero sin permisos |
| 404 Not Found | Recurso no encontrado |
| 409 Conflict | Conflicto (ej: transacción duplicada) |
| 429 Too Many Requests | Rate limit excedido |
| 500 Internal Server Error | Error interno del servidor |

---

## 🧪 Testing (Postman/Curl)

### Ejemplo: Calcular Opciones de Recarga
```bash
curl -X GET https://api.sorteos.club/api/v1/wallet/recharge-options
```

### Ejemplo: Consultar Saldo
```bash
curl -X GET https://api.sorteos.club/api/v1/wallet/balance \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Ejemplo: Agregar Fondos
```bash
curl -X POST https://api.sorteos.club/api/v1/wallet/add-funds \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -d '{
    "amount": "100.00",
    "payment_method": "stripe"
  }'
```

### Ejemplo: Listar Transacciones
```bash
curl -X GET "https://api.sorteos.club/api/v1/wallet/transactions?limit=10&offset=0" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 🚧 Pendiente (Próximas Fases)

- [ ] **Webhook de Stripe** (`POST /api/v1/wallet/webhook/stripe`)
  - Confirma pagos de add-funds
  - Sin autenticación (firma de Stripe)
  - Actualiza transacción pending → completed

- [ ] **Retiros** (`POST /api/v1/wallet/withdraw`)
  - Retirar fondos a cuenta bancaria
  - Requiere KYC verificado
  - Período de hold (3-7 días)

- [ ] **Integración con Pagos de Sorteos**
  - Usar billetera como método de pago en compra de boletos
  - Débito automático con locks de concurrencia

---

**Versión**: 1.0
**Última actualización**: 2025-11-18
**Estado**: MVP - Endpoints core implementados, pendiente integración Stripe completa
