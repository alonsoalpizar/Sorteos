import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { CreditCard, AlertCircle, Eye, DollarSign } from "lucide-react";
import { useAdminPayments } from "../../hooks/useAdminPayments";
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
import type { PaymentFilters } from "../../types";

export function PaymentsListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState<PaymentFilters>({});

  const { data, isLoading, error } = useAdminPayments(filters, {
    page,
    limit: 20,
  });

  const handleFilterChange = (key: keyof PaymentFilters, value: any) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setPage(1);
  };

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      succeeded: "bg-green-100 text-green-700",
      pending: "bg-yellow-100 text-yellow-700",
      failed: "bg-red-100 text-red-700",
      refunded: "bg-gray-100 text-gray-700",
      disputed: "bg-orange-100 text-orange-700",
    };
    const labels: Record<string, string> = {
      succeeded: "Exitoso",
      pending: "Pendiente",
      failed: "Fallido",
      refunded: "Reembolsado",
      disputed: "Disputado",
    };
    return (
      <Badge className={styles[status] || "bg-slate-100 text-slate-700"}>
        {labels[status] || status}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900">Gestión de Pagos</h1>
        <p className="text-slate-600 mt-2">
          Administra transacciones, reembolsos y disputas
        </p>
      </div>

      {/* Stats Cards */}
      {data && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Total Procesado</p>
                <p className="text-2xl font-bold text-slate-900 mt-1">
                  {formatCurrency(data.total_amount || 0)}
                </p>
              </div>
              <DollarSign className="w-8 h-8 text-green-600" />
            </div>
          </Card>

          <Card className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Exitosos</p>
                <p className="text-2xl font-bold text-green-600 mt-1">
                  {data.succeeded_count || 0}
                </p>
              </div>
              <CreditCard className="w-8 h-8 text-green-600" />
            </div>
          </Card>

          <Card className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Reembolsados</p>
                <p className="text-2xl font-bold text-gray-600 mt-1">
                  {data.refunded_count || 0}
                </p>
              </div>
              <CreditCard className="w-8 h-8 text-gray-600" />
            </div>
          </Card>

          <Card className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Fallidos</p>
                <p className="text-2xl font-bold text-red-600 mt-1">
                  {data.failed_count || 0}
                </p>
              </div>
              <AlertCircle className="w-8 h-8 text-red-600" />
            </div>
          </Card>
        </div>
      )}

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
                handleFilterChange("status", e.target.value || undefined)
              }
            >
              <option value="">Todos</option>
              <option value="succeeded">Exitoso</option>
              <option value="pending">Pendiente</option>
              <option value="failed">Fallido</option>
              <option value="refunded">Reembolsado</option>
              <option value="disputed">Disputado</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Proveedor
            </label>
            <select
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={filters.provider || ""}
              onChange={(e) =>
                handleFilterChange("provider", e.target.value || undefined)
              }
            >
              <option value="">Todos</option>
              <option value="stripe">Stripe</option>
              <option value="paypal">PayPal</option>
              <option value="mercadopago">MercadoPago</option>
              <option value="pagadito">Pagadito</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Buscar
            </label>
            <Input
              placeholder="Payment Intent, Order ID..."
              value={filters.search || ""}
              onChange={(e) => handleFilterChange("search", e.target.value)}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Opciones
            </label>
            <div className="flex items-center gap-2 mt-2">
              <input
                type="checkbox"
                id="include_refund"
                checked={filters.include_refund || false}
                onChange={(e) =>
                  handleFilterChange("include_refund", e.target.checked)
                }
                className="w-4 h-4 text-blue-600 border-slate-300 rounded focus:ring-blue-500"
              />
              <label htmlFor="include_refund" className="text-sm text-slate-700">
                Incluir reembolsados
              </label>
            </div>
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
              title="Error al cargar pagos"
              description={(error as Error).message}
            />
          </div>
        ) : !data || !data.payments || data.payments.length === 0 ? (
          <div className="p-6">
            <EmptyState
              icon={<CreditCard className="w-12 h-12 text-slate-400" />}
              title="No se encontraron pagos"
              description="Intenta ajustar los filtros de búsqueda"
            />
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID Pago</TableHead>
                  <TableHead>Usuario</TableHead>
                  <TableHead>Rifa</TableHead>
                  <TableHead>Organizador</TableHead>
                  <TableHead className="text-right">Monto</TableHead>
                  <TableHead>Estado</TableHead>
                  <TableHead>Proveedor</TableHead>
                  <TableHead>Fecha</TableHead>
                  <TableHead className="text-center">Acciones</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.payments.map((item) => (
                  <TableRow
                    key={item.payment.id}
                    onClick={() => navigate(`/admin/payments/${item.payment.id}`)}
                    className="cursor-pointer hover:bg-slate-50"
                  >
                    <TableCell className="font-mono text-sm">
                      {item.payment.id.substring(0, 8)}...
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="text-sm font-medium text-slate-900">
                          {item.user_name}
                        </p>
                        <p className="text-xs text-slate-500">{item.user_email}</p>
                      </div>
                    </TableCell>
                    <TableCell className="text-sm text-slate-900">
                      {item.raffle_title}
                    </TableCell>
                    <TableCell className="text-sm text-slate-600">
                      {item.organizer_name}
                    </TableCell>
                    <TableCell className="text-right font-medium text-slate-900">
                      {formatCurrency(item.payment.amount)}
                    </TableCell>
                    <TableCell>{getStatusBadge(item.payment.status)}</TableCell>
                    <TableCell className="text-sm text-slate-600 capitalize">
                      {item.payment.provider || "—"}
                    </TableCell>
                    <TableCell className="text-sm text-slate-600">
                      {item.payment.created_at
                        ? format(new Date(item.payment.created_at), "dd/MM/yyyy HH:mm")
                        : "—"}
                    </TableCell>
                    <TableCell className="text-center">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={(e) => {
                          e.stopPropagation();
                          navigate(`/admin/payments/${item.payment.id}`);
                        }}
                      >
                        <Eye className="w-4 h-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {/* Pagination */}
            {data && data.total_pages > 1 && (
              <div className="flex items-center justify-between px-6 py-4 border-t border-slate-200">
                <p className="text-sm text-slate-600">
                  Mostrando {data.payments.length} de {data.total} pagos
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
                    Página {page} de {data.total_pages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => p + 1)}
                    disabled={page >= data.total_pages}
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
