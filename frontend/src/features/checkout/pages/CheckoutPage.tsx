import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCartStore } from '../../../store/cartStore';
import { useAuth } from '../../../hooks/useAuth';
import { useRaffleDetail } from '../../../hooks/useRaffles';
import { useCreateReservation } from '../../../hooks/useReservations';
import { useCreatePaymentIntent } from '../../../hooks/usePayments';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { ReservationTimer } from '../../../components/ReservationTimer';
import { formatCurrency } from '../../../lib/utils';

type CheckoutStep = 'review' | 'reserving' | 'reserved' | 'creating_payment' | 'payment_ready' | 'expired';

export function CheckoutPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [step, setStep] = useState<CheckoutStep>('review');

  const {
    currentRaffleId,
    selectedNumbers,
    getSelectedCount,
    getTotalAmount,
    setReservation,
    activeReservation,
    clearNumbers,
    clearReservation,
  } = useCartStore();

  // Fetch raffle details
  const { data: raffleData } = useRaffleDetail(currentRaffleId || '', {
    includeNumbers: false,
    includeImages: false,
  });

  const createReservationMutation = useCreateReservation();
  const createPaymentIntentMutation = useCreatePaymentIntent();

  // Redirect if no numbers selected
  useEffect(() => {
    if (!user) {
      navigate('/login?redirect=/checkout');
      return;
    }

    if (getSelectedCount() === 0 && !activeReservation) {
      navigate('/raffles');
    }
  }, [user, getSelectedCount, activeReservation, navigate]);

  const handleCreateReservation = async () => {
    if (!currentRaffleId || getSelectedCount() === 0) {
      return;
    }

    setStep('reserving');

    try {
      const sessionId = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

      const result = await createReservationMutation.mutateAsync({
        raffle_id: currentRaffleId,
        number_ids: selectedNumbers.map(n => n.id),
        session_id: sessionId,
      });

      // Save reservation to cart store
      setReservation(result.reservation);
      setStep('reserved');
    } catch (error) {
      console.error('Error creating reservation:', error);
      alert(error instanceof Error ? error.message : 'Error al crear la reserva');
      setStep('review');
    }
  };

  const handleCreatePaymentIntent = async () => {
    if (!activeReservation) {
      return;
    }

    setStep('creating_payment');

    try {
      const result = await createPaymentIntentMutation.mutateAsync({
        reservation_id: activeReservation.id,
        return_url: window.location.origin + '/payment/success',
        cancel_url: window.location.origin + '/checkout',
      });

      // Redirect to PayPal approval URL
      if (result.payment_intent.client_secret) {
        // This is the PayPal approval URL
        window.location.href = result.payment_intent.client_secret;
      }
    } catch (error) {
      console.error('Error creating payment intent:', error);
      alert(error instanceof Error ? error.message : 'Error al iniciar el pago');
      setStep('reserved');
    }
  };

  const handleReservationExpire = () => {
    setStep('expired');
    clearReservation();
    clearNumbers();

    setTimeout(() => {
      navigate('/raffles');
    }, 3000);
  };

  const handleCancel = () => {
    if (confirm('¿Estás seguro de cancelar tu compra? Los números seleccionados se perderán.')) {
      clearNumbers();
      clearReservation();
      navigate('/raffles');
    }
  };

  if (!raffleData || !user) {
    return <LoadingSpinner text="Cargando..." />;
  }

  const { raffle } = raffleData;
  const pricePerNumber = Number(raffle.price_per_number);
  const totalAmount = activeReservation
    ? Number(activeReservation.total_amount)
    : getTotalAmount(pricePerNumber);

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900 dark:text-white mb-2">
          Checkout
        </h1>
        <p className="text-slate-600 dark:text-slate-400">
          Completa tu compra de números del sorteo
        </p>
      </div>

      {/* Reservation Timer */}
      {activeReservation && step !== 'expired' && (
        <ReservationTimer
          expiresAt={activeReservation.expires_at}
          onExpire={handleReservationExpire}
        />
      )}

      {/* Expired Message */}
      {step === 'expired' && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6 text-center">
          <svg
            className="w-16 h-16 text-red-600 dark:text-red-400 mx-auto mb-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <h2 className="text-2xl font-bold text-red-900 dark:text-red-100 mb-2">
            Reserva Expirada
          </h2>
          <p className="text-red-700 dark:text-red-300 mb-4">
            Tu reserva ha expirado. Serás redirigido al listado de sorteos...
          </p>
        </div>
      )}

      {step !== 'expired' && (
        <>
          {/* Order Summary */}
          <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
            <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-4">
              Resumen del Pedido
            </h2>

            {/* Raffle Info */}
            <div className="mb-6 pb-6 border-b border-slate-200 dark:border-slate-700">
              <h3 className="font-medium text-slate-900 dark:text-white mb-2">
                {raffle.title}
              </h3>
              <p className="text-sm text-slate-600 dark:text-slate-400">
                {raffle.description}
              </p>
            </div>

            {/* Numbers */}
            <div className="mb-6">
              <div className="flex items-center justify-between mb-3">
                <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
                  Números seleccionados ({activeReservation?.number_ids.length || selectedNumbers.length})
                </span>
              </div>
              <div className="flex flex-wrap gap-2">
                {(activeReservation?.number_ids || selectedNumbers.map(n => n.id))
                  .sort((a: string, b: string) => Number(a) - Number(b))
                  .map((num: string) => (
                    <span
                      key={num}
                      className="px-3 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-200 rounded-full text-sm font-mono font-semibold"
                    >
                      {num}
                    </span>
                  ))}
              </div>
            </div>

            {/* Price Breakdown */}
            <div className="space-y-3">
              <div className="flex items-center justify-between text-sm">
                <span className="text-slate-600 dark:text-slate-400">
                  Precio por número
                </span>
                <span className="font-medium text-slate-900 dark:text-white">
                  {formatCurrency(pricePerNumber)}
                </span>
              </div>
              <div className="flex items-center justify-between text-sm">
                <span className="text-slate-600 dark:text-slate-400">
                  Cantidad
                </span>
                <span className="font-medium text-slate-900 dark:text-white">
                  {activeReservation?.number_ids.length || selectedNumbers.length}
                </span>
              </div>
              <div className="pt-3 border-t border-slate-200 dark:border-slate-700">
                <div className="flex items-center justify-between">
                  <span className="text-lg font-semibold text-slate-900 dark:text-white">
                    Total
                  </span>
                  <span className="text-2xl font-bold text-blue-600 dark:text-blue-400">
                    {formatCurrency(totalAmount)}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
            {step === 'review' && (
              <div className="space-y-4">
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Al hacer clic en "Crear Reserva", tus números serán reservados por 5 minutos
                  para que completes el pago.
                </p>
                <div className="flex gap-3">
                  <Button
                    onClick={handleCreateReservation}
                    disabled={createReservationMutation.isPending}
                    className="flex-1"
                    size="lg"
                  >
                    {createReservationMutation.isPending ? 'Creando reserva...' : 'Crear Reserva'}
                  </Button>
                  <Button
                    variant="outline"
                    onClick={handleCancel}
                    size="lg"
                  >
                    Cancelar
                  </Button>
                </div>
              </div>
            )}

            {step === 'reserving' && (
              <div className="text-center py-8">
                <LoadingSpinner text="Creando tu reserva..." />
              </div>
            )}

            {step === 'reserved' && (
              <div className="space-y-4">
                <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
                  <div className="flex items-center gap-3">
                    <svg
                      className="w-6 h-6 text-green-600 dark:text-green-400"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    <div>
                      <p className="font-semibold text-green-900 dark:text-green-100">
                        ¡Reserva creada exitosamente!
                      </p>
                      <p className="text-sm text-green-700 dark:text-green-300">
                        Tus números están reservados. Procede al pago para confirmar tu compra.
                      </p>
                    </div>
                  </div>
                </div>
                <Button
                  onClick={handleCreatePaymentIntent}
                  disabled={createPaymentIntentMutation.isPending}
                  className="w-full"
                  size="lg"
                >
                  {createPaymentIntentMutation.isPending ? 'Redirigiendo a PayPal...' : 'Proceder al Pago con PayPal'}
                </Button>
              </div>
            )}

            {step === 'creating_payment' && (
              <div className="text-center py-8">
                <LoadingSpinner text="Preparando el pago..." />
                <p className="text-sm text-slate-600 dark:text-slate-400 mt-4">
                  Serás redirigido a PayPal en un momento...
                </p>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
