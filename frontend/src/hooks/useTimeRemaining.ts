import { useState, useEffect, useRef } from 'react';

interface TimeRemaining {
  minutes: number;
  seconds: number;
  total: number; // Total en segundos
  isExpired: boolean;
  isUrgent: boolean; // < 2 minutos
}

export function useTimeRemaining(expiresAt: string | Date | null | undefined): TimeRemaining {
  const [timeRemaining, setTimeRemaining] = useState<TimeRemaining>({
    minutes: 0,
    seconds: 0,
    total: 0,
    isExpired: true,
    isUrgent: false,
  });

  const intervalRef = useRef<NodeJS.Timeout>();

  useEffect(() => {
    if (!expiresAt) {
      setTimeRemaining({
        minutes: 0,
        seconds: 0,
        total: 0,
        isExpired: true,
        isUrgent: false,
      });
      return;
    }

    const calculateTimeRemaining = () => {
      const now = new Date().getTime();
      const expiryTime = new Date(expiresAt).getTime();
      const difference = expiryTime - now;

      if (difference <= 0) {
        setTimeRemaining({
          minutes: 0,
          seconds: 0,
          total: 0,
          isExpired: true,
          isUrgent: false,
        });
        if (intervalRef.current) {
          clearInterval(intervalRef.current);
        }
        return;
      }

      const totalSeconds = Math.floor(difference / 1000);
      const minutes = Math.floor(totalSeconds / 60);
      const seconds = totalSeconds % 60;
      const isUrgent = totalSeconds < 120; // Menos de 2 minutos

      setTimeRemaining({
        minutes,
        seconds,
        total: totalSeconds,
        isExpired: false,
        isUrgent,
      });
    };

    // Calcular inmediatamente
    calculateTimeRemaining();

    // Actualizar cada segundo
    intervalRef.current = setInterval(calculateTimeRemaining, 1000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [expiresAt]);

  return timeRemaining;
}
