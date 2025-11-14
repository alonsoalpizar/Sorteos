import { useEffect, useRef, useCallback } from 'react';
import { useAuthStore } from '@/store/authStore';

interface UseInactivityTimeoutOptions {
  /**
   * Tiempo de inactividad en milisegundos antes de hacer logout
   * @default 1800000 (30 minutos)
   */
  timeout?: number;

  /**
   * Tiempo en milisegundos antes del timeout para mostrar advertencia
   * @default 120000 (2 minutos)
   */
  warningTime?: number;

  /**
   * Callback cuando se detecta inactividad (antes del logout)
   */
  onWarning?: () => void;

  /**
   * Callback cuando se ejecuta el logout por inactividad
   */
  onTimeout?: () => void;
}

/**
 * Hook para manejar timeout de inactividad del usuario
 *
 * Detecta inactividad basándose en:
 * - Movimiento del mouse
 * - Clicks
 * - Pulsaciones de teclado
 * - Scroll
 * - Touch events (móviles)
 *
 * @example
 * ```tsx
 * useInactivityTimeout({
 *   timeout: 30 * 60 * 1000, // 30 minutos
 *   warningTime: 2 * 60 * 1000, // 2 minutos antes
 *   onWarning: () => toast.warning('Tu sesión expirará pronto'),
 *   onTimeout: () => toast.error('Sesión expirada por inactividad')
 * });
 * ```
 */
export function useInactivityTimeout({
  timeout = 30 * 60 * 1000, // 30 minutos por defecto
  warningTime = 2 * 60 * 1000, // 2 minutos por defecto
  onWarning,
  onTimeout,
}: UseInactivityTimeoutOptions = {}) {
  const { isAuthenticated, logout } = useAuthStore();
  const timeoutRef = useRef<NodeJS.Timeout>();
  const warningRef = useRef<NodeJS.Timeout>();
  const lastActivityRef = useRef<number>(Date.now());
  const warningShownRef = useRef<boolean>(false);

  // Función para hacer logout
  const handleTimeout = useCallback(() => {
    if (isAuthenticated) {
      console.log('[InactivityTimeout] Usuario inactivo por', timeout / 1000 / 60, 'minutos. Cerrando sesión...');
      logout();
      onTimeout?.();
    }
  }, [isAuthenticated, logout, onTimeout, timeout]);

  // Función para mostrar advertencia
  const handleWarning = useCallback(() => {
    if (isAuthenticated && !warningShownRef.current) {
      console.log('[InactivityTimeout] Mostrando advertencia de inactividad');
      warningShownRef.current = true;
      onWarning?.();
    }
  }, [isAuthenticated, onWarning]);

  // Función para resetear los timers
  const resetTimers = useCallback(() => {
    // Limpiar timers existentes
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    if (warningRef.current) {
      clearTimeout(warningRef.current);
    }

    // Resetear flag de advertencia
    warningShownRef.current = false;
    lastActivityRef.current = Date.now();

    // Solo configurar timers si el usuario está autenticado
    if (isAuthenticated) {
      // Timer para advertencia (timeout - warningTime)
      const warningDelay = timeout - warningTime;
      warningRef.current = setTimeout(handleWarning, warningDelay);

      // Timer para logout
      timeoutRef.current = setTimeout(handleTimeout, timeout);
    }
  }, [isAuthenticated, timeout, warningTime, handleWarning, handleTimeout]);

  // Eventos que indican actividad del usuario
  useEffect(() => {
    if (!isAuthenticated) {
      return;
    }

    const events = [
      'mousedown',
      'mousemove',
      'keypress',
      'scroll',
      'touchstart',
      'click',
    ];

    // Throttle para evitar demasiadas llamadas
    let throttleTimeout: NodeJS.Timeout;
    const throttledResetTimers = () => {
      if (!throttleTimeout) {
        throttleTimeout = setTimeout(() => {
          resetTimers();
          throttleTimeout = undefined!;
        }, 1000); // Throttle de 1 segundo
      }
    };

    // Agregar listeners
    events.forEach((event) => {
      window.addEventListener(event, throttledResetTimers);
    });

    // Iniciar timers
    resetTimers();

    // Cleanup
    return () => {
      events.forEach((event) => {
        window.removeEventListener(event, throttledResetTimers);
      });
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      if (warningRef.current) {
        clearTimeout(warningRef.current);
      }
      if (throttleTimeout) {
        clearTimeout(throttleTimeout);
      }
    };
  }, [isAuthenticated, resetTimers]);

  // Retornar función para extender manualmente la sesión
  return {
    /**
     * Resetea manualmente el timer de inactividad
     * Útil después de acciones importantes del usuario
     */
    extendSession: resetTimers,

    /**
     * Tiempo transcurrido desde la última actividad (en ms)
     */
    getInactivityTime: () => Date.now() - lastActivityRef.current,
  };
}
