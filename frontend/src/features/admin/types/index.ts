// Admin-specific types extending base User type
import type { User, UserRole, KYCLevel, UserStatus } from "@/types/auth";

// Re-export base types from auth
export type { UserRole, KYCLevel, UserStatus };

// ===========================
// User Management Types
// ===========================

export interface AdminUserListItem extends User {
  suspension_reason?: string;
  suspended_by?: number;
  suspended_at?: string;
  last_kyc_review?: string;
  kyc_reviewer?: number;
  last_login_at?: string;
}

export interface AdminUserDetail extends AdminUserListItem {
  raffle_stats: {
    total_raffles: number;
    active_raffles: number;
    completed_raffles: number;
    total_revenue: number;
  };
  payment_stats: {
    total_payments: number;
    total_spent: number;
    refund_count: number;
  };
  recent_audit_logs: AuditLog[];
}

export interface UpdateUserStatusRequest {
  new_status: UserStatus;
  reason?: string;
}

export interface UpdateUserKYCRequest {
  new_kyc_level: KYCLevel;
  notes?: string;
}

// ===========================
// Organizer Management Types
// ===========================

export interface OrganizerProfile {
  id: number;
  user_id: number;
  business_name: string;
  tax_id: string;
  bank_account_number?: string;
  bank_name?: string;
  commission_override?: number;
  total_payouts: number;
  pending_payout: number;
  verified: boolean;
  payout_schedule: "weekly" | "biweekly" | "monthly";
  created_at: string;
  updated_at: string;
}

export interface OrganizerListItem {
  user: AdminUserListItem;
  profile: OrganizerProfile;
  metrics: {
    total_raffles: number;
    total_revenue: number;
    pending_payout: number;
  };
}

export interface OrganizerDetail extends OrganizerListItem {
  raffle_list: AdminRaffleListItem[];
  settlement_history: Settlement[];
  revenue_breakdown: {
    gross_revenue: number;
    platform_fees: number;
    net_revenue: number;
  };
}

export interface UpdateCommissionRequest {
  commission_override: number;
}

export interface VerifyOrganizerRequest {
  verified: boolean;
  notes?: string;
}

// ===========================
// Raffle Management Types
// ===========================

export interface AdminRaffleListItem {
  id: number;
  uuid: string;
  title: string;
  user_id: number;
  organizer_name: string;
  category_id: number;
  category_name: string;
  status: "draft" | "active" | "suspended" | "completed" | "cancelled";
  total_numbers: number;
  sold_count: number;
  price: number;
  revenue: number;
  platform_fee: number;
  draw_date: string;
  suspended_by?: number;
  suspended_at?: string;
  suspension_reason?: string;
  admin_notes?: string;
  created_at: string;
}

export interface ForceStatusChangeRequest {
  new_status: "active" | "suspended" | "cancelled";
  reason: string;
}

export interface AddAdminNotesRequest {
  notes: string;
}

export interface ManualDrawRequest {
  winner_number?: string;
}

export interface CancelWithRefundRequest {
  reason: string;
}

// ===========================
// Payment Management Types
// ===========================

export interface AdminPaymentListItem {
  id: number;
  uuid: string;
  user_id: number;
  user_email: string;
  raffle_id: number;
  raffle_title: string;
  amount: number;
  status: "pending" | "succeeded" | "failed" | "refunded" | "disputed";
  payment_method: "stripe" | "paypal" | "cash";
  payment_intent_id?: string;
  created_at: string;
}

export interface ProcessRefundRequest {
  amount?: number;
  reason: string;
}

export interface ManageDisputeRequest {
  action: "open" | "resolve" | "reject";
  notes: string;
}

// ===========================
// Settlement Types
// ===========================

export interface Settlement {
  id: number;
  raffle_id: number;
  raffle_title: string;
  organizer_id: number;
  organizer_name: string;
  gross_revenue: number;
  platform_fee_percentage: number;
  platform_fee_amount: number;
  net_payout: number;
  status: "pending" | "approved" | "paid" | "rejected";
  approved_by?: number;
  approved_at?: string;
  payment_method?: string;
  payment_reference?: string;
  paid_at?: string;
  notes?: string;
  created_at: string;
}

export interface CreateSettlementRequest {
  raffle_ids: number[];
  organizer_id?: number;
}

export interface ApproveSettlementRequest {
  notes?: string;
}

