import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Package, AlertCircle, Eye } from "lucide-react";
import { useAdminRaffles } from "../../hooks/useAdminRaffles";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/Table";
import { formatCurrency } from "@/lib/currency";
import { format } from "date-fns";
import type { RaffleFilters } from "../../types";

export function RafflesListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState<RaffleFilters>({});

  const { data, isLoading, error } = useAdminRaffles(filters, {
    page,
    limit: 20,
  });

  const handleFilterChange = (key: keyof RaffleFilters, value: any) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setPage(1); // Reset to first page when filters change
  };

  const getStatusBadge = (status: string) => {
    const styles = {
      draft: "bg-slate-100 text-slate-700",
      active: "bg-green-100 text-green-700",
      suspended: "bg-red-100 text-red-700",
      completed: "bg-blue-100 text-blue-700",
      cancelled: "bg-gray-100 text-gray-700",
    };
    const labels = {
      draft: "Borrador",
      active: "Activa",
      suspended: "Suspendida",
      completed: "Completada",
      cancelled: "Cancelada",
    };
    return (
      <Badge className={styles[status as keyof typeof styles] || ""}>
        {labels[status as keyof typeof labels] || status}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900">Gestión de Rifas</h1>
        <p className="text-slate-600 mt-2">
          Administra rifas, suspensiones y sorteos manuales
        </p>
      </div>

      {/* Filters */}
      <Card className="p-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Estado
            </label>
            <select
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={filters.status || ""}
              onChange={(e) =>
                handleFilterChange(
                  "status",
                  e.target.value || undefined
                )
              }
            >
              <option value="">Todos los estados</option>
              <option value="draft">Borrador</option>
              <option value="active">Activa</option>
              <option value="suspended">Suspendida</option>
              <option value="completed">Completada</option>
              <option value="cancelled">Cancelada</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Buscar
            </label>
            <Input
              placeholder="Buscar por título..."
              value={filters.search || ""}
              onChange={(e) => handleFilterChange("search", e.target.value)}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              ID Organizador
            </label>
            <Input
              type="number"
              placeholder="ID organizador..."
              value={filters.organizer_id || ""}
              onChange={(e) =>
                handleFilterChange(
                  "organizer_id",
                  e.target.value ? parseInt(e.target.value) : undefined
                )
              }
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              ID Categoría
            </label>
            <Input
              type="number"
              placeholder="ID categoría..."
              value={filters.category_id || ""}
              onChange={(e) =>
                handleFilterChange(
                  "category_id",
                  e.target.value ? parseInt(e.target.value) : undefined
                )
              }
            />
          </div>
        </div>
      </Card>

      {/* Table */}
      <Card>
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <LoadingSpinner />
          </div>
        ) : error ? (
          <div className="p-6">
            <EmptyState
              icon={<AlertCircle className="w-12 h-12 text-red-500" />}
              title="Error al cargar rifas"
              description={error.message}
            />
          </div>
        ) : !data || !data.data || data.data.length === 0 ? (
          <div className="p-6">
            <EmptyState
              icon={<Package className="w-12 h-12 text-slate-400" />}
              title="No se encontraron rifas"
              description="Intenta ajustar los filtros de búsqueda"
            />
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Título</TableHead>
                  <TableHead>Organizador</TableHead>
                  <TableHead>Estado</TableHead>
                  <TableHead>Fecha Sorteo</TableHead>
                  <TableHead className="text-right">Vendidos</TableHead>
                  <TableHead className="text-right">Conversión</TableHead>
                  <TableHead className="text-right">Ingresos</TableHead>
                  <TableHead className="text-center">Acciones</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.data.map((item) => (
                  <TableRow
                    key={item.raffle.ID}
                    onClick={() => navigate(`/admin/raffles/${item.raffle.ID}`)}
                    className="cursor-pointer hover:bg-slate-50"
                  >
                    <TableCell className="font-mono text-sm">
                      #{item.raffle.ID}
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="font-medium text-slate-900">
                          {item.raffle.Title}
                        </p>
                        <p className="text-sm text-slate-500">
                          {item.sold_count}/{item.raffle.TotalNumbers} boletos
                        </p>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="text-sm font-medium text-slate-900">
                          {item.organizer_name}
                        </p>
                        <p className="text-xs text-slate-500">
                          {item.organizer_email}
                        </p>
                      </div>
                    </TableCell>
                    <TableCell>{getStatusBadge(item.raffle.Status)}</TableCell>
                    <TableCell className="text-sm text-slate-600">
                      {item.raffle.DrawDate ? format(new Date(item.raffle.DrawDate), "dd/MM/yyyy HH:mm") : "N/A"}
                    </TableCell>
                    <TableCell className="text-right text-sm">
                      <span className="font-medium text-slate-900">
                        {item.sold_count}
                      </span>
                      <span className="text-slate-500">
                        /{item.raffle.TotalNumbers}
                      </span>
                    </TableCell>
                    <TableCell className="text-right">
                      <span
                        className={`font-medium ${
                          item.conversion_rate >= 80
                            ? "text-green-600"
                            : item.conversion_rate >= 50
                            ? "text-blue-600"
                            : "text-amber-600"
                        }`}
                      >
                        {item.conversion_rate.toFixed(1)}%
                      </span>
                    </TableCell>
                    <TableCell className="text-right font-medium text-slate-900">
                      {formatCurrency(item.total_revenue)}
                    </TableCell>
                    <TableCell className="text-center">
                      <div className="flex items-center justify-center gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/admin/raffles/${item.raffle.ID}`);
                          }}
                        >
                          <Eye className="w-4 h-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {/* Pagination */}
            {data.pagination && data.pagination.total_pages > 1 && (
              <div className="flex items-center justify-between px-6 py-4 border-t border-slate-200">
                <p className="text-sm text-slate-600">
                  Mostrando {data.data.length} de {data.pagination.total} rifas
                </p>
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                  >
                    Anterior
                  </Button>
                  <span className="px-4 py-2 text-sm text-slate-700">
                    Página {page} de {data.pagination.total_pages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => p + 1)}
                    disabled={page >= data.pagination.total_pages}
                  >
                    Siguiente
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
      </Card>
    </div>
  );
}
