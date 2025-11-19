import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { adminSystemApi } from "../api/adminSystemApi";
import type { GetSystemSettingsInput, UpdateSystemSettingsInput } from "../types";
import { toast } from "sonner";

// Query Keys
export const systemKeys = {
  all: ["admin", "system"] as const,
  settings: (input: GetSystemSettingsInput = {}) => [...systemKeys.all, "settings", input] as const,
};

// Get System Settings
export function useSystemSettings(input: GetSystemSettingsInput = {}, enabled = true) {
  return useQuery({
    queryKey: systemKeys.settings(input),
    queryFn: () => adminSystemApi.getSystemSettings(input),
    enabled,
    staleTime: 300000, // 5 minutos
  });
}

// Update System Setting
export function useUpdateSystemSetting() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UpdateSystemSettingsInput) => adminSystemApi.updateSystemSetting(input),
    onSuccess: (data) => {
      toast.success("Configuración actualizada", {
        description: `${data.key} actualizado exitosamente`,
      });
      // Invalidar todas las queries de settings para que se recarguen
      queryClient.invalidateQueries({ queryKey: systemKeys.all });
    },
    onError: (error: any) => {
      toast.error("Error al actualizar configuración", {
        description: error.response?.data?.error?.message || error.message,
      });
    },
  });
}
