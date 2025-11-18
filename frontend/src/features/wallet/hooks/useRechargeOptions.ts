import { useQuery } from '@tanstack/react-query';
import { getRechargeOptions } from '../../../api/wallet';

/**
 * Hook para obtener las opciones predefinidas de recarga
 * Este endpoint es público y se cachea por 5 minutos
 */
export const useRechargeOptions = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['wallet', 'recharge-options'],
    queryFn: getRechargeOptions,
    staleTime: 5 * 60 * 1000, // 5 minutos
    gcTime: 10 * 60 * 1000, // 10 minutos (antes era cacheTime)
  });

  return {
    options: data?.options || [],
    currency: data?.currency || 'CRC',
    note: data?.note || '',
    isLoading,
    error,
  };
};
