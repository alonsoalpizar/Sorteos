import { useQuery } from "@tanstack/react-query";
import { adminAuditApi } from "../api/adminAuditApi";
import type { ListAuditLogsInput } from "../types";

// Query Keys
export const auditKeys = {
  all: ["admin", "audit"] as const,
  logs: (input: ListAuditLogsInput = {}) => [...auditKeys.all, "logs", input] as const,
};

// List Audit Logs
export function useAuditLogs(input: ListAuditLogsInput = {}, enabled = true) {
  return useQuery({
    queryKey: auditKeys.logs(input),
    queryFn: () => adminAuditApi.listAuditLogs(input),
    enabled,
    staleTime: 60000, // 1 minuto (datos de auditoría cambian frecuentemente)
  });
}
