import { createContext, useContext, useState, useEffect, ReactNode } from 'react';

export type UserMode = 'participant' | 'organizer';

interface UserModeContextType {
  mode: UserMode;
  setMode: (mode: UserMode) => void;
  toggleMode: () => void;
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
  }, [mode]);

  const setMode = (newMode: UserMode) => {
    setModeState(newMode);
  };

  const toggleMode = () => {
    setModeState(prev => prev === 'participant' ? 'organizer' : 'participant');
  };

  return (
    <UserModeContext.Provider value={{ mode, setMode, toggleMode }}>
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
