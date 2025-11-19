import { useParams, useNavigate } from "react-router-dom";
import { ArrowLeft, AlertCircle, CheckCircle, XCircle, Percent } from "lucide-react";
import { useAdminOrganizerDetail, useVerifyOrganizer, useUpdateOrganizerCommission } from "../../hooks/useAdminOrganizers";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import { formatCurrency } from "@/lib/currency";

export function OrganizerDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const organizerId = parseInt(id || "0", 10);

  const { data: organizer, isLoading, error } = useAdminOrganizerDetail(organizerId);
  const verifyOrganizer = useVerifyOrganizer();
  const updateCommission = useUpdateOrganizerCommission();

  const handleToggleVerification = () => {
    if (!organizer) return;

    const newStatus = !organizer.profile.verified;
    const action = newStatus ? "verificar" : "remover verificación de";

    if (!confirm(`¿Confirmas ${action} este organizador?`)) return;

    verifyOrganizer.mutate({
      organizerId,
      data: {
        verified: newStatus,
        notes: newStatus ? "Verificado por admin" : "Verificación removida por admin",
      },
    });
  };

  const handleUpdateCommission = () => {
    const newCommission = prompt("Ingrese la nueva comisión personalizada (0-100):");
    if (newCommission === null) return;

    const commissionValue = parseFloat(newCommission);
    if (isNaN(commissionValue) || commissionValue < 0 || commissionValue > 100) {
      alert("Comisión inválida. Debe estar entre 0 y 100.");
      return;
    }

    updateCommission.mutate({
      organizerId,
      data: { commission_override: commissionValue },
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>
    );
  }

  if (error || !organizer) {
    return (
      <div className="p-6">
        <EmptyState
          icon={<AlertCircle className="w-12 h-12 text-red-500" />}
          title="Error al cargar organizador"
          description={error?.message || "Organizador no encontrado"}
        />
        <div className="mt-4 flex justify-center">
          <Button onClick={() => navigate("/admin/organizers")}>
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
            onClick={() => navigate("/admin/organizers")}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Volver
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">
              {organizer.profile.business_name}
            </h1>
            <p className="text-slate-600 mt-1">
              {organizer.user.first_name} {organizer.user.last_name} · {organizer.user.email}
            </p>
          </div>
        </div>

        {organizer.profile.verified ? (
          <Badge className="bg-green-100 text-green-700">
            <CheckCircle className="w-4 h-4 mr-2" />
            Verificado
          </Badge>
        ) : (
          <Badge className="bg-amber-100 text-amber-700">
            <XCircle className="w-4 h-4 mr-2" />
            No Verificado
          </Badge>
        )}
      </div>

      {/* Profile Info & Actions */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="p-6 lg:col-span-2">
          <h2 className="text-xl font-semibold text-slate-900 mb-4">
            Información del Perfil
          </h2>
          <dl className="grid grid-cols-2 gap-4">
            <div>
              <dt className="text-sm font-medium text-slate-600">ID Usuario</dt>
              <dd className="text-sm text-slate-900 mt-1">{organizer.user.id}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Tax ID</dt>
              <dd className="text-sm text-slate-900 mt-1 font-mono">
                {organizer.profile.tax_id || "No proporcionado"}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Banco</dt>
              <dd className="text-sm text-slate-900 mt-1">
                {organizer.profile.bank_name || "No proporcionado"}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Cuenta Bancaria</dt>
              <dd className="text-sm text-slate-900 mt-1 font-mono">
                {organizer.profile.bank_account_number || "No proporcionado"}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Comisión Personalizada</dt>
              <dd className="text-sm text-slate-900 mt-1">
                {organizer.profile.commission_override
                  ? `${organizer.profile.commission_override}%`
                  : "Por defecto"}
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-slate-600">Frecuencia de Pago</dt>
              <dd className="text-sm text-slate-900 mt-1 capitalize">
                {organizer.profile.payout_schedule || "Mensual"}
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
            <Button
              variant="outline"
              className={`w-full justify-start ${
                organizer.profile.verified
                  ? "text-amber-600 hover:text-amber-700 hover:bg-amber-50 hover:border-amber-600"
                  : "text-green-600 hover:text-green-700 hover:bg-green-50 hover:border-green-600"
              }`}
              onClick={handleToggleVerification}
            >
              {organizer.profile.verified ? (
                <>
                  <XCircle className="w-4 h-4 mr-2" />
                  Remover Verificación
                </>
              ) : (
                <>
                  <CheckCircle className="w-4 h-4 mr-2" />
                  Verificar Organizador
                </>
              )}
            </Button>

            <Button
              variant="outline"
              className="w-full justify-start text-blue-600 hover:text-blue-700 hover:bg-blue-50 hover:border-blue-600"
              onClick={handleUpdateCommission}
            >
              <Percent className="w-4 h-4 mr-2" />
              Actualizar Comisión
            </Button>
          </div>

          <div className="mt-6 pt-6 border-t border-slate-200">
            <p className="text-xs text-slate-500">
              Las acciones administrativas se registran en los logs de auditoría.
            </p>
          </div>
        </Card>
      </div>

      {/* Revenue Breakdown */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-slate-900 mb-4">
          Desglose de Ingresos
        </h2>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <p className="text-sm font-medium text-slate-600">Ingresos Brutos</p>
            <p className="text-2xl font-bold text-slate-900 mt-1">
              {formatCurrency(organizer.revenue_breakdown.gross_revenue)}
            </p>
          </div>
          <div>
            <p className="text-sm font-medium text-slate-600">Comisión Plataforma</p>
            <p className="text-2xl font-bold text-red-600 mt-1">
              -{formatCurrency(organizer.revenue_breakdown.platform_fees)}
            </p>
          </div>
          <div>
            <p className="text-sm font-medium text-slate-600">Ingresos Netos</p>
            <p className="text-2xl font-bold text-green-600 mt-1">
              {formatCurrency(organizer.revenue_breakdown.net_revenue)}
            </p>
          </div>
        </div>
      </Card>

      {/* Stats */}
      <Card className="p-6">
        <h2 className="text-xl font-semibold text-slate-900 mb-4">
          Métricas de Rifas
        </h2>
        <div className="grid grid-cols-3 gap-4">
          <div>
            <p className="text-sm font-medium text-slate-600">Total Rifas</p>
            <p className="text-2xl font-bold text-slate-900 mt-1">
              {organizer.metrics.total_raffles}
            </p>
          </div>
          <div>
            <p className="text-sm font-medium text-slate-600">Pago Pendiente</p>
            <p className="text-2xl font-bold text-blue-600 mt-1">
              {formatCurrency(organizer.metrics.pending_payout)}
            </p>
          </div>
          <div>
            <p className="text-sm font-medium text-slate-600">Total Pagado</p>
            <p className="text-2xl font-bold text-green-600 mt-1">
              {formatCurrency(organizer.profile.total_payouts)}
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
