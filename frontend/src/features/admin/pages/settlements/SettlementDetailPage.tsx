import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  ArrowLeft,
  DollarSign,
  User,
  Package,
  Clock,
  AlertCircle,
  CheckCircle,
  XCircle,
  CreditCard,
} from "lucide-react";
import {
  useAdminSettlementDetail,
  useApproveSettlement,
  useRejectSettlement,
  useMarkSettlementPaid,
} from "../../hooks/useAdminSettlements";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import { formatCurrency } from "@/lib/currency";
import { format } from "date-fns";
import type {
  ApproveSettlementRequest,
  RejectSettlementRequest,
  MarkSettlementPaidRequest,
} from "../../types";

export function SettlementDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const settlementId = parseInt(id || "0");

  const [showApproveModal, setShowApproveModal] = useState(false);
  const [showRejectModal, setShowRejectModal] = useState(false);
  const [showPayoutModal, setShowPayoutModal] = useState(false);

  const { data, isLoading, error } = useAdminSettlementDetail(settlementId);
  const approveMutation = useApproveSettlement();
  const rejectMutation = useRejectSettlement();
  const payoutMutation = useMarkSettlementPaid();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="p-6">
        <EmptyState
          icon={<AlertCircle className="w-12 h-12 text-red-500" />}
          title="Error al cargar liquidación"
          description={(error as Error)?.message || "Liquidación no encontrada"}
        />
      </div>
    );
  }

  const { settlement, raffle, payments_summary, timeline, bank_account } = data;

  const getStatusBadge = (status: string) => {
    const styles: Record<string, string> = {
      pending: "bg-yellow-100 text-yellow-700",
      approved: "bg-blue-100 text-blue-700",
      paid: "bg-green-100 text-green-700",
      rejected: "bg-red-100 text-red-700",
    };
    const labels: Record<string, string> = {
      pending: "Pendiente",
      approved: "Aprobada",
      paid: "Pagada",
      rejected: "Rechazada",
    };
    return (
      <Badge className={styles[status] || "bg-slate-100 text-slate-700"}>
        {labels[status] || status}
      </Badge>
    );
  };

  const getEventIcon = (type: string) => {
    switch (type) {
      case "calculated":
        return <Clock className="w-5 h-5 text-blue-600" />;
      case "approved":
        return <CheckCircle className="w-5 h-5 text-green-600" />;
      case "rejected":
        return <XCircle className="w-5 h-5 text-red-600" />;
      case "paid":
        return <DollarSign className="w-5 h-5 text-green-600" />;
      default:
        return <Clock className="w-5 h-5 text-slate-600" />;
    }
  };

  const canApprove = settlement.status === "pending";
  const canReject = settlement.status === "pending" || settlement.status === "approved";
  const canMarkPaid = settlement.status === "approved";

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="outline" size="sm" onClick={() => navigate("/admin/settlements")}>
            <ArrowLeft className="w-4 h-4" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-slate-900">Liquidación #{settlement.id}</h1>
            <p className="text-slate-600 mt-1">
              Detalles completos de la liquidación
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          {canApprove && (
            <Button
              onClick={() => setShowApproveModal(true)}
              className="bg-green-600 hover:bg-green-700"
            >
              <CheckCircle className="w-4 h-4 mr-2" />
              Aprobar
            </Button>
          )}
          {canReject && (
            <Button
              variant="outline"
              onClick={() => setShowRejectModal(true)}
              className="text-red-600 border-red-600 hover:bg-red-50"
            >
              <XCircle className="w-4 h-4 mr-2" />
              Rechazar
            </Button>
          )}
          {canMarkPaid && (
            <Button
              onClick={() => setShowPayoutModal(true)}
              className="bg-blue-600 hover:bg-blue-700"
            >
              <DollarSign className="w-4 h-4 mr-2" />
              Marcar Como Pagada
            </Button>
          )}
        </div>
      </div>

      {/* Main Info Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Settlement Info */}
        <Card className="p-6">
          <div className="flex items-center gap-3 mb-4">
            <DollarSign className="w-6 h-6 text-blue-600" />
            <h2 className="text-lg font-semibold text-slate-900">Liquidación</h2>
          </div>
          <div className="space-y-3">
            <div>
              <p className="text-sm text-slate-600">Estado</p>
              <div className="mt-1">{getStatusBadge(settlement.status)}</div>
            </div>
            <div>
              <p className="text-sm text-slate-600">Ingresos Totales</p>
              <p className="text-lg font-bold text-slate-900">
                {formatCurrency(settlement.total_revenue)}
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-600">Comisión Plataforma</p>
              <p className="text-lg font-medium text-red-600">
                - {formatCurrency(settlement.platform_fee)}
              </p>
            </div>
            <div className="pt-3 border-t border-slate-200">
              <p className="text-sm text-slate-600">Monto Neto a Pagar</p>
              <p className="text-2xl font-bold text-green-600">
                {formatCurrency(settlement.net_amount)}
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-600">Fecha de Creación</p>
              <p className="text-sm font-medium text-slate-900">
                {format(new Date(settlement.created_at), "dd/MM/yyyy HH:mm")}
              </p>
            </div>
          </div>
        </Card>

        {/* Organizer Info */}
        <Card className="p-6">
          <div className="flex items-center gap-3 mb-4">
            <User className="w-6 h-6 text-green-600" />
            <h2 className="text-lg font-semibold text-slate-900">Organizador</h2>
          </div>
          <div className="space-y-3">
            <div>
              <p className="text-sm text-slate-600">Nombre</p>
              <p className="text-sm font-medium text-slate-900">
                {settlement.organizer_name}
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-600">Email</p>
              <p className="text-sm font-medium text-slate-900">
                {settlement.organizer_email}
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-600">Nivel KYC</p>
              <Badge className={
                settlement.organizer_kyc_level === "verified" || settlement.organizer_kyc_level === "enhanced"
                  ? "bg-green-100 text-green-700"
                  : "bg-yellow-100 text-yellow-700"
              }>
                {settlement.organizer_kyc_level}
              </Badge>
            </div>
            {bank_account && (
              <>
                <div className="pt-3 border-t border-slate-200">
                  <p className="text-sm font-medium text-slate-700 mb-2">Cuenta Bancaria</p>
                  <p className="text-xs text-slate-600">Banco: {bank_account.bank_name}</p>
                  <p className="text-xs text-slate-600">
                    Cuenta: {bank_account.account_number}
                  </p>
                  <p className="text-xs text-slate-600">
                    Titular: {bank_account.account_holder}
                  </p>
                  {bank_account.verified_at && (
                    <Badge className="bg-green-100 text-green-700 mt-2">
                      Verificada
                    </Badge>
                  )}
                </div>
              </>
            )}
          </div>
        </Card>

        {/* Raffle Info */}
        <Card className="p-6">
          <div className="flex items-center gap-3 mb-4">
            <Package className="w-6 h-6 text-purple-600" />
            <h2 className="text-lg font-semibold text-slate-900">Rifa</h2>
          </div>
          <div className="space-y-3">
            <div>
              <p className="text-sm text-slate-600">Título</p>
              <p className="text-sm font-medium text-slate-900">{settlement.raffle_title}</p>
            </div>
            {raffle && (
              <>
                <div>
                  <p className="text-sm text-slate-600">Estado</p>
                  <Badge>{raffle.status}</Badge>
                </div>
                <div>
                  <p className="text-sm text-slate-600">Precio por Número</p>
                  <p className="text-sm font-medium text-slate-900">
                    {formatCurrency(raffle.ticket_price || 0)}
                  </p>
                </div>
              </>
            )}
          </div>
        </Card>
      </div>

      {/* Payments Summary */}
      {payments_summary && (
        <Card className="p-6">
          <div className="flex items-center gap-3 mb-6">
            <CreditCard className="w-6 h-6 text-blue-600" />
            <h2 className="text-lg font-semibold text-slate-900">Resumen de Pagos</h2>
          </div>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-sm text-slate-600">Total Pagos</p>
              <p className="text-2xl font-bold text-slate-900">
                {payments_summary.total_payments}
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-600">Exitosos</p>
              <p className="text-2xl font-bold text-green-600">
                {payments_summary.succeeded_payments}
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-600">Ingresos Brutos</p>
              <p className="text-xl font-bold text-slate-900">
                {formatCurrency(payments_summary.total_revenue)}
              </p>
            </div>
            <div>
              <p className="text-sm text-slate-600">
                Comisión ({payments_summary.platform_fee_percent}%)
              </p>
              <p className="text-xl font-bold text-red-600">
                {formatCurrency(payments_summary.platform_fee_amount)}
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* Timeline */}
      <Card className="p-6">
        <div className="flex items-center gap-3 mb-6">
          <Clock className="w-6 h-6 text-blue-600" />
          <h2 className="text-lg font-semibold text-slate-900">Timeline de Eventos</h2>
        </div>
        <div className="space-y-4">
          {timeline && timeline.length > 0 ? (
            timeline.map((event, index) => (
              <div key={index} className="flex gap-4">
                <div className="flex-shrink-0">{getEventIcon(event.type)}</div>
                <div className="flex-1">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-slate-900">{event.details}</p>
                      {event.actor && (
                        <p className="text-xs text-slate-500">Por: {event.actor}</p>
                      )}
                    </div>
                    <p className="text-xs text-slate-500">
                      {format(new Date(event.timestamp), "dd/MM/yyyy HH:mm")}
                    </p>
                  </div>
                  {event.metadata && Object.keys(event.metadata).length > 0 && (
                    <div className="mt-1 text-xs text-slate-600">
                      {Object.entries(event.metadata).map(([key, value]) => (
                        <span key={key} className="mr-3">
                          <span className="font-medium">{key}:</span>{" "}
                          {typeof value === "number" && key.includes("amount")
                            ? formatCurrency(value)
                            : String(value)}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ))
          ) : (
            <p className="text-sm text-slate-500">No hay eventos en el timeline</p>
          )}
        </div>
      </Card>

      {/* Admin Notes */}
      {settlement.admin_notes && (
        <Card className="p-6 bg-yellow-50 border-yellow-200">
          <div className="flex items-center gap-3">
            <AlertCircle className="w-6 h-6 text-yellow-600" />
            <div className="flex-1">
              <h3 className="text-sm font-semibold text-yellow-900">Notas Administrativas</h3>
              <p className="text-sm text-yellow-700 mt-1 whitespace-pre-wrap">
                {settlement.admin_notes}
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* Rejection Reason */}
      {settlement.rejection_reason && (
        <Card className="p-6 bg-red-50 border-red-200">
          <div className="flex items-center gap-3">
            <XCircle className="w-6 h-6 text-red-600" />
            <div>
              <h3 className="text-sm font-semibold text-red-900">Razón del Rechazo</h3>
              <p className="text-sm text-red-700 mt-1">{settlement.rejection_reason}</p>
            </div>
          </div>
        </Card>
      )}

      {/* Modals */}
      {showApproveModal && (
        <ApproveModal
          settlementId={settlementId}
          onClose={() => setShowApproveModal(false)}
          onSubmit={(data) => {
            approveMutation.mutate(
              { settlementId, data },
              {
                onSuccess: () => setShowApproveModal(false),
              }
            );
          }}
        />
      )}

      {showRejectModal && (
        <RejectModal
          settlementId={settlementId}
          onClose={() => setShowRejectModal(false)}
          onSubmit={(data) => {
            rejectMutation.mutate(
              { settlementId, data },
              {
                onSuccess: () => setShowRejectModal(false),
              }
            );
          }}
        />
      )}

      {showPayoutModal && (
        <PayoutModal
          settlementId={settlementId}
          netAmount={settlement.net_amount}
          onClose={() => setShowPayoutModal(false)}
          onSubmit={(data) => {
            payoutMutation.mutate(
              { settlementId, data },
              {
                onSuccess: () => setShowPayoutModal(false),
              }
            );
          }}
        />
      )}
    </div>
  );
}

// Approve Modal Component
interface ApproveModalProps {
  settlementId: number;
  onClose: () => void;
  onSubmit: (data: ApproveSettlementRequest) => void;
}

function ApproveModal({ onClose, onSubmit }: ApproveModalProps) {
  const [notes, setNotes] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ notes });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <Card className="w-full max-w-md p-6">
        <h2 className="text-xl font-bold text-slate-900 mb-4">Aprobar Liquidación</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Notas (Opcional)
            </label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={3}
              placeholder="Notas administrativas..."
            />
          </div>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p className="text-sm text-blue-900">
              <strong>Nota:</strong> Esta acción aprobará la liquidación. El organizador deberá tener:
            </p>
            <ul className="text-sm text-blue-800 mt-2 list-disc list-inside">
              <li>KYC nivel verificado o superior</li>
              <li>Cuenta bancaria verificada</li>
            </ul>
          </div>

          <div className="flex gap-3 justify-end">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" className="bg-green-600 hover:bg-green-700">
              Aprobar Liquidación
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}

// Reject Modal Component
interface RejectModalProps {
  settlementId: number;
  onClose: () => void;
  onSubmit: (data: RejectSettlementRequest) => void;
}

function RejectModal({ onClose, onSubmit }: RejectModalProps) {
  const [reason, setReason] = useState("");
  const [notes, setNotes] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!reason.trim()) {
      return;
    }
    onSubmit({ reason, notes });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <Card className="w-full max-w-md p-6">
        <h2 className="text-xl font-bold text-slate-900 mb-4">Rechazar Liquidación</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Razón del Rechazo *
            </label>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={3}
              required
              placeholder="Explica por qué se rechaza..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Notas Adicionales
            </label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={2}
              placeholder="Notas internas..."
            />
          </div>

          <div className="flex gap-3 justify-end">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancelar
            </Button>
            <Button
              type="submit"
              className="bg-red-600 hover:bg-red-700"
              disabled={!reason.trim()}
            >
              Rechazar Liquidación
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}

// Payout Modal Component
interface PayoutModalProps {
  settlementId: number;
  netAmount: number;
  onClose: () => void;
  onSubmit: (data: MarkSettlementPaidRequest) => void;
}

function PayoutModal({ netAmount, onClose, onSubmit }: PayoutModalProps) {
  const [paymentMethod, setPaymentMethod] = useState("bank_transfer");
  const [paymentReference, setPaymentReference] = useState("");
  const [notes, setNotes] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      payment_method: paymentMethod,
      payment_reference: paymentReference || undefined,
      notes: notes || undefined,
    });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <Card className="w-full max-w-md p-6">
        <h2 className="text-xl font-bold text-slate-900 mb-4">Marcar Como Pagada</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="bg-green-50 border border-green-200 rounded-lg p-4">
            <p className="text-sm text-green-900">
              Monto a pagar: <strong>{formatCurrency(netAmount)}</strong>
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Método de Pago *
            </label>
            <select
              value={paymentMethod}
              onChange={(e) => setPaymentMethod(e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option value="bank_transfer">Transferencia Bancaria</option>
              <option value="paypal">PayPal</option>
              <option value="stripe">Stripe</option>
              <option value="cash">Efectivo</option>
              <option value="check">Cheque</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Referencia de Pago
            </label>
            <input
              type="text"
              value={paymentReference}
              onChange={(e) => setPaymentReference(e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Número de transferencia, Transaction ID..."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Notas
            </label>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows={2}
              placeholder="Notas sobre el pago..."
            />
          </div>

          <div className="flex gap-3 justify-end">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancelar
            </Button>
            <Button type="submit" className="bg-blue-600 hover:bg-blue-700">
              Marcar Como Pagada
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
