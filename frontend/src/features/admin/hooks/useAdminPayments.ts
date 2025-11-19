import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { adminPaymentsApi } from "../api/adminPaymentsApi";
import type {
  PaymentFilters,
  ProcessRefundRequest,
  ManageDisputeRequest,
  PaginationParams,
} from "../types";

// Query keys
export const paymentKeys = {
  all: ["admin", "payments"] as const,
  lists: () => [...paymentKeys.all, "list"] as const,
  list: (filters: PaymentFilters, pagination: PaginationParams) =>
    [...paymentKeys.lists(), filters, pagination] as const,
  details: () => [...paymentKeys.all, "detail"] as const,
  detail: (id: string) => [...paymentKeys.details(), id] as const,
};

// List payments
export function useAdminPayments(
  filters: PaymentFilters = {},
  pagination: PaginationParams = { page: 1, limit: 20 }
) {
  return useQuery({
    queryKey: paymentKeys.list(filters, pagination),
    queryFn: () => adminPaymentsApi.getAll(filters, pagination),
    staleTime: 30000, // 30 seconds
  });
}

// Get payment detail
export function useAdminPaymentDetail(paymentId: string) {
  return useQuery({
    queryKey: paymentKeys.detail(paymentId),
    queryFn: () => adminPaymentsApi.getDetail(paymentId),
    enabled: !!paymentId,
    staleTime: 60000, // 1 minute
  });
}

// Process refund
export function useProcessRefund() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      paymentId,
      data,
    }: {
      paymentId: string;
      data: ProcessRefundRequest;
    }) => adminPaymentsApi.processRefund(paymentId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: paymentKeys.lists() });
      queryClient.invalidateQueries({ queryKey: paymentKeys.details() });
      toast.success("Reembolso procesado exitosamente");
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message || "Error al procesar reembolso"
      );
    },
  });
}

// Manage dispute
export function useManageDispute() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      paymentId,
      data,
    }: {
      paymentId: string;
      data: ManageDisputeRequest;
    }) => adminPaymentsApi.manageDispute(paymentId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: paymentKeys.lists() });
      queryClient.invalidateQueries({ queryKey: paymentKeys.details() });
      toast.success("Disputa gestionada exitosamente");
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message || "Error al gestionar disputa"
      );
    },
  });
}
