import { useUser, useLogout } from "@/hooks/useAuth";
import { Button } from "@/components/ui/Button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";

export const DashboardPage = () => {
  const user = useUser();
  const logoutMutation = useLogout();

  const handleLogout = () => {
    logoutMutation.mutate();
  };

  if (!user) {
    return null;
  }

  const kycLevelLabels: Record<typeof user.kyc_level, string> = {
    none: "Sin verificar",
    email_verified: "Email verificado",
    phone_verified: "Teléfono verificado",
    cedula_verified: "Cédula verificada",
    full_kyc: "KYC completo",
  };

  const roleLabels: Record<typeof user.role, string> = {
    user: "Usuario",
    admin: "Administrador",
    super_admin: "Super Administrador",
  };

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-900">
      <header className="border-b border-border bg-white dark:bg-slate-950">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-primary">Sorteos Platform</h1>
          <Button
            variant="outline"
            onClick={handleLogout}
            loading={logoutMutation.isPending}
          >
            Cerrar Sesión
          </Button>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto space-y-6">
          <div>
            <h2 className="text-3xl font-bold mb-2">
              Bienvenido, {user.first_name}!
            </h2>
            <p className="text-muted-foreground">
              Esta es tu área de usuario. Aquí podrás gestionar tus sorteos y
              participaciones.
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Información de la Cuenta</CardTitle>
              <CardDescription>
                Detalles de tu perfil y nivel de verificación
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-muted-foreground">Nombre completo</p>
                  <p className="font-medium">
                    {user.first_name} {user.last_name}
                  </p>
                </div>

                <div>
                  <p className="text-sm text-muted-foreground">Email</p>
                  <p className="font-medium">{user.email}</p>
                </div>

                {user.phone && (
                  <div>
                    <p className="text-sm text-muted-foreground">Teléfono</p>
                    <p className="font-medium">{user.phone}</p>
                  </div>
                )}

                {user.cedula && (
                  <div>
                    <p className="text-sm text-muted-foreground">Cédula</p>
                    <p className="font-medium">{user.cedula}</p>
                  </div>
                )}

                <div>
                  <p className="text-sm text-muted-foreground">Rol</p>
                  <Badge variant="secondary">{roleLabels[user.role]}</Badge>
                </div>

                <div>
                  <p className="text-sm text-muted-foreground">
                    Nivel de verificación
                  </p>
                  <Badge
                    variant={
                      user.kyc_level === "full_kyc"
                        ? "success"
                        : user.kyc_level === "none"
                        ? "destructive"
                        : "default"
                    }
                  >
                    {kycLevelLabels[user.kyc_level]}
                  </Badge>
                </div>

                <div>
                  <p className="text-sm text-muted-foreground">Estado</p>
                  <Badge
                    variant={
                      user.status === "active"
                        ? "success"
                        : user.status === "suspended"
                        ? "warning"
                        : "destructive"
                    }
                  >
                    {user.status === "active"
                      ? "Activo"
                      : user.status === "suspended"
                      ? "Suspendido"
                      : "Baneado"}
                  </Badge>
                </div>

                <div>
                  <p className="text-sm text-muted-foreground">UUID</p>
                  <p className="font-mono text-xs">{user.uuid}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Próximamente</CardTitle>
              <CardDescription>
                Funcionalidades en desarrollo
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ul className="list-disc list-inside space-y-2 text-muted-foreground">
                <li>Crear y gestionar sorteos</li>
                <li>Participar en sorteos activos</li>
                <li>Ver historial de participaciones</li>
                <li>Completar verificación KYC</li>
                <li>Gestionar métodos de pago</li>
                <li>Ver estadísticas y reportes</li>
              </ul>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
};
