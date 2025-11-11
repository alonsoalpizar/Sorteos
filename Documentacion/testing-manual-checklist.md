# Checklist de Testing Manual - Sprint 5-6

**Fecha:** 2025-11-11
**Tester:** _________________
**Duración:** ~30-45 minutos
**Entorno:** Docker local + PayPal Sandbox

---

## Pre-requisitos

### 1. Entorno levantado
```bash
cd /opt/Sorteos
docker compose up -d
docker compose logs -f api  # Verificar que no haya errores
```

- [ ] Postgres corriendo en puerto 5432
- [ ] Redis corriendo en puerto 6379
- [ ] API corriendo en puerto 8080
- [ ] Frontend corriendo en puerto 5173
- [ ] Migraciones aplicadas (ver logs del API)

### 2. PayPal Sandbox configurado
- [ ] Cuenta Business creada en https://developer.paypal.com
- [ ] Cuenta Personal creada (comprador)
- [ ] Client ID y Secret configurados en `.env`
- [ ] `CONFIG_PAYMENT_SANDBOX=true` en `.env`

### 3. Navegador preparado
- [ ] Chrome/Firefox con DevTools abierto
- [ ] Network tab abierto para ver requests
- [ ] Console abierto para ver errores JS

---

## Test Suite 1: Autenticación (5 min)

### TC-1.1: Registro de Usuario
**URL:** http://localhost:5173/register

- [ ] Abrir página de registro
- [ ] Llenar formulario:
  - Email: `test+buyer@example.com`
  - Password: `Test123!`
  - Nombre: `Test Buyer`
- [ ] Click "Registrarse"
- [ ] **Expected:** Redirect a `/verify-email` con mensaje
- [ ] **Actual:** _________________

**Notas:** _________________

### TC-1.2: Verificación de Email (mock)
**Nota:** En desarrollo, verificar automáticamente en DB

```sql
-- Ejecutar en postgres
UPDATE users SET email_verified = true
WHERE email = 'test+buyer@example.com';
```

- [ ] Email marcado como verificado
- [ ] **Expected:** `email_verified = true` en DB
- [ ] **Actual:** _________________

### TC-1.3: Login
**URL:** http://localhost:5173/login

- [ ] Ingresar email: `test+buyer@example.com`
- [ ] Ingresar password: `Test123!`
- [ ] Click "Iniciar Sesión"
- [ ] **Expected:** Redirect a `/dashboard` con navbar mostrando nombre
- [ ] **Actual:** _________________
- [ ] **Token en localStorage:** Verificar en DevTools > Application > Local Storage

**Notas:** _________________

---

## Test Suite 2: Creación de Sorteo (10 min)

### TC-2.1: Crear Sorteo como Draft
**URL:** http://localhost:5173/raffles/create

- [ ] Navegar a "Crear Sorteo"
- [ ] Llenar formulario:
  - Título: `Test Sorteo - iPhone 15 Pro`
  - Descripción: `Sorteo de prueba con PayPal`
  - Precio por número: `5000` (₡5,000)
  - Total números: `100`
  - Fecha de sorteo: `2025-12-25`
  - Método: `Lotería Nacional`
- [ ] Click "Guardar como Borrador"
- [ ] **Expected:** Redirect a `/my-raffles` con sorteo en estado "Borrador"
- [ ] **Actual:** _________________

**Raffle ID (anotar):** _________________

### TC-2.2: Publicar Sorteo
**URL:** http://localhost:5173/my-raffles

- [ ] Click en sorteo creado
- [ ] Verificar estado "Borrador"
- [ ] Click "Publicar Sorteo"
- [ ] Confirmar en modal
- [ ] **Expected:** Estado cambia a "Activo" con badge verde
- [ ] **Actual:** _________________
- [ ] **Network:** Verificar `PATCH /api/v1/raffles/{id}/publish` → 200 OK

**Notas:** _________________

---

## Test Suite 3: Selección y Reserva (15 min)

### TC-3.1: Ver Sorteo Público
**URL:** http://localhost:5173/raffles

- [ ] Click "Ver Sorteos" en navbar
- [ ] Verificar que el sorteo publicado aparece en listado
- [ ] Click en el sorteo
- [ ] **Expected:** Página de detalle con grid de 100 números
- [ ] **Actual:** _________________

### TC-3.2: Selección de 1 Número
- [ ] Click en número "0001"
- [ ] **Expected:**
  - Número se marca en azul
  - Contador "Números seleccionados: 1" aparece
  - Total: ₡5,000
  - Botón "Proceder al Pago" visible
- [ ] **Actual:** _________________

