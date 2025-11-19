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
    active_raffles: number;
    completed_raffles: number;
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
  raffle: {
    ID: number;
    UUID: string;
    UserID: number;
    Title: string;
    Description?: string;
    Status: "draft" | "active" | "suspended" | "completed" | "cancelled";
    CategoryID?: number;
    PricePerNumber: string;
    TotalNumbers: number;
    MinNumber: number;
    MaxNumber: number;
    DrawDate: string;
    DrawMethod: string;
    WinnerNumber?: string;
    WinnerUserID?: number;
    SoldCount: number;
    ReservedCount: number;
    TotalRevenue: string;
    PlatformFeePercentage: string;
    PlatformFeeAmount: string;
    NetAmount: string;
    SettledAt?: string;
    SettlementStatus: string;
    Metadata?: any;
    SuspensionReason?: string;
    SuspendedBy?: number;
    SuspendedAt?: string;
    AdminNotes?: string;
    CreatedAt: string;
    UpdatedAt: string;
    PublishedAt?: string;
    CompletedAt?: string;
    DeletedAt?: string;
  };
  organizer_name: string;
  organizer_email: string;
  sold_count: number;
  reserved_count: number;
  available_count: number;
  total_revenue: number;
  platform_fee: number;
  net_revenue: number;
  conversion_rate: number;
}

export interface AdminRaffleDetail extends AdminRaffleListItem {
  timeline: TransactionEvent[];
  transaction_metrics: {
    total_reservations: number;
    total_payments: number;
    total_refunds: number;
    conversion_rate: number;
    refund_rate: number;
    total_revenue: number;
    total_refunded: number;
    net_revenue: number;
  };
}

