import { useState } from "react";
import { Settings, Save, Database, Mail, Shield, Zap, DollarSign, RefreshCw, Plus, X } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { LoadingSpinner } from "@/components/ui/LoadingSpinner";
import { EmptyState } from "@/components/ui/EmptyState";
import { useSystemSettings, useUpdateSystemSetting } from "../../hooks/useAdminSystem";
import type { SystemSetting } from "../../types";

export function SystemConfigPage() {
  const [selectedCategory, setSelectedCategory] = useState<string | undefined>(undefined);
  const [editingKey, setEditingKey] = useState<string | null>(null);
  const [editValue, setEditValue] = useState<any>(null);
  const [editCategory, setEditCategory] = useState<string>("business");
  const [editValueType, setEditValueType] = useState<"string" | "int" | "float" | "bool" | "json">("string");
  const [editDescription, setEditDescription] = useState<string>("");

  // New parameter state
  const [showNewForm, setShowNewForm] = useState(false);
  const [newKey, setNewKey] = useState("");
  const [newValue, setNewValue] = useState("");
  const [newCategory, setNewCategory] = useState("business");
  const [newValueType, setNewValueType] = useState<"string" | "int" | "float" | "bool" | "json">("string");
  const [newDescription, setNewDescription] = useState("");

  const { data, isLoading, refetch } = useSystemSettings({ category: selectedCategory });
  const updateMutation = useUpdateSystemSetting();

  const handleEdit = (setting: SystemSetting) => {
    setEditingKey(setting.key);
    // Si el valor es un objeto, lo convertimos a JSON string para editar
    if (typeof setting.value === "object") {
      setEditValue(JSON.stringify(setting.value, null, 2));
    } else {
      setEditValue(setting.value);
    }
    // Cargar valores actuales para edición
    setEditCategory(setting.category || "business");
    setEditValueType(setting.value_type || "string");
    setEditDescription(setting.description || "");
  };

  const handleSave = async (setting: SystemSetting) => {
    let finalValue = editValue;

    // Intentar parsear como JSON si parece ser un objeto/array
    if (typeof editValue === "string") {
      const trimmed = editValue.trim();
      if ((trimmed.startsWith("{") && trimmed.endsWith("}")) ||
          (trimmed.startsWith("[") && trimmed.endsWith("]"))) {
        try {
          finalValue = JSON.parse(editValue);
        } catch (e) {
          // Si falla, dejarlo como string
        }
      } else if (trimmed === "true") {
        finalValue = true;
      } else if (trimmed === "false") {
        finalValue = false;
      } else if (!isNaN(Number(trimmed)) && trimmed !== "") {
        finalValue = Number(trimmed);
      }
    }

    await updateMutation.mutateAsync({
      key: setting.key,
      value: finalValue,
      category: editCategory,
      value_type: editValueType,
      description: editDescription || undefined,
    });

    setEditingKey(null);
    setEditValue(null);
    setEditCategory("business");
    setEditValueType("string");
    setEditDescription("");
  };

  const handleCancel = () => {
    setEditingKey(null);
    setEditValue(null);
    setEditCategory("business");
    setEditValueType("string");
    setEditDescription("");
  };

  const handleCreateNew = async () => {
    if (!newKey.trim() || !newValue.trim()) {
      return;
    }

    // El valor se envía como string, el backend lo guarda y valida según value_type
    await updateMutation.mutateAsync({
      key: newKey.trim(),
      value: newValue.trim(),
      category: newCategory,
      value_type: newValueType,
      description: newDescription.trim() || undefined,
    });

    // Reset form
    setNewKey("");
    setNewValue("");
    setNewCategory("business");
    setNewValueType("string");
    setNewDescription("");
    setShowNewForm(false);
  };

  const getCategoryIcon = (category: string) => {
    if (!category) return <Settings className="w-5 h-5 text-slate-600" />;

    switch (category.toLowerCase()) {
      case "database":
        return <Database className="w-5 h-5 text-blue-600" />;
      case "email":
      case "smtp":
        return <Mail className="w-5 h-5 text-blue-600" />;
      case "security":
      case "auth":
        return <Shield className="w-5 h-5 text-red-600" />;
      case "performance":
        return <Zap className="w-5 h-5 text-yellow-600" />;
      case "payment":
      case "billing":
        return <DollarSign className="w-5 h-5 text-green-600" />;
      default:
        return <Settings className="w-5 h-5 text-slate-600" />;
    }
  };

  const renderValue = (value: any, isEditing: boolean) => {
    if (isEditing) {
      if (typeof value === "object") {
        return (
          <textarea
            className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
            rows={6}
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
          />
        );
      } else {
        return (
          <input
            type="text"
            className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            value={editValue}
            onChange={(e) => setEditValue(e.target.value)}
          />
        );
      }
    }

    // Display mode
    if (typeof value === "boolean") {
      return (
        <span
          className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
            value
              ? "bg-green-100 text-green-700"
              : "bg-red-100 text-red-700"
          }`}
        >
          {value ? "Habilitado" : "Deshabilitado"}
        </span>
      );
    }

    if (typeof value === "object") {
      return (
        <pre className="bg-slate-50 p-2 rounded text-xs font-mono overflow-x-auto">
          {JSON.stringify(value, null, 2)}
        </pre>
      );
    }

    return <span className="text-sm text-slate-900">{String(value)}</span>;
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-slate-900">Configuración del Sistema</h1>
          <p className="text-slate-600 mt-2">Gestiona los parámetros de configuración de la plataforma</p>
        </div>
        <div className="flex items-center gap-3">
          <Button onClick={() => setShowNewForm(!showNewForm)} disabled={isLoading}>
            {showNewForm ? (
              <>
                <X className="w-4 h-4 mr-2" />
                Cancelar
              </>
            ) : (
              <>
                <Plus className="w-4 h-4 mr-2" />
                Nuevo Parámetro
              </>
            )}
          </Button>
          <Button variant="outline" onClick={() => refetch()} disabled={isLoading}>
            <RefreshCw className={`w-4 h-4 mr-2 ${isLoading ? "animate-spin" : ""}`} />
            Recargar
          </Button>
        </div>
      </div>

      {/* New Parameter Form */}
      {showNewForm && (
        <Card className="p-6 bg-blue-50 border-blue-200">
          <h3 className="text-lg font-semibold text-slate-900 mb-4">Crear Nuevo Parámetro</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Clave (Key) *
              </label>
              <input
                type="text"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="ej: max_upload_size"
                value={newKey}
                onChange={(e) => setNewKey(e.target.value)}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Valor *
              </label>
              <input
                type="text"
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="ej: 100 o true o {}"
                value={newValue}
                onChange={(e) => setNewValue(e.target.value)}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Tipo de Dato *
              </label>
              <select
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                value={newValueType}
                onChange={(e) => setNewValueType(e.target.value as "string" | "int" | "float" | "bool" | "json")}
              >
                <option value="string">string (texto)</option>
                <option value="int">int (entero)</option>
                <option value="float">float (decimal)</option>
                <option value="bool">bool (true/false)</option>
                <option value="json">json (objeto/array)</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">
                Categoría *
              </label>
              <select
                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                value={newCategory}
                onChange={(e) => setNewCategory(e.target.value)}
              >
                <option value="business">business</option>
                <option value="email">email</option>
                <option value="payment">payment</option>
                <option value="security">security</option>
                <option value="performance">performance</option>
                <option value="database">database</option>
              </select>
            </div>
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-slate-700 mb-2">
              Descripción (opcional)
            </label>
            <input
              type="text"
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="ej: Tamaño máximo de archivo en MB"
              value={newDescription}
              onChange={(e) => setNewDescription(e.target.value)}
            />
          </div>
          <div className="flex items-center gap-3">
            <Button
              onClick={handleCreateNew}
              disabled={updateMutation.isPending || !newKey.trim() || !newValue.trim()}
            >
              {updateMutation.isPending ? (
                <>
                  <div className="w-4 h-4 mr-2 inline-block">
                    <LoadingSpinner />
                  </div>
                  Creando...
                </>
              ) : (
                <>
                  <Save className="w-4 h-4 mr-2" />
                  Crear Parámetro
                </>
              )}
            </Button>
            <Button variant="outline" onClick={() => setShowNewForm(false)}>
              Cancelar
            </Button>
          </div>
        </Card>
      )}

      {/* Category Filter */}
      {data && data.Categories && data.Categories.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap">
          <Button
            variant={selectedCategory === undefined ? "default" : "outline"}
            size="sm"
            onClick={() => setSelectedCategory(undefined)}
          >
            Todas ({data.TotalSettings})
          </Button>
          {data.Categories.map((category) => (
            <Button
              key={category}
              variant={selectedCategory === category ? "default" : "outline"}
              size="sm"
              onClick={() => setSelectedCategory(category)}
            >
              {category}
            </Button>
          ))}
        </div>
      )}

      {/* Loading */}
      {isLoading && (
        <div className="flex items-center justify-center py-12">
          <LoadingSpinner />
        </div>
      )}

      {/* Settings List */}
      {!isLoading && data && (
        <>
          {data.Settings.length === 0 ? (
            <EmptyState
              icon={<Settings className="w-12 h-12 text-slate-400" />}
              title="No hay configuraciones"
              description="No se encontraron configuraciones para esta categoría"
            />
          ) : (
            <div className="space-y-4">
              {data.Settings.map((setting) => {
                const isEditing = editingKey === setting.key;

                return (
                  <Card key={setting.key} className="p-6">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-3">
                          {getCategoryIcon(setting.category)}
                          <div>
                            <h3 className="text-lg font-semibold text-slate-900">{setting.key}</h3>
                            <p className="text-xs text-slate-500">
                              Categoría: {setting.category} • Tipo: {setting.value_type || "string"} • Actualizado:{" "}
                              {new Date(setting.updated_at).toLocaleString("es-CR")}
                            </p>
                          </div>
                        </div>

                        <div className="mt-4">
                          <label className="block text-sm font-medium text-slate-700 mb-2">Valor</label>
                          {renderValue(setting.value, isEditing)}
                        </div>

                        {/* Mostrar descripción actual si existe y no estamos editando */}
                        {!isEditing && setting.description && (
                          <p className="text-sm text-slate-500 mt-2 italic">{setting.description}</p>
                        )}

                        {/* Campos adicionales de edición */}
                        {isEditing && (
                          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4 p-4 bg-slate-50 rounded-lg">
                            <div>
                              <label className="block text-sm font-medium text-slate-700 mb-2">
                                Tipo de Dato
                              </label>
                              <select
                                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                                value={editValueType}
                                onChange={(e) => setEditValueType(e.target.value as "string" | "int" | "float" | "bool" | "json")}
                              >
                                <option value="string">string (texto)</option>
                                <option value="int">int (entero)</option>
                                <option value="float">float (decimal)</option>
                                <option value="bool">bool (true/false)</option>
                                <option value="json">json (objeto/array)</option>
                              </select>
                            </div>
                            <div>
                              <label className="block text-sm font-medium text-slate-700 mb-2">
                                Categoría
                              </label>
                              <select
                                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                                value={editCategory}
                                onChange={(e) => setEditCategory(e.target.value)}
                              >
                                <option value="business">business</option>
                                <option value="email">email</option>
                                <option value="payment">payment</option>
                                <option value="security">security</option>
                                <option value="performance">performance</option>
                                <option value="database">database</option>
                              </select>
                            </div>
                            <div>
                              <label className="block text-sm font-medium text-slate-700 mb-2">
                                Descripción
                              </label>
                              <input
                                type="text"
                                className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                                placeholder="Descripción del parámetro"
                                value={editDescription}
                                onChange={(e) => setEditDescription(e.target.value)}
                              />
                            </div>
                          </div>
                        )}

                        {isEditing && (
                          <div className="flex items-center gap-3 mt-4">
                            <Button
                              size="sm"
                              onClick={() => handleSave(setting)}
                              disabled={updateMutation.isPending}
                            >
                              {updateMutation.isPending ? (
                                <>
                                  <div className="w-4 h-4 mr-2 inline-block">
                                    <LoadingSpinner />
                                  </div>
                                  Guardando...
                                </>
                              ) : (
                                <>
                                  <Save className="w-4 h-4 mr-2" />
                                  Guardar
                                </>
                              )}
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={handleCancel}
                              disabled={updateMutation.isPending}
                            >
                              Cancelar
                            </Button>
                            <p className="text-xs text-slate-500 ml-2">
                              Tip: Puedes usar true/false para booleanos, números directamente, o JSON para
                              objetos
                            </p>
                          </div>
                        )}
                      </div>

                      {!isEditing && (
                        <Button variant="outline" size="sm" onClick={() => handleEdit(setting)}>
                          Editar
                        </Button>
                      )}
                    </div>
                  </Card>
                );
              })}
            </div>
          )}
        </>
      )}

      {/* Info Card */}
      <Card className="p-6 bg-blue-50 border-blue-200">
        <div className="flex items-start gap-3">
          <Settings className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-semibold text-blue-900 mb-1">Importante</h3>
            <p className="text-sm text-blue-800">
              Los cambios en la configuración del sistema se aplican inmediatamente. Ten cuidado al
              modificar parámetros críticos como configuración de base de datos, autenticación o pagos.
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
