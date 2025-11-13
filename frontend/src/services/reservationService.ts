import { api } from '@/lib/api';

export interface Reservation {
  id: string;
  raffle_id: string;
  user_id: string;
  number_ids: string[];
  status: 'pending' | 'confirmed' | 'expired' | 'cancelled';
  phase: 'selection' | 'checkout' | 'completed' | 'expired';
  session_id: string;
  selection_started_at: string;
  checkout_started_at?: string;
  total_amount: number;
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface CreateReservationRequest {
  raffle_id: string;
  number_ids: string[];
  session_id: string;
}

export interface CreateReservationResponse {
  reservation: Reservation;
}

export const reservationService = {
  /**
   * Crear una nueva reserva
   */
  async create(data: CreateReservationRequest): Promise<Reservation> {
    const response = await api.post<CreateReservationResponse>(
      '/reservations',
      data
    );
    return response.data.reservation;
  },

  /**
   * Obtener una reserva por ID
   */
  async getById(id: string): Promise<Reservation> {
    const response = await api.get<{ reservation: Reservation }>(
      `/reservations/${id}`
    );
    return response.data.reservation;
  },

  /**
   * Obtener mis reservas
   */
  async getMyReservations(): Promise<Reservation[]> {
    const response = await api.get<{ reservations: Reservation[] }>(
      '/reservations/me'
    );
    return response.data.reservations;
  },

  /**
   * Mover reserva a fase de checkout
   */
  async moveToCheckout(id: string): Promise<Reservation> {
    const response = await api.post<{ reservation: Reservation }>(
      `/reservations/${id}/move-to-checkout`
    );
    return response.data.reservation;
  },

  /**
   * Cancelar una reserva
   */
  async cancel(id: string): Promise<void> {
    await api.post(`/reservations/${id}/cancel`);
  },

  /**
   * Confirmar una reserva (pago completado)
   */
  async confirm(id: string): Promise<void> {
    await api.post(`/reservations/${id}/confirm`);
  },

  /**
   * Agregar un número a una reserva existente
   */
  async addNumber(id: string, numberId: string): Promise<Reservation> {
    const response = await api.post<{ reservation: Reservation }>(
      `/reservations/${id}/add-number`,
      { number_id: numberId }
    );
    return response.data.reservation;
  },

  /**
   * Obtener reserva activa del usuario para un sorteo específico
   */
  async getActiveForRaffle(raffleId: string): Promise<Reservation | null> {
    try {
      const response = await api.get<{ success: boolean; data: Reservation }>(
        `/raffles/${raffleId}/my-reservation`
      );
      return response.data.data;
    } catch (error: any) {
      if (error.response?.status === 404) {
        return null; // No hay reserva activa
      }
      throw error;
    }
  },
};
