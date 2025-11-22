import { ReactNode } from 'react';
import { Link } from 'react-router-dom';
import { Navbar } from './Navbar';
import { useUserMode } from '@/contexts/UserModeContext';

interface MainLayoutProps {
  children: ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const { mode } = useUserMode();
  const isOrganizer = mode === 'organizer';

  // Colores dinámicos según modo
  const logoColor = isOrganizer ? 'bg-teal-600' : 'bg-blue-600';
  const hoverColor = isOrganizer
    ? 'hover:text-teal-600 dark:hover:text-teal-400'
    : 'hover:text-blue-600 dark:hover:text-blue-400';

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-900">
      <Navbar />

      <main className="container mx-auto px-4 py-8">
        {children}
      </main>

      {/* Footer */}
      <footer className="mt-auto border-t border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800">
        <div className="container mx-auto px-4 py-8">
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center gap-2 mb-3">
                <div className={`w-8 h-8 ${logoColor} rounded-lg flex items-center justify-center`}>
                  <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
                  </svg>
                </div>
                <span className="font-semibold text-slate-900 dark:text-white">Sorteos.club</span>
              </div>
              <p className="text-sm text-slate-600 dark:text-slate-400">
                La plataforma más confiable para sorteos en línea
              </p>
            </div>

            <div>
              <h4 className="font-semibold text-slate-900 dark:text-white mb-3">Plataforma</h4>
              <ul className="space-y-2 text-sm text-slate-600 dark:text-slate-400">
                <li>
                  <Link to="/explore" className={`${hoverColor} transition-colors`}>
                    Ver Sorteos
                  </Link>
                </li>
                <li>
                  <Link to="/organizer/raffles/new" className={`${hoverColor} transition-colors`}>
                    Crear Sorteo
                  </Link>
                </li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-slate-900 dark:text-white mb-3">Legal</h4>
              <ul className="space-y-2 text-sm text-slate-600 dark:text-slate-400">
                <li>
                  <Link to="#" className={`${hoverColor} transition-colors`}>
                    Términos y Condiciones
                  </Link>
                </li>
                <li>
                  <Link to="#" className={`${hoverColor} transition-colors`}>
                    Política de Privacidad
                  </Link>
                </li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold text-slate-900 dark:text-white mb-3">Soporte</h4>
              <ul className="space-y-2 text-sm text-slate-600 dark:text-slate-400">
                <li>
                  <Link to="#" className={`${hoverColor} transition-colors`}>
                    Centro de Ayuda
                  </Link>
                </li>
                <li>
                  <Link to="#" className={`${hoverColor} transition-colors`}>
                    Contacto
                  </Link>
                </li>
              </ul>
            </div>
          </div>

          <div className="border-t border-slate-200 dark:border-slate-700 mt-8 pt-6 text-center">
            <p className="text-sm text-slate-600 dark:text-slate-400">
              &copy; {new Date().getFullYear()} Sorteos.club. Todos los derechos reservados.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
