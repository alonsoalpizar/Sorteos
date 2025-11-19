import { useState } from "react";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import { ChevronLeft, ChevronRight, RefreshCw, History } from "lucide-react";
import { useWalletTransactions } from "../hooks/useWallet";
import type { TransactionType, TransactionStatus } from "../types";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";

// Helper para formatear CRC
function formatCRC(amount: number | string): string {
  const num = typeof amount === "string" ? parseFloat(amount) : amount;
  return new Intl.NumberFormat("es-CR", {
    style: "currency",
    currency: "CRC",
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(num);
}

// Helper para traducir tipos de transacción
function translateTransactionType(type: TransactionType): string {
  const types: Record<TransactionType, string> = {
    deposit: "Recarga",
    withdrawal: "Retiro",
    purchase: "Compra de boletos",
    refund: "Reembolso",
    prize_claim: "Premio ganado",
    settlement_payout: "Pago de liquidación",
    adjustment: "Ajuste",
  };
  return types[type] || type;
}

// Helper para traducir estados
function translateTransactionStatus(status: TransactionStatus): string {
  const statuses: Record<TransactionStatus, string> = {
    pending: "Pendiente",
    completed: "Completada",
    failed: "Fallida",
    reversed: "Revertida",
  };
  return statuses[status] || status;
}

// Helper para obtener badge color según estado
function getStatusBadgeClass(status: TransactionStatus): string {
  const classes: Record<TransactionStatus, string> = {
    pending: "bg-yellow-100 text-yellow-700",
    completed: "bg-green-100 text-green-700",
    failed: "bg-red-100 text-red-700",
    reversed: "bg-slate-100 text-slate-700",
  };
  return classes[status] || "bg-slate-100 text-slate-700";
}

export const TransactionHistory = () => {
  const [page, setPage] = useState(0);
  const limit = 20;

  const { data, isLoading, error, refetch } = useWalletTransactions({
    limit,
    offset: page * limit,
  });

  if (error) {
    return (
      <Card className="p-6">
        <div className="text-center text-red-600">
          <p>Error al cargar las transacciones</p>
          <Button variant="outline" size="sm" onClick={() => refetch()} className="mt-4">
            Reintentar
          </Button>
        </div>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center py-8">
          <LoadingSpinner />
        </div>
      </Card>
    );
  }

  if (!data || data.transactions.length === 0) {
    return (
      <Card className="p-6">
        <EmptyState
          icon={<History className="w-12 h-12 text-slate-400" />}
          title="No hay transacciones"
          description="Aún no has realizado ninguna transacción en tu billetera"
        />
      </Card>
    );
  }

  const totalPages = Math.ceil(data.pagination.total / limit);
  const hasNextPage = page < totalPages - 1;
  const hasPreviousPage = page > 0;

  return (
    <div className="space-y-4">
      {/* Header con botón de refrescar */}
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-slate-900">
          Historial de Transacciones ({data.pagination.total})
        </h2>
        <Button variant="ghost" size="sm" onClick={() => refetch()}>
          <RefreshCw className="w-4 h-4" />
        </Button>
      </div>

      {/* Tabla de transacciones (Desktop) */}
      <Card className="overflow-hidden hidden md:block">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                  Fecha
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                  Tipo
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                  Monto
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                  Estado
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                  Saldo después
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-slate-200">
              {data.transactions.map((tx) => {
                const isDebit = ["purchase", "withdrawal", "adjustment"].includes(tx.type);
                const date = new Date(tx.created_at);

                return (
                  <tr key={tx.id} className="hover:bg-slate-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-600">
                      {format(date, "dd MMM yyyy, HH:mm", { locale: es })}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-900">
                      {translateTransactionType(tx.type)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span
                        className={`font-semibold ${
                          isDebit ? "text-red-600" : "text-green-600"
                        }`}
                      >
                        {isDebit ? "-" : "+"} {formatCRC(tx.amount)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusBadgeClass(
                          tx.status
                        )}`}
                      >
                        {translateTransactionStatus(tx.status)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-600">
                      {formatCRC(tx.balance_after)}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </Card>

      {/* Lista de transacciones (Mobile) */}
      <div className="md:hidden space-y-3">
        {data.transactions.map((tx) => {
          const isDebit = ["purchase", "withdrawal", "adjustment"].includes(tx.type);
          const date = new Date(tx.created_at);

          return (
            <Card key={tx.id} className="p-4">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <p className="font-medium text-slate-900">{translateTransactionType(tx.type)}</p>
                  <p className="text-xs text-slate-500">
                    {format(date, "dd MMM yyyy, HH:mm", { locale: es })}
                  </p>
                </div>
                <span
                  className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusBadgeClass(
                    tx.status
                  )}`}
                >
                  {translateTransactionStatus(tx.status)}
                </span>
              </div>
              <div className="flex items-center justify-between mt-3 pt-3 border-t border-slate-200">
                <span
                  className={`text-lg font-bold ${isDebit ? "text-red-600" : "text-green-600"}`}
                >
                  {isDebit ? "-" : "+"} {formatCRC(tx.amount)}
                </span>
                <span className="text-sm text-slate-600">
                  Saldo: {formatCRC(tx.balance_after)}
                </span>
              </div>
            </Card>
          );
        })}
      </div>

      {/* Paginación */}
      {totalPages > 1 && (
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
