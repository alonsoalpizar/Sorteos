import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { adminWalletsApi } from "../api/walletsApi";
import type {
  ListWalletsInput,
  ListWalletTransactionsInput,
  FreezeWalletInput,
  UnfreezeWalletInput,
} from "../types";
import { toast } from "sonner";

// ==================== QUERY KEYS ====================

export const adminWalletKeys = {
  all: ["admin", "wallets"] as const,
  lists: () => [...adminWalletKeys.all, "list"] as const,
  list: (input: ListWalletsInput) => [...adminWalletKeys.lists(), input] as const,
  details: () => [...adminWalletKeys.all, "detail"] as const,
  detail: (id: number) => [...adminWalletKeys.details(), id] as const,
  transactions: (walletId: number, input: ListWalletTransactionsInput) =>
    [...adminWalletKeys.all, "transactions", walletId, input] as const,
};

// ==================== QUERIES ====================

// List all wallets with pagination and filters
export function useAdminWallets(input: ListWalletsInput = {}, enabled = true) {
  return useQuery({
    queryKey: adminWalletKeys.list(input),
    queryFn: () => adminWalletsApi.listWallets(input),
    enabled,
    staleTime: 30000, // 30 seconds
  });
}

// View wallet details
export function useAdminWalletDetails(walletId: number, enabled = true) {
  return useQuery({
    queryKey: adminWalletKeys.detail(walletId),
    queryFn: () => adminWalletsApi.viewWalletDetails(walletId),
    enabled,
    staleTime: 30000, // 30 seconds
  });
}

// List wallet transactions
export function useAdminWalletTransactions(
  input: ListWalletTransactionsInput,
  enabled = true
) {
  return useQuery({
    queryKey: adminWalletKeys.transactions(input.wallet_id, input),
    queryFn: () => adminWalletsApi.listWalletTransactions(input),
    enabled,
    staleTime: 30000, // 30 seconds
  });
}

// ==================== MUTATIONS ====================

// Freeze wallet
export function useFreezeWallet() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: FreezeWalletInput) => adminWalletsApi.freezeWallet(input),
    onSuccess: (data, variables) => {
      toast.success("Billetera congelada", {
        description: data.message,
      });
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: adminWalletKeys.all });
      queryClient.invalidateQueries({ queryKey: adminWalletKeys.detail(variables.wallet_id) });
    },
    onError: (error: any) => {
      toast.error("Error al congelar billetera", {
        description: error.response?.data?.error?.message || error.message,
      });
    },
  });
}

// Unfreeze wallet
export function useUnfreezeWallet() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: UnfreezeWalletInput) => adminWalletsApi.unfreezeWallet(input),
    onSuccess: (data, variables) => {
      toast.success("Billetera descongelada", {
        description: data.message,
      });
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: adminWalletKeys.all });
      queryClient.invalidateQueries({ queryKey: adminWalletKeys.detail(variables.wallet_id) });
    },
    onError: (error: any) => {
      toast.error("Error al descongelar billetera", {
        description: error.response?.data?.error?.message || error.message,
      });
    },
  });
}
