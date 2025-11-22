import { useState } from "react";
import { useGoogleLogin } from "@react-oauth/google";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import * as authApi from "@/api/auth";
import { useAuthStore } from "@/store/authStore";
import { getErrorMessage } from "@/lib/api";
import { Button } from "@/components/ui/Button";
import { GoogleLinkModal } from "./GoogleLinkModal";
import type { GoogleAuthRequiresLinkingResponse } from "@/types/auth";

// Google SVG icon
const GoogleIcon = () => (
  <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
    <path
      fill="#4285F4"
      d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
    />
    <path
      fill="#34A853"
      d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
    />
    <path
      fill="#FBBC05"
      d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
    />
    <path
      fill="#EA4335"
      d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
    />
  </svg>
);

interface GoogleAuthButtonProps {
  mode: "login" | "register";
  className?: string;
}

export const GoogleAuthButton = ({ mode, className = "" }: GoogleAuthButtonProps) => {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);
  const [linkingData, setLinkingData] = useState<{ email: string; idToken: string } | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Mutation para autenticación con Google
  const googleAuthMutation = useMutation({
    mutationFn: authApi.googleAuth,
    onSuccess: (data) => {
      // Si requiere vinculación
      if ("requires_linking" in data && data.requires_linking) {
        const linkData = data as GoogleAuthRequiresLinkingResponse;
        setLinkingData({
          email: linkData.email,
          idToken: localStorage.getItem("pending_google_token") || "",
        });
        return;
      }

      // Login exitoso - necesitamos cast porque ya verificamos que no es requires_linking
      const authData = data as import("@/types/auth").GoogleAuthResponse;
      setAuth(authData.user, authData.access_token, authData.refresh_token);
      navigate("/dashboard");
    },
    onError: (err) => {
      setError(getErrorMessage(err));
    },
  });

  // Hook de Google Login usando access_token (implicit flow)
  const googleLogin = useGoogleLogin({
    onSuccess: async (tokenResponse) => {
      setError(null);
      // Usar el access_token para obtener info del usuario de Google
      // El backend verificará el token
      try {
        // Guardamos el token temporalmente por si necesitamos vinculación
        localStorage.setItem("pending_google_token", tokenResponse.access_token);
        googleAuthMutation.mutate({ id_token: tokenResponse.access_token });
      } catch {
        setError("Error al procesar la autenticación con Google");
      }
    },
    onError: (error) => {
      console.error("Google Login Error:", error);
      setError("Error al conectar con Google. Por favor intenta de nuevo.");
    },
    flow: "implicit",
  });

  const handleGoogleClick = () => {
    setError(null);
    googleLogin();
  };

  const handleLinkSuccess = () => {
    setLinkingData(null);
    localStorage.removeItem("pending_google_token");
    navigate("/dashboard");
  };

  const handleLinkClose = () => {
    setLinkingData(null);
    localStorage.removeItem("pending_google_token");
  };

  const buttonText = mode === "login" ? "Continuar con Google" : "Registrarse con Google";

  return (
    <>
      <div className={className}>
        {error && (
          <p className="text-sm text-red-600 mb-2 text-center">{error}</p>
        )}
        <Button
          type="button"
          variant="outline"
          className="w-full"
          onClick={handleGoogleClick}
          disabled={googleAuthMutation.isPending}
        >
          {googleAuthMutation.isPending ? (
            <span className="flex items-center">
              <svg className="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
              </svg>
              Conectando...
            </span>
          ) : (
            <>
              <GoogleIcon />
              {buttonText}
            </>
          )}
        </Button>
      </div>

      {/* Modal de vinculación de cuenta */}
      {linkingData && (
        <GoogleLinkModal
          email={linkingData.email}
          idToken={linkingData.idToken}
          onSuccess={handleLinkSuccess}
          onClose={handleLinkClose}
        />
      )}
    </>
  );
};
