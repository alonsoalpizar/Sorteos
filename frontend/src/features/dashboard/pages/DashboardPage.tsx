import { useNavigate } from 'react-router-dom';
import { useUser } from '@/hooks/useAuth';
import { useUserMode } from '@/contexts/UserModeContext';
import { useEffect } from 'react';
import { LoadingSpinner } from '@/components/ui/LoadingSpinner';

export const DashboardPage = () => {
  const user = useUser();
  const { mode } = useUserMode();
  const navigate = useNavigate();

  // Redirect to appropriate dashboard based on mode
  useEffect(() => {
    if (user) {
      if (mode === 'participant') {
        navigate('/explore', { replace: true });
      } else {
        navigate('/organizer', { replace: true });
      }
    }
  }, [user, mode, navigate]);

  if (!user) {
    return <LoadingSpinner text="Cargando dashboard..." />;
  }

  // Show loading while redirecting
  return <LoadingSpinner text="Redirigiendo..." />;
};