export interface TransactionEvent {
  type: "reservation" | "payment" | "refund" | "status_change" | "note";
  timestamp: string;
  user_id?: number;
  user_name?: string;
  amount?: number;
  status?: string;
  details: string;
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

export interface Payment {
  id: string; // UUID
  reservation_id: string;
  user_id: string; // UUID reference
  raffle_id: string; // UUID reference
  stripe_payment_intent_id: string;
  stripe_client_secret: string;
  amount: number;
  currency: string;
  status: string;
  payment_method?: string;
  error_message?: string;
  created_at: string;
  updated_at: string;
  paid_at?: string;
  provider?: string;
  refunded_at?: string;
  refunded_by?: number;
  admin_notes?: string;
}

export interface AdminPaymentListItem {
  payment: Payment;
  user_name: string;
  user_email: string;
  raffle_title: string;
  organizer_name: string;
}

export interface PaymentDetailEvent {
  type: "created" | "webhook" | "status_change" | "refund" | "note";
  timestamp: string;
  details: string;
  metadata?: Record<string, any>;
}

export interface RefundRecord {
  refunded_at: string;
  refunded_by: number;
  amount: number;
  type: "full" | "partial";
  reason: string;
  notes: string;
}

export interface WebhookEvent {
  received_at: string;
  provider: string;
  event_type: string;
  status: string;
  data?: Record<string, any>;
}

export interface AdminPaymentDetail {
  payment: Payment;
  user: any; // domain.User
  raffle: any; // domain.Raffle
  organizer: any; // domain.User
  numbers: string[];
  timeline: PaymentDetailEvent[];
  refund_history?: RefundRecord[];
  webhook_events?: WebhookEvent[];
}

export interface ProcessRefundRequest {
  reason: string;
  amount?: number;
  notes?: string;
}

export interface ManageDisputeRequest {
  action: "open" | "update" | "close" | "escalate";
  dispute_reason?: string;
  dispute_evidence?: string;
  resolution?: "accepted" | "rejected" | "refunded";
  admin_notes?: string;
  metadata?: Record<string, any>;
}

// ===========================
// Settlement Types
// ===========================

// Settlement base with details
export interface SettlementWithDetails {
  id: number;
  raffle_id: number;
  organizer_id: number;
  total_revenue: number;
  platform_fee: number;
  net_amount: number;
  status: "pending" | "approved" | "paid" | "rejected";
  created_at: string;
  approved_at?: string;
  approved_by?: number;
  rejected_at?: string;
  rejected_by?: number;
  rejection_reason?: string;
  paid_at?: string;
  payment_reference?: string;
  payment_method?: string;
  admin_notes?: string;
  updated_at: string;
  // Detalles adicionales
  raffle_title: string;
  organizer_name: string;
  organizer_email: string;
  organizer_kyc_level: string;
}

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

export interface PaymentsSummary {
  total_payments: number;
  succeeded_payments: number;
  refunded_payments: number;
  total_revenue: number;
  total_refunded: number;
  net_revenue: number;
  platform_fee_percent: number;
  platform_fee_amount: number;
}

export interface SettlementEvent {
  type: "calculated" | "approved" | "rejected" | "paid";
  timestamp: string;
  actor?: string;
  details: string;
  metadata?: Record<string, any>;
}

export interface OrganizerBankAccount {
  account_holder: string;
  bank_name: string;
  account_number: string;
  account_type: string;
  iban?: string;
  swift?: string;
  verified_at?: string;
}

export interface SettlementFullDetails {
  settlement: SettlementWithDetails;
  raffle: any;
  organizer: any;
  payments_summary: PaymentsSummary;
  timeline: SettlementEvent[];
  bank_account?: OrganizerBankAccount;
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
  notes?: string;
}

export interface MarkSettlementPaidRequest {
  payment_method: string;
  payment_reference?: string;
  notes?: string;
}

export interface AutoCreateSettlementsRequest {
  days_after_completion?: number;
  dry_run?: boolean;
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
  icon_url?: string;
  description?: string;
  is_active?: boolean;
}

export interface UpdateCategoryRequest {
  name?: string;
  icon_url?: string;
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
  status?: "draft" | "active" | "suspended" | "completed" | "cancelled";
  organizer_id?: number;
  category_id?: number;
  search?: string;
  date_from?: string;
  date_to?: string;
  order_by?: string;
  include_all?: boolean;
}

export interface PaymentFilters {
  status?: string;
  user_id?: number;
  raffle_id?: number;
  organizer_id?: number;
  provider?: string;
  date_from?: string;
  date_to?: string;
  min_amount?: number;
  max_amount?: number;
  search?: string;
  include_refund?: boolean;
}

export interface SettlementFilters {
  status?: "pending" | "approved" | "paid" | "rejected";
  organizer_id?: number;
  raffle_id?: number;
  date_from?: string;
  date_to?: string;
  min_amount?: number;
  max_amount?: number;
  search?: string;
  kyc_level?: string;
  pending_only?: boolean;
  order_by?: string;
}

// ===========================
// Reports Types
// ===========================

// Dashboard KPIs
export interface DashboardKPIs {
  // Usuarios
  total_users: number;
  active_users: number;
  suspended_users: number;
  banned_users: number;
  new_users_today: number;
  new_users_week: number;
  new_users_month: number;
  // Organizadores
  total_organizers: number;
  verified_organizers: number;
  pending_organizers: number;
  // Rifas
  total_raffles: number;
  active_raffles: number;
  completed_raffles: number;
  suspended_raffles: number;
  draft_raffles: number;
  // Revenue
  revenue_today: number;
  revenue_week: number;
  revenue_month: number;
  revenue_year: number;
  revenue_all_time: number;
  // Platform Fees
  platform_fees_today: number;
  platform_fees_month: number;
  platform_fees_all_time: number;
  // Settlements
  pending_settlements_count: number;
  pending_settlements_amount: number;
  approved_settlements_count: number;
  approved_settlements_amount: number;
  // Payments
  total_payments: number;
  succeeded_payments: number;
  pending_payments: number;
  failed_payments: number;
  refunded_payments: number;
  total_payments_amount: number;
  // Activity (últimas 24h)
  recent_users: number;
  recent_raffles: number;
  recent_payments: number;
  recent_settlements: number;
}

// Revenue Report
export interface RevenueDataPoint {
  date: string;
  gross_revenue: number;
  platform_fees: number;
  net_revenue: number;
  payment_count: number;
  raffle_count: number;
}

export interface RevenueReportInput {
  date_from: string;
  date_to: string;
  organizer_id?: number;
  category_id?: number;
  group_by?: "day" | "week" | "month";
}

export interface RevenueReportOutput {
  DataPoints: RevenueDataPoint[];
  TotalGrossRevenue: number;
  TotalPlatformFees: number;
  TotalNetRevenue: number;
  TotalPayments: number;
  TotalRaffles: number;
  AverageRevenuePerDay: number;
  AverageRevenuePerRaffle: number;
}

// Liquidations Report
export interface RaffleLiquidationRow {
  raffle_id: number;
  raffle_title: string;
  organizer_id: number;
  organizer_name: string;
  organizer_email: string;
  completed_at: string;
  gross_revenue: number;
  platform_fee_percent: number;
  platform_fee: number;
  net_revenue: number;
  settlement_id?: number;
  settlement_status?: string;
  paid_at?: string;
}

export interface RaffleLiquidationsReportInput {
  date_from: string;
  date_to: string;
  organizer_id?: number;
  category_id?: number;
  settlement_status?: string;
  order_by?: string;
}

export interface RaffleLiquidationsReportOutput {
  Rows: RaffleLiquidationRow[];
  Total: number;
  TotalGrossRevenue: number;
  TotalPlatformFees: number;
  TotalNetRevenue: number;
  WithSettlement: number;
  WithoutSettlement: number;
  PendingCount: number;
  ApprovedCount: number;
  PaidCount: number;
  RejectedCount: number;
}

// Export Data
export interface ExportDataInput {
  report_type: "users" | "organizers" | "raffles" | "payments" | "settlements";
  format: "csv" | "excel" | "pdf";
  date_from?: string;
  date_to?: string;
  filters?: Record<string, any>;
}

export interface ExportDataOutput {
  file_url: string;
  file_name: string;
  file_size: number;
  expires_at: string;
}

// ============================================================================
// NOTIFICATIONS MODULE
// ============================================================================

// Email Recipient
export interface EmailRecipient {
  email: string;
  name?: string;
}

// Send Email
export interface SendEmailInput {
  to: EmailRecipient[];
  cc?: EmailRecipient[];
  bcc?: EmailRecipient[];
  subject: string;
  body: string;
  template_id?: number;
  variables?: Record<string, any>;
  priority: "low" | "normal" | "high";
  scheduled_at?: string; // ISO 8601
}

export interface SendEmailOutput {
  notification_id: number;
  status: "queued" | "scheduled" | "sent" | "failed";
  sent_at?: string;
  scheduled_at?: string;
  recipients: number;
  message: string;
}

// Send Bulk Email
export interface SendBulkEmailInput {
  target_audience: "all_users" | "all_organizers" | "custom";
  custom_filters?: Record<string, any>;
  subject: string;
  body: string;
  template_id?: number;
  variables?: Record<string, any>;
  priority: "low" | "normal" | "high";
  scheduled_at?: string;
}

export interface SendBulkEmailOutput {
  notification_id: number;
  status: string;
  total_recipients: number;
  sent_count: number;
  failed_count: number;
  message: string;
}

// View Notification History
export interface ViewNotificationHistoryInput {
  type?: "email" | "sms" | "push" | "announcement";
  status?: "queued" | "sent" | "failed" | "scheduled";
  priority?: "low" | "normal" | "high" | "critical";
  admin_id?: number;
  date_from?: string;
  date_to?: string;
  search?: string;
  limit: number;
  offset: number;
}

export interface NotificationHistoryItem {
  id: number;
  type: string;
  subject?: string;
  recipients?: EmailRecipient[];
  recipient_count: number;
  priority: string;
  status: string;
  sent_at?: string;
  scheduled_at?: string;
  provider_status?: string;
  error?: string;
  admin_id: number;
  admin_email: string;
  created_at: string;
  metadata?: Record<string, any>;
}

export interface NotificationStatistics {
  total_sent: number;
  total_failed: number;
  total_queued: number;
  total_scheduled: number;
  success_rate: number;
  average_per_day: number;
  last_sent_at?: string;
}

export interface ViewNotificationHistoryOutput {
  notifications: NotificationHistoryItem[];
  total_count: number;
  statistics: NotificationStatistics;
}

// ============================================================================
// SYSTEM CONFIGURATION MODULE
// ============================================================================

// System Setting (snake_case - HAS json tags in Go)
export interface SystemSetting {
  key: string;
  value: any; // Can be string, number, boolean, object, etc.
  category: string;
  updated_at: string;
  updated_by?: number;
}

// Get System Settings
export interface GetSystemSettingsInput {
  category?: string;
  key?: string;
}

// Output wrapper uses PascalCase (no json tags), but Settings use snake_case
export interface GetSystemSettingsOutput {
  Settings: SystemSetting[];
  Categories: string[];
  TotalSettings: number;
}

// Update System Settings
export interface UpdateSystemSettingsInput {
  key: string;
  value: any;
  category: string;
}

// Output uses snake_case (HAS json tags in Go struct)
export interface UpdateSystemSettingsOutput {
  key: string;
  value: any;
  category: string;
  updated_at: string;
}

// ============================================================================
// AUDIT LOGS MODULE
// ============================================================================

// Audit Log
export interface AuditLog {
  id: number;
  admin_id: number;
  admin_name: string;
  admin_email: string;
  action: string;
  entity_type: string; // user, raffle, payment, settlement, etc.
  entity_id?: number;
  description: string;
  severity: string; // info, warning, error, critical
  ip_address?: string;
  user_agent?: string;
  metadata?: string; // JSON string
  created_at: string;
}

// List Audit Logs
export interface ListAuditLogsInput {
  page?: number;
  page_size?: number;
  admin_id?: number;
  action?: string;
  entity_type?: string;
  entity_id?: number;
  severity?: string;
  date_from?: string;
  date_to?: string;
  search?: string;
  order_by?: string;
}

export interface ListAuditLogsOutput {
  Logs: AuditLog[];
  Total: number;
  Page: number;
  PageSize: number;
  TotalPages: number;
  // Statistics
  InfoCount: number;
  WarningCount: number;
  ErrorCount: number;
  CriticalCount: number;
}
