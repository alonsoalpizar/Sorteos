# Fix: Mensaje Confuso "Números Ya Reservados" (Reserva Propia)

**Fecha:** 2025-11-13
**Issue:** Usuario recibe mensaje confuso cuando regresa después de ir a checkout
**Prioridad:** Alta - Impacta UX negativamente

---

## 🐛 Problema Original

### Flujo que causaba confusión:

```
Usuario selecciona números
    ↓
Reserva creada (fase: selection)
    ↓
Usuario hace clic "Proceder al Pago"
    ↓
Va a página de checkout (fase: checkout)
    ↓
Usuario presiona ATRÁS en navegador
    ↓
Vuelve a RaffleDetailPage
    ↓
Frontend intenta crear NUEVA reserva con los mismos números
    ↓
Backend: "❌ Uno o más números ya están reservados"
    ↓
Usuario: "¿QUÉ? ¡Si son MÍOS!" 😠
```

### Por qué pasaba:

El frontend **no detectaba** que el usuario YA TENÍA una reserva activa para ese sorteo, e intentaba crear una nueva reserva con los mismos números que ya estaban reservados por él mismo.

---

## ✅ Solución Implementada

### Estrategia: **Opción 1 + Opción 3 Combinadas**

1. **Backend:** Endpoint para obtener reserva activa (ya existía) ✅
2. **Frontend:** Hook `useActiveReservation` para cargar reserva al montar ✅
3. **Frontend:** Banner de Alert mostrando reserva activa ✅
4. **Frontend:** Prevención de creación duplicada ✅

---

## 🔧 Cambios Implementados

### 1. Backend - Endpoint existente (sin cambios)

**Endpoint:**
```
GET /api/v1/raffles/:id/my-reservation
```

**Respuesta exitosa (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "raffle_id": "uuid",
    "user_id": "uuid",
    "number_ids": ["001", "042", "123"],
    "status": "pending",
    "phase": "checkout",
    "expires_at": "2025-11-13T03:20:00Z",
    ...
  }
}
```

**Respuesta sin reserva (404):**
```json
{
  "message": "no active reservation"
}
```

---

### 2. Frontend - Hook `useActiveReservation`

**Archivo:** `frontend/src/hooks/useActiveReservation.ts`

**Funcionalidad:**
- Carga automáticamente la reserva activa al montar
- Verifica si está expirada y la cancela automáticamente
- Proporciona métodos para crear y cancelar reservas
- Retorna estado de carga y datos de reserva

**Uso:**
```typescript
const {
  reservation: activeReservation,
  isLoading,
  createReservation,
  cancelReservation,
  refreshReservation,
} = useActiveReservation(raffleId);
```

---

### 3. Frontend - RaffleDetailPage con Banner

**Archivo:** `frontend/src/features/raffles/pages/RaffleDetailPage.tsx`

**Cambios principales:**

#### a) Importar hook y componentes
```typescript
import { useActiveReservation } from '../../../hooks/useActiveReservation';
import { Alert, AlertTitle, AlertDescription } from '../../../components/ui/Alert';
```

#### b) Usar el hook
```typescript
const {
  reservation: activeReservation,
  cancelReservation: cancelActiveReservation,
} = useActiveReservation(data?.raffle?.uuid || '');
```

#### c) Restaurar números al cargar reserva activa
```typescript
useEffect(() => {
  if (activeReservation && activeReservation.number_ids) {
    // Restaurar números al carrito
    clearNumbers();
    activeReservation.number_ids.forEach((numberId) => {
      toggleNumber({
        id: numberId,
        displayNumber: numberId,
      });
    });

    // Actualizar estado de reserva actual
    setCurrentReservation({
      id: activeReservation.id,
      expires_at: activeReservation.expires_at,
    });
  }
}, [activeReservation]);
```

#### d) Prevenir creación duplicada
```typescript
const createOrUpdateReservation = useCallback(async () => {
  // ... validaciones existentes

  // ✅ PREVENIR DUPLICADOS
  if (activeReservation) {
    toast.warning('Ya tienes números reservados', {
      description: 'Cancela tu reserva actual primero si quieres seleccionar otros números',
    });
    return;
  }

  // ... resto de la lógica
}, [..., activeReservation]);
```

#### e) Banner de Alert
```tsx
{/* Active Reservation Banner */}
{activeReservation && !isOwner && (
  <Alert variant="info" className="border-blue-500">
    <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    <AlertTitle>Tienes una reserva activa</AlertTitle>
    <AlertDescription>
      <div className="space-y-3">
        <p className="text-sm">
          Has reservado <strong>{activeReservation.number_ids?.length || 0}</strong> números para este sorteo.
        </p>

        <div className="text-sm">
          <strong>Números reservados:</strong> {activeReservation.number_ids?.join(', ')}
        </div>

        <div className="flex flex-wrap gap-2 pt-2">
          <Button
            size="sm"
            onClick={() => navigate('/checkout')}
            className="bg-blue-600 hover:bg-blue-700 text-white"
          >
            <svg className="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
            </svg>
            Ir a Pagar
          </Button>

          <Button
            size="sm"
            variant="outline"
            onClick={async () => {
              if (confirm('¿Seguro que quieres cancelar tu reserva? Los números quedarán disponibles de nuevo.')) {
                await cancelActiveReservation();
                clearNumbers();
                toast.success('Reserva cancelada', {
                  description: 'Los números están disponibles de nuevo',
                });
              }
            }}
            className="border-red-500 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/10"
          >
            Cancelar Reserva
          </Button>
        </div>
      </div>
    </AlertDescription>
  </Alert>
)}
```

---

### 4. Frontend - Tipo actualizado

**Archivo:** `frontend/src/services/reservationService.ts`

**Cambio:**
```typescript
export interface CreateReservationRequest {
  raffle_id: string;
  number_ids: string[];
  session_id: string;  // ✅ Agregado
}
```

**Motivo:** El hook `useActiveReservation` necesita pasar `session_id` al crear reservas.

---

## 🎯 Flujo Corregido

### Escenario: Usuario regresa después de ir a checkout

```
1. Usuario entra al sorteo
    ↓