### TC-3.3: Selección de Múltiples Números
- [ ] Click en números: "0002", "0003", "0010", "0042"
- [ ] **Expected:**
  - Contador "Números seleccionados: 5"
  - Total: ₡25,000
  - Números ordenados: 0001, 0002, 0003, 0010, 0042
- [ ] **Actual:** _________________

### TC-3.4: Limpiar Selección
- [ ] Click "Limpiar selección"
- [ ] **Expected:**
  - Números deseleccionados (vuelven a blanco)
  - Contador desaparece
  - Botón "Proceder al Pago" desaparece
- [ ] **Actual:** _________________

### TC-3.5: Proceder al Checkout
- [ ] Seleccionar 3 números: "0005", "0015", "0025"
- [ ] Click "Proceder al Pago"
- [ ] **Expected:** Redirect a `/checkout`
- [ ] **Actual:** _________________

---

## Test Suite 4: Checkout y Reserva (10 min)

### TC-4.1: Página de Checkout
**URL:** http://localhost:5173/checkout

- [ ] **Resumen visible:**
  - Título del sorteo
  - 3 números seleccionados
  - Precio por número: ₡5,000
  - Total: ₡15,000
- [ ] Botón "Crear Reserva" visible
- [ ] **Actual:** _________________

### TC-4.2: Crear Reserva
- [ ] Click "Crear Reserva"
- [ ] **Expected:**
  - Loading spinner "Creando tu reserva..."
  - POST `/api/v1/reservations` → 201 Created
  - Timer aparece: "Tiempo restante: 4:59"
  - Mensaje "¡Reserva creada exitosamente!"
  - Botón "Proceder al Pago con PayPal"
- [ ] **Actual:** _________________
- [ ] **Network:** Copiar `reservation_id` del response: _________________

### TC-4.3: Timer de Expiración
- [ ] Observar timer contando hacia abajo
- [ ] **Expected:**
  - Timer actualiza cada segundo
  - Cuando llega a <1:00, fondo cambia a amarillo/rojo (urgente)
- [ ] **Actual:** _________________

**Nota:** No esperar los 5 minutos completos, continuar con pago.

