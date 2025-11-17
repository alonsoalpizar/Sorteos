import { useQuery } from '@tanstack/react-query';
import { Category, CategoryListResponse } from '@/types/category';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

// Fetch de categorías (público, no requiere autenticación)
const fetchCategories = async (): Promise<Category[]> => {
  const response = await fetch(`${API_BASE_URL}/categories`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error('Error al obtener categorías');
  }

  const data: CategoryListResponse = await response.json();
  return data.categories;
};

// Hook para obtener todas las categorías
export const useCategories = () => {
  return useQuery<Category[], Error>({
    queryKey: ['categories'],
    queryFn: fetchCategories,
    staleTime: 1000 * 60 * 30, // 30 minutos (las categorías cambian poco)
    gcTime: 1000 * 60 * 60, // 1 hora en caché
  });
};
