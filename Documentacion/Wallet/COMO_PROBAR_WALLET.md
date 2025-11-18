# 🧪 Cómo Probar el Sistema de Billetera

## 📋 Requisitos Previos

1. ✅ Backend corriendo en `http://localhost:8080`
2. ✅ Base de datos PostgreSQL con migraciones aplicadas
3. ✅ Frontend corriendo en `http://localhost:5173`
4. ✅ Usuario registrado y autenticado

---

## 🚀 Pasos para Probar

### 1. Levantar el Backend

```bash
cd /opt/Sorteos/backend
go run cmd/api/main.go
```

Verificar que veas en los logs:
```
✓ Wallet routes configured successfully
```

### 2. Levantar el Frontend

```bash
cd /opt/Sorteos/frontend
npm run dev
```

Debería abrir en `http://localhost:5173`

### 3. Registrar un Usuario (si no tienes uno)

1. Ir a `http://localhost:5173/register`
2. Completar formulario:
   - Email: `test@example.com`
   - Password: `TestPassword123!`
   - Aceptar términos y privacidad
3. Click "Registrarse"
4. El sistema **auto-crea la billetera** al registrar

### 4. Verificar que la Billetera fue Creada

**Opción A: Revisar en la base de datos**
```sql
-- Conectarse a PostgreSQL
psql -U sorteos_user -d sorteos_db

-- Ver billeteras creadas
SELECT id, uuid, user_id, balance, currency, status
FROM wallets
WHERE user_id = (SELECT id FROM users WHERE email = 'test@example.com');

-- Deberías ver:
-- id | uuid | user_id | balance | currency | status
-- 1  | xxx  | 1       | 0.00    | CRC      | active
```

**Opción B: Llamar al API**
```bash
# 1. Login para obtener token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"TestPassword123!"}'

# Copiar el access_token de la respuesta

# 2. Consultar saldo
curl http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer <access_token>"

# Respuesta esperada:
# {
#   "success": true,
#   "data": {
#     "wallet_id": 1,
#     "balance": "0.00",
#     "pending_balance": "0.00",
#     "currency": "CRC",
#     "status": "active"
#   }
# }
```

---

## 🎨 Probar el Frontend

### Acceder a la Billetera

1. **Login** en `http://localhost:5173/login`
2. En el **navbar superior**, verás el enlace **"💰 Billetera"**
3. Click en "💰 Billetera" → Te lleva a `/wallet`

### Tab 1: Mi Saldo

✅ **Qué deberías ver:**
- Card con saldo: `₡0` (cero colones)
- Saldo pendiente: `₡0`
- Moneda: `CRC`
- Estado: `Activa` (verde)
- Botón de refrescar (icono)
- 2 botones grandes:
  - "Recargar saldo" → Cambia al tab de recarga
  - "Ver historial" → Cambia al tab de historial
- Info box con explicación de cómo funciona

### Tab 2: Recargar

✅ **Qué deberías ver:**
- Alert azul con nota informativa
- **5 cards** con opciones de recarga:
  - ₡1,000 → Total a pagar: ₡1,155.79
  - ₡5,000 → Total a pagar: ₡5,378.95
  - ₡10,000 → Total a pagar: ₡10,640.00
  - ₡15,000 → Total a pagar: ₡15,901.05
  - ₡20,000 → Total a pagar: ₡21,162.11

✅ **Probar selección de opción:**
1. Click en el card de **₡5,000**
2. El card debe resaltarse con borde azul
3. Aparece **checkmark verde** ✓
4. Abajo aparece **desglose detallado**:
   - Crédito deseado: ₡5,000.00
   - Tarifa fija: ₡100.00
   - Comisión procesador (3%): ₡268.42
   - Comisión plataforma (2%): ₡110.53
   - **Total a pagar: ₡5,378.95**

✅ **Probar métodos de pago:**
1. Aparecen 3 botones:
   - 💳 Tarjeta
   - 💸 SINPE Móvil
   - 🏦 Transferencia
2. Click en cada uno → Se resalta con borde azul

✅ **Confirmar recarga:**
1. Click en botón azul grande: **"Recargar ₡5,000"**
2. Botón cambia a "Procesando..." con spinner
3. Después de ~1 segundo, aparece **alert verde de éxito**:
   - ✓ ¡Transacción creada exitosamente!
   - ID de transacción: `xxx-xxx-xxx`
   - Monto: ₡5,000.00
   - Estado: `pending`
   - Nota: "En esta fase de desarrollo, el pago real aún no está habilitado"
   - Botón: "Realizar otra recarga"

