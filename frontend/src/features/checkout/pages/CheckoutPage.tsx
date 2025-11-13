import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { LoadingSpinner } from '../../../components/ui/LoadingSpinner';

/**
 * CheckoutPage - YA NO SE USA
 *
 * El flujo de pago ahora es directo desde la grilla de números.
 * Esta página solo redirige a /raffles.
 */
export function CheckoutPage() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate('/raffles');
  }, [navigate]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <LoadingSpinner text="Redirigiendo..." />
    </div>
  );
}
