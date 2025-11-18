import { useState } from 'react';
import { CreditCard, Info, ArrowRight } from 'lucide-react';
import { useRechargeOptions } from '../hooks/useRechargeOptions';
import { useWallet } from '../hooks/useWallet';
import { formatCRC, parseAmount } from '../../../types/wallet';
import { Card } from '../../../components/ui/Card';
import { Button } from '../../../components/ui/Button';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { Alert } from '../../../components/ui/Alert';

export const RechargeOptions = () => {
  const { options, note, isLoading: optionsLoading } = useRechargeOptions();
  const { addFunds, isAddingFunds, addFundsError, addFundsSuccess, addFundsData } = useWallet();
  const [selectedOptionIndex, setSelectedOptionIndex] = useState<number | null>(null);
  const [paymentMethod, setPaymentMethod] = useState<'card' | 'sinpe' | 'transfer'>('card');

  const handleRecharge = () => {
    if (selectedOptionIndex === null) return;

    const selectedOption = options[selectedOptionIndex];
    addFunds({
      amount: selectedOption.desired_credit,
      payment_method: paymentMethod,
    });
  };

  if (optionsLoading) {
    return (
      <Card className="p-6">
        <div className="flex items-center justify-center">
          <LoadingSpinner size="md" />
        </div>
      </Card>
    );
  }

  // Mostrar mensaje de éxito si la recarga fue creada
  if (addFundsSuccess && addFundsData) {
    return (
      <Card className="p-6">
        <Alert variant="success">
          <div>
            <h3 className="font-semibold mb-2">¡Transacción creada exitosamente!</h3>
            <p className="text-sm mb-4">{addFundsData.message}</p>
            <div className="bg-white p-4 rounded-lg border border-green-200 text-sm space-y-2">
              <p>
                <span className="font-medium">ID de transacción:</span> {addFundsData.data.transaction_uuid}
              </p>
              <p>
                <span className="font-medium">Monto:</span> {formatCRC(addFundsData.data.amount)}
              </p>
              <p>
                <span className="font-medium">Estado:</span> {addFundsData.data.status}
              </p>
            </div>
            {addFundsData.data.payment_url && (
              <Button
                className="mt-4"
                onClick={() => window.location.href = addFundsData.data.payment_url!}
              >
                Continuar al pago
                <ArrowRight className="w-4 h-4 ml-2" />
              </Button>
            )}
            <p className="text-xs text-slate-600 mt-4">
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
        </Alert>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Info note */}
      {note && (
        <Alert variant="info">
          <Info className="w-4 h-4" />
          <p className="text-sm">{note}</p>
        </Alert>
      )}

      {/* Error alert */}
      {addFundsError && (
        <Alert variant="error">
          <p className="text-sm">
            Error al crear la recarga:{' '}
            {addFundsError instanceof Error ? addFundsError.message : 'Error desconocido'}
          </p>
        </Alert>
      )}

      {/* Opciones de recarga */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {options.map((option: any, index: number) => {
          const isSelected = selectedOptionIndex === index;
          const desiredCredit = parseAmount(option.desired_credit);
          const chargeAmount = parseAmount(option.charge_amount);
          const totalFees = parseAmount(option.total_fees);

          return (
            <Card
              key={index}
              className={`p-5 cursor-pointer transition-all ${
                isSelected ? 'ring-2 ring-blue-500 bg-blue-50' : 'hover:shadow-md'
              }`}
              onClick={() => setSelectedOptionIndex(index)}
            >
              <div className="text-center">
                {/* Crédito que recibirá */}
                <div className="mb-3">
                  <p className="text-sm text-slate-600 font-medium">Recibirás</p>
                  <p className="text-3xl font-bold text-blue-600">{formatCRC(desiredCredit)}</p>
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
                    <div className="bg-blue-600 text-white rounded-full w-6 h-6 flex items-center justify-center mx-auto">
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
              <span className="font-medium">{formatCRC(options[selectedOptionIndex].desired_credit)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-600">Comisión por servicio:</span>
              <span className="font-medium">{formatCRC(options[selectedOptionIndex].total_fees)}</span>
            </div>
            <div className="border-t border-slate-300 my-2"></div>
            <div className="flex justify-between font-semibold text-base">
              <span>Total a pagar:</span>
              <span className="text-blue-600">{formatCRC(options[selectedOptionIndex].charge_amount)}</span>
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
              onClick={() => setPaymentMethod('card')}
              className={`p-4 border-2 rounded-lg transition-all ${
                paymentMethod === 'card'
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-slate-200 hover:border-slate-300'
              }`}
            >
              <CreditCard className="w-6 h-6 mx-auto mb-2 text-blue-600" />
              <p className="text-sm font-medium">Tarjeta</p>
            </button>
            <button
              onClick={() => setPaymentMethod('sinpe')}
              className={`p-4 border-2 rounded-lg transition-all ${
                paymentMethod === 'sinpe'
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-slate-200 hover:border-slate-300'
              }`}
            >
              <div className="w-6 h-6 mx-auto mb-2 text-2xl">💸</div>
              <p className="text-sm font-medium">SINPE Móvil</p>
            </button>
            <button
              onClick={() => setPaymentMethod('transfer')}
              className={`p-4 border-2 rounded-lg transition-all ${
                paymentMethod === 'transfer'
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-slate-200 hover:border-slate-300'
              }`}
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
          disabled={isAddingFunds}
          className="w-full"
          size="lg"
        >
          {isAddingFunds ? (
            <>
              <LoadingSpinner size="sm" />
              <span className="ml-2">Procesando...</span>
            </>
          ) : (
            <>
              Recargar {formatCRC(options[selectedOptionIndex].desired_credit)}
              <ArrowRight className="w-5 h-5 ml-2" />
            </>
          )}
        </Button>
      )}
    </div>
  );
};
