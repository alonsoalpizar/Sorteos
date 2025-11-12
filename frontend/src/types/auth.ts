// User types
export type UserRole = "user" | "admin" | "super_admin";
export type KYCLevel = "none" | "email_verified" | "phone_verified" | "cedula_verified" | "full_kyc";
export type UserStatus = "active" | "suspended" | "banned" | "deleted";

export interface User {
  id: number;
  uuid: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  cedula?: string;
  role: UserRole;
  kyc_level: KYCLevel;
  status: UserStatus;
  email_verified: boolean;
  phone_verified: boolean;
  created_at: string;
  updated_at: string;
}

// Auth API request types
export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
  accepted_terms: boolean;
  accepted_privacy: boolean;
  accepted_marketing?: boolean;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface VerifyEmailRequest {
  user_id: number;
  code: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

// Auth API response types
export interface RegisterResponse {
  user: User;
  verification_code_sent: boolean;
}

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

export interface VerifyEmailResponse {
  success: boolean;
  access_token: string;
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
}

// API wrapper response
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}
