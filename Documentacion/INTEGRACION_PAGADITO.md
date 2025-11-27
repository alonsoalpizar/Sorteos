# Integración de Pagadito para Recarga de Créditos

## 📋 Resumen

Este documento describe la integración de **Pagadito** (procesador de pagos costarricense) con el sistema de billetera de Sorteos para permitir la recarga de créditos.

## ✅ Implementación Completada

### 1. SDK de Pagadito (Go)
**Ubicación:** `/opt/Sorteos/backend/internal/infrastructure/pagadito/`

Archivos creados:
- `types.go` - Tipos de datos y interfaces
- `errors.go` - Mapeo de códigos de error Pagadito
- `client.go` - Cliente HTTP para API de Pagadito

**Características:**
- Autenticación con `Connect()` (tokens con TTL de 30 min)
- Creación de transacciones con `CreateTransaction()`
- Consulta de estado con `GetStatus()`
- Soporte para sandbox y producción
- Manejo de errores estándar de Pagadito (PG1001-PG3006)

### 2. Base de Datos
**Migración:** `000022_create_credit_purchases.up.sql`

**Tabla `credit_purchases`:**
- Registro de todas las compras de créditos
- Estados: pending, processing, completed, failed, expired
- Integración con Pagadito (ERN, token, reference)
- Desglose de comisiones para transparencia
- Idempotencia mediante `idempotency_key` y `ern`
- TTL de 30 minutos (campo `expires_at`)

### 3. Entidad de Dominio
**Archivo:** `internal/domain/credit_purchase.go`

**Características:**
- Entidad `CreditPurchase` con estados y validaciones
- Generación de ERN (External Reference Number) único
- Métodos de transición de estado (MarkAsProcessing, MarkAsCompleted, etc.)
- Repositorio con interfaz bien definida

### 4. Repositorio
**Archivo:** `internal/adapters/db/credit_purchase_repository.go`

**Implementación PostgreSQL con:**
- CRUD completo
- Búsqueda por ERN, UUID, token de Pagadito
- Paginación para historial de usuario
- Método `MarkExpired()` para cron jobs

### 5. Casos de Uso

#### `PurchaseCreditsUseCase`
**Archivo:** `internal/usecase/credits/purchase_credits.go`

**Flujo:**
1. Valida monto (mín/máx según moneda)
2. Verifica idempotencia
3. Obtiene billetera del usuario
4. Calcula comisiones con `RechargeCalculator`
5. Genera ERN único
6. Crea registro en DB (estado: pending)
7. Conecta con Pagadito
8. Crea transacción en Pagadito
9. Retorna URL de pago para redirección

#### `ProcessPagaditoCallbackUseCase`
**Archivo:** `internal/usecase/credits/process_callback.go`

**Flujo:**
1. Busca compra por ERN (token del callback)
2. Verifica idempotencia (si ya fue procesada)
3. Consulta estado real en Pagadito
4. Según estado:
   - **COMPLETED**: Acredita créditos con `AddFundsUseCase` → Éxito
   - **VERIFYING**: Mantiene en procesamiento → Espera verificación manual
   - **REGISTERED/FAILED/REVOKED**: Marca como fallida → Error
5. Actualiza compra en DB
6. Crea log de auditoría
7. Retorna URL de redirección al frontend

## 🔄 Flujo Completo de Usuario

