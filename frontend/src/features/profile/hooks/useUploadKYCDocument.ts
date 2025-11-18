import { useMutation, useQueryClient } from '@tanstack/react-query';
import { uploadKYCDocument } from '../../../api/profile';

/**
 * Hook para subir documentos KYC
 */
export const useUploadKYCDocument = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      documentType,
      fileUrl,
    }: {
      documentType: 'cedula_front' | 'cedula_back' | 'selfie';
      fileUrl: string;
    }) => uploadKYCDocument(documentType, fileUrl),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};
