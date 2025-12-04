import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { rafflesApi } from '../api/raffles';
import type {
  CreateRaffleInput,
  UpdateRaffleInput,
  RaffleFilters,
} from '../types/raffle';

// Query keys
export const raffleKeys = {
  all: ['raffles'] as const,
  lists: () => [...raffleKeys.all, 'list'] as const,
  list: (filters?: RaffleFilters) => [...raffleKeys.lists(), filters] as const,
  details: () => [...raffleKeys.all, 'detail'] as const,
  detail: (id: number | string) => [...raffleKeys.details(), id] as const,
  myTickets: (page?: number) => [...raffleKeys.all, 'my-tickets', page] as const,
  buyers: (raffleId: string) => [...raffleKeys.all, 'buyers', raffleId] as const,
};

/**
 * Hook para listar sorteos
 */
export function useRafflesList(filters?: RaffleFilters) {
  return useQuery({
    queryKey: raffleKeys.list(filters),
    queryFn: () => rafflesApi.list(filters),
  });
}

/**
 * Hook para obtener detalle de sorteo
 */
export function useRaffleDetail(
  id: number | string,
  options?: { includeNumbers?: boolean; includeImages?: boolean }
) {
  return useQuery({
    queryKey: raffleKeys.detail(id),
    queryFn: () => rafflesApi.getDetail(id, options),
    enabled: !!id,
  });
}

/**
 * Hook para crear sorteo
 */
export function useCreateRaffle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateRaffleInput) => rafflesApi.create(input),
    onSuccess: () => {
      // Invalidar lista de sorteos
      queryClient.invalidateQueries({ queryKey: raffleKeys.lists() });
    },
  });
}

/**
 * Hook para actualizar sorteo
 */
export function useUpdateRaffle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateRaffleInput }) =>
      rafflesApi.update(id, input),
    onSuccess: (_, variables) => {
      // Invalidar detalle y lista
      queryClient.invalidateQueries({ queryKey: raffleKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: raffleKeys.lists() });
    },
  });
}

/**
 * Hook para publicar sorteo
 */
export function usePublishRaffle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => rafflesApi.publish(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: raffleKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: raffleKeys.lists() });
    },
  });
}

/**
 * Hook para eliminar sorteo
 */
export function useDeleteRaffle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => rafflesApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: raffleKeys.lists() });
    },
  });
}

/**
 * Hook para suspender sorteo (admin only)
 */
export function useSuspendRaffle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason: string }) =>
      rafflesApi.suspend(id, reason),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: raffleKeys.detail(variables.id) });
      queryClient.invalidateQueries({ queryKey: raffleKeys.lists() });
    },
  });
}

/**
 * Hook para obtener tickets del usuario autenticado
 */
export function useMyTickets(page: number = 1, pageSize: number = 20) {
  return useQuery({
    queryKey: raffleKeys.myTickets(page),
    queryFn: () => rafflesApi.getMyTickets(page, pageSize),
  });
}

/**
 * Hook para obtener lista de compradores de un sorteo (solo owner)
 */
export function useRaffleBuyers(
  raffleId: string,
  options?: { includeSold?: boolean; includeReserved?: boolean; enabled?: boolean }
) {
  return useQuery({
    queryKey: raffleKeys.buyers(raffleId),
    queryFn: () => rafflesApi.getBuyers(raffleId, {
      includeSold: options?.includeSold ?? true,
      includeReserved: options?.includeReserved ?? true,
    }),
    enabled: options?.enabled !== false && !!raffleId,
  });
}
