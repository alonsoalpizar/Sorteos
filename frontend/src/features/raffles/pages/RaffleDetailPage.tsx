import { useParams, useNavigate, Link } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { useRaffleDetail, usePublishRaffle, useDeleteRaffle, useRaffleBuyers } from '../../../hooks/useRaffles';
import { useAuth } from '../../../hooks/useAuth';
import { useRaffleWebSocket } from '../../../hooks/useRaffleWebSocket';
import { useUserMode } from '../../../contexts/UserModeContext';
import { NumberGrid } from '../components/NumberGrid';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { FloatingCheckoutButton } from '../../../components/ui/FloatingCheckoutButton';
import { RaffleImageGallery } from '../../../components/RaffleImageGallery';
import { reservationService, Reservation } from '../../../services/reservationService';
import { toast } from 'sonner';
import {
  cn,
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
  const { colors } = useUserMode();

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

  // Cargar lista de compradores (solo para owner)
  const { data: buyersData, isLoading: isLoadingBuyers } = useRaffleBuyers(
    data?.raffle?.uuid || '',
    { enabled: !!isOwner && !!data?.raffle?.uuid }
  );

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

    if (isLoadingReservation) {
      return;
    }

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
            refetch(); // Refrescar datos del raffle
          } else {
            // Remover número específico
            const updatedReservation = await reservationService.removeNumber(
              activeReservation.id,
              numberStr
            );

            setActiveReservation(updatedReservation);
            setSelectedNumbers(prev => prev.filter(n => n !== numberStr));

            toast.success('Número liberado', {
              description: `Has des-reservado el número ${numberStr}`,
            });
            refetch(); // Refrescar datos del raffle
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

          setActiveReservation(reservation);
          setSelectedNumbers([numberStr]);
          await refetch();

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
      } else if (error.response?.status === 400 && error.response?.data?.code === 'CHECKOUT_PHASE') {
        toast.error('No puedes des-reservar en fase de pago', {
          description: 'Cancela la reserva completa o completa el pago',
        });
      } else if (error.response?.status === 400 && error.response?.data?.code === 'CANNOT_REMOVE_LAST') {
        toast.info('Último número', {
          description: 'No puedes remover el último número. La reserva se cancelará automáticamente',
        });
      } else {
        toast.error('Error al procesar', {
          description: 'No se pudo procesar la acción. Intenta de nuevo',
        });
      }
    } finally {
      setIsLoadingReservation(false);
    }
  };

  const handlePublish = async () => {
    if (!data?.raffle?.id) return;

    const confirmed = confirm(
      '⚠️ IMPORTANTE: Una vez publicado, el sorteo estará visible para todos los usuarios y NO podrás:\n\n' +
      '• Modificar el título\n' +
      '• Cambiar el precio por número\n' +
      '• Alterar la cantidad de números\n' +
      '• Eliminar el sorteo (solo suspender si hay problemas)\n\n' +
      'Solo podrás modificar la descripción y la fecha del sorteo.\n\n' +
      '¿Estás seguro de publicar este sorteo?'
    );

    if (!confirmed) return;

    try {
      await publishMutation.mutateAsync(data.raffle.id);
      toast.success('Sorteo publicado exitosamente', {
        description: 'Ahora es visible para todos los usuarios'
      });
      refetch();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error al publicar sorteo');
    }
  };

  const handleDelete = async () => {
    if (!data?.raffle?.id || !confirm('¿Estás seguro de eliminar este sorteo? Esta acción no se puede deshacer.'))
      return;

    try {
      await deleteMutation.mutateAsync(data.raffle.id);
      toast.success('Sorteo eliminado exitosamente');
      navigate('/raffles');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error al eliminar sorteo');
    }
  };

  const handleSuspend = async () => {
    if (!id) return;

    const reason = prompt('¿Por qué deseas suspender este sorteo?\n(Esta información será visible para los compradores)');
    if (!reason) return;

    try {
      // TODO: Implementar useSuspendRaffle hook y API
      toast.info('Funcionalidad de suspender sorteo en desarrollo');
      // await suspendMutation.mutateAsync({ id: Number(id), reason });
      // toast.success('Sorteo suspendido');
      // refetch();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error al suspender sorteo');
    }
  };

  const handleReactivate = async () => {
    if (!id || !confirm('¿Deseas reactivar este sorteo?')) return;

    try {
      // TODO: Implementar reactivar sorteo en el backend
      toast.info('Funcionalidad de reactivar sorteo en desarrollo');
      // await reactivateMutation.mutateAsync(Number(id));
      // toast.success('Sorteo reactivado');
      // refetch();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error al reactivar sorteo');
    }
  };

  const handleExtendDate = async () => {
    if (!id) return;

    const newDate = prompt('Ingresa la nueva fecha del sorteo (formato: YYYY-MM-DD HH:mm):');
    if (!newDate) return;

    try {
      // TODO: Implementar extender fecha
      toast.info('Funcionalidad de extender fecha en desarrollo');
      // const isoDate = new Date(newDate).toISOString();
      // await updateMutation.mutateAsync({ id: Number(id), input: { draw_date: isoDate } });
      // toast.success('Fecha extendida exitosamente');
      // refetch();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error al extender fecha');
    }
  };

  const handleCloseDraw = async () => {
    if (!id || !confirm('¿Estás seguro de cerrar este sorteo sin realizar el sorteo?\nEsta acción cancelará el sorteo y se devolverá el dinero a los compradores.')) return;

    try {
      // TODO: Implementar cerrar sorteo sin ganador
      toast.info('Funcionalidad de cerrar sorteo en desarrollo');
      // await closeMutation.mutateAsync(Number(id));
      // toast.success('Sorteo cerrado. Se procesarán las devoluciones.');
      // refetch();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Error al cerrar sorteo');
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

  const handleReserve = async () => {
    console.log('handleReserve - activeReservation:', activeReservation);
    console.log('handleReserve - selectedNumbers:', selectedNumbers);

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
      // Confirmar reserva (marca como 'confirmed' con timeout de 24 horas)
      await reservationService.confirm(activeReservation.id);

      // Limpiar estado local
      setActiveReservation(null);
      setSelectedNumbers([]);

      // Mostrar mensaje de éxito con instrucciones
      toast.success('¡Números reservados exitosamente!', {
        description: `Tienes 24 horas para coordinar el pago con el organizador. De lo contrario, los números se liberarán automáticamente.`,
        duration: 8000,
      });

      // Redirigir a /my-tickets después de 2 segundos
      setTimeout(() => {
        navigate('/my-tickets');
      }, 2000);
    } catch (error) {
      console.error('Error al confirmar reserva:', error);
      toast.error('Error al confirmar la reserva');
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

      {/* Hero Section - color dinámico según modo */}
      <div className={cn(
        "rounded-xl overflow-hidden",
        colors.gradient,
        colors.gradientDark
      )}>
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

              <p className={cn("text-lg mb-6 max-w-2xl", colors.textMuted)}>
                {raffle.description}
              </p>

              {/* Price */}
              <div className="inline-flex flex-col bg-white/10 backdrop-blur-sm rounded-lg p-4 border border-white/20">
                <span className={cn("text-sm mb-1", colors.textMuted)}>Precio por número</span>
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
                      <p className={cn("text-sm mb-1", colors.textMuted)}>Números reservados</p>
                      <p className="text-3xl font-bold text-white">{selectedNumbers.length}</p>
                      <p className={cn("text-sm mt-2", colors.textMuted)}>
                        Total: {formatCurrency(selectedNumbers.length * Number(raffle.price_per_number))}
                      </p>
                    </div>
                    <Button
                      size="lg"
                      onClick={handleReserve}
                      disabled={isLoadingReservation}
                      className="bg-white text-blue-600 hover:bg-blue-50 shadow-lg w-full"
                    >
                      <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      Reservar
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
                    <p className={cn("text-sm mb-3", colors.textMuted)}>
                      Selecciona números en la grilla
                    </p>
                    <div className={cn("flex items-center justify-center gap-2 text-xs", colors.textMuted)}>
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
            {isOwner && (
              <div className="flex flex-col gap-2">
                {/* Draft: Editar, Publicar, Eliminar */}
                {raffle.status === 'draft' && (
                  <>
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
                      {publishMutation.isPending ? 'Publicando...' : 'Publicar'}
                    </Button>
                    {raffle.sold_count === 0 && (
                      <Button
                        variant="outline"
                        onClick={handleDelete}
                        disabled={deleteMutation.isPending}
                        className="w-full bg-red-600/10 border-red-400/20 text-red-100 hover:bg-red-600/20"
                      >
                        {deleteMutation.isPending ? 'Eliminando...' : 'Eliminar'}
                      </Button>
                    )}
                  </>
                )}

                {/* Active: Solo suspender si hay problemas */}
                {raffle.status === 'active' && (
                  <Button
                    variant="outline"
                    onClick={handleSuspend}
                    className="w-full bg-yellow-600/10 border-yellow-400/20 text-yellow-100 hover:bg-yellow-600/20"
                  >
                    Suspender Sorteo
                  </Button>
                )}

                {/* Suspended: Reactivar o eliminar si no hay ventas */}
                {raffle.status === 'suspended' && (
                  <>
                    <Button
                      onClick={handleReactivate}
                      className="w-full bg-green-600/10 border-green-400/20 text-green-100 hover:bg-green-600/20"
                    >
                      Reactivar Sorteo
                    </Button>
                    {raffle.sold_count === 0 && (
                      <Button
                        variant="outline"
                        onClick={handleDelete}
                        disabled={deleteMutation.isPending}
                        className="w-full bg-red-600/10 border-red-400/20 text-red-100 hover:bg-red-600/20"
                      >
                        {deleteMutation.isPending ? 'Eliminando...' : 'Eliminar'}
                      </Button>
                    )}
                  </>
                )}

                {/* Completed sin ganador: Extender o cerrar */}
                {raffle.status === 'completed' && !raffle.winner_number && (
                  <>
                    <Button
                      onClick={handleExtendDate}
                      className="w-full bg-white/10 border-white/20 text-white hover:bg-white/20"
                    >
                      Extender Fecha
                    </Button>
                    <Button
                      variant="outline"
                      onClick={handleCloseDraw}
                      className="w-full bg-slate-600/10 border-slate-400/20 text-slate-100 hover:bg-slate-600/20"
                    >
                      Cerrar Sorteo
                    </Button>
                  </>
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
          <div className={cn("flex items-center justify-between text-sm mb-2", colors.textMuted)}>
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
            <span className="text-sm text-slate-600 dark:text-slate-400">
              {raffle.my_total_spent ? 'Mi Inversión' : raffle.total_revenue ? 'Recaudación' : 'Información'}
            </span>
            <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-3xl font-bold text-slate-900 dark:text-white">
            {raffle.my_total_spent
              ? formatCurrency(Number(raffle.my_total_spent))
              : raffle.total_revenue
              ? formatCurrency(Number(raffle.total_revenue))
              : '-'}
          </p>
          {raffle.my_total_spent && raffle.my_numbers_count && (
            <p className="text-sm text-slate-600 dark:text-slate-400 mt-1">
              {raffle.my_numbers_count} número(s) comprado(s)
            </p>
          )}
        </div>
      </div>

      {/* Image Gallery */}
      {data.images && data.images.length > 0 && (
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <h2 className="text-xl font-semibold text-slate-900 dark:text-white mb-6">
            Galería de Imágenes
          </h2>
          <RaffleImageGallery images={data.images} />
        </div>
      )}

      {/* Organizer Info - Importante para coordinar pago */}
      {raffle.organizer && !isOwner && (
        <div className="bg-gradient-to-r from-blue-50 to-teal-50 dark:from-slate-800 dark:to-slate-800 rounded-lg border border-blue-200 dark:border-slate-700 p-6">
          <div className="flex items-start gap-4">
            {/* Avatar */}
            <div className="flex-shrink-0">
              <div className="w-14 h-14 rounded-full bg-gradient-to-br from-blue-500 to-teal-500 flex items-center justify-center text-white text-xl font-bold shadow-lg">
                {raffle.organizer.name.charAt(0).toUpperCase()}
              </div>
            </div>

            {/* Info */}
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <h3 className="text-lg font-semibold text-slate-900 dark:text-white">
                  {raffle.organizer.name}
                </h3>
                {raffle.organizer.verified && (
                  <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 text-xs font-medium rounded-full">
                    <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    Verificado
                  </span>
                )}
              </div>
              <p className="text-sm text-slate-600 dark:text-slate-400 mb-3">
                Organizador de este sorteo
              </p>

              {/* Info box para coordinación de pago */}
              <div className="bg-white/60 dark:bg-slate-700/50 rounded-lg p-3 border border-blue-100 dark:border-slate-600">
                <div className="flex items-start gap-2">
                  <svg className="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p className="text-sm text-slate-700 dark:text-slate-300">
                    <span className="font-medium">Coordina el pago directamente</span> con el organizador.
                    Tienes <span className="font-semibold text-blue-600 dark:text-blue-400">24 horas</span> después de reservar para completar el pago.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

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

      {/* Lista de Compradores - Solo para Owner */}
      {isOwner && (
        <div className="bg-white dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700 p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold text-slate-900 dark:text-white">
              Compradores y Reservaciones
            </h2>
            {buyersData && (
              <div className="flex items-center gap-4 text-sm">
                <span className="text-green-600 dark:text-green-400">
                  {buyersData.total_sold} vendidos
                </span>
                <span className="text-yellow-600 dark:text-yellow-400">
                  {buyersData.total_reserved} reservados
                </span>
              </div>
            )}
          </div>

          {isLoadingBuyers ? (
            <div className="flex items-center justify-center py-8">
              <LoadingSpinner />
            </div>
          ) : !buyersData?.buyers || buyersData.buyers.length === 0 ? (
            <div className="text-center py-8">
              <svg className="w-12 h-12 text-slate-300 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
              <p className="text-slate-500 font-medium">No hay compradores aún</p>
              <p className="text-sm text-slate-400 mt-1">
                Cuando alguien reserve o compre números, aparecerá aquí
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {buyersData.buyers.map((buyer) => (
                <div
                  key={buyer.user_id}
                  className={cn(
                    "rounded-lg p-4 border",
                    buyer.status === 'sold'
                      ? "bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800"
                      : buyer.status === 'reserved'
                      ? "bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800"
                      : "bg-slate-50 dark:bg-slate-700/50 border-slate-200 dark:border-slate-600"
                  )}
                >
                  <div className="flex items-start justify-between gap-4">
                    {/* Info del comprador */}
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="font-semibold text-slate-900 dark:text-white truncate">
                          {buyer.name}
                        </h3>
                        <span className={cn(
                          "px-2 py-0.5 text-xs font-medium rounded-full",
                          buyer.status === 'sold'
                            ? "bg-green-100 text-green-700 dark:bg-green-800 dark:text-green-200"
                            : "bg-yellow-100 text-yellow-700 dark:bg-yellow-800 dark:text-yellow-200"
                        )}>
                          {buyer.status === 'sold' ? 'Vendido' : 'Reservado'}
                        </span>
                      </div>

                      {/* Email y teléfono copiables */}
                      <div className="flex flex-wrap items-center gap-3 text-sm text-slate-600 dark:text-slate-400 mb-2">
                        <button
                          onClick={() => {
                            navigator.clipboard.writeText(buyer.email);
                            toast.success('Email copiado');
                          }}
                          className="flex items-center gap-1 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                          title="Click para copiar"
                        >
                          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                          </svg>
                          {buyer.email}
                        </button>
                        {buyer.phone && (
                          <button
                            onClick={() => {
                              navigator.clipboard.writeText(buyer.phone!);
                              toast.success('Teléfono copiado');
                            }}
                            className="flex items-center gap-1 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                            title="Click para copiar"
                          >
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                            </svg>
                            {buyer.phone}
                          </button>
                        )}
                      </div>

                      {/* Números */}
                      <div className="flex flex-wrap gap-1">
                        {buyer.numbers.sort((a, b) => Number(a) - Number(b)).map((num) => (
                          <span
                            key={num}
                            className={cn(
                              "px-2 py-0.5 text-xs font-mono rounded",
                              buyer.status === 'sold'
                                ? "bg-green-200 text-green-800 dark:bg-green-700 dark:text-green-100"
                                : "bg-yellow-200 text-yellow-800 dark:bg-yellow-700 dark:text-yellow-100"
                            )}
                          >
                            {num}
                          </span>
                        ))}
                      </div>

                      {/* Expiración para reservados */}
                      {buyer.status === 'reserved' && buyer.expires_at && (
                        <p className="text-xs text-yellow-600 dark:text-yellow-400 mt-2">
                          Expira: {new Date(buyer.expires_at).toLocaleString('es-CR')}
                        </p>
                      )}
                    </div>

                    {/* Total */}
                    <div className="text-right flex-shrink-0">
                      <p className="text-lg font-bold text-slate-900 dark:text-white">
                        {formatCurrency(Number(buyer.total_amount))}
                      </p>
                      <p className="text-xs text-slate-500">
                        {buyer.numbers.length} número{buyer.numbers.length !== 1 ? 's' : ''}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Floating Checkout Button */}
      {!isOwner && raffle.status === 'active' && selectedNumbers.length > 0 && (
        <FloatingCheckoutButton
          selectedCount={selectedNumbers.length}
          selectedNumbers={selectedNumbers}
          totalAmount={selectedNumbers.length * Number(raffle.price_per_number)}
          onCheckout={handleReserve}
          onCancel={handleClearSelection}
          disabled={!user || user?.kyc_level === 'none' || isLoadingReservation}
        />
      )}
    </div>
  );
}
