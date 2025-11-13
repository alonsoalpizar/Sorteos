import { Clock, AlertCircle } from 'lucide-react';
import { useTimeRemaining } from '@/hooks/useTimeRemaining';
import { cn } from '@/lib/utils';

interface ReservationTimerProps {
  expiresAt: string | Date | null | undefined;
  phase?: 'selection' | 'checkout';
  onExpire?: () => void;
  className?: string;
}

export function ReservationTimer({
  expiresAt,
  phase = 'selection',
  onExpire,
  className
}: ReservationTimerProps) {
  const { minutes, seconds, isExpired, isUrgent } = useTimeRemaining(expiresAt);

  // Ejecutar callback cuando expire
  if (isExpired && onExpire) {
    onExpire();
  }

  if (!expiresAt || isExpired) {
    return null;
  }

  const phaseText = phase === 'selection' ? 'Selección' : 'Pago';

  return (
    <div
      className={cn(
        'flex items-center gap-3 px-4 py-3 rounded-lg border-2 transition-all duration-300',
        isUrgent
          ? 'bg-red-50 border-red-500 animate-pulse'
          : 'bg-blue-50 border-blue-500',
        className
      )}
    >
      {isUrgent ? (
        <AlertCircle className="w-5 h-5 text-red-600" />
      ) : (
        <Clock className="w-5 h-5 text-blue-600" />
      )}

      <div className="flex-1">
        <p className={cn(
          'text-sm font-medium',
          isUrgent ? 'text-red-900' : 'text-blue-900'
        )}>
          Tiempo restante ({phaseText})
        </p>
        <p className={cn(
          'text-xs',
          isUrgent ? 'text-red-700' : 'text-blue-700'
        )}>
          {isUrgent ? '¡Última oportunidad!' : 'Completa tu reserva antes de que expire'}
        </p>
      </div>

      <div className={cn(
        'text-2xl font-bold tabular-nums',
        isUrgent ? 'text-red-600' : 'text-blue-600'
      )}>
        {String(minutes).padStart(2, '0')}:{String(seconds).padStart(2, '0')}
      </div>
    </div>
  );
}
