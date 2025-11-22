import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useCreateRaffle } from '../../../hooks/useRaffles';
import { useCategories } from '../../../hooks/useCategories';
import { Button } from '../../../components/ui/Button';
import { Input } from '../../../components/ui/Input';
import { Label } from '../../../components/ui/Label';
import { Alert } from '../../../components/ui/Alert';
import type { CreateRaffleInput, DrawMethod } from '../../../types/raffle';

export function CreateRafflePage() {
  const navigate = useNavigate();
  const createMutation = useCreateRaffle();
  const { data: categories, isLoading: categoriesLoading } = useCategories();

  const [formData, setFormData] = useState<CreateRaffleInput>({
    title: '',
    description: '',
    price_per_number: 0,
    total_numbers: 100,
    draw_date: '',
    draw_method: 'loteria_nacional_cr',
    category_id: undefined,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.title || formData.title.length < 5) {
      newErrors.title = 'El título debe tener al menos 5 caracteres';
    }

    if (!formData.description || formData.description.length < 20) {
      newErrors.description = 'La descripción debe tener al menos 20 caracteres';
    }

    if (!formData.price_per_number || formData.price_per_number <= 0) {
      newErrors.price_per_number = 'El precio debe ser mayor a 0';
    }

    if (!formData.total_numbers || formData.total_numbers < 10 || formData.total_numbers > 10000) {
      newErrors.total_numbers = 'El total de números debe estar entre 10 y 10,000';
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

    if (!validate()) return;

    try {
      // Convertir fecha local a ISO 8601 (RFC3339) con timezone
      const drawDate = new Date(formData.draw_date);
      const isoDate = drawDate.toISOString();

      const payload = {
        ...formData,
        draw_date: isoDate,
      };

      const result = await createMutation.mutateAsync(payload);
      alert('Sorteo creado exitosamente');
      navigate(`/raffles/${result.raffle.id}`);
    } catch (error) {
      alert(error instanceof Error ? error.message : 'Error al crear sorteo');
    }
  };

  const handleChange = (
    field: keyof CreateRaffleInput,
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
        <Link to="/organizer" className="inline-flex items-center text-teal-600 hover:text-teal-700 mb-4">
          <svg className="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Volver al panel
        </Link>

        <h1 className="text-3xl font-bold text-slate-900 dark:text-white">
          Crear Nuevo Sorteo
        </h1>
        <p className="text-slate-600 dark:text-slate-400 mt-2">
          Completa la información del sorteo. Podrás editarla antes de publicarlo.
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
            className="w-full px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-teal-500 min-h-[120px]"
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
            className="w-full px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-teal-500"
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

        {/* Price and Total Numbers */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <Label htmlFor="price_per_number">
              Precio por Número <span className="text-red-500">*</span>
            </Label>
            <Input
              id="price_per_number"
              type="number"
              placeholder="5000"
              value={formData.price_per_number || ''}
              onChange={(e) => handleChange('price_per_number', Number(e.target.value))}
              error={errors.price_per_number}
              min="1"
              step="1"
            />
          </div>

          <div>
            <Label htmlFor="total_numbers">
              Total de Números <span className="text-red-500">*</span>
            </Label>
            <Input
              id="total_numbers"
              type="number"
              placeholder="100"
              value={formData.total_numbers || ''}
              onChange={(e) => handleChange('total_numbers', Number(e.target.value))}
              error={errors.total_numbers}
              min="10"
              max="10000"
              step="1"
            />
            <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
              Entre 10 y 10,000
            </p>
          </div>
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
              className="w-full px-4 py-2 rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-teal-500"
              value={formData.draw_method}
              onChange={(e) => handleChange('draw_method', e.target.value as DrawMethod)}
            >
              <option value="loteria_nacional_cr">Lotería Nacional CR</option>
              <option value="manual">Sorteo Manual</option>
              <option value="random">Sorteo Aleatorio</option>
            </select>
          </div>
        </div>

        {/* Info Alert */}
        <Alert variant="info">
          <p className="font-semibold">Importante</p>
          <p className="text-sm mt-1">
            El sorteo se creará en estado "borrador". Podrás editarlo y agregar imágenes antes de publicarlo.
            Una vez publicado, solo podrás modificar la descripción y la fecha.
          </p>
        </Alert>

        {/* Summary */}
        {formData.total_numbers > 0 && formData.price_per_number > 0 && (
          <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
            <p className="text-sm font-medium text-slate-900 dark:text-white mb-2">
              Resumen
            </p>
            <div className="space-y-1 text-sm text-slate-600 dark:text-slate-400">
              <p>
                Total de números: <strong>{formData.total_numbers}</strong>
              </p>
              <p>
                Precio por número: <strong>₡{formData.price_per_number.toLocaleString()}</strong>
              </p>
              <p className="text-base font-semibold text-blue-600 pt-2">
                Recaudación potencial: ₡{(formData.total_numbers * formData.price_per_number).toLocaleString()}
              </p>
              <p className="text-xs text-slate-500 dark:text-slate-400">
                (Comisión de plataforma 10%: ₡{((formData.total_numbers * formData.price_per_number) * 0.1).toLocaleString()})
              </p>
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-4">
          <Button
            type="submit"
            disabled={createMutation.isPending}
            className="flex-1"
          >
            {createMutation.isPending ? 'Creando...' : 'Crear Sorteo'}
          </Button>
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate('/organizer')}
          >
            Cancelar
          </Button>
        </div>
      </form>
    </div>
  );
}
