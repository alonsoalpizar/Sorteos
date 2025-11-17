import { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { useRaffleDetail, useUpdateRaffle } from '../../../hooks/useRaffles';
import { useCategories } from '../../../hooks/useCategories';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';
import { Label } from '../../../components/ui/Label';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';
import { Alert } from '../../../components/ui/Alert';
import { ImageUploader } from '../../../components/ImageUploader';
import type { UpdateRaffleInput, DrawMethod } from '../../../types/raffle';

export function EditRafflePage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const updateMutation = useUpdateRaffle();
  const { data: categories, isLoading: categoriesLoading } = useCategories();

  const { data, isLoading, error } = useRaffleDetail(id!, {
    includeNumbers: false,
    includeImages: true,
  });

  const [formData, setFormData] = useState<UpdateRaffleInput>({
    title: '',
    description: '',
    draw_date: '',
    draw_method: 'loteria_nacional_cr',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  // Cargar datos del sorteo cuando se obtienen
  useEffect(() => {
    if (data?.raffle) {
      const drawDate = new Date(data.raffle.draw_date);
      // Convertir a formato datetime-local (YYYY-MM-DDTHH:mm)
      const localDate = new Date(drawDate.getTime() - drawDate.getTimezoneOffset() * 60000)
        .toISOString()
        .slice(0, 16);

      setFormData({
        title: data.raffle.title,
        description: data.raffle.description,
        draw_date: localDate,
        draw_method: data.raffle.draw_method as DrawMethod,
        category_id: data.raffle.category_id,
      });
    }
  }, [data]);

  if (isLoading) {
    return <LoadingSpinner text="Cargando sorteo..." />;
  }

  if (error || !data) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">Error al cargar el sorteo</p>
      </div>
    );
  }

  const raffle = data.raffle;

  // Solo el dueño puede editar y solo si está en draft
  if (raffle.status !== 'draft') {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <Alert variant="warning">
          <p className="font-semibold">No se puede editar</p>
          <p className="text-sm mt-1">
            Solo se pueden editar sorteos en estado borrador. Este sorteo ya ha sido publicado.
          </p>
        </Alert>
        <Button onClick={() => navigate(`/raffles/${id}`)} className="mt-4">
          Volver al sorteo
        </Button>
      </div>
    );
  }

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.title || formData.title.length < 5) {
      newErrors.title = 'El título debe tener al menos 5 caracteres';
    }

    if (!formData.description || formData.description.length < 20) {
      newErrors.description = 'La descripción debe tener al menos 20 caracteres';
    }

    if (!formData.draw_date) {
      newErrors.draw_date = 'La fecha del sorteo es requerida';
    } else {
      const drawDate = new Date(formData.draw_date);
      const now = new Date();
      if (drawDate <= now) {
        newErrors.draw_date = 'La fecha del sorteo debe ser en el futuro';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate() || !data?.raffle) return;

    try {
      // Convertir fecha local a ISO 8601 (RFC3339) con timezone
      const drawDate = new Date(formData.draw_date!);
      const isoDate = drawDate.toISOString();

      const payload: UpdateRaffleInput = {
        ...formData,
        draw_date: isoDate,
      };

      // Usar el ID numérico del raffle, no el parámetro de la URL (que puede ser UUID)
      await updateMutation.mutateAsync({ id: data.raffle.id, input: payload });

      alert('Sorteo actualizado exitosamente');

      // Navegar usando el ID o UUID del sorteo actualizado
      // Forzar recarga de la página para que se vean los cambios
      window.location.href = `/raffles/${id}`;
    } catch (error) {
      alert(error instanceof Error ? error.message : 'Error al actualizar sorteo');
    }
  };

  const handleChange = (
    field: keyof UpdateRaffleInput,
    value: string | number | undefined
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error for this field
    if (errors[field]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      {/* Header */}
      <div className="mb-8">
        <Link to={`/raffles/${id}`} className="inline-flex items-center text-blue-600 hover:text-blue-700 mb-4">
          <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Volver al sorteo
        </Link>

        <h1 className="text-3xl font-bold text-slate-900 dark:text-white">
          Editar Sorteo
        </h1>
        <p className="text-slate-600 dark:text-slate-400 mt-2">
          Modifica la información del sorteo antes de publicarlo.
        </p>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="bg-white dark:bg-slate-800 rounded-lg shadow p-6 space-y-6">
        {/* Title */}
        <div>
          <Label htmlFor="title">
            Título del Sorteo <span className="text-red-500">*</span>
          </Label>
          <Input
            id="title"
            type="text"
            placeholder="Ej: Sorteo iPhone 15 Pro Max 256GB"
            value={formData.title}
            onChange={(e) => handleChange('title', e.target.value)}
            error={errors.title}
            maxLength={255}
          />
          <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
            Mínimo 5 caracteres, máximo 255
          </p>
        </div>

        {/* Description */}
        <div>
          <Label htmlFor="description">
            Descripción <span className="text-red-500">*</span>
          </Label>
          <textarea
            id="description"
            className="w-full px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500 min-h-[120px]"
            placeholder="Describe detalladamente el premio y las condiciones del sorteo..."
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
          />
          {errors.description && (
            <p className="text-sm text-red-500 mt-1">{errors.description}</p>
          )}
          <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
            Mínimo 20 caracteres
          </p>
        </div>

        {/* Category */}
        <div>
          <Label htmlFor="category_id">
            Categoría <span className="text-red-500">*</span>
          </Label>
          <select
            id="category_id"
            className="w-full px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={formData.category_id || ''}
            onChange={(e) => handleChange('category_id', e.target.value ? Number(e.target.value) : undefined)}
            disabled={categoriesLoading}
          >
            <option value="">Selecciona una categoría</option>
            {categories?.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.icon} {cat.name}
              </option>
            ))}
          </select>
          {errors.category_id && (
            <p className="text-sm text-red-500 mt-1">{errors.category_id}</p>
          )}
          <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
            Ayuda a los usuarios a encontrar tu sorteo
          </p>
        </div>

        {/* Draw Date and Method */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <Label htmlFor="draw_date">
              Fecha del Sorteo <span className="text-red-500">*</span>
            </Label>
            <Input
              id="draw_date"
              type="datetime-local"
              value={formData.draw_date}
              onChange={(e) => handleChange('draw_date', e.target.value)}
              error={errors.draw_date}
            />
          </div>

          <div>
            <Label htmlFor="draw_method">
              Método de Sorteo <span className="text-red-500">*</span>
            </Label>
            <select
              id="draw_method"
              className="w-full px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={formData.draw_method}
              onChange={(e) => handleChange('draw_method', e.target.value as DrawMethod)}
            >
              <option value="loteria_nacional_cr">Lotería Nacional CR</option>
              <option value="manual">Sorteo Manual</option>
              <option value="random">Sorteo Aleatorio</option>
            </select>
          </div>
        </div>

        {/* Read-only fields info */}
        <Alert variant="info">
          <p className="font-semibold">Nota</p>
          <p className="text-sm mt-1">
            El precio por número y la cantidad total de números no se pueden modificar una vez creado el sorteo.
            Estas son: <strong>{raffle.total_numbers} números</strong> a <strong>₡{Number(raffle.price_per_number).toLocaleString()}</strong> cada uno.
          </p>
        </Alert>

        {/* Images Section */}
        <div className="border-t pt-6">
          <h3 className="text-lg font-semibold text-slate-900 dark:text-white mb-4">
            Imágenes del Sorteo
          </h3>
          <p className="text-sm text-slate-600 dark:text-slate-400 mb-4">
            Agrega hasta 5 imágenes. La primera imagen será la principal por defecto.
          </p>
          <ImageUploader
            raffleId={raffle.id}
            maxImages={5}
          />
        </div>

        {/* Actions */}
        <div className="flex gap-4">
          <Button
            type="submit"
            disabled={updateMutation.isPending}
            className="flex-1"
          >
            {updateMutation.isPending ? 'Guardando...' : 'Guardar Cambios'}
          </Button>
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate(`/raffles/${id}`)}
          >
            Cancelar
          </Button>
        </div>
      </form>
    </div>
  );
}
