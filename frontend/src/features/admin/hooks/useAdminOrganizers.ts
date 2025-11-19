import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { adminOrganizersApi } from "../api/adminApi";
import type {
  OrganizerFilters,
  PaginationParams,
  VerifyOrganizerRequest,
  UpdateCommissionRequest,
} from "../types";

// Query keys
export const adminOrganizersKeys = {
  all: ["admin", "organizers"] as const,
  lists: () => [...adminOrganizersKeys.all, "list"] as const,
  list: (filters?: OrganizerFilters, pagination?: PaginationParams) =>
    [...adminOrganizersKeys.lists(), { filters, pagination }] as const,
  details: () => [...adminOrganizersKeys.all, "detail"] as const,
  detail: (id: number) => [...adminOrganizersKeys.details(), id] as const,
};

// ===========================
// Queries
// ===========================

/**
 * Hook to fetch paginated list of organizers with filters
 */
export function useAdminOrganizers(
  filters?: OrganizerFilters,
  pagination?: PaginationParams
) {
  return useQuery({
    queryKey: adminOrganizersKeys.list(filters, pagination),
    queryFn: () => adminOrganizersApi.list(filters, pagination),
    staleTime: 30 * 1000, // 30 seconds
  });
}

/**
 * Hook to fetch organizer detail
 */
export function useAdminOrganizerDetail(organizerId: number) {
  return useQuery({
    queryKey: adminOrganizersKeys.detail(organizerId),
    queryFn: () => adminOrganizersApi.getDetail(organizerId),
    enabled: !!organizerId,
    staleTime: 60 * 1000, // 1 minute
  });
}

// ===========================
// Mutations
// ===========================

/**
 * Hook to verify/unverify organizer
 */
export function useVerifyOrganizer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      organizerId,
      data,
    }: {
      organizerId: number;
      data: VerifyOrganizerRequest;
    }) => adminOrganizersApi.verify(organizerId, data),

    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: adminOrganizersKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: adminOrganizersKeys.detail(variables.organizerId),
      });

      toast.success("Estado de verificación actualizado correctamente");
    },

    onError: (error: Error) => {
      toast.error("Error al actualizar verificación", {
        description: error.message,
      });
    },
  });
}

/**
 * Hook to update organizer commission
 */
export function useUpdateOrganizerCommission() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      organizerId,
      data,
    }: {
      organizerId: number;
      data: UpdateCommissionRequest;
    }) => adminOrganizersApi.updateCommission(organizerId, data),

    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: adminOrganizersKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: adminOrganizersKeys.detail(variables.organizerId),
      });

      toast.success("Comisión actualizada correctamente");
    },

    onError: (error: Error) => {
      toast.error("Error al actualizar comisión", {
        description: error.message,
      });
    },
  });
}
