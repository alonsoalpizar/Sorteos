import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/queryClient";

// Auth pages
import { LoginPage } from "@/features/auth/pages/LoginPage";
import { RegisterPage } from "@/features/auth/pages/RegisterPage";
import { VerifyEmailPage } from "@/features/auth/pages/VerifyEmailPage";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { DashboardPage } from "@/features/dashboard/pages/DashboardPage";

// Raffle pages
import { RafflesListPage } from "@/features/raffles/pages/RafflesListPage";
import { RaffleDetailPage } from "@/features/raffles/pages/RaffleDetailPage";
import { CreateRafflePage } from "@/features/raffles/pages/CreateRafflePage";

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Public routes */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/verify-email" element={<VerifyEmailPage />} />

          {/* Raffle routes (public) */}
          <Route path="/raffles" element={<RafflesListPage />} />
          <Route path="/raffles/:id" element={<RaffleDetailPage />} />

          {/* Raffle routes (protected) */}
          <Route
            path="/raffles/create"
            element={
              <ProtectedRoute>
                <CreateRafflePage />
              </ProtectedRoute>
            }
          />

          {/* Protected routes */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <DashboardPage />
              </ProtectedRoute>
            }
          />

          {/* Default redirect */}
          <Route path="/" element={<Navigate to="/raffles" replace />} />
          <Route path="*" element={<Navigate to="/raffles" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
