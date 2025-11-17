import { useMutation, useQueryClient } from '@tanstack/react-query';
import { imagesApi } from '../api/images';
import { raffleKeys } from './useRaffles';

/**
 * Hook para subir imagen a un sorteo
 */
export function useUploadImage() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ raffleId, file }: { raffleId: number; file: File }) =>
      imagesApi.upload(raffleId, file),
    onSuccess: async (response, variables) => {
      console.log('Upload success - invalidating queries for raffle:', variables.raffleId);
      console.log('Upload response:', response);

      // Invalidar el query detail para forzar refetch
      await queryClient.invalidateQueries({
        queryKey: raffleKeys.detail(variables.raffleId),
      });

      // También forzar un refetch inmediato
      await queryClient.refetchQueries({
        queryKey: raffleKeys.detail(variables.raffleId),
        type: 'active',
      });

      // Debug: verificar el estado del cache después del refetch
      const queryData = queryClient.getQueryData(raffleKeys.detail(variables.raffleId));
      console.log('Query data after refetch:', queryData);

      console.log('Queries invalidated and refetched');
    },
  });
}

/**
 * Hook para eliminar imagen de un sorteo
 */
export function useDeleteImage() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ raffleId, imageId }: { raffleId: number; imageId: number }) =>
      imagesApi.delete(raffleId, imageId),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: raffleKeys.detail(variables.raffleId),
        refetchType: 'active',
      });
    },
  });
}

/**
 * Hook para establecer imagen como primaria
 */
export function useSetPrimaryImage() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ raffleId, imageId }: { raffleId: number; imageId: number }) =>
      imagesApi.setPrimary(raffleId, imageId),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: raffleKeys.detail(variables.raffleId),
        refetchType: 'active',
      });
    },
  });
}
