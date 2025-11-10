import { Link } from 'react-router-dom';
import { Card } from '../../../components/ui/Card';
import { Badge } from '../../../components/ui/Badge';
import type { Raffle } from '../../../types/raffle';
import { formatCurrency, formatDate, getStatusColor, getStatusLabel } from '../../../lib/utils';

interface RaffleCardProps {
  raffle: Raffle;
}

export function RaffleCard({ raffle }: RaffleCardProps) {
  const availableCount = raffle.total_numbers - raffle.sold_count - raffle.reserved_count;
  const soldPercentage = (raffle.sold_count / raffle.total_numbers) * 100;

  return (
    <Link to={`/raffles/${raffle.id}`}>
      <Card className="hover:shadow-lg transition-shadow cursor-pointer h-full">
        <div className="p-6 space-y-4">
          {/* Header */}
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <h3 className="text-lg font-semibold text-slate-900 dark:text-white line-clamp-2">
                {raffle.title}
              </h3>
              <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                {formatDate(raffle.draw_date)}
              </p>
            </div>
            <Badge variant={getStatusColor(raffle.status)}>
              {getStatusLabel(raffle.status)}
            </Badge>
          </div>

          {/* Description */}
          <p className="text-sm text-slate-600 dark:text-slate-300 line-clamp-2">
            {raffle.description}
          </p>

          {/* Stats */}
          <div className="space-y-2">
            {/* Progress bar */}
            <div className="w-full bg-slate-200 dark:bg-slate-700 rounded-full h-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all"
                style={{ width: `${soldPercentage}%` }}
              />
            </div>

            {/* Numbers info */}
            <div className="flex items-center justify-between text-sm">
              <span className="text-slate-600 dark:text-slate-400">
                {raffle.sold_count} de {raffle.total_numbers} vendidos
              </span>
              <span className="text-slate-600 dark:text-slate-400">
                {soldPercentage.toFixed(0)}%
              </span>
            </div>

            {/* Available numbers */}
            <div className="flex items-center gap-2 text-sm">
              <span className="text-slate-600 dark:text-slate-400">Disponibles:</span>
              <span className="font-semibold text-slate-900 dark:text-white">
                {availableCount}
              </span>
            </div>
          </div>

          {/* Price */}
          <div className="pt-4 border-t border-slate-200 dark:border-slate-700">
            <div className="flex items-center justify-between">
              <span className="text-sm text-slate-600 dark:text-slate-400">
                Precio por número
              </span>
              <span className="text-xl font-bold text-blue-600">
                {formatCurrency(Number(raffle.price_per_number))}
              </span>
            </div>
          </div>

          {/* Draw method */}
          <div className="flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
            <svg
              className="w-4 h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span>
              {raffle.draw_method === 'loteria_nacional_cr'
                ? 'Lotería Nacional CR'
                : raffle.draw_method === 'manual'
                ? 'Sorteo Manual'
                : 'Sorteo Aleatorio'}
            </span>
          </div>
        </div>
      </Card>
    </Link>
  );
}
