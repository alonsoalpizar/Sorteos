import { useParams, useNavigate } from "react-router-dom";
import { ArrowLeft, Shield, AlertCircle, Ban, CheckCircle } from "lucide-react";
import { useAdminUserDetail, useUpdateUserStatus, useUpdateUserKYC } from "../../hooks/useAdminUsers";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import type { UserStatus, KYCLevel } from "../../types";
import { format } from "date-fns";

export function UserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const userId = parseInt(id || "0", 10);

  const { data: user, isLoading, error } = useAdminUserDetail(userId);
  const updateStatus = useUpdateUserStatus();
  const updateKYC = useUpdateUserKYC();

  const handleUpdateStatus = (newStatus: UserStatus, reason?: string) => {
    if (!confirm(`¿Confirmas cambiar el estado a "${newStatus}"?`)) return;

    updateStatus.mutate({
      userId,
      data: { new_status: newStatus, reason },
    });
  };

  const handleUpdateKYC = (newLevel: KYCLevel) => {
    if (!confirm(`¿Confirmas cambiar el nivel KYC a "${newLevel}"?`)) return;

    updateKYC.mutate({
      userId,
      data: { new_kyc_level: newLevel },
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="p-6">
        <EmptyState
          icon={<AlertCircle className="w-12 h-12 text-red-500" />}
          title="Error al cargar usuario"
          description={error?.message || "Usuario no encontrado"}
        />
        <div className="mt-4 flex justify-center">
          <Button onClick={() => navigate("/admin/users")}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Volver a la lista
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate("/admin/users")}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Volver
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">
              {user.first_name} {user.last_name}
            </h1>
            <p className="text-slate-600 mt-1">{user.email}</p>
          </div>
        </div>
      </div>

      {/* User Info */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="p-6 lg:col-span-2">
          <h2 className="text-xl font-semibold text-slate-900 mb-4">
            Información del Usuario
          </h2>
          <dl className="grid grid-cols-2 gap-4">
            <div>
              <dt className="text-sm font-medium text-slate-600">ID</dt>
              <dd className="text-sm text-slate-900 mt-1">{user.id}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">UUID</dt>
              <dd className="text-sm text-slate-900 mt-1 font-mono">{user.uuid}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Email</dt>
              <dd className="text-sm text-slate-900 mt-1 flex items-center gap-2">
                {user.email}
                {user.email_verified && (
                  <CheckCircle className="w-4 h-4 text-green-600" />
                )}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Teléfono</dt>
              <dd className="text-sm text-slate-900 mt-1">
                {user.phone || "No proporcionado"}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Cédula</dt>
              <dd className="text-sm text-slate-900 mt-1">
                {user.cedula || "No proporcionado"}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Rol</dt>
              <dd className="text-sm text-slate-900 mt-1">
                <Badge className="bg-blue-100 text-blue-700">{user.role}</Badge>
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Estado</dt>
              <dd className="text-sm text-slate-900 mt-1">
                <Badge className={
                  user.status === "active" ? "bg-green-100 text-green-700" :
                  user.status === "suspended" ? "bg-amber-100 text-amber-700" :
                  "bg-red-100 text-red-700"
                }>
                  {user.status}
                </Badge>
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Nivel KYC</dt>
              <dd className="text-sm text-slate-900 mt-1">
                <Badge className="bg-blue-100 text-blue-700">{user.kyc_level}</Badge>
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Fecha de Registro</dt>
              <dd className="text-sm text-slate-900 mt-1">
                {format(new Date(user.created_at), "dd/MM/yyyy HH:mm")}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Última Actualización</dt>
              <dd className="text-sm text-slate-900 mt-1">
                {format(new Date(user.updated_at), "dd/MM/yyyy HH:mm")}
              </dd>
            </div>
          </dl>
        </Card>

        {/* Actions */}
        <Card className="p-6">
          <h2 className="text-xl font-semibold text-slate-900 mb-4">
            Acciones Administrativas
          </h2>
          <div className="space-y-3">
            {user.status === "active" ? (
              <Button
                variant="outline"
                className="w-full justify-start text-amber-600 hover:text-amber-700 hover:bg-amber-50 hover:border-amber-600"
                onClick={() => handleUpdateStatus("suspended", "Suspendido por admin")}
              >
                <Ban className="w-4 h-4 mr-2" />
                Suspender Usuario
              </Button>
            ) : (
              <Button
                variant="outline"
                className="w-full justify-start text-green-600 hover:text-green-700 hover:bg-green-50 hover:border-green-600"
                onClick={() => handleUpdateStatus("active")}
              >
                <CheckCircle className="w-4 h-4 mr-2" />
                Activar Usuario
              </Button>
            )}

            <Button
              variant="outline"
              className="w-full justify-start text-blue-600 hover:text-blue-700 hover:bg-blue-50 hover:border-blue-600"
              onClick={() => {
                const levels: KYCLevel[] = ["none", "email_verified", "phone_verified", "cedula_verified", "full_kyc"];
                const currentIndex = levels.indexOf(user.kyc_level);
                const nextLevel = levels[Math.min(currentIndex + 1, levels.length - 1)];
                handleUpdateKYC(nextLevel);
              }}
            >
              <Shield className="w-4 h-4 mr-2" />
              Actualizar KYC
            </Button>
          </div>

          <div className="mt-6 pt-6 border-t border-slate-200">
            <p className="text-xs text-slate-500">
              Las acciones administrativas se registran en los logs de auditoría.
            </p>
          </div>
        </Card>
      </div>

      {/* Stats */}
      {user.raffle_stats && (
        <Card className="p-6">
          <h2 className="text-xl font-semibold text-slate-900 mb-4">
            Estadísticas de Rifas
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-sm font-medium text-slate-600">Total Rifas</p>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {user.raffle_stats.total_raffles}
              </p>
            </div>
            <div>
              <p className="text-sm font-medium text-slate-600">Rifas Activas</p>
              <p className="text-2xl font-bold text-blue-600 mt-1">
                {user.raffle_stats.active_raffles}
              </p>
            </div>
            <div>
              <p className="text-sm font-medium text-slate-600">Rifas Completadas</p>
              <p className="text-2xl font-bold text-green-600 mt-1">
                {user.raffle_stats.completed_raffles}
              </p>
            </div>
            <div>
              <p className="text-sm font-medium text-slate-600">Ingresos Totales</p>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                ${user.raffle_stats.total_revenue.toLocaleString()}
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* Payment Stats */}
      {user.payment_stats && (
        <Card className="p-6">
          <h2 className="text-xl font-semibold text-slate-900 mb-4">
            Estadísticas de Pagos
          </h2>
          <div className="grid grid-cols-3 gap-4">
            <div>
              <p className="text-sm font-medium text-slate-600">Total Pagos</p>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {user.payment_stats.total_payments}
              </p>
            </div>
            <div>
              <p className="text-sm font-medium text-slate-600">Total Gastado</p>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                ${user.payment_stats.total_spent.toLocaleString()}
              </p>
            </div>
            <div>
              <p className="text-sm font-medium text-slate-600">Refunds</p>
              <p className="text-2xl font-bold text-amber-600 mt-1">
                {user.payment_stats.refund_count}
              </p>
            </div>
          </div>
        </Card>
      )}
    </div>
  );
}
