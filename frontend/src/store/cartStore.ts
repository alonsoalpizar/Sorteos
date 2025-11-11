import { create } from 'zustand';
import { persist } from 'zustand/middleware';
// Import Reservation type from API hook to ensure consistency
import type { Reservation } from '../hooks/useReservations';

export interface CartNumber {
  id: string; // e.g., "0001", "0042"
  displayNumber: string; // formatted for display
}

export interface CartItem {
  raffleId: string;
  raffleName: string;
  pricePerNumber: number;
  numbers: CartNumber[];
  imageUrl?: string;
}

interface CartStore {
  // Cart state
  currentRaffleId: string | null;
  selectedNumbers: CartNumber[];

  // Reservation state
  activeReservation: Reservation | null;
  reservationExpiry: Date | null;

  // Actions
  setCurrentRaffle: (raffleId: string) => void;
  addNumber: (number: CartNumber) => void;
  removeNumber: (numberId: string) => void;
  toggleNumber: (number: CartNumber) => void;
  clearNumbers: () => void;

  // Reservation actions
  setReservation: (reservation: Reservation) => void;
  clearReservation: () => void;
  isReservationActive: () => boolean;

  // Computed
  getSelectedCount: () => number;
  getSelectedNumberIds: () => string[];
  getTotalAmount: (pricePerNumber: number) => number;
}

export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      // Initial state
      currentRaffleId: null,
      selectedNumbers: [],
      activeReservation: null,
      reservationExpiry: null,

      // Cart actions
      setCurrentRaffle: (raffleId) => {
        const current = get().currentRaffleId;
        // If changing raffle, clear selected numbers
        if (current && current !== raffleId) {
          set({ currentRaffleId: raffleId, selectedNumbers: [] });
        } else {
          set({ currentRaffleId: raffleId });
        }
      },

      addNumber: (number) => {
        const { selectedNumbers } = get();
        // Check if already selected
        if (!selectedNumbers.find(n => n.id === number.id)) {
          set({ selectedNumbers: [...selectedNumbers, number] });
        }
      },

      removeNumber: (numberId) => {
        set(state => ({
          selectedNumbers: state.selectedNumbers.filter(n => n.id !== numberId)
        }));
      },

      toggleNumber: (number) => {
        const { selectedNumbers } = get();
        const exists = selectedNumbers.find(n => n.id === number.id);

        if (exists) {
          set({
            selectedNumbers: selectedNumbers.filter(n => n.id !== number.id)
          });
        } else {
          set({
            selectedNumbers: [...selectedNumbers, number]
          });
        }
      },

      clearNumbers: () => {
        set({ selectedNumbers: [] });
      },

      // Reservation actions
      setReservation: (reservation) => {
        set({
          activeReservation: reservation,
          reservationExpiry: new Date(reservation.expires_at),
          // Clear selected numbers once reserved
          selectedNumbers: []
        });
      },

      clearReservation: () => {
        set({
          activeReservation: null,
          reservationExpiry: null
        });
      },

      isReservationActive: () => {
        const { activeReservation, reservationExpiry } = get();

        if (!activeReservation || !reservationExpiry) {
          return false;
        }

        // Check if reservation is still pending and not expired
        const isExpired = new Date() > reservationExpiry;
        const isPending = activeReservation.status === 'pending';

        return isPending && !isExpired;
      },

      // Computed
      getSelectedCount: () => {
        return get().selectedNumbers.length;
      },

      getSelectedNumberIds: () => {
        return get().selectedNumbers.map(n => n.id);
      },

      getTotalAmount: (pricePerNumber) => {
        return get().selectedNumbers.length * pricePerNumber;
      },
    }),
    {
      name: 'sorteos-cart-storage', // localStorage key
      partialize: (state) => ({
        // Only persist these fields
        currentRaffleId: state.currentRaffleId,
        selectedNumbers: state.selectedNumbers,
        activeReservation: state.activeReservation,
        reservationExpiry: state.reservationExpiry,
      }),
    }
  )
);
