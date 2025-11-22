import { useState, useMemo } from 'react';
import { Wallet, TrendingUp, History, DollarSign } from 'lucide-react';
import { WalletBalance } from '../components/WalletBalance';
import { RechargeOptions } from '../components/RechargeOptions';
import { TransactionHistory } from '../components/TransactionHistory';
import { Earnings } from '../components/Earnings';
import { useUserMode } from '@/contexts/UserModeContext';
import { cn } from '@/lib/utils';

type Tab = 'balance' | 'recharge' | 'history' | 'earnings';

export const WalletPage = () => {
  const { mode } = useUserMode();
  const isOrganizer = mode === 'organizer';

  // Tab default según modo: Organizador = Ganancias, Participante = Recargar
  const defaultTab: Tab = isOrganizer ? 'earnings' : 'recharge';
  const [activeTab, setActiveTab] = useState<Tab>(defaultTab);

  // Tabs ordenados según modo
  const tabs = useMemo(() => {
    if (isOrganizer) {
      // Organizador: Ganancias primero, luego Saldo, Historial, Recargar
      return [
        { id: 'earnings' as Tab, label: 'Mis Ganancias', icon: DollarSign },
        { id: 'balance' as Tab, label: 'Mi Saldo', icon: Wallet },
        { id: 'history' as Tab, label: 'Historial', icon: History },
        { id: 'recharge' as Tab, label: 'Recargar', icon: TrendingUp },
      ];
    } else {
      // Participante: Recargar primero, luego Saldo, Historial, Ganancias
      return [
        { id: 'recharge' as Tab, label: 'Recargar', icon: TrendingUp },
        { id: 'balance' as Tab, label: 'Mi Saldo', icon: Wallet },
        { id: 'history' as Tab, label: 'Historial', icon: History },
        { id: 'earnings' as Tab, label: 'Mis Ganancias', icon: DollarSign },
      ];
    }
  }, [isOrganizer]);

  return (
    <div className="min-h-screen bg-slate-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900 mb-2">Mi Billetera</h1>
          <p className="text-slate-600">
            Gestiona tu saldo, recarga créditos y revisa tu historial de transacciones
          </p>
        </div>

        {/* Tabs */}
        <div className="mb-6 border-b border-slate-200">
          <nav className="-mb-px flex space-x-8">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;

              // Colores dinámicos según modo
              const activeColor = isOrganizer ? 'border-teal-500 text-teal-600' : 'border-blue-500 text-blue-600';
              const activeIconColor = isOrganizer ? 'text-teal-600' : 'text-blue-600';

              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={cn(
                    "group inline-flex items-center py-4 px-1 border-b-2 font-medium text-sm transition-colors",
                    isActive
                      ? activeColor
                      : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
                  )}
                >
                  <Icon
                    className={cn(
                      "-ml-0.5 mr-2 h-5 w-5 transition-colors",
                      isActive ? activeIconColor : 'text-slate-400 group-hover:text-slate-500'
                    )}
                  />
                  {tab.label}
                </button>
              );
            })}
          </nav>
        </div>

        {/* Tab Content */}
        <div className="mt-6">
          {activeTab === 'balance' && (
            <div className="space-y-6">
              <WalletBalance showRefreshButton={true} compact={false} />

              {/* Quick actions */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <button
                  onClick={() => setActiveTab('recharge')}
                  className={cn(
                    "p-6 bg-white border-2 rounded-lg transition-colors text-left",
                    isOrganizer
                      ? "border-teal-500 hover:bg-teal-50"
                      : "border-blue-500 hover:bg-blue-50"
                  )}
                >
                  <TrendingUp className={cn("w-6 h-6 mb-2", isOrganizer ? "text-teal-600" : "text-blue-600")} />
                  <h3 className="font-semibold text-slate-900 mb-1">Recargar saldo</h3>
                  <p className="text-sm text-slate-600">
                    Agrega créditos a tu billetera con tus métodos de pago preferidos
                  </p>
                </button>

                <button
                  onClick={() => setActiveTab('history')}
                  className="p-6 bg-white border-2 border-slate-200 rounded-lg hover:bg-slate-50 transition-colors text-left"
                >
                  <History className="w-6 h-6 text-slate-600 mb-2" />
                  <h3 className="font-semibold text-slate-900 mb-1">Ver historial</h3>
                  <p className="text-sm text-slate-600">
                    Revisa todas tus transacciones y movimientos de saldo
                  </p>
                </button>
              </div>

              {/* Info sobre el uso de la billetera */}
              <div className={cn(
                "rounded-lg p-4 border",
                isOrganizer
                  ? "bg-teal-50 border-teal-200"
                  : "bg-blue-50 border-blue-200"
              )}>
                <h3 className={cn("font-semibold mb-2", isOrganizer ? "text-teal-900" : "text-blue-900")}>
                  {isOrganizer ? "💰 Tu billetera de organizador" : "💡 ¿Cómo funciona?"}
                </h3>
                <ul className={cn("text-sm space-y-1", isOrganizer ? "text-teal-800" : "text-blue-800")}>
                  {isOrganizer ? (
                    <>
                      <li>• Recibe automáticamente las ganancias de tus sorteos completados</li>
                      <li>• Solicita retiros cuando lo necesites</li>
                      <li>• Consulta el desglose de comisiones y ganancias netas</li>
                    </>
                  ) : (
                    <>
                      <li>• Recarga créditos una vez y úsalos para comprar boletos en todos los sorteos</li>
                      <li>• Sin comisiones adicionales al pagar con tu saldo</li>
                      <li>• Transacciones instantáneas y seguras</li>
                      <li>• Consulta tu historial completo en cualquier momento</li>
                    </>
                  )}
                </ul>
              </div>
            </div>
          )}

          {activeTab === 'recharge' && (
            <div className="space-y-6">
              {/* Saldo compacto arriba de las opciones de recarga */}
              <WalletBalance showRefreshButton={true} compact={true} />
              <RechargeOptions />
            </div>
          )}

          {activeTab === 'history' && (
            <div>
              <TransactionHistory />
            </div>
          )}

          {activeTab === 'earnings' && (
            <div>
              <Earnings />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
