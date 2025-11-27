import { useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate, useNavigate } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { GoogleOAuthProvider } from "@react-oauth/google";
import { Toaster, toast } from "sonner";
import { queryClient } from "@/lib/queryClient";
import { UserModeProvider } from "@/contexts/UserModeContext";
import { useInactivityTimeout } from "@/hooks/useInactivityTimeout";
import { useAuthStore } from "@/store/authStore";
import { setSessionExpiredCallback } from "@/lib/api";

// Google OAuth Client ID - debe configurarse en .env
const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID || "";

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
import { RaffleDetailPage } from "@/features/raffles/pages/RaffleDetailPage";
import { CreateRafflePage } from "@/features/raffles/pages/CreateRafflePage";
import { EditRafflePage } from "@/features/raffles/pages/EditRafflePage";
import { MyRafflesPage } from "@/features/raffles/pages/MyRafflesPage";
import { MyPurchasesPage } from "@/features/raffles/pages/MyPurchasesPage";

// Checkout pages
import { CheckoutPage } from "@/features/checkout/pages/CheckoutPage";
import { PaymentSuccessPage } from "@/features/checkout/pages/PaymentSuccessPage";
import { PaymentCancelPage } from "@/features/checkout/pages/PaymentCancelPage";

// Wallet pages
import { WalletPage } from "@/features/wallet/pages/WalletPage";
import { CreditSuccess } from "@/features/wallet/pages/CreditSuccess";
import { CreditFailed } from "@/features/wallet/pages/CreditFailed";
import { CreditVerifying } from "@/features/wallet/pages/CreditVerifying";

// Profile pages
import { ProfilePage } from "@/features/profile/components/ProfilePage";

// Admin pages
import { AdminRoute } from "@/features/admin/components/AdminRoute";
import { AdminLayout } from "@/features/admin/components/AdminLayout";
import { AdminDashboardPage } from "@/features/admin/pages/dashboard/AdminDashboardPage";
import { UsersListPage } from "@/features/admin/pages/users/UsersListPage";
import { UserDetailPage } from "@/features/admin/pages/users/UserDetailPage";
import { OrganizersListPage } from "@/features/admin/pages/organizers/OrganizersListPage";
import { OrganizerDetailPage } from "@/features/admin/pages/organizers/OrganizerDetailPage";
import { RafflesListPage as AdminRafflesListPage } from "@/features/admin/pages/raffles/RafflesListPage";
import { RaffleDetailPage as AdminRaffleDetailPage } from "@/features/admin/pages/raffles/RaffleDetailPage";
import { CategoriesPage } from "@/features/admin/pages/categories/CategoriesPage";
import { PaymentsListPage } from "@/features/admin/pages/payments/PaymentsListPage";
import { PaymentDetailPage } from "@/features/admin/pages/payments/PaymentDetailPage";
import { SettlementsListPage } from "@/features/admin/pages/settlements/SettlementsListPage";
import { SettlementDetailPage } from "@/features/admin/pages/settlements/SettlementDetailPage";
import { WalletsListPage } from "@/features/admin/wallets/components/WalletsListPage";
import { WalletDetailPage } from "@/features/admin/wallets/components/WalletDetailPage";
import { ReportsPage } from "@/features/admin/pages/reports/ReportsPage";
import { RevenueReportPage } from "@/features/admin/pages/reports/RevenueReportPage";
import { LiquidationsReportPage } from "@/features/admin/pages/reports/LiquidationsReportPage";
import { NotificationsPage } from "@/features/admin/pages/notifications/NotificationsPage";
import { SystemConfigPage } from "@/features/admin/pages/system/SystemConfigPage";
import { AuditLogsPage } from "@/features/admin/pages/audit/AuditLogsPage";

