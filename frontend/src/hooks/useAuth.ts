import React from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as authApi from "@/api/auth";
import { useAuthStore } from "@/store/authStore";
import type {
  RegisterRequest,
  LoginRequest,
  VerifyEmailRequest,
} from "@/types/auth";
import { getErrorMessage } from "@/lib/api";

/**
 * Hook for user registration
 */
export const useRegister = () => {
  const setUser = useAuthStore((state) => state.setUser);

  return useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
    onSuccess: (data) => {
      // After registration, user is not authenticated yet (needs email verification)
      setUser(data.user);
    },
  });
};

/**
 * Hook for user login
 */
export const useLogin = () => {
  const setAuth = useAuthStore((state) => state.setAuth);

  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: (data) => {
      setAuth(data.user, data.access_token, data.refresh_token);
    },
  });
};

/**
 * Hook for email verification
 */
export const useVerifyEmail = () => {
  const setAuth = useAuthStore((state) => state.setAuth);
  const user = useAuthStore((state) => state.user);

  return useMutation({
    mutationFn: (data: VerifyEmailRequest) => authApi.verifyEmail(data),
    onSuccess: (data) => {
      // After verification, user is automatically logged in
      if (user) {
        const updatedUser = {
          ...user,
          email_verified: true,
          kyc_level: "email_verified" as const,
        };
        setAuth(updatedUser, data.access_token, data.refresh_token);
      }
    },
  });
};

/**
 * Hook for logout
 */
export const useLogout = () => {
  const logout = useAuthStore((state) => state.logout);
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      logout();
      queryClient.clear();
    },
    onError: () => {
      // Even if API call fails, logout locally
      logout();
      queryClient.clear();
    },
  });
};

/**
 * Hook to get current user from API
 */
export const useCurrentUser = () => {
  const { isAuthenticated, setUser } = useAuthStore();

  const query = useQuery({
    queryKey: ["currentUser"],
    queryFn: () => authApi.getCurrentUser(),
    enabled: isAuthenticated,
  });

  // Handle side effects with useEffect
  React.useEffect(() => {
    if (query.data) {
      setUser(query.data);
    }
    if (query.error) {
      console.error("Failed to fetch current user:", getErrorMessage(query.error));
      // Don't logout on error, token might still be valid
    }
  }, [query.data, query.error, setUser]);

  return query;
};

/**
 * Hook to check authentication status
 */
export const useIsAuthenticated = () => {
  return useAuthStore((state) => state.isAuthenticated);
};

/**
 * Hook to check if user is admin
 */
export const useIsAdmin = () => {
  return useAuthStore((state) => state.isAdmin());
};

/**
 * Hook to check if email is verified
 */
export const useIsEmailVerified = () => {
  return useAuthStore((state) => state.isEmailVerified());
};

/**
 * Hook to get current user
 */
export const useUser = () => {
  return useAuthStore((state) => state.user);
};
