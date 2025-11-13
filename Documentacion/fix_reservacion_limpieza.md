# Fix: Limpieza de Reservaciones No Completadas

## Problema Identificado

**Comportamiento actual (INCORRECTO):**
- Usuario reserva números [00, 01, 02]
- Usuario sale de la página sin pagar
- Números quedan reservados indefinidamente
- Usuario no puede volver a seleccionar esos números

**Impacto:**
- ❌ Números bloqueados sin venta real
- ❌ Mala experiencia de usuario
- ❌ Pérdida potencial de ventas

## Solución: 3 Capas de Protección

### 1. Job de Expiración Automática (Backend)
**Frecuencia:** Cada 30 segundos
**Acción:** Libera números de reservas expiradas

### 2. Cancelación al Salir (Frontend)
**Trigger:** Usuario cierra tab o sale de página
**Acción:** Cancela reserva automáticamente

### 3. Limpieza al Volver (Frontend)
**Trigger:** Usuario regresa al sorteo
**Acción:** Verifica y limpia reservas expiradas

## Implementación

### Backend: Job de Expiración

**Archivo:** `backend/cmd/api/jobs.go`

Agregar nueva función:

```go
// startReservationExpirationJob inicia el job de expiración de reservas
func startReservationExpirationJob(
    gormDB *gorm.DB,
    rdb *redis.Client,
    wsHub *websocket.Hub,
    cfg *config.Config,
    log *logger.Logger,
) {
    // Repositorios
    reservationRepo := db.NewReservationRepository(gormDB)
    numberRepo := db.NewRaffleNumberRepository(gormDB)

    // Lock service
    lockService := redisinfra.NewLockService(rdb)

    // Use case
    reservationUC := usecases.NewReservationUseCases(
        reservationRepo,
        db.NewRaffleRepository(gormDB),
        lockService,
        wsHub,
    )

    // Job goroutine
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        log.Info("Reservation expiration job started")

        for range ticker.C {
            ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)

            count, err := reservationUC.ExpireOldReservations(ctx)
            if err != nil {
                log.Error("Error expiring reservations", logger.Error(err))
            } else if count > 0 {
                log.Info("Expired reservations", logger.Int("count", count))
            }

            cancel()
        }
    }()
}
```

### Backend: Use Case de Expiración

**Archivo:** `backend/internal/usecases/reservation_usecases.go`

Agregar método:

```go
// ExpireOldReservations busca y expira reservas vencidas
func (uc *ReservationUseCases) ExpireOldReservations(ctx context.Context) (int, error) {
    // 1. Buscar reservas pendientes expiradas
    expiredReservations, err := uc.reservationRepo.FindExpired(ctx)
    if err != nil {
        return 0, fmt.Errorf("find expired: %w", err)
    }

    if len(expiredReservations) == 0 {
        return 0, nil
    }

    count := 0

    for _, reservation := range expiredReservations {
        // 2. Transacción para liberar números
        err := uc.reservationRepo.WithTransaction(ctx, func(txCtx context.Context) error {
            // Actualizar estado de reserva
            reservation.Status = entities.ReservationStatusExpired
            reservation.Phase = entities.ReservationPhaseExpired
            reservation.UpdatedAt = time.Now()

            if err := uc.reservationRepo.Update(txCtx, reservation); err != nil {
                return fmt.Errorf("update reservation: %w", err)
            }

            // Liberar cada número
            for _, numberID := range reservation.NumberIDs {
                if err := uc.raffleRepo.UpdateNumberStatus(txCtx, numberID, "available"); err != nil {
                    return fmt.Errorf("update number %s: %w", numberID, err)
                }
            }

            return nil
        })

        if err != nil {
            continue // Continuar con siguientes reservas
        }

        // 3. Notificar vía WebSocket
        for _, numberID := range reservation.NumberIDs {
            uc.wsHub.BroadcastNumberUpdate(
                reservation.RaffleID.String(),
                numberID,
                "available",
                nil, // Sin user_id (liberado)
            )
        }

        count++
    }

    return count, nil
}
```

### Backend: Repository Method

**Archivo:** `backend/internal/adapters/db/reservation_repository.go`

Agregar método:

