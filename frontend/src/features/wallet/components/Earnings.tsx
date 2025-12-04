import { DollarSign, AlertCircle, BarChart3 } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useEarnings } from "../hooks/useWallet";

// ===========================================
// NOTA: Comisión de plataforma desactivada
// ===========================================
// Sorteos.club es ahora solo plataforma de gestión.
// No se cobra comisión por los sorteos.
// El código original que mostraba comisiones está comentado abajo.
// ===========================================

// Helper para formatear CRC
function formatCRC(amount: number | string): string {
  const num = typeof amount === "string" ? parseFloat(amount) : amount;
  return new Intl.NumberFormat("es-CR", {
    style: "currency",
    currency: "CRC",
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(num);
}

export const Earnings = () => {
  const { data, isLoading } = useEarnings();

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center">
          <LoadingSpinner />
        </div>
      </Card>
    );
  }

  const totalCollected = parseFloat(data?.total_collected || "0");
  const completedRafflesCount = data?.completed_raffles || 0;

  return (
    <div className="space-y-6">
      {/* Info alert */}
      <div className="rounded-lg p-4 flex items-start gap-3 border bg-teal-50 border-teal-200">
        <AlertCircle className="w-5 h-5 flex-shrink-0 mt-0.5 text-teal-600" />
        <div className="text-sm text-teal-900">
          <p className="font-medium mb-1">Resumen de tus sorteos</p>
          <p>
            Aquí puedes ver el <strong>historial de ventas de tus sorteos completados</strong>.
            El monto mostrado es el total recaudado por la venta de números.
            Sorteos.club no cobra comisión - ¡el 100% es tuyo!
          </p>
        </div>
      </div>

      {/* Resumen de ganancias - Solo 2 cards sin comisión */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Total recolectado */}
        <Card className="p-6 bg-green-50 border-green-200">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-green-700">Total Recaudado</span>
            <DollarSign className="w-5 h-5 text-green-600" />
          </div>
          <p className="text-2xl font-bold text-green-900">{formatCRC(totalCollected)}</p>
          <p className="text-xs text-green-600 mt-1">
            De {completedRafflesCount} sorteo{completedRafflesCount !== 1 ? "s" : ""} completado
            {completedRafflesCount !== 1 ? "s" : ""}
          </p>
        </Card>

        {/* Sorteos completados */}
        <Card className="p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-slate-600">Sorteos Completados</span>
            <BarChart3 className="w-5 h-5 text-teal-600" />
          </div>
          <p className="text-2xl font-bold text-slate-900">{completedRafflesCount}</p>
          <p className="text-xs text-slate-500 mt-1">Total de sorteos finalizados</p>
        </Card>
      </div>

      {/* Desglose por sorteo - Sin columna de comisión */}
      <Card className="p-6">
        <h3 className="font-semibold text-slate-900 mb-4">Desglose por Sorteo</h3>

        {!data?.raffles || data.raffles.length === 0 ? (
          <div className="text-center py-8">
            <DollarSign className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <p className="text-slate-500 font-medium">No tienes sorteos completados aún</p>
            <p className="text-sm text-slate-400 mt-1">
              Cuando tus sorteos finalicen, verás aquí el desglose de tus ventas
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
                  <th className="text-center py-3 px-2 text-sm font-semibold text-slate-700">Estado</th>
                </tr>
              </thead>
              <tbody>
                {data.raffles.map((raffle) => {
                  const revenue = parseFloat(raffle.total_revenue);
                  const drawDate = new Date(raffle.draw_date).toLocaleDateString("es-CR", {
                    day: "2-digit",
                    month: "short",
                    year: "numeric",
                  });

                  return (
                    <tr key={raffle.raffle_id} className="border-b border-slate-100 hover:bg-slate-50">
                      <td className="py-3 px-2">
                        <div className="font-medium text-slate-900">{raffle.title}</div>
                        <div className="text-xs text-slate-500">ID: {raffle.raffle_uuid.substring(0, 8)}</div>
                      </td>
                      <td className="py-3 px-2 text-sm text-slate-600">{drawDate}</td>
                      <td className="py-3 px-2 text-right font-semibold text-green-700">
                        {formatCRC(revenue)}
                      </td>
                      <td className="py-3 px-2 text-center">
                        <span
                          className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                            raffle.settlement_status === "completed"
                              ? "bg-green-100 text-green-800"
                              : "bg-amber-100 text-amber-800"
                          }`}
                        >
                          {raffle.settlement_status === "completed" ? "Completado" : "Pendiente"}
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

      {/* CÓDIGO ORIGINAL CON COMISIONES - COMENTADO

      // Imports originales:
      // import { Percent } from "lucide-react";
      // import { useUserMode } from "@/contexts/UserModeContext";
      // import { cn } from "@/lib/utils";

      // Variables originales:
      // const platformCommission = parseFloat(data?.platform_commission || "0");
      // const netEarnings = parseFloat(data?.net_earnings || "0");
      // const { mode } = useUserMode();
      // const isOrganizer = mode === 'organizer';

      // Card de comisión original:
      // <Card className="p-6 bg-orange-50 border-orange-200">
      //   <div className="flex items-center justify-between mb-2">
      //     <span className="text-sm font-medium text-orange-700">Comisión Plataforma</span>
      //     <Percent className="w-5 h-5 text-orange-600" />
      //   </div>
      //   <p className="text-2xl font-bold text-orange-900">-{formatCRC(platformCommission)}</p>
      //   <p className="text-xs text-orange-600 mt-1">Según tarifas de plataforma</p>
      // </Card>

      // Columna de comisión en tabla:
      // <th className="text-right py-3 px-2 text-sm font-semibold text-slate-700">Comisión</th>
      // <td className="py-3 px-2 text-right font-semibold text-orange-600">
      //   -{formatCRC(commission)}
      // </td>

      */}
    </div>
  );
};
