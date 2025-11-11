# Testing de API - Scripts cURL

**Fecha:** 2025-11-11
**Sprint:** 5-6 (Reservas y Pagos)
**Herramienta:** cURL / bash
**Duración:** 1-2 horas

---

## Setup

### Variables de Entorno
```bash
# Guardar en ~/.bashrc o ejecutar en terminal
export API_URL="http://localhost:8080/api/v1"
export TOKEN=""  # Se llenará después del login
export RAFFLE_ID=""
export RESERVATION_ID=""
export PAYMENT_ID=""
```

### Funciones Helper
```bash
# Guardar en ~/sorteos-test-helpers.sh

# Pretty print JSON responses
alias pj='python3 -m json.tool'

# Extract field from JSON
jqe() {
  echo "$1" | jq -r "$2"
}

# POST with auth
post_auth() {
  curl -X POST "$API_URL$1" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$2" \
    -w "\n%{http_code}\n" \
    -s
}

# GET with auth
get_auth() {
  curl -X GET "$API_URL$1" \
    -H "Authorization: Bearer $TOKEN" \
    -w "\n%{http_code}\n" \
    -s
}

# Source helpers
source ~/sorteos-test-helpers.sh
```

---

## Test Suite 1: Autenticación

### 1.1 Registro de Usuario
```bash
curl -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-api@example.com",
    "password": "Test123!",
    "first_name": "Test",
    "last_name": "API User"
  }' | pj

# Expected: 201 Created
# {
#   "data": {
#     "user": {
#       "id": 1,
#       "email": "test-api@example.com",
#       "email_verified": false
#     }
#   }
# }
```

### 1.2 Verificar Email (Mock en DB)
```sql
-- Ejecutar en postgres
UPDATE users SET email_verified = true
WHERE email = 'test-api@example.com';
```

### 1.3 Login
```bash
RESPONSE=$(curl -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-api@example.com",
    "password": "Test123!"
  }' -s)

echo "$RESPONSE" | pj

# Extraer token
export TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')
echo "Token saved: ${TOKEN:0:20}..."

# Expected: 200 OK con tokens
```

### 1.4 Get Current User
```bash
get_auth "/auth/me" | pj

# Expected: 200 OK
# {
#   "data": {
#     "user": {
#       "id": 1,
#       "email": "test-api@example.com",
#       "email_verified": true
#     }
#   }
# }
```

---

## Test Suite 2: CRUD de Sorteos

### 2.1 Crear Sorteo (Draft)
```bash
RESPONSE=$(post_auth "/raffles" '{
  "title": "Test API Raffle - iPhone 15",
  "description": "Testing con scripts cURL",
  "price_per_number": 5000,
  "total_numbers": 100,
  "draw_date": "2025-12-25T20:00:00Z",
  "draw_method": "loteria_nacional"
}')

echo "$RESPONSE" | pj

# Extraer UUID
export RAFFLE_ID=$(echo "$RESPONSE" | jq -r '.data.raffle.uuid')
echo "Raffle ID: $RAFFLE_ID"

# Expected: 201 Created
```

### 2.2 Obtener Sorteo
```bash
get_auth "/raffles/$RAFFLE_ID?include_numbers=true" | pj

# Expected: 200 OK con 100 números (status: available)
```

### 2.3 Publicar Sorteo
```bash
post_auth "/raffles/$RAFFLE_ID/publish" '{}' | pj

# Expected: 200 OK
# { "data": { "raffle": { "status": "active" } } }
```

### 2.4 Listar Sorteos Activos
```bash
get_auth "/raffles?status=active" | pj

# Expected: 200 OK con array de raffles
```

---

## Test Suite 3: Reservas

### 3.1 Crear Reserva (Happy Path)
```bash
SESSION_ID="test-session-$(date +%s)"

RESPONSE=$(post_auth "/reservations" "{
  \"raffle_id\": \"$RAFFLE_ID\",
  \"number_ids\": [\"0001\", \"0002\", \"0003\"],
  \"session_id\": \"$SESSION_ID\"
}")

echo "$RESPONSE" | pj

export RESERVATION_ID=$(echo "$RESPONSE" | jq -r '.data.reservation.id')
echo "Reservation ID: $RESERVATION_ID"

# Expected: 201 Created
# {
#   "data": {
#     "reservation": {
#       "id": "uuid",
#       "status": "pending",
#       "number_ids": ["0001", "0002", "0003"],
#       "total_amount": 15000,
#       "expires_at": "2025-11-11T02:10:00Z"
#     }
#   }
# }
```

