import { useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { useRafflesList } from '../../../hooks/useRaffles';
import { RaffleCard } from '../components/RaffleCard';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { EmptyState } from '../../../components/ui/EmptyState';
import type { RaffleFilters } from '../../../types/raffle';

export function RafflesListPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchInput, setSearchInput] = useState(searchParams.get('search') || '');

  const [filters, setFilters] = useState<RaffleFilters>({
    page: 1,
    page_size: 12,
    search: searchParams.get('search') || undefined,
    status: (searchParams.get('status') as RaffleFilters['status']) || undefined,
  });

  const { data, isLoading, error } = useRafflesList(filters);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setFilters((prev) => ({ ...prev, search: searchInput || undefined, page: 1 }));
    if (searchInput) {
      searchParams.set('search', searchInput);
    } else {
      searchParams.delete('search');
    }
    setSearchParams(searchParams);
  };

  const handleStatusFilter = (status?: string) => {
    setFilters((prev) => ({
      ...prev,
      status: status as RaffleFilters['status'],
      page: 1,
    }));
    if (status) {
      searchParams.set('status', status);
    } else {
      searchParams.delete('status');
    }
    setSearchParams(searchParams);
  };

  const handlePageChange = (page: number) => {
    setFilters((prev) => ({ ...prev, page }));
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400 font-medium mb-2">
          Error al cargar sorteos
        </p>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          {error instanceof Error ? error.message : 'Error desconocido'}
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900 dark:text-white">
            Explorar Sorteos
          </h1>
          <p className="text-slate-600 dark:text-slate-400 mt-2">
            Descubre y participa en sorteos activos
          </p>
        </div>

        <Link to="/raffles/create">
          <Button className="hidden md:flex">
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

      {/* Search and Filters */}
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
        <div className="space-y-4">
          {/* Search Bar */}
          <form onSubmit={handleSearch}>
            <div className="relative">
              <input
                type="text"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                placeholder="Buscar por título o descripción..."
                className="w-full pl-10 pr-4 py-3 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-900 text-slate-900 dark:text-white placeholder-slate-500 dark:placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
              />
              <svg
                className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
              {searchInput && (
                <button
                  type="button"
                  onClick={() => {
                    setSearchInput('');
                    setFilters((prev) => ({ ...prev, search: undefined, page: 1 }));
                    searchParams.delete('search');
                    setSearchParams(searchParams);
                  }}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300"
                >
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
            </div>
          </form>

          {/* Status Filters */}
          <div>
            <p className="text-sm font-medium text-slate-700 dark:text-slate-300 mb-3">
              Filtrar por estado:
            </p>
            <div className="flex flex-wrap gap-2">
              <button
                onClick={() => handleStatusFilter(undefined)}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  !filters.status
                    ? 'bg-blue-600 text-white'
                    : 'bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600'
                }`}
              >
                Todos
              </button>
              <button
                onClick={() => handleStatusFilter('active')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  filters.status === 'active'
                    ? 'bg-blue-600 text-white'
                    : 'bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600'
                }`}
              >
                <span className="inline-block w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                Activos
              </button>
              <button
                onClick={() => handleStatusFilter('completed')}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  filters.status === 'completed'
                    ? 'bg-blue-600 text-white'
                    : 'bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600'
                }`}
              >
                Completados
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Results Count */}
      {!isLoading && data && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-slate-600 dark:text-slate-400">
            {data.pagination.total === 0
              ? 'No se encontraron sorteos'
              : `${data.pagination.total} ${data.pagination.total === 1 ? 'sorteo encontrado' : 'sorteos encontrados'}`}
          </p>
        </div>
      )}

      {/* Loading State */}
      {isLoading && <LoadingSpinner text="Cargando sorteos..." />}

      {/* Results Grid */}
      {!isLoading && data && (
        <>
          {data.raffles.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data.raffles.map((raffle) => (
                <RaffleCard key={raffle.id} raffle={raffle} />
              ))}
            </div>
          ) : (
            <EmptyState
              icon={
                <svg className="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
              }
              title="No se encontraron sorteos"
              description={
                filters.search || filters.status
                  ? 'Intenta cambiar los filtros de búsqueda'
                  : 'No hay sorteos disponibles en este momento'
              }
              action={
                filters.search || filters.status
                  ? {
                      label: 'Limpiar Filtros',
                      onClick: () => {
                        setSearchInput('');
                        setFilters({ page: 1, page_size: 12 });
                        setSearchParams({});
                      },
                    }
                  : undefined
              }
            />
          )}

          {/* Pagination */}
          {data.pagination.total_pages > 1 && (
            <div className="flex items-center justify-center gap-2 pt-4">
              <Button
                variant="outline"
                onClick={() => handlePageChange(data.pagination.page - 1)}
                disabled={data.pagination.page === 1}
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
                <span className="ml-2 hidden sm:inline">Anterior</span>
              </Button>

              <div className="flex items-center gap-2">
                {Array.from({ length: Math.min(5, data.pagination.total_pages) }, (_, i) => {
                  let pageNum;
                  if (data.pagination.total_pages <= 5) {
                    pageNum = i + 1;
                  } else if (data.pagination.page <= 3) {
                    pageNum = i + 1;
                  } else if (data.pagination.page >= data.pagination.total_pages - 2) {
                    pageNum = data.pagination.total_pages - 4 + i;
                  } else {
                    pageNum = data.pagination.page - 2 + i;
                  }

                  return (
                    <button
                      key={pageNum}
                      onClick={() => handlePageChange(pageNum)}
                      className={`w-10 h-10 rounded-lg font-medium text-sm transition-colors ${
                        data.pagination.page === pageNum
                          ? 'bg-blue-600 text-white'
                          : 'bg-slate-100 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-600'
                      }`}
                    >
                      {pageNum}
                    </button>
                  );
                })}
              </div>

              <Button
                variant="outline"
                onClick={() => handlePageChange(data.pagination.page + 1)}
                disabled={data.pagination.page === data.pagination.total_pages}
              >
                <span className="mr-2 hidden sm:inline">Siguiente</span>
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </Button>
            </div>
          )}
        </>
      )}

      {/* Mobile Create Button */}
      <Link to="/raffles/create" className="md:hidden fixed bottom-6 right-6 z-40">
        <Button size="lg" className="rounded-full shadow-lg">
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </Button>
      </Link>
    </div>
  );
}
