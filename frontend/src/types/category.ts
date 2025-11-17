// Tipos para categorías de sorteos

export interface Category {
  id: number;
  name: string;
  slug: string;
  icon: string;
  description?: string;
  display_order: number;
}

export interface CategoryListResponse {
  categories: Category[];
}
