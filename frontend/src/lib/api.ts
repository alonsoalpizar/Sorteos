import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";

// Base API URL - uses Vite proxy in development
const API_BASE_URL = import.meta.env.VITE_API_URL || "/api/v1";

// Create axios instance
export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30 seconds - increased for slow operations like registration with email sending
  headers: {
    "Content-Type": "application/json",
  },
});

// Token management
let accessToken: string | null = null;
let refreshToken: string | null = null;

export const setTokens = (access: string, refresh: string) => {
  accessToken = access;
  refreshToken = refresh;
  localStorage.setItem("access_token", access);
  localStorage.setItem("refresh_token", refresh);
};

export const getAccessToken = (): string | null => {
  if (!accessToken) {
    accessToken = localStorage.getItem("access_token");
  }
  return accessToken;
};

export const getRefreshToken = (): string | null => {
  if (!refreshToken) {
    refreshToken = localStorage.getItem("refresh_token");
  }
  return refreshToken;
};

export const clearTokens = () => {
  accessToken = null;
  refreshToken = null;
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
};

// Request interceptor to add auth token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAccessToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle token refresh
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: Error | null, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // If error is 401 and we haven't tried to refresh yet
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // If already refreshing, queue this request
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
            }
            return api(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshTokenValue = getRefreshToken();

      if (!refreshTokenValue) {
        // No refresh token, redirect to login
        clearTokens();
        window.location.href = "/login";
        return Promise.reject(error);
      }

      try {
        // Try to refresh the token
        const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
          refresh_token: refreshTokenValue,
        });

        const { access_token, refresh_token } = response.data.data;
        setTokens(access_token, refresh_token);

        // Update the failed requests with new token
        processQueue(null, access_token);

        // Retry the original request
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
        }
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        processQueue(refreshError as Error, null);
        clearTokens();
        window.location.href = "/login";
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

// API Error type
export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

// Helper to extract error message from API response
export const getErrorMessage = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const apiError = error.response?.data as ApiError | undefined;
    return apiError?.message || error.message || "An error occurred";
  }
  if (error instanceof Error) {
    return error.message;
  }
  return "An unknown error occurred";
};
