import { api } from "@/lib/api";
import type {
  AdminUserListItem,
  AdminUserDetail,
  UpdateUserStatusRequest,
  UpdateUserKYCRequest,
  OrganizerListItem,
  OrganizerDetail,
  VerifyOrganizerRequest,
  UpdateCommissionRequest,
  AdminRaffleListItem,
  AdminRaffleDetail,
  ForceStatusChangeRequest,
  AddAdminNotesRequest,
  ManualDrawRequest,
  CancelWithRefundRequest,
  PaginatedResponse,
  UserFilters,
  OrganizerFilters,
  RaffleFilters,
  PaginationParams,
} from "../types";

// ===========================
// User Management API
// ===========================

export const adminUsersApi = {
  // GET /api/v1/admin/users
  list: async (
    filters?: UserFilters,
    pagination?: PaginationParams
  ): Promise<PaginatedResponse<AdminUserListItem>> => {
    const params = new URLSearchParams();

    if (pagination) {
      params.append("page", pagination.page.toString());
      params.append("limit", pagination.limit.toString());
    }

    if (filters?.role) params.append("role", filters.role);
    if (filters?.status) params.append("status", filters.status);
    if (filters?.kyc_level) params.append("kyc_level", filters.kyc_level);
    if (filters?.search) params.append("search", filters.search);

    const response = await api.get(`/admin/users?${params.toString()}`);

    // El backend devuelve: { success, data: { Users, Total, Page, PageSize, TotalPages } }
    // Mapeamos a la estructura esperada por el frontend
    const backendData = response.data.data;
    return {
      data: backendData.Users || [],
      pagination: {
        page: backendData.Page || 1,
        limit: backendData.PageSize || 20,
        total: backendData.Total || 0,
        total_pages: backendData.TotalPages || 0,
      },
    };
  },

  // GET /api/v1/admin/users/:id
  getDetail: async (userId: number): Promise<AdminUserDetail> => {
    const response = await api.get(`/admin/users/${userId}`);

    // El backend devuelve: { success, data: { User, Stats, RecentRaffles } }
    const backendData = response.data.data;
    const user = backendData.User || backendData.user || {};
    const stats = backendData.Stats || backendData.stats || {};

    // Mapear explícitamente todos los campos
    return {
      id: user.id,
      uuid: user.uuid,
      email: user.email,
      first_name: user.first_name || "",
      last_name: user.last_name || "",
      phone: user.phone,
      cedula: user.cedula,
      role: user.role,
      kyc_level: user.kyc_level,
      status: user.status,
      email_verified: user.email_verified,
      phone_verified: user.phone_verified,
      created_at: user.created_at,
      updated_at: user.updated_at,
      raffle_stats: {
        total_raffles: stats.TotalRaffles || stats.total_raffles || 0,
        active_raffles: stats.ActiveRaffles || stats.active_raffles || 0,
        completed_raffles: stats.CompletedRaffles || stats.completed_raffles || 0,
        total_revenue: stats.TotalRevenue || stats.total_revenue || 0,
      },
      payment_stats: {
        total_payments: stats.TotalTicketsBought || stats.total_tickets_bought || 0,
        total_spent: stats.TotalSpent || stats.total_spent || 0,
        refund_count: 0, // TODO: agregar al backend
      },
      recent_audit_logs: [], // TODO: implementar en backend
    };
  },

  // PUT /api/v1/admin/users/:id/status
  updateStatus: async (
    userId: number,
    data: UpdateUserStatusRequest
  ): Promise<void> => {
    // Backend espera action (suspend/activate/ban/unban), no new_status
    const actionMap: Record<string, string> = {
      suspended: "suspend",
      active: "activate",
      banned: "ban",
    };
    await api.put(`/admin/users/${userId}/status`, {
      action: actionMap[data.new_status] || data.new_status,
      reason: data.reason,
    });
  },

  // PUT /api/v1/admin/users/:id/kyc
  updateKYC: async (
    userId: number,
    data: UpdateUserKYCRequest
  ): Promise<void> => {
    // Backend espera kyc_level, no new_kyc_level
    await api.put(`/admin/users/${userId}/kyc`, {
      kyc_level: data.new_kyc_level,
      notes: data.notes,
    });
  },

  // DELETE /api/v1/admin/users/:id
  deleteUser: async (userId: number): Promise<void> => {
    await api.delete(`/admin/users/${userId}`);
  },

  // POST /api/v1/admin/users/:id/reset-password
  resetPassword: async (userId: number): Promise<{ reset_token?: string }> => {
    const response = await api.post(`/admin/users/${userId}/reset-password`);
    return response.data;
  },
};

