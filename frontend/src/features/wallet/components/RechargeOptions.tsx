import { useState } from "react";
import { CreditCard, Info, ArrowRight } from "lucide-react";
import { useRechargeOptions, useAddFunds, usePurchaseCredits } from "../hooks/useWallet";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useUserMode } from "@/contexts/UserModeContext";
import { cn } from "@/lib/utils";

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

// Helper para generar idempotency key
function generateIdempotencyKey(): string {
  return `${Date.now()}-${Math.random().toString(36).substring(7)}`;
}

export const RechargeOptions = () => {
  const { data: optionsData, isLoading: optionsLoading } = useRechargeOptions();
  const addFundsMutation = useAddFunds();
  const purchaseCreditsMutation = usePurchaseCredits();
  const { mode } = useUserMode();
  const isOrganizer = mode === 'organizer';

  const [selectedOptionIndex, setSelectedOptionIndex] = useState<number | null>(null);
  const [paymentMethod, setPaymentMethod] = useState<"card" | "sinpe" | "transfer">("card");

  const handleRecharge = () => {
    if (selectedOptionIndex === null || !optionsData) return;

    const selectedOption = optionsData.options[selectedOptionIndex];

    // Si el método de pago es tarjeta, usar Pagadito
    if (paymentMethod === "card") {
      purchaseCreditsMutation.mutate({
        desired_credit: selectedOption.desired_credit,
        currency: optionsData.currency || "CRC",
        idempotency_key: generateIdempotencyKey(),
      });
    } else {
      // Para otros métodos de pago, usar el flujo anterior
      addFundsMutation.mutate({
        amount: selectedOption.desired_credit,
        payment_method: paymentMethod,
        idempotency_key: generateIdempotencyKey(),
      });
    }
  };

  if (optionsLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center">
          <LoadingSpinner />
        </div>
      </Card>
    );
  }

  if (!optionsData) {
    return (
      <Card className="p-6">
        <p className="text-sm text-slate-600 text-center">No se pudieron cargar las opciones de recarga</p>
      </Card>
    );
  }

  // Mostrar mensaje de éxito si la recarga fue creada
  if (addFundsMutation.isSuccess && addFundsMutation.data) {
    return (
      <Card className="p-6">
        <div className="bg-green-50 border border-green-200 rounded-lg p-6">
          <div>
            <h3 className="font-semibold text-green-900 mb-2">¡Transacción creada exitosamente!</h3>
            <p className="text-sm text-green-800 mb-4">Tu recarga está siendo procesada</p>
            <div className="bg-white p-4 rounded-lg border border-green-200 text-sm space-y-2">
              <p>
                <span className="font-medium">ID de transacción:</span> {addFundsMutation.data.transaction_uuid}
              </p>
              <p>
                <span className="font-medium">Monto:</span> {formatCRC(addFundsMutation.data.amount)}
              </p>
              <p>
                <span className="font-medium">Estado:</span> {addFundsMutation.data.status}
              </p>
            </div>
            <p className="text-xs text-green-700 mt-4">
              * En esta fase de desarrollo, el pago real aún no está habilitado. La transacción quedará pendiente.
            </p>
            <Button
              variant="outline"
              size="sm"
              className="mt-3"
              onClick={() => window.location.reload()}
            >
              Realizar otra recarga
            </Button>
          </div>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Info note */}
      {optionsData.note && (
        <div className={cn(
          "rounded-lg p-4 flex items-start gap-3 border",
          isOrganizer ? "bg-teal-50 border-teal-200" : "bg-blue-50 border-blue-200"
        )}>
          <Info className={cn("w-5 h-5 flex-shrink-0 mt-0.5", isOrganizer ? "text-teal-600" : "text-blue-600")} />
          <p className={cn("text-sm", isOrganizer ? "text-teal-900" : "text-blue-900")}>{optionsData.note}</p>
        </div>
      )}

      {/* Error alert */}
      {(addFundsMutation.isError || purchaseCreditsMutation.isError) && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-sm text-red-900">
            Error al crear la recarga:{" "}
            {addFundsMutation.error instanceof Error
              ? addFundsMutation.error.message
              : purchaseCreditsMutation.error instanceof Error
              ? purchaseCreditsMutation.error.message
              : "Error desconocido"}
          </p>
        </div>
      )}

      {/* Opciones de recarga */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {optionsData.options.map((option, index) => {
          const isSelected = selectedOptionIndex === index;
          const desiredCredit = parseFloat(option.desired_credit);
          const chargeAmount = parseFloat(option.charge_amount);
          const totalFees = parseFloat(option.total_fees);

          return (
            <Card
              key={index}
              className={cn(
                "p-5 cursor-pointer transition-all",
                isSelected
                  ? isOrganizer
                    ? "ring-2 ring-teal-500 bg-teal-50"
                    : "ring-2 ring-blue-500 bg-blue-50"
                  : "hover:shadow-md"
              )}
              onClick={() => setSelectedOptionIndex(index)}
            >
              <div className="text-center">
                {/* Crédito que recibirá */}
                <div className="mb-3">
                  <p className="text-sm text-slate-600 font-medium">Recibirás</p>
                  <p className={cn("text-3xl font-bold", isOrganizer ? "text-teal-600" : "text-blue-600")}>
                    {formatCRC(desiredCredit)}
                  </p>
                </div>

                {/* Divider */}
                <div className="border-t border-slate-200 my-3"></div>

                {/* Monto a pagar */}
                <div className="space-y-1">
                  <div className="flex justify-between text-sm">
                    <span className="text-slate-600">Total a pagar:</span>
                    <span className="font-semibold text-slate-900">{formatCRC(chargeAmount)}</span>
                  </div>
                  <div className="flex justify-between text-xs text-slate-500">
                    <span>Comisión por servicio:</span>
                    <span>{formatCRC(totalFees)}</span>
                  </div>
                </div>

                {/* Checkmark si está seleccionado */}
                {isSelected && (
                  <div className="mt-3">
                    <div className={cn(
                      "text-white rounded-full w-6 h-6 flex items-center justify-center mx-auto",
                      isOrganizer ? "bg-teal-600" : "bg-blue-600"
                    )}>
                      ✓
                    </div>
                  </div>
                )}
              </div>
            </Card>
          );
        })}
      </div>

      {/* Desglose simplificado de la opción seleccionada */}
      {selectedOptionIndex !== null && (
        <Card className="p-6 bg-slate-50">
          <h3 className="font-semibold text-slate-900 mb-4">Desglose</h3>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-slate-600">Crédito deseado:</span>
              <span className="font-medium">{formatCRC(optionsData.options[selectedOptionIndex].desired_credit)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-600">Comisión por servicio:</span>
              <span className="font-medium">{formatCRC(optionsData.options[selectedOptionIndex].total_fees)}</span>
            </div>
            <div className="border-t border-slate-300 my-2"></div>
            <div className="flex justify-between font-semibold text-base">
              <span>Total a pagar:</span>
              <span className={isOrganizer ? "text-teal-600" : "text-blue-600"}>
                {formatCRC(optionsData.options[selectedOptionIndex].charge_amount)}
              </span>
            </div>
          </div>
        </Card>
      )}

      {/* Métodos de pago */}
      {selectedOptionIndex !== null && (
        <Card className="p-6">
          <h3 className="font-semibold text-slate-900 mb-4">Método de pago</h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <button
              onClick={() => setPaymentMethod("card")}
              className={cn(
                "p-4 border-2 rounded-lg transition-all",
                paymentMethod === "card"
                  ? isOrganizer
                    ? "border-teal-500 bg-teal-50"
                    : "border-blue-500 bg-blue-50"
                  : "border-slate-200 hover:border-slate-300"
              )}
            >
              <CreditCard className={cn("w-6 h-6 mx-auto mb-2", isOrganizer ? "text-teal-600" : "text-blue-600")} />
              <p className="text-sm font-medium">Tarjeta</p>
            </button>
            <button
              onClick={() => setPaymentMethod("sinpe")}
              className={cn(
                "p-4 border-2 rounded-lg transition-all",
                paymentMethod === "sinpe"
                  ? isOrganizer
                    ? "border-teal-500 bg-teal-50"
                    : "border-blue-500 bg-blue-50"
                  : "border-slate-200 hover:border-slate-300"
              )}
            >
              <div className="w-6 h-6 mx-auto mb-2 text-2xl">💸</div>
              <p className="text-sm font-medium">SINPE Móvil</p>
            </button>
            <button
              onClick={() => setPaymentMethod("transfer")}
              className={cn(
                "p-4 border-2 rounded-lg transition-all",
                paymentMethod === "transfer"
                  ? isOrganizer
                    ? "border-teal-500 bg-teal-50"
                    : "border-blue-500 bg-blue-50"
                  : "border-slate-200 hover:border-slate-300"
              )}
            >
              <div className="w-6 h-6 mx-auto mb-2 text-2xl">🏦</div>
              <p className="text-sm font-medium">Transferencia</p>
            </button>
          </div>
        </Card>
      )}

      {/* Botón de confirmar */}
      {selectedOptionIndex !== null && (
        <Button
          onClick={handleRecharge}
          disabled={addFundsMutation.isPending || purchaseCreditsMutation.isPending}
          className="w-full"
          size="lg"
        >
          {(addFundsMutation.isPending || purchaseCreditsMutation.isPending) ? (
            <>
              <div className="w-5 h-5 mr-2 inline-block">
                <LoadingSpinner />
              </div>
              <span>Procesando...</span>
            </>
          ) : (
            <>
              Recargar {formatCRC(optionsData.options[selectedOptionIndex].desired_credit)}
              <ArrowRight className="w-5 h-5 ml-2" />
            </>
          )}
        </Button>
      )}
    </div>
  );
};
