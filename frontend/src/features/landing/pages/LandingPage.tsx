import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { useIsAuthenticated } from '@/hooks/useAuth';
import { Shield, Eye, Zap, Gift, Ticket, Trophy, Ban, CheckCircle, UserCheck, FileCheck, Clock, ShieldX } from 'lucide-react';

export function LandingPage() {
  const isAuthenticated = useIsAuthenticated();

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">
      {/* Header/Nav */}
      <header className="border-b border-slate-200 dark:border-slate-700 bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
                <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
                </svg>
              </div>
              <span className="text-2xl font-bold text-slate-900 dark:text-white">
                Sorteos.club
              </span>
            </div>

            <nav className="hidden md:flex items-center gap-6">
              {isAuthenticated ? (
                <Link to="/dashboard">
                  <Button>
                    Ir al Panel
                  </Button>
                </Link>
              ) : (
                <>
                  <Link to="/login">
                    <Button variant="outline">
                      Iniciar Sesión
                    </Button>
                  </Link>
                  <Link to="/register">
                    <Button>
                      Registrarse
                    </Button>
                  </Link>
                </>
              )}
            </nav>

            {/* Mobile menu button */}
            <div className="md:hidden flex gap-2">
              {isAuthenticated ? (
                <Link to="/dashboard">
                  <Button size="sm">
                    Mi Panel
                  </Button>
                </Link>
              ) : (
                <Link to="/login">
                  <Button variant="outline" size="sm">
                    Ingresar
                  </Button>
                </Link>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        {/* Background decorations */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-blue-400/20 rounded-full blur-3xl" />
          <div className="absolute top-20 -left-20 w-60 h-60 bg-indigo-400/20 rounded-full blur-3xl" />
          <div className="absolute bottom-0 right-1/4 w-40 h-40 bg-teal-400/20 rounded-full blur-3xl" />
        </div>

        <div className="container mx-auto px-4 py-16 md:py-24 relative">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            {/* Left: Text content */}
            <div className="text-center lg:text-left">
              {/* Badge */}
              <div className="inline-flex items-center gap-2 bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300 px-4 py-2 rounded-full text-sm font-medium mb-6">
                <span className="flex h-2 w-2 rounded-full bg-green-500 animate-pulse" />
                Nuevo en Costa Rica
              </div>

              <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-slate-900 dark:text-white mb-6 leading-tight">
                Sorteos{' '}
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600">
                  transparentes
                </span>{' '}
                y seguros
              </h1>

              <p className="text-lg md:text-xl text-slate-600 dark:text-slate-300 mb-8 max-w-xl mx-auto lg:mx-0">
                Gestiona y realiza sorteos verificables.
                La forma más confiable de organizar y participar en sorteos en línea.
              </p>

              <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
                {isAuthenticated ? (
                  <Link to="/explore">
                    <Button size="lg" className="w-full sm:w-auto shadow-lg shadow-blue-500/25">
                      Ver Sorteos Disponibles
                      <Ticket className="w-5 h-5 ml-2" />
                    </Button>
                  </Link>
                ) : (
                  <>
                    <Link to="/register">
                      <Button size="lg" className="w-full sm:w-auto shadow-lg shadow-blue-500/25">
                        Crear Cuenta Gratis
                        <Gift className="w-5 h-5 ml-2" />
                      </Button>
                    </Link>
                    <Link to="/login">
                      <Button size="lg" variant="outline" className="w-full sm:w-auto">
                        Ya tengo cuenta
                      </Button>
                    </Link>
                  </>
                )}
              </div>

              {/* Trust indicators */}
              <div className="flex flex-wrap items-center justify-center lg:justify-start gap-6 mt-10 text-sm text-slate-500 dark:text-slate-400">
                <div className="flex items-center gap-2">
                  <Shield className="w-4 h-4 text-green-500" />
                  <span>Datos protegidos</span>
                </div>
                <div className="flex items-center gap-2">
                  <Eye className="w-4 h-4 text-blue-500" />
                  <span>100% auditable</span>
                </div>
                <div className="flex items-center gap-2">
                  <Zap className="w-4 h-4 text-amber-500" />
                  <span>Resultados instantáneos</span>
                </div>
              </div>
            </div>

            {/* Right: Visual illustration */}
            <div className="relative hidden lg:block">
              {/* Main card - Raffle preview mockup */}
              <div className="relative z-10 bg-white dark:bg-slate-800 rounded-2xl shadow-2xl p-6 transform rotate-2 hover:rotate-0 transition-transform duration-500">
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center">
                    <Trophy className="w-6 h-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-slate-900 dark:text-white">iPhone 15 Pro Max</h3>
                    <p className="text-sm text-slate-500">Sorteo activo</p>
                  </div>
                </div>

                {/* Ticket grid preview */}
                <div className="grid grid-cols-5 gap-2 mb-4">
                  {[...Array(15)].map((_, i) => {
                    const isSold = [0, 2, 4, 7, 9, 11, 13].includes(i);
                    const isWinner = i === 7;
                    return (
                      <div
                        key={i}
                        className={`
                          aspect-square rounded-lg flex items-center justify-center text-sm font-medium
                          ${isWinner
                            ? 'bg-gradient-to-br from-amber-400 to-orange-500 text-white animate-pulse'
                            : isSold
                              ? 'bg-blue-100 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400'
                              : 'bg-slate-100 dark:bg-slate-700 text-slate-400'
                          }
                        `}
                      >
                        {String(i).padStart(2, '0')}
                      </div>
                    );
                  })}
                </div>

                <div className="flex items-center justify-between text-sm">
                  <span className="text-slate-500">7 de 100 vendidos</span>
                  <span className="font-semibold text-green-600">¢5,000 / boleto</span>
                </div>
              </div>

              {/* Floating elements */}
              <div className="absolute -top-6 right-8 bg-green-500 text-white px-4 py-2 rounded-full text-sm font-medium shadow-lg animate-bounce z-20">
                +₡50,000 recaudado
              </div>

              <div className="absolute bottom-4 -left-4 bg-white dark:bg-slate-800 rounded-xl shadow-lg p-3 transform rotate-3 z-20">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 bg-amber-100 dark:bg-amber-900/50 rounded-full flex items-center justify-center">
                    <Trophy className="w-4 h-4 text-amber-600" />
                  </div>
                  <div>
                    <p className="text-xs text-slate-500">Ganador</p>
                    <p className="font-semibold text-slate-900 dark:text-white text-sm">#08 - María G.</p>
                  </div>
                </div>
              </div>

              {/* Background decoration card - hidden to avoid visual clutter */}
            </div>
          </div>

        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-white dark:bg-slate-800">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 dark:text-white mb-4">
              ¿Por qué elegir Sorteos.club?
            </h2>
            <p className="text-slate-600 dark:text-slate-300 max-w-2xl mx-auto">
              Una plataforma diseñada pensando en la transparencia, seguridad y facilidad de uso
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {/* Feature 1 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-green-400 to-emerald-500 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-green-500/25">
                <Shield className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                100% Seguro
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Tus datos y transacciones están protegidos con los más altos estándares de seguridad.
              </p>
            </div>

            {/* Feature 2 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-indigo-500 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-blue-500/25">
                <Eye className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Transparente
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Todos los sorteos son auditables. Ve en tiempo real el estado de cada número vendido.
              </p>
            </div>

            {/* Feature 3 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-amber-400 to-orange-500 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-amber-500/25">
                <Zap className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Rápido y Fácil
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Crea un sorteo en minutos. Interfaz intuitiva y proceso simplificado.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How it Works */}
      <section className="py-20 bg-slate-50 dark:bg-slate-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 dark:text-white mb-4">
              ¿Cómo funciona?
            </h2>
            <p className="text-slate-600 dark:text-slate-300 max-w-2xl mx-auto">
              En solo 3 pasos puedes crear tu sorteo o participar en uno existente
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-12 max-w-5xl mx-auto">
            {/* Step 1 */}
            <div className="text-center">
              <div className="w-16 h-16 bg-blue-600 text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4">
                1
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Regístrate
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Crea tu cuenta de forma gratuita y verifica tu correo electrónico
              </p>
            </div>

            {/* Step 2 */}
            <div className="text-center">
              <div className="w-16 h-16 bg-blue-600 text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4">
                2
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Crea o Participa
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Crea tu propio sorteo o elige números en sorteos activos
              </p>
            </div>

            {/* Step 3 */}
            <div className="text-center">
              <div className="w-16 h-16 bg-blue-600 text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4">
                3
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Gana
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Espera el sorteo y verifica los resultados de forma transparente
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Platform Rules Section */}
      <section className="py-20 bg-white dark:bg-slate-800">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 dark:text-white mb-4">
              Normas que garantizan tu seguridad
            </h2>
            <p className="text-slate-600 dark:text-slate-300 max-w-2xl mx-auto">
              Reglas claras para sorteos justos, transparentes y legales
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-6xl mx-auto">
            {/* Rule 1 - No money */}
            <div className="flex items-start gap-4 p-6 bg-red-50 dark:bg-red-900/20 rounded-xl border border-red-200 dark:border-red-800">
              <div className="w-10 h-10 bg-red-100 dark:bg-red-900/50 rounded-lg flex items-center justify-center flex-shrink-0">
                <Ban className="w-5 h-5 text-red-600 dark:text-red-400" />
              </div>
              <div>
                <h3 className="font-semibold text-slate-900 dark:text-white mb-1">Solo productos, no dinero</h3>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Solo se permite sortear artículos y productos. No se pueden rifar premios en efectivo.
                </p>
              </div>
            </div>

            {/* Rule 2 - Prohibited items */}
            <div className="flex items-start gap-4 p-6 bg-rose-50 dark:bg-rose-900/20 rounded-xl border border-rose-200 dark:border-rose-800">
              <div className="w-10 h-10 bg-rose-100 dark:bg-rose-900/50 rounded-lg flex items-center justify-center flex-shrink-0">
                <ShieldX className="w-5 h-5 text-rose-600 dark:text-rose-400" />
              </div>
              <div>
                <h3 className="font-semibold text-slate-900 dark:text-white mb-1">Productos prohibidos</h3>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  No se permiten armas, sustancias ilícitas, artículos robados ni cualquier producto que infrinja la ley.
                </p>
              </div>
            </div>

            {/* Rule 3 - Verified source */}
            <div className="flex items-start gap-4 p-6 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200 dark:border-blue-800">
              <div className="w-10 h-10 bg-blue-100 dark:bg-blue-900/50 rounded-lg flex items-center justify-center flex-shrink-0">
                <CheckCircle className="w-5 h-5 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <h3 className="font-semibold text-slate-900 dark:text-white mb-1">Fuente verificable</h3>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Los resultados se basan en la Lotería Nacional de Costa Rica, garantizando imparcialidad.
                </p>
              </div>
            </div>

            {/* Rule 4 - Verified users */}
            <div className="flex items-start gap-4 p-6 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-200 dark:border-green-800">
              <div className="w-10 h-10 bg-green-100 dark:bg-green-900/50 rounded-lg flex items-center justify-center flex-shrink-0">
                <UserCheck className="w-5 h-5 text-green-600 dark:text-green-400" />
              </div>
              <div>
                <h3 className="font-semibold text-slate-900 dark:text-white mb-1">Usuarios verificados</h3>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Todos los participantes y organizadores deben verificar su identidad.
                </p>
              </div>
            </div>

            {/* Rule 5 - Moderation */}
            <div className="flex items-start gap-4 p-6 bg-amber-50 dark:bg-amber-900/20 rounded-xl border border-amber-200 dark:border-amber-800">
              <div className="w-10 h-10 bg-amber-100 dark:bg-amber-900/50 rounded-lg flex items-center justify-center flex-shrink-0">
                <FileCheck className="w-5 h-5 text-amber-600 dark:text-amber-400" />
              </div>
              <div>
                <h3 className="font-semibold text-slate-900 dark:text-white mb-1">Sorteos moderados</h3>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Cada sorteo es revisado antes de publicarse. Contenido inapropiado será suspendido.
                </p>
              </div>
            </div>

            {/* Rule 6 - Full traceability */}
            <div className="flex items-start gap-4 p-6 bg-slate-100 dark:bg-slate-700/50 rounded-xl border border-slate-200 dark:border-slate-600">
              <div className="w-10 h-10 bg-slate-200 dark:bg-slate-600 rounded-lg flex items-center justify-center flex-shrink-0">
                <Clock className="w-5 h-5 text-slate-600 dark:text-slate-300" />
              </div>
              <div>
                <h3 className="font-semibold text-slate-900 dark:text-white mb-1">Trazabilidad completa</h3>
                <p className="text-sm text-slate-600 dark:text-slate-400">
                  Todas las transacciones quedan registradas. Historial completo disponible para auditorías.
                </p>
              </div>
            </div>

          </div>

          {/* Platform disclaimer */}
          <div className="mt-12 max-w-4xl mx-auto">
            <div className="bg-slate-800 dark:bg-slate-900 rounded-xl p-6 text-center">
              <p className="text-slate-300 text-sm leading-relaxed">
                <span className="font-semibold text-white">Sorteos.club</span> es una plataforma tecnológica que facilita la gestión de sorteos entre usuarios.
                La responsabilidad de cada sorteo recae exclusivamente en el organizador que lo publica.
                Recomendamos a los participantes verificar siempre la información del organizador antes de participar.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-blue-600 to-blue-700 dark:from-blue-700 dark:to-blue-800">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            {isAuthenticated ? "Explora los sorteos disponibles" : "¿Listo para comenzar?"}
          </h2>
          <p className="text-blue-100 mb-8 max-w-2xl mx-auto text-lg">
            {isAuthenticated
              ? "Descubre los sorteos activos y participa por increíbles premios"
              : "Únete a Sorteos.club hoy y comienza a participar en sorteos transparentes y seguros"
            }
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            {isAuthenticated ? (
              <Link to="/explore">
                <Button size="lg" variant="outline" className="w-full sm:w-auto bg-white text-blue-600 hover:bg-slate-50 border-white">
                  Ver Sorteos Activos
                </Button>
              </Link>
            ) : (
              <>
                <Link to="/register">
                  <Button size="lg" variant="outline" className="w-full sm:w-auto bg-white text-blue-600 hover:bg-slate-50 border-white">
                    Crear Cuenta Gratis
                  </Button>
                </Link>
                <Link to="/login">
                  <Button size="lg" variant="outline" className="w-full sm:w-auto bg-white text-blue-600 hover:bg-slate-50 border-white">
                    Ya tengo cuenta
                  </Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-slate-900 dark:bg-slate-950 text-slate-400 py-12">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
                  </svg>
                </div>
                <span className="text-white font-semibold">Sorteos.club</span>
              </div>
              <p className="text-sm">
                La plataforma más confiable para sorteos en línea
              </p>
            </div>

            <div>
              <h4 className="text-white font-semibold mb-4">Plataforma</h4>
              <ul className="space-y-2 text-sm">
                {isAuthenticated ? (
                  <>
                    <li><Link to="/explore" className="hover:text-white transition-colors">Ver Sorteos</Link></li>
                    <li><Link to="/dashboard" className="hover:text-white transition-colors">Mi Panel</Link></li>
                    <li><Link to="/wallet" className="hover:text-white transition-colors">Mi Billetera</Link></li>
                  </>
                ) : (
                  <>
                    <li><Link to="/register" className="hover:text-white transition-colors">Crear Cuenta</Link></li>
                    <li><Link to="/login" className="hover:text-white transition-colors">Iniciar Sesión</Link></li>
                  </>
                )}
              </ul>
            </div>

            <div>
              <h4 className="text-white font-semibold mb-4">Legal</h4>
              <ul className="space-y-2 text-sm">
                <li><a href="#" className="hover:text-white transition-colors">Términos y Condiciones</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Política de Privacidad</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Política de Cookies</a></li>
              </ul>
            </div>

            <div>
              <h4 className="text-white font-semibold mb-4">Soporte</h4>
              <ul className="space-y-2 text-sm">
                <li><a href="#" className="hover:text-white transition-colors">Centro de Ayuda</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Contacto</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Estado del Servicio</a></li>
              </ul>
            </div>
          </div>

          <div className="border-t border-slate-800 mt-12 pt-8 text-center text-sm">
            <p>&copy; {new Date().getFullYear()} Sorteos.club. Todos los derechos reservados.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
