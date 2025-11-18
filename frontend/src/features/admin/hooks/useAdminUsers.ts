import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { adminUsersApi } from "../api/adminApi";
import type {
  UserFilters,
  PaginationParams,
  UpdateUserStatusRequest,
  UpdateUserKYCRequest,
} from "../types";

// Query keys
export const adminUsersKeys = {
  all: ["admin", "users"] as const,
  lists: () => [...adminUsersKeys.all, "list"] as const,
  list: (filters?: UserFilters, pagination?: PaginationParams) =>
    [...adminUsersKeys.lists(), { filters, pagination }] as const,
  details: () => [...adminUsersKeys.all, "detail"] as const,
  detail: (id: number) => [...adminUsersKeys.details(), id] as const,
};

// ===========================
// Queries
// ===========================

/**
 * Hook to fetch paginated list of users with filters
 */
export function useAdminUsers(
  filters?: UserFilters,
  pagination?: PaginationParams
) {
  return useQuery({
    queryKey: adminUsersKeys.list(filters, pagination),
    queryFn: () => adminUsersApi.list(filters, pagination),
    staleTime: 30 * 1000, // 30 seconds
  });
}

/**
 * Hook to fetch user detail
 */
export function useAdminUserDetail(userId: number) {
  return useQuery({
    queryKey: adminUsersKeys.detail(userId),
    queryFn: () => adminUsersApi.getDetail(userId),
    enabled: !!userId,
    staleTime: 60 * 1000, // 1 minute
  });
}

// ===========================
// Mutations
// ===========================

/**
 * Hook to update user status (suspend, activate, ban)
 */
export function useUpdateUserStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      userId,
      data,
    }: {
      userId: number;
      data: UpdateUserStatusRequest;
    }) => adminUsersApi.updateStatus(userId, data),

    onSuccess: (_, variables) => {
      // Invalidate queries
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: adminUsersKeys.detail(variables.userId),
      });

      toast.success("Estado del usuario actualizado correctamente");
    },

    onError: (error: Error) => {
      toast.error("Error al actualizar estado", {
        description: error.message,
      });
    },
  });
}

/**
 * Hook to update user KYC level
 */
export function useUpdateUserKYC() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      userId,
      data,
    }: {
      userId: number;
      data: UpdateUserKYCRequest;
    }) => adminUsersApi.updateKYC(userId, data),

    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.lists() });
      queryClient.invalidateQueries({
        queryKey: adminUsersKeys.detail(variables.userId),
      });

      toast.success("Nivel KYC actualizado correctamente");
    },

    onError: (error: Error) => {
      toast.error("Error al actualizar KYC", {
        description: error.message,
      });
    },
  });
}

/**
 * Hook to delete user (soft delete)
 */
export function useDeleteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: number) => adminUsersApi.deleteUser(userId),

    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminUsersKeys.lists() });

      toast.success("Usuario eliminado correctamente");
    },

    onError: (error: Error) => {
      toast.error("Error al eliminar usuario", {
        description: error.message,
      });
    },
  });
}
