import { api } from "@/lib/api";
import type {
  AdminUserListItem,
  AdminUserDetail,
  UpdateUserStatusRequest,
  UpdateUserKYCRequest,
  PaginatedResponse,
  UserFilters,
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
    return response.data.data;
  },

  // PUT /api/v1/admin/users/:id/status
  updateStatus: async (
    userId: number,
    data: UpdateUserStatusRequest
  ): Promise<void> => {
    await api.put(`/admin/users/${userId}/status`, data);
  },

  // PUT /api/v1/admin/users/:id/kyc
  updateKYC: async (
    userId: number,
    data: UpdateUserKYCRequest
  ): Promise<void> => {
    await api.put(`/admin/users/${userId}/kyc`, data);
  },

  // DELETE /api/v1/admin/users/:id
  deleteUser: async (userId: number): Promise<void> => {
    await api.delete(`/admin/users/${userId}`);
  },
};
