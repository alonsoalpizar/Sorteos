import { useEffect, useRef, useState, useCallback } from 'react';
import { toast } from 'sonner';

interface NumberUpdate {
  number_id: string;
  status: 'available' | 'reserved' | 'sold';
  user_id?: string;
}

interface ReservationExpired {
  number_ids: string[];
}

interface WebSocketMessage {
  type: 'number_update' | 'reservation_expired' | 'reservation_created';
  raffle_id: string;
  data: NumberUpdate | ReservationExpired;
}

interface UseRaffleWebSocketReturn {
  isConnected: boolean;
  connectionError: string | null;
  onNumberUpdate: (callback: (update: NumberUpdate) => void) => () => void;
  onReservationExpired: (callback: (data: ReservationExpired) => void) => () => void;
}

export function useRaffleWebSocket(raffleId: string | undefined): UseRaffleWebSocketReturn {
  const ws = useRef<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [connectionError, setConnectionError] = useState<string | null>(null);
  const reconnectTimeout = useRef<NodeJS.Timeout>();
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;

  const numberUpdateCallbacks = useRef<Set<(update: NumberUpdate) => void>>(new Set());
  const reservationExpiredCallbacks = useRef<Set<(data: ReservationExpired) => void>>(new Set());

  const connect = useCallback(() => {
    if (!raffleId) {
      console.warn('[WebSocket] No raffle ID provided, skipping connection');
      return;
    }

    // Obtener la URL del WebSocket desde variables de entorno
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = import.meta.env.VITE_API_URL?.replace(/^https?:\/\//, '') || window.location.host;
    const wsUrl = `${wsProtocol}//${wsHost}/api/v1/raffles/${raffleId}/ws`;

    console.log('[WebSocket] Connecting to:', wsUrl);

    try {
      ws.current = new WebSocket(wsUrl);

      ws.current.onopen = () => {
        console.log('[WebSocket] Connected to raffle', raffleId);
        setIsConnected(true);
        setConnectionError(null);
        reconnectAttempts.current = 0;

        toast.success('Conexión en tiempo real activada', {
          duration: 2000,
          icon: '🔌',
        });
      };

      ws.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);

          console.log('[WebSocket] Message received:', message.type, message.data);

          if (message.type === 'number_update') {
            const data = message.data as NumberUpdate;
            numberUpdateCallbacks.current.forEach(cb => cb(data));
          } else if (message.type === 'reservation_expired') {
            const data = message.data as ReservationExpired;
            reservationExpiredCallbacks.current.forEach(cb => cb(data));

            toast.info(`${data.number_ids.length} números han sido liberados`, {
              duration: 3000,
              icon: '🔓',
            });
          }
        } catch (error) {
          console.error('[WebSocket] Error parsing message:', error);
        }
      };

      ws.current.onclose = (event) => {
        console.log('[WebSocket] Disconnected', event.code, event.reason);
        setIsConnected(false);

        // Auto-reconnect con backoff exponencial
        if (reconnectAttempts.current < maxReconnectAttempts) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000);
          console.log(`[WebSocket] Reconnecting in ${delay}ms (attempt ${reconnectAttempts.current + 1}/${maxReconnectAttempts})`);

          reconnectTimeout.current = setTimeout(() => {
            reconnectAttempts.current++;
            connect();
          }, delay);
        } else {
          setConnectionError('No se pudo establecer conexión en tiempo real');
          toast.error('Conexión perdida. Por favor recarga la página.', {
            duration: 5000,
          });
        }
      };

      ws.current.onerror = (error) => {
        console.error('[WebSocket] Error:', error);
        setConnectionError('Error de conexión');
      };
    } catch (error) {
      console.error('[WebSocket] Failed to create connection:', error);
      setConnectionError('Error al crear conexión WebSocket');
    }
  }, [raffleId]);

  useEffect(() => {
    if (!raffleId) return;

    connect();

    return () => {
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }

      if (ws.current) {
        console.log('[WebSocket] Closing connection');
        ws.current.close();
        ws.current = null;
      }

      // Limpiar callbacks
      numberUpdateCallbacks.current.clear();
      reservationExpiredCallbacks.current.clear();
    };
  }, [connect, raffleId]);

  const onNumberUpdate = useCallback((callback: (update: NumberUpdate) => void) => {
    numberUpdateCallbacks.current.add(callback);

    // Retornar función de cleanup
    return () => {
      numberUpdateCallbacks.current.delete(callback);
    };
  }, []);

  const onReservationExpired = useCallback((callback: (data: ReservationExpired) => void) => {
    reservationExpiredCallbacks.current.add(callback);

    // Retornar función de cleanup
    return () => {
      reservationExpiredCallbacks.current.delete(callback);
    };
  }, []);

  return {
    isConnected,
    connectionError,
    onNumberUpdate,
    onReservationExpired,
  };
}