### 3.2 Obtener Reserva
```bash
get_auth "/reservations/$RESERVATION_ID" | pj

# Expected: 200 OK
```

### 3.3 Intentar Reservar Números Ya Reservados (Conflict)
```bash
post_auth "/reservations" "{
  \"raffle_id\": \"$RAFFLE_ID\",
  \"number_ids\": [\"0001\", \"0010\"],
  \"session_id\": \"test-conflict-$(date +%s)\"
}" | pj

# Expected: 409 Conflict
# {
#   "error": {
#     "code": "NUMBERS_NOT_AVAILABLE",
#     "message": "Some numbers are not available"
#   }
# }
```

### 3.4 Listar Mis Reservas
```bash
get_auth "/reservations/me" | pj

# Expected: 200 OK con array de reservations
```

---

## Test Suite 4: Pagos

### 4.1 Crear Payment Intent
```bash
RESPONSE=$(post_auth "/payments/intent" "{
  \"reservation_id\": \"$RESERVATION_ID\",
  \"return_url\": \"http://localhost:5173/payment/success\",
  \"cancel_url\": \"http://localhost:5173/payment/cancel\"
}")

echo "$RESPONSE" | pj

# Extraer URL de PayPal
PAYPAL_URL=$(echo "$RESPONSE" | jq -r '.data.payment_intent.client_secret')
echo "PayPal URL: $PAYPAL_URL"

export PAYMENT_ID=$(echo "$RESPONSE" | jq -r '.data.payment_intent.id')
echo "Payment ID: $PAYMENT_ID"

# Expected: 201 Created
# {
#   "data": {
#     "payment_intent": {
#       "id": "uuid",
#       "status": "pending",
#       "amount": 15000,
#       "currency": "CRC",
#       "client_secret": "https://www.sandbox.paypal.com/..."
#     }
#   }
# }
```

### 4.2 Obtener Payment
```bash
get_auth "/payments/$PAYMENT_ID" | pj

# Expected: 200 OK
```

### 4.3 Listar Mis Pagos
```bash
get_auth "/payments/me" | pj

# Expected: 200 OK con array de payments
```

---

## Test Suite 5: Testing de Concurrencia

### 5.1 Distributed Locks Test (10 requests simultáneas)
```bash
# Crear script de test
cat > /tmp/test-concurrency.sh << 'EOF'
#!/bin/bash

API_URL="http://localhost:8080/api/v1"
TOKEN="$1"
RAFFLE_ID="$2"

# Función para hacer POST
test_request() {
  SESSION_ID="concurrent-$1-$(date +%s%N)"

  RESPONSE=$(curl -X POST "$API_URL/reservations" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"raffle_id\": \"$RAFFLE_ID\",
      \"number_ids\": [\"0050\"],
      \"session_id\": \"$SESSION_ID\"
    }" -w "\n%{http_code}" -s 2>&1)

  HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
  BODY=$(echo "$RESPONSE" | head -n-1)

  echo "[$1] HTTP $HTTP_CODE: $(echo "$BODY" | jq -r '.error.code // .data.reservation.id // "unknown"')"
}

export -f test_request
export API_URL TOKEN RAFFLE_ID

# 10 requests en paralelo
seq 10 | xargs -P10 -I {} bash -c 'test_request {}'

echo ""
echo "Expected: Solo 1 request con 201 Created, resto 409 Conflict"
EOF

chmod +x /tmp/test-concurrency.sh

# Ejecutar test
/tmp/test-concurrency.sh "$TOKEN" "$RAFFLE_ID"

# Expected output:
# [1] HTTP 201: abc-123-uuid
# [2] HTTP 409: NUMBERS_NOT_AVAILABLE
# [3] HTTP 409: NUMBERS_NOT_AVAILABLE
# ...
```

### 5.2 Verificar en DB que Solo Hay 1 Reserva
```sql
SELECT COUNT(*), number_ids
FROM reservations
WHERE raffle_id = '<RAFFLE_ID>'
  AND '0050' = ANY(number_ids)
GROUP BY number_ids;

-- Expected: COUNT = 1
```

---

## Test Suite 6: Idempotency

