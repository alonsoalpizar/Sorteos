import { DollarSign } from 'lucide-react';
import { Earnings } from '../components/Earnings';

// ===========================================
// NOTA: Tabs desactivados temporalmente
// ===========================================
// Sorteos.club ahora es solo plataforma de gestión sin monetización directa.
// Los imports y funcionalidad originales están comentados abajo por si se
// necesitan reactivar en el futuro:
//
// import { useState, useMemo } from 'react';
// import { Wallet, TrendingUp, History } from 'lucide-react';
// import { WalletBalance } from '../components/WalletBalance';
// import { RechargeOptions } from '../components/RechargeOptions';
// import { TransactionHistory } from '../components/TransactionHistory';
// import { useUserMode } from '@/contexts/UserModeContext';
// import { cn } from '@/lib/utils';
//
// Tabs desactivados:
// - 'balance' (Mi Saldo)
// - 'recharge' (Recargar)
// - 'history' (Historial)
// ===========================================

export const WalletPage = () => {
  return (
    <div className="min-h-screen bg-slate-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <DollarSign className="w-8 h-8 text-teal-600" />
            <h1 className="text-3xl font-bold text-slate-900">Mis Ganancias</h1>
          </div>
          <p className="text-slate-600">
            Consulta las ventas de tus sorteos completados
          </p>
        </div>

        {/* Contenido - Solo Earnings (sin comisión de plataforma) */}
        <div className="mt-6">
          <Earnings />
        </div>
      </div>
    </div>
  );
};
