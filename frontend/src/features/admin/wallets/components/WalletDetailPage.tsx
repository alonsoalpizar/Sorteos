import { useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { format } from "date-fns";
import { es } from "date-fns/locale";
import {
  ArrowLeft,
  Wallet,
  RefreshCw,
  Lock,
  Unlock,
  History,
  ChevronLeft,
  ChevronRight,
  AlertCircle,
  User,
} from "lucide-react";
import {
  useAdminWalletDetails,
  useAdminWalletTransactions,
  useFreezeWallet,
  useUnfreezeWallet,
} from "../hooks/useAdminWallets";
import type { TransactionType, TransactionStatus } from "../types";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";

// Helper para formatear CRC
function formatCRC(amount: string): string {
  const num = parseFloat(amount);
  return new Intl.NumberFormat("es-CR", {
    style: "currency",
    currency: "CRC",
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(num);
}

// Helper para traducir tipos de transacción
function translateTransactionType(type: TransactionType): string {
  const types: Record<TransactionType, string> = {
    deposit: "Recarga",
    withdrawal: "Retiro",
    purchase: "Compra de boletos",
    refund: "Reembolso",
    prize_claim: "Premio ganado",
    settlement_payout: "Pago de liquidación",
    adjustment: "Ajuste",
  };
  return types[type] || type;
}

// Helper para traducir estados
function translateTransactionStatus(status: TransactionStatus): string {
  const statuses: Record<TransactionStatus, string> = {
    pending: "Pendiente",
    completed: "Completada",
    failed: "Fallida",
    reversed: "Revertida",
  };
  return statuses[status] || status;
}

// Helper para obtener badge color según estado
function getStatusBadgeClass(status: TransactionStatus): string {
  const classes: Record<TransactionStatus, string> = {
    pending: "bg-yellow-100 text-yellow-700",
    completed: "bg-green-100 text-green-700",
    failed: "bg-red-100 text-red-700",
    reversed: "bg-slate-100 text-slate-700",
  };
  return classes[status] || "bg-slate-100 text-slate-700";
}

export const WalletDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const walletId = parseInt(id || "0", 10);

  const [page, setPage] = useState(0);
  const [limit] = useState(20);
  const [showFreezeModal, setShowFreezeModal] = useState(false);
  const [freezeReason, setFreezeReason] = useState("");

  const { data: walletData, isLoading: walletLoading, error: walletError, refetch: refetchWallet } = useAdminWalletDetails(walletId);
  const { data: transactionsData, isLoading: transactionsLoading, refetch: refetchTransactions } = useAdminWalletTransactions({
    wallet_id: walletId,
    page,
    limit,
  });

  const freezeMutation = useFreezeWallet();
  const unfreezeMutation = useUnfreezeWallet();

  const handleFreeze = async () => {
    if (!freezeReason.trim()) {
      alert("Debes proporcionar una razón para congelar la billetera");
      return;
    }

    await freezeMutation.mutateAsync({
      wallet_id: walletId,
      reason: freezeReason,
    });

    setShowFreezeModal(false);
    setFreezeReason("");
    refetchWallet();
  };

  const handleUnfreeze = async () => {
    if (confirm("¿Estás seguro de que deseas descongelar esta billetera?")) {
      await unfreezeMutation.mutateAsync({ wallet_id: walletId });
      refetchWallet();
    }
  };

  if (walletError) {
    return (
      <div className="p-6">
        <Card className="p-6">
          <div className="text-center text-red-600">
            <AlertCircle className="w-12 h-12 mx-auto mb-4" />
            <p className="font-semibold mb-2">Error al cargar la billetera</p>
            <p className="text-sm">{walletError instanceof Error ? walletError.message : "Error desconocido"}</p>
            <Button variant="outline" size="sm" onClick={() => navigate("/admin/wallets")} className="mt-4">
              Volver a la lista
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  if (walletLoading) {
    return (
      <div className="p-6">
        <Card className="p-6">
          <div className="flex items-center justify-center py-8">
            <LoadingSpinner />
          </div>
        </Card>
      </div>
    );
  }

  if (!walletData) {
    return null;
  }

  const wallet = walletData.wallet;
  const isFrozen = wallet.status === "frozen";
  const totalPages = transactionsData ? Math.ceil(transactionsData.pagination.total / limit) : 0;
  const hasNextPage = page < totalPages - 1;
  const hasPreviousPage = page > 0;

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link to="/admin/wallets">
            <Button variant="ghost" size="sm">
              <ArrowLeft className="w-4 h-4" />
            </Button>
          </Link>
          <div className="p-3 bg-blue-100 rounded-lg">
            <Wallet className="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-900">Billetera #{wallet.id}</h1>
            <p className="text-sm text-slate-600">{wallet.user_email}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="ghost" size="sm" onClick={() => { refetchWallet(); refetchTransactions(); }}>
            <RefreshCw className="w-4 h-4 mr-2" />
            Actualizar
          </Button>
          {isFrozen ? (
            <Button
              variant="outline"
              size="sm"
              onClick={handleUnfreeze}
              disabled={unfreezeMutation.isPending}
            >
              <Unlock className="w-4 h-4 mr-2" />
              {unfreezeMutation.isPending ? "Procesando..." : "Descongelar"}
            </Button>
          ) : (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowFreezeModal(true)}
              className="text-red-600 border-red-300 hover:bg-red-50"
            >
              <Lock className="w-4 h-4 mr-2" />
              Congelar
            </Button>
          )}
        </div>
      </div>

      {/* Wallet Status Alert */}
      {isFrozen && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-center gap-2 text-red-900">
            <AlertCircle className="w-5 h-5" />
            <p className="font-semibold">Esta billetera está congelada</p>
          </div>
          <p className="text-sm text-red-800 mt-1">
            El usuario no puede realizar transacciones mientras la billetera esté congelada.
          </p>
        </div>
      )}

      {/* User Info */}
      <Card className="p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4 flex items-center gap-2">
          <User className="w-5 h-5" />
          Información del Usuario
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-slate-600">Nombre</p>
            <p className="font-medium text-slate-900">{wallet.user_name}</p>
          </div>
          <div>
            <p className="text-sm text-slate-600">Email</p>
            <p className="font-medium text-slate-900">{wallet.user_email}</p>
          </div>
          <div>
            <p className="text-sm text-slate-600">ID de Usuario</p>
            <p className="font-medium text-slate-900">{wallet.user_id}</p>
          </div>
          <div>
            <p className="text-sm text-slate-600">UUID Billetera</p>
            <p className="font-mono text-xs text-slate-900">{wallet.uuid}</p>
          </div>
        </div>
      </Card>

      {/* Balance Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-4">
          <p className="text-sm text-slate-600 mb-1">Saldo Disponible</p>
          <p className="text-2xl font-bold text-blue-600">
            {formatCRC(wallet.balance_available)}
          </p>
          <p className="text-xs text-slate-500 mt-1">No retirable</p>
        </Card>
        <Card className="p-4">
          <p className="text-sm text-slate-600 mb-1">Ganancias</p>
          <p className="text-2xl font-bold text-green-600">
            {formatCRC(wallet.earnings_balance)}
          </p>
          <p className="text-xs text-slate-500 mt-1">Retirable</p>
        </Card>
        <Card className="p-4">
          <p className="text-sm text-slate-600 mb-1">Pendiente</p>
          <p className="text-2xl font-bold text-yellow-600">
            {formatCRC(wallet.pending_balance)}
          </p>
          <p className="text-xs text-slate-500 mt-1">En proceso</p>
        </Card>
        <Card className="p-4">
          <p className="text-sm text-slate-600 mb-1">Total</p>
          <p className="text-2xl font-bold text-slate-900">
            {formatCRC(wallet.total_balance)}
          </p>
          <p className="text-xs text-slate-500 mt-1">Suma total</p>
        </Card>
      </div>

      {/* Wallet Info */}
      <Card className="p-6">
        <h2 className="text-lg font-semibold text-slate-900 mb-4">Información de la Billetera</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <p className="text-sm text-slate-600">Moneda</p>
            <p className="font-medium text-slate-900">{wallet.currency}</p>
          </div>
          <div>
            <p className="text-sm text-slate-600">Estado</p>
            <p className="font-medium text-slate-900">
              {wallet.status === "active" ? "Activa" : wallet.status === "frozen" ? "Congelada" : "Cerrada"}
            </p>
          </div>
          <div>
            <p className="text-sm text-slate-600">Fecha de Creación</p>
            <p className="font-medium text-slate-900">
              {format(new Date(wallet.created_at), "dd MMM yyyy, HH:mm", { locale: es })}
            </p>
          </div>
        </div>
      </Card>

      {/* Transaction History */}
      <div className="space-y-4">
        <h2 className="text-lg font-semibold text-slate-900 flex items-center gap-2">
          <History className="w-5 h-5" />
          Historial de Transacciones
        </h2>

        {transactionsLoading ? (
          <Card className="p-6">
            <div className="flex items-center justify-center py-8">
              <LoadingSpinner />
            </div>
          </Card>
        ) : transactionsData && transactionsData.transactions.length > 0 ? (
          <>
            {/* Desktop Table */}
            <Card className="overflow-hidden hidden md:block">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-slate-50 border-b border-slate-200">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                        Fecha
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                        Tipo
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                        Monto
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                        Estado
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-600 uppercase tracking-wider">
                        Saldo Después
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-slate-200">
                    {transactionsData.transactions.map((tx) => {
                      const isDebit = ["purchase", "withdrawal", "adjustment"].includes(tx.type);
                      const date = new Date(tx.created_at);

                      return (
                        <tr key={tx.id} className="hover:bg-slate-50">
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-600">
                            {format(date, "dd MMM yyyy, HH:mm", { locale: es })}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <p className="text-sm font-medium text-slate-900">
                              {translateTransactionType(tx.type)}
                            </p>
                            {tx.description && (
                              <p className="text-xs text-slate-500">{tx.description}</p>
                            )}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm">
                            <span
                              className={`font-semibold ${
                                isDebit ? "text-red-600" : "text-green-600"
                              }`}
                            >
                              {isDebit ? "-" : "+"} {formatCRC(tx.amount)}
                            </span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <span
                              className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusBadgeClass(
                                tx.status
                              )}`}
                            >
                              {translateTransactionStatus(tx.status)}
                            </span>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-600">
                            {formatCRC(tx.balance_after)}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            </Card>

            {/* Mobile Cards */}
            <div className="md:hidden space-y-3">
              {transactionsData.transactions.map((tx) => {
                const isDebit = ["purchase", "withdrawal", "adjustment"].includes(tx.type);
                const date = new Date(tx.created_at);

                return (
                  <Card key={tx.id} className="p-4">
                    <div className="flex items-start justify-between mb-2">
                      <div>
                        <p className="font-medium text-slate-900">{translateTransactionType(tx.type)}</p>
                        <p className="text-xs text-slate-500">
                          {format(date, "dd MMM yyyy, HH:mm", { locale: es })}
                        </p>
                        {tx.description && (
                          <p className="text-xs text-slate-600 mt-1">{tx.description}</p>
                        )}
                      </div>
                      <span
                        className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getStatusBadgeClass(
                          tx.status
                        )}`}
                      >
                        {translateTransactionStatus(tx.status)}
                      </span>
                    </div>
                    <div className="flex items-center justify-between mt-3 pt-3 border-t border-slate-200">
                      <span
                        className={`text-lg font-bold ${isDebit ? "text-red-600" : "text-green-600"}`}
                      >
                        {isDebit ? "-" : "+"} {formatCRC(tx.amount)}
                      </span>
                      <span className="text-sm text-slate-600">
                        Saldo: {formatCRC(tx.balance_after)}
                      </span>
                    </div>
                  </Card>
                );
              })}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <Card className="p-4">
                <div className="flex items-center justify-between">
                  <div className="text-sm text-slate-600">
                    Página {page + 1} de {totalPages}
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage((p) => p - 1)}
                      disabled={!hasPreviousPage}
                    >
                      <ChevronLeft className="w-4 h-4" />
                      Anterior
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage((p) => p + 1)}
                      disabled={!hasNextPage}
                    >
                      Siguiente
                      <ChevronRight className="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </Card>
            )}
          </>
        ) : (
          <Card className="p-6">
            <div className="text-center text-slate-600">
              <History className="w-12 h-12 mx-auto mb-4 text-slate-400" />
              <p className="font-semibold mb-2">No hay transacciones</p>
              <p className="text-sm">Esta billetera aún no tiene transacciones</p>
            </div>
          </Card>
        )}
      </div>

      {/* Freeze Modal */}
      {showFreezeModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <Card className="max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-slate-900 mb-4">Congelar Billetera</h3>
            <p className="text-sm text-slate-600 mb-4">
              Esta acción impedirá que el usuario realice transacciones hasta que se descongele la billetera.
            </p>
            <div className="mb-4">
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Razón para congelar <span className="text-red-600">*</span>
              </label>
              <textarea
                value={freezeReason}
                onChange={(e) => setFreezeReason(e.target.value)}
                placeholder="Ej: Actividad sospechosa detectada"
                rows={3}
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
              />
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                onClick={() => {
                  setShowFreezeModal(false);
                  setFreezeReason("");
                }}
                className="flex-1"
              >
                Cancelar
              </Button>
              <Button
                onClick={handleFreeze}
                disabled={freezeMutation.isPending || !freezeReason.trim()}
                className="flex-1 bg-red-600 hover:bg-red-700 text-white"
              >
                {freezeMutation.isPending ? "Procesando..." : "Congelar"}
              </Button>
            </div>
          </Card>
        </div>
      )}
    </div>
  );
};
