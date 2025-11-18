import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getWalletBalance, addFunds, generateIdempotencyKey } from '../../../api/wallet';
import type { AddFundsRequest } from '../../../types/wallet';

/**
 * Hook para gestionar el estado de la billetera del usuario
 */
export const useWallet = () => {
  const queryClient = useQueryClient();

  // Query para obtener el saldo
  const {
    data: walletData,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['wallet', 'balance'],
    queryFn: async () => {
      const response = await getWalletBalance();
      return response.data;
    },
    staleTime: 30000, // 30 segundos
    refetchInterval: 60000, // Refetch cada minuto
  });

  // Mutation para agregar fondos
  const addFundsMutation = useMutation({
    mutationFn: async (request: AddFundsRequest) => {
      const idempotencyKey = generateIdempotencyKey();
      return addFunds(request, idempotencyKey);
    },
    onSuccess: () => {
      // Invalidar el cache de balance para refetch automático
      queryClient.invalidateQueries({ queryKey: ['wallet', 'balance'] });
      queryClient.invalidateQueries({ queryKey: ['wallet', 'transactions'] });
    },
  });

  return {
    wallet: walletData,
    balance: walletData?.balance || '0',
    pendingBalance: walletData?.pending_balance || '0',
    currency: walletData?.currency || 'CRC',
    status: walletData?.status || 'active',
    isLoading,
    error,
    refetch,
    addFunds: addFundsMutation.mutate,
    isAddingFunds: addFundsMutation.isPending,
    addFundsError: addFundsMutation.error,
    addFundsSuccess: addFundsMutation.isSuccess,
    addFundsData: addFundsMutation.data,
  };
};

/**
 * Hook para verificar si el usuario tiene saldo suficiente
 */
export const useHasSufficientBalance = (requiredAmount: number): boolean => {
  const { balance } = useWallet();
  const currentBalance = parseFloat(balance) || 0;
  return currentBalance >= requiredAmount;
};
