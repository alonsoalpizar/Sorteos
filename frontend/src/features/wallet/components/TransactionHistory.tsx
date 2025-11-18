import { format } from 'date-fns';
import { es } from 'date-fns/locale';
import { ChevronLeft, ChevronRight, RefreshCw } from 'lucide-react';
import { useTransactionHistory } from '../hooks/useTransactionHistory';
import {
  formatCRC,
  translateTransactionType,
  translateTransactionStatus,
  getStatusColor,
} from '../../../types/wallet';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';
import { Badge } from '../../../components/ui/Badge';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { EmptyState } from '../../../components/ui/EmptyState';

export const TransactionHistory = () => {
  const {
    transactions,
    pagination,
    isLoading,
    error,
    refetch,
    hasNextPage,
    hasPreviousPage,
    nextPage,
    previousPage,
    currentPage,
    totalPages,
  } = useTransactionHistory(20);

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
          <LoadingSpinner size="lg" />
        </div>
      </Card>
    );
  }

  if (transactions.length === 0) {
    return (
      <Card className="p-6">
        <EmptyState
          icon="📊"
          title="No hay transacciones"
          description="Aún no has realizado ninguna transacción en tu billetera"
        />
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header con botón de refrescar */}
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-slate-900">
          Historial de Transacciones ({pagination.total})
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
              {transactions.map((tx) => {
                const isDebit = ['purchase', 'withdrawal', 'adjustment'].includes(tx.type);
                const date = new Date(tx.created_at);

                return (
                  <tr key={tx.id} className="hover:bg-slate-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-600">
                      {format(date, 'dd MMM yyyy, HH:mm', { locale: es })}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-900">
                      {translateTransactionType(tx.type)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span
                        className={`font-semibold ${
                          isDebit ? 'text-red-600' : 'text-green-600'
                        }`}
                      >
                        {isDebit ? '-' : '+'} {formatCRC(tx.amount)}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <Badge variant={getStatusColor(tx.status)}>
                        {translateTransactionStatus(tx.status)}
                      </Badge>
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
        {transactions.map((tx) => {
          const isDebit = ['purchase', 'withdrawal', 'adjustment'].includes(tx.type);
          const date = new Date(tx.created_at);

          return (
            <Card key={tx.id} className="p-4">
              <div className="flex items-start justify-between mb-2">
                <div>
                  <p className="font-medium text-slate-900">{translateTransactionType(tx.type)}</p>
                  <p className="text-xs text-slate-500">
                    {format(date, 'dd MMM yyyy, HH:mm', { locale: es })}
                  </p>
                </div>
                <Badge variant={getStatusColor(tx.status)}>
                  {translateTransactionStatus(tx.status)}
                </Badge>
              </div>
              <div className="flex items-center justify-between mt-3 pt-3 border-t border-slate-200">
                <span
                  className={`text-lg font-bold ${isDebit ? 'text-red-600' : 'text-green-600'}`}
                >
                  {isDebit ? '-' : '+'} {formatCRC(tx.amount)}
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
              Página {currentPage + 1} de {totalPages}
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={previousPage}
                disabled={!hasPreviousPage}
              >
                <ChevronLeft className="w-4 h-4" />
                Anterior
              </Button>
              <Button variant="outline" size="sm" onClick={nextPage} disabled={!hasNextPage}>
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
