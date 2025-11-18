import { Link, useNavigate } from 'react-router-dom';
import { UserMenu } from './UserMenu';
import { useAuthStore } from '@/store/authStore';
import { useUserMode } from '@/contexts/UserModeContext';
import { useState } from 'react';

export function Navbar() {
  const { user } = useAuthStore();
  const { mode } = useUserMode();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/raffles?search=${encodeURIComponent(searchQuery.trim())}`);
      setSearchQuery('');
    }
  };

  return (
    <header className="sticky top-0 z-50 bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-700 shadow-sm">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2 group">
            <div className="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center transition-transform group-hover:scale-105">
              <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
              </svg>
            </div>
            <span className="text-xl font-bold text-slate-900 dark:text-white hidden sm:block">
              Sorteos.club
            </span>
          </Link>

          {/* Search bar (only for authenticated users) */}
          {user && (
            <form onSubmit={handleSearch} className="flex-1 max-w-md mx-4 hidden md:block">
              <div className="relative">
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Buscar sorteos..."
                  className="w-full pl-10 pr-4 py-2 border border-slate-300 dark:border-slate-600 rounded-lg bg-slate-50 dark:bg-slate-800 text-slate-900 dark:text-white placeholder-slate-500 dark:placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
                />
                <svg
                  className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </form>
          )}

          {/* Navigation links and user menu */}
          <div className="flex items-center gap-4">
            {/* Main navigation for authenticated users */}
            {user && (
              <nav className="hidden lg:flex items-center gap-1">
                {mode === 'participant' ? (
                  <>
                    <Link
                      to="/explore"
                      className="px-3 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                      Explorar
                    </Link>
                    <Link
                      to="/my-tickets"
                      className="px-3 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                      Mis Números
                    </Link>
                    <Link
                      to="/wallet"
                      className="px-3 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                      💰 Billetera
                    </Link>
                  </>
                ) : (
                  <>
                    <Link
                      to="/organizer"
                      className="px-3 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                      Panel
                    </Link>
                    <Link
                      to="/organizer/raffles"
                      className="px-3 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                      Sorteos
                    </Link>
                    <Link
                      to="/wallet"
                      className="px-3 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
                    >
                      💰 Billetera
                    </Link>
                    <Link
                      to="/organizer/raffles/new"
                      className="px-3 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
                    >
                      + Crear
                    </Link>
                  </>
                )}
              </nav>
            )}

            {/* User menu */}
            <UserMenu />
          </div>
        </div>
      </div>

      {/* Mobile navigation tabs (only for authenticated users) */}
      {user && (
        <div className="border-t border-slate-200 dark:border-slate-700 lg:hidden">
          <nav className="container mx-auto px-4 flex items-center gap-1 overflow-x-auto">
            {mode === 'participant' ? (
              <>
                <Link
                  to="/explore"
                  className="px-4 py-3 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors whitespace-nowrap"
                >
                  Explorar
                </Link>
                <Link
                  to="/my-tickets"
                  className="px-4 py-3 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors whitespace-nowrap"
                >
                  Mis Números
                </Link>
              </>
            ) : (
              <>
                <Link
                  to="/organizer"
                  className="px-4 py-3 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors whitespace-nowrap"
                >
                  Panel
                </Link>
                <Link
                  to="/organizer/raffles"
                  className="px-4 py-3 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors whitespace-nowrap"
                >
                  Sorteos
                </Link>
                <Link
                  to="/organizer/raffles/new"
                  className="px-4 py-3 text-sm font-medium text-blue-600 dark:text-blue-400 transition-colors whitespace-nowrap"
                >
                  + Crear
                </Link>
              </>
            )}
          </nav>
        </div>
      )}
    </header>
  );
}
