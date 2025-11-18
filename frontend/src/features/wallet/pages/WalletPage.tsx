import { useState } from 'react';
import { Wallet, TrendingUp, History, DollarSign } from 'lucide-react';
import { WalletBalance } from '../components/WalletBalance';
import { RechargeOptions } from '../components/RechargeOptions';
import { TransactionHistory } from '../components/TransactionHistory';
import { Earnings } from '../components/Earnings';

type Tab = 'balance' | 'recharge' | 'history' | 'earnings';

export const WalletPage = () => {
  const [activeTab, setActiveTab] = useState<Tab>('balance');

  const tabs = [
    { id: 'balance' as Tab, label: 'Mi Saldo', icon: Wallet },
    { id: 'recharge' as Tab, label: 'Recargar', icon: TrendingUp },
    { id: 'history' as Tab, label: 'Historial', icon: History },
    { id: 'earnings' as Tab, label: 'Mis Ganancias', icon: DollarSign },
  ];

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

              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`
                    group inline-flex items-center py-4 px-1 border-b-2 font-medium text-sm
                    ${
                      isActive
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
                    }
                  `}
                >
                  <Icon
                    className={`
                      -ml-0.5 mr-2 h-5 w-5
                      ${isActive ? 'text-blue-600' : 'text-slate-400 group-hover:text-slate-500'}
                    `}
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
                  className="p-6 bg-white border-2 border-blue-500 rounded-lg hover:bg-blue-50 transition-colors text-left"
                >
                  <TrendingUp className="w-6 h-6 text-blue-600 mb-2" />
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
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h3 className="font-semibold text-blue-900 mb-2">💡 ¿Cómo funciona?</h3>
                <ul className="text-sm text-blue-800 space-y-1">
                  <li>• Recarga créditos una vez y úsalos para comprar boletos en todos los sorteos</li>
                  <li>• Sin comisiones adicionales al pagar con tu saldo</li>
                  <li>• Transacciones instantáneas y seguras</li>
                  <li>• Consulta tu historial completo en cualquier momento</li>
                </ul>
              </div>
            </div>
          )}

          {activeTab === 'recharge' && (
            <div>
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