// ===========================
// Organizer Management API
// ===========================

export const adminOrganizersApi = {
  // GET /api/v1/admin/organizers
  list: async (
    filters?: OrganizerFilters,
    pagination?: PaginationParams
  ): Promise<PaginatedResponse<OrganizerListItem>> => {
    const params = new URLSearchParams();

    if (pagination) {
      params.append("page", pagination.page.toString());
      params.append("limit", pagination.limit.toString());
    }

    if (filters?.verified !== undefined) params.append("verified", filters.verified.toString());
    if (filters?.min_revenue) params.append("min_revenue", filters.min_revenue.toString());
    if (filters?.max_revenue) params.append("max_revenue", filters.max_revenue.toString());

    const response = await api.get(`/admin/organizers?${params.toString()}`);

    // El backend devuelve: { success, data: { Organizers, Total, Page, PageSize, TotalPages } }
    const backendData = response.data.data;
    return {
      data: (backendData.Organizers || []).map((item: any) => ({
        user: item.User,
        profile: item.Profile,
        metrics: {
          total_raffles: item.TotalRaffles || 0,
          active_raffles: item.ActiveRaffles || 0,
          completed_raffles: item.CompletedRaffles || 0,
          total_revenue: item.TotalRevenue || 0,
          pending_payout: item.PendingPayout || 0,
        },
      })),
      pagination: {
        page: backendData.Page || 1,
        limit: backendData.PageSize || 20,
        total: backendData.Total || 0,
        total_pages: backendData.TotalPages || 0,
      },
    };
  },

  // GET /api/v1/admin/organizers/:id
  getDetail: async (organizerId: number): Promise<OrganizerDetail> => {
    const response = await api.get(`/admin/organizers/${organizerId}`);

    // El backend devuelve: { success, data: { Profile, User, Revenue } }
    const backendData = response.data.data;

    return {
      user: backendData.User,
      profile: backendData.Profile,
      metrics: {
        total_raffles: 0, // TODO: agregar al backend
        active_raffles: 0, // TODO: agregar al backend
        completed_raffles: 0, // TODO: agregar al backend
        total_revenue: backendData.Revenue?.total_revenue || 0,
        pending_payout: backendData.Revenue?.pending_payout || 0,
      },
      raffle_list: [], // TODO: agregar al backend
      settlement_history: [], // TODO: agregar al backend
      revenue_breakdown: {
        gross_revenue: backendData.Revenue?.gross_revenue || 0,
        platform_fees: backendData.Revenue?.platform_fees || 0,
        net_revenue: backendData.Revenue?.net_revenue || 0,
      },
    };
  },

  // PUT /api/v1/admin/organizers/:id/verify
  verify: async (
    organizerId: number,
    data: VerifyOrganizerRequest
  ): Promise<void> => {
    await api.put(`/admin/organizers/${organizerId}/verify`, data);
  },

  // PUT /api/v1/admin/organizers/:id/commission
  updateCommission: async (
    organizerId: number,
    data: UpdateCommissionRequest
  ): Promise<void> => {
    await api.put(`/admin/organizers/${organizerId}/commission`, data);
  },
};

// ===========================
// Raffle Management API
// ===========================

