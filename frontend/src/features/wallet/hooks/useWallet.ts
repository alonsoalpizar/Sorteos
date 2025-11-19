import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { walletApi } from "../api/walletApi";
import type { ListTransactionsInput, AddFundsInput } from "../types";
import { toast } from "sonner";

// Query Keys
export const walletKeys = {
  all: ["wallet"] as const,
  balance: () => [...walletKeys.all, "balance"] as const,
  transactions: (input: ListTransactionsInput = {}) => [...walletKeys.all, "transactions", input] as const,
  rechargeOptions: () => [...walletKeys.all, "recharge-options"] as const,
  earnings: () => [...walletKeys.all, "earnings"] as const,
};

// Get Balance
export function useWalletBalance(enabled = true) {
  return useQuery({
    queryKey: walletKeys.balance(),
    queryFn: () => walletApi.getBalance(),
    enabled,
    staleTime: 30000, // 30 segundos
    refetchInterval: 60000, // Refetch cada minuto
  });
}

// List Transactions
export function useWalletTransactions(input: ListTransactionsInput = {}, enabled = true) {
  return useQuery({
    queryKey: walletKeys.transactions(input),
    queryFn: () => walletApi.listTransactions(input),
    enabled,
    staleTime: 30000, // 30 segundos
  });
}

// Get Recharge Options (público, no requiere autenticación)
export function useRechargeOptions(enabled = true) {
  return useQuery({
    queryKey: walletKeys.rechargeOptions(),
    queryFn: () => walletApi.getRechargeOptions(),
    enabled,
    staleTime: 300000, // 5 minutos (las opciones no cambian frecuentemente)
  });
}

// Get Earnings
export function useEarnings(enabled = true) {
  return useQuery({
    queryKey: walletKeys.earnings(),
    queryFn: () => walletApi.getEarnings(),
    enabled,
    staleTime: 60000, // 1 minuto
  });
}

// Add Funds Mutation
export function useAddFunds() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: AddFundsInput) => walletApi.addFunds(input),
    onSuccess: (data) => {
      toast.success("Recarga iniciada", {
        description: `Se está procesando tu recarga de ₡${data.amount}`,
      });
      // Invalidar balance y transacciones para refrescar
      queryClient.invalidateQueries({ queryKey: walletKeys.balance() });
      queryClient.invalidateQueries({ queryKey: walletKeys.all });
    },
    onError: (error: any) => {
      toast.error("Error al procesar recarga", {
        description: error.response?.data?.error?.message || error.message,
      });
    },
  });
}

// Hook de conveniencia para verificar saldo suficiente
export function useHasSufficientBalance(requiredAmount: number): boolean {
  const { data } = useWalletBalance();
  if (!data) return false;
  const currentBalance = parseFloat(data.balance) || 0;
  return currentBalance >= requiredAmount;
}
