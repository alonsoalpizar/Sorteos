import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User } from "@/types/auth";
import { setTokens, clearTokens } from "@/lib/api";
import { logoutApi } from "@/features/auth/api/authApi";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setUser: (user: User | null) => void;
  setAuth: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;

  // Helpers
  hasMinimumKYC: (level: User["kyc_level"]) => boolean;
  isAdmin: () => boolean;
  isEmailVerified: () => boolean;
}

const kycLevels: Record<User["kyc_level"], number> = {
  none: 0,
  email_verified: 1,
  phone_verified: 2,
  cedula_verified: 3,
  full_kyc: 4,
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,

      setUser: (user) => {
        set({
          user,
          isAuthenticated: !!user,
          isLoading: false,
        });
      },

      setAuth: (user, accessToken, refreshToken) => {
        setTokens(accessToken, refreshToken);
        set({
          user,
          isAuthenticated: true,
          isLoading: false,
        });
      },

      logout: async () => {
        // Call backend to invalidate tokens (add to blacklist)
        try {
          await logoutApi();
        } catch (error) {
          console.error('Error calling logout API:', error);
          // Continue with local logout even if backend call fails
        }

        // Clear local tokens and state
        clearTokens();
        // Clear cart storage
        localStorage.removeItem('sorteos-cart-storage');
        set({
          user: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },

      setLoading: (loading) => {
        set({ isLoading: loading });
      },

      hasMinimumKYC: (level) => {
        const user = get().user;
        if (!user) return false;
        return kycLevels[user.kyc_level] >= kycLevels[level];
      },

      isAdmin: () => {
        const user = get().user;
        return user?.role === "admin" || user?.role === "super_admin";
      },

      isEmailVerified: () => {
        const user = get().user;
        return user?.email_verified ?? false;
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
