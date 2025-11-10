import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useRafflesList } from '../../../hooks/useRaffles';
import { RaffleCard } from '../components/RaffleCard';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';
import { Alert } from '../../../components/ui/Alert';
import type { RaffleFilters } from '../../../types/raffle';

export function RafflesListPage() {
  const [filters, setFilters] = useState<RaffleFilters>({
    page: 1,
    page_size: 20,
  });

  const { data, isLoading, error } = useRafflesList(filters);

  const handleSearch = (search: string) => {
    setFilters((prev) => ({ ...prev, search, page: 1 }));
  };

  const handleStatusFilter = (status?: string) => {
    setFilters((prev) => ({
      ...prev,
      status: status as RaffleFilters['status'],
      page: 1,
    }));
  };

  const handlePageChange = (page: number) => {
    setFilters((prev) => ({ ...prev, page }));
  };

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Alert variant="error">
          <p className="font-semibold">Error al cargar sorteos</p>
          <p className="text-sm mt-1">
            {error instanceof Error ? error.message : 'Error desconocido'}
          </p>
        </Alert>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-white">
            Sorteos
          </h1>
          <p className="text-slate-600 dark:text-slate-400 mt-2">
            Descubre y participa en sorteos activos
          </p>
        </div>

        <Link to="/raffles/create">
          <Button>
            <svg
              className="w-5 h-5 mr-2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4v16m8-8H4"
              />
            </svg>
            Crear Sorteo
          </Button>
        </Link>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-slate-800 rounded-lg shadow p-6 mb-8">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Search */}
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
              Buscar
            </label>
            <Input
              type="text"
              placeholder="Buscar por título o descripción..."
              value={filters.search || ''}
              onChange={(e) => handleSearch(e.target.value)}
            />
          </div>

          {/* Status filter */}
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
              Estado
            </label>
            <select
              className="w-full px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={filters.status || ''}
              onChange={(e) => handleStatusFilter(e.target.value || undefined)}
            >
              <option value="">Todos los estados</option>
              <option value="active">Activos</option>
              <option value="draft">Borradores</option>
              <option value="completed">Completados</option>
              <option value="cancelled">Cancelados</option>
            </select>
          </div>
        </div>
      </div>

      {/* Loading state */}
      {isLoading && (
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      )}

      {/* Results */}
      {!isLoading && data && (
        <>
          {/* Grid */}
          {data.raffles.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
              {data.raffles.map((raffle) => (
                <RaffleCard key={raffle.id} raffle={raffle} />
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <svg
                className="w-16 h-16 mx-auto text-slate-400 mb-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
              <p className="text-lg text-slate-600 dark:text-slate-400">
                No se encontraron sorteos
              </p>
              <p className="text-sm text-slate-500 dark:text-slate-500 mt-2">
                Intenta cambiar los filtros o crear un nuevo sorteo
              </p>
            </div>
          )}

          {/* Pagination */}
          {data.pagination.total_pages > 1 && (
            <div className="flex items-center justify-center gap-2">
              <Button
                variant="outline"
                onClick={() => handlePageChange(data.pagination.page - 1)}
                disabled={data.pagination.page === 1}
              >
                Anterior
              </Button>

              <span className="px-4 py-2 text-sm text-slate-600 dark:text-slate-400">
                Página {data.pagination.page} de {data.pagination.total_pages}
              </span>

              <Button
                variant="outline"
                onClick={() => handlePageChange(data.pagination.page + 1)}
                disabled={data.pagination.page === data.pagination.total_pages}
              >
                Siguiente
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
