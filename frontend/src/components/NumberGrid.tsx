import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/Button';
import { cn } from '@/lib/utils';
import { Check, Lock, User } from 'lucide-react';

interface RaffleNumber {
  id: string;
  number: string;
  status: 'available' | 'reserved' | 'sold';
  user_id?: string;
}

interface NumberGridProps {
  numbers: RaffleNumber[];
  selectedNumbers: string[];
  onSelectNumber: (numberId: string) => void;
  onNumberUpdate?: (numberId: string, status: 'available' | 'reserved' | 'sold') => void;
  disabled?: boolean;
  currentUserId?: string;
  className?: string;
}

export function NumberGrid({
  numbers,
  selectedNumbers,
  onSelectNumber,
  onNumberUpdate,
  disabled = false,
  currentUserId,
  className
}: NumberGridProps) {
  const [localNumbers, setLocalNumbers] = useState<RaffleNumber[]>(numbers);

  // Actualizar números locales cuando cambie el prop
  useEffect(() => {
    setLocalNumbers(numbers);
  }, [numbers]);

  // Callback para actualizaciones externas (WebSocket)
  useEffect(() => {
    // Este efecto mantiene la referencia sincronizada
    // El callback será llamado desde el componente padre
  }, [onNumberUpdate]);

  const handleClick = (number: RaffleNumber) => {
    // No permitir seleccionar si está deshabilitado
    if (disabled) return;

    // Permitir deseleccionar números propios
    if (selectedNumbers.includes(number.id)) {
      onSelectNumber(number.id);
      return;
    }

    // Solo permitir seleccionar números disponibles
    if (number.status === 'available') {
      onSelectNumber(number.id);
    }
  };

  const getButtonVariant = (number: RaffleNumber): 'default' | 'outline' | 'secondary' | 'destructive' => {
    if (selectedNumbers.includes(number.id)) return 'default';
    if (number.status === 'available') return 'outline';
    return 'secondary';
  };

  const getButtonIcon = (number: RaffleNumber) => {
    if (selectedNumbers.includes(number.id)) {
      return <Check className="w-4 h-4" />;
    }
    if (number.status === 'sold') {
      return <Lock className="w-4 h-4" />;
    }
    if (number.status === 'reserved' && number.user_id !== currentUserId) {
      return <User className="w-4 h-4" />;
    }
    return null;
  };

  const isDisabled = (number: RaffleNumber): boolean => {
    if (disabled) return true;
    if (selectedNumbers.includes(number.id)) return false;
    if (number.status === 'reserved' && number.user_id === currentUserId) return false;
    return number.status !== 'available';
  };

  return (
    <div className={cn('grid grid-cols-5 sm:grid-cols-10 gap-2', className)}>
      {localNumbers.map((number) => {
        const icon = getButtonIcon(number);
        const buttonDisabled = isDisabled(number);

        return (
          <Button
            key={number.id}
            variant={getButtonVariant(number)}
            size="lg"
            disabled={buttonDisabled}
            onClick={() => handleClick(number)}
            className={cn(
              'relative h-14 transition-all duration-200 font-semibold text-base',
              selectedNumbers.includes(number.id) && 'ring-2 ring-blue-500 ring-offset-2',
              number.status === 'sold' && 'opacity-50 cursor-not-allowed',
              number.status === 'reserved' && !selectedNumbers.includes(number.id) && 'opacity-60',
              !buttonDisabled && 'hover:scale-105 active:scale-95'
            )}
          >
            <div className="flex flex-col items-center justify-center gap-1">
              <span>{number.number}</span>
              {icon && <span className="text-xs">{icon}</span>}
            </div>

            {/* Indicador de estado en la esquina */}
            {number.status !== 'available' && !selectedNumbers.includes(number.id) && (
              <div
                className={cn(
                  'absolute top-1 right-1 w-2 h-2 rounded-full',
                  number.status === 'sold' ? 'bg-red-500' : 'bg-yellow-500'
                )}
              />
            )}
          </Button>
        );
      })}
    </div>
  );
}
