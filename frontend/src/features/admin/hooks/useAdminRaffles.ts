import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { adminRafflesApi } from "../api/adminApi";
import type {
  RaffleFilters,
  PaginationParams,
  ForceStatusChangeRequest,
  AddAdminNotesRequest,
  ManualDrawRequest,
  CancelWithRefundRequest,
} from "../types";

// Query keys centralizados
export const adminRaffleKeys = {
  all: ["admin", "raffles"] as const,
  lists: () => [...adminRaffleKeys.all, "list"] as const,
  list: (filters?: RaffleFilters, pagination?: PaginationParams) =>
    [...adminRaffleKeys.lists(), { filters, pagination }] as const,
  details: () => [...adminRaffleKeys.all, "detail"] as const,
  detail: (id: number) => [...adminRaffleKeys.details(), id] as const,
};

// ===========================
// Queries
// ===========================

export function useAdminRaffles(
  filters?: RaffleFilters,
  pagination?: PaginationParams
) {
  return useQuery({
    queryKey: adminRaffleKeys.list(filters, pagination),
    queryFn: () => adminRafflesApi.list(filters, pagination),
    staleTime: 30 * 1000, // 30 segundos
  });
}

export function useAdminRaffleDetail(raffleId: number) {
  return useQuery({
    queryKey: adminRaffleKeys.detail(raffleId),
    queryFn: () => adminRafflesApi.getDetail(raffleId),
    enabled: raffleId > 0,
    staleTime: 60 * 1000, // 1 minuto
  });
}

// ===========================
// Mutations
// ===========================

export function useForceStatusChange() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      raffleId,
      data,
    }: {
      raffleId: number;
      data: ForceStatusChangeRequest;
    }) => adminRafflesApi.forceStatusChange(raffleId, data),
    onSuccess: (_, variables) => {
      toast.success("Estado de rifa actualizado exitosamente");
      queryClient.invalidateQueries({ queryKey: adminRaffleKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: adminRaffleKeys.detail(variables.raffleId),
      });
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message ||
          "Error al actualizar estado de la rifa"
      );
    },
  });
}

export function useAddAdminNotes() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      raffleId,
      data,
    }: {
      raffleId: number;
      data: AddAdminNotesRequest;
    }) => adminRafflesApi.addAdminNotes(raffleId, data),
    onSuccess: (_, variables) => {
      toast.success("Notas agregadas exitosamente");
      queryClient.invalidateQueries({
        queryKey: adminRaffleKeys.detail(variables.raffleId),
      });
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message || "Error al agregar notas"
      );
    },
  });
}

export function useManualDraw() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      raffleId,
      data,
    }: {
      raffleId: number;
      data: ManualDrawRequest;
    }) => adminRafflesApi.manualDraw(raffleId, data),
    onSuccess: (_, variables) => {
      toast.success("Sorteo manual realizado exitosamente");
      queryClient.invalidateQueries({ queryKey: adminRaffleKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: adminRaffleKeys.detail(variables.raffleId),
      });
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message || "Error al realizar sorteo manual"
      );
    },
  });
}

export function useCancelWithRefund() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      raffleId,
      data,
    }: {
      raffleId: number;
      data: CancelWithRefundRequest;
    }) => adminRafflesApi.cancelWithRefund(raffleId, data),
    onSuccess: (_, variables) => {
      toast.success("Rifa cancelada y reembolsos iniciados");
      queryClient.invalidateQueries({ queryKey: adminRaffleKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: adminRaffleKeys.detail(variables.raffleId),
      });
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message ||
          "Error al cancelar rifa con reembolso"
      );
    },
  });
}
