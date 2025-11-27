import { useNavigate, useSearchParams } from "react-router-dom";
import { XCircle, RefreshCw } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";

export const CreditFailed = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const amount = searchParams.get("amount") || "0";
  const reference = searchParams.get("reference") || "";
  const reason = searchParams.get("reason") || "El pago no pudo ser procesado";

  const formatCRC = (amount: string): string => {
    const num = parseFloat(amount);
    return new Intl.NumberFormat("es-CR", {
      style: "currency",
      currency: "CRC",
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    }).format(num);
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-50 to-white flex items-center justify-center p-4">
      <Card className="max-w-md w-full p-8">
        <div className="text-center">
          {/* Error Icon */}
          <div className="flex justify-center mb-6">
            <div className="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center">
              <XCircle className="w-12 h-12 text-red-600" />
            </div>
          </div>

          {/* Error Message */}
          <h1 className="text-2xl font-bold text-slate-900 mb-2">
            Pago rechazado
          </h1>
          <p className="text-slate-600 mb-6">{reason}</p>

          {/* Amount Display */}
          {amount !== "0" && (
            <div className="bg-slate-50 rounded-lg p-6 mb-6">
              <p className="text-sm text-slate-600 mb-2">Monto intentado</p>
              <p className="text-3xl font-bold text-slate-900">
                {formatCRC(amount)}
              </p>
            </div>
          )}

          {/* Reference */}
          {reference && (
            <div className="mb-6">
              <p className="text-xs text-slate-500">
                Referencia: <span className="font-mono">{reference}</span>
              </p>
            </div>
          )}

          {/* Help Message */}
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6 text-left">
            <p className="text-sm text-blue-900 font-medium mb-2">
              ¿Qué puedes hacer?
            </p>
            <ul className="text-sm text-blue-800 space-y-1 list-disc list-inside">
              <li>Verifica que tu método de pago tenga fondos suficientes</li>
              <li>Asegúrate de que los datos ingresados sean correctos</li>
              <li>Intenta con otro método de pago</li>
              <li>
                Contacta a tu banco si el problema persiste
              </li>
            </ul>
          </div>

          {/* Actions */}
          <div className="space-y-3">
            <Button
              onClick={() => navigate("/wallet")}
              className="w-full"
              size="lg"
            >
              <RefreshCw className="w-5 h-5 mr-2" />
              Intentar nuevamente
            </Button>
            <Button
              onClick={() => navigate("/explore")}
              variant="outline"
              className="w-full"
              size="lg"
            >
              Volver a explorar
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
};
