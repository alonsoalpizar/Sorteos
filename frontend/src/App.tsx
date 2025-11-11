import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/lib/queryClient";

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

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Public routes (no layout) */}
          <Route path="/" element={<LandingPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/verify-email" element={<VerifyEmailPage />} />

          {/* Protected routes with layout */}
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
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
