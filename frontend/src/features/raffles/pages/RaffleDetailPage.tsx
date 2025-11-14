import { useParams, useNavigate, Link } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { useRaffleDetail, usePublishRaffle, useDeleteRaffle } from '../../../hooks/useRaffles';
import { useAuth } from '../../../hooks/useAuth';
import { useRaffleWebSocket } from '../../../hooks/useRaffleWebSocket';
import { NumberGrid } from '../components/NumberGrid';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { FloatingCheckoutButton } from '../../../components/ui/FloatingCheckoutButton';
import { reservationService, Reservation } from '../../../services/reservationService';
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

  const { data, isLoading, error, refetch } = useRaffleDetail(id!, {
    includeNumbers: true,
    includeImages: true,
  });

  const publishMutation = usePublishRaffle();
  const deleteMutation = useDeleteRaffle();

  // Reservation state (única fuente de verdad)
  const [activeReservation, setActiveReservation] = useState<Reservation | null>(null);
  const [selectedNumbers, setSelectedNumbers] = useState<string[]>([]);
  const [isLoadingReservation, setIsLoadingReservation] = useState(false);
  const [sessionId] = useState(() => `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`);

  const isOwner = user && data?.raffle && user.id === data.raffle.user_id;
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';

  // WebSocket connection for real-time updates
  const { isConnected, onNumberUpdate, onReservationExpired } = useRaffleWebSocket(data?.raffle?.uuid);

  // Listen for WebSocket number updates (available, reserved, sold)
  useEffect(() => {
    if (!isConnected) return;

    const unsubscribeNumberUpdate = onNumberUpdate((update) => {
      console.log('[WebSocket] Number update:', update);

      // Solo refrescar si la actualización NO es del usuario actual
      // Si es del usuario actual, el estado ya se actualizó localmente
      const isMyUpdate = update.user_id === user?.id;

      if (!isMyUpdate) {
        // Refrescar los datos del sorteo para actualizar la grilla
        refetch();

        // Mostrar notificación según el estado
        if (update.status === 'sold') {
          toast.info(`Número ${update.number_id} vendido`);
        } else if (update.status === 'reserved') {
          toast.info(`Número ${update.number_id} reservado`);
        } else if (update.status === 'available') {
          toast.info(`Número ${update.number_id} disponible`);
        }
      }
    });

    const unsubscribeExpired = onReservationExpired((data) => {
      console.log('[WebSocket] Reservation expired:', data);

      // Refrescar datos
      refetch();
    });

    return () => {
      unsubscribeNumberUpdate();
      unsubscribeExpired();
    };
  }, [isConnected, onNumberUpdate, onReservationExpired, refetch, user?.id]);

  // Al montar componente: Cargar reserva activa si existe
  useEffect(() => {
    const loadOrCleanup = async () => {
      if (!data || !user || isOwner) return;

      try {
        const prevReservation = await reservationService.getActiveForRaffle(data.raffle.uuid);

        if (prevReservation) {
          // Ya tiene reserva activa - cargarla en lugar de cancelarla
          setActiveReservation(prevReservation);
          setSelectedNumbers(prevReservation.number_ids);
          console.log('Reserva activa cargada:', prevReservation.id);
        }
      } catch (error) {
        console.error('Error al cargar reserva activa:', error);
      }
    };

    loadOrCleanup();

    // NO hacemos cleanup automático aquí
    // La reserva se cancela explícitamente con el botón "Limpiar selección"
    // o expira automáticamente después de 10 minutos
  }, [data, user, isOwner]);

  // Monitorear timeout de reserva (10 minutos)
  useEffect(() => {
    if (!activeReservation) return;

    const checkExpiration = () => {
      const expiresAt = new Date(activeReservation.expires_at);
      const now = new Date();
      const timeLeft = expiresAt.getTime() - now.getTime();

      // Si ya expiró
      if (timeLeft <= 0) {
        toast.error('Tu reserva ha expirado', {
          description: 'Los números han sido liberados',
        });
        setActiveReservation(null);
        setSelectedNumbers([]);
        return;
      }

      // Alerta 1 minuto antes de expirar
      if (timeLeft <= 60 * 1000 && timeLeft > 59 * 1000) {
        toast.warning('¡Queda 1 minuto!', {
          description: 'Tu reserva está por expirar. Completa tu compra ahora.',
          duration: 10000,
        });
      }

      // Alerta 30 segundos antes de expirar
      if (timeLeft <= 30 * 1000 && timeLeft > 29 * 1000) {
        toast.warning('¡30 segundos!', {
          description: 'Tu reserva expirará pronto',
          duration: 10000,
        });
      }
    };

    // Verificar cada segundo
    const interval = setInterval(checkExpiration, 1000);

    return () => clearInterval(interval);
  }, [activeReservation]);

  // Manejar selección de números
  const handleNumberSelect = async (numberStr: string) => {
    // No permitir selección si es owner o no está autenticado
    if (isOwner || !user) {
      if (!user) {
        toast.info('Inicia sesión para reservar números');
      }
      return;
    }

    if (isLoadingReservation) return;

    const isAlreadySelected = selectedNumbers.includes(numberStr);

    try {
      setIsLoadingReservation(true);

      if (isAlreadySelected) {
        // REMOVER número de reserva
        if (activeReservation) {
          // Si es el último número, cancelar toda la reserva
          if (selectedNumbers.length === 1) {
            await reservationService.cancel(activeReservation.id);
            setActiveReservation(null);
            setSelectedNumbers([]);
            toast.info('Reserva cancelada');
          } else {
            // TODO: Implementar endpoint para remover número específico
            // Por ahora, no permitimos remover números individuales
            toast.warning('Por ahora no puedes desseleccionar números individuales. Usa "Limpiar selección"');
            return;
          }
        } else {
          // Solo está en estado local
          setSelectedNumbers(prev => prev.filter(n => n !== numberStr));
        }
      } else {
        // AGREGAR número
        const isFirstNumber = selectedNumbers.length === 0;

        if (isFirstNumber) {
          // CREAR NUEVA RESERVA con primer número
          const reservation = await reservationService.create({
            raffle_id: data!.raffle.uuid,
            number_ids: [numberStr],
            session_id: sessionId,
          });

          console.log('Reserva creada:', reservation);
          setActiveReservation(reservation);
          setSelectedNumbers([numberStr]);

          toast.success('Número reservado', {
            description: 'Tienes 10 minutos para completar tu compra',
          });
        } else {
          // AGREGAR a reserva existente
          if (!activeReservation) {
            throw new Error('No hay reserva activa');
          }

          const updatedReservation = await reservationService.addNumber(
            activeReservation.id,
            numberStr
          );

          setActiveReservation(updatedReservation);
          setSelectedNumbers(prev => [...prev, numberStr]);
        }
      }
    } catch (error: any) {
      console.error('Error al manejar selección:', error);

      // Manejo de errores específicos
      if (error.response?.status === 403) {
        toast.error('Email no verificado', {
          description: 'Verifica tu email para poder reservar números',
        });
      } else if (error.response?.status === 409) {
        toast.error('Número no disponible', {
          description: 'Este número ya está reservado por otro usuario',
        });
      } else {
        toast.error('Error al reservar', {
          description: 'No se pudo reservar el número. Intenta de nuevo',
        });
      }
    } finally {
      setIsLoadingReservation(false);
    }
  };

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

  const handleClearSelection = async () => {
    if (activeReservation) {
      try {
        await reservationService.cancel(activeReservation.id);
        setActiveReservation(null);
        setSelectedNumbers([]);
        toast.info('Reserva cancelada');
      } catch (error) {
        console.error('Error al cancelar reserva:', error);
        toast.error('Error al cancelar reserva');
      }
    } else {
      setSelectedNumbers([]);
    }
  };

  const handlePayNow = async () => {
    console.log('handlePayNow - activeReservation:', activeReservation);
    console.log('handlePayNow - selectedNumbers:', selectedNumbers);

    if (!activeReservation) {
      toast.error('No tienes números reservados');
      return;
    }

    if (!user) {
      toast.info('Inicia sesión para continuar');
      navigate(`/login?redirect=/raffles/${id}`);
      return;
    }

    try {
      // TODO: Implementar pago desde wallet (deducir balance)

      // Confirmar reserva (marca como 'confirmed' y deja de contar timeout)
      await reservationService.confirm(activeReservation.id);

      // Limpiar estado local
      setActiveReservation(null);
      setSelectedNumbers([]);

      // Mostrar mensaje de éxito
      toast.success('¡Gracias por tu compra!', {
        description: `Has comprado ${selectedNumbers.length} número(s). Redirigiendo a tus tickets...`,
        duration: 3000,
      });

      // Redirigir a /my-tickets después de 2 segundos
      setTimeout(() => {
        navigate('/my-tickets');
      }, 2000);
    } catch (error) {
      console.error('Error al procesar pago:', error);
      toast.error('Error al procesar el pago');
    }
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
      <button
        onClick={() => navigate(-1)}
        className="inline-flex items-center text-blue-600 hover:text-blue-700 transition-colors"
      >
        <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        Volver al listado
      </button>

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
                {selectedNumbers.length > 0 ? (
                  <div className="space-y-3">
                    <div className="bg-white/10 backdrop-blur-sm rounded-lg p-4 border border-white/20">
                      <p className="text-blue-100 text-sm mb-1">Números reservados</p>
                      <p className="text-3xl font-bold text-white">{selectedNumbers.length}</p>
                      <p className="text-blue-100 text-sm mt-2">
                        Total: {formatCurrency(selectedNumbers.length * Number(raffle.price_per_number))}
                      </p>
                    </div>
                    <Button
                      size="lg"
                      onClick={handlePayNow}
                      disabled={isLoadingReservation}
                      className="bg-white text-blue-600 hover:bg-blue-50 shadow-lg w-full"
                    >
                      <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                      </svg>
                      Pagar Ahora
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={handleClearSelection}
                      disabled={isLoadingReservation}
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
            selectedNumbers={selectedNumbers}
            onNumberSelect={handleNumberSelect}
            readonly={isOwner || raffle.status !== 'active'}
          />
        </div>
      )}

      {/* Floating Checkout Button */}
      {!isOwner && raffle.status === 'active' && selectedNumbers.length > 0 && (
        <FloatingCheckoutButton
          selectedCount={selectedNumbers.length}
          totalAmount={selectedNumbers.length * Number(raffle.price_per_number)}
          onCheckout={handlePayNow}
          onCancel={handleClearSelection}
          disabled={!user || user?.kyc_level === 'none' || isLoadingReservation}
        />
      )}
    </div>
  );
}
