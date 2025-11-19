// ==================== WALLET MANAGEMENT ====================

export type WalletStatus = "active" | "frozen" | "closed";

// Wallet summary for list view
export interface WalletSummary {
  id: number;
  uuid: string;
  user_id: number;
  user_email: string;
  user_name: string;
  balance_available: string;
  earnings_balance: string;
  pending_balance: string;
  total_balance: string;
  currency: string;
  status: WalletStatus;
  created_at: string;
  updated_at: string;
}

// Wallet details (full information)
export interface WalletDetails {
  id: number;
  uuid: string;
  user_id: number;
  user_email: string;
  user_name: string;
  balance_available: string;
  earnings_balance: string;
  pending_balance: string;
  total_balance: string;
  currency: string;
  status: WalletStatus;
  created_at: string;
  updated_at: string;
}

// Transaction types (from user wallet)
export type TransactionType =
  | "deposit"
  | "withdrawal"
  | "purchase"
  | "refund"
  | "prize_claim"
  | "settlement_payout"
  | "adjustment";

export type TransactionStatus = "pending" | "completed" | "failed" | "reversed";

// Transaction summary for list
export interface TransactionSummary {
  id: number;
  uuid: string;
  wallet_id: number;
  type: TransactionType;
  amount: string;
  status: TransactionStatus;
  balance_before: string;
  balance_after: string;
  description: string;
  metadata: Record<string, any> | null;
  created_at: string;
  updated_at: string;
}

// ==================== API INPUTS ====================

export interface ListWalletsInput {
  page?: number;
  limit?: number;
  status?: WalletStatus;
  user_email?: string;
}

export interface ListWalletTransactionsInput {
  wallet_id: number;
  page?: number;
  limit?: number;
  type?: TransactionType;
  status?: TransactionStatus;
}

export interface FreezeWalletInput {
  wallet_id: number;
  reason: string;
}

export interface UnfreezeWalletInput {
  wallet_id: number;
}

// ==================== API OUTPUTS ====================

export interface ListWalletsOutput {
  wallets: WalletSummary[];
  pagination: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

export interface ViewWalletDetailsOutput {
  wallet: WalletDetails;
}

export interface ListWalletTransactionsOutput {
  transactions: TransactionSummary[];
  pagination: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

export interface FreezeWalletOutput {
  success: boolean;
  message: string;
}

export interface UnfreezeWalletOutput {
  success: boolean;
  message: string;
}
