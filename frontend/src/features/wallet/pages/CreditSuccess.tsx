import { useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { CheckCircle2 } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { useUserMode } from "@/contexts/UserModeContext";
import { cn } from "@/lib/utils";

export const CreditSuccess = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { mode } = useUserMode();
  const isOrganizer = mode === "organizer";

  const amount = searchParams.get("amount") || "0";
  const reference = searchParams.get("reference") || "";

  useEffect(() => {
    // Confetti animation would go here if library is available
  }, []);

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
          {/* Success Icon */}
          <div className="flex justify-center mb-6">
            <div
              className={cn(
                "w-20 h-20 rounded-full flex items-center justify-center",
                isOrganizer ? "bg-teal-100" : "bg-green-100"
              )}
            >
              <CheckCircle2
                className={cn(
                  "w-12 h-12",
                  isOrganizer ? "text-teal-600" : "text-green-600"
                )}
              />
            </div>
          </div>

          {/* Success Message */}
          <h1 className="text-2xl font-bold text-slate-900 mb-2">
            ¡Recarga exitosa!
          </h1>
          <p className="text-slate-600 mb-6">
            Tu pago ha sido procesado correctamente
          </p>

          {/* Amount Display */}
          <div className="bg-slate-50 rounded-lg p-6 mb-6">
            <p className="text-sm text-slate-600 mb-2">Crédito agregado</p>
            <p
              className={cn(
                "text-4xl font-bold",
                isOrganizer ? "text-teal-600" : "text-blue-600"
              )}
            >
              {formatCRC(amount)}
            </p>
          </div>

          {/* Reference */}
          {reference && (
            <div className="mb-6">
              <p className="text-xs text-slate-500">
                Referencia: <span className="font-mono">{reference}</span>
              </p>
            </div>
          )}

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
              Explorar sorteos
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
};