### 6.1 Request Duplicado con Mismo Session ID
```bash
SESSION_ID="idempotency-test-$(date +%s)"

# Primera request
RESPONSE1=$(post_auth "/reservations" "{
  \"raffle_id\": \"$RAFFLE_ID\",
  \"number_ids\": [\"0070\"],
  \"session_id\": \"$SESSION_ID\"
}")

RESERVATION_ID1=$(echo "$RESPONSE1" | jq -r '.data.reservation.id')
echo "First request - Reservation ID: $RESERVATION_ID1"

# Segunda request (EXACTAMENTE el mismo payload)
sleep 1
RESPONSE2=$(post_auth "/reservations" "{
  \"raffle_id\": \"$RAFFLE_ID\",
  \"number_ids\": [\"0070\"],
  \"session_id\": \"$SESSION_ID\"
}")

RESERVATION_ID2=$(echo "$RESPONSE2" | jq -r '.data.reservation.id')
echo "Second request - Reservation ID: $RESERVATION_ID2"

# Verificar que son el mismo ID
if [ "$RESERVATION_ID1" == "$RESERVATION_ID2" ]; then
  echo "✅ Idempotency works! Same reservation returned"
else
  echo "❌ Idempotency failed! Different reservations created"
fi

# Expected: Mismo reservation ID en ambas requests
```

---

## Test Suite 7: Validaciones y Errores

### 7.1 Crear Reserva Sin Auth
```bash
curl -X POST "$API_URL/reservations" \
  -H "Content-Type: application/json" \
  -d "{
    \"raffle_id\": \"$RAFFLE_ID\",
    \"number_ids\": [\"0080\"],
    \"session_id\": \"test\"
  }" | pj

# Expected: 401 Unauthorized
```

### 7.2 Crear Reserva con Raffle Inválido
```bash
post_auth "/reservations" '{
  "raffle_id": "00000000-0000-0000-0000-000000000000",
  "number_ids": ["0001"],
  "session_id": "test"
}' | pj

# Expected: 404 Not Found
```

### 7.3 Crear Reserva con Números Inválidos
```bash
post_auth "/reservations" "{
  \"raffle_id\": \"$RAFFLE_ID\",
  \"number_ids\": [\"9999\"],
  \"session_id\": \"test\"
}" | pj

# Expected: 400 Bad Request - INVALID_NUMBERS
```

### 7.4 Crear Payment Intent con Reserva Expirada
```bash
# Primero, marcar reserva como expirada en DB
# UPDATE reservations SET status = 'expired' WHERE id = '$RESERVATION_ID';

post_auth "/payments/intent" "{
  \"reservation_id\": \"$RESERVATION_ID\"
}" | pj

# Expected: 400 Bad Request - RESERVATION_EXPIRED
```

---

## Test Suite 8: Performance

### 8.1 Benchmark - Crear Reserva
```bash
# Usar Apache Bench
ab -n 100 -c 10 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -p /tmp/reservation-payload.json \
  "$API_URL/reservations"

# Expected:
# - 95% requests < 500ms
# - 0% failures
```

### 8.2 Benchmark - GET Raffle
```bash
ab -n 1000 -c 50 \
  -H "Authorization: Bearer $TOKEN" \
  "$API_URL/raffles/$RAFFLE_ID"

# Expected:
# - 95% requests < 200ms
# - 0% failures
```

---

## Test Suite 9: Cleanup

### 9.1 Limpiar Test Data
```sql
-- Ejecutar en postgres
DELETE FROM payments WHERE reservation_id IN (
  SELECT id FROM reservations WHERE user_id IN (
    SELECT id FROM users WHERE email LIKE 'test-%@example.com'
  )
);

DELETE FROM reservations WHERE user_id IN (
  SELECT id FROM users WHERE email LIKE 'test-%@example.com'
);

DELETE FROM raffles WHERE user_id IN (
  SELECT id FROM users WHERE email LIKE 'test-%@example.com'
);

DELETE FROM users WHERE email LIKE 'test-%@example.com';
```

---

## Resultado Esperado

### ✅ Casos de Éxito
- POST /auth/register → 201
- POST /auth/login → 200
- POST /raffles → 201
- POST /raffles/{id}/publish → 200
- POST /reservations → 201 (primera request)
- POST /reservations (duplicate) → 409 Conflict
- POST /reservations (idempotent) → 200 (mismo ID)
- POST /payments/intent → 201

### ⏱️ Performance
- Crear reserva: < 500ms (p95)
- GET raffle: < 200ms (p95)
- Concurrencia: 0 duplicados

### 🔒 Seguridad
- Sin auth → 401
- Invalid data → 400/404

---

## Próximo Paso

Una vez completado este testing:

1. **Si todos los tests pasan:** ✅ Documentar como "API estable"
2. **Si hay failures:** 🔧 Resolver bugs y re-ejecutar
3. **Performance issues:** ⚡ Optimizar queries/locks

**Total tests:** 30
**Duración:** 1-2 horas
