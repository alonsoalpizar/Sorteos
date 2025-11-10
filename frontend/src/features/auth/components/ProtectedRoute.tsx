import { Navigate } from "react-router-dom";
import { useIsAuthenticated, useIsEmailVerified } from "@/hooks/useAuth";

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireEmailVerification?: boolean;
}

export const ProtectedRoute = ({
  children,
  requireEmailVerification = true,
}: ProtectedRouteProps) => {
  const isAuthenticated = useIsAuthenticated();
  const isEmailVerified = useIsEmailVerified();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requireEmailVerification && !isEmailVerified) {
    return <Navigate to="/verify-email" replace />;
  }

  return <>{children}</>;
};
