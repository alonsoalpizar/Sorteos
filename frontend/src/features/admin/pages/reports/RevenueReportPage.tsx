import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  ArrowLeft,
  TrendingUp,
  Calendar,
  Users,
  Package,
  DollarSign,
  BarChart3,
} from "lucide-react";
import { useRevenueReport } from "../../hooks/useAdminReports";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import { formatCurrency } from "@/lib/currency";
import { format, subDays, parseISO } from "date-fns";
import { es } from "date-fns/locale";
import type { RevenueReportInput } from "../../types";

export function RevenueReportPage() {
  const navigate = useNavigate();

  // Filtros por defecto: últimos 30 días
  const [filters, setFilters] = useState<RevenueReportInput>({
    date_from: format(subDays(new Date(), 30), "yyyy-MM-dd"),
    date_to: format(new Date(), "yyyy-MM-dd"),
    group_by: "day",
  });

  const { data, isLoading, error } = useRevenueReport(filters);

  const handleFilterChange = (key: keyof RevenueReportInput, value: any) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  const handleClearOrganizerFilter = () => {
    setFilters((prev) => {
      const newFilters = { ...prev };
      delete newFilters.organizer_id;
      return newFilters;
    });
  };

  const handleClearCategoryFilter = () => {
    setFilters((prev) => {
      const newFilters = { ...prev };
      delete newFilters.category_id;
      return newFilters;
    });
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="outline" onClick={() => navigate("/admin/reports")}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Volver
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Reporte de Ingresos</h1>
            <p className="text-slate-600 mt-2">
              Análisis de ingresos y métricas de rendimiento
            </p>
          </div>
        </div>
      </div>

      {/* Filtros */}
      <Card className="p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Calendar className="w-5 h-5" />
          Filtros
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {/* Fecha desde */}
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Fecha Desde
            </label>
            <input
              type="date"
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              value={filters.date_from}
              onChange={(e) => handleFilterChange("date_from", e.target.value)}
            />
          </div>

          {/* Fecha hasta */}
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Fecha Hasta
            </label>
            <input
              type="date"
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              value={filters.date_to}
              onChange={(e) => handleFilterChange("date_to", e.target.value)}
            />
          </div>

          {/* Agrupar por */}
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Agrupar Por
            </label>
            <select
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              value={filters.group_by}
              onChange={(e) =>
                handleFilterChange("group_by", e.target.value as "day" | "week" | "month")
              }
            >
              <option value="day">Día</option>
              <option value="week">Semana</option>
              <option value="month">Mes</option>
            </select>
          </div>

          {/* Organizador (opcional) */}
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Organizador (Opcional)
            </label>
            <div className="flex gap-2">
              <input
                type="number"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="ID del organizador"
                value={filters.organizer_id || ""}
                onChange={(e) =>
                  handleFilterChange(
                    "organizer_id",
                    e.target.value ? parseInt(e.target.value) : undefined
                  )
                }
              />
              {filters.organizer_id && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleClearOrganizerFilter}
                >
                  Limpiar
                </Button>
              )}
            </div>
          </div>
        </div>

        {/* Filtro de categoría (segunda fila) */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mt-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Categoría (Opcional)
            </label>
            <div className="flex gap-2">
              <input
                type="number"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="ID de categoría"
                value={filters.category_id || ""}
                onChange={(e) =>
                  handleFilterChange(
                    "category_id",
                    e.target.value ? parseInt(e.target.value) : undefined
                  )
                }
              />
              {filters.category_id && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleClearCategoryFilter}
                >
                  Limpiar
                </Button>
              )}
            </div>
          </div>
        </div>
      </Card>

      {/* Loading */}
      {isLoading && (
        <div className="flex items-center justify-center py-12">
          <LoadingSpinner />
        </div>
      )}

      {/* Error */}
      {error && !isLoading && (
        <EmptyState
          icon={<TrendingUp className="w-12 h-12 text-red-500" />}
          title="Error al cargar reporte"
          description={(error as Error)?.message || "No se pudo cargar el reporte de ingresos"}
        />
      )}

      {/* Datos */}
      {data && !isLoading && (
        <>
          {/* KPIs Resumen */}
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
            <Card className="p-4">
              <p className="text-sm font-medium text-slate-600">Ingresos Brutos</p>
              <p className="text-2xl font-bold text-green-600 mt-2">
                {formatCurrency(data.TotalGrossRevenue)}
              </p>
            </Card>

            <Card className="p-4">
              <p className="text-sm font-medium text-slate-600">Comisiones Plataforma</p>
              <p className="text-2xl font-bold text-blue-600 mt-2">
                {formatCurrency(data.TotalPlatformFees)}
              </p>
            </Card>

            <Card className="p-4">
              <p className="text-sm font-medium text-slate-600">Ingresos Netos</p>
              <p className="text-2xl font-bold text-slate-900 mt-2">
                {formatCurrency(data.TotalNetRevenue)}
              </p>
            </Card>

            <Card className="p-4">
              <p className="text-sm font-medium text-slate-600">Total Pagos</p>
              <p className="text-2xl font-bold text-purple-600 mt-2">
                {data.TotalPayments.toLocaleString()}
              </p>
              <p className="text-xs text-slate-500 mt-1">
                Promedio: {formatCurrency(data.TotalGrossRevenue / (data.TotalPayments || 1))}
              </p>
            </Card>

            <Card className="p-4">
              <p className="text-sm font-medium text-slate-600">Total Rifas</p>
              <p className="text-2xl font-bold text-orange-600 mt-2">
                {data.TotalRaffles.toLocaleString()}
              </p>
              <p className="text-xs text-slate-500 mt-1">
                Promedio: {formatCurrency(data.AverageRevenuePerRaffle)}
              </p>
            </Card>
          </div>

          {/* Chart y Estadísticas */}
          <Card className="p-6">
            <h2 className="text-xl font-semibold text-slate-900 mb-6 flex items-center gap-2">
              <BarChart3 className="w-5 h-5" />
              Tendencia de Ingresos
            </h2>

            {data.DataPoints.length === 0 ? (
              <EmptyState
                icon={<TrendingUp className="w-12 h-12 text-slate-400" />}
                title="No hay datos"
                description="No se encontraron datos para el período seleccionado"
              />
            ) : (
              <div className="space-y-6">
                {/* Simple Bar Chart con divs (sin librería) */}
                <div className="space-y-2">
                  {data.DataPoints.map((point, index) => {
                    const maxRevenue = Math.max(
                      ...data.DataPoints.map((p) => p.gross_revenue)
                    );
                    const widthPercent = (point.gross_revenue / maxRevenue) * 100;

                    return (
                      <div key={index} className="space-y-1">
                        <div className="flex items-center justify-between text-sm">
                          <span className="font-medium text-slate-700">
                            {point.date ? format(parseISO(point.date), "PPP", { locale: es }) : "Fecha inválida"}
                          </span>
                          <div className="flex items-center gap-4">
                            <span className="text-green-600 font-semibold">
                              {formatCurrency(point.gross_revenue)}
                            </span>
                            <span className="text-blue-600 text-xs">
                              Fee: {formatCurrency(point.platform_fees)}
                            </span>
                            <span className="text-slate-600 text-xs">
                              {point.payment_count} pagos
                            </span>
                          </div>
                        </div>
                        <div className="h-8 bg-slate-100 rounded-lg overflow-hidden">
                          <div
                            className="h-full bg-gradient-to-r from-green-500 to-green-600 rounded-lg transition-all duration-300"
                            style={{ width: `${widthPercent}%` }}
                          />
                        </div>
                      </div>
                    );
                  })}
                </div>

                {/* Tabla de datos */}
                <div className="mt-8">
                  <h3 className="text-lg font-semibold text-slate-900 mb-4">
                    Datos Detallados
                  </h3>
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead className="bg-slate-50 border-b-2 border-slate-200">
                        <tr>
                          <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                            Fecha
                          </th>
                          <th className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase">
                            Ingresos Brutos
                          </th>
                          <th className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase">
                            Comisión
                          </th>
                          <th className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase">
                            Ingresos Netos
                          </th>
                          <th className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase">
                            Pagos
                          </th>
                          <th className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase">
                            Rifas
                          </th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-200">
                        {data.DataPoints.map((point, index) => (
                          <tr key={index} className="hover:bg-slate-50">
                            <td className="px-4 py-3 text-sm text-slate-900">
                              {point.date ? format(parseISO(point.date), "PPP", { locale: es }) : "Fecha inválida"}
                            </td>
                            <td className="px-4 py-3 text-sm text-right font-semibold text-green-600">
                              {formatCurrency(point.gross_revenue)}
                            </td>
                            <td className="px-4 py-3 text-sm text-right text-blue-600">
                              {formatCurrency(point.platform_fees)}
                            </td>
                            <td className="px-4 py-3 text-sm text-right font-semibold text-slate-900">
                              {formatCurrency(point.net_revenue)}
                            </td>
                            <td className="px-4 py-3 text-sm text-right text-slate-600">
                              {point.payment_count}
                            </td>
                            <td className="px-4 py-3 text-sm text-right text-slate-600">
                              {point.raffle_count}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            )}
          </Card>

          {/* Métricas adicionales */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card className="p-4">
              <div className="flex items-center gap-3 mb-2">
                <DollarSign className="w-5 h-5 text-green-600" />
                <p className="text-sm font-medium text-slate-600">
                  Ingreso Promedio por Día
                </p>
              </div>
              <p className="text-2xl font-bold text-slate-900">
                {formatCurrency(data.AverageRevenuePerDay)}
              </p>
            </Card>

            <Card className="p-4">
              <div className="flex items-center gap-3 mb-2">
                <Package className="w-5 h-5 text-orange-600" />
                <p className="text-sm font-medium text-slate-600">
                  Ingreso Promedio por Rifa
                </p>
              </div>
              <p className="text-2xl font-bold text-slate-900">
                {formatCurrency(data.AverageRevenuePerRaffle)}
              </p>
            </Card>

            <Card className="p-4">
              <div className="flex items-center gap-3 mb-2">
                <Users className="w-5 h-5 text-blue-600" />
                <p className="text-sm font-medium text-slate-600">
                  Tasa de Comisión Promedio
                </p>
              </div>
              <p className="text-2xl font-bold text-slate-900">
                {((data.TotalPlatformFees / (data.TotalGrossRevenue || 1)) * 100).toFixed(
                  2
                )}
                %
              </p>
            </Card>
          </div>
        </>
      )}
    </div>
  );
}
