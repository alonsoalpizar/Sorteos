import { useQuery } from '@tanstack/react-query';
import { getEarnings } from '../../../api/wallet';

export interface RaffleEarning {
  raffle_id: number;
  raffle_uuid: string;
  title: string;
  draw_date: string;
  completed_at: string | null;
  total_revenue: string;
  platform_fee_percent: string;
  platform_fee_amount: string;
  net_amount: string;
  settlement_status: 'pending' | 'completed';
  settled_at: string | null;
}

export interface EarningsData {
  total_collected: string;
  platform_commission: string;
  net_earnings: string;
  completed_raffles: number;
  raffles: RaffleEarning[];
}

export const useEarnings = (limit = 0, offset = 0) => {
  return useQuery({
    queryKey: ['wallet', 'earnings', limit, offset],
    queryFn: async (): Promise<EarningsData> => {
      console.log('🔍 useEarnings - Fetching with limit:', limit, 'offset:', offset);
      const data = await getEarnings(limit, offset);
      console.log('🔍 useEarnings - Data received:', data);
      return data;
    },
    staleTime: 30000, // 30 seconds
    gcTime: 300000,   // 5 minutes
  });
};
