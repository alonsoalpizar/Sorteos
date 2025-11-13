import { CreditCard, X } from 'lucide-react';
import { Button } from './Button';
import { cn } from '@/lib/utils';

interface FloatingCheckoutButtonProps {
  selectedCount: number;
  totalAmount: number;
  onCheckout: () => void;
  onCancel: () => void;
  disabled?: boolean;
  className?: string;
}

export function FloatingCheckoutButton({
  selectedCount,
  totalAmount,
  onCheckout,
  onCancel,
  disabled = false,
  className
}: FloatingCheckoutButtonProps) {
  if (selectedCount === 0) return null;

  return (
    <div
      className={cn(
        'fixed bottom-6 left-1/2 -translate-x-1/2 z-50 animate-in slide-in-from-bottom-10 duration-300',
        className
      )}
    >
      <div className="bg-white rounded-full shadow-2xl border-2 border-blue-500 px-6 py-4 flex items-center gap-6">
        {/* Contador de números */}
        <div className="flex items-center gap-3">
          <div className="bg-blue-500 text-white rounded-full w-10 h-10 flex items-center justify-center font-bold">
            {selectedCount}
          </div>
          <div>
            <p className="text-sm text-gray-600 font-medium">
              {selectedCount === 1 ? '1 número' : `${selectedCount} números`}
            </p>
            <p className="text-lg font-bold text-gray-900">
              ₡{totalAmount.toLocaleString()}
            </p>
          </div>
        </div>

        {/* Botones de acción */}
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={onCancel}
            className="rounded-full"
          >
            <X className="w-4 h-4 mr-1" />
            Cancelar
          </Button>

          <Button
            size="lg"
            onClick={onCheckout}
            disabled={disabled}
            className="rounded-full bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 shadow-lg"
          >
            <CreditCard className="w-5 h-5 mr-2" />
            Pagar Ahora
          </Button>
        </div>
      </div>
    </div>
  );
}
