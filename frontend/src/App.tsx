import { BrowserRouter, Routes, Route, Navigate, useNavigate } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { Toaster, toast } from "sonner";
import { queryClient } from "@/lib/queryClient";
import { UserModeProvider } from "@/contexts/UserModeContext";
import { useInactivityTimeout } from "@/hooks/useInactivityTimeout";
import { useAuthStore } from "@/store/authStore";

// Layout
import { MainLayout } from "@/components/layout/MainLayout";

// Auth pages
import { LoginPage } from "@/features/auth/pages/LoginPage";
import { RegisterPage } from "@/features/auth/pages/RegisterPage";
import { VerifyEmailPage } from "@/features/auth/pages/VerifyEmailPage";
import { ProtectedRoute } from "@/features/auth/components/ProtectedRoute";
import { DashboardPage } from "@/features/dashboard/pages/DashboardPage";

// Landing page
import { LandingPage } from "@/features/landing/pages/LandingPage";

// Participant pages
import { ExplorePage } from "@/features/raffles/pages/ExplorePage";
import { MyTicketsPage } from "@/features/raffles/pages/MyTicketsPage";

// Organizer pages
import { OrganizerDashboardPage } from "@/features/organizer/pages/OrganizerDashboardPage";

// Raffle pages
import { RafflesListPage } from "@/features/raffles/pages/RafflesListPage";
import { RaffleDetailPage } from "@/features/raffles/pages/RaffleDetailPage";
import { CreateRafflePage } from "@/features/raffles/pages/CreateRafflePage";
import { MyRafflesPage } from "@/features/raffles/pages/MyRafflesPage";
import { MyPurchasesPage } from "@/features/raffles/pages/MyPurchasesPage";

// Checkout pages
import { CheckoutPage } from "@/features/checkout/pages/CheckoutPage";
import { PaymentSuccessPage } from "@/features/checkout/pages/PaymentSuccessPage";
import { PaymentCancelPage } from "@/features/checkout/pages/PaymentCancelPage";

// Componente interno para usar hooks que requieren Router context
function AppRoutes() {
  const navigate = useNavigate();
  const logout = useAuthStore((state) => state.logout);

  // Timeout de inactividad: 30 minutos
  useInactivityTimeout({
    timeout: 30 * 60 * 1000, // 30 minutos
    warningTime: 2 * 60 * 1000, // Advertir 2 minutos antes
    onWarning: () => {
      toast.warning('Tu sesión expirará pronto', {
        description: 'Tienes 2 minutos de inactividad. Interactúa con la página para mantener tu sesión activa.',
        duration: 10000,
      });
    },
    onTimeout: async () => {
      toast.error('Sesión expirada', {
        description: 'Tu sesión ha expirado por inactividad. Por favor, inicia sesión nuevamente.',
        duration: 5000,
      });
      // Call logout to invalidate tokens on the backend
      await logout();
      navigate('/login');
    },
  });

  return (
          <Routes>
            {/* Public routes (no layout) */}
            <Route path="/" element={<LandingPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/verify-email" element={<VerifyEmailPage />} />

            {/* Participant routes (protected, with layout) */}
            <Route
              path="/explore"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <ExplorePage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/my-tickets"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <MyTicketsPage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />

            {/* Organizer routes (protected, with layout) */}
            <Route
              path="/organizer"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <OrganizerDashboardPage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/organizer/raffles"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <MyRafflesPage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/organizer/raffles/new"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <CreateRafflePage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />

            {/* Legacy dashboard route - redirect based on mode */}
            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <MainLayout>
                    <DashboardPage />
                  </MainLayout>
                </ProtectedRoute>
              }
            />

          {/* Raffle routes with layout */}
          <Route
            path="/raffles"
            element={
              <MainLayout>
                <RafflesListPage />
              </MainLayout>
            }
          />
          <Route
            path="/raffles/:id"
            element={
              <MainLayout>
                <RaffleDetailPage />
              </MainLayout>
            }
          />
          <Route
            path="/raffles/create"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <CreateRafflePage />
                </MainLayout>
              </ProtectedRoute>
            }
          />

          {/* My Raffles (protected, with layout) */}
          <Route
            path="/my-raffles"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <MyRafflesPage />
                </MainLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/my-purchases"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <MyPurchasesPage />
                </MainLayout>
              </ProtectedRoute>
            }
          />

          {/* Checkout routes (protected, with layout) */}
          <Route
            path="/checkout"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <CheckoutPage />
                </MainLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/payment/success"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <PaymentSuccessPage />
                </MainLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/payment/cancel"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <PaymentCancelPage />
                </MainLayout>
              </ProtectedRoute>
            }
          />

          {/* 404 redirect */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <UserModeProvider>
        <Toaster position="top-right" richColors closeButton />
        <BrowserRouter>
          <AppRoutes />
        </BrowserRouter>
      </UserModeProvider>
    </QueryClientProvider>
  );
}

export default App;
