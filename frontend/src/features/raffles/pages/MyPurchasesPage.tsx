import { useNavigate } from 'react-router-dom';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';

// Placeholder interface - will be replaced with real data from backend
interface Purchase {
  id: number;
  raffle_title: string;
  raffle_uuid: string;
  numbers: string[];
  total_amount: string;
  purchase_date: string;
  raffle_status: 'active' | 'completed' | 'cancelled';
  draw_date: string;
}

export function MyPurchasesPage() {
  const navigate = useNavigate();

  // Mock data - will be replaced with real API call
  const purchases: Purchase[] = [];
  const isLoading = false;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="w-8 h-8 border-4 border-slate-200 dark:border-slate-700 border-t-blue-600 rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-white">
          Mis Compras
        </h1>
        <p className="text-slate-600 dark:text-slate-400 mt-2">
          Historial de números que has comprado
        </p>
      </div>

      {/* Purchases List */}
      {purchases.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-12 h-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          }
          title="No has comprado números aún"
          description="Explora los sorteos activos y compra números para participar"
          action={{
            label: 'Explorar Sorteos',
            onClick: () => navigate('/raffles'),
          }}
        />
      ) : (
        <div className="space-y-4">
          {purchases.map((purchase) => {
            const isPending = purchase.raffle_status === 'active';
            const isCompleted = purchase.raffle_status === 'completed';

            return (
              <div
                key={purchase.id}
                className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6 hover:shadow-md transition-shadow cursor-pointer"
                onClick={() => navigate(`/raffles/${purchase.raffle_uuid}`)}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <h3 className="text-lg font-semibold text-slate-900 dark:text-white">
                        {purchase.raffle_title}
                      </h3>
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          isPending
                            ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
                            : isCompleted
                            ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                            : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                        }`}
                      >
                        {isPending ? 'Pendiente' : isCompleted ? 'Completado' : 'Cancelado'}
                      </span>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
                      <div>
                        <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">
                          Números comprados
                        </p>
                        <div className="flex flex-wrap gap-2">
                          {purchase.numbers.map((number) => (
                            <span
                              key={number}
                              className="inline-flex items-center justify-center w-10 h-10 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 font-semibold rounded-lg text-sm"
                            >
                              {number}
                            </span>
                          ))}
                        </div>
                      </div>

                      <div>
                        <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">
                          Monto total
                        </p>
                        <p className="text-xl font-bold text-slate-900 dark:text-white">
                          ₡{parseFloat(purchase.total_amount).toLocaleString()}
                        </p>
                        <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                          Comprado el {new Date(purchase.purchase_date).toLocaleDateString()}
                        </p>
                      </div>

                      <div>
                        <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">
                          Fecha del sorteo
                        </p>
                        <p className="font-medium text-slate-900 dark:text-white">
                          {new Date(purchase.draw_date).toLocaleDateString()}
                        </p>
                        {isPending && (
                          <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                            {Math.ceil(
                              (new Date(purchase.draw_date).getTime() - new Date().getTime()) /
                                (1000 * 60 * 60 * 24)
                            )}{' '}
                            días restantes
                          </p>
                        )}
                      </div>
                    </div>
                  </div>

                  <Button
                    size="sm"
                    variant="outline"
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/raffles/${purchase.raffle_uuid}`);
                    }}
                  >
                    Ver Sorteo
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Stats Summary */}
      {purchases.length > 0 && (
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <h2 className="text-lg font-semibold text-slate-900 dark:text-white mb-4">
            Resumen de Compras
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">
                Total invertido
              </p>
              <p className="text-2xl font-bold text-slate-900 dark:text-white">
                ₡
                {purchases
                  .reduce((sum, p) => sum + parseFloat(p.total_amount), 0)
                  .toLocaleString()}
              </p>
            </div>

            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">
                Números comprados
              </p>
              <p className="text-2xl font-bold text-slate-900 dark:text-white">
                {purchases.reduce((sum, p) => sum + p.numbers.length, 0)}
              </p>
            </div>

            <div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">
                Sorteos activos
              </p>
              <p className="text-2xl font-bold text-slate-900 dark:text-white">
                {purchases.filter((p) => p.raffle_status === 'active').length}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