```
┌─────────────────────────────────────────────────────────────────────┐
│ 1. Usuario en Frontend: Click "Recargar Créditos"                  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 2. POST /api/v1/credits/purchase                                   │
│    Body: {desired_credit: 5000, currency: "CRC"}                   │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 3. Backend: PurchaseCreditsUseCase                                 │
│    - Calcula comisiones (₡5,000 → ₡5,500 con fees)                │
│    - Crea credit_purchase en DB                                    │
│    - Llama Pagadito API                                            │
│    - Obtiene payment_url                                           │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 4. Response: {payment_url: "https://pagadito.com/pay/xyz"}        │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 5. Frontend: window.location.href = payment_url                    │
│    (Usuario redirigido a Pagadito)                                 │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 6. Usuario en página de Pagadito:                                  │
│    - Ingresa datos de tarjeta / selecciona método de pago         │
│    - Confirma pago                                                 │
│    - Pagadito procesa transacción                                  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 7. Pagadito redirige: GET /api/v1/credits/callback?token=ERN_XXX  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 8. Backend: ProcessPagaditoCallbackUseCase                        │
│    - Consulta estado en Pagadito (GetStatus)                      │
│    - Si COMPLETED:                                                 │
│      * Ejecuta AddFundsUseCase                                     │
│      * Acredita ₡5,000 a billetera                                │
│      * Marca compra como completed                                 │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 9. Redirect: /credits/success?purchase_id=UUID&amount=5000        │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│ 10. Frontend: Muestra mensaje de éxito                             │
│     "¡Créditos acreditados! Nuevo saldo: ₡10,000"                 │
└─────────────────────────────────────────────────────────────────────┘
```

## 🚧 Pendiente de Implementar

### 1. Handlers API
**Archivos a crear:**
- `cmd/api/handlers/credits_handler.go`
  - `PurchaseCreditsHandler.Handle()` → POST /credits/purchase
  - `CallbackHandler.Handle()` → GET /credits/callback
  - `GetPackagesHandler.Handle()` → GET /credits/packages
  - `GetPurchaseStatusHandler.Handle()` → GET /credits/purchase/:id

### 2. Rutas API
**Archivo a modificar:** `cmd/api/routes.go`

```go
creditsGroup := v1.Group("/credits")
{
    creditsGroup.GET("/packages", getPackagesHandler)
    creditsGroup.POST("/purchase", authMiddleware, purchaseCreditsHandler)
    creditsGroup.GET("/callback", callbackHandler) // Sin auth (público)
    creditsGroup.GET("/purchase/:id", authMiddleware, getPurchaseStatusHandler)
}
```

### 3. Configuración de Pagadito (Admin)
**Requisito:** Configurar UID, WSK y otros parámetros desde Admin Dashboard

**Opciones:**
- **Opción A**: Usar tabla `payment_processors` existente
- **Opción B**: Crear sección en `system_parameters`
- **Opción C**: CRUD específico en admin

**Campos necesarios:**
- `uid` (Merchant ID)
- `wsk` (Secret Key)
- `sandbox_mode` (boolean)
- `currency` (CRC/USD)
- `callback_url`

### 4. Frontend

#### Componentes a crear:
- `components/Credits/RechargeModal.tsx` - Modal para comprar créditos
- `components/Credits/PackageCard.tsx` - Card de paquete predefinido
- `pages/Credits/Success.tsx` - Página de éxito
- `pages/Credits/Failed.tsx` - Página de error
- `pages/Credits/Verifying.tsx` - Página de verificación pendiente

#### Flujo sugerido:
```tsx
// En cualquier parte del frontend
<Button onClick={() => setShowRechargeModal(true)}>
  Recargar Créditos
</Button>

<RechargeModal
  show={showRechargeModal}
  onClose={() => setShowRechargeModal(false)}
  onPurchase={(amount) => handlePurchase(amount)}
/>

// Handler
const handlePurchase = async (desiredCredit) => {
  const response = await api.post('/credits/purchase', {
    desired_credit: desiredCredit,
    currency: 'CRC',
  })

  // Redirigir a Pagadito
  window.location.href = response.data.payment_url
}
```

### 5. Cron Jobs

#### Expirar compras pendientes
**Archivo:** `internal/jobs/expire_credit_purchases.go`

```go
// Ejecutar cada 5 minutos
func ExpireCreditPurchases() {
    count, err := creditPurchaseRepo.MarkExpired()
    if err != nil {
        logger.Error("Error expirando compras", logger.Error(err))
        return
    }
    if count > 0 {
        logger.Info("Compras expiradas", logger.Int64("count", count))
    }
}
```

