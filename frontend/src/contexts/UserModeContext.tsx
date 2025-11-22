import { createContext, useContext, useState, useEffect, ReactNode } from 'react';

export type UserMode = 'participant' | 'organizer';

// Colores por modo - Participante: Azul, Organizador: Teal
export const modeColors = {
  participant: {
    primary: 'blue',
    bg: 'bg-blue-600',
    bgHover: 'hover:bg-blue-700',
    text: 'text-blue-600',
    textHover: 'hover:text-blue-600',
    border: 'border-blue-600',
    ring: 'ring-blue-500',
    // Variantes claras
    bgLight: 'bg-blue-50',
    textLight: 'text-blue-700',
    // Gradientes para headers
    gradient: 'bg-gradient-to-r from-blue-600 to-blue-700',
    gradientDark: 'dark:from-blue-700 dark:to-blue-800',
    // Texto sobre gradiente
    textMuted: 'text-blue-100',
  },
  organizer: {
    primary: 'teal',
    bg: 'bg-teal-600',
    bgHover: 'hover:bg-teal-700',
    text: 'text-teal-600',
    textHover: 'hover:text-teal-600',
    border: 'border-teal-600',
    ring: 'ring-teal-500',
    // Variantes claras
    bgLight: 'bg-teal-50',
    textLight: 'text-teal-700',
    // Gradientes para headers
    gradient: 'bg-gradient-to-r from-teal-600 to-teal-700',
    gradientDark: 'dark:from-teal-700 dark:to-teal-800',
    // Texto sobre gradiente
    textMuted: 'text-teal-100',
  },
} as const;

export type ModeColors = typeof modeColors[UserMode];

interface UserModeContextType {
  mode: UserMode;
  setMode: (mode: UserMode) => void;
  toggleMode: () => void;
  colors: ModeColors;
}

const UserModeContext = createContext<UserModeContextType | undefined>(undefined);

const USER_MODE_KEY = 'sorteos_user_mode';

export function UserModeProvider({ children }: { children: ReactNode }) {
  const [mode, setModeState] = useState<UserMode>(() => {
    // Load from localStorage on mount
    const saved = localStorage.getItem(USER_MODE_KEY);
    return (saved === 'organizer' ? 'organizer' : 'participant') as UserMode;
  });

  useEffect(() => {
    // Save to localStorage whenever mode changes
    localStorage.setItem(USER_MODE_KEY, mode);

    // Actualizar data-mode en el elemento HTML para CSS variables dinámicas
    document.documentElement.setAttribute('data-mode', mode);
  }, [mode]);

  // Establecer data-mode inicial al montar
  useEffect(() => {
    document.documentElement.setAttribute('data-mode', mode);
  }, []);

  const setMode = (newMode: UserMode) => {
    setModeState(newMode);
  };

  const toggleMode = () => {
    setModeState(prev => prev === 'participant' ? 'organizer' : 'participant');
  };

  // Colores actuales según el modo
  const colors = modeColors[mode];

  return (
    <UserModeContext.Provider value={{ mode, setMode, toggleMode, colors }}>
      {children}
    </UserModeContext.Provider>
  );
}

export function useUserMode() {
  const context = useContext(UserModeContext);
  if (context === undefined) {
    throw new Error('useUserMode must be used within a UserModeProvider');
  }
  return context;
}
