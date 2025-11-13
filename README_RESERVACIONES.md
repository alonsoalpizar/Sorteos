# Sistema de Reservaciones con Doble Timeout

Este documento describe el sistema de reservaciones en tiempo real implementado para la plataforma de sorteos, que incluye WebSocket, doble timeout (selección + checkout), y concurrencia con Redis.

## Características Principales

- **Doble Timeout**: 10 minutos para selección de números + 5 minutos para checkout
- **Actualizaciones en Tiempo Real**: WebSocket para sincronización instantánea entre usuarios
- **Concurrencia Zero**: Redis distributed locks + PostgreSQL transactions
- **Reconexión Automática**: Frontend reconecta automáticamente con backoff exponencial
- **Proveedores de Pago Modulares**: Sistema flexible para habilitar/deshabilitar proveedores

## Arquitectura

### Backend (Go)

```
backend/
├── internal/
│   ├── domain/
│   │   └── entities/
│   │       └── reservation.go          # Entidad con lógica de doble timeout
│   ├── usecases/
│   │   └── reservation_usecases.go     # Lógica de negocio + WebSocket
│   ├── infrastructure/
│   │   ├── websocket/
│   │   │   ├── hub.go                  # Hub para gestionar conexiones
│   │   │   ├── client.go               # Cliente WebSocket individual
│   │   │   └── message.go              # Tipos de mensajes
│   │   └── redis/
│   │       └── lock_service.go         # Distributed locks
│   └── adapters/
│       └── http/
│           └── handler/
│               └── websocket/
│                   └── websocket_handler.go
└── migrations/
    └── 009_enhance_reservations_double_timeout.up.sql
```

### Frontend (React + TypeScript)

```
frontend/
└── src/
    ├── hooks/
    │   ├── useRaffleWebSocket.ts       # Hook de conexión WebSocket
    │   └── useTimeRemaining.ts         # Hook de countdown timer
    ├── components/
    │   ├── NumberGrid.tsx              # Grid de números con estados
    │   ├── ReservationTimer.tsx        # Timer visual con urgencia
    │   └── ui/
    │       └── FloatingCheckoutButton.tsx
    └── services/
        └── reservationService.ts       # API client para reservaciones
```

## Flujo de Usuario

### 1. Selección de Números (10 minutos)

```typescript
// Usuario selecciona números en el grid
const { onNumberUpdate } = useRaffleWebSocket(raffleId);

onNumberUpdate((data) => {
  // Actualizar estado del número en tiempo real
  updateNumberStatus(data.number_id, data.status);
});
```

**Acciones disponibles**:
- Seleccionar números disponibles (hasta 10)
- Deseleccionar números propios
- Ver actualizaciones de otros usuarios en tiempo real
- Cancelar reserva

### 2. Checkout (5 minutos)

```typescript
// Mover a fase de checkout
await reservationService.moveToCheckout(reservationId);

// Timer se actualiza automáticamente a 5 minutos
<ReservationTimer
  expiresAt={reservation.expires_at}
  phase="checkout"
  onExpire={() => handleExpiration()}
/>
```

**Acciones disponibles**:
- Seleccionar proveedor de pago habilitado
- Completar pago (números quedan como "sold")
- Cancelar (números vuelven a "available")

### 3. Expiración Automática

Si el usuario no completa en el tiempo límite:
- Backend ejecuta job cada 30 segundos
- Detecta reservaciones expiradas
- Libera números automáticamente
- Notifica vía WebSocket a todos los clientes

## API Endpoints

### Reservaciones

#### Crear Reservación
```bash
POST /api/v1/reservations
Content-Type: application/json

{
  "raffle_id": "uuid",
  "number_ids": ["uuid1", "uuid2"]
}

# Response: 201 Created
{
  "reservation": {
    "id": "uuid",
    "phase": "selection",
    "expires_at": "2025-11-12T10:10:00Z",
    "numbers": [...]
  }
}
```

#### Mover a Checkout
```bash
POST /api/v1/reservations/:id/move-to-checkout

# Response: 200 OK
{
  "reservation": {
    "id": "uuid",
    "phase": "checkout",
    "expires_at": "2025-11-12T10:05:00Z"  # +5 minutos desde ahora
  }
}
```

#### Agregar Número (en fase selección)
```bash
POST /api/v1/reservations/:id/add-number
Content-Type: application/json

{
  "number_id": "uuid"
}
```

#### Cancelar Reservación
```bash
POST /api/v1/reservations/:id/cancel

# Response: 200 OK
```

### WebSocket

