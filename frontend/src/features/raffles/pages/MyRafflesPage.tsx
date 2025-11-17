import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRafflesList } from '@/hooks/useRaffles';
import { useUser } from '@/hooks/useAuth';
import { Button } from '@/components/ui/Button';
import { EmptyState } from '@/components/ui/EmptyState';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import type { RaffleStatus, Raffle } from '@/types/raffle';

const statusLabels: Record<RaffleStatus, string> = {
  draft: 'Borrador',
  active: 'Activo',
  suspended: 'Suspendido',
  completed: 'Completado',
  cancelled: 'Cancelado',
};

const statusColors: Record<RaffleStatus, string> = {
  draft: 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300',
  active: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  suspended: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
  completed: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
  cancelled: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
};

export function MyRafflesPage() {
  const navigate = useNavigate();
  const user = useUser();
  const [filterStatus, setFilterStatus] = useState<RaffleStatus | 'all'>('all');

  // Fetch raffles created by current user
  const { data: rafflesData, isLoading, error } = useRafflesList({
    user_id: user?.id,
    status: filterStatus === 'all' ? undefined : filterStatus,
  });

  if (isLoading) {
    return <LoadingSpinner text="Cargando tus sorteos..." />;
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">
          Error al cargar tus sorteos. Por favor, intenta nuevamente.
        </p>
      </div>
    );
  }

  const raffles = rafflesData?.raffles || [];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-white">
            Mis Sorteos
          </h1>
          <p className="text-slate-600 dark:text-slate-400 mt-2">
            Gestiona los sorteos que has creado
          </p>
        </div>
        <Button onClick={() => navigate('/raffles/create')}>
          <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Crear Sorteo
        </Button>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-4">
        <div className="flex items-center gap-2 overflow-x-auto">
          <span className="text-sm font-medium text-slate-700 dark:text-slate-300 whitespace-nowrap">
            Filtrar por:
          </span>
          <button
            onClick={() => setFilterStatus('all')}
            className={`px-3 py-1.5 text-sm font-medium rounded-lg transition-colors whitespace-nowrap ${
              filterStatus === 'all'
                ? 'bg-blue-600 text-white'
                : 'bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600'
            }`}
          >
            Todos
          </button>
          {Object.entries(statusLabels).map(([status, label]) => (
            <button
              key={status}
              onClick={() => setFilterStatus(status as RaffleStatus)}
              className={`px-3 py-1.5 text-sm font-medium rounded-lg transition-colors whitespace-nowrap ${
                filterStatus === status
                  ? 'bg-blue-600 text-white'
                  : 'bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600'
              }`}
            >
              {label}
            </button>
          ))}
        </div>
      </div>

      {/* Raffles List */}
      {raffles.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
            </svg>
          }
          title={filterStatus === 'all' ? 'No tienes sorteos' : `No tienes sorteos ${statusLabels[filterStatus]?.toLowerCase()}`}
          description={filterStatus === 'all' ? 'Crea tu primer sorteo para comenzar' : 'Intenta cambiar el filtro para ver otros sorteos'}
          action={
            filterStatus === 'all'
              ? {
                  label: 'Crear Sorteo',
                  onClick: () => navigate('/raffles/create'),
                }
              : undefined
          }
        />
      ) : (
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                    Sorteo
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                    Estado
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                    Progreso
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                    Ingresos
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                    Sorteo
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-slate-700 dark:text-slate-300 uppercase tracking-wider">
                    Acciones
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-200 dark:divide-slate-700">
                {raffles.map((raffle: Raffle) => {
                  const soldPercentage = (raffle.sold_count / raffle.total_numbers) * 100;
                  const daysUntilDraw = Math.ceil(
                    (new Date(raffle.draw_date).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24)
                  );

                  return (
                    <tr
                      key={raffle.id}
                      className="hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors cursor-pointer"
                      onClick={() => navigate(`/raffles/${raffle.uuid}`)}
                    >
                      <td className="px-6 py-4">
                        <div>
                          <p className="font-medium text-slate-900 dark:text-white">
                            {raffle.title}
                          </p>
                          <p className="text-sm text-slate-500 dark:text-slate-400">
                            {raffle.total_numbers} números
                          </p>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            statusColors[raffle.status]
                          }`}
                        >
                          {statusLabels[raffle.status]}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <div className="space-y-1">
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-slate-600 dark:text-slate-400">
                              {raffle.sold_count} vendidos
                            </span>
                            <span className="font-medium text-slate-900 dark:text-white">
                              {soldPercentage.toFixed(0)}%
                            </span>
                          </div>
                          <div className="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-2">
                            <div
                              className="bg-blue-600 h-2 rounded-full transition-all"
                              style={{ width: `${soldPercentage}%` }}
                            />
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div>
                          <p className="font-semibold text-slate-900 dark:text-white">
                            ₡{parseFloat(raffle.total_revenue || '0').toLocaleString()}
                          </p>
                          <p className="text-xs text-slate-500 dark:text-slate-400">
                            ₡{parseFloat(raffle.price_per_number).toFixed(2)}/número
                          </p>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div>
                          <p className="text-sm text-slate-900 dark:text-white">
                            {new Date(raffle.draw_date).toLocaleDateString()}
                          </p>
                          {raffle.status === 'active' && daysUntilDraw > 0 && (
                            <p className="text-xs text-slate-500 dark:text-slate-400">
                              En {daysUntilDraw} {daysUntilDraw === 1 ? 'día' : 'días'}
                            </p>
                          )}
                        </div>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/raffles/${raffle.uuid}`);
                          }}
                        >
                          Ver Detalles
                        </Button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Pagination */}
      {rafflesData && rafflesData.pagination.total_pages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={rafflesData.pagination.page === 1}
          >
            Anterior
          </Button>
          <span className="text-sm text-slate-600 dark:text-slate-400">
            Página {rafflesData.pagination.page} de {rafflesData.pagination.total_pages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={rafflesData.pagination.page === rafflesData.pagination.total_pages}
          >
            Siguiente
          </Button>
        </div>
      )}
    </div>
  );
}
