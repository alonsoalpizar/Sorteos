import { Link } from 'react-router-dom';
import { Button } from '../../../components/ui/Button';
import { useCartStore } from '../../../store/cartStore';

export function PaymentCancelPage() {
  const { activeReservation } = useCartStore();

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-8 text-center">
        {/* Cancel Icon */}
        <div className="mb-6">
          <div className="w-24 h-24 mx-auto bg-yellow-100 dark:bg-yellow-900/30 rounded-full flex items-center justify-center">
            <svg
              className="w-12 h-12 text-yellow-600 dark:text-yellow-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
          </div>
        </div>

        {/* Message */}
        <h1 className="text-3xl font-bold text-slate-900 dark:text-white mb-4">
          Pago Cancelado
        </h1>

        <p className="text-lg text-slate-600 dark:text-slate-400 mb-6">
          Has cancelado el proceso de pago.
          {activeReservation && ' Tu reserva sigue activa, puedes intentar pagar nuevamente.'}
        </p>

        {/* Reservation Status */}
        {activeReservation && (
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 mb-6">
            <p className="text-sm text-blue-900 dark:text-blue-100 font-medium mb-2">
              Tu reserva sigue activa
            </p>
            <p className="text-sm text-blue-800 dark:text-blue-200">
              Puedes volver al checkout para completar tu pago antes de que expire la reserva.
            </p>
          </div>
        )}

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-3 justify-center">
          {activeReservation ? (
            <>
              <Link to="/checkout">
                <Button className="w-full sm:w-auto">
                  Volver al Checkout
                </Button>
              </Link>
              <Link to="/raffles">
                <Button variant="outline" className="w-full sm:w-auto">
                  Ver Sorteos
                </Button>
              </Link>
            </>
          ) : (
            <>
              <Link to="/raffles">
                <Button className="w-full sm:w-auto">
                  Ver Sorteos
                </Button>
              </Link>
              <Link to="/dashboard">
                <Button variant="outline" className="w-full sm:w-auto">
                  Ir al Dashboard
                </Button>
              </Link>
            </>
          )}
        </div>
      </div>

      {/* Help Text */}
      <div className="mt-6 text-center">
        <p className="text-sm text-slate-600 dark:text-slate-400">
          ¿Necesitas ayuda?{' '}
          <Link to="/support" className="text-blue-600 dark:text-blue-400 hover:underline">
            Contáctanos
          </Link>
        </p>
      </div>
    </div>
  );
}
