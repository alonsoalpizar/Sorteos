import { useEffect, useState } from 'react';

interface ReservationTimerProps {
  expiresAt: Date | string;
  onExpire?: () => void;
}

export function ReservationTimer({ expiresAt, onExpire }: ReservationTimerProps) {
  const [timeLeft, setTimeLeft] = useState<number>(0);

  useEffect(() => {
    const calculateTimeLeft = () => {
      const expiry = typeof expiresAt === 'string' ? new Date(expiresAt) : expiresAt;
      const now = new Date();
      const diff = expiry.getTime() - now.getTime();

      if (diff <= 0) {
        setTimeLeft(0);
        if (onExpire) {
          onExpire();
        }
        return 0;
      }

      setTimeLeft(diff);
      return diff;
    };

    // Calculate immediately
    calculateTimeLeft();

    // Update every second
    const interval = setInterval(() => {
      calculateTimeLeft();
    }, 1000);

    return () => clearInterval(interval);
  }, [expiresAt, onExpire]);

  const minutes = Math.floor(timeLeft / 1000 / 60);
  const seconds = Math.floor((timeLeft / 1000) % 60);

  const isUrgent = timeLeft < 60000; // Less than 1 minute
  const isExpired = timeLeft === 0;

  if (isExpired) {
    return (
      <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
        <div className="flex items-center gap-3">
          <svg
            className="w-6 h-6 text-red-600 dark:text-red-400"
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
          <div>
            <p className="font-semibold text-red-900 dark:text-red-100">
              Reserva expirada
            </p>
            <p className="text-sm text-red-700 dark:text-red-300">
              Tu reserva ha expirado. Los números han sido liberados.
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className={`border rounded-lg p-4 transition-colors ${
        isUrgent
          ? 'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800'
          : 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800'
      }`}
    >
      <div className="flex items-center gap-3">
        <svg
          className={`w-6 h-6 ${
            isUrgent
              ? 'text-yellow-600 dark:text-yellow-400'
              : 'text-blue-600 dark:text-blue-400'
          }`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <div className="flex-1">
          <p
            className={`text-sm ${
              isUrgent
                ? 'text-yellow-700 dark:text-yellow-300'
                : 'text-blue-700 dark:text-blue-300'
            }`}
          >
            {isUrgent ? '¡Apresúrate!' : 'Tiempo restante para completar tu compra'}
          </p>
          <p
            className={`text-2xl font-bold font-mono ${
              isUrgent
                ? 'text-yellow-900 dark:text-yellow-100'
                : 'text-blue-900 dark:text-blue-100'
            }`}
          >
            {String(minutes).padStart(2, '0')}:{String(seconds).padStart(2, '0')}
          </p>
        </div>
      </div>
    </div>
  );
}