export interface RejectSettlementRequest {
  reason: string;
}

export interface MarkSettlementPaidRequest {
  payment_method: string;
  payment_reference: string;
  notes?: string;
}

// ===========================
// Category Types
// ===========================

export interface Category {
  id: number;
  name: string;
  slug: string;
  icon?: string;
  description?: string;
  display_order: number;
  is_active: boolean;
  raffle_count?: number;
  created_at: string;
  updated_at: string;
}

export interface CreateCategoryRequest {
  name: string;
  icon?: string;
  description?: string;
}

export interface UpdateCategoryRequest {
  name?: string;
  icon?: string;
  description?: string;
  is_active?: boolean;
}

export interface ReorderCategoriesRequest {
  order: number[];
}

// ===========================
// System Configuration Types
// ===========================

export interface SystemParameter {
  id: number;
  key: string;
  value: string;
  value_type: "string" | "int" | "float" | "bool" | "json";
  category: "business" | "security" | "payment" | "email" | "notifications";
  description?: string;
  is_sensitive: boolean;
  updated_by?: number;
  updated_at?: string;
}

export interface UpdateParameterRequest {
  value: string;
}

export interface CompanySettings {
  id: number;
  company_name: string;
  tax_id: string;
  email: string;
  phone: string;
  address?: string;
  city?: string;
  country?: string;
  logo_url?: string;
  website_url?: string;
  updated_at: string;
}

export interface UpdateCompanySettingsRequest {
  company_name?: string;
  tax_id?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  country?: string;
  logo_url?: string;
  website_url?: string;
}

// ===========================
// Audit Log Types
// ===========================

export interface AuditLog {
  id: number;
  action: string;
  entity_type: string;
  entity_id?: number;
  user_id?: number;
  admin_id?: number;
  ip_address?: string;
  user_agent?: string;
  severity: "info" | "warning" | "error" | "critical";
  metadata?: Record<string, unknown>;
  created_at: string;
}

// ===========================
// Reports Types
// ===========================

export interface DashboardKPIs {
  total_users: {
    total: number;
    active: number;
    suspended: number;
    banned: number;
  };
  total_organizers: {
    total: number;
    verified: number;
    pending: number;
  };
  total_raffles: {
    total: number;
    active: number;
    completed: number;
    suspended: number;
  };
  revenue: {
    today: number;
    this_week: number;
    this_month: number;
    this_year: number;
    all_time: number;
  };
  platform_fees: {
    this_month: number;
    all_time: number;
  };
  pending_settlements: {
    count: number;
    total_amount: number;
  };
}

export interface RevenueReportItem {
  date: string;
  gross_revenue: number;
  platform_fees: number;
  net_revenue: number;
}

export interface LiquidationReportItem {
  raffle_id: number;
  raffle_title: string;
  organizer_name: string;
  gross_revenue: number;
  platform_fee: number;
  net_payout: number;
  settlement_status: Settlement["status"];
  completed_at: string;
}

// ===========================
// Notification Types
// ===========================

export interface SendEmailRequest {
  user_ids: number[];
  subject: string;
  body: string;
  template_id?: string;
  variables?: Record<string, string>;
}

export interface SendBulkEmailRequest {
  filters: {
    role?: UserRole;
    kyc_level?: KYCLevel;
    status?: UserStatus;
  };
  subject: string;
  body: string;
  template_id?: string;
}

export interface NotificationHistory {
  id: number;
  type: "email" | "sms";
  recipient_count: number;
  subject: string;
  status: "pending" | "sent" | "failed";
  sent_by: number;
  sent_at: string;
}

// ===========================
// Pagination & Filtering
// ===========================

export interface PaginationParams {
  page: number;
  limit: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface UserFilters {
  role?: UserRole;
  status?: UserStatus;
  kyc_level?: KYCLevel;
  search?: string;
}

export interface OrganizerFilters {
  verified?: boolean;
  min_revenue?: number;
  max_revenue?: number;
}

export interface RaffleFilters {
  status?: AdminRaffleListItem["status"];
  organizer_id?: number;
  category_id?: number;
  search?: string;
}

export interface PaymentFilters {
  status?: AdminPaymentListItem["status"];
  user_id?: number;
  raffle_id?: number;
  payment_method?: AdminPaymentListItem["payment_method"];
}

export interface SettlementFilters {
  status?: Settlement["status"];
  organizer_id?: number;
}
