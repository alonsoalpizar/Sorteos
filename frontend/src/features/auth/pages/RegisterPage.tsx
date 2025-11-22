import { useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useRegister, useIsAuthenticated } from "@/hooks/useAuth";
import { getErrorMessage } from "@/lib/api";

import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Label } from "@/components/ui/Label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/Card";
import { Alert, AlertDescription } from "@/components/ui/Alert";
import { PasswordStrength } from "@/components/ui/PasswordStrength";
import { GoogleAuthButton } from "../components/GoogleAuthButton";

const registerSchema = z
  .object({
    email: z.string().email("Email inválido"),
    password: z
      .string()
      .min(12, "La contraseña debe tener al menos 12 caracteres")
      .regex(/[A-Z]/, "Debe contener al menos una mayúscula")
      .regex(/[a-z]/, "Debe contener al menos una minúscula")
      .regex(/[0-9]/, "Debe contener al menos un número")
      .regex(/[^A-Za-z0-9]/, "Debe contener al menos un símbolo"),
    confirmPassword: z.string(),
    first_name: z.string().min(2, "El nombre debe tener al menos 2 caracteres"),
    last_name: z.string().min(2, "El apellido debe tener al menos 2 caracteres"),
    phone: z.string().optional(),
    accept_terms: z.boolean().refine((val) => val === true, {
      message: "Debes aceptar los términos y condiciones",
    }),
    accept_privacy: z.boolean().refine((val) => val === true, {
      message: "Debes aceptar la política de privacidad",
    }),
    accept_marketing: z.boolean().optional(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Las contraseñas no coinciden",
    path: ["confirmPassword"],
  });

type RegisterFormData = z.infer<typeof registerSchema>;

export const RegisterPage = () => {
  const navigate = useNavigate();
  const isAuthenticated = useIsAuthenticated();
  const registerMutation = useRegister();

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    mode: "onBlur",
    defaultValues: {
      accept_terms: false,
      accept_privacy: false,
      accept_marketing: false,
    },
  });

  const password = watch("password", "");

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate("/dashboard");
    }
  }, [isAuthenticated, navigate]);

  const onSubmit = (data: RegisterFormData) => {
    const { confirmPassword, accept_terms, accept_privacy, accept_marketing, ...rest } = data;

    // Transformar nombres de campos para coincidir con backend
    const registerData = {
      ...rest,
      accepted_terms: accept_terms,
      accepted_privacy: accept_privacy,
      accepted_marketing: accept_marketing,
    };

    registerMutation.mutate(registerData, {
      onSuccess: () => {
        navigate("/verify-email");
      },
    });
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 dark:bg-slate-900 p-4">
      <Card className="w-full max-w-2xl my-8">
        <CardHeader className="space-y-1">
          <CardTitle className="text-3xl font-bold text-center">
            Crear Cuenta
          </CardTitle>
          <CardDescription className="text-center">
            Completa el formulario para registrarte en Sorteos Platform
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Google OAuth Button - Registro rápido */}
          <div className="mb-6">
            <GoogleAuthButton mode="register" />
            <p className="text-xs text-center text-slate-500 mt-2">
              Al registrarte con Google, aceptas automáticamente los términos y condiciones
            </p>
          </div>

          {/* Divider */}
          <div className="relative mb-6">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t border-slate-200" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-white px-2 text-slate-500">O regístrate con email</span>
            </div>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {registerMutation.isError && (
              <Alert variant="destructive">
                <AlertDescription>
                  {getErrorMessage(registerMutation.error)}
                </AlertDescription>
              </Alert>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="first_name" required>
                  Nombre
                </Label>
                <Input
                  id="first_name"
                  placeholder="Juan"
                  error={errors.first_name?.message}
                  {...register("first_name")}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="last_name" required>
                  Apellido
                </Label>
                <Input
                  id="last_name"
                  placeholder="Pérez"
                  error={errors.last_name?.message}
                  {...register("last_name")}
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="email" required>
                Email
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="tu@email.com"
                error={errors.email?.message}
                {...register("email")}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="phone">Teléfono (opcional)</Label>
              <Input
                id="phone"
                type="tel"
                placeholder="+573001234567"
                error={errors.phone?.message}
                {...register("phone")}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password" required>
                Contraseña
              </Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••••••"
                error={errors.password?.message}
                {...register("password")}
              />
              <PasswordStrength password={password} />
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirmPassword" required>
                Confirmar Contraseña
              </Label>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="••••••••••••"
                error={errors.confirmPassword?.message}
                {...register("confirmPassword")}
              />
            </div>

            <div className="space-y-3 pt-2">
              <div className="flex items-start space-x-2">
                <input
                  type="checkbox"
                  id="accept_terms"
                  className="mt-1 h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                  {...register("accept_terms")}
                />
                <div className="flex-1">
                  <label
                    htmlFor="accept_terms"
                    className="text-sm text-muted-foreground"
                  >
                    Acepto los{" "}
                    <Link to="/terms" className="text-primary hover:underline">
                      términos y condiciones
                    </Link>
                  </label>
                  {errors.accept_terms && (
                    <p className="text-xs text-destructive mt-1">
                      {errors.accept_terms.message}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-start space-x-2">
                <input
                  type="checkbox"
                  id="accept_privacy"
                  className="mt-1 h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                  {...register("accept_privacy")}
                />
                <div className="flex-1">
                  <label
                    htmlFor="accept_privacy"
                    className="text-sm text-muted-foreground"
                  >
                    Acepto la{" "}
                    <Link to="/privacy" className="text-primary hover:underline">
                      política de privacidad
                    </Link>
                  </label>
                  {errors.accept_privacy && (
                    <p className="text-xs text-destructive mt-1">
                      {errors.accept_privacy.message}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-start space-x-2">
                <input
                  type="checkbox"
                  id="accept_marketing"
                  className="mt-1 h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                  {...register("accept_marketing")}
                />
                <label
                  htmlFor="accept_marketing"
                  className="text-sm text-muted-foreground"
                >
                  Acepto recibir comunicaciones de marketing (opcional)
                </label>
              </div>
            </div>

            <Button
              type="submit"
              className="w-full"
              loading={registerMutation.isPending}
              disabled={registerMutation.isPending}
            >
              Crear Cuenta
            </Button>

            <div className="text-center text-sm">
              <span className="text-muted-foreground">
                ¿Ya tienes una cuenta?{" "}
              </span>
              <Link
                to="/login"
                className="text-primary hover:underline font-medium"
              >
                Inicia sesión
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};
