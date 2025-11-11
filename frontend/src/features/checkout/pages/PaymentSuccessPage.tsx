import { useEffect } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { Button } from '../../../components/ui/Button';
import { useCartStore } from '../../../store/cartStore';

export function PaymentSuccessPage() {
  const [searchParams] = useSearchParams();
  const { clearNumbers, clearReservation } = useCartStore();

  const paymentId = searchParams.get('payment_id');
  const reservationId = searchParams.get('reservation_id');

  useEffect(() => {
    // Clear cart and reservation on successful payment
    clearNumbers();
    clearReservation();
  }, [clearNumbers, clearReservation]);

  return (
    <div className="max-w-2xl mx-auto">
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-8 text-center">
        {/* Success Icon with animation */}
        <div className="mb-6 relative">
          <div className="w-24 h-24 mx-auto bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center animate-bounce">
            <svg
              className="w-12 h-12 text-green-600 dark:text-green-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
          </div>
        </div>

        {/* Success Message */}
        <h1 className="text-3xl font-bold text-slate-900 dark:text-white mb-4">
          ¡Pago Exitoso!
        </h1>

        <p className="text-lg text-slate-600 dark:text-slate-400 mb-6">
          Tu pago ha sido procesado correctamente. Tus números han sido confirmados.
        </p>

        {/* Payment Details */}
        {(paymentId || reservationId) && (
          <div className="bg-slate-50 dark:bg-slate-900/50 rounded-lg p-4 mb-6 space-y-2">
            {reservationId && (
              <div className="text-sm">
                <span className="text-slate-600 dark:text-slate-400">ID de Reserva: </span>
                <span className="font-mono text-slate-900 dark:text-white">{reservationId}</span>
              </div>
            )}
            {paymentId && (
              <div className="text-sm">
                <span className="text-slate-600 dark:text-slate-400">ID de Pago: </span>
                <span className="font-mono text-slate-900 dark:text-white">{paymentId}</span>
              </div>
            )}
          </div>
        )}

        {/* Next Steps */}
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 mb-6">
          <p className="text-sm text-blue-900 dark:text-blue-100 font-medium mb-2">
            ¿Qué sigue?
          </p>
          <ul className="text-sm text-blue-800 dark:text-blue-200 text-left space-y-1">
            <li>✓ Recibirás un correo de confirmación</li>
            <li>✓ Puedes ver tus números en "Mis Compras"</li>
            <li>✓ Te notificaremos cuando se realice el sorteo</li>
          </ul>
        </div>

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-3 justify-center">
          <Link to="/dashboard/purchases">
            <Button className="w-full sm:w-auto">
              Ver Mis Compras
            </Button>
          </Link>
          <Link to="/raffles">
            <Button variant="outline" className="w-full sm:w-auto">
              Ver Más Sorteos
            </Button>
          </Link>
        </div>
      </div>

      {/* Celebration Effect */}
      <div className="mt-8 text-center">
        <p className="text-6xl animate-pulse">🎉</p>
      </div>
    </div>
  );
}
