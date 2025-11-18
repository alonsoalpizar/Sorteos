import { useQuery } from '@tanstack/react-query';
import { getProfile } from '../../../api/profile';
import type { ProfileData } from '../types/profile';

/**
 * Hook para obtener el perfil completo del usuario
 * Retorna: user data, KYC documents, wallet info, can_withdraw status
 */
export const useProfile = () => {
  return useQuery<ProfileData>({
    queryKey: ['profile'],
    queryFn: getProfile,
    staleTime: 60000, // 1 minute
    gcTime: 300000,   // 5 minutes
  });
};