2. Frontend ejecuta useActiveReservation hook
    ↓
3. Hook hace GET /raffles/:id/my-reservation
    ↓
4. Backend responde: "Sí, tienes reserva activa"
    ↓
5. Frontend:
   ├─ Restaura números seleccionados al carrito
   ├─ Actualiza estado de currentReservation
   └─ Muestra banner de Alert
    ↓
6. Usuario ve:
   ✅ Banner azul: "Tienes una reserva activa"
   ✅ Sus números TODAVÍA seleccionados en el grid
   ✅ Lista de números reservados
   ✅ Botón "IR A PAGAR"
   ✅ Botón "CANCELAR RESERVA"
    ↓
7. Si usuario intenta seleccionar más números:
   ⚠️  Toast: "Ya tienes números reservados"
   ⚠️  No permite crear nueva reserva
    ↓
8. ¡Sin mensajes confusos! 🎉
```

---

## 📊 Casos de Uso Cubiertos

### ✅ Caso 1: Usuario regresa desde checkout

**Antes:**
- ❌ Mensaje: "Números ya reservados"
- ❌ Usuario confundido

**Ahora:**
- ✅ Banner claro: "Tienes una reserva activa"
- ✅ Números restaurados automáticamente
- ✅ Botón "Ir a Pagar" prominente

---

### ✅ Caso 2: Reserva expirada cuando usuario regresa

**Antes:**
- ❌ Números aparecían seleccionados pero no se podían reservar

**Ahora:**
- ✅ Hook detecta expiración automáticamente
- ✅ Cancela reserva expirada en backend
- ✅ Limpia selección en frontend
- ✅ Usuario puede seleccionar números nuevamente

---

### ✅ Caso 3: Usuario intenta seleccionar más números teniendo una reserva

**Antes:**
- ❌ Permitía seleccionar pero fallaba al crear reserva
- ❌ Mensaje confuso: "Números ya reservados"

**Ahora:**
- ✅ Detecta reserva activa ANTES de intentar crear nueva
- ✅ Toast informativo: "Ya tienes números reservados"
- ✅ Sugiere cancelar reserva actual primero

---

### ✅ Caso 4: Usuario cancela reserva desde banner

**Antes:**
- ❌ No había forma fácil de cancelar desde la página del sorteo

**Ahora:**
- ✅ Botón "Cancelar Reserva" visible en el banner
- ✅ Confirmación antes de cancelar
- ✅ Limpia carrito automáticamente
- ✅ Toast de éxito
- ✅ Usuario puede seleccionar nuevos números

---

## 🎨 UI/UX Mejorado

### Banner de Reserva Activa

**Diseño:**
- 🎨 Color azul (variante `info`)
- 🕐 Icono de reloj
- 📝 Título claro: "Tienes una reserva activa"
- 🔢 Muestra cantidad y lista de números
- 🎯 Dos acciones principales:
  - **Ir a Pagar** (azul, prominente)
  - **Cancelar Reserva** (rojo, outline)

**Posición:**
- Entre el botón "Volver al listado" y la sección Hero
- Visible inmediatamente al cargar la página
- No se puede perder de vista

---

## 🧪 Testing Recomendado

### Checklist Manual

- [ ] **Test 1:** Crear reserva → Ir a checkout → Volver atrás
  - Verificar: Banner aparece
  - Verificar: Números están seleccionados
  - Verificar: No mensaje de error

- [ ] **Test 2:** Tener reserva activa → Intentar seleccionar más números
  - Verificar: Toast de advertencia
  - Verificar: No permite crear nueva reserva

- [ ] **Test 3:** Tener reserva activa → Hacer clic "Cancelar Reserva"
  - Verificar: Confirmación aparece
  - Verificar: Reserva se cancela
  - Verificar: Números se liberan
  - Verificar: Banner desaparece

- [ ] **Test 4:** Tener reserva expirada → Volver a la página
  - Verificar: No aparece banner
  - Verificar: Números no están seleccionados
  - Verificar: Puede seleccionar nuevos números

- [ ] **Test 5:** Tener reserva activa → Hacer clic "Ir a Pagar"
  - Verificar: Navega a /checkout
  - Verificar: Reserva sigue activa

---

## 📝 Archivos Modificados

### Backend
- ✅ (Sin cambios - endpoint ya existía)

### Frontend
1. **`frontend/src/hooks/useActiveReservation.ts`**
   - Actualizado signature de `createReservation` para aceptar `sessionId`

2. **`frontend/src/services/reservationService.ts`**
   - Agregado `session_id` a interfaz `CreateReservationRequest`

3. **`frontend/src/features/raffles/pages/RaffleDetailPage.tsx`**
   - Importado `useActiveReservation` y componentes Alert
   - Agregado lógica de restauración de reserva activa
   - Agregado prevención de duplicados en `createOrUpdateReservation`
   - Agregado banner de Alert para mostrar reserva activa

---

## 🚀 Despliegue

**Comandos ejecutados:**
```bash
# 1. Compilar frontend localmente (verificación)
cd /opt/Sorteos/frontend && npm run build

