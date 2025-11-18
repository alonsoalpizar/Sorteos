import { Card } from "@/components/ui/Card";
import { Users, UserCog, Ticket, DollarSign } from "lucide-react";

export function AdminDashboardPage() {
  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900">Panel de Administración</h1>
        <p className="text-slate-600 mt-2">Vista general del sistema Sorteos.club</p>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card className="p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-slate-600">Total Usuarios</p>
              <p className="text-3xl font-bold text-slate-900 mt-2">-</p>
              <p className="text-xs text-slate-500 mt-1">Próximamente</p>
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
              <p className="text-3xl font-bold text-slate-900 mt-2">-</p>
              <p className="text-xs text-slate-500 mt-1">Próximamente</p>
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
              <p className="text-3xl font-bold text-slate-900 mt-2">-</p>
              <p className="text-xs text-slate-500 mt-1">Próximamente</p>
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
              <p className="text-3xl font-bold text-slate-900 mt-2">-</p>
              <p className="text-xs text-slate-500 mt-1">Próximamente</p>
            </div>
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
              <DollarSign className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </Card>
      </div>

      {/* Welcome Message */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-slate-900 mb-4">
          Bienvenido al Panel de Administración
        </h2>
        <div className="prose prose-slate max-w-none">
          <p className="text-slate-600">
            Este es el módulo Almighty Admin de Sorteos.club. Desde aquí puedes gestionar todos los aspectos
            de la plataforma:
          </p>
          <ul className="text-slate-600 space-y-2 mt-4">
            <li><strong>Usuarios:</strong> Gestión completa de usuarios, KYC, suspensiones</li>
            <li><strong>Organizadores:</strong> Perfiles, comisiones, verificación</li>
            <li><strong>Rifas:</strong> Control administrativo, suspensiones, sorteos manuales</li>
            <li><strong>Pagos:</strong> Procesamiento de refunds y disputas</li>
            <li><strong>Liquidaciones:</strong> Aprobación y pagos a organizadores</li>
            <li><strong>Categorías:</strong> CRUD completo de categorías de rifas</li>
            <li><strong>Reportes:</strong> Métricas financieras y operacionales</li>
            <li><strong>Notificaciones:</strong> Envío de emails administrativos</li>
            <li><strong>Configuración:</strong> Parámetros del sistema</li>
            <li><strong>Auditoría:</strong> Logs de todas las acciones administrativas</li>
          </ul>

          <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <p className="text-sm text-blue-800">
              <strong>Estado del desarrollo:</strong> Fase 1 - Setup base + Users Management en progreso.
              Los módulos se irán activando progresivamente según el roadmap de 7-8 semanas.
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
