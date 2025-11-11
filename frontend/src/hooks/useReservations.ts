import { useMutation, useQuery } from '@tanstack/react-query';
import { apiClient } from '../lib/apiClient';

export interface CreateReservationInput {
  raffle_id: string;
  number_ids: string[];
  session_id: string;
}

export interface Reservation {
  id: string;
  raffle_id: string;
  user_id: string;
  number_ids: string[];
  status: 'pending' | 'confirmed' | 'expired' | 'cancelled';
  session_id: string;
  total_amount: number;
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface CreateReservationResponse {
  reservation: Reservation;
}

export function useCreateReservation() {
  return useMutation({
    mutationFn: async (input: CreateReservationInput): Promise<CreateReservationResponse> => {
      // Generate session ID if not provided (for idempotency)
      const sessionId = input.session_id || `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

      const { data } = await apiClient.post<CreateReservationResponse>('/reservations', {
        raffle_id: input.raffle_id,
        number_ids: input.number_ids,
        session_id: sessionId,
      });

      return data;
    },
  });
}

export function useGetReservation(reservationId: string | null) {
  return useQuery({
    queryKey: ['reservation', reservationId],
    queryFn: async (): Promise<Reservation> => {
      if (!reservationId) throw new Error('No reservation ID provided');

      const { data } = await apiClient.get<{ reservation: Reservation }>(`/reservations/${reservationId}`);
      return data.reservation;
    },
    enabled: !!reservationId,
    refetchInterval: (query) => {
      // Refetch every 5 seconds if reservation is still pending
      if (query.state.data?.status === 'pending') {
        return 5000;
      }
      return false;
    },
  });
}

export function useGetMyReservations() {
  return useQuery({
    queryKey: ['myReservations'],
    queryFn: async (): Promise<Reservation[]> => {
      const { data } = await apiClient.get<{ reservations: Reservation[] }>('/reservations/me');
      return data.reservations;
    },
  });
}
