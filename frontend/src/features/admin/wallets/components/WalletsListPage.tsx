import { useState } from "react";
import { Link } from "react-router-dom";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import {
  Wallet,
  Search,
  RefreshCw,
  ChevronLeft,
  ChevronRight,
  Eye,
  Filter,
  AlertCircle,
} from "lucide-react";
import { useAdminWallets } from "../hooks/useAdminWallets";
import type { WalletStatus } from "../types";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";

// Helper para formatear CRC
function formatCRC(amount: string): string {
  const num = parseFloat(amount);
  return new Intl.NumberFormat("es-CR", {
    style: "currency",
    currency: "CRC",
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(num);
}

// Helper para traducir status
function translateStatus(status: WalletStatus): string {
  const statuses: Record<WalletStatus, string> = {
    active: "Activa",
    frozen: "Congelada",
    closed: "Cerrada",
  };
  return statuses[status] || status;
}

// Helper para obtener badge color según status
function getStatusBadgeClass(status: WalletStatus): string {
  const classes: Record<WalletStatus, string> = {
    active: "bg-green-100 text-green-700",
    frozen: "bg-red-100 text-red-700",
    closed: "bg-slate-100 text-slate-700",
  };
  return classes[status] || "bg-slate-100 text-slate-700";
}

export const WalletsListPage = () => {
  const [page, setPage] = useState(0);
  const [limit] = useState(20);
  const [statusFilter, setStatusFilter] = useState<WalletStatus | "">("");
  const [emailSearch, setEmailSearch] = useState("");
  const [emailInput, setEmailInput] = useState("");

  const { data, isLoading, error, refetch } = useAdminWallets({
    page,
    limit,
    status: statusFilter || undefined,
    user_email: emailSearch || undefined,
  });

  const handleSearch = () => {
    setEmailSearch(emailInput);
    setPage(0); // Reset to first page on search
  };

  const handleClearFilters = () => {
    setEmailInput("");
    setEmailSearch("");
    setStatusFilter("");
    setPage(0);
  };

  if (error) {
    return (
      <div className="p-6">
        <Card className="p-6">
          <div className="text-center text-red-600">
            <AlertCircle className="w-12 h-12 mx-auto mb-4" />
            <p className="font-semibold mb-2">Error al cargar las billeteras</p>
            <p className="text-sm">{error instanceof Error ? error.message : "Error desconocido"}</p>
            <Button variant="outline" size="sm" onClick={() => refetch()} className="mt-4">
              Reintentar
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  const totalPages = data ? Math.ceil(data.pagination.total / limit) : 0;
  const hasNextPage = page < totalPages - 1;
  const hasPreviousPage = page > 0;
  const hasActiveFilters = statusFilter !== "" || emailSearch !== "";

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-blue-100 rounded-lg">
            <Wallet className="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-900">Billeteras</h1>
            <p className="text-sm text-slate-600">
              Administración de billeteras de usuarios
            </p>
          </div>
        </div>
        <Button variant="ghost" size="sm" onClick={() => refetch()}>
          <RefreshCw className="w-4 h-4 mr-2" />
          Actualizar
        </Button>
      </div>

      {/* Filters */}
      <Card className="p-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Email search */}
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Buscar por email
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                placeholder="usuario@ejemplo.com"
                value={emailInput}
                onChange={(e) => setEmailInput(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSearch()}
                className="flex-1 px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
              />
              <Button size="sm" onClick={handleSearch}>
                <Search className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Status filter */}
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Estado
            </label>
            <select
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value as WalletStatus | "");
                setPage(0);
              }}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            >
              <option value="">Todos</option>
              <option value="active">Activas</option>
              <option value="frozen">Congeladas</option>
              <option value="closed">Cerradas</option>
            </select>
          </div>

          {/* Clear filters */}
          <div className="flex items-end">
            {hasActiveFilters && (
              <Button variant="outline" size="sm" onClick={handleClearFilters} className="w-full">
                <Filter className="w-4 h-4 mr-2" />
                Limpiar filtros
              </Button>
            )}
          </div>
        </div>
      </Card>

      {/* Results count */}
      {data && (
        <div className="text-sm text-slate-600">
          Mostrando {data.wallets.length} de {data.pagination.total} billeteras
        </div>
      )}

      {/* Loading state */}
      {isLoading && (
        <Card className="p-6">
          <div className="flex items-center justify-center py-8">
            <LoadingSpinner />
          </div>
        </Card>
      )}

      {/* Table (Desktop) */}
      {!isLoading && data && (
        <>
          <Card className="overflow-hidden hidden md:block">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-slate-50 border-b border-slate-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                      Usuario
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                      Saldo Disponible
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                      Ganancias
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                      Total
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                      Estado
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                      Acciones
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-slate-200">
                  {data.wallets.map((wallet) => {
                    const createdDate = new Date(wallet.created_at);

                    return (
                      <tr key={wallet.id} className="hover:bg-slate-50">
                        <td className="px-6 py-4">
                          <div>
                            <p className="text-sm font-medium text-slate-900">{wallet.user_name}</p>
                            <p className="text-xs text-slate-500">{wallet.user_email}</p>
                            <p className="text-xs text-slate-400 mt-1">
                              Creada: {format(createdDate, "dd MMM yyyy", { locale: es })}
                            </p>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-900">
                          {formatCRC(wallet.balance_available)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-900">
                          {formatCRC(wallet.earnings_balance)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-slate-900">
                          {formatCRC(wallet.total_balance)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span
                            className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusBadgeClass(
                              wallet.status
                            )}`}
                          >
                            {translateStatus(wallet.status)}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <Link to={`/admin/wallets/${wallet.id}`}>
                            <Button variant="ghost" size="sm">
                              <Eye className="w-4 h-4 mr-1" />
                              Ver detalles
                            </Button>
                          </Link>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </Card>

          {/* Cards (Mobile) */}
          <div className="md:hidden space-y-3">
            {data.wallets.map((wallet) => {
              const createdDate = new Date(wallet.created_at);

              return (
                <Card key={wallet.id} className="p-4">
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <p className="font-medium text-slate-900">{wallet.user_name}</p>
                      <p className="text-xs text-slate-500">{wallet.user_email}</p>
                      <p className="text-xs text-slate-400 mt-1">
                        Creada: {format(createdDate, "dd MMM yyyy", { locale: es })}
                      </p>
                    </div>
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusBadgeClass(
                        wallet.status
                      )}`}
                    >
                      {translateStatus(wallet.status)}
                    </span>
                  </div>

                  <div className="space-y-2 text-sm mb-3">
                    <div className="flex justify-between">
                      <span className="text-slate-600">Saldo disponible:</span>
                      <span className="font-medium">{formatCRC(wallet.balance_available)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-600">Ganancias:</span>
                      <span className="font-medium">{formatCRC(wallet.earnings_balance)}</span>
                    </div>
                    <div className="flex justify-between border-t border-slate-200 pt-2">
                      <span className="font-medium">Total:</span>
                      <span className="font-semibold">{formatCRC(wallet.total_balance)}</span>
                    </div>
                  </div>

                  <Link to={`/admin/wallets/${wallet.id}`} className="block">
                    <Button variant="outline" size="sm" className="w-full">
                      <Eye className="w-4 h-4 mr-2" />
                      Ver detalles
                    </Button>
                  </Link>
                </Card>
              );
            })}
          </div>
        </>
      )}

      {/* Empty state */}
      {!isLoading && data && data.wallets.length === 0 && (
        <Card className="p-6">
          <div className="text-center text-slate-600">
            <Wallet className="w-12 h-12 mx-auto mb-4 text-slate-400" />
            <p className="font-semibold mb-2">No se encontraron billeteras</p>
            <p className="text-sm">Intenta ajustar los filtros de búsqueda</p>
          </div>
        </Card>
      )}

      {/* Pagination */}
      {data && totalPages > 1 && (
        <Card className="p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm text-slate-600">
              Página {page + 1} de {totalPages}
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => p - 1)}
                disabled={!hasPreviousPage}
              >
                <ChevronLeft className="w-4 h-4" />
                Anterior
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => p + 1)}
                disabled={!hasNextPage}
              >
                Siguiente
                <ChevronRight className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </Card>
      )}
    </div>
  );
};
