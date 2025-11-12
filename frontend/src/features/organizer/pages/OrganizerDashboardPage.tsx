import { useNavigate } from 'react-router-dom';
import { useUser } from '@/hooks/useAuth';
import { GradientButton } from '@/components/ui/GradientButton';
import { EmptyState } from '@/components/ui/EmptyState';
import { StatsCard } from '@/components/ui/StatsCard';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { Plus, TrendingUp, Users, DollarSign, Package } from 'lucide-react';

export const OrganizerDashboardPage = () => {
  const user = useUser();
  const navigate = useNavigate();

  if (!user) {
    return <LoadingSpinner text="Cargando panel..." />;
  }

  // Mock data - in a real app, this would come from API calls
  const stats = {
    activeRaffles: 0,
    totalSales: 0,
    totalParticipants: 0,
    completedRaffles: 0,
  };

  const hasRaffles = stats.activeRaffles > 0 || stats.completedRaffles > 0;

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
            title="Ventas del Mes"
            value={`₡${stats.totalSales.toLocaleString()}`}
            icon={<DollarSign className="w-6 h-6" />}
            description="Ingresos generados"
          />

          <StatsCard
            title="Participantes"
            value={stats.totalParticipants}
            icon={<Users className="w-6 h-6" />}
            description="Compradores únicos"
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
              Sorteos Recientes
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
          {/* TODO: Add raffle list here */}
          <p className="text-slate-600 dark:text-slate-400">
            Lista de sorteos próximamente...
          </p>
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
