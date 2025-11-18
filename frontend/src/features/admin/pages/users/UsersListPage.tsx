import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Search, UserX, UserCheck, AlertCircle } from "lucide-react";
import { useAdminUsers } from "../../hooks/useAdminUsers";
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
import type { UserFilters, UserRole, UserStatus, KYCLevel } from "../../types";
import { format } from "date-fns";

export function UsersListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState<UserFilters>({});
  const [searchInput, setSearchInput] = useState("");

  const { data, isLoading, error } = useAdminUsers(filters, {
    page,
    limit: 20,
  });

  const handleSearch = () => {
    setFilters({ ...filters, search: searchInput });
    setPage(1);
  };

  const handleFilterChange = (key: keyof UserFilters, value: string) => {
    setFilters({ ...filters, [key]: value || undefined });
    setPage(1);
  };

  const handleRowClick = (userId: number) => {
    navigate(`/admin/users/${userId}`);
  };

  const getStatusBadge = (status: UserStatus) => {
    const variants: Record<UserStatus, { color: string; label: string }> = {
      active: { color: "bg-green-100 text-green-700", label: "Activo" },
      suspended: { color: "bg-amber-100 text-amber-700", label: "Suspendido" },
      banned: { color: "bg-red-100 text-red-700", label: "Bloqueado" },
      deleted: { color: "bg-slate-100 text-slate-700", label: "Eliminado" },
    };

    const variant = variants[status];
    return (
      <Badge className={variant.color}>
        {variant.label}
      </Badge>
    );
  };

  const getRoleBadge = (role: UserRole) => {
    const variants: Record<UserRole, { color: string; label: string }> = {
      user: { color: "bg-slate-100 text-slate-700", label: "Usuario" },
      admin: { color: "bg-blue-100 text-blue-700", label: "Admin" },
      super_admin: { color: "bg-blue-100 text-blue-700", label: "Super Admin" },
    };

    const variant = variants[role];
    return <Badge className={variant.color}>{variant.label}</Badge>;
  };

  const getKYCBadge = (level: KYCLevel) => {
    const variants: Record<KYCLevel, { color: string; label: string }> = {
      none: { color: "bg-slate-100 text-slate-700", label: "Sin KYC" },
      email_verified: { color: "bg-blue-100 text-blue-700", label: "Email" },
      phone_verified: { color: "bg-blue-100 text-blue-700", label: "Teléfono" },
      cedula_verified: { color: "bg-green-100 text-green-700", label: "Cédula" },
      full_kyc: { color: "bg-green-100 text-green-700", label: "Full KYC" },
    };

    const variant = variants[level];
    return <Badge className={variant.color}>{variant.label}</Badge>;
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-slate-900">Gestión de Usuarios</h1>
        <p className="text-slate-600 mt-2">
          Administra usuarios, KYC y suspensiones
        </p>
      </div>

      {/* Filters */}
      <Card className="p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* Search */}
          <div className="lg:col-span-2">
            <div className="flex gap-2">
              <Input
                placeholder="Buscar por nombre, email o cédula..."
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              />
              <Button onClick={handleSearch}>
                <Search className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Role filter */}
          <div>
            <select
              className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={filters.role || ""}
              onChange={(e) => handleFilterChange("role", e.target.value)}
            >
              <option value="">Todos los roles</option>
              <option value="user">Usuario</option>
              <option value="admin">Admin</option>
              <option value="super_admin">Super Admin</option>
            </select>
          </div>

          {/* Status filter */}
          <div>
            <select
              className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={filters.status || ""}
              onChange={(e) => handleFilterChange("status", e.target.value)}
            >
              <option value="">Todos los estados</option>
              <option value="active">Activo</option>
              <option value="suspended">Suspendido</option>
              <option value="banned">Bloqueado</option>
              <option value="deleted">Eliminado</option>
            </select>
          </div>
        </div>

        {/* Active filters indicator */}
        {Object.keys(filters).length > 0 && (
          <div className="mt-4 flex items-center gap-2">
            <span className="text-sm text-slate-600">Filtros activos:</span>
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setFilters({});
                setSearchInput("");
              }}
            >
              Limpiar filtros
            </Button>
          </div>
        )}
      </Card>

      {/* Results */}
      <Card>
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <LoadingSpinner />
          </div>
        ) : error ? (
          <div className="p-6">
            <EmptyState
              icon={<AlertCircle className="w-12 h-12 text-red-500" />}
              title="Error al cargar usuarios"
              description={error.message}
            />
          </div>
        ) : !data || !data.data || data.data.length === 0 ? (
          <div className="p-6">
            <EmptyState
              icon={<UserX className="w-12 h-12 text-slate-400" />}
              title="No se encontraron usuarios"
              description="Intenta ajustar los filtros de búsqueda"
            />
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Usuario</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Rol</TableHead>
                  <TableHead>Estado</TableHead>
                  <TableHead>KYC</TableHead>
                  <TableHead>Registro</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.data.map((user) => (
                  <TableRow
                    key={user.id}
                    onClick={() => handleRowClick(user.id)}
                    className="hover:bg-blue-50"
                  >
                    <TableCell>
                      <div>
                        <p className="font-medium text-slate-900">
                          {user.first_name} {user.last_name}
                        </p>
                        {user.cedula && (
                          <p className="text-xs text-slate-500">
                            Cédula: {user.cedula}
                          </p>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <span>{user.email}</span>
                        {user.email_verified && (
                          <UserCheck className="w-4 h-4 text-green-600" />
                        )}
                      </div>
                    </TableCell>
                    <TableCell>{getRoleBadge(user.role)}</TableCell>
                    <TableCell>{getStatusBadge(user.status)}</TableCell>
                    <TableCell>{getKYCBadge(user.kyc_level)}</TableCell>
                    <TableCell className="text-slate-600">
                      {format(new Date(user.created_at), "dd/MM/yyyy")}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {/* Pagination */}
            {data.pagination.total_pages > 1 && (
              <div className="px-6 py-4 border-t border-slate-200 flex items-center justify-between">
                <div className="text-sm text-slate-600">
                  Página {data.pagination.page} de {data.pagination.total_pages}
                  {" · "}
                  {data.pagination.total} usuarios totales
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage(page - 1)}
                    disabled={page === 1}
                  >
                    Anterior
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage(page + 1)}
                    disabled={page === data.pagination.total_pages}
                  >
                    Siguiente
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
      </Card>
    </div>
  );
}
