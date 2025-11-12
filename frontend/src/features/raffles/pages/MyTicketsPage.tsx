import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUser } from '@/hooks/useAuth';
import { EmptyState } from '@/components/ui/EmptyState';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { Ticket, TrendingUp } from 'lucide-react';

export const MyTicketsPage = () => {
  const user = useUser();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'active' | 'finished' | 'won'>('active');

  if (!user) {
    return <LoadingSpinner text="Cargando participaciones..." />;
  }

  // Mock data - in a real app, this would come from API
  const tickets: any[] = [];

  const tabs = [
    { id: 'active' as const, label: 'Activos', count: 0 },
    { id: 'finished' as const, label: 'Finalizados', count: 0 },
    { id: 'won' as const, label: 'Ganados 🎉', count: 0 },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="animate-fade-in">
        <h1 className="text-4xl font-bold text-slate-900 dark:text-white">
          Mis Números
        </h1>
        <p className="text-lg text-slate-600 dark:text-slate-400 mt-2">
          Gestiona tus participaciones en sorteos
        </p>
      </div>

      {/* Stats Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Total Participaciones</p>
              <p className="text-3xl font-bold text-slate-900 dark:text-white">0</p>
            </div>
            <div className="w-12 h-12 bg-primary-100 dark:bg-primary-900/20 rounded-lg flex items-center justify-center">
              <Ticket className="w-6 h-6 text-primary-600 dark:text-primary-400" />
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Invertido Total</p>
              <p className="text-3xl font-bold text-slate-900 dark:text-white">₡0</p>
            </div>
            <div className="w-12 h-12 bg-success-100 dark:bg-success-900/20 rounded-lg flex items-center justify-center">
              <TrendingUp className="w-6 h-6 text-success-600 dark:text-success-400" />
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Sorteos Ganados</p>
              <p className="text-3xl font-bold text-slate-900 dark:text-white">0</p>
            </div>
            <div className="w-12 h-12 bg-warning-100 dark:bg-warning-900/20 rounded-lg flex items-center justify-center">
              <span className="text-2xl">🏆</span>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 overflow-hidden">
        <div className="border-b border-slate-200 dark:border-slate-700">
          <nav className="flex gap-0">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex-1 px-6 py-4 text-sm font-medium transition-colors relative ${
                  activeTab === tab.id
                    ? 'text-primary-600 dark:text-primary-400'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
                }`}
              >
                {tab.label}
                {tab.count > 0 && (
                  <span className="ml-2 px-2 py-0.5 text-xs rounded-full bg-slate-100 dark:bg-slate-700">
                    {tab.count}
                  </span>
                )}
                {activeTab === tab.id && (
                  <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary-600 dark:bg-primary-400" />
                )}
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <div className="p-6">
          {tickets.length === 0 ? (
            <EmptyState
              icon={<Ticket className="w-12 h-12" />}
              title={
                activeTab === 'active'
                  ? '¡Empieza a participar!'
                  : activeTab === 'finished'
                  ? 'No hay sorteos finalizados'
                  : '¡Aún no has ganado!'
              }
              description={
                activeTab === 'active'
                  ? 'Explora los sorteos activos y compra números para participar. ¡La suerte te espera!'
                  : activeTab === 'finished'
                  ? 'Los sorteos en los que has participado y ya finalizaron aparecerán aquí.'
                  : 'Sigue participando en sorteos. ¡Tu momento llegará!'
              }
              action={{
                label: 'Explorar sorteos',
                onClick: () => navigate('/explore'),
              }}
            />
          ) : (
            <div className="space-y-4">
              {/* TODO: Map tickets here */}
              <p className="text-slate-600 dark:text-slate-400">Tickets próximamente...</p>
            </div>
          )}
        </div>
      </div>

      {/* Help Section */}
      <div className="bg-gradient-to-br from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20 rounded-lg border border-primary-200 dark:border-primary-800 p-6">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white mb-2">
          ¿Cómo funcionan los sorteos?
        </h3>
        <p className="text-slate-600 dark:text-slate-400 mb-4">
          Todos nuestros sorteos están basados en Lotería Nacional de Costa Rica, garantizando total transparencia y verificabilidad.
        </p>
        <div className="flex gap-4 text-sm">
          <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
            <span className="w-2 h-2 bg-success-500 rounded-full"></span>
            100% Verificable
          </div>
          <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
            <span className="w-2 h-2 bg-success-500 rounded-full"></span>
            Transparente
          </div>
          <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
            <span className="w-2 h-2 bg-success-500 rounded-full"></span>
            Seguro
          </div>
        </div>
      </div>
    </div>
  );
};
