import { Navigate } from "react-router-dom";
import { useAuthStore } from "@/store/authStore";
import { useIsAuthenticated } from "@/hooks/useAuth";

interface AdminRouteProps {
  children: React.ReactNode;
}

/**
 * AdminRoute - Protects admin routes
 *
 * Requirements:
 * 1. User must be authenticated
 * 2. User must have admin or super_admin role
 *
 * If not authenticated → redirect to /login
 * If authenticated but not admin → redirect to /dashboard
 */
export const AdminRoute = ({ children }: AdminRouteProps) => {
  const isAuthenticated = useIsAuthenticated();
  const isAdmin = useAuthStore((state) => state.isAdmin);

  // Check authentication first
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Check admin role
  if (!isAdmin()) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
};
