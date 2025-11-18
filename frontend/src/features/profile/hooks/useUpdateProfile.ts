import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateProfile } from '../../../api/profile';
import type { UpdateProfileRequest } from '../types/profile';

/**
 * Hook para actualizar la información personal del usuario
 */
export const useUpdateProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => updateProfile(data),
    onSuccess: () => {
      // Invalidar el cache del perfil para refrescar los datos
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};
