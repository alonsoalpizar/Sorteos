// Tipos para gestión de sorteos

export type RaffleStatus = 'draft' | 'active' | 'suspended' | 'completed' | 'cancelled';
export type DrawMethod = 'loteria_nacional_cr' | 'manual' | 'random';
export type SettlementStatus = 'pending' | 'processing' | 'completed' | 'failed';
export type RaffleNumberStatus = 'available' | 'reserved' | 'sold';

// Información del organizador (sin exponer user_id)
export interface OrganizerInfo {
  name: string;
  verified: boolean;
}

// PublicRaffle - Para usuarios NO autenticados o sin compras
// Oculta información financiera y sensible
export interface PublicRaffle {
  id: number;
  uuid: string;
  organizer: OrganizerInfo; // En lugar de user_id
  title: string;
  description: string;
  status: RaffleStatus;
  price_per_number: string;
  total_numbers: number;
  draw_date: string;
  draw_method: DrawMethod;
  sold_count: number;
  reserved_count: number;
  available_count: number;
  created_at: string;
  published_at?: string;
}

// BuyerRaffle - Para usuarios autenticados que HAN comprado
// Incluye gasto personal pero NO la información financiera del organizador
export interface BuyerRaffle extends PublicRaffle {
  my_total_spent: string;     // Cuánto ha gastado este usuario
  my_numbers_count: number;   // Cuántos números tiene
}

// OwnerRaffle - Para el organizador del sorteo y admins
// Incluye TODA la información financiera
export interface OwnerRaffle extends PublicRaffle {
  user_id: number; // Solo visible para owner/admin
  total_revenue: string;
  platform_fee_percentage: string;
  platform_fee_amount: string;
  net_amount: string;
  settlement_status: SettlementStatus;
}

// Raffle genérico - DEPRECATED: usar PublicRaffle, BuyerRaffle o OwnerRaffle
// Se mantiene temporalmente para compatibilidad con código existente
export interface Raffle {
  id: number;
  uuid: string;
  user_id?: number;
  title: string;
  description: string;
  status: RaffleStatus;
  price_per_number: string;
  total_numbers: number;
  draw_date: string;
  draw_method: DrawMethod;
  sold_count: number;
  reserved_count: number;
  available_count?: number;
  total_revenue?: string;
  platform_fee_percentage?: string;
  platform_fee_amount?: string;
  net_amount?: string;
  settlement_status?: SettlementStatus;
  winner_number?: string;
  winner_user_id?: number;
  published_at?: string;
  completed_at?: string;
  cancelled_at?: string;
  created_at: string;
  updated_at?: string;
  organizer?: OrganizerInfo;
  my_total_spent?: string;
  my_numbers_count?: number;
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

// Tipos para tickets del usuario
export interface UserTicketNumber {
  id: number;
  number: string;
  price: string;
  sold_at: string;
  payment_id?: number;
}

export interface TicketGroup {
  raffle: Raffle;
  numbers: UserTicketNumber[];
  total_numbers: number;
  total_spent: string;
}

export interface UserTicketsResponse {
  tickets: TicketGroup[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}
