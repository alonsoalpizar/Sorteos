# 🎯 Floating Checkout Button + Auto-Reserva

## Descripción

Sistema de **burbuja flotante** (Floating Action Button - FAB) con **reserva automática** de números y **timer de expiración** para mejorar la experiencia de compra en sorteos.

---

## 🎨 Características Implementadas

### 1. **Floating Action Button (FAB)**
Botón flotante que aparece cuando el usuario selecciona números en un sorteo.

**Ubicación**: Esquina inferior derecha de la pantalla
**Comportamiento**:
- ✅ Sigue al usuario mientras hace scroll
- ✅ Aparece con animación de escala y fade-in
- ✅ Tiene un anillo pulsante (ping effect) para llamar la atención
- ✅ Se oculta automáticamente cuando no hay números seleccionados

**Contenido de la burbuja:**
```
┌────────────────────────────┐
│  [X]                       │  ← Botón cerrar
│                            │
│  🛒  Números seleccionados │
│      3                     │
│                            │
│  Total a pagar             │
│  ₡15,000                   │
│                            │
│  🕐 Reservado por 14:35    │  ← Timer
│                            │
│  [ Proceder al Pago ]      │
│                            │
│  Limpiar selección         │
└────────────────────────────┘
```

---

### 2. **Reserva Automática**
Los números se reservan **automáticamente** 500ms después de la última selección.

**Flujo:**
1. Usuario selecciona número(s)
2. Se espera 500ms (debounce)
3. Se crea reserva automática en el backend
4. Se muestra notificación toast: "3 número(s) reservado(s) por 15 minutos"
5. Los números quedan bloqueados para otros usuarios

**Ventajas:**
- ✅ No requiere botón "Reservar"
- ✅ Protege la selección del usuario inmediatamente
- ✅ Evita que otros usuarios tomen los mismos números
- ✅ UX más fluida y natural

---

### 3. **Timer de Expiración**
Contador regresivo que muestra el tiempo restante de la reserva.

**Estados:**

**a) Normal (> 2 minutos restantes):**
```
🕐 Reservado por 14:35
```

**b) Advertencia (< 2 minutos):**
```
┌─────────────────────────────┐
│ 🕐 ¡Reserva expira en 1:45! │  ← Burbuja amarilla
└─────────────────────────────┘
      ↓ (animación bounce)
[ Proceder al Pago ]
```

**c) Expirado (0:00):**
```
┌─────────────────────────┐
│    ❌ Reserva expirada  │  ← Botón rojo
└─────────────────────────┘
```

**Comportamiento:**
- Timer actualiza cada segundo
- Alerta visual y animación bounce cuando quedan < 2 min
- Botón se deshabilita cuando expira
- Números se liberan automáticamente en el backend

---

## 📱 Responsive Design

### Desktop (> 1024px)
- FAB en esquina inferior derecha
- Sin overlay de fondo
- Width: 280px

### Mobile (< 1024px)
- FAB en esquina inferior derecha
- Overlay semi-transparente con blur
- Tap en overlay cierra el FAB
- FAB ocupa ancho completo en pantallas pequeñas

---

## 🔧 Implementación Técnica

### Archivos Creados

**1. `FloatingCheckoutButton.tsx`**
```typescript
interface FloatingCheckoutButtonProps {
  selectedCount: number;        // Cantidad de números seleccionados
  totalAmount: number;          // Monto total a pagar
  expiresAt?: string | null;    // Fecha de expiración (ISO 8601)
  onCheckout: () => void;       // Callback al hacer checkout
  onClear: () => void;          // Callback al limpiar
  show: boolean;                // Mostrar/ocultar FAB
}
```

**Características:**
- ✅ Timer con `useState` + `useEffect`
- ✅ Cálculo de tiempo restante en segundos
- ✅ Animaciones con Tailwind CSS
- ✅ Toast notifications con `sonner`
- ✅ Iconos con `lucide-react`

### Archivos Modificados

**2. `RaffleDetailPage.tsx`**

**Nuevo estado:**
```typescript
const [currentReservation, setCurrentReservation] = useState<{
  id: string;
  expires_at: string;
} | null>(null);

const [sessionId] = useState(() =>
  `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
);
```

**Auto-reserva con debounce:**
```typescript
useEffect(() => {
  if (getSelectedCount() === 0) {
    setCurrentReservation(null);
    return;
  }

  const timer = setTimeout(() => {
    createOrUpdateReservation();
  }, 500); // Espera 500ms después de la última selección

  return () => clearTimeout(timer);
}, [getSelectedCount, createOrUpdateReservation]);
```

**Renderizado del FAB:**
```typescript
{!isOwner && raffle.status === 'active' && (
  <FloatingCheckoutButton
    selectedCount={getSelectedCount()}
    totalAmount={getTotalAmount(Number(raffle.price_per_number))}
    expiresAt={currentReservation?.expires_at}
    onCheckout={handleProceedToCheckout}
    onClear={clearNumbers}
    show={getSelectedCount() > 0}
  />
)}
```

---

## 🎯 Casos de Uso

### Escenario 1: Usuario selecciona números rápidamente
```
1. Click en número 42
2. Click en número 57
3. Click en número 89
   ↓ (espera 500ms)
