// Tipos para gestión de sorteos

export type RaffleStatus = 'draft' | 'active' | 'suspended' | 'completed' | 'cancelled';
export type DrawMethod = 'loteria_nacional_cr' | 'manual' | 'random';
export type SettlementStatus = 'pending' | 'processing' | 'completed' | 'failed';
export type RaffleNumberStatus = 'available' | 'reserved' | 'sold';

export interface Raffle {
  id: number;
  uuid: string;
  user_id: number;
  title: string;
  description: string;
  status: RaffleStatus;
  price_per_number: string;
  total_numbers: number;
  draw_date: string;
  draw_method: DrawMethod;
  sold_count: number;
  reserved_count: number;
  total_revenue: string;
  platform_fee_percentage: string;
  platform_fee_amount: string;
  net_amount: string;
  settlement_status: SettlementStatus;
  winner_number?: string;
  winner_user_id?: number;
  published_at?: string;
  completed_at?: string;
  cancelled_at?: string;
  created_at: string;
  updated_at?: string;
}

export interface RaffleNumber {
  id: number;
  raffle_id: number;
  number: string;
  status: RaffleNumberStatus;
  user_id?: number;
  reserved_by?: number;
  reserved_until?: string;
  reservation_id?: number;
  purchased_at?: string;
}

export interface RaffleImage {
  id: number;
  raffle_id: number;
  filename: string;
  file_size: number;
  mime_type: string;
  width?: number;
  height?: number;
  alt_text?: string;
  display_order: number;
  is_primary: boolean;
}

export interface RaffleDetail {
  raffle: Raffle;
  numbers?: RaffleNumber[];
  images?: RaffleImage[];
  available_count: number;
  reserved_count: number;
  sold_count: number;
}

export interface RaffleListResponse {
  raffles: Raffle[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

export interface CreateRaffleInput {
  title: string;
  description: string;
  price_per_number: number;
  total_numbers: number;
  draw_date: string;
  draw_method: DrawMethod;
}

export interface UpdateRaffleInput {
  title?: string;
  description?: string;
  draw_date?: string;
  draw_method?: DrawMethod;
}

export interface RaffleFilters {
  status?: RaffleStatus;
  search?: string;
  user_id?: number;
  page?: number;
  page_size?: number;
}
