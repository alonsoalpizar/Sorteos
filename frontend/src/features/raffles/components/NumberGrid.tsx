import { cn } from '../../../lib/utils';
import type { RaffleNumber } from '../../../types/raffle';

interface NumberGridProps {
  numbers: RaffleNumber[];
  selectedNumbers?: string[];
  onNumberSelect?: (number: string) => void;
  readonly?: boolean;
}

export function NumberGrid({
  numbers,
  selectedNumbers = [],
  onNumberSelect,
  readonly = false,
}: NumberGridProps) {
  const getNumberStyle = (number: RaffleNumber) => {
    const isSelected = selectedNumbers.includes(number.number);

    if (number.status === 'sold') {
      return 'bg-slate-400 text-white cursor-not-allowed';
    }

    if (number.status === 'reserved') {
      return 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200 cursor-not-allowed';
    }

    if (isSelected) {
      return 'bg-blue-600 text-white border-blue-600 cursor-pointer hover:bg-blue-700 hover:border-blue-700';
    }

    if (readonly) {
      return 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200 cursor-default';
    }

    return 'bg-white dark:bg-slate-800 text-slate-900 dark:text-white border-slate-200 dark:border-slate-700 hover:border-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/20 cursor-pointer';
  };

  const handleNumberClick = (number: RaffleNumber) => {
    if (readonly) return;
    if (!onNumberSelect) return;

    const isSelected = selectedNumbers.includes(number.number);

    // Permitir des-seleccionar si ya está seleccionado
    if (isSelected) {
      onNumberSelect(number.number);
      return;
    }

    // Solo permitir seleccionar si está disponible
    if (number.status !== 'available') return;
    onNumberSelect(number.number);
  };

  return (
    <div className="w-full">
      {/* Legend */}
      <div className="flex flex-wrap items-center gap-4 mb-6 text-sm">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded border-2 bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700"></div>
          <span className="text-slate-600 dark:text-slate-400">Disponible</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded bg-blue-600"></div>
          <span className="text-slate-600 dark:text-slate-400">Seleccionado</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded bg-yellow-100 dark:bg-yellow-900/30"></div>
          <span className="text-slate-600 dark:text-slate-400">Reservado</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded bg-slate-400"></div>
          <span className="text-slate-600 dark:text-slate-400">Vendido</span>
        </div>
      </div>

      {/* Grid */}
      <div className="grid grid-cols-10 gap-2">
        {numbers.map((number) => {
          const isSelected = selectedNumbers.includes(number.number);
          const isDisabled = readonly || (number.status !== 'available' && !isSelected);

          return (
            <button
              key={number.id}
              onClick={() => handleNumberClick(number)}
              disabled={isDisabled}
              className={cn(
                'aspect-square rounded border-2 font-mono font-semibold text-sm transition-all',
                'flex items-center justify-center',
                getNumberStyle(number)
              )}
              title={`Número ${number.number} - ${
                number.status === 'sold'
                  ? 'Vendido'
                  : number.status === 'reserved'
                  ? 'Reservado'
                  : isSelected
                  ? 'Seleccionado (clic para des-reservar)'
                  : 'Disponible (clic para reservar)'
              }`}
            >
              {number.number}
            </button>
          );
        })}
      </div>

      {/* Selection summary */}
      {!readonly && selectedNumbers.length > 0 && (
        <div className="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-slate-900 dark:text-white">
                Números seleccionados: {selectedNumbers.length}
              </p>
              <p className="text-xs text-slate-600 dark:text-slate-400 mt-1">
                {selectedNumbers.sort((a, b) => Number(a) - Number(b)).join(', ')}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Mobile view note */}
      <p className="mt-4 text-xs text-slate-500 dark:text-slate-400 sm:hidden">
        Desliza horizontalmente para ver todos los números
      </p>
    </div>
  );
}
