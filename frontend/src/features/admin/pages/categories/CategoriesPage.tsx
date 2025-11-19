import { useState } from "react";
import { Package, AlertCircle, Plus, Edit, Trash2, GripVertical } from "lucide-react";
import {
  useAdminCategories,
  useCreateCategory,
  useUpdateCategory,
  useDeleteCategory,
} from "../../hooks/useAdminCategories";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/Table";
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from "../../types";
import { format } from "date-fns";

export function CategoriesPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [isActiveFilter, setIsActiveFilter] = useState<boolean | undefined>(undefined);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);

  const { data, isLoading, error } = useAdminCategories(
    {
      search: search || undefined,
      is_active: isActiveFilter,
      order_by: "name",
    },
    { page, limit: 50 }
  );

  const createMutation = useCreateCategory();
  const updateMutation = useUpdateCategory();
  const deleteMutation = useDeleteCategory();

  const handleCreate = (formData: CreateCategoryRequest) => {
    createMutation.mutate(formData, {
      onSuccess: () => {
        setShowCreateModal(false);
      },
    });
  };

  const handleUpdate = (categoryId: number, formData: UpdateCategoryRequest) => {
    updateMutation.mutate(
      { categoryId, data: formData },
      {
        onSuccess: () => {
          setEditingCategory(null);
        },
      }
    );
  };

  const handleDelete = (category: Category) => {
    if (
      !confirm(
        `¿Confirmas eliminar la categoría "${category.name}"?\n\nEsta acción no se puede deshacer.`
      )
    )
      return;

    deleteMutation.mutate(category.id);
  };

  const handleToggleActive = (category: Category) => {
    updateMutation.mutate({
      categoryId: category.id,
      data: { is_active: !category.is_active },
    });
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">Gestión de Categorías</h1>
          <p className="text-slate-600 mt-2">
            Administra las categorías de rifas del sistema
          </p>
        </div>
        <Button
          onClick={() => setShowCreateModal(true)}
          className="bg-blue-600 hover:bg-blue-700"
        >
          <Plus className="w-4 h-4 mr-2" />
          Nueva Categoría
        </Button>
      </div>

      {/* Filters */}
      <Card className="p-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Buscar
            </label>
            <Input
              placeholder="Buscar por nombre o descripción..."
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
                setPage(1);
              }}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Estado
            </label>
            <select
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={
                isActiveFilter === undefined ? "" : isActiveFilter ? "true" : "false"
              }
              onChange={(e) => {
                const value = e.target.value;
                setIsActiveFilter(
                  value === "" ? undefined : value === "true"
                );
                setPage(1);
              }}
            >
              <option value="">Todos</option>
              <option value="true">Activas</option>
              <option value="false">Inactivas</option>
            </select>
          </div>
        </div>
      </Card>

      {/* Table */}
      <Card>
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <LoadingSpinner />
          </div>
        ) : error ? (
          <div className="p-6">
            <EmptyState
              icon={<AlertCircle className="w-12 h-12 text-red-500" />}
              title="Error al cargar categorías"
              description={(error as Error).message}
            />
          </div>
        ) : !data || !data.data || data.data.length === 0 ? (
          <div className="p-6">
            <EmptyState
              icon={<Package className="w-12 h-12 text-slate-400" />}
              title="No se encontraron categorías"
              description="Intenta ajustar los filtros de búsqueda"
            />
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-12"> </TableHead>
                  <TableHead>Nombre</TableHead>
                  <TableHead>Descripción</TableHead>
                  <TableHead>Icono</TableHead>
                  <TableHead>Rifas</TableHead>
                  <TableHead>Estado</TableHead>
                  <TableHead>Fecha Creación</TableHead>
                  <TableHead className="text-center">Acciones</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.data.map((category) => (
                  <TableRow key={category.id}>
                    <TableCell>
                      <GripVertical className="w-4 h-4 text-slate-400 cursor-move" />
                    </TableCell>
                    <TableCell className="font-medium text-slate-900">
                      {category.name}
                    </TableCell>
                    <TableCell className="text-sm text-slate-600 max-w-md truncate">
                      {category.description || "—"}
                    </TableCell>
                    <TableCell className="text-sm text-slate-600">
                      {category.icon || "—"}
                    </TableCell>
                    <TableCell>
                      <Badge className="bg-blue-100 text-blue-700">
                        {category.raffle_count || 0} rifas
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge
                        className={
                          category.is_active
                            ? "bg-green-100 text-green-700"
                            : "bg-gray-100 text-gray-700"
                        }
                      >
                        {category.is_active ? "Activa" : "Inactiva"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-sm text-slate-600">
                      {category.created_at
                        ? format(new Date(category.created_at), "dd/MM/yyyy")
                        : "—"}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-center gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => setEditingCategory(category)}
                        >
                          <Edit className="w-4 h-4" />
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleToggleActive(category)}
                          className={
                            category.is_active
                              ? "text-amber-600 hover:text-amber-700"
                              : "text-green-600 hover:text-green-700"
                          }
                        >
                          {category.is_active ? "Desactivar" : "Activar"}
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleDelete(category)}
                          className="text-red-600 hover:text-red-700"
                        >
                          <Trash2 className="w-4 h-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {/* Pagination */}
            {data.pagination && data.pagination.total_pages > 1 && (
              <div className="flex items-center justify-between px-6 py-4 border-t border-slate-200">
                <p className="text-sm text-slate-600">
                  Mostrando {data.data.length} de {data.pagination.total} categorías
                </p>
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                  >
                    Anterior
                  </Button>
                  <span className="px-4 py-2 text-sm text-slate-700">
                    Página {page} de {data.pagination.total_pages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => p + 1)}
                    disabled={page >= data.pagination.total_pages}
                  >
                    Siguiente
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
      </Card>

      {/* Create/Edit Modal */}
      {(showCreateModal || editingCategory) && (
        <CategoryFormModal
          category={editingCategory}
          onClose={() => {
            setShowCreateModal(false);
            setEditingCategory(null);
          }}
          onCreate={handleCreate}
          onUpdate={handleUpdate}
          isCreating={createMutation.isPending}
          isUpdating={updateMutation.isPending}
        />
      )}
    </div>
  );
}

