import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { adminCategoriesApi, type CategoryFilters } from "../api/adminCategoriesApi";
import type {
  CreateCategoryRequest,
  UpdateCategoryRequest,
  ReorderCategoriesRequest,
  PaginationParams,
} from "../types";

// Query keys
export const categoryKeys = {
  all: ["admin", "categories"] as const,
  lists: () => [...categoryKeys.all, "list"] as const,
  list: (filters: CategoryFilters, pagination: PaginationParams) =>
    [...categoryKeys.lists(), filters, pagination] as const,
};

// List categories
export function useAdminCategories(
  filters: CategoryFilters = {},
  pagination: PaginationParams = { page: 1, limit: 50 }
) {
  return useQuery({
    queryKey: categoryKeys.list(filters, pagination),
    queryFn: () => adminCategoriesApi.getAll(filters, pagination),
    staleTime: 30000, // 30 seconds
  });
}

// Create category
export function useCreateCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateCategoryRequest) =>
      adminCategoriesApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.lists() });
      toast.success("Categoría creada exitosamente");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Error al crear categoría");
    },
  });
}

// Update category
export function useUpdateCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      categoryId,
      data,
    }: {
      categoryId: number;
      data: UpdateCategoryRequest;
    }) => adminCategoriesApi.update(categoryId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.lists() });
      toast.success("Categoría actualizada exitosamente");
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message || "Error al actualizar categoría"
      );
    },
  });
}

// Delete category
export function useDeleteCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (categoryId: number) => adminCategoriesApi.delete(categoryId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.lists() });
      toast.success("Categoría eliminada exitosamente");
    },
    onError: (error: any) => {
      const message = error.response?.data?.message;
      if (message?.includes("in use") || message?.includes("being used")) {
        toast.error("No se puede eliminar: la categoría está siendo usada por rifas activas");
      } else {
        toast.error(message || "Error al eliminar categoría");
      }
    },
  });
}

// Reorder categories
export function useReorderCategories() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ReorderCategoriesRequest) =>
      adminCategoriesApi.reorder(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.lists() });
      toast.success("Orden actualizado exitosamente");
    },
    onError: (error: any) => {
      toast.error(
        error.response?.data?.message || "Error al reordenar categorías"
      );
    },
  });
}