4. ✅ Se crea reserva automática
5. 🔔 Toast: "3 número(s) reservado(s) por 15 minutos"
6. 🎈 FAB aparece con timer: "14:59"
```

### Escenario 2: Usuario deja la página abierta
```
1. Usuario tiene 3 números seleccionados
2. Timer: 14:59 → 14:58 → ... → 2:00
3. ⚠️ Alerta amarilla bouncing: "¡Reserva expira en 1:59!"
4. Timer: 1:59 → 1:58 → ... → 0:00
5. ❌ Botón se pone rojo: "Reserva expirada"
6. Backend libera los números automáticamente
```

### Escenario 3: Usuario cambia de opinión
```
1. Usuario tiene números seleccionados
2. Click en "Limpiar selección" o "X"
3. FAB desaparece con animación
4. Reserva se cancela (opcional: implementar cancelación explícita)
```

---

## 🚀 Mejoras Futuras (Opcionales)

### 1. Vibración en Mobile
```typescript
// Cuando quedan < 1 minuto
if ('vibrate' in navigator) {
  navigator.vibrate([200, 100, 200]);
}
```

### 2. Sonido de Advertencia
```typescript
const audio = new Audio('/sounds/alert.mp3');
if (timeLeft === 60) audio.play();
```

### 3. Persistencia en LocalStorage
```typescript
// Guardar reserva en localStorage
localStorage.setItem('pending_reservation', JSON.stringify({
  reservation_id: currentReservation.id,
  expires_at: currentReservation.expires_at,
}));

// Recuperar al volver a la página
```

### 4. Extensión de Tiempo
```typescript
// Botón para extender reserva
<button onClick={extendReservation}>
  + 5 minutos más
</button>
```

### 5. Multi-selección Rápida
```typescript
// Selección por rango
"Números 10-20" → Selecciona todos
```

---

## 🎨 Design Tokens

### Colores
```css
/* Normal */
background: gradient-to-br from-primary-600 to-primary-700
text: white

/* Alerta (< 2 min) */
background: yellow-500
text: white
animation: bounce

/* Expirado */
background: red-500
text: white
```

### Animaciones
```css
/* Entrada */
@keyframes fadeIn {
  from: opacity-0, scale-95, translateY(20px)
  to: opacity-100, scale-100, translateY(0)
}

/* Ping effect */
@keyframes ping {
  75%, 100%: opacity-0, scale-2
}

/* Bounce (advertencia) */
@keyframes bounce {
  0%, 100%: translateY(0)
  50%: translateY(-10px)
}
```

---

## 📊 Métricas de UX

### Antes (sin FAB)
- Usuario debe scrollear hacia arriba para ver botón "Proceder al Pago"
- No hay feedback inmediato de reserva
- Usuario puede perder números si otro compra primero

### Después (con FAB)
- ✅ Botón siempre visible (sticky)
- ✅ Feedback inmediato con toast
- ✅ Números protegidos automáticamente
- ✅ Timer visible reduce ansiedad
- ✅ UX más profesional y moderna

---

## 🧪 Testing

### Test Manual
1. Abrir sorteo activo
2. Seleccionar 1 número
3. ✅ Verificar que FAB aparece en < 500ms
4. ✅ Verificar toast de confirmación
5. ✅ Verificar timer contando hacia atrás
6. Seleccionar más números
7. ✅ Verificar que cantidad se actualiza
8. Esperar hasta < 2 minutos
9. ✅ Verificar alerta amarilla bouncing
10. Click en "Limpiar"
11. ✅ Verificar que FAB desaparece

### Test en Mobile
1. Abrir en móvil
2. Seleccionar números
3. ✅ Verificar overlay con blur
4. ✅ Tap fuera del FAB lo cierra
5. ✅ FAB responsive al ancho de pantalla

---

## 📝 Notas Técnicas

### Debounce
Se usa `setTimeout` con 500ms para evitar crear múltiples reservas mientras el usuario selecciona rápido.

### Session ID
Se genera un ID único por sesión para idempotencia en el backend. Si se envía la misma petición 2 veces, el backend sabe que es la misma reserva.

### Cleanup
El `useEffect` limpia el timer cuando el componente se desmonta para evitar memory leaks:
```typescript
return () => clearTimeout(timer);
```

### Toast Library
Usamos `sonner` (ya instalado) para notificaciones elegantes y no intrusivas.

---

## ✅ Estado Actual

- ✅ Componente `FloatingCheckoutButton` creado
- ✅ Auto-reserva implementada
- ✅ Timer de expiración funcionando
- ✅ Alertas visuales (amarillo < 2min, rojo = expirado)
- ✅ Animaciones y transiciones
- ✅ Responsive design
- ✅ Toast notifications
- ✅ Compilado y desplegado

**URL de prueba**: https://sorteos.club/raffles/1
(o cualquier sorteo activo)
