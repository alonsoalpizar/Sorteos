import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import * as authApi from "@/api/auth";
import { useAuthStore } from "@/store/authStore";
import { getErrorMessage } from "@/lib/api";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Label } from "@/components/ui/Label";
import { X } from "lucide-react";

interface GoogleLinkModalProps {
  email: string;
  idToken: string;
  onSuccess: () => void;
  onClose: () => void;
}

export const GoogleLinkModal = ({
  email,
  idToken,
  onSuccess,
  onClose,
}: GoogleLinkModalProps) => {
  const setAuth = useAuthStore((state) => state.setAuth);
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  const linkMutation = useMutation({
    mutationFn: () =>
      authApi.googleLink({
        id_token: idToken,
        password: password,
      }),
    onSuccess: (data) => {
      setAuth(data.user, data.access_token, data.refresh_token);
      onSuccess();
    },
    onError: (err) => {
      setError(getErrorMessage(err));
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!password.trim()) {
      setError("La contraseña es requerida");
      return;
    }

    linkMutation.mutate();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative bg-white rounded-lg shadow-xl w-full max-w-md mx-4 p-6 z-10">
        {/* Close button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-slate-400 hover:text-slate-600 transition-colors"
        >
          <X className="w-5 h-5" />
        </button>

        {/* Header */}
        <div className="text-center mb-6">
          <div className="mx-auto w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center mb-4">
            <svg className="w-6 h-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            </svg>
          </div>
          <h2 className="text-xl font-semibold text-slate-900">
            Vincular cuenta de Google
          </h2>
          <p className="text-sm text-slate-600 mt-2">
            Ya existe una cuenta registrada con el email:
          </p>
          <p className="text-sm font-medium text-slate-900 mt-1">{email}</p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          <p className="text-sm text-slate-600">
            Ingresa tu contraseña actual para vincular tu cuenta de Google y poder iniciar sesión
            con ambos métodos.
          </p>

          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-3">
              <p className="text-sm text-red-700">{error}</p>
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="link-password" required>
              Contraseña de tu cuenta
            </Label>
            <Input
              id="link-password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoFocus
            />
          </div>

          <div className="flex gap-3">
            <Button
              type="button"
              variant="outline"
              className="flex-1"
              onClick={onClose}
              disabled={linkMutation.isPending}
            >
              Cancelar
            </Button>
            <Button
              type="submit"
              className="flex-1"
              loading={linkMutation.isPending}
              disabled={linkMutation.isPending || !password.trim()}
            >
              Vincular cuenta
            </Button>
          </div>
        </form>

        {/* Footer note */}
        <p className="text-xs text-slate-500 text-center mt-4">
          Una vez vinculada, podrás iniciar sesión con tu contraseña o con Google.
        </p>
      </div>
    </div>
  );
};
