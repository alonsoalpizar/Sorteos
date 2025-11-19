import { useNavigate } from "react-router-dom";
import {
  Users,
  Package,
  DollarSign,
  TrendingUp,
  AlertCircle,
  CheckCircle,
  Clock,
  BarChart3,
  FileText,
} from "lucide-react";
import { useAdminDashboard } from "../../hooks/useAdminReports";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import { formatCurrency } from "@/lib/currency";

export function ReportsPage() {
  const navigate = useNavigate();
  const { data: kpis, isLoading, error } = useAdminDashboard();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>
    );
  }

  if (error || !kpis) {
    return (
      <div className="p-6">
        <EmptyState
          icon={<AlertCircle className="w-12 h-12 text-red-500" />}
          title="Error al cargar dashboard"
          description={(error as Error)?.message || "No se pudieron cargar las métricas"}
        />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">Dashboard de Reportes</h1>
          <p className="text-slate-600 mt-2">
            Métricas globales y reportes del sistema
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={() => navigate("/admin/reports/revenue")}
          >
            <BarChart3 className="w-4 h-4 mr-2" />
            Reporte de Ingresos
          </Button>
          <Button
            variant="outline"
            onClick={() => navigate("/admin/reports/liquidations")}
          >
            <FileText className="w-4 h-4 mr-2" />
            Reporte de Liquidaciones
          </Button>
        </div>
      </div>

      {/* Usuarios */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Users className="w-5 h-5" />
          Usuarios
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Total Usuarios</p>
            <p className="text-3xl font-bold text-slate-900 mt-2">
              {kpis.total_users.toLocaleString()}
            </p>
            <div className="flex items-center gap-2 mt-2">
              <span className="text-xs text-green-600">
                +{kpis.new_users_today} hoy
              </span>
              <span className="text-xs text-slate-500">
                +{kpis.new_users_month} este mes
              </span>
            </div>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Activos</p>
            <p className="text-3xl font-bold text-green-600 mt-2">
              {kpis.active_users.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Suspendidos</p>
            <p className="text-3xl font-bold text-yellow-600 mt-2">
              {kpis.suspended_users.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Baneados</p>
            <p className="text-3xl font-bold text-red-600 mt-2">
              {kpis.banned_users.toLocaleString()}
            </p>
          </Card>
        </div>
      </div>

      {/* Organizadores */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <CheckCircle className="w-5 h-5" />
          Organizadores
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Total</p>
            <p className="text-3xl font-bold text-slate-900 mt-2">
              {kpis.total_organizers.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Verificados</p>
            <p className="text-3xl font-bold text-green-600 mt-2">
              {kpis.verified_organizers.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Pendientes</p>
            <p className="text-3xl font-bold text-yellow-600 mt-2">
              {kpis.pending_organizers.toLocaleString()}
            </p>
          </Card>
        </div>
      </div>

      {/* Rifas */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Package className="w-5 h-5" />
          Rifas
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Total</p>
            <p className="text-3xl font-bold text-slate-900 mt-2">
              {kpis.total_raffles.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Activas</p>
            <p className="text-3xl font-bold text-blue-600 mt-2">
              {kpis.active_raffles.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Completadas</p>
            <p className="text-3xl font-bold text-green-600 mt-2">
              {kpis.completed_raffles.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Borradores</p>
            <p className="text-3xl font-bold text-slate-500 mt-2">
              {kpis.draft_raffles.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Suspendidas</p>
            <p className="text-3xl font-bold text-red-600 mt-2">
              {kpis.suspended_raffles.toLocaleString()}
            </p>
          </Card>
        </div>
      </div>

      {/* Revenue */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <TrendingUp className="w-5 h-5" />
          Ingresos
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Hoy</p>
            <p className="text-2xl font-bold text-green-600 mt-2">
              {formatCurrency(kpis.revenue_today)}
            </p>
            <p className="text-xs text-slate-500 mt-1">
              Comisión: {formatCurrency(kpis.platform_fees_today)}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Esta Semana</p>
            <p className="text-2xl font-bold text-green-600 mt-2">
              {formatCurrency(kpis.revenue_week)}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Este Mes</p>
            <p className="text-2xl font-bold text-green-600 mt-2">
              {formatCurrency(kpis.revenue_month)}
            </p>
            <p className="text-xs text-slate-500 mt-1">
              Comisión: {formatCurrency(kpis.platform_fees_month)}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Total Histórico</p>
            <p className="text-2xl font-bold text-green-600 mt-2">
              {formatCurrency(kpis.revenue_all_time)}
            </p>
            <p className="text-xs text-slate-500 mt-1">
              Comisión: {formatCurrency(kpis.platform_fees_all_time)}
            </p>
          </Card>
        </div>
      </div>

      {/* Liquidaciones */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <DollarSign className="w-5 h-5" />
          Liquidaciones
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Card className="p-4">
            <div className="flex items-center justify-between mb-2">
              <p className="text-sm font-medium text-slate-600">Pendientes</p>
              <Clock className="w-5 h-5 text-yellow-600" />
            </div>
            <p className="text-2xl font-bold text-yellow-600">
              {kpis.pending_settlements_count} liquidaciones
            </p>
            <p className="text-lg text-slate-700 mt-1">
              {formatCurrency(kpis.pending_settlements_amount)}
            </p>
          </Card>

          <Card className="p-4">
            <div className="flex items-center justify-between mb-2">
              <p className="text-sm font-medium text-slate-600">Aprobadas</p>
              <CheckCircle className="w-5 h-5 text-blue-600" />
            </div>
            <p className="text-2xl font-bold text-blue-600">
              {kpis.approved_settlements_count} liquidaciones
            </p>
            <p className="text-lg text-slate-700 mt-1">
              {formatCurrency(kpis.approved_settlements_amount)}
            </p>
          </Card>
        </div>
      </div>

      {/* Pagos */}
      <div>
        <h2 className="text-xl font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <BarChart3 className="w-5 h-5" />
          Pagos
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Total</p>
            <p className="text-2xl font-bold text-slate-900 mt-2">
              {kpis.total_payments.toLocaleString()}
            </p>
            <p className="text-sm text-slate-600 mt-1">
              {formatCurrency(kpis.total_payments_amount)}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Exitosos</p>
            <p className="text-2xl font-bold text-green-600 mt-2">
              {kpis.succeeded_payments.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Pendientes</p>
            <p className="text-2xl font-bold text-yellow-600 mt-2">
              {kpis.pending_payments.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Fallidos</p>
            <p className="text-2xl font-bold text-red-600 mt-2">
              {kpis.failed_payments.toLocaleString()}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm font-medium text-slate-600">Reembolsados</p>
            <p className="text-2xl font-bold text-slate-600 mt-2">
              {kpis.refunded_payments.toLocaleString()}
            </p>
          </Card>
        </div>
      </div>

      {/* Actividad Reciente (últimas 24h) */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <Clock className="w-5 h-5" />
          Actividad Reciente (Últimas 24 horas)
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div>
            <p className="text-sm font-medium text-slate-600">Nuevos Usuarios</p>
            <p className="text-3xl font-bold text-blue-600 mt-2">
              {kpis.recent_users.toLocaleString()}
            </p>
          </div>
          <div>
            <p className="text-sm font-medium text-slate-600">Nuevas Rifas</p>
            <p className="text-3xl font-bold text-purple-600 mt-2">
              {kpis.recent_raffles.toLocaleString()}
            </p>
          </div>
          <div>
            <p className="text-sm font-medium text-slate-600">Pagos Procesados</p>
            <p className="text-3xl font-bold text-green-600 mt-2">
              {kpis.recent_payments.toLocaleString()}
            </p>
          </div>
          <div>
            <p className="text-sm font-medium text-slate-600">Liquidaciones Creadas</p>
            <p className="text-3xl font-bold text-orange-600 mt-2">
              {kpis.recent_settlements.toLocaleString()}
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
