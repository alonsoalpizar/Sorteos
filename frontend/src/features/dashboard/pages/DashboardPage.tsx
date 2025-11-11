import { useNavigate } from 'react-router-dom';
import { useUser } from '@/hooks/useAuth';
import { Button } from '@/components/ui/Button';
import { StatsCard } from '@/components/ui/StatsCard';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';

export const DashboardPage = () => {
  const user = useUser();
  const navigate = useNavigate();

  if (!user) {
    return <LoadingSpinner text="Cargando dashboard..." />;
  }

  // Mock data - in a real app, this would come from API calls
  const stats = {
    activeRaffles: 0,
    totalSales: 0,
    pendingPurchases: 0,
    totalParticipations: 0,
  };

  const fullName = [user.first_name, user.last_name].filter(Boolean).join(' ');

  return (
    <div className="space-y-8">
      {/* Welcome Section */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-white">
          Bienvenido, {user.first_name || 'Usuario'}!
        </h1>
        <p className="text-slate-600 dark:text-slate-400 mt-2">
          Gestiona tus sorteos y participaciones desde tu panel de control
        </p>
      </div>

      {/* Quick Actions */}
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
        <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-4">
          Acciones Rápidas
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <Button
            onClick={() => navigate('/raffles/create')}
            className="h-auto py-4 flex flex-col items-center gap-2"
          >
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            <span className="font-semibold">Crear Sorteo</span>
            <span className="text-xs opacity-90">Organiza un nuevo sorteo</span>
          </Button>

          <Button
            variant="outline"
            onClick={() => navigate('/raffles')}
            className="h-auto py-4 flex flex-col items-center gap-2"
          >
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <span className="font-semibold">Explorar Sorteos</span>
            <span className="text-xs opacity-70">Encuentra sorteos activos</span>
          </Button>

          <Button
            variant="outline"
            onClick={() => navigate('/my-raffles')}
            className="h-auto py-4 flex flex-col items-center gap-2"
          >
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
            </svg>
            <span className="font-semibold">Mis Sorteos</span>
            <span className="text-xs opacity-70">Gestiona tus sorteos</span>
          </Button>
        </div>
      </div>

      {/* Stats Overview */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-4">
          Resumen
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatsCard
            title="Sorteos Activos"
            value={stats.activeRaffles}
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
              </svg>
            }
            description="Sorteos que has creado"
          />

          <StatsCard
            title="Ventas Totales"
            value={`$${stats.totalSales.toLocaleString()}`}
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            }
            description="De tus sorteos"
          />

          <StatsCard
            title="Compras Pendientes"
            value={stats.pendingPurchases}
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            }
            description="Números reservados"
          />

          <StatsCard
            title="Participaciones"
            value={stats.totalParticipations}
            icon={
              <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
              </svg>
            }
            description="Números comprados"
          />
        </div>
      </div>

      {/* Recent Activity */}
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
        <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-4">
          Actividad Reciente
        </h2>
        <div className="text-center py-8">
          <div className="inline-flex items-center justify-center w-12 h-12 bg-slate-100 dark:bg-slate-700 rounded-full mb-3">
            <svg className="w-6 h-6 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          <p className="text-slate-600 dark:text-slate-400">
            No hay actividad reciente
          </p>
          <p className="text-sm text-slate-500 dark:text-slate-500 mt-1">
            Crea tu primer sorteo o participa en uno existente para ver actividad aquí
          </p>
        </div>
      </div>

      {/* Account Info */}
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
        <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-4">
          Información de la Cuenta
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Nombre completo</p>
            <p className="font-medium text-slate-900 dark:text-white">
              {fullName || 'No especificado'}
            </p>
          </div>

          <div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Email</p>
            <p className="font-medium text-slate-900 dark:text-white">{user.email}</p>
          </div>

          {user.phone && (
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Teléfono</p>
              <p className="font-medium text-slate-900 dark:text-white">{user.phone}</p>
            </div>
          )}

          <div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Estado de verificación</p>
            <div className="flex items-center gap-2">
              {user.email_verified ? (
                <>
                  <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                  <span className="text-sm font-medium text-green-600 dark:text-green-400">
                    Email verificado
                  </span>
                </>
              ) : (
                <>
                  <span className="w-2 h-2 bg-yellow-500 rounded-full"></span>
                  <span className="text-sm font-medium text-yellow-600 dark:text-yellow-400">
                    Pendiente de verificación
                  </span>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