# 2. Rebuild y reiniciar Docker
cd /opt/Sorteos
docker compose build api && docker compose up -d api

# 3. Verificar logs
docker logs sorteos-api --tail 40

# 4. Health checks
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

**Estado:** ✅ Desplegado exitosamente

---

## 📊 Impacto

### Antes:
- ❌ Mensaje confuso: "Números ya reservados"
- ❌ Usuario no sabía que eran SUS números
- ❌ No había forma fácil de continuar
- ❌ UX frustrante

### Ahora:
- ✅ Mensaje claro: "Tienes una reserva activa"
- ✅ Usuario ve exactamente qué números tiene reservados
- ✅ Acceso directo a checkout
- ✅ Opción de cancelar si cambió de opinión
- ✅ UX fluida y sin confusión

---

## 🔮 Mejoras Futuras (Opcional)

1. **Timer visual en el banner**
   - Mostrar countdown de tiempo restante
   - Cambiar color cuando quedan < 2 minutos

2. **Auto-refresh al expirar**
   - Actualizar banner automáticamente cuando expire
   - No requerir refresh manual

3. **Animación al restaurar números**
   - Highlight visual en números restaurados
   - Indicar que fueron cargados automáticamente

4. **Historial de reservas**
   - Mostrar reservas previas (expiradas/canceladas)
   - Permitir "re-reservar" mismos números

---

**Última actualización:** 2025-11-13 03:05 UTC
**Versión:** 1.3 - Fix mensaje confuso implementado
**Estado:** ✅ Completado y desplegado
