// ============================================================================
// WALLET TYPES
// ============================================================================

// Transaction Types
export type TransactionType =
  | "deposit"            // Compra de créditos vía procesador
  | "withdrawal"         // Retiro a cuenta bancaria
  | "purchase"           // Pago de sorteo
  | "refund"             // Devolución de compra
  | "prize_claim"        // Premio ganado
  | "settlement_payout"  // Pago de liquidación a organizador
  | "adjustment";        // Ajuste manual (admin)

export type TransactionStatus = "pending" | "completed" | "failed" | "reversed";

// Wallet Balance
export interface WalletBalance {
  wallet_id: number;
  wallet_uuid: string;
  balance: string;
  pending_balance: string;
  currency: string;
  status: string;
}

export interface GetBalanceOutput {
  wallet_id: number;
  wallet_uuid: string;
  balance: string;
  pending_balance: string;
  currency: string;
  status: string;
}

// Wallet Transaction
export interface WalletTransaction {
  id: number;
  uuid: string;
  type: TransactionType;
  amount: string;
  status: TransactionStatus;
  balance_before: string;
  balance_after: string;
  reference_type: string | null;
  reference_id: number | null;
  notes: string | null;
  created_at: string;
  completed_at: string | null;
}

export interface ListTransactionsInput {
  limit?: number;
  offset?: number;
}

export interface ListTransactionsOutput {
  transactions: WalletTransaction[];
  pagination: {
    total: number;
    limit: number;
    offset: number;
  };
}

// Recharge Options
export interface RechargeBreakdown {
  desired_credit: string;   // Crédito que el usuario recibirá
  fixed_fee: string;        // Tarifa fija del procesador (₡200)
  processor_rate: string;   // Tasa porcentual del procesador (0.0425 = 4.25%)
  processor_fee: string;    // Comisión calculada del procesador
  platform_fee_rate: string; // Tasa de la plataforma (0.0 = sin comisión adicional)
  platform_fee: string;     // Comisión de la plataforma
  total_fees: string;       // Total de comisiones
  charge_amount: string;    // Monto total a cobrar al usuario
}

export interface GetRechargeOptionsOutput {
  options: RechargeBreakdown[];
  currency: string;
  note: string;
}

// Earnings
export interface RaffleEarning {
  raffle_id: number;
  raffle_uuid: string;
  title: string;
  draw_date: string;
  completed_at: string | null;
  total_revenue: string;        // Ingresos totales
  platform_fee_percent: string; // Porcentaje de comisión de plataforma
  platform_fee_amount: string;  // Monto de comisión de plataforma
  net_amount: string;           // Monto neto para el organizador
  settlement_status: string;    // pending, completed, etc.
  settled_at: string | null;
}

export interface GetEarningsOutput {
  total_collected: string;     // Total recaudado
  platform_commission: string; // Comisión total de plataforma
  net_earnings: string;        // Ganancias netas
  completed_raffles: number;   // Cantidad de sorteos completados
  raffles: RaffleEarning[];    // Lista de sorteos
}

// Add Funds
export interface AddFundsInput {
  amount: string;
  payment_method: string;
  idempotency_key: string;
}

export interface AddFundsOutput {
  transaction_id: number;
  transaction_uuid: string;
  amount: string;
  status: TransactionStatus;
  created_at: string;
}

// Purchase Credits (Pagadito)
export interface PurchaseCreditsInput {
  desired_credit: string;
  currency: string;
  idempotency_key: string;
}

export interface PurchaseCreditsOutput {
  purchase_id: number;
  purchase_uuid: string;
  payment_url: string;
  ern: string;
  desired_credit: string;
  charge_amount: string;
  expires_at: string;
}

// Credit Purchase Status
export type CreditPurchaseStatus = "pending" | "processing" | "completed" | "failed" | "expired";

export interface CreditPurchase {
  id: number;
  uuid: string;
  desired_credit: string;
  charge_amount: string;
  currency: string;
  status: CreditPurchaseStatus;
  pagadito_reference: string | null;
  created_at: string;
  completed_at: string | null;
}
