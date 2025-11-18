import { useQuery } from '@tanstack/react-query';
import { getTransactionHistory } from '../../../api/wallet';
import { useState } from 'react';

/**
 * Hook para gestionar el historial de transacciones con paginación
 */
export const useTransactionHistory = (initialLimit: number = 20) => {
  const [limit] = useState(initialLimit);
  const [offset, setOffset] = useState(0);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['wallet', 'transactions', limit, offset],
    queryFn: () => getTransactionHistory(limit, offset),
    staleTime: 30000, // 30 segundos
  });

  const transactions = data?.data.transactions || [];
  const pagination = data?.data.pagination || { total: 0, limit, offset };

  const hasNextPage = offset + limit < pagination.total;
  const hasPreviousPage = offset > 0;

  const nextPage = () => {
    if (hasNextPage) {
      setOffset((prev) => prev + limit);
    }
  };

  const previousPage = () => {
    if (hasPreviousPage) {
      setOffset((prev) => Math.max(0, prev - limit));
    }
  };

  const goToPage = (page: number) => {
    const newOffset = page * limit;
    if (newOffset >= 0 && newOffset < pagination.total) {
      setOffset(newOffset);
    }
  };

  const currentPage = Math.floor(offset / limit);
  const totalPages = Math.ceil(pagination.total / limit);

  return {
    transactions,
    pagination,
    isLoading,
    error,
    refetch,
    hasNextPage,
    hasPreviousPage,
    nextPage,
    previousPage,
    goToPage,
    currentPage,
    totalPages,
  };
};
