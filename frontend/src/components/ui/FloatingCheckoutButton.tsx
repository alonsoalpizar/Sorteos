import { useState, useEffect } from 'react';
import { ShoppingCart, Clock, X } from 'lucide-react';
import { formatCurrency } from '@/lib/utils';
import { cn } from '@/lib/utils';

interface FloatingCheckoutButtonProps {
  selectedCount: number;
  totalAmount: number;
  expiresAt?: string | null;
  onCheckout: () => void;
  onClear: () => void;
  show: boolean;
}

export function FloatingCheckoutButton({
  selectedCount,
  totalAmount,
  expiresAt,
  onCheckout,
  onClear,
  show,
}: FloatingCheckoutButtonProps) {
  const [timeLeft, setTimeLeft] = useState<number | null>(null);
  const [isExpiring, setIsExpiring] = useState(false);

  useEffect(() => {
    if (!expiresAt) {
      setTimeLeft(null);
      return;
    }

    const calculateTimeLeft = () => {
      const now = new Date().getTime();
      const expiry = new Date(expiresAt).getTime();
      const diff = expiry - now;

      if (diff <= 0) {
        setTimeLeft(0);
        return;
      }

      setTimeLeft(Math.floor(diff / 1000)); // seconds

      // Warn when less than 2 minutes left
      if (diff < 120000) {
        setIsExpiring(true);
      }
    };

    calculateTimeLeft();
    const interval = setInterval(calculateTimeLeft, 1000);

    return () => clearInterval(interval);
  }, [expiresAt]);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  if (!show || selectedCount === 0) return null;

  return (
    <>
      {/* Overlay for mobile */}
      <div
        className={cn(
          'fixed inset-0 bg-black/20 backdrop-blur-sm z-40 transition-opacity lg:hidden',
          show ? 'opacity-100' : 'opacity-0 pointer-events-none'
        )}
        onClick={onClear}
      />

      {/* Floating Button */}
      <div
        className={cn(
          'fixed bottom-6 right-6 z-50 transition-all duration-300 transform',
          show ? 'translate-y-0 opacity-100 scale-100' : 'translate-y-20 opacity-0 scale-95 pointer-events-none'
        )}
      >
        <div className="relative">
          {/* Timer Warning */}
          {timeLeft !== null && timeLeft > 0 && isExpiring && (
            <div className="absolute -top-14 right-0 left-0 mx-auto w-max px-4 py-2 bg-yellow-500 text-white rounded-lg shadow-lg animate-bounce">
              <div className="flex items-center gap-2 text-sm font-medium">
                <Clock className="w-4 h-4" />
                <span>¡Reserva expira en {formatTime(timeLeft)}!</span>
              </div>
              <div className="absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-1/2 rotate-45 w-2 h-2 bg-yellow-500"></div>
            </div>
          )}

          {/* Main Button Container */}
          <div className="bg-gradient-to-br from-primary-600 to-primary-700 rounded-2xl shadow-2xl border border-primary-400/20 overflow-hidden min-w-[280px]">
            {/* Close button */}
            <button
              onClick={onClear}
              className="absolute top-2 right-2 p-1.5 bg-white/10 hover:bg-white/20 rounded-full transition-colors"
              aria-label="Limpiar selección"
            >
              <X className="w-4 h-4 text-white" />
            </button>

            {/* Content */}
            <div className="p-4 space-y-3">
              {/* Header */}
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 bg-white/10 rounded-xl flex items-center justify-center backdrop-blur-sm">
                  <ShoppingCart className="w-6 h-6 text-white" />
                </div>
                <div className="flex-1">
                  <p className="text-xs text-primary-100 font-medium">Números seleccionados</p>
                  <p className="text-2xl font-bold text-white">{selectedCount}</p>
                </div>
              </div>

              {/* Total Amount */}
              <div className="bg-white/10 backdrop-blur-sm rounded-xl p-3 border border-white/20">
                <p className="text-xs text-primary-100 mb-1">Total a pagar</p>
                <p className="text-2xl font-bold text-white">{formatCurrency(totalAmount)}</p>
              </div>

              {/* Timer Display (if active) */}
              {timeLeft !== null && timeLeft > 0 && !isExpiring && (
                <div className="flex items-center justify-center gap-2 text-xs text-primary-100">
                  <Clock className="w-3.5 h-3.5" />
                  <span>Reservado por {formatTime(timeLeft)}</span>
                </div>
              )}

              {/* Checkout Button */}
              <button
                onClick={onCheckout}
                className={cn(
                  'w-full py-3 px-4 rounded-xl font-semibold transition-all transform hover:scale-105 active:scale-95',
                  'flex items-center justify-center gap-2 shadow-lg',
                  timeLeft === 0
                    ? 'bg-red-500 hover:bg-red-600 text-white'
                    : 'bg-white text-primary-600 hover:bg-primary-50'
                )}
                disabled={timeLeft === 0}
              >
                {timeLeft === 0 ? (
                  <>
                    <X className="w-5 h-5" />
                    <span>Reserva expirada</span>
                  </>
                ) : (
                  <>
                    <ShoppingCart className="w-5 h-5" />
                    <span>Proceder al Pago</span>
                  </>
                )}
              </button>

              {/* Clear button */}
              <button
                onClick={onClear}
                className="w-full py-2 text-xs text-primary-100 hover:text-white transition-colors"
              >
                Limpiar selección
              </button>
            </div>

            {/* Pulse animation ring */}
            <div className="absolute inset-0 rounded-2xl border-2 border-primary-400 animate-ping opacity-20 pointer-events-none" />
          </div>
        </div>
      </div>
    </>
  );
}
