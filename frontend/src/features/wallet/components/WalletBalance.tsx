import { Wallet, RefreshCw } from "lucide-react";
import { useWalletBalance } from "../hooks/useWallet";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";

interface WalletBalanceProps {
  showRefreshButton?: boolean;
  compact?: boolean;
}

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

export const WalletBalance = ({ showRefreshButton = true, compact = false }: WalletBalanceProps) => {
  const { data, isLoading, refetch } = useWalletBalance();

  if (isLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center">
          <LoadingSpinner />
        </div>
      </Card>
    );
  }

  if (!data) {
    return (
      <Card className="p-6">
        <p className="text-sm text-slate-600 text-center">No se pudo cargar el saldo</p>
      </Card>
    );
  }

  const balanceAmount = parseFloat(data.balance) || 0;
  const pendingAmount = parseFloat(data.pending_balance) || 0;
  const isFrozen = data.status === "frozen";

  return (
    <Card className={`${compact ? "p-4" : "p-6"} ${isFrozen ? "border-red-500" : ""}`}>
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-blue-100 rounded-lg">
            <Wallet className="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <p className="text-sm text-slate-600 font-medium">Saldo Disponible</p>
            <p className={`${compact ? "text-2xl" : "text-3xl"} font-bold text-slate-900`}>
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
            <span className="font-medium text-slate-900">{data.currency}</span>
          </div>
          <div className="flex items-center justify-between text-sm mt-2">
            <span className="text-slate-600">Estado</span>
            <span
              className={`font-medium ${
                data.status === "active"
                  ? "text-green-600"
                  : data.status === "frozen"
                  ? "text-red-600"
                  : "text-slate-600"
              }`}
            >
              {data.status === "active" ? "Activa" : data.status === "frozen" ? "Congelada" : "Cerrada"}
            </span>
          </div>
        </div>
      )}
    </Card>
  );
};