export const adminRafflesApi = {
  // GET /api/v1/admin/raffles
  list: async (
    filters?: RaffleFilters,
    pagination?: PaginationParams
  ): Promise<PaginatedResponse<AdminRaffleListItem>> => {
    const params = new URLSearchParams();

    if (pagination) {
      params.append("page", pagination.page.toString());
      params.append("page_size", pagination.limit.toString());
    }

    if (filters?.status) params.append("status", filters.status);
    if (filters?.organizer_id) params.append("organizer_id", filters.organizer_id.toString());
    if (filters?.category_id) params.append("category_id", filters.category_id.toString());
    if (filters?.search) params.append("search", filters.search);
    if (filters?.date_from) params.append("date_from", filters.date_from);
    if (filters?.date_to) params.append("date_to", filters.date_to);
    if (filters?.order_by) params.append("order_by", filters.order_by);
    if (filters?.include_all) params.append("include_all", filters.include_all.toString());

    const response = await api.get(`/admin/raffles?${params.toString()}`);

    // El backend devuelve: { success, data: { Raffles, Total, Page, PageSize, TotalPages } }
    const backendData = response.data.data;
    return {
      data: backendData.Raffles || [],
      pagination: {
        page: backendData.Page || 1,
        limit: backendData.PageSize || 20,
        total: backendData.Total || 0,
        total_pages: backendData.TotalPages || 0,
      },
    };
  },

  // GET /api/v1/admin/raffles/:id/transactions
  getDetail: async (raffleId: number): Promise<AdminRaffleDetail> => {
    const response = await api.get(`/admin/raffles/${raffleId}/transactions`);

    // El backend devuelve: { success, data: { raffle, timeline, metrics } }
    const backendData = response.data.data;
    const raffle = backendData.raffle;

    // Calcular métricas básicas desde el raffle
    const totalRevenue = parseFloat(raffle?.TotalRevenue || 0);
    const platformFee = parseFloat(raffle?.PlatformFeeAmount || 0);
    const netRevenue = parseFloat(raffle?.NetAmount || 0);
    const conversionRate = raffle?.TotalNumbers > 0
      ? (raffle.SoldCount / raffle.TotalNumbers) * 100
      : 0;

    return {
      raffle: raffle,
      organizer_name: "", // TODO: El backend debería incluir esto
      organizer_email: "", // TODO: El backend debería incluir esto
      sold_count: raffle?.SoldCount || 0,
      reserved_count: raffle?.ReservedCount || 0,
      available_count: (raffle?.TotalNumbers || 0) - (raffle?.SoldCount || 0) - (raffle?.ReservedCount || 0),
      total_revenue: totalRevenue,
      platform_fee: platformFee,
      net_revenue: netRevenue,
      conversion_rate: conversionRate,
      timeline: backendData.timeline || [],
      transaction_metrics: backendData.metrics || {
        total_reservations: 0,
        total_payments: 0,
        total_refunds: 0,
        conversion_rate: 0,
        refund_rate: 0,
        total_revenue: 0,
        total_refunded: 0,
        net_revenue: 0,
      },
    };
  },

  // PUT /api/v1/admin/raffles/:id/force-status
  forceStatusChange: async (
    raffleId: number,
    data: ForceStatusChangeRequest
  ): Promise<void> => {
    await api.put(`/admin/raffles/${raffleId}/force-status`, data);
  },

  // PUT /api/v1/admin/raffles/:id/notes
  addAdminNotes: async (
    raffleId: number,
    data: AddAdminNotesRequest
  ): Promise<void> => {
    await api.put(`/admin/raffles/${raffleId}/notes`, data);
  },

  // POST /api/v1/admin/raffles/:id/manual-draw
  manualDraw: async (
    raffleId: number,
    data: ManualDrawRequest
  ): Promise<void> => {
    await api.post(`/admin/raffles/${raffleId}/manual-draw`, data);
  },

  // POST /api/v1/admin/raffles/:id/cancel-refund
  cancelWithRefund: async (
    raffleId: number,
    data: CancelWithRefundRequest
  ): Promise<void> => {
    await api.post(`/admin/raffles/${raffleId}/cancel-refund`, data);
  },
};
