// Tipos para el perfil de usuario

export interface ProfileData {
  user: User;
  kyc_documents: KYCDocument[];
  wallet: Wallet;
  can_withdraw: boolean;
}

export interface User {
  id: number;
  uuid: string;
  email: string;
  email_verified: boolean;
  first_name?: string;
  last_name?: string;
  phone?: string;
  phone_verified: boolean;
  cedula?: string;
  date_of_birth?: string; // ISO date string
  profile_photo_url?: string;

  // Address
  address_line1?: string;
  address_line2?: string;
  city?: string;
  state?: string;
  postal_code?: string;
  country: string;

  // Banking
  iban?: string;

  // KYC & Status
  role: 'user' | 'admin' | 'super_admin';
  kyc_level: 'none' | 'email_verified' | 'phone_verified' | 'cedula_verified' | 'full_kyc';
  status: 'active' | 'suspended' | 'banned' | 'deleted';

  created_at: string;
  updated_at: string;
}

export interface KYCDocument {
  id: number;
  user_id: number;
  document_type: 'cedula_front' | 'cedula_back' | 'selfie';
  file_url: string;
  verification_status: 'pending' | 'approved' | 'rejected';
  verified_at?: string;
  verified_by?: number;
  rejected_reason?: string;
  uploaded_at: string;
}

export interface Wallet {
  id: number;
  uuid: string;
  user_id: number;
  balance_available: string;
  pending_balance: string;
  earnings_balance: string;
  currency: string;
  status: string;
  created_at: string;
  updated_at: string;
}

// Request types
export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  date_of_birth?: string; // YYYY-MM-DD format
  phone?: string;
  cedula?: string;
  address_line1?: string;
  address_line2?: string;
  city?: string;
  state?: string;
  postal_code?: string;
}

export interface ConfigureIBANRequest {
  iban: string;
}

export interface UploadPhotoRequest {
  photo_url: string;
}

export interface UploadKYCDocumentRequest {
  file_url: string;
}
