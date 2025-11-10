import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useVerifyEmail, useUser } from "@/hooks/useAuth";
import { getErrorMessage } from "@/lib/api";

import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Label } from "@/components/ui/Label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/Card";
import { Alert, AlertDescription } from "@/components/ui/Alert";

export const VerifyEmailPage = () => {
  const navigate = useNavigate();
  const user = useUser();
  const verifyEmailMutation = useVerifyEmail();
  const [code, setCode] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    if (!user) {
      navigate("/login");
    } else if (user.email_verified) {
      navigate("/dashboard");
    }
  }, [user, navigate]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!code || code.length !== 6) {
      setError("El código debe tener 6 dígitos");
      return;
    }

    if (!user?.id) {
      setError("Error: Usuario no encontrado");
      return;
    }

    verifyEmailMutation.mutate(
      {
        user_id: user.id,
        code: code,
      },
      {
        onSuccess: () => {
          navigate("/dashboard");
        },
        onError: (err) => {
          setError(getErrorMessage(err));
        },
      }
    );
  };

  const handleCodeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, "").slice(0, 6);
    setCode(value);
    setError("");
  };

  if (!user) {
    return null;
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 dark:bg-slate-900 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-3xl font-bold text-center">
            Verificar Email
          </CardTitle>
          <CardDescription className="text-center">
            Ingresa el código de 6 dígitos que enviamos a{" "}
            <span className="font-semibold text-foreground">{user.email}</span>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {verifyEmailMutation.isSuccess && (
              <Alert variant="success">
                <AlertDescription>
                  ¡Email verificado exitosamente! Redirigiendo...
                </AlertDescription>
              </Alert>
            )}

            {(error || verifyEmailMutation.isError) && (
              <Alert variant="destructive">
                <AlertDescription>
                  {error || getErrorMessage(verifyEmailMutation.error)}
                </AlertDescription>
              </Alert>
            )}

            <div className="space-y-2">
              <Label htmlFor="code" required>
                Código de Verificación
              </Label>
              <Input
                id="code"
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                placeholder="123456"
                value={code}
                onChange={handleCodeChange}
                maxLength={6}
                className="text-center text-2xl tracking-widest font-mono"
                autoComplete="off"
              />
              <p className="text-xs text-muted-foreground text-center">
                El código expira en 15 minutos
              </p>
            </div>

            <Button
              type="submit"
              className="w-full"
              loading={verifyEmailMutation.isPending}
              disabled={
                verifyEmailMutation.isPending ||
                code.length !== 6 ||
                verifyEmailMutation.isSuccess
              }
            >
              Verificar Email
            </Button>

            <div className="text-center">
              <Button
                type="button"
                variant="link"
                className="text-sm"
                onClick={() => {
                  // TODO: Implement resend verification code
                  alert("Funcionalidad de reenvío próximamente");
                }}
              >
                ¿No recibiste el código? Reenviar
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};