### TC-4.4: Crear Payment Intent
- [ ] Click "Proceder al Pago con PayPal"
- [ ] **Expected:**
  - Loading "Preparando el pago..."
  - POST `/api/v1/payments/intent` → 201 Created
  - Redirect a PayPal sandbox (URL: https://www.sandbox.paypal.com/...)
- [ ] **Actual:** _________________

**PayPal Order ID (anotar):** _________________

---

## Test Suite 5: Pago con PayPal Sandbox (10 min)

### TC-5.1: Login en PayPal Sandbox
**URL:** https://www.sandbox.paypal.com (redirect automático)

- [ ] Página de PayPal carga correctamente
- [ ] Campos de login visibles
- [ ] Ingresar credenciales de **Personal Account** (comprador)
  - Email: _________________
  - Password: _________________
- [ ] Click "Log In"
- [ ] **Expected:** Página de revisión de orden
- [ ] **Actual:** _________________

### TC-5.2: Aprobar Pago
- [ ] Verificar detalles de la orden:
  - Merchant: "Sorteos Platform"
  - Amount: $15.00 (o equivalente según currency)
  - Description: menciona reservation_id
- [ ] Click "Pay Now" o "Complete Purchase"
- [ ] **Expected:**
  - Loading/spinner
  - Redirect de vuelta a `http://localhost:5173/payment/success?payment_id=xxx`
- [ ] **Actual:** _________________

### TC-5.3: Payment Success Page
**URL:** http://localhost:5173/payment/success

- [ ] **Expected:**
  - Confetti animation 🎉
  - Mensaje "¡Pago completado exitosamente!"
  - Payment ID visible
  - Reservation ID visible
  - Botón "Ver Mis Compras"
  - Botón "Ver Sorteos"
  - Cart limpio (verificar localStorage)
- [ ] **Actual:** _________________

**Payment ID:** _________________

### TC-5.4: Verificar Estado en DB
```sql
-- Ejecutar en postgres
SELECT status, payment_intent_id FROM payments
WHERE id = '<payment_id>';

SELECT status FROM reservations
WHERE id = '<reservation_id>';

SELECT status FROM raffle_numbers
WHERE id IN ('0005', '0015', '0025');
```

- [ ] **Expected:**
  - Payment status: `succeeded`
  - Reservation status: `confirmed`
  - Numbers status: `sold`
- [ ] **Actual:** _________________

---

## Test Suite 6: Casos de Error (10 min)

### TC-6.1: Cancelar Pago en PayPal
**Setup:** Repetir TC-3 y TC-4 con nuevos números

- [ ] Seleccionar números: "0030", "0031"
- [ ] Crear reserva
- [ ] Click "Proceder al Pago con PayPal"
- [ ] En PayPal, click "Cancel and return to merchant"
- [ ] **Expected:**
  - Redirect a `/payment/cancel`
  - Mensaje "Pago cancelado"
  - Opción "Volver al Checkout"
  - Reserva sigue activa con timer
- [ ] **Actual:** _________________

### TC-6.2: Expiración de Reserva
**Setup:** Modificar timer a 30 segundos para testing rápido

Option 1 - Modificar código temporalmente:
```typescript
// ReservationTimer.tsx
const TIMEOUT_MS = 30000; // 30 segundos en lugar de 5 minutos
```

Option 2 - Esperar 5 minutos (no recomendado para test rápido)

- [ ] Crear reserva con números: "0040", "0041"
- [ ] NO proceder al pago
- [ ] Esperar que timer llegue a 0:00
- [ ] **Expected:**
  - Mensaje "Reserva Expirada"
  - "Serás redirigido al listado de sorteos..."
  - Redirect automático después de 3 segundos
  - Números liberados en DB
- [ ] **Actual:** _________________

### TC-6.3: Número Ya Vendido (Race Condition)
**Setup:** Usar 2 navegadores o sesiones

**Navegador A:**
- [ ] Login como `test+buyer@example.com`
- [ ] Seleccionar número "0050"
- [ ] Crear reserva
- [ ] Completar pago con PayPal

**Navegador B (sin cerrar A):**
- [ ] Login como `test+buyer2@example.com` (crear cuenta nueva)
- [ ] Intentar seleccionar número "0050"
- [ ] **Expected:**
  - Número aparece en gris (status: sold)
  - Cursor: not-allowed
  - No se puede seleccionar
  - Tooltip: "Número 0050 - Vendido"
- [ ] **Actual:** _________________

### TC-6.4: Reserva Duplicada (Idempotency)
**Setup:** Usar DevTools para reenviar request

- [ ] Crear reserva con números "0060", "0061"
- [ ] En DevTools > Network, buscar POST `/api/v1/reservations`
- [ ] Right-click → "Replay XHR" o "Copy as cURL" y ejecutar 2 veces
- [ ] **Expected:**
  - Primera request: 201 Created
  - Segunda request: 200 OK (misma reserva, idempotency key match)
  - NO se crea reserva duplicada
- [ ] **Actual:** _________________

---

## Test Suite 7: Mis Compras (5 min)

### TC-7.1: Ver Historial de Compras
**URL:** http://localhost:5173/my-purchases

- [ ] Navegar a "Mis Compras"
- [ ] **Expected:**
  - Listado con todas las compras del usuario
  - Para cada compra:
    - Título del sorteo
    - Números comprados
    - Monto pagado
    - Estado: "Pagado" con badge verde
    - Fecha de compra
- [ ] **Actual:** _________________

### TC-7.2: Ver Detalle de Compra
- [ ] Click en una compra
- [ ] **Expected:**
  - Página de detalle del sorteo
  - Números comprados destacados en verde (readonly)
  - NO se puede seleccionar más números (readonly mode)
- [ ] **Actual:** _________________

---

## Resumen de Resultados

### ✅ Tests Exitosos: ___ / 30

### ❌ Tests Fallidos: ___ / 30

### Bugs Críticos Encontrados:
1. _________________
2. _________________
3. _________________

### Bugs Menores:
1. _________________
2. _________________

### Mejoras Sugeridas:
1. _________________
2. _________________

---

## Observaciones Generales

**Performance:**
- [ ] Todas las páginas cargan en < 2 segundos
- [ ] No hay errores en console de navegador
- [ ] API responde en < 500ms (verificar en Network tab)

**UX:**
- [ ] Loading states claros en todas las acciones
- [ ] Mensajes de error informativos
- [ ] Navegación fluida sin recargas innecesarias
- [ ] Timer visible y claro

**Seguridad:**
- [ ] Tokens JWT presentes en requests
- [ ] Protected routes redirigen a login si no autenticado
- [ ] No se pueden reservar números de otros usuarios

---

## Siguiente Paso

Una vez completado este checklist:

1. **Si 0 bugs críticos:** ✅ Proceder a Nivel 2 (Testing de API)
2. **Si 1-3 bugs críticos:** 🔧 Resolver bugs y re-ejecutar checklist
3. **Si >3 bugs críticos:** 🚨 Revisar implementación completa

**Tiempo total:** _____ minutos

**Feedback del tester:**
_________________
_________________
_________________
