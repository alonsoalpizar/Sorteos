import { useNavigate, useSearchParams } from "react-router-dom";
import { Clock, AlertCircle } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useUserMode } from "@/contexts/UserModeContext";
import { cn } from "@/lib/utils";

export const CreditVerifying = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { mode } = useUserMode();
  const isOrganizer = mode === "organizer";

  const amount = searchParams.get("amount") || "0";
  const reference = searchParams.get("reference") || "";

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
          {/* Verifying Icon */}
          <div className="flex justify-center mb-6">
            <div
              className={cn(
                "w-20 h-20 rounded-full flex items-center justify-center",
                isOrganizer ? "bg-amber-100" : "bg-blue-100"
              )}
            >
              <div className="relative">
                <Clock
                  className={cn(
                    "w-12 h-12",
                    isOrganizer ? "text-amber-600" : "text-blue-600"
                  )}
                />
                <div className="absolute -top-1 -right-1">
                  <LoadingSpinner />
                </div>
              </div>
            </div>
          </div>

          {/* Verifying Message */}
          <h1 className="text-2xl font-bold text-slate-900 mb-2">
            Pago en verificación
          </h1>
          <p className="text-slate-600 mb-6">
            Tu pago está siendo verificado por el procesador
          </p>

          {/* Amount Display */}
          {amount !== "0" && (
            <div className="bg-slate-50 rounded-lg p-6 mb-6">
              <p className="text-sm text-slate-600 mb-2">Monto</p>
              <p
                className={cn(
                  "text-3xl font-bold",
                  isOrganizer ? "text-teal-600" : "text-blue-600"
                )}
              >
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

          {/* Info Message */}
          <div
            className={cn(
              "border rounded-lg p-4 mb-6 text-left",
              isOrganizer
                ? "bg-amber-50 border-amber-200"
                : "bg-blue-50 border-blue-200"
            )}
          >
            <div className="flex items-start gap-3">
              <AlertCircle
                className={cn(
                  "w-5 h-5 flex-shrink-0 mt-0.5",
                  isOrganizer ? "text-amber-600" : "text-blue-600"
                )}
              />
              <div>
                <p
                  className={cn(
                    "text-sm font-medium mb-2",
                    isOrganizer ? "text-amber-900" : "text-blue-900"
                  )}
                >
                  ¿Qué significa esto?
                </p>
                <p
                  className={cn(
                    "text-sm",
                    isOrganizer ? "text-amber-800" : "text-blue-800"
                  )}
                >
                  El procesador de pagos está verificando tu transacción. Este
                  proceso puede tomar algunos minutos. Te notificaremos cuando
                  el pago sea confirmado y los créditos se agreguen a tu
                  billetera.
                </p>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="space-y-3">
            <Button
              onClick={() => navigate("/wallet")}
              className="w-full"
              size="lg"
            >
              Ver mi billetera
            </Button>
            <Button
              onClick={() => navigate("/explore")}
              variant="outline"
              className="w-full"
              size="lg"
            >
              Continuar explorando
            </Button>
          </div>

          {/* Note */}
          <p className="text-xs text-slate-500 mt-6">
            Puedes cerrar esta página. Recibirás una notificación cuando el
            pago sea procesado.
          </p>
        </div>
      </Card>
    </div>
  );
};
