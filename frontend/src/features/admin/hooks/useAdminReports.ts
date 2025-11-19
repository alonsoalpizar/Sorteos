import { useQuery, useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { adminReportsApi } from "../api/adminReportsApi";
import type {
  RevenueReportInput,
  RaffleLiquidationsReportInput,
  ExportDataInput,
} from "../types";

// Query keys
export const reportsKeys = {
  all: ["admin", "reports"] as const,
  dashboard: () => [...reportsKeys.all, "dashboard"] as const,
  revenue: (input: RevenueReportInput) => [...reportsKeys.all, "revenue", input] as const,
  liquidations: (input: RaffleLiquidationsReportInput) =>
    [...reportsKeys.all, "liquidations", input] as const,
};

// Dashboard KPIs
export function useAdminDashboard() {
  return useQuery({
    queryKey: reportsKeys.dashboard(),
    queryFn: () => adminReportsApi.getDashboard(),
    staleTime: 60000, // 1 minute
  });
}

// Revenue Report
export function useRevenueReport(input: RevenueReportInput, enabled = true) {
  return useQuery({
    queryKey: reportsKeys.revenue(input),
    queryFn: () => adminReportsApi.getRevenueReport(input),
    enabled: enabled && !!input.date_from && !!input.date_to,
    staleTime: 300000, // 5 minutes
  });
}

// Liquidations Report
export function useLiquidationsReport(input: RaffleLiquidationsReportInput, enabled = true) {
  return useQuery({
    queryKey: reportsKeys.liquidations(input),
    queryFn: () => adminReportsApi.getLiquidationsReport(input),
    enabled: enabled && !!input.date_from && !!input.date_to,
    staleTime: 300000, // 5 minutes
  });
}

// Export Data
export function useExportData() {
  return useMutation({
    mutationFn: (input: ExportDataInput) => adminReportsApi.exportData(input),
    onSuccess: (data) => {
      toast.success("Exportación completada");
      // Abrir el archivo en una nueva pestaña
      window.open(data.file_url, "_blank");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Error al exportar datos");
    },
  });
}
