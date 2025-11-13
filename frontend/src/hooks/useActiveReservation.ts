import { useEffect, useState, useCallback } from 'react';
import { reservationService, Reservation } from '@/services/reservationService';
import { toast } from 'sonner';

interface UseActiveReservationReturn {
  reservation: Reservation | null;
  isLoading: boolean;
  createReservation: (numberIds: string[], sessionId: string) => Promise<Reservation>;
  cancelReservation: () => Promise<void>;
  refreshReservation: () => Promise<void>;
}

export function useActiveReservation(raffleId: string): UseActiveReservationReturn {
  const [reservation, setReservation] = useState<Reservation | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const loadActiveReservation = useCallback(async () => {
    // No hacer petición si raffleId está vacío
    if (!raffleId) {
      setIsLoading(false);
      setReservation(null);
      return;
    }

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

  const createReservation = async (
    numberIds: string[],
    sessionId: string
  ): Promise<Reservation> => {
    const newReservation = await reservationService.create({
      raffle_id: raffleId,
      number_ids: numberIds,
      session_id: sessionId,
    });

    setReservation(newReservation);
    return newReservation;
  };

  const cancelReservation = async (): Promise<void> => {
    if (!reservation) return;

    try {
      await reservationService.cancel(reservation.id);
      setReservation(null);
      toast.success('Reserva cancelada', {
        description: 'Los números han sido liberados',
      });
    } catch (error) {
      toast.error('Error al cancelar la reserva', {
        description: 'No se pudo cancelar la reserva',
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
