import { useParams, useNavigate, Link } from 'react-router-dom';
import { useEffect, useState, useCallback } from 'react';
import { useRaffleDetail, usePublishRaffle, useDeleteRaffle } from '../../../hooks/useRaffles';
import { useCreateReservation } from '../../../hooks/useReservations';
import { useAuth } from '../../../hooks/useAuth';
import { NumberGrid } from '../components/NumberGrid';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { FloatingCheckoutButton } from '../../../components/ui/FloatingCheckoutButton';
import { useCartStore } from '../../../store/cartStore';
import { toast } from 'sonner';
import {
  formatCurrency,
  formatDateTime,
  getStatusLabel,
  getDrawMethodLabel,
} from '../../../lib/utils';

const statusColors = {
  draft: 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300',
  active: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  suspended: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
  completed: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
  cancelled: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
};

export function RaffleDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();

  const { data, isLoading, error } = useRaffleDetail(id!, {
    includeNumbers: true,
    includeImages: true,
  });

  const publishMutation = usePublishRaffle();
  const deleteMutation = useDeleteRaffle();
  const createReservation = useCreateReservation();

  // Cart store integration
  const {
    setCurrentRaffle,
    toggleNumber,
    selectedNumbers,
    getSelectedCount,
    getTotalAmount,
    clearNumbers,
  } = useCartStore();

  // Reservation state
  const [currentReservation, setCurrentReservation] = useState<{
    id: string;
    expires_at: string;
  } | null>(null);
  const [sessionId] = useState(() => `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`);

  const isOwner = user && data?.raffle && user.id === data.raffle.user_id;
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';

  // Set current raffle when component mounts or id changes
  useEffect(() => {
    if (id && data?.raffle.uuid) {
      setCurrentRaffle(data.raffle.uuid);
    }
  }, [id, data?.raffle.uuid, setCurrentRaffle]);

  // Auto-create reservation when numbers are selected
  const createOrUpdateReservation = useCallback(async () => {
    if (!data?.raffle.uuid || getSelectedCount() === 0 || isOwner) return;

    // Don't create reservation if user is not authenticated
    if (!user) {
      console.log('User not authenticated, skipping auto-reservation');
      return;
    }

    // Don't create reservation if already processing
    if (createReservation.isPending) return;

    try {
      const response = await createReservation.mutateAsync({
        raffle_id: data.raffle.uuid,
        number_ids: selectedNumbers.map(n => n.id),
        session_id: sessionId,
      });

      setCurrentReservation({
        id: response.reservation.id,
        expires_at: response.reservation.expires_at,
      });

      toast.success(`${getSelectedCount()} número(s) reservado(s) por 15 minutos`, {
        description: 'Completa tu compra antes de que expire la reserva',
      });
    } catch (error) {
      console.error('Error creating reservation:', error);

      // Check for specific error types
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as { response?: { status?: number; data?: { message?: string } } };

        // 403: Email not verified
        if (axiosError.response?.status === 403) {
          toast.error('Email no verificado', {
            description: 'Debes verificar tu email para reservar números. Revisa tu bandeja de entrada.',
            duration: 5000,
          });
          return;
        }

        // 409: Numbers already reserved (possibly by another user or stale reservation)
        if (axiosError.response?.status === 409) {
          toast.warning('Números no disponibles', {
            description: 'Uno o más números ya están reservados. Por favor refresca la página.',
            duration: 4000,
          });
          // Don't clear automatically - let user see what they selected
          return;
        }

        // 429: Rate limit exceeded (too many requests)
        if (axiosError.response?.status === 429) {
          toast.info('Demasiadas solicitudes', {
            description: 'Por favor espera un momento antes de intentar de nuevo.',
            duration: 3000,
          });
          return;
        }
      }

      toast.error('Error al reservar números', {
        description: error instanceof Error ? error.message : 'Intenta de nuevo',
      });
      // Don't clear numbers on error - just let user try again
    }
  }, [data?.raffle.uuid, selectedNumbers, sessionId, getSelectedCount, isOwner, user, createReservation]);

  // DISABLED: Auto-reservation causes issues when user selects multiple numbers
  // Instead, reservation is created only when user clicks "Proceder al Pago"
  // useEffect(() => {
  //   if (getSelectedCount() === 0) {
  //     setCurrentReservation(null);
  //     return;
  //   }

  //   // Don't create if already have a reservation
  //   if (currentReservation) return;

  //   const timer = setTimeout(() => {
  //     createOrUpdateReservation();
  //   }, 1500); // Wait 1.5 seconds after last selection

  //   return () => clearTimeout(timer);
  // }, [selectedNumbers, currentReservation, createOrUpdateReservation, getSelectedCount]);

  const handlePublish = async () => {
    if (!id || !confirm('¿Estás seguro de publicar este sorteo?')) return;

    try {
      await publishMutation.mutateAsync(Number(id));
      alert('Sorteo publicado exitosamente');
    } catch (error) {
      alert(error instanceof Error ? error.message : 'Error al publicar sorteo');
    }
  };

  const handleDelete = async () => {
    if (!id || !confirm('¿Estás seguro de eliminar este sorteo? Esta acción no se puede deshacer.'))
      return;

    try {
      await deleteMutation.mutateAsync(Number(id));
      alert('Sorteo eliminado exitosamente');
      navigate('/raffles');
    } catch (error) {
      alert(error instanceof Error ? error.message : 'Error al eliminar sorteo');
    }
  };

  const handleNumberSelect = (numberStr: string) => {
    toggleNumber({
      id: numberStr,
      displayNumber: numberStr,
    });
  };

  const handleProceedToCheckout = async () => {
    if (getSelectedCount() === 0) {
      toast.error('Por favor selecciona al menos un número');
      return;
    }

    // Check if user is authenticated
    if (!user) {
      toast.info('Inicia sesión para continuar con tu compra', {
        description: 'Te redirigiremos al login',
      });
      // Redirect to login with return URL
      navigate(`/login?redirect=/raffles/${id}`);
      return;
    }

    // Create reservation before going to checkout
    if (!currentReservation) {
      await createOrUpdateReservation();
    }

    navigate('/checkout');
  };

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400 font-medium mb-2">
          Error al cargar el sorteo
        </p>
        <p className="text-sm text-slate-600 dark:text-slate-400 mb-4">
          {error instanceof Error ? error.message : 'Error desconocido'}
        </p>
        <Link to="/raffles">
          <Button variant="outline">Volver al listado</Button>
        </Link>
      </div>
    );
  }

  if (isLoading || !data) {
    return <LoadingSpinner text="Cargando sorteo..." />;
  }

  const { raffle, numbers = [], available_count, reserved_count, sold_count } = data;
  const soldPercentage = (sold_count / raffle.total_numbers) * 100;
  const daysUntilDraw = Math.ceil(
    (new Date(raffle.draw_date).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24)
  );

  return (
    <div className="space-y-8">
      {/* Back button */}
      <Link to="/raffles" className="inline-flex items-center text-blue-600 hover:text-blue-700 transition-colors">
        <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        Volver al listado
      </Link>

      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-700 dark:from-blue-700 dark:to-blue-800 rounded-xl overflow-hidden">
        <div className="p-8 md:p-12">
          <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-6">
            {/* Title and Status */}
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-4">
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${statusColors[raffle.status]}`}>
                  {getStatusLabel(raffle.status)}
                </span>
                {raffle.status === 'active' && daysUntilDraw > 0 && (
                  <span className="px-3 py-1 bg-white/20 text-white rounded-full text-sm font-medium backdrop-blur-sm">
                    {daysUntilDraw} {daysUntilDraw === 1 ? 'día' : 'días'} restantes
                  </span>
                )}
              </div>

              <h1 className="text-3xl md:text-4xl font-bold text-white mb-4">
                {raffle.title}
              </h1>

              <p className="text-blue-100 text-lg mb-6 max-w-2xl">
                {raffle.description}
              </p>

              {/* Price */}
              <div className="inline-flex flex-col bg-white/10 backdrop-blur-sm rounded-lg p-4 border border-white/20">
                <span className="text-blue-100 text-sm mb-1">Precio por número</span>
                <span className="text-3xl font-bold text-white">
                  {formatCurrency(Number(raffle.price_per_number))}
                </span>
              </div>
            </div>

            {/* CTA */}
            {raffle.status === 'active' && available_count > 0 && !isOwner && (
              <div className="flex-shrink-0">
                {getSelectedCount() > 0 ? (
                  <div className="space-y-3">
                    <div className="bg-white/10 backdrop-blur-sm rounded-lg p-4 border border-white/20">
                      <p className="text-blue-100 text-sm mb-1">Números seleccionados</p>
                      <p className="text-3xl font-bold text-white">{getSelectedCount()}</p>
                      <p className="text-blue-100 text-sm mt-2">
                        Total: {formatCurrency(getTotalAmount(Number(raffle.price_per_number)))}
                      </p>
                    </div>
                    <Button
                      size="lg"
                      onClick={handleProceedToCheckout}
                      className="bg-white text-blue-600 hover:bg-blue-50 shadow-lg w-full"
                    >
                      <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                      </svg>
                      Proceder al Pago
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={clearNumbers}
                      className="w-full bg-white/10 border-white/20 text-white hover:bg-white/20"
                    >
                      Limpiar selección
                    </Button>
                  </div>
                ) : (
                  <div className="text-center">
                    <p className="text-blue-100 text-sm mb-3">
                      Selecciona números en la grilla
                    </p>
                    <div className="flex items-center justify-center gap-2 text-blue-100 text-xs">
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
                      </svg>
                      <span>{available_count} números disponibles</span>
                    </div>
                  </div>
                )}
              </div>
            )}

            {/* Owner actions */}
            {isOwner && raffle.status === 'draft' && (
              <div className="flex flex-col gap-2">
                <Link to={`/raffles/${id}/edit`}>
                  <Button variant="outline" className="w-full bg-white/10 border-white/20 text-white hover:bg-white/20">
                    Editar
                  </Button>
                </Link>
                <Button
                  onClick={handlePublish}
                  disabled={publishMutation.isPending}
                  className="w-full bg-white text-blue-600 hover:bg-blue-50"
                >
                  Publicar
                </Button>
                {raffle.sold_count === 0 && (
                  <Button
                    variant="outline"
                    onClick={handleDelete}
                    disabled={deleteMutation.isPending}
                    className="w-full bg-red-600/10 border-red-400/20 text-red-100 hover:bg-red-600/20"
                  >
                    Eliminar
                  </Button>
                )}
              </div>
            )}

            {/* Admin actions (only for suspended raffles) */}
            {isAdmin && !isOwner && (raffle.status === 'draft' || raffle.status === 'suspended') && raffle.sold_count === 0 && (
              <div className="flex flex-col gap-2">
                <Button
                  variant="outline"
                  onClick={handleDelete}
                  disabled={deleteMutation.isPending}
                  className="w-full bg-red-600/10 border-red-400/20 text-red-100 hover:bg-red-600/20"
                >
                  Eliminar (Admin)
                </Button>
              </div>
            )}
          </div>
        </div>

        {/* Progress bar */}
        <div className="bg-white/10 backdrop-blur-sm px-8 md:px-12 py-4">
          <div className="flex items-center justify-between text-sm text-blue-100 mb-2">
            <span>Progreso de ventas</span>
            <span className="font-semibold">{soldPercentage.toFixed(1)}%</span>
          </div>
          <div className="w-full bg-white/20 rounded-full h-3">
            <div
              className="bg-white rounded-full h-3 transition-all duration-500"
              style={{ width: `${soldPercentage}%` }}
            />
          </div>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-slate-600 dark:text-slate-400">Disponibles</span>
            <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-3xl font-bold text-slate-900 dark:text-white">{available_count}</p>
        </div>

        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-slate-600 dark:text-slate-400">Vendidos</span>
            <svg className="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
            </svg>
          </div>
          <p className="text-3xl font-bold text-slate-900 dark:text-white">{sold_count}</p>
        </div>

        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-slate-600 dark:text-slate-400">Reservados</span>
            <svg className="w-5 h-5 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-3xl font-bold text-slate-900 dark:text-white">{reserved_count}</p>
        </div>

        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-slate-600 dark:text-slate-400">Recaudación</span>
            <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-3xl font-bold text-slate-900 dark:text-white">
            {formatCurrency(Number(raffle.total_revenue))}
          </p>
        </div>
      </div>

      {/* Raffle Info */}
      <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
        <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-6">
          Información del Sorteo
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Fecha del sorteo</p>
            <p className="font-medium text-slate-900 dark:text-white">
              {formatDateTime(raffle.draw_date)}
            </p>
          </div>
          <div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Método de sorteo</p>
            <p className="font-medium text-slate-900 dark:text-white">
              {getDrawMethodLabel(raffle.draw_method)}
            </p>
          </div>
          <div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">Total de números</p>
            <p className="font-medium text-slate-900 dark:text-white">
              {raffle.total_numbers}
            </p>
          </div>
          <div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mb-1">UUID</p>
            <p className="font-mono text-xs text-slate-600 dark:text-slate-400">
              {raffle.uuid}
            </p>
          </div>
        </div>
      </div>

      {/* Numbers Grid */}
      {numbers.length > 0 && (
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-6">
            Números del Sorteo
          </h2>
          <NumberGrid
            numbers={numbers}
            selectedNumbers={selectedNumbers.map((n) => n.id)}
            onNumberSelect={handleNumberSelect}
            readonly={isOwner || raffle.status !== 'active'}
          />
        </div>
      )}

      {/* Floating Checkout Button */}
      {!isOwner && raffle.status === 'active' && (
        <FloatingCheckoutButton
          selectedCount={getSelectedCount()}
          totalAmount={getTotalAmount(Number(raffle.price_per_number))}
          expiresAt={currentReservation?.expires_at}
          onCheckout={handleProceedToCheckout}
          onClear={clearNumbers}
          show={getSelectedCount() > 0}
          isAuthenticated={!!user}
          isEmailVerified={user?.kyc_level !== 'none'}
        />
      )}
    </div>
  );
}
