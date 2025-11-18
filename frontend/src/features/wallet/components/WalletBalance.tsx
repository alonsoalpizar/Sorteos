import { Wallet, RefreshCw } from 'lucide-react';
import { useWallet } from '../hooks/useWallet';
import { formatCRC } from '../../../types/wallet';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';

interface WalletBalanceProps {
  showRefreshButton?: boolean;
  compact?: boolean;
}

export const WalletBalance = ({ showRefreshButton = true, compact = false }: WalletBalanceProps) => {
  const { balance, pendingBalance, currency, status, isLoading, refetch } = useWallet();

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center">
          <LoadingSpinner size="md" />
        </div>
      </Card>
    );
  }

  const balanceAmount = parseFloat(balance) || 0;
  const pendingAmount = parseFloat(pendingBalance) || 0;
  const isFrozen = status === 'frozen';

  return (
    <Card className={`${compact ? 'p-4' : 'p-6'} ${isFrozen ? 'border-red-500' : ''}`}>
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-blue-100 rounded-lg">
            <Wallet className="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <p className="text-sm text-slate-600 font-medium">Saldo Disponible</p>
            <p className={`${compact ? 'text-2xl' : 'text-3xl'} font-bold text-slate-900`}>
              {formatCRC(balanceAmount)}
            </p>
            {pendingAmount > 0 && (
              <p className="text-sm text-slate-500 mt-1">
                Pendiente: {formatCRC(pendingAmount)}
              </p>
            )}
            {isFrozen && (
              <p className="text-sm text-red-600 font-medium mt-1">⚠️ Billetera congelada</p>
            )}
          </div>
        </div>

        {showRefreshButton && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => refetch()}
            className="text-slate-600 hover:text-slate-900"
          >
            <RefreshCw className="w-4 h-4" />
          </Button>
        )}
      </div>

      {!compact && (
        <div className="mt-4 pt-4 border-t border-slate-200">
          <div className="flex items-center justify-between text-sm">
            <span className="text-slate-600">Moneda</span>
            <span className="font-medium text-slate-900">{currency}</span>
          </div>
          <div className="flex items-center justify-between text-sm mt-2">
            <span className="text-slate-600">Estado</span>
            <span
              className={`font-medium ${
                status === 'active'
                  ? 'text-green-600'
                  : status === 'frozen'
                  ? 'text-red-600'
                  : 'text-slate-600'
              }`}
            >
              {status === 'active' ? 'Activa' : status === 'frozen' ? 'Congelada' : 'Cerrada'}
            </span>
          </div>
        </div>
      )}
    </Card>
  );
};