#### Conectar
```javascript
const ws = new WebSocket('ws://62.171.188.255:8080/api/v1/raffles/:raffle_id/ws');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  switch(message.type) {
    case 'number_update':
      // message.data = { number_id, status, user_id? }
      break;
    case 'reservation_expired':
      // message.data = { reservation_id }
      break;
  }
};
```

#### Estadísticas (Admin)
```bash
GET /api/v1/raffles/:id/ws/stats
Authorization: Bearer <token>

# Response:
{
  "total_connections": 15,
  "connections_by_raffle": {
    "raffle-uuid-1": 10,
    "raffle-uuid-2": 5
  }
}
```

## Configuración

### Backend

**Variables de entorno** (`backend/.env`):
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=sorteos_user
DB_PASSWORD=your_password
DB_NAME=sorteos_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=your-secret-key
CORS_ORIGINS=http://localhost:5173,http://62.171.188.255:3000
```

**Ejecutar migraciones**:
```bash
cd backend
migrate -path migrations -database "postgresql://user:password@localhost:5432/sorteos_db?sslmode=disable" up
```

**Iniciar servidor**:
```bash
cd backend
go run cmd/api/main.go
```

### Frontend

**Variables de entorno** (`frontend/.env.local`):
```bash
# API URL (sin protocolo ws://)
VITE_API_URL=http://62.171.188.255:8080

# WebSocket URL se construye automáticamente desde VITE_API_URL
# ws:// para http:// y wss:// para https://
```

**Instalar dependencias**:
```bash
cd frontend
npm install
```

**Modo desarrollo**:
```bash
npm run dev
```

**Build producción**:
```bash
npm run build
```

## Uso en Componentes

### Ejemplo Completo: Página de Detalle de Rifa

```typescript
import { useRaffleWebSocket } from '@/hooks/useRaffleWebSocket';
import { NumberGrid } from '@/components/NumberGrid';
import { ReservationTimer } from '@/components/ReservationTimer';
import { FloatingCheckoutButton } from '@/components/ui/FloatingCheckoutButton';

export function RaffleDetailPage() {
  const { raffleId } = useParams();
  const [selectedNumbers, setSelectedNumbers] = useState<string[]>([]);
  const [reservation, setReservation] = useState<Reservation | null>(null);

  // Conectar WebSocket
  const { isConnected, onNumberUpdate } = useRaffleWebSocket(raffleId);

  // Manejar actualizaciones en tiempo real
  onNumberUpdate((data) => {
    setNumbers(prev =>
      prev.map(n => n.id === data.number_id
        ? { ...n, status: data.status, user_id: data.user_id }
        : n
      )
    );
  });

  // Crear reservación
  const handleCreateReservation = async () => {
    const res = await reservationService.create({
      raffle_id: raffleId,
      number_ids: selectedNumbers
    });
    setReservation(res);
  };

  // Ir a checkout
  const handleCheckout = async () => {
    if (!reservation) return;
    const updated = await reservationService.moveToCheckout(reservation.id);
    setReservation(updated);
    navigate(`/checkout/${reservation.id}`);
  };

  return (
    <div>
      {/* Timer de reservación */}
      {reservation && (
        <ReservationTimer
          expiresAt={reservation.expires_at}
          phase={reservation.phase}
          onExpire={() => setReservation(null)}
        />
      )}

      {/* Grid de números */}
      <NumberGrid
        numbers={numbers}
        selectedNumbers={selectedNumbers}
        onSelectNumber={(id) => {
          setSelectedNumbers(prev =>
            prev.includes(id)
              ? prev.filter(n => n !== id)
              : [...prev, id]
          );
        }}
        currentUserId={user?.id}
      />

      {/* Botón flotante */}
      <FloatingCheckoutButton
        selectedCount={selectedNumbers.length}
        totalAmount={selectedNumbers.length * raffle.ticket_price}
        onCheckout={handleCheckout}
        onCancel={() => setSelectedNumbers([])}
        disabled={!reservation}
      />
    </div>
  );
}
```

## Lógica de Concurrencia

### Redis Distributed Locks

```go
// Adquirir lock antes de reservar
lock, err := uc.lockService.AcquireLock(
    ctx,
    fmt.Sprintf("raffle:number:%s", numberID),
    10*time.Second, // TTL
)
if err != nil {
    return ErrNumberAlreadyReserved
}
defer lock.Release(ctx)

// Verificar disponibilidad en DB
number, err := uc.raffleRepo.GetNumberByID(ctx, numberID)
if number.Status != "available" {
    return ErrNumberAlreadyReserved
}

