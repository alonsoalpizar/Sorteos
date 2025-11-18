// Wallet types for the Sorteos platform

export interface Wallet {
  wallet_id: number;
  wallet_uuid: string;
  balance: string; // decimal as string
  pending_balance: string; // decimal as string
  currency: string; // "CRC"
  status: WalletStatus;
}

export type WalletStatus = 'active' | 'frozen' | 'closed';

export interface WalletTransaction {
  id: number;
  uuid: string;
  type: TransactionType;
  amount: string; // decimal as string
  status: TransactionStatus;
  balance_before: string;
  balance_after: string;
  reference_type: string | null;
  reference_id: number | null;
  notes: string | null;
  created_at: string; // ISO 8601
  completed_at: string | null; // ISO 8601
}

export type TransactionType =
  | 'deposit'
  | 'withdrawal'
  | 'purchase'
  | 'refund'
  | 'prize_claim'
  | 'settlement_payout'
  | 'adjustment';

export type TransactionStatus = 'pending' | 'completed' | 'failed' | 'reversed';

export interface RechargeOption {
  desired_credit: string; // ₡1,000, ₡5,000, etc.
  fixed_fee: string;
  processor_rate: string; // "0.03"
  processor_fee: string;
  platform_fee_rate: string; // "0.02"
  platform_fee: string;
  total_fees: string;
  charge_amount: string; // Total a cobrar al usuario
}

export interface RechargeOptionsResponse {
  options: RechargeOption[];
  currency: string;
  note: string;
}

export interface WalletBalanceResponse {
  success: boolean;
  data: Wallet;
}

export interface TransactionHistoryResponse {
  success: boolean;
  data: {
    transactions: WalletTransaction[];
    pagination: {
      total: number;
      limit: number;
      offset: number;
    };
  };
}

export interface AddFundsRequest {
  amount: string; // Crédito deseado (no el charge_amount)
  payment_method: 'card' | 'sinpe' | 'transfer';
}

export interface AddFundsResponse {
  success: boolean;
  message: string;
  data: {
    transaction_id: number;
    transaction_uuid: string;
    amount: string;
    status: TransactionStatus;
    payment_method: string;
    payment_url?: string; // URL del procesador de pagos (opcional, para fase 2)
    idempotency_key: string;
  };
}

// Helper para formatear montos en CRC
export const formatCRC = (amount: string | number): string => {
  const numAmount = typeof amount === 'string' ? parseFloat(amount) : amount;
  return new Intl.NumberFormat('es-CR', {
    style: 'currency',
    currency: 'CRC',
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(numAmount);
};

// Helper para parsear montos
export const parseAmount = (amount: string): number => {
  return parseFloat(amount) || 0;
};

// Helper para traducir tipos de transacción
export const translateTransactionType = (type: TransactionType): string => {
  const translations: Record<TransactionType, string> = {
    deposit: 'Recarga',
    withdrawal: 'Retiro',
    purchase: 'Compra de boletos',
    refund: 'Devolución',
    prize_claim: 'Premio reclamado',
    settlement_payout: 'Pago de liquidación',
    adjustment: 'Ajuste',
  };
  return translations[type] || type;
};

// Helper para traducir estados de transacción
export const translateTransactionStatus = (status: TransactionStatus): string => {
  const translations: Record<TransactionStatus, string> = {
    pending: 'Pendiente',
    completed: 'Completado',
    failed: 'Fallido',
    reversed: 'Revertido',
  };
  return translations[status] || status;
};

// Helper para obtener color del badge de estado
export const getStatusColor = (
  status: TransactionStatus
): 'default' | 'success' | 'error' | 'warning' => {
  const colors: Record<TransactionStatus, 'default' | 'success' | 'error' | 'warning'> = {
    pending: 'warning',
    completed: 'success',
    failed: 'error',
    reversed: 'default',
  };
  return colors[status] || 'default';
};
