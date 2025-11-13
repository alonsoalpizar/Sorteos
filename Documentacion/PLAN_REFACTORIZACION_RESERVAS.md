# Plan de Refactorización: Flujo de Reservas

**Fecha:** 2025-11-13
**Estado:** PENDIENTE DE APROBACIÓN
**Prioridad:** ALTA - Sistema actual tiene bugs críticos

---

## Problema Actual

### Flujo Confuso y Buggy
```
Usuario selecciona números → ¿Cuándo se reservan?
    ↓
"Proceder al Pago" → ¿Crea reserva aquí?
    ↓
CheckoutPage → ¿También crea reserva aquí?
    ↓
ERROR: "Números ya reservados" (pero son MÍOS!)
```

### Issues Específicos

1. **Duplicación de lógica de creación de reserva**
   - RaffleDetailPage a veces crea reserva
   - CheckoutPage también intenta crear reserva
   - No hay claridad sobre CUÁNDO se debe crear

2. **Validación incorrecta en backend**
   - Backend dice "números ya reservados"
   - No valida si el usuario es el MISMO que tiene la reserva
   - Error 409 confunde al usuario

3. **Estado duplicado**
   - `useCartStore` tiene `activeReservation`
   - `useActiveReservation` consulta servidor
   - No están sincronizados
   - localStorage puede tener datos obsoletos

4. **Al entrar al sorteo**
   - Frontend valida si hay reserva activa
   - ¿Por qué? Si está entrando debería estar limpio
   - Genera confusión innecesaria

---

## Flujo Correcto (Intención Original del Usuario)

### 1. Usuario Entra al Sorteo
```
→ NO debe tener reservas previas
→ Si existe alguna, debe ser cancelada automáticamente
→ Pantalla limpia, lista para seleccionar
```

### 2. Usuario Selecciona PRIMER Número
```
→ Se crea reserva INMEDIATAMENTE en backend
→ Timeout: 10 minutos (fase: selection)
→ Número aparece seleccionado en grilla
→ WebSocket notifica a otros usuarios
```

### 3. Usuario Selecciona MÁS Números
```
→ Se agregan a la reserva EXISTENTE
→ Endpoint: POST /reservations/{id}/add-number
→ Timeout sigue corriendo desde el inicio
→ NO se crea nueva reserva
```

### 4. Usuario Hace Clic "Proceder al Pago"
```
→ YA TIENE reserva activa (creada en paso 2)
→ NO crea nada nuevo
→ Solo navega a /checkout
→ Opcionalmente: Mover a fase "checkout" (extiende timeout 5 min más)
```

### 5. CheckoutPage Se Monta
```
→ Carga reserva EXISTENTE del servidor
→ Muestra números reservados
→ Botón "Pagar con PayPal"
→ NO intenta crear nueva reserva
```

### 6. Timeouts y Expiración
```
- Si 10+ minutos en grilla → Expira, libera todo
- Si sale de cualquier pantalla → Cancela reserva, libera todo
- Si regresa después → NO debe tener reserva activa
```

---

## Implementación

### Backend: Validación de Duplicados

**Archivo:** `backend/internal/usecases/reservation_usecases.go`

**Problema Actual:**
```go
// Valida si CUALQUIERA tiene el número reservado
count, err := uc.reservationRepo.CountActiveReservationsForNumbers(...)
if count > 0 {
    return nil, ErrNumbersAlreadyReserved  // ❌ Error genérico
}
```

**Solución:**
```go
// Validar si ESTE USUARIO ya tiene reserva activa para este sorteo
existingReservation, err := uc.reservationRepo.FindActiveByUserAndRaffle(ctx, input.UserID, input.RaffleID)
if existingReservation != nil {
    // Usuario YA tiene reserva → retornar la existente (idempotencia)
    return existingReservation, nil
}

// Validar si OTROS usuarios tienen los números
count, err := uc.reservationRepo.CountActiveReservationsForNumbers(...)
if count > 0 {
    return nil, ErrNumbersAlreadyReserved  // ✅ Solo si OTROS los tienen
}
```

---

### Frontend: RaffleDetailPage

**Archivo:** `frontend/src/features/raffles/pages/RaffleDetailPage.tsx`

#### Cambios Principales

1. **Al montar componente: Cancelar reservas previas**
```typescript
useEffect(() => {
  const cancelPreviousReservation = async () => {
    const activeReservation = await reservationService.getActiveForRaffle(raffleId);
    if (activeReservation) {
      await reservationService.cancel(activeReservation.id);
    }
  };

  if (user && raffleId) {
    cancelPreviousReservation();
  }
}, [raffleId, user]);
```

2. **Al seleccionar PRIMER número: Crear reserva**
```typescript
const handleNumberToggle = async (number: RaffleNumber) => {
  const currentlySelected = selectedNumbers.find(n => n.id === number.id);

  if (currentlySelected) {
    // Deseleccionar (remover de reserva)
    await removeNumberFromReservation(number.id);
  } else {
    // Seleccionar
    const isFirstNumber = selectedNumbers.length === 0;

    if (isFirstNumber) {
      // CREAR NUEVA RESERVA con este primer número
      await createReservation([number.id]);
    } else {
      // AGREGAR a reserva existente
      await addNumberToReservation(number.id);
    }
  }
};
```

3. **"Proceder al Pago": Solo navegar**
```typescript
const handleProceedToCheckout = () => {
  if (!activeReservation) {
    toast.error('No tienes números reservados');
    return;
  }

  // Solo navegar (reserva ya existe)
  navigate('/checkout');
};
```

---

### Frontend: CheckoutPage

