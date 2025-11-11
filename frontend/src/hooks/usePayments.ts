import { useMutation, useQuery } from '@tanstack/react-query';
import { apiClient } from '../lib/apiClient';

export interface CreatePaymentIntentInput {
  reservation_id: string;
  return_url?: string;
  cancel_url?: string;
}

export interface PaymentIntent {
  id: string;
  amount: number;
  currency: string;
  status: string;
  client_secret: string; // For Stripe or approval_url for PayPal
  approval_url?: string; // PayPal specific
  metadata: Record<string, string>;
}

export interface Payment {
  id: string;
  reservation_id: string;
  user_id: string;
  raffle_id: string;
  amount: number;
  currency: string;
  status: 'pending' | 'succeeded' | 'failed' | 'cancelled';
  payment_intent_id: string;
  client_secret: string;
  metadata: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface CreatePaymentIntentResponse {
  payment: Payment;
  payment_intent: PaymentIntent;
}

export function useCreatePaymentIntent() {
  return useMutation({
    mutationFn: async (input: CreatePaymentIntentInput): Promise<CreatePaymentIntentResponse> => {
      const { data } = await apiClient.post<CreatePaymentIntentResponse>('/payments/intent', {
        reservation_id: input.reservation_id,
        return_url: input.return_url || window.location.origin + '/payment/success',
        cancel_url: input.cancel_url || window.location.origin + '/payment/cancel',
      });

      return data;
    },
  });
}

export function useGetPayment(paymentId: string | null) {
  return useQuery({
    queryKey: ['payment', paymentId],
    queryFn: async (): Promise<Payment> => {
      if (!paymentId) throw new Error('No payment ID provided');

      const { data } = await apiClient.get<{ payment: Payment }>(`/payments/${paymentId}`);
      return data.payment;
    },
    enabled: !!paymentId,
  });
}

export function useGetMyPayments() {
  return useQuery({
    queryKey: ['myPayments'],
    queryFn: async (): Promise<Payment[]> => {
      const { data } = await apiClient.get<{ payments: Payment[] }>('/payments/me');
      return data.payments;
    },
  });
}
