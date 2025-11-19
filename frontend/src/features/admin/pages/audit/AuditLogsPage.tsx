import { useState } from "react";
import { Shield, Search, Filter, RefreshCw, Calendar, User, AlertTriangle, Info, XCircle, AlertCircle } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import { useAuditLogs } from "../../hooks/useAdminAudit";
import type { ListAuditLogsInput } from "../../types";
import { format, parseISO } from "date-fns";
import { es } from "date-fns/locale";

export function AuditLogsPage() {
  const [filters, setFilters] = useState<ListAuditLogsInput>({
    page: 1,
    page_size: 20,
    order_by: "-created_at",
  });

  const [searchTerm, setSearchTerm] = useState("");
  const [showFilters, setShowFilters] = useState(false);

  const { data, isLoading, refetch } = useAuditLogs(filters);

  const handleSearch = () => {
    setFilters({
      ...filters,
      search: searchTerm || undefined,
      page: 1,
    });
  };

  const handleFilterChange = (key: keyof ListAuditLogsInput, value: any) => {
    setFilters({
      ...filters,
      [key]: value || undefined,
      page: 1,
    });
  };

  const handlePageChange = (newPage: number) => {
    setFilters({
      ...filters,
      page: newPage,
    });
  };

  const getSeverityBadge = (severity: string) => {
    switch (severity) {
      case "info":
        return (
          <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-700">
            <Info className="w-3 h-3" />
            Info
          </span>
        );
      case "warning":
        return (
          <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-700">
            <AlertTriangle className="w-3 h-3" />
            Advertencia
          </span>
        );
      case "error":
        return (
          <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-700">
            <XCircle className="w-3 h-3" />
            Error
          </span>
        );
      case "critical":
        return (
          <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-red-200 text-red-900">
            <AlertCircle className="w-3 h-3" />
            Crítico
          </span>
        );
      default:
        return <span className="text-xs text-slate-500">{severity}</span>;
    }
  };

  const getEntityTypeBadge = (entityType: string) => {
    const colors: Record<string, string> = {
      user: "bg-purple-100 text-purple-700",
      organizer: "bg-green-100 text-green-700",
      raffle: "bg-blue-100 text-blue-700",
      payment: "bg-yellow-100 text-yellow-700",
      settlement: "bg-indigo-100 text-indigo-700",
      category: "bg-slate-100 text-slate-700",
      system: "bg-red-100 text-red-700",
    };

    const colorClass = colors[entityType.toLowerCase()] || "bg-slate-100 text-slate-700";

    return (
      <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${colorClass}`}>
        {entityType}
      </span>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">Auditoría</h1>
          <p className="text-slate-600 mt-2">Registro de acciones administrativas del sistema</p>
        </div>
        <Button variant="outline" onClick={() => refetch()} disabled={isLoading}>
          <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? "animate-spin" : ""}`} />
          Recargar
        </Button>
      </div>

      {/* Statistics */}
      {data && (
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Total Registros</p>
            <p className="text-2xl font-bold text-slate-900 mt-2">{data.Total}</p>
          </Card>
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Info</p>
            <p className="text-2xl font-bold text-blue-600 mt-2">{data.InfoCount}</p>
          </Card>
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Advertencias</p>
            <p className="text-2xl font-bold text-yellow-600 mt-2">{data.WarningCount}</p>
          </Card>
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Errores</p>
            <p className="text-2xl font-bold text-red-600 mt-2">{data.ErrorCount}</p>
          </Card>
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Críticos</p>
            <p className="text-2xl font-bold text-red-900 mt-2">{data.CriticalCount}</p>
          </Card>
        </div>
      )}

      {/* Search and Filters */}
      <Card className="p-4">
        <div className="space-y-4">
          {/* Search Bar */}
          <div className="flex items-center gap-3">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-slate-400" />
              <input
                type="text"
                className="w-full pl-10 pr-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Buscar en descripción, admin, acción..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                onKeyPress={(e) => e.key === "Enter" && handleSearch()}
              />
            </div>
            <Button onClick={handleSearch}>
              <Search className="w-4 h-4 mr-2" />
              Buscar
            </Button>
            <Button
              variant="outline"
              onClick={() => setShowFilters(!showFilters)}
            >
              <Filter className="w-4 h-4 mr-2" />
              Filtros
            </Button>
          </div>

          {/* Advanced Filters */}
          {showFilters && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4 border-t border-slate-200">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">Severidad</label>
                <select
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  value={filters.severity || ""}
                  onChange={(e) => handleFilterChange("severity", e.target.value)}
                >
                  <option value="">Todas</option>
                  <option value="info">Info</option>
                  <option value="warning">Advertencia</option>
                  <option value="error">Error</option>
                  <option value="critical">Crítico</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">Tipo de Entidad</label>
                <select
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  value={filters.entity_type || ""}
                  onChange={(e) => handleFilterChange("entity_type", e.target.value)}
                >
                  <option value="">Todas</option>
                  <option value="user">Usuario</option>
                  <option value="organizer">Organizador</option>
                  <option value="raffle">Rifa</option>
                  <option value="payment">Pago</option>
                  <option value="settlement">Liquidación</option>
                  <option value="category">Categoría</option>
                  <option value="system">Sistema</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">Acción</label>
                <input
                  type="text"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="ej: CREATE, UPDATE, DELETE"
                  value={filters.action || ""}
                  onChange={(e) => handleFilterChange("action", e.target.value)}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">Fecha Desde</label>
                <input
                  type="datetime-local"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  value={filters.date_from || ""}
                  onChange={(e) => handleFilterChange("date_from", e.target.value)}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">Fecha Hasta</label>
                <input
                  type="datetime-local"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  value={filters.date_to || ""}
                  onChange={(e) => handleFilterChange("date_to", e.target.value)}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-2">ID de Admin</label>
                <input
                  type="number"
                  className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="ej: 1"
                  value={filters.admin_id || ""}
                  onChange={(e) => handleFilterChange("admin_id", e.target.value ? parseInt(e.target.value) : undefined)}
                />
              </div>

              <div className="md:col-span-3 flex items-center gap-3">
                <Button
                  variant="outline"
                  onClick={() => {
                    setFilters({
                      page: 1,
                      page_size: 20,
                      order_by: "-created_at",
                    });
                    setSearchTerm("");
                  }}
                >
                  Limpiar Filtros
                </Button>
              </div>
            </div>
          )}
        </div>
      </Card>

      {/* Audit Logs Table */}
      <Card className="p-6">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <LoadingSpinner />
          </div>
        ) : data && data.Logs && data.Logs.length > 0 ? (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-slate-50 border-b-2 border-slate-200">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                      ID
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                      Fecha
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                      Admin
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                      Acción
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                      Entidad
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                      Descripción
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase">
                      Severidad
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-200">
                  {data.Logs.map((log) => (
                    <tr key={log.id} className="hover:bg-slate-50">
                      <td className="px-4 py-3 text-sm text-slate-900">#{log.id}</td>
                      <td className="px-4 py-3 text-sm text-slate-600">
                        <div className="flex items-center gap-1">
                          <Calendar className="w-3 h-3" />
                          {format(parseISO(log.created_at), "PPp", { locale: es })}
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-1">
                          <User className="w-3 h-3 text-slate-400" />
                          <div>
                            <p className="text-sm font-medium text-slate-900">{log.admin_name}</p>
                            <p className="text-xs text-slate-500">{log.admin_email}</p>
                          </div>
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <code className="px-2 py-1 bg-slate-100 rounded text-xs font-mono text-slate-900">
                          {log.action}
                        </code>
                      </td>
                      <td className="px-4 py-3">
                        <div className="space-y-1">
                          {getEntityTypeBadge(log.entity_type)}
                          {log.entity_id && (
                            <p className="text-xs text-slate-500">ID: {log.entity_id}</p>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <p className="text-sm text-slate-900 max-w-md">{log.description}</p>
                        {log.metadata && (
                          <details className="mt-1">
                            <summary className="text-xs text-blue-600 cursor-pointer hover:text-blue-700">
                              Ver metadatos
                            </summary>
                            <pre className="mt-2 p-2 bg-slate-50 rounded text-xs font-mono overflow-x-auto">
                              {log.metadata}
                            </pre>
                          </details>
                        )}
                      </td>
                      <td className="px-4 py-3">{getSeverityBadge(log.severity)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {data.TotalPages > 1 && (
              <div className="flex items-center justify-between mt-6 pt-4 border-t border-slate-200">
                <p className="text-sm text-slate-600">
                  Página {data.Page} de {data.TotalPages} • {data.Total} registros totales
                </p>
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(data.Page - 1)}
                    disabled={data.Page === 1}
                  >
                    Anterior
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(data.Page + 1)}
                    disabled={data.Page >= data.TotalPages}
                  >
                    Siguiente
                  </Button>
                </div>
              </div>
            )}
          </>
        ) : (
          <EmptyState
            icon={<Shield className="w-12 h-12 text-slate-400" />}
            title="No hay registros de auditoría"
            description="No se encontraron registros que coincidan con los filtros aplicados"
          />
        )}
      </Card>

      {/* Info Card */}
      <Card className="p-6 bg-blue-50 border-blue-200">
        <div className="flex items-start gap-3">
          <Shield className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-blue-900 mb-1">Información</h3>
            <p className="text-sm text-blue-800">
              Los registros de auditoría muestran todas las acciones administrativas realizadas en el
              sistema. Esta información es crítica para seguridad y cumplimiento normativo.
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