// Componente interno para usar hooks que requieren Router context
function AppRoutes() {
  const navigate = useNavigate();
  const logout = useAuthStore((state) => state.logout);
  const setUser = useAuthStore((state) => state.setUser);

  // Registrar callback para cuando el interceptor detecte sesión expirada (401)
  useEffect(() => {
    setSessionExpiredCallback(() => {
      // Limpiar el estado de Zustand sin llamar al backend (ya falló)
      setUser(null);
    });
  }, [setUser]);

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

          {/* Redirect /raffles to /explore */}
          <Route
            path="/raffles"
            element={<Navigate to="/explore" replace />}
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
            path="/raffles/:id/edit"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <EditRafflePage />
                </MainLayout>
              </ProtectedRoute>
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

          {/* Wallet routes (protected, with layout) */}
          <Route
            path="/wallet"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <WalletPage />
                </MainLayout>
              </ProtectedRoute>
            }
          />

          {/* Credit Purchase Result Pages (no layout for full screen experience) */}
          <Route
            path="/credits/success"
            element={
              <ProtectedRoute>
                <CreditSuccess />
              </ProtectedRoute>
            }
          />
          <Route
            path="/credits/failed"
            element={
              <ProtectedRoute>
                <CreditFailed />
              </ProtectedRoute>
            }
          />
          <Route
            path="/credits/verifying"
            element={
              <ProtectedRoute>
                <CreditVerifying />
              </ProtectedRoute>
            }
          />

          {/* Profile routes (protected, with layout) */}
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <MainLayout>
                  <ProfilePage />
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

          {/* Admin routes (protected, admin-only, with AdminLayout) */}
          <Route
            path="/admin/dashboard"
            element={
              <AdminRoute>
                <AdminLayout>
                  <AdminDashboardPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/users"
            element={
              <AdminRoute>
                <AdminLayout>
                  <UsersListPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/users/:id"
            element={
              <AdminRoute>
                <AdminLayout>
                  <UserDetailPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/organizers"
            element={
              <AdminRoute>
                <AdminLayout>
                  <OrganizersListPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/organizers/:id"
            element={
              <AdminRoute>
                <AdminLayout>
                  <OrganizerDetailPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/raffles"
            element={
              <AdminRoute>
                <AdminLayout>
                  <AdminRafflesListPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/raffles/:id"
            element={
              <AdminRoute>
                <AdminLayout>
                  <AdminRaffleDetailPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin: Categories */}
          <Route
            path="/admin/categories"
            element={
              <AdminRoute>
                <AdminLayout>
                  <CategoriesPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin Payments */}
          <Route
            path="/admin/payments"
            element={
              <AdminRoute>
                <AdminLayout>
                  <PaymentsListPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/payments/:id"
            element={
              <AdminRoute>
                <AdminLayout>
                  <PaymentDetailPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin Settlements */}
          <Route
            path="/admin/settlements"
            element={
              <AdminRoute>
                <AdminLayout>
                  <SettlementsListPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/settlements/:id"
            element={
              <AdminRoute>
                <AdminLayout>
                  <SettlementDetailPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin Wallets */}
          <Route
            path="/admin/wallets"
            element={
              <AdminRoute>
                <AdminLayout>
                  <WalletsListPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/wallets/:id"
            element={
              <AdminRoute>
                <AdminLayout>
                  <WalletDetailPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin Reports */}
          <Route
            path="/admin/reports"
            element={
              <AdminRoute>
                <AdminLayout>
                  <ReportsPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/reports/revenue"
            element={
              <AdminRoute>
                <AdminLayout>
                  <RevenueReportPage />
                </AdminLayout>
              </AdminRoute>
            }
          />
          <Route
            path="/admin/reports/liquidations"
            element={
              <AdminRoute>
                <AdminLayout>
                  <LiquidationsReportPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin Notifications */}
          <Route
            path="/admin/notifications"
            element={
              <AdminRoute>
                <AdminLayout>
                  <NotificationsPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin System Configuration */}
          <Route
            path="/admin/system"
            element={
              <AdminRoute>
                <AdminLayout>
                  <SystemConfigPage />
                </AdminLayout>
              </AdminRoute>
            }
          />

          {/* Admin Audit Logs */}
          <Route
            path="/admin/audit"
            element={
              <AdminRoute>
                <AdminLayout>
                  <AuditLogsPage />
                </AdminLayout>
              </AdminRoute>
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
      <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
        <UserModeProvider>
          <Toaster position="top-right" richColors closeButton />
          <BrowserRouter>
            <AppRoutes />
          </BrowserRouter>
        </UserModeProvider>
      </GoogleOAuthProvider>
    </QueryClientProvider>
  );
}

export default App;
