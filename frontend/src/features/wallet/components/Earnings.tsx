import { DollarSign, TrendingUp, Percent, AlertCircle } from 'lucide-react';
import { Card } from '../../../components/ui/Card';
import { Alert } from '../../../components/ui/Alert';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { formatCRC } from '../../../types/wallet';
import { useEarnings } from '../hooks/useEarnings';

export const Earnings = () => {
  const { data, isLoading, error } = useEarnings();

  // DEBUG: Ver qué datos llegan
  console.log('🔍 Earnings Component - Raw data:', data);
  console.log('🔍 Earnings Component - isLoading:', isLoading);
  console.log('🔍 Earnings Component - error:', error);

  // Usar datos reales del backend
  const totalCollected = parseFloat(data?.total_collected || '0');
  const platformCommission = parseFloat(data?.platform_commission || '0');
  const netEarnings = parseFloat(data?.net_earnings || '0');
  const completedRafflesCount = data?.completed_raffles || 0;

  console.log('🔍 Earnings Component - Parsed values:', {
    totalCollected,
    platformCommission,
    netEarnings,
    completedRafflesCount
  });

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center">
          <LoadingSpinner size="md" />
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Info alert */}
      <Alert variant="info">
        <AlertCircle className="w-4 h-4" />
        <div className="text-sm">
          <p className="font-medium mb-1">¿Cómo funcionan las ganancias?</p>
          <p>Aquí puedes ver las <strong>ganancias estimadas de tus sorteos activos</strong> que tienen ventas. El monto mostrado es el total recolectado menos una comisión del 10% por el uso de la plataforma. Las ganancias se depositan automáticamente en tu billetera cuando el sorteo finaliza y el ganador confirma la recepción del premio.</p>
        </div>
      </Alert>

      {/* Resumen de ganancias */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Total recolectado */}
        <Card className="p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-slate-600">Total Recolectado</span>
            <DollarSign className="w-5 h-5 text-blue-600" />
          </div>
          <p className="text-2xl font-bold text-slate-900">{formatCRC(totalCollected)}</p>
          <p className="text-xs text-slate-500 mt-1">
            De {completedRafflesCount} sorteo{completedRafflesCount !== 1 ? 's' : ''} activo{completedRafflesCount !== 1 ? 's' : ''}
          </p>
        </Card>

        {/* Comisión de plataforma */}
        <Card className="p-6 bg-orange-50 border-orange-200">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-orange-700">Comisión Plataforma</span>
            <Percent className="w-5 h-5 text-orange-600" />
          </div>
          <p className="text-2xl font-bold text-orange-900">-{formatCRC(platformCommission)}</p>
          <p className="text-xs text-orange-600 mt-1">10% por uso de la plataforma</p>
        </Card>

        {/* Ganancias netas */}
        <Card className="p-6 bg-green-50 border-green-200">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-green-700">Ganancias Netas</span>
            <TrendingUp className="w-5 h-5 text-green-600" />
          </div>
          <p className="text-2xl font-bold text-green-900">{formatCRC(netEarnings)}</p>
          <p className="text-xs text-green-600 mt-1">Depositado en tu billetera</p>
        </Card>
      </div>

      {/* Desglose por sorteo */}
      <Card className="p-6">
        <h3 className="font-semibold text-slate-900 mb-4">Desglose por Sorteo</h3>

        {!data?.raffles || data.raffles.length === 0 ? (
          <div className="text-center py-8">
            <DollarSign className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <p className="text-slate-500 font-medium">No tienes sorteos activos con ventas aún</p>
            <p className="text-sm text-slate-400 mt-1">
              Cuando tus sorteos tengan ventas, verás aquí el desglose de tus ganancias estimadas
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-slate-200">
                  <th className="text-left py-3 px-2 text-sm font-semibold text-slate-700">Sorteo</th>
                  <th className="text-left py-3 px-2 text-sm font-semibold text-slate-700">Fecha</th>
                  <th className="text-right py-3 px-2 text-sm font-semibold text-slate-700">Recaudado</th>
                  <th className="text-right py-3 px-2 text-sm font-semibold text-slate-700">Comisión</th>
                  <th className="text-right py-3 px-2 text-sm font-semibold text-slate-700">Ganancia Neta</th>
                  <th className="text-center py-3 px-2 text-sm font-semibold text-slate-700">Estado</th>
                </tr>
              </thead>
              <tbody>
                {data.raffles.map((raffle) => {
                  const revenue = parseFloat(raffle.total_revenue);
                  const commission = parseFloat(raffle.platform_fee_amount);
                  const net = parseFloat(raffle.net_amount);
                  const drawDate = new Date(raffle.draw_date).toLocaleDateString('es-CR', {
                    day: '2-digit',
                    month: 'short',
                    year: 'numeric'
                  });

                  return (
                    <tr key={raffle.raffle_id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="py-3 px-2">
                        <div className="font-medium text-slate-900">{raffle.title}</div>
                        <div className="text-xs text-slate-500">ID: {raffle.raffle_uuid.substring(0, 8)}</div>
                      </td>
                      <td className="py-3 px-2 text-sm text-slate-600">{drawDate}</td>
                      <td className="py-3 px-2 text-right font-semibold text-slate-900">{formatCRC(revenue)}</td>
                      <td className="py-3 px-2 text-right font-semibold text-orange-600">-{formatCRC(commission)}</td>
                      <td className="py-3 px-2 text-right font-semibold text-green-700">{formatCRC(net)}</td>
                      <td className="py-3 px-2 text-center">
                        <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                          raffle.settlement_status === 'completed'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-amber-100 text-amber-800'
                        }`}>
                          {raffle.settlement_status === 'completed' ? 'Liquidado' : 'Pendiente'}
                        </span>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  );
};