```go
// FindExpired busca reservas pendientes que ya expiraron
func (r *ReservationRepository) FindExpired(ctx context.Context) ([]*entities.Reservation, error) {
    var reservations []*entities.Reservation

    err := r.db.WithContext(ctx).
        Where("status = ?", entities.ReservationStatusPending).
        Where("expires_at < ?", time.Now()).
        Order("expires_at ASC").
        Limit(100). // Procesar máximo 100 por ejecución
        Find(&reservations).
        Error

    if err != nil {
        return nil, fmt.Errorf("find expired reservations: %w", err)
    }

    return reservations, nil
}

// WithTransaction ejecuta una función dentro de una transacción
func (r *ReservationRepository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        txCtx := context.WithValue(ctx, "tx", tx)
        return fn(txCtx)
    })
}
```

### Frontend: Hook de Reserva Activa

**Archivo:** `frontend/src/hooks/useActiveReservation.ts`

```typescript
import { useEffect, useState, useCallback } from 'react';
import { reservationService, Reservation } from '@/services/reservationService';
import { useToast } from '@/components/ui/use-toast';

interface UseActiveReservationReturn {
  reservation: Reservation | null;
  isLoading: boolean;
  createReservation: (numberIds: string[]) => Promise<Reservation>;
  cancelReservation: () => Promise<void>;
  refreshReservation: () => Promise<void>;
}

export function useActiveReservation(raffleId: string): UseActiveReservationReturn {
  const [reservation, setReservation] = useState<Reservation | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const { toast } = useToast();

  const loadActiveReservation = useCallback(async () => {
    try {
      setIsLoading(true);
      const active = await reservationService.getActiveForRaffle(raffleId);

      if (active) {
        // Verificar si está expirada
        const isExpired = new Date(active.expires_at) < new Date();

        if (isExpired) {
          // Expirada: cancelar y limpiar
          await reservationService.cancel(active.id).catch(() => {});
          setReservation(null);
        } else {
          setReservation(active);
        }
      } else {
        setReservation(null);
      }
    } catch (error) {
      console.error('Error loading active reservation:', error);
      setReservation(null);
    } finally {
      setIsLoading(false);
    }
  }, [raffleId]);

  useEffect(() => {
    loadActiveReservation();
  }, [loadActiveReservation]);

  const createReservation = async (numberIds: string[]): Promise<Reservation> => {
    const newReservation = await reservationService.create({
      raffle_id: raffleId,
      number_ids: numberIds,
    });

    setReservation(newReservation);
    return newReservation;
  };

  const cancelReservation = async (): Promise<void> => {
    if (!reservation) return;

    try {
      await reservationService.cancel(reservation.id);
      setReservation(null);
      toast({
        title: 'Reserva cancelada',
        description: 'Los números han sido liberados',
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'No se pudo cancelar la reserva',
        variant: 'destructive',
      });
      throw error;
    }
  };

  const refreshReservation = async (): Promise<void> => {
    await loadActiveReservation();
  };

  return {
    reservation,
    isLoading,
    createReservation,
    cancelReservation,
    refreshReservation,
  };
}
```

### Frontend: Actualizar Servicio

**Archivo:** `frontend/src/services/reservationService.ts`

Agregar método:

```typescript
/**
 * Obtener reserva activa del usuario para un sorteo específico
 */
async getActiveForRaffle(raffleId: string): Promise<Reservation | null> {
  try {
    const response = await api.get<{ reservation: Reservation }>(
      `/raffles/${raffleId}/my-reservation`
    );
    return response.data.reservation;
  } catch (error: any) {
    if (error.response?.status === 404) {
      return null; // No hay reserva activa
    }
    throw error;
  }
},
```

### Frontend: Cleanup al Salir

**Archivo:** `frontend/src/features/raffles/pages/RaffleDetailPage.tsx`

Agregar lógica de cleanup:

```typescript
// Al final del componente, antes del return
useEffect(() => {
  // Cleanup: cancelar reserva si el usuario sale sin completar
  return () => {
    if (currentReservation &&
        currentReservation.phase === 'selection' &&
        currentReservation.status === 'pending') {
      // Usuario salió sin ir a checkout
      reservationService.cancel(currentReservation.id).catch(() => {
        // Ignorar errores en cleanup
      });
    }
  };
}, [currentReservation]);

// Warning al cerrar tab
useEffect(() => {
  const handleBeforeUnload = (e: BeforeUnloadEvent) => {
    if (currentReservation && currentReservation.phase === 'selection') {
      e.preventDefault();
      e.returnValue = '¿Seguro que quieres salir? Perderás tu reserva.';
    }
  };

  window.addEventListener('beforeunload', handleBeforeUnload);
  return () => window.removeEventListener('beforeunload', handleBeforeUnload);
}, [currentReservation]);
```