// Actualizar en transacción
tx.Begin()
tx.UpdateNumberStatus(numberID, "reserved")
tx.CreateReservation(...)
tx.Commit()

// Notificar vía WebSocket
uc.wsHub.BroadcastNumberUpdate(raffleID, numberID, "reserved", &userID)
```

### WebSocket Broadcast

```go
// Hub gestiona conexiones por raffle
type Hub struct {
    raffles map[string]map[*Client]bool  // raffleID -> clients
}

// Broadcast solo a clientes del raffle específico
func (h *Hub) BroadcastNumberUpdate(raffleID, numberID, status string) {
    h.mu.RLock()
    clients := h.raffles[raffleID]
    h.mu.RUnlock()

    message := &Message{
        Type: MessageTypeNumberUpdate,
        RaffleID: raffleID,
        Data: map[string]interface{}{
            "number_id": numberID,
            "status": status,
        },
    }

    for client := range clients {
        client.Send <- message
    }
}
```

## Jobs en Background

### Job de Expiración de Reservaciones

```go
// Ejecuta cada 30 segundos
func StartReservationExpirationJob(db *sql.DB, wsHub *websocket.Hub) {
    ticker := time.NewTicker(30 * time.Second)

    for range ticker.C {
        // Buscar reservaciones expiradas
        expired := findExpiredReservations(db)

        for _, res := range expired {
            // Liberar números
            releaseNumbers(db, res.ID)

            // Actualizar fase
            updateReservation(db, res.ID, "expired")

            // Notificar vía WebSocket
            wsHub.BroadcastReservationExpired(res.RaffleID, res.ID)

            for _, number := range res.Numbers {
                wsHub.BroadcastNumberUpdate(res.RaffleID, number.ID, "available", nil)
            }
        }
    }
}
```

## Testing

### Pruebas de Concurrencia

```bash
# Test de 1000 usuarios simultáneos
cd backend/tests
go test -v -run TestConcurrentReservations -count=1
```

### Pruebas de WebSocket

```bash
# Test de broadcast a 100 clientes
go test -v -run TestWebSocketBroadcast -count=1
```

### Pruebas de Locks

```bash
# Test de adquisición de locks Redis
go test -v -run TestRedisLocks -count=1
```

## Monitoreo

### Logs del Backend

```bash
# Ver logs de WebSocket
tail -f backend/logs/app.log | grep "WebSocket"

# Ver logs de reservaciones
tail -f backend/logs/app.log | grep "Reservation"
```

### Redis CLI

```bash
# Ver locks activos
redis-cli KEYS "raffle:number:*"

# Ver TTL de un lock
redis-cli TTL "raffle:number:uuid-123"
```

### PostgreSQL

```sql
-- Reservaciones activas por fase
SELECT phase, COUNT(*)
FROM reservations
WHERE phase IN ('selection', 'checkout')
GROUP BY phase;

-- Números reservados vs disponibles
SELECT status, COUNT(*)
FROM raffle_numbers
GROUP BY status;
```

## Troubleshooting

### WebSocket no conecta

1. Verificar CORS en backend:
```go
// cmd/api/main.go
config.AllowOrigins = []string{
    "http://localhost:5173",
    "http://62.171.188.255:3000"
}
```

2. Verificar URL en frontend:
```bash
# .env.local
VITE_API_URL=http://62.171.188.255:8080  # Sin ws://
```

3. Verificar Hub está corriendo:
```bash
# En logs del backend debe aparecer:
[INFO] WebSocket Hub initialized
```

### Números no se liberan al expirar

1. Verificar job está corriendo:
```bash
# En logs debe aparecer cada 30 segundos:
[INFO] Running reservation expiration job
```

2. Verificar zona horaria del servidor:
```sql
SELECT NOW(), expires_at FROM reservations WHERE id = 'uuid';
```

### Lock de Redis no se libera

Los locks tienen TTL automático de 10 segundos. Si persiste:

```bash
# Eliminar lock manualmente
redis-cli DEL "raffle:number:uuid-123"
```

## Próximas Mejoras

- [ ] Implementar página de checkout completa
- [ ] Agregar tests de integración frontend
- [ ] Implementar sistema de notificaciones push
- [ ] Agregar métricas de Prometheus
- [ ] Implementar circuit breaker para Redis
- [ ] Agregar soporte para múltiples idiomas
- [ ] Implementar modo offline con sincronización

## Soporte

Para dudas o problemas:
- Email: support@sorteos.com
- GitHub Issues: https://github.com/alonsoalpizar/Sorteos/issues
- Documentación: [/opt/Sorteos/Documentacion/](file:///opt/Sorteos/Documentacion/)
