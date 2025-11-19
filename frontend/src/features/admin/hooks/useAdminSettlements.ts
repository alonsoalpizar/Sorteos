import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { adminSettlementsApi } from "../api/adminSettlementsApi";
import type {
  SettlementFilters,
  CreateSettlementRequest,
  ApproveSettlementRequest,
  RejectSettlementRequest,
  MarkSettlementPaidRequest,
  AutoCreateSettlementsRequest,
  PaginationParams,
} from "../types";

// Query keys
export const settlementKeys = {
  all: ["admin", "settlements"] as const,
  lists: () => [...settlementKeys.all, "list"] as const,
  list: (filters: SettlementFilters, pagination: PaginationParams) =>
    [...settlementKeys.lists(), filters, pagination] as const,
  details: () => [...settlementKeys.all, "detail"] as const,
  detail: (id: number) => [...settlementKeys.details(), id] as const,
};

// List settlements
export function useAdminSettlements(
  filters: SettlementFilters = {},
  pagination: PaginationParams = { page: 1, limit: 20 }
) {
  return useQuery({
    queryKey: settlementKeys.list(filters, pagination),
    queryFn: () => adminSettlementsApi.getAll(filters, pagination),
    staleTime: 30000, // 30 seconds
  });
}

// Get settlement details
export function useAdminSettlementDetail(settlementId: number) {
  return useQuery({
    queryKey: settlementKeys.detail(settlementId),
    queryFn: () => adminSettlementsApi.getDetail(settlementId),
    enabled: !!settlementId,
    staleTime: 60000, // 1 minute
  });
}

// Create settlement
export function useCreateSettlement() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSettlementRequest) => adminSettlementsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settlementKeys.lists() });
      toast.success("Liquidación creada exitosamente");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Error al crear liquidación");
    },
  });
}

// Approve settlement
export function useApproveSettlement() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ settlementId, data }: { settlementId: number; data: ApproveSettlementRequest }) =>
      adminSettlementsApi.approve(settlementId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: settlementKeys.lists() });
      queryClient.invalidateQueries({ queryKey: settlementKeys.detail(variables.settlementId) });
      toast.success("Liquidación aprobada exitosamente");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Error al aprobar liquidación");
    },
  });
}

// Reject settlement
export function useRejectSettlement() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ settlementId, data }: { settlementId: number; data: RejectSettlementRequest }) =>
      adminSettlementsApi.reject(settlementId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: settlementKeys.lists() });
      queryClient.invalidateQueries({ queryKey: settlementKeys.detail(variables.settlementId) });
      toast.success("Liquidación rechazada");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Error al rechazar liquidación");
    },
  });
}

// Mark settlement as paid
export function useMarkSettlementPaid() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ settlementId, data }: { settlementId: number; data: MarkSettlementPaidRequest }) =>
      adminSettlementsApi.markPaid(settlementId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: settlementKeys.lists() });
      queryClient.invalidateQueries({ queryKey: settlementKeys.detail(variables.settlementId) });
      toast.success("Liquidación marcada como pagada");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Error al marcar como pagada");
    },
  });
}

// Auto-create settlements
export function useAutoCreateSettlements() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AutoCreateSettlementsRequest) => adminSettlementsApi.autoCreate(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: settlementKeys.lists() });
      toast.success("Liquidaciones creadas automáticamente");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Error al auto-crear liquidaciones");
    },
  });
}
