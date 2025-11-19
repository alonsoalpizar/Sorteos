import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { adminNotificationsApi } from "../api/adminNotificationsApi";
import type {
  SendEmailInput,
  SendBulkEmailInput,
  ViewNotificationHistoryInput,
} from "../types";
import { toast } from "sonner";

// Query Keys
export const notificationsKeys = {
  all: ["admin", "notifications"] as const,
  history: (input: ViewNotificationHistoryInput) =>
    [...notificationsKeys.all, "history", input] as const,
};

// View Notification History
export function useNotificationHistory(input: ViewNotificationHistoryInput, enabled = true) {
  return useQuery({
    queryKey: notificationsKeys.history(input),
    queryFn: () => adminNotificationsApi.viewHistory(input),
    enabled,
    staleTime: 30000, // 30 segundos
  });
}

// Send Email
export function useSendEmail() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: SendEmailInput) => adminNotificationsApi.sendEmail(input),
    onSuccess: (data) => {
      toast.success(`Email ${data.status} exitosamente`, {
        description: `Enviado a ${data.recipients} destinatario(s)`,
      });
      // Invalidar historial para que se recargue
      queryClient.invalidateQueries({ queryKey: notificationsKeys.all });
    },
    onError: (error: any) => {
      toast.error("Error al enviar email", {
        description: error.response?.data?.error?.message || error.message,
      });
    },
  });
}

// Send Bulk Email
export function useSendBulkEmail() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: SendBulkEmailInput) => adminNotificationsApi.sendBulkEmail(input),
    onSuccess: (data) => {
      toast.success("Email masivo enviado", {
        description: `${data.sent_count} de ${data.total_recipients} emails enviados`,
      });
      queryClient.invalidateQueries({ queryKey: notificationsKeys.all });
    },
    onError: (error: any) => {
      toast.error("Error al enviar email masivo", {
        description: error.response?.data?.error?.message || error.message,
      });
    },
  });
}