// Category Form Modal Component
interface CategoryFormModalProps {
  category?: Category | null;
  onClose: () => void;
  onCreate: (data: CreateCategoryRequest) => void;
  onUpdate: (categoryId: number, data: UpdateCategoryRequest) => void;
  isCreating: boolean;
  isUpdating: boolean;
}

function CategoryFormModal({
  category,
  onClose,
  onCreate,
  onUpdate,
  isCreating,
  isUpdating,
}: CategoryFormModalProps) {
  const [formData, setFormData] = useState({
    name: category?.name || "",
    description: category?.description || "",
    icon: category?.icon || "",
    is_active: category?.is_active ?? true,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (category) {
      // Update
      const updateData: UpdateCategoryRequest = {};
      if (formData.name !== category.name) updateData.name = formData.name;
      if (formData.description !== category.description)
        updateData.description = formData.description;
      if (formData.icon !== category.icon) updateData.icon_url = formData.icon;
      if (formData.is_active !== category.is_active)
        updateData.is_active = formData.is_active;

      onUpdate(category.id, updateData);
    } else {
      // Create
      onCreate({
        name: formData.name,
        description: formData.description || undefined,
        icon_url: formData.icon || undefined,
        is_active: formData.is_active,
      });
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <Card className="w-full max-w-md">
        <form onSubmit={handleSubmit}>
          <div className="p-6 border-b border-slate-200">
            <h2 className="text-xl font-semibold text-slate-900">
              {category ? "Editar Categoría" : "Nueva Categoría"}
            </h2>
          </div>

          <div className="p-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Nombre *
              </label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Ej: Electrónica"
                required
                maxLength={100}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Descripción
              </label>
              <textarea
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={3}
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                placeholder="Descripción de la categoría"
                maxLength={500}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Icono (emoji o URL)
              </label>
              <Input
                value={formData.icon}
                onChange={(e) => setFormData({ ...formData, icon: e.target.value })}
                placeholder="Ej: 📱 o https://..."
              />
            </div>

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="is_active"
                checked={formData.is_active}
                onChange={(e) =>
                  setFormData({ ...formData, is_active: e.target.checked })
                }
                className="w-4 h-4 text-blue-600 border-slate-300 rounded focus:ring-blue-500"
              />
              <label htmlFor="is_active" className="text-sm font-medium text-slate-700">
                Categoría activa
              </label>
            </div>
          </div>

          <div className="p-6 border-t border-slate-200 flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancelar
            </Button>
            <Button
              type="submit"
              className="bg-blue-600 hover:bg-blue-700"
              disabled={isCreating || isUpdating}
            >
              {isCreating || isUpdating ? "Guardando..." : category ? "Actualizar" : "Crear"}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
