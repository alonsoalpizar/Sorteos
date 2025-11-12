import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUser } from '@/hooks/useAuth';
import { useRafflesList } from '@/hooks/useRaffles';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { RaffleCard } from '@/features/raffles/components/RaffleCard';
import { Search, Filter, Sparkles, TrendingUp, Zap } from 'lucide-react';
import type { RaffleFilters } from '@/types/raffle';

export const ExplorePage = () => {
  const user = useUser();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('Todos');

  // Fetch active raffles
  const [filters] = useState<RaffleFilters>({
    page: 1,
    page_size: 12,
    status: 'active',
  });

  const { data, isLoading, error } = useRafflesList(filters);

  const greeting = user ? (
    new Date().getHours() < 12
      ? `Buenos días, ${user.first_name || 'Usuario'}`
      : new Date().getHours() < 19
      ? `Buenas tardes, ${user.first_name || 'Usuario'}`
      : `Buenas noches, ${user.first_name || 'Usuario'}`
  ) : 'Bienvenido';

  const categories = [
    { id: 'Todos', icon: '🎯', label: 'Todos', count: data?.pagination.total || 0 },
    { id: 'Electrónica', icon: '📱', label: 'Electrónica', count: 0 },
    { id: 'Vehículos', icon: '🏍️', label: 'Vehículos', count: 0 },
    { id: 'Hogar', icon: '🏠', label: 'Hogar', count: 0 },
    { id: 'Otros', icon: '🎁', label: 'Otros', count: 0 },
  ];

  // Calculate stats from data
  const activeCount = data?.pagination.total || 0;
  const endingToday = 0; // TODO: Calculate from draw_date
  const newToday = 0; // TODO: Calculate from created_at

  // Handle error state
  if (error) {
    return (
      <div className="space-y-6">
        <div className="animate-fade-in">
          <h1 className="text-4xl font-bold text-slate-900 dark:text-white">
            {greeting} 👋
          </h1>
          <p className="text-lg text-slate-600 dark:text-slate-400 mt-2">
            Descubre sorteos verificables y participa para ganar increíbles premios
          </p>
        </div>
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <EmptyState
            icon={<Sparkles className="w-12 h-12" />}
            title="Error al cargar sorteos"
            description={error instanceof Error ? error.message : 'Ocurrió un error al cargar los sorteos. Por favor intenta de nuevo.'}
          />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="animate-fade-in">
        <h1 className="text-4xl font-bold text-slate-900 dark:text-white">
          {greeting} 👋
        </h1>
        <p className="text-lg text-slate-600 dark:text-slate-400 mt-2">
          Descubre sorteos verificables y participa para ganar increíbles premios
        </p>
      </div>

      {/* Stats Banner */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-gradient-to-br from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20 rounded-lg border border-primary-200 dark:border-primary-800 p-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center">
              <TrendingUp className="w-5 h-5 text-white" />
            </div>
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400">Sorteos Activos</p>
              <p className="text-2xl font-bold text-slate-900 dark:text-white">
                {isLoading ? '...' : activeCount}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-gradient-to-br from-success-50 to-success-100 dark:from-success-900/20 dark:to-success-800/20 rounded-lg border border-success-200 dark:border-success-800 p-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-success-600 rounded-lg flex items-center justify-center">
              <Zap className="w-5 h-5 text-white" />
            </div>
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400">Finalizan Hoy</p>
              <p className="text-2xl font-bold text-slate-900 dark:text-white">
                {isLoading ? '...' : endingToday}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-gradient-to-br from-warning-50 to-warning-100 dark:from-warning-900/20 dark:to-warning-800/20 rounded-lg border border-warning-200 dark:border-warning-800 p-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-warning-600 rounded-lg flex items-center justify-center">
              <span className="text-xl">🏆</span>
            </div>
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400">Nuevos Hoy</p>
              <p className="text-2xl font-bold text-slate-900 dark:text-white">
                {isLoading ? '...' : newToday}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Search and Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 w-5 h-5" />
          <input
            type="text"
            placeholder="Buscar sorteos..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-3 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
        </div>
        <Button variant="outline" className="sm:w-auto">
          <Filter className="w-4 h-4 mr-2" />
          Filtros
        </Button>
      </div>

      {/* Categories */}
      <div className="flex gap-2 overflow-x-auto pb-2">
        {categories.map((category) => (
          <button
            key={category.id}
            onClick={() => setSelectedCategory(category.id)}
            className={`px-4 py-2 rounded-lg whitespace-nowrap transition-all flex items-center gap-2 ${
              selectedCategory === category.id
                ? 'bg-primary-500 text-white shadow-lg shadow-primary-500/30'
                : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 border border-slate-300 dark:border-slate-600 hover:border-primary-500 dark:hover:border-primary-500 hover:shadow-md'
            }`}
          >
            <span>{category.icon}</span>
            <span>{category.label}</span>
            {category.count > 0 && (
              <span className={`ml-1 px-2 py-0.5 text-xs rounded-full ${
                selectedCategory === category.id
                  ? 'bg-white/20'
                  : 'bg-slate-100 dark:bg-slate-700'
              }`}>
                {category.count}
              </span>
            )}
          </button>
        ))}
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="py-12">
          <LoadingSpinner text="Cargando sorteos activos..." />
        </div>
      )}

      {/* Raffles Grid or Empty State */}
      {!isLoading && data && (
        <>
          {data.raffles.length === 0 ? (
            <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
              <EmptyState
                icon={<Sparkles className="w-12 h-12" />}
                title="No hay sorteos activos aún"
                description="Sé el primero en crear un sorteo o vuelve pronto para ver nuevas oportunidades."
                action={{
                  label: "Crear sorteo",
                  onClick: () => navigate('/organizer')
                }}
              />
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data.raffles.map((raffle) => (
                <RaffleCard key={raffle.id} raffle={raffle} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
};