## 📊 Datos de Ejemplo

### Paquetes Predefinidos Recomendados

| Crédito | Comisión | Total a Pagar | Badge |
|---------|----------|---------------|-------|
| ₡1,000  | ₡300     | ₡1,300        | -     |
| ₡5,000  | ₡500     | ₡5,500        | POPULAR |
| ₡10,000 | ₡900     | ₡10,900       | -     |
| ₡15,000 | ₡1,300   | ₡16,300       | -     |
| ₡20,000 | ₡1,700   | ₡21,700       | MEJOR VALOR |
| ₡30,000 | ₡2,500   | ₡32,500       | -     |

**Nota:** Las comisiones se calculan dinámicamente usando `RechargeCalculator` basado en `system_parameters`.

### Estados de Transacción

| Estado | Descripción | Acción Usuario |
|--------|-------------|----------------|
| `pending` | Compra iniciada, esperando redirección | - |
| `processing` | Usuario en Pagadito | Completar pago |
| `completed` | Pago exitoso, créditos acreditados | Ver saldo |
| `failed` | Pago fallido o cancelado | Reintentar |
| `expired` | Expiró sin completarse (30 min) | Nueva compra |

### Estados de Pagadito

| Estado | Significado | Acción Sistema |
|--------|-------------|----------------|
| `COMPLETED` | Pago aprobado | Acreditar créditos |
| `REGISTERED` | Usuario canceló | Marcar como failed |
| `VERIFYING` | En verificación manual | Esperar decisión admin |
| `REVOKED` | Rechazado por Pagadito | Marcar como failed |
| `FAILED` | Error de procesamiento | Marcar como failed |

## 🔒 Seguridad

### Idempotencia
- `idempotency_key` en todas las operaciones de compra
- `ERN` único por transacción (formato: `CP_{user_id}_{timestamp}_{random}`)
- Verificación antes de procesar callbacks

### Validaciones
- Montos mínimo/máximo según moneda
- Estado de billetera (debe estar activa)
- Existencia de usuario y billetera
- Verificación de estado con Pagadito antes de acreditar

### Logs de Auditoría
- Compra iniciada (`credit_purchase_initiated`)
- Compra completada (`credit_purchase_completed`)
- Compra fallida (`credit_purchase_failed`)

## 🧪 Testing

### Sandbox de Pagadito
**URL:** https://sandbox.pagadito.com

**Credenciales de prueba:**
- UID: (solicitar a Pagadito)
- WSK: (solicitar a Pagadito)

### Tarjetas de prueba
- Éxito: 4111111111111111
- Fallo: 4242424242424242

### Escenarios a probar:
1. ✅ Compra exitosa (COMPLETED)
2. ✅ Usuario cancela en Pagadito (REGISTERED)
3. ✅ Pago en verificación (VERIFYING)
4. ✅ Compra duplicada (idempotencia)
5. ✅ Callback duplicado (idempotencia)
6. ✅ Compra expira sin completarse
7. ✅ Error de conexión con Pagadito

## 📝 Próximos Pasos

1. **Implementar Handlers** (30 min)
2. **Agregar Rutas** (15 min)
3. **Configurar Admin CRUD para Pagadito** (1 hora)
4. **Crear Frontend** (2-3 horas)
5. **Testing en Sandbox** (1 hora)
6. **Deployment a Producción** (30 min)

**Tiempo total estimado:** 1 día de trabajo

## 🔗 Referencias

- [Documentación Pagadito](https://dev.pagadito.com/)
- [API Reference](https://dev.pagadito.com/index.php?mod=docs&hac=apipg)
- Ejemplos analizados:
  - `/opt/Sorteos/mitiendapagadito4` (Java)
  - `/opt/Sorteos/mitiendapagadito_1.4.1` (PHP)

---

**Generado el:** 2025-11-26
**Autor:** Claude Code
**Versión:** 1.0
