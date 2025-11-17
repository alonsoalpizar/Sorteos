import { useNavigate } from 'react-router-dom';
import { useUser } from '@/hooks/useAuth';
import { useRafflesList } from '@/hooks/useRaffles';
import { GradientButton } from '@/components/ui/GradientButton';
import { EmptyState } from '@/components/ui/EmptyState';
import { StatsCard } from '@/components/ui/StatsCard';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { Plus, TrendingUp, Users, DollarSign, Package, Calendar } from 'lucide-react';
import { formatCurrency, formatDateTime } from '@/lib/utils';

export const OrganizerDashboardPage = () => {
  const user = useUser();
  const navigate = useNavigate();

  // Obtener sorteos del usuario actual
  const { data: myRaffles, isLoading } = useRafflesList({
    user_id: user?.id,
    page_size: 100
  });

  if (!user) {
    return <LoadingSpinner text="Cargando panel..." />;
  }

  if (isLoading) {
    return <LoadingSpinner text="Cargando estadísticas..." />;
  }

  const raffles = myRaffles?.raffles || [];

  // Calcular estadísticas desde los sorteos del usuario
  const activeRaffles = raffles.filter(r => r.status === 'active').length;
  const completedRaffles = raffles.filter(r => r.status === 'completed').length;
  const draftRaffles = raffles.filter(r => r.status === 'draft').length;

  // Calcular ventas totales (solo de sorteos activos y completados)
  const totalSales = raffles
    .filter(r => r.status === 'active' || r.status === 'completed')
    .reduce((sum, r) => sum + (parseFloat(r.total_revenue || '0')), 0);

  const stats = {
    activeRaffles,
    totalSales,
    draftRaffles,
    completedRaffles,
  };

  const hasRaffles = raffles.length > 0;

  return (
    <div className="space-y-8">
      {/* Welcome Section */}
      <div className="animate-fade-in">
        <h1 className="text-4xl font-bold text-slate-900 dark:text-white">
          Panel de Organizador
        </h1>
        <p className="text-lg text-slate-600 dark:text-slate-400 mt-2">
          Gestiona tus sorteos y monitorea tus ventas
        </p>
      </div>

      {/* Stats Overview */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-4">
          Estadísticas
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatsCard
            title="Sorteos Activos"
            value={stats.activeRaffles}
            icon={<Package className="w-6 h-6" />}
            description="En venta actualmente"
          />

          <StatsCard
            title="Total Recaudado"
            value={formatCurrency(stats.totalSales)}
            icon={<DollarSign className="w-6 h-6" />}
            description="Ingresos generados"
          />

          <StatsCard
            title="Borradores"
            value={stats.draftRaffles}
            icon={<Users className="w-6 h-6" />}
            description="Sorteos sin publicar"
          />

          <StatsCard
            title="Completados"
            value={stats.completedRaffles}
            icon={<TrendingUp className="w-6 h-6" />}
            description="Sorteos finalizados"
          />
        </div>
      </div>

      {/* Quick Actions or Empty State */}
      {hasRaffles ? (
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold text-slate-900 dark:text-white">
              Mis Sorteos
            </h2>
            <GradientButton
              onClick={() => navigate('/organizer/raffles/new')}
              variant="primary"
              size="sm"
            >
              <Plus className="w-4 h-4 mr-2" />
              Crear Sorteo
            </GradientButton>
          </div>

          <div className="space-y-4">
            {raffles.slice(0, 5).map((raffle) => (
              <div
                key={raffle.id}
                onClick={() => navigate(`/raffles/${raffle.id}`)}
                className="flex items-center justify-between p-4 rounded-lg border border-slate-200 dark:border-slate-700 hover:border-blue-300 dark:hover:border-blue-700 transition-colors cursor-pointer"
              >
                <div className="flex-1">
                  <h3 className="font-semibold text-slate-900 dark:text-white">
                    {raffle.title}
                  </h3>
                  <div className="flex items-center gap-4 mt-1 text-sm text-slate-600 dark:text-slate-400">
                    <span className="flex items-center gap-1">
                      <Calendar className="w-4 h-4" />
                      {formatDateTime(raffle.draw_date)}
                    </span>
                    <span>
                      {raffle.sold_count}/{raffle.total_numbers} vendidos
                    </span>
                  </div>
                </div>

                <div className="flex items-center gap-4">
                  <div className="text-right">
                    <p className="text-sm text-slate-600 dark:text-slate-400">Recaudado</p>
                    <p className="font-semibold text-slate-900 dark:text-white">
                      {formatCurrency(parseFloat(raffle.total_revenue || '0'))}
                    </p>
                  </div>

                  <span
                    className={`px-3 py-1 rounded-full text-sm font-medium ${
                      raffle.status === 'active'
                        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                        : raffle.status === 'draft'
                        ? 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300'
                        : raffle.status === 'completed'
                        ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
                        : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
                    }`}
                  >
                    {raffle.status === 'active' ? 'Activo' :
                     raffle.status === 'draft' ? 'Borrador' :
                     raffle.status === 'completed' ? 'Completado' : 'Suspendido'}
                  </span>
                </div>
              </div>
            ))}
          </div>

          {raffles.length > 5 && (
            <div className="mt-4 text-center">
              <button
                onClick={() => navigate('/organizer/raffles')}
                className="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium"
              >
                Ver todos los sorteos ({raffles.length})
              </button>
            </div>
          )}
        </div>
      ) : (
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <EmptyState
            icon={<Package className="w-12 h-12" />}
            title="¡Comienza tu primer sorteo!"
            description="Crea sorteos verificables y transparentes basados en Lotería Nacional. Es fácil, rápido y seguro."
            action={{
              label: "Crear mi primer sorteo",
              onClick: () => navigate('/organizer/raffles/new')
            }}
          />
        </div>
      )}

      {/* Help Section */}
      <div className="bg-gradient-to-br from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20 rounded-lg border border-primary-200 dark:border-primary-800 p-6">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white mb-2">
          ¿Necesitas ayuda?
        </h3>
        <p className="text-slate-600 dark:text-slate-400 mb-4">
          Aprende cómo crear sorteos exitosos, configurar premios y gestionar participantes.
        </p>
        <GradientButton
          onClick={() => navigate('/help')}
          variant="primary"
          size="sm"
        >
          Ver guía de inicio
        </GradientButton>
      </div>
    </div>
  );
};
