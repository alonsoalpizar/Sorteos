import { useMutation, useQueryClient } from '@tanstack/react-query';
import { configureIBAN } from '../../../api/profile';

/**
 * Hook para configurar el IBAN del usuario
 * Requiere kyc_level >= cedula_verified
 */
export const useConfigureIBAN = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (iban: string) => configureIBAN(iban),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};