### Backend: Endpoint para Obtener Reserva Activa

**Archivo:** `backend/cmd/api/routes.go`

Agregar ruta:

```go
// En setupReservationAndPaymentRoutes, agregar:
router.GET("/api/v1/raffles/:id/my-reservation",
    authMiddleware.Authenticate(),
    func(c *gin.Context) {
        raffleID := c.Param("id")
        userIDInt, _ := middleware.GetUserID(c)
        userUUID, _ := getUserUUID(userRepo, userIDInt)

        reservation, err := reservationUseCases.GetActiveReservation(
            c.Request.Context(),
            userUUID,
            raffleID,
        )

        if err != nil || reservation == nil {
            c.JSON(http.StatusNotFound, gin.H{"message": "no active reservation"})
            return
        }

        // Verificar si expiró
        if reservation.IsExpired() {
            reservationUseCases.CancelReservation(c.Request.Context(), reservation.ID)
            c.JSON(http.StatusNotFound, gin.H{"message": "reservation expired"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"reservation": reservation})
    },
)
```

### Backend: Use Case Method

**Archivo:** `backend/internal/usecases/reservation_usecases.go`

```go
// GetActiveReservation obtiene la reserva activa de un usuario para un sorteo
func (uc *ReservationUseCases) GetActiveReservation(
    ctx context.Context,
    userID uuid.UUID,
    raffleID string,
) (*entities.Reservation, error) {
    raffleUUID, err := uuid.Parse(raffleID)
    if err != nil {
        return nil, fmt.Errorf("invalid raffle id: %w", err)
    }

    reservation, err := uc.reservationRepo.FindActiveByUserAndRaffle(ctx, userID, raffleUUID)
    if err != nil {
        return nil, fmt.Errorf("find active reservation: %w", err)
    }

    return reservation, nil
}
```

## Checklist de Implementación

### Backend
- [ ] Agregar `FindExpired()` al repository
- [ ] Agregar `ExpireOldReservations()` al use case
- [ ] Crear `startReservationExpirationJob()` en jobs.go
- [ ] Llamar job desde main.go
- [ ] Agregar endpoint `GET /raffles/:id/my-reservation`
- [ ] Agregar método `GetActiveReservation()` al use case
- [ ] Tests unitarios de expiración

### Frontend
- [ ] Crear hook `useActiveReservation`
- [ ] Agregar método `getActiveForRaffle()` al service
- [ ] Implementar cleanup en `useEffect`
- [ ] Agregar warning en `beforeunload`
- [ ] Callback `onExpire` en timer
- [ ] Tests de integración

### Testing
- [ ] Usuario sale sin pagar → números liberados
- [ ] Timer expira → job libera números
- [ ] Usuario vuelve → grid limpio
- [ ] WebSocket notifica cambios en tiempo real
- [ ] Múltiples tabs sincronizados

## Métricas de Éxito

- ✅ 0% de números bloqueados indefinidamente
- ✅ 100% de reservas expiradas liberadas en < 1 minuto
- ✅ UX mejorada: usuarios pueden volver a intentar
- ✅ Aumento en tasa de conversión

## Notas Técnicas

### Frecuencia del Job
30 segundos es un balance entre:
- Liberar números rápidamente (bueno para UX)
- No sobrecargar DB con queries frecuentes

### Límite de 100 Reservas por Ejecución
Previene queries muy pesadas si hay muchas expiradas.
Si hay > 100, se procesan en la siguiente ejecución.

### Cleanup en Frontend
El `useEffect` cleanup NO siempre se ejecuta (el navegador puede matar el proceso).
Por eso necesitamos el job de backend como respaldo.

### beforeunload Warning
Solo funciona si el usuario interactuó con la página (seguridad del navegador).

---

**Prioridad:** ALTA
**Complejidad:** Media
**Tiempo estimado:** 2-3 horas
