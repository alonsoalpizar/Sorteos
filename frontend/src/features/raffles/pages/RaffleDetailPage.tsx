import { useParams, useNavigate, Link } from 'react-router-dom';
import { useRaffleDetail, usePublishRaffle, useDeleteRaffle } from '../../../hooks/useRaffles';
import { useAuth } from '../../../hooks/useAuth';
import { NumberGrid } from '../components/NumberGrid';
import { Button } from '../../../components/ui/Button';
import { Badge } from '../../../components/ui/Badge';
import { Alert } from '../../../components/ui/Alert';
import { Card } from '../../../components/ui/Card';
import {
  formatCurrency,
  formatDateTime,
  getStatusColor,
  getStatusLabel,
  getDrawMethodLabel,
} from '../../../lib/utils';

export function RaffleDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();

  const { data, isLoading, error } = useRaffleDetail(id!, {
    includeNumbers: true,
    includeImages: true,
  });

  const publishMutation = usePublishRaffle();
  const deleteMutation = useDeleteRaffle();

  const isOwner = user && data?.raffle && user.id === data.raffle.user_id;
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';
  const canEdit = isOwner || isAdmin;

  const handlePublish = async () => {
    if (!id || !confirm('¿Estás seguro de publicar este sorteo?')) return;

    try {
      await publishMutation.mutateAsync(Number(id));
      alert('Sorteo publicado exitosamente');
    } catch (error) {
      alert(error instanceof Error ? error.message : 'Error al publicar sorteo');
    }
  };

  const handleDelete = async () => {
    if (!id || !confirm('¿Estás seguro de eliminar este sorteo? Esta acción no se puede deshacer.'))
      return;

    try {
      await deleteMutation.mutateAsync(Number(id));
      alert('Sorteo eliminado exitosamente');
      navigate('/raffles');
    } catch (error) {
      alert(error instanceof Error ? error.message : 'Error al eliminar sorteo');
    }
  };

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Alert variant="error">
          <p className="font-semibold">Error al cargar el sorteo</p>
          <p className="text-sm mt-1">
            {error instanceof Error ? error.message : 'Error desconocido'}
          </p>
        </Alert>
        <Link to="/raffles" className="mt-4 inline-block">
          <Button variant="outline">Volver al listado</Button>
        </Link>
      </div>
    );
  }

  if (isLoading || !data) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  const { raffle, numbers = [], available_count, reserved_count, sold_count } = data;

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Back button */}
      <Link to="/raffles" className="inline-flex items-center text-blue-600 hover:text-blue-700 mb-6">
        <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        Volver al listado
      </Link>

      {/* Header */}
      <div className="bg-white dark:bg-slate-800 rounded-lg shadow p-6 mb-8">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <h1 className="text-3xl font-bold text-slate-900 dark:text-white mb-2">
              {raffle.title}
            </h1>
            <Badge variant={getStatusColor(raffle.status)}>
              {getStatusLabel(raffle.status)}
            </Badge>
          </div>

          {/* Actions */}
          {canEdit && (
            <div className="flex gap-2">
              {raffle.status === 'draft' && (
                <>
                  <Link to={`/raffles/${id}/edit`}>
                    <Button variant="outline" size="sm">
                      Editar
                    </Button>
                  </Link>
                  <Button
                    size="sm"
                    onClick={handlePublish}
                    disabled={publishMutation.isPending}
                  >
                    Publicar
                  </Button>
                </>
              )}
              {(raffle.status === 'draft' || raffle.status === 'suspended') && raffle.sold_count === 0 && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleDelete}
                  disabled={deleteMutation.isPending}
                  className="text-red-600 border-red-600 hover:bg-red-50"
                >
                  Eliminar
                </Button>
              )}
            </div>
          )}
        </div>

        <p className="text-slate-600 dark:text-slate-300 mb-6">{raffle.description}</p>

        {/* Stats grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Card className="p-4">
            <p className="text-sm text-slate-600 dark:text-slate-400">Precio por número</p>
            <p className="text-2xl font-bold text-blue-600">
              {formatCurrency(Number(raffle.price_per_number))}
            </p>
          </Card>

          <Card className="p-4">
            <p className="text-sm text-slate-600 dark:text-slate-400">Disponibles</p>
            <p className="text-2xl font-bold text-green-600">{available_count}</p>
          </Card>

          <Card className="p-4">
            <p className="text-sm text-slate-600 dark:text-slate-400">Vendidos</p>
            <p className="text-2xl font-bold text-slate-900 dark:text-white">{sold_count}</p>
          </Card>

          <Card className="p-4">
            <p className="text-sm text-slate-600 dark:text-slate-400">Reservados</p>
            <p className="text-2xl font-bold text-yellow-600">{reserved_count}</p>
          </Card>
        </div>

        {/* Info */}
        <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-slate-600 dark:text-slate-400">Fecha del sorteo:</span>
            <span className="ml-2 text-slate-900 dark:text-white font-medium">
              {formatDateTime(raffle.draw_date)}
            </span>
          </div>
          <div>
            <span className="text-slate-600 dark:text-slate-400">Método:</span>
            <span className="ml-2 text-slate-900 dark:text-white font-medium">
              {getDrawMethodLabel(raffle.draw_method)}
            </span>
          </div>
          <div>
            <span className="text-slate-600 dark:text-slate-400">Total de números:</span>
            <span className="ml-2 text-slate-900 dark:text-white font-medium">
              {raffle.total_numbers}
            </span>
          </div>
          <div>
            <span className="text-slate-600 dark:text-slate-400">Recaudación total:</span>
            <span className="ml-2 text-slate-900 dark:text-white font-medium">
              {formatCurrency(Number(raffle.total_revenue))}
            </span>
          </div>
        </div>
      </div>

      {/* Numbers grid */}
      {numbers.length > 0 && (
        <div className="bg-white dark:bg-slate-800 rounded-lg shadow p-6">
          <h2 className="text-xl font-bold text-slate-900 dark:text-white mb-6">
            Números del Sorteo
          </h2>
          <NumberGrid numbers={numbers} readonly />
        </div>
      )}
    </div>
  );
}