**Archivo:** `frontend/src/features/checkout/pages/CheckoutPage.tsx`

#### Simplificación Radical

```typescript
export function CheckoutPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { currentRaffleId } = useCartStore();

  // ✅ Cargar reserva del servidor (única fuente de verdad)
  const {
    reservation: activeReservation,
    isLoading,
  } = useActiveReservation(currentRaffleId || '');

  // ✅ Si no hay reserva → volver atrás
  useEffect(() => {
    if (!isLoading && !activeReservation) {
      toast.error('No tienes números reservados');
      navigate('/raffles');
    }
  }, [activeReservation, isLoading, navigate]);

  // ✅ NO crear reserva aquí
  // ✅ Solo mostrar botón de pago
  const handlePay = async () => {
    if (!activeReservation) return;

    const result = await createPaymentIntent({
      reservation_id: activeReservation.id,
      return_url: window.location.origin + '/payment/success',
      cancel_url: window.location.origin + '/checkout',
    });

    // Redirigir a PayPal
    window.location.href = result.payment_intent.client_secret;
  };

  // ... resto del componente (solo UI)
}
```

**Eliminar:**
- ❌ `handleCreateReservation` (no se usa)
- ❌ Estados `step: 'review' | 'reserving' | 'reserved'` (simplificar)
- ❌ Lógica de sincronización entre cart y server
- ❌ useCreateReservation mutation

---

### Frontend: Limpieza del CartStore

**Archivo:** `frontend/src/store/cartStore.ts`

#### Opción A: Mantener CartStore Simple
```typescript
interface CartStore {
  // Solo números seleccionados localmente
  currentRaffleId: string | null;
  selectedNumbers: CartNumber[];

  // Actions básicas
  toggleNumber: (number: CartNumber) => void;
  clearNumbers: () => void;

  // NO almacenar activeReservation aquí
  // (usar hook useActiveReservation en su lugar)
}
```

#### Opción B: Eliminar CartStore Completamente
- Usar solo `useActiveReservation` hook
- Estado local en RaffleDetailPage para selección temporal
- Reserva se crea inmediatamente al seleccionar primer número

---

## Endpoints Backend Necesarios

### 1. ✅ Ya Existe: Crear Reserva
```
POST /api/v1/reservations
Body: { raffle_id, number_ids[], session_id }
```

### 2. ✅ Ya Existe: Agregar Número
```
POST /api/v1/reservations/:id/add-number
Body: { number_id }
```

### 3. ✅ Ya Existe: Cancelar Reserva
```
POST /api/v1/reservations/:id/cancel
```

### 4. ✅ Ya Existe: Obtener Reserva Activa
```
GET /api/v1/raffles/:id/my-reservation
Response: { reservation } | 404
```

### 5. ⚠️ FALTA: Remover Número de Reserva
```
DELETE /api/v1/reservations/:id/numbers/:number_id
```

---

## Orden de Implementación

### Fase 1: Backend
1. ✅ Fix PostgreSQL array (ya completado)
2. Mejorar validación de duplicados en `CreateReservation`
3. Implementar endpoint `DELETE /reservations/:id/numbers/:number_id`

### Fase 2: Frontend - RaffleDetailPage
1. Agregar cancelación de reservas previas al montar
2. Implementar creación de reserva al primer número
3. Implementar agregar/remover números a reserva existente
4. Simplificar "Proceder al Pago" (solo navegar)

### Fase 3: Frontend - CheckoutPage
1. Eliminar lógica de creación de reserva
2. Simplificar a solo cargar reserva y pagar
3. Remover estados innecesarios

### Fase 4: Testing
1. Test: Entrar al sorteo sin reservas previas
2. Test: Seleccionar primer número → crea reserva
3. Test: Seleccionar más números → agrega a reserva
4. Test: Desseleccionar número → remueve de reserva
5. Test: Proceder al pago → navega sin errores
6. Test: Timeout de 10 minutos → libera todo
7. Test: Salir de la página → cancela reserva

---

## Riesgos y Consideraciones

### 1. Performance
- Crear reserva al primer click puede ser lento
- Solución: Mostrar loading spinner optimista

### 2. Race Conditions
- Dos usuarios seleccionan mismo número al mismo tiempo
- Solución: Redis locks (ya implementados)

### 3. UX: Feedback Inmediato
- Usuario debe ver confirmación visual inmediata
- Solución: Optimistic UI updates + rollback on error

### 4. Reservas Huérfanas
- Usuario selecciona número y cierra navegador
- Solución: Job de expiración cada 30 segundos (ya existe)

---

## Beneficios de Este Flujo

✅ **Simple y Claro**
- Un solo lugar crea reservas (RaffleDetailPage)
- CheckoutPage solo consume reserva existente

✅ **No Más 409 Errors**
- Backend valida correctamente si es el mismo usuario
- Frontend no intenta crear duplicados

✅ **Estado Consistente**
- Servidor es única fuente de verdad
- No hay sincronización compleja

✅ **UX Mejorado**
- Usuario ve números reservados inmediatamente
- Sin errores confusos
- Flujo lineal y predecible

---

## Preguntas para el Usuario

1. **¿Aprobar este plan?**
   - ¿Proceder con la implementación?

2. **CartStore:**
   - ¿Mantener simple o eliminar completamente?

3. **Optimistic UI:**
   - ¿Mostrar número seleccionado antes de confirmar reserva?
   - ¿O esperar respuesta del servidor?

4. **Endpoint de remover número:**
   - ¿Implementarlo o solo permitir cancelar reserva completa?

---

**Pendiente de Aprobación del Usuario**
