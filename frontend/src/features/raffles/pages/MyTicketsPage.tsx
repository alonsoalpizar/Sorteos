import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useUser } from '@/hooks/useAuth';
import { useMyTickets } from '@/hooks/useRaffles';
import { EmptyState } from '@/components/ui/EmptyState';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';
import { Ticket, TrendingUp, Calendar, DollarSign } from 'lucide-react';
import { formatCurrency, formatDateTime } from '@/lib/utils';

export const MyTicketsPage = () => {
  const user = useUser();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'active' | 'finished' | 'won'>('active');
  const [page] = useState(1);

  const { data, isLoading, error } = useMyTickets(page, 20);

  if (!user) {
    return <LoadingSpinner text="Cargando participaciones..." />;
  }

  if (isLoading) {
    return <LoadingSpinner text="Cargando tus tickets..." />;
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">Error al cargar tus tickets</p>
      </div>
    );
  }

  const tickets = data?.tickets || [];

  // Filtrar tickets según el tab activo
  const filteredTickets = tickets.filter((ticket) => {
    const now = new Date();
    const drawDate = new Date(ticket.raffle.draw_date);

    if (activeTab === 'active') {
      return ticket.raffle.status === 'active' && drawDate > now;
    } else if (activeTab === 'finished') {
      return ticket.raffle.status === 'completed' || (ticket.raffle.status === 'active' && drawDate <= now);
    } else if (activeTab === 'won') {
      return ticket.raffle.winner_user_id === user.id;
    }
    return false;
  });

  // Calcular estadísticas
  const totalTickets = tickets.reduce((sum, t) => sum + t.total_numbers, 0);
  const totalSpent = tickets.reduce((sum, t) => sum + parseFloat(t.total_spent || '0'), 0);
  const totalWon = tickets.filter(t => t.raffle.winner_user_id === user.id).length;

  const tabs = [
    { id: 'active' as const, label: 'Activos', count: tickets.filter(t => t.raffle.status === 'active').length },
    { id: 'finished' as const, label: 'Finalizados', count: tickets.filter(t => t.raffle.status === 'completed').length },
    { id: 'won' as const, label: 'Ganados 🎉', count: totalWon },
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
              <p className="text-3xl font-bold text-slate-900 dark:text-white">{totalTickets}</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
              <Ticket className="w-6 h-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Invertido Total</p>
              <p className="text-3xl font-bold text-slate-900 dark:text-white">
                {formatCurrency(totalSpent)}
              </p>
            </div>
            <div className="w-12 h-12 bg-green-100 dark:bg-green-900/20 rounded-lg flex items-center justify-center">
              <TrendingUp className="w-6 h-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Sorteos Ganados</p>
              <p className="text-3xl font-bold text-slate-900 dark:text-white">{totalWon}</p>
            </div>
            <div className="w-12 h-12 bg-amber-100 dark:bg-amber-900/20 rounded-lg flex items-center justify-center">
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
                    ? 'text-blue-600 dark:text-blue-400'
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
                  <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-600 dark:bg-blue-400" />
                )}
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <div className="p-6">
          {filteredTickets.length === 0 ? (
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
              {filteredTickets.map((ticket) => (
                <div
                  key={ticket.raffle.id}
                  className="border border-slate-200 dark:border-slate-700 rounded-lg p-6 hover:border-blue-300 dark:hover:border-blue-700 transition-colors"
                >
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                      <Link
                        to={`/raffles/${ticket.raffle.id}`}
                        className="text-lg font-semibold text-slate-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                      >
                        {ticket.raffle.title}
                      </Link>
                      <p className="text-sm text-slate-600 dark:text-slate-400 mt-1">
                        {ticket.raffle.description}
                      </p>
                    </div>
                    <span
                      className={`px-3 py-1 rounded-full text-sm font-medium ${
                        ticket.raffle.status === 'active'
                          ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                          : ticket.raffle.status === 'completed'
                          ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
                          : 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300'
                      }`}
                    >
                      {ticket.raffle.status === 'active' ? 'Activo' : 'Finalizado'}
                    </span>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-400">
                      <Calendar className="w-4 h-4" />
                      <span>Sorteo: {formatDateTime(ticket.raffle.draw_date)}</span>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-400">
                      <Ticket className="w-4 h-4" />
                      <span>{ticket.total_numbers} número(s)</span>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-400">
                      <DollarSign className="w-4 h-4" />
                      <span>Invertido: {formatCurrency(parseFloat(ticket.total_spent))}</span>
                    </div>
                  </div>

                  {/* Números comprados */}
                  <div className="border-t border-slate-200 dark:border-slate-700 pt-4">
                    <p className="text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
                      Tus números:
                    </p>
                    <div className="flex flex-wrap gap-2">
                      {ticket.numbers.map((num) => (
                        <span
                          key={num.id}
                          className={`px-3 py-1 rounded-md text-sm font-mono font-semibold ${
                            ticket.raffle.winner_number === num.number
                              ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400 ring-2 ring-amber-400'
                              : 'bg-slate-100 text-slate-700 dark:bg-slate-700 dark:text-slate-300'
                          }`}
                        >
                          {num.number}
                          {ticket.raffle.winner_number === num.number && ' 🎉'}
                        </span>
                      ))}
                    </div>
                  </div>

                  {ticket.raffle.winner_user_id === user.id && (
                    <div className="mt-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-4">
                      <p className="text-amber-800 dark:text-amber-200 font-semibold">
                        🎉 ¡Felicidades! Has ganado este sorteo
                      </p>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Help Section */}
      <div className="bg-gradient-to-br from-blue-50 to-blue-100 dark:from-blue-900/20 dark:to-blue-800/20 rounded-lg border border-blue-200 dark:border-blue-800 p-6">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white mb-2">
          ¿Cómo funcionan los sorteos?
        </h3>
        <p className="text-slate-600 dark:text-slate-400 mb-4">
          Todos nuestros sorteos están basados en Lotería Nacional de Costa Rica, garantizando total transparencia y verificabilidad.
        </p>
        <div className="flex gap-4 text-sm">
          <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
            <span className="w-2 h-2 bg-green-500 rounded-full"></span>
            100% Verificable
          </div>
          <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
            <span className="w-2 h-2 bg-green-500 rounded-full"></span>
            Transparente
          </div>
          <div className="flex items-center gap-2 text-slate-700 dark:text-slate-300">
            <span className="w-2 h-2 bg-green-500 rounded-full"></span>
            Seguro
          </div>
        </div>
      </div>
    </div>
  );
};
