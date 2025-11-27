import { Card } from "@/components/ui/Card";
import { Link } from "react-router-dom";
import {
  Users,
  UserCog,
  Ticket,
  DollarSign,
  CreditCard,
  FolderTree,
  BarChart3,
  Bell,
  Settings,
  FileText,
} from "lucide-react";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { useAdminDashboard } from "../../hooks/useAdminReports";
import { formatCurrency } from "@/lib/currency";

export function AdminDashboardPage() {
  const { data, isLoading, error } = useAdminDashboard();

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900">Panel de Administración</h1>
        <p className="text-slate-600 mt-2">Vista general del sistema Sorteos.club</p>
      </div>

      {/* KPI Cards */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <LoadingSpinner />
        </div>
      ) : error ? (
        <div className="text-center py-12">
          <p className="text-red-600">Error al cargar los datos del dashboard</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Total Usuarios</p>
                <p className="text-3xl font-bold text-slate-900 mt-2">
                  {data?.total_users.toLocaleString() || 0}
                </p>
                <p className="text-xs text-slate-500 mt-1">
                  {data?.active_users || 0} activos
                </p>
              </div>
              <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                <Users className="w-6 h-6 text-blue-600" />
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Organizadores</p>
                <p className="text-3xl font-bold text-slate-900 mt-2">
                  {data?.total_organizers.toLocaleString() || 0}
                </p>
                <p className="text-xs text-slate-500 mt-1">
                  {data?.verified_organizers || 0} verificados
                </p>
              </div>
              <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
                <UserCog className="w-6 h-6 text-green-600" />
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Rifas Activas</p>
                <p className="text-3xl font-bold text-slate-900 mt-2">
                  {data?.active_raffles.toLocaleString() || 0}
                </p>
                <p className="text-xs text-slate-500 mt-1">
                  {data?.total_raffles || 0} totales
                </p>
              </div>
              <div className="w-12 h-12 bg-amber-100 rounded-lg flex items-center justify-center">
                <Ticket className="w-6 h-6 text-amber-600" />
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-slate-600">Ingresos del Mes</p>
                <p className="text-3xl font-bold text-slate-900 mt-2">
                  {formatCurrency(data?.revenue_month || 0)}
                </p>
                <p className="text-xs text-slate-500 mt-1">
                  {formatCurrency(data?.revenue_all_time || 0)} totales
                </p>
              </div>
              <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                <DollarSign className="w-6 h-6 text-blue-600" />
              </div>
            </div>
          </Card>
        </div>
      )}

      {/* Welcome Message */}
      <Card className="p-8 bg-gradient-to-br from-blue-50 to-white border-blue-100">
        <div className="flex items-start gap-4">
          <div className="w-12 h-12 bg-blue-600 rounded-xl flex items-center justify-center flex-shrink-0">
            <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div className="flex-1">
            <h2 className="text-2xl font-bold text-slate-900 mb-2">
              Bienvenido al Panel de Administración
            </h2>
            <p className="text-slate-600 text-lg mb-6">
              Este es el módulo <span className="font-semibold text-blue-700">Almighty Admin</span> de Sorteos.club. Desde aquí puedes gestionar todos los aspectos de la plataforma:
            </p>
          </div>
        </div>

        {/* Módulos Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6">
          <Link to="/admin/users" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-blue-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <Users className="w-5 h-5 text-blue-600" />
                <h3 className="font-semibold text-slate-900">Usuarios</h3>
              </div>
              <p className="text-sm text-slate-600">Gestión completa de usuarios, KYC, suspensiones</p>
            </div>
          </Link>

          <Link to="/admin/organizers" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-green-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <UserCog className="w-5 h-5 text-green-600" />
                <h3 className="font-semibold text-slate-900">Organizadores</h3>
              </div>
              <p className="text-sm text-slate-600">Perfiles, comisiones, verificación</p>
            </div>
          </Link>

          <Link to="/admin/raffles" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-amber-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <Ticket className="w-5 h-5 text-amber-600" />
                <h3 className="font-semibold text-slate-900">Rifas</h3>
              </div>
              <p className="text-sm text-slate-600">Control administrativo, suspensiones, sorteos manuales</p>
            </div>
          </Link>

          <Link to="/admin/payments" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-purple-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <CreditCard className="w-5 h-5 text-purple-600" />
                <h3 className="font-semibold text-slate-900">Pagos</h3>
              </div>
              <p className="text-sm text-slate-600">Procesamiento de refunds y disputas</p>
            </div>
          </Link>

          <Link to="/admin/settlements" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-blue-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <DollarSign className="w-5 h-5 text-blue-600" />
                <h3 className="font-semibold text-slate-900">Liquidaciones</h3>
              </div>
              <p className="text-sm text-slate-600">Aprobación y pagos a organizadores</p>
            </div>
          </Link>

          <Link to="/admin/categories" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-orange-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <FolderTree className="w-5 h-5 text-orange-600" />
                <h3 className="font-semibold text-slate-900">Categorías</h3>
              </div>
              <p className="text-sm text-slate-600">CRUD completo de categorías de rifas</p>
            </div>
          </Link>

          <Link to="/admin/reports" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-green-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <BarChart3 className="w-5 h-5 text-green-600" />
                <h3 className="font-semibold text-slate-900">Reportes</h3>
              </div>
              <p className="text-sm text-slate-600">Métricas financieras y operacionales</p>
            </div>
          </Link>

          <Link to="/admin/notifications" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-indigo-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <Bell className="w-5 h-5 text-indigo-600" />
                <h3 className="font-semibold text-slate-900">Notificaciones</h3>
              </div>
              <p className="text-sm text-slate-600">Envío de emails administrativos</p>
            </div>
          </Link>

          <Link to="/admin/system" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-slate-400 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <Settings className="w-5 h-5 text-slate-600" />
                <h3 className="font-semibold text-slate-900">Configuración</h3>
              </div>
              <p className="text-sm text-slate-600">Parámetros del sistema</p>
            </div>
          </Link>

          <Link to="/admin/audit" className="block">
            <div className="p-4 bg-white rounded-lg border border-slate-200 hover:border-blue-300 hover:shadow-md transition-all cursor-pointer">
              <div className="flex items-center gap-3 mb-2">
                <FileText className="w-5 h-5 text-blue-600" />
                <h3 className="font-semibold text-slate-900">Auditoría</h3>
              </div>
              <p className="text-sm text-slate-600">Logs de todas las acciones administrativas</p>
            </div>
          </Link>
        </div>

        {/* Estado del desarrollo */}
        <div className="mt-6 p-5 bg-gradient-to-r from-blue-600 to-blue-700 rounded-xl text-white shadow-lg">
          <div className="flex items-start gap-3">
            <svg className="w-6 h-6 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
            </svg>
            <div>
              <p className="font-semibold text-lg mb-1">Estado del desarrollo</p>
              <p className="text-blue-100">
                <strong className="text-white">Fase actual:</strong> Reportes completado ✓
              </p>
              <p className="text-sm text-blue-100 mt-2">
                Los módulos se están activando progresivamente. Backend 100% completo con 52 endpoints administrativos.
              </p>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}
