import { Check, X } from "lucide-react";
import { cn } from "@/lib/utils";

interface PasswordStrengthProps {
  password: string;
}

export const PasswordStrength = ({ password }: PasswordStrengthProps) => {
  const requirements = [
    {
      test: password.length >= 12,
      label: "Mínimo 12 caracteres",
    },
    {
      test: /[A-Z]/.test(password),
      label: "Una letra mayúscula",
    },
    {
      test: /[a-z]/.test(password),
      label: "Una letra minúscula",
    },
    {
      test: /[0-9]/.test(password),
      label: "Un número",
    },
    {
      test: /[!@#$%^&*]/.test(password),
      label: "Un símbolo especial (!@#$%^&*)",
    },
  ];

  return (
    <div className="space-y-2 mt-2">
      <p className="text-xs font-medium text-muted-foreground">
        La contraseña debe contener:
      </p>
      <div className="space-y-1">
        {requirements.map((req, index) => (
          <div key={index} className="flex items-center gap-2">
            {req.test ? (
              <Check className="w-3.5 h-3.5 text-green-600" />
            ) : (
              <X className="w-3.5 h-3.5 text-gray-300" />
            )}
            <span
              className={cn(
                "text-xs",
                req.test ? "text-green-600 font-medium" : "text-muted-foreground"
              )}
            >
              {req.label}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
};