### Tab 3: Historial

✅ **Primera vez (sin transacciones):**
- Icono 📊
- Texto: "No hay transacciones"
- Descripción: "Aún no has realizado ninguna transacción en tu billetera"

✅ **Con transacciones (después de crear una recarga):**
- Header: "Historial de Transacciones (1)"
- Botón de refrescar
- **Desktop**: Tabla con columnas:
  - Fecha | Tipo | Monto | Estado | Saldo después
- **Mobile**: Cards compactos
- Transacción mostrada:
  - Tipo: "Recarga"
  - Monto: `+₡5,000.00` (verde)
  - Estado: Badge amarillo "Pendiente"
  - Saldo después: ₡0.00 (porque está pending)

### Paginación (si hay >20 transacciones)

- Footer con: "Página 1 de X"
- Botones: "← Anterior" (deshabilitado) y "Siguiente →"

---

## 🧪 Probar Endpoints Directamente

### 1. Opciones de Recarga (Público)

```bash
curl http://localhost:8080/api/v1/wallet/recharge-options
```

Respuesta esperada: JSON con 5 opciones

### 2. Consultar Saldo (Autenticado)

```bash
# Primero hacer login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"TestPassword123!"}' \
  | jq -r '.data.access_token')

# Consultar saldo
curl http://localhost:8080/api/v1/wallet/balance \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Agregar Fondos (Autenticado)

```bash
curl -X POST http://localhost:8080/api/v1/wallet/add-funds \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(uuidgen)" \
  -d '{
    "amount": "5000",
    "payment_method": "card"
  }'
```

Respuesta esperada:
```json
{
  "success": true,
  "message": "Transacción de depósito creada. Complete el pago con el procesador.",
  "data": {
    "transaction_id": 1,
    "transaction_uuid": "xxx-xxx",
    "amount": "5000.00",
    "status": "pending",
    "payment_method": "card",
    "idempotency_key": "xxx"
  }
}
```

### 4. Ver Transacciones (Autenticado)

```bash
curl "http://localhost:8080/api/v1/wallet/transactions?limit=20&offset=0" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🐛 Troubleshooting

### Error: "Billetera no encontrada"

**Causa**: El usuario no tiene billetera creada automáticamente

**Solución**:
```sql
-- Crear billetera manualmente
INSERT INTO wallets (uuid, user_id, balance, pending_balance, currency, status, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  (SELECT id FROM users WHERE email = 'test@example.com'),
  0,
  0,
  'CRC',
  'active',
  NOW(),
  NOW()
);
```

### Error: "CORS blocked"

**Causa**: Frontend en puerto diferente al backend

**Solución**: Verificar que el backend tenga CORS habilitado para `http://localhost:5173`

### Error: "Network request failed"

**Causa**: Backend no está corriendo

**Solución**: Levantar backend con `go run cmd/api/main.go`

### Transacciones no aparecen en historial

**Causa**: Las transacciones PENDING no afectan el balance hasta ser confirmadas

**Solución**: Normal. En producción, el webhook del procesador las confirmará.

---

## ✅ Checklist de Prueba

- [ ] Backend corriendo y respondiendo
- [ ] Frontend corriendo
- [ ] Usuario registrado
- [ ] Billetera auto-creada al registrar
- [ ] Enlace "💰 Billetera" visible en navbar
- [ ] Tab "Mi Saldo" muestra ₡0
- [ ] Tab "Recargar" muestra 5 opciones
- [ ] Selección de opción resalta el card
- [ ] Desglose de comisiones aparece
- [ ] Métodos de pago seleccionables
- [ ] Botón "Recargar" crea transacción
- [ ] Alert de éxito aparece
- [ ] Tab "Historial" muestra la transacción
- [ ] Badge de estado "Pendiente" amarillo
- [ ] Monto con signo "+" en verde
- [ ] Botón refrescar funciona
- [ ] Responsive en mobile

---

## 🎯 Próximos Pasos

Una vez validado que todo funciona:

1. **Integrar con checkout de sorteos** → Usar saldo para pagar boletos
2. **Implementar webhook** del procesador de pagos local (BAC/BCR/SINPE)
3. **Confirmar transacciones** → Cambiar status de PENDING → COMPLETED
4. **Acreditar saldo** cuando el webhook confirme el pago

---

**Versión**: 1.0
**Fecha**: 2025-11-18
**Estado**: ✅ Listo para pruebas funcionales

