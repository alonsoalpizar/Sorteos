import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { useIsAuthenticated } from '@/hooks/useAuth';
import {
  ClipboardList,
  Users,
  Bell,
  Globe,
  BarChart3,
  Shield,
  CheckCircle,
  ArrowRight,
  Sparkles,
  FileText,
  Smartphone,
  Trophy,
  Eye,
  Clock,
  Shuffle
} from 'lucide-react';

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
                <ClipboardList className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold text-slate-900 dark:text-white">
                Sorteos.club
              </span>
            </div>

            <nav className="hidden md:flex items-center gap-6">
              <a href="#features" className="text-slate-600 hover:text-blue-600 dark:text-slate-300 dark:hover:text-blue-400 transition-colors">
                Características
              </a>
              <a href="#how-it-works" className="text-slate-600 hover:text-blue-600 dark:text-slate-300 dark:hover:text-blue-400 transition-colors">
                Cómo Funciona
              </a>
              {/* Pricing link - desactivado temporalmente
              <a href="#pricing" className="text-slate-600 hover:text-blue-600 dark:text-slate-300 dark:hover:text-blue-400 transition-colors">
                Planes
              </a>
              */}
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
                      Comenzar Gratis
                    </Button>
                  </Link>
                </>
              )}
            </nav>

            {/* Mobile menu */}
            <div className="md:hidden flex gap-2">
              {isAuthenticated ? (
                <Link to="/dashboard">
                  <Button size="sm">Mi Panel</Button>
                </Link>
              ) : (
                <Link to="/register">
                  <Button size="sm">Comenzar</Button>
                </Link>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section - ORGANIZADORES */}
      <section className="relative overflow-hidden">
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
                <Sparkles className="w-4 h-4" />
                Plataforma de Gestión de Sorteos
              </div>

              <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-slate-900 dark:text-white mb-6 leading-tight">
                Administra tus{' '}
                <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600">
                  sorteos
                </span>{' '}
                de forma profesional
              </h1>

              <p className="text-lg md:text-xl text-slate-600 dark:text-slate-300 mb-8 max-w-xl mx-auto lg:mx-0">
                Deja el cuaderno y WhatsApp. Registra compradores, gestiona pagos,
                comparte tu sorteo con un link público y notifica al ganador automáticamente.
              </p>

              <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
                {isAuthenticated ? (
                  <Link to="/organizer/raffles/new">
                    <Button size="lg" className="w-full sm:w-auto shadow-lg shadow-blue-500/25">
                      Crear Nuevo Sorteo
                      <ArrowRight className="w-5 h-5 ml-2" />
                    </Button>
                  </Link>
                ) : (
                  <>
                    <Link to="/register">
                      <Button size="lg" className="w-full sm:w-auto shadow-lg shadow-blue-500/25">
                        Crear Mi Primer Sorteo
                        <ArrowRight className="w-5 h-5 ml-2" />
                      </Button>
                    </Link>
                    <a href="#how-it-works">
                      <Button size="lg" variant="outline" className="w-full sm:w-auto">
                        Ver cómo funciona
                      </Button>
                    </a>
                  </>
                )}
              </div>

              {/* Trust indicators */}
              <div className="flex flex-wrap items-center justify-center lg:justify-start gap-6 mt-10 text-sm text-slate-500 dark:text-slate-400">
                <div className="flex items-center gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500" />
                  <span>100% Gratis para empezar</span>
                </div>
                <div className="flex items-center gap-2">
                  <Shield className="w-4 h-4 text-blue-500" />
                  <span>Sorteos verificables</span>
                </div>
                <div className="flex items-center gap-2">
                  <Globe className="w-4 h-4 text-indigo-500" />
                  <span>Link público para compartir</span>
                </div>
              </div>
            </div>

            {/* Right: Dashboard mockup */}
            <div className="relative hidden lg:block">
              {/* Main card - Dashboard preview */}
              <div className="relative z-10 bg-white dark:bg-slate-800 rounded-2xl shadow-2xl p-6 transform hover:scale-[1.02] transition-transform duration-500">
                {/* Header del dashboard */}
                <div className="flex items-center justify-between mb-6">
                  <div>
                    <h3 className="font-semibold text-slate-900 dark:text-white">Rifa iPhone 15 Pro</h3>
                    <p className="text-sm text-green-600 flex items-center gap-1">
                      <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                      Activo
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-2xl font-bold text-slate-900 dark:text-white">47/100</p>
                    <p className="text-xs text-slate-500">números vendidos</p>
                  </div>
                </div>

                {/* Stats */}
                <div className="grid grid-cols-3 gap-4 mb-6">
                  <div className="bg-blue-50 dark:bg-blue-900/30 rounded-lg p-3 text-center">
                    <p className="text-lg font-bold text-blue-600">₡235,000</p>
                    <p className="text-xs text-slate-500">Recaudado</p>
                  </div>
                  <div className="bg-green-50 dark:bg-green-900/30 rounded-lg p-3 text-center">
                    <p className="text-lg font-bold text-green-600">42</p>
                    <p className="text-xs text-slate-500">Pagados</p>
                  </div>
                  <div className="bg-amber-50 dark:bg-amber-900/30 rounded-lg p-3 text-center">
                    <p className="text-lg font-bold text-amber-600">5</p>
                    <p className="text-xs text-slate-500">Pendientes</p>
                  </div>
                </div>

                {/* Recent buyers list */}
                <div className="space-y-2">
                  <p className="text-xs font-medium text-slate-500 uppercase">Últimos registros</p>
                  {[
                    { name: 'María García', number: '#23', status: 'paid' },
                    { name: 'Carlos Rodríguez', number: '#45', status: 'paid' },
                    { name: 'Ana Mora', number: '#67', status: 'pending' },
                  ].map((buyer, i) => (
                    <div key={i} className="flex items-center justify-between py-2 border-b border-slate-100 dark:border-slate-700 last:border-0">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 bg-slate-200 dark:bg-slate-600 rounded-full flex items-center justify-center text-xs font-medium">
                          {buyer.name.split(' ').map(n => n[0]).join('')}
                        </div>
                        <div>
                          <p className="text-sm font-medium text-slate-900 dark:text-white">{buyer.name}</p>
                          <p className="text-xs text-slate-500">Número {buyer.number}</p>
                        </div>
                      </div>
                      <span className={`text-xs px-2 py-1 rounded-full ${
                        buyer.status === 'paid'
                          ? 'bg-green-100 text-green-700 dark:bg-green-900/50 dark:text-green-400'
                          : 'bg-amber-100 text-amber-700 dark:bg-amber-900/50 dark:text-amber-400'
                      }`}>
                        {buyer.status === 'paid' ? 'Pagado' : 'Pendiente'}
                      </span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Floating element - Link */}
              <div className="absolute -top-4 right-4 bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium shadow-lg z-20 flex items-center gap-2">
                <Globe className="w-4 h-4" />
                sorteos.club/r/iphone15
              </div>

              {/* Floating element - Notification */}
              <div className="absolute bottom-8 -left-4 bg-white dark:bg-slate-800 rounded-xl shadow-lg p-3 z-20 border border-slate-200 dark:border-slate-700">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 bg-green-100 dark:bg-green-900/50 rounded-full flex items-center justify-center">
                    <Bell className="w-4 h-4 text-green-600" />
                  </div>
                  <div>
                    <p className="text-xs text-slate-500">Nuevo pago</p>
                    <p className="text-sm font-medium text-slate-900 dark:text-white">+₡5,000</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Problems Section */}
      <section className="py-16 bg-slate-900 dark:bg-slate-950">
        <div className="container mx-auto px-4">
          <div className="text-center mb-12">
            <h2 className="text-2xl md:text-3xl font-bold text-white mb-4">
              ¿Te suena familiar?
            </h2>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-5xl mx-auto">
            {[
              { icon: FileText, text: "Anoto en un cuaderno quién me compró" },
              { icon: Smartphone, text: "Busco en WhatsApp quién me pagó" },
              { icon: Users, text: "No sé cuántos números vendí realmente" },
              { icon: Bell, text: "Tengo que avisar uno por uno del resultado" },
            ].map((problem, i) => (
              <div key={i} className="bg-slate-800 rounded-xl p-6 text-center border border-slate-700">
                <problem.icon className="w-8 h-8 text-red-400 mx-auto mb-3" />
                <p className="text-slate-300 text-sm">{problem.text}</p>
              </div>
            ))}
          </div>

          <div className="text-center mt-10">
            <p className="text-xl text-white font-medium">
              Sorteos.club <span className="text-blue-400">resuelve todo eso</span>
            </p>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 bg-white dark:bg-slate-800">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 dark:text-white mb-4">
              Todo lo que necesitas para gestionar tu sorteo
            </h2>
            <p className="text-slate-600 dark:text-slate-300 max-w-2xl mx-auto">
              Una plataforma completa para organizadores que quieren dejar de improvisar
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-6xl mx-auto">
            {/* Feature 1 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-blue-500/25">
                <ClipboardList className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Registro Digital
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Registra cada comprador con nombre, contacto y número asignado. Olvídate del cuaderno.
              </p>
            </div>

            {/* Feature 2 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-green-400 to-emerald-500 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-green-500/25">
                <BarChart3 className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Control de Pagos
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Marca quién pagó y quién no. Ve en tiempo real cuánto has recaudado.
              </p>
            </div>

            {/* Feature 3 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-indigo-400 to-indigo-600 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-indigo-500/25">
                <Globe className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Página Pública
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Cada sorteo tiene su link único para compartir. Los participantes ven qué números hay disponibles.
              </p>
            </div>

            {/* Feature 4 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-amber-400 to-orange-500 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-amber-500/25">
                <Shuffle className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Sorteo Verificable
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Resultado público y auditable. Transparencia total para ti y tus participantes.
              </p>
            </div>

            {/* Feature 5 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-purple-400 to-purple-600 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-purple-500/25">
                <Bell className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Notificaciones
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Avisa automáticamente al ganador y a todos los participantes cuando finalice el sorteo.
              </p>
            </div>

            {/* Feature 6 */}
            <div className="bg-slate-50 dark:bg-slate-900 p-8 rounded-2xl border border-slate-200 dark:border-slate-700 hover:shadow-lg hover:border-blue-200 dark:hover:border-blue-800 transition-all">
              <div className="w-12 h-12 bg-gradient-to-br from-teal-400 to-teal-600 rounded-xl flex items-center justify-center mb-4 shadow-lg shadow-teal-500/25">
                <Clock className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                Historial Completo
              </h3>
              <p className="text-slate-600 dark:text-slate-400">
                Todos tus sorteos anteriores guardados. Trazabilidad completa para futuras referencias.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How it Works */}
      <section id="how-it-works" className="py-20 bg-slate-50 dark:bg-slate-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 dark:text-white mb-4">
              Así de fácil es usar Sorteos.club
            </h2>
            <p className="text-slate-600 dark:text-slate-300 max-w-2xl mx-auto">
              En 5 simples pasos, gestiona tu rifa de principio a fin
            </p>
          </div>

          <div className="max-w-4xl mx-auto">
            {[
              {
                step: 1,
                title: "Crea tu sorteo",
                description: "Define el premio, cantidad de números, precio por número y fecha del sorteo.",
                color: "blue"
              },
              {
                step: 2,
                title: "Comparte tu link",
                description: "Recibe un link único (sorteos.club/r/tu-sorteo) para compartir en redes sociales.",
                color: "indigo"
              },
              {
                step: 3,
                title: "Registra compradores",
                description: "Cuando alguien te compre un número, regístralo en el sistema con sus datos.",
                color: "green"
              },
              {
                step: 4,
                title: "Marca los pagos",
                description: "Conforme te paguen por SINPE o transferencia, marca cada pago como recibido.",
                color: "amber"
              },
              {
                step: 5,
                title: "Sorteo y notificación",
                description: "El sorteo se realiza y se notifica al ganador automáticamente.",
                color: "purple"
              }
            ].map((item, i) => (
              <div key={i} className="flex gap-6 mb-8 last:mb-0">
                <div className={`w-12 h-12 bg-${item.color}-600 text-white rounded-full flex items-center justify-center text-xl font-bold flex-shrink-0`}>
                  {item.step}
                </div>
                <div className="pt-2">
                  <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-1">
                    {item.title}
                  </h3>
                  <p className="text-slate-600 dark:text-slate-400">
                    {item.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* For Participants Section */}
      <section className="py-20 bg-white dark:bg-slate-800">
        <div className="container mx-auto px-4">
          <div className="grid lg:grid-cols-2 gap-12 items-center max-w-6xl mx-auto">
            <div>
              <div className="inline-flex items-center gap-2 bg-green-100 dark:bg-green-900/40 text-green-700 dark:text-green-300 px-4 py-2 rounded-full text-sm font-medium mb-6">
                <Users className="w-4 h-4" />
                Para participantes
              </div>
              <h2 className="text-3xl md:text-4xl font-bold text-slate-900 dark:text-white mb-6">
                ¿Te invitaron a un sorteo?
              </h2>
              <p className="text-lg text-slate-600 dark:text-slate-300 mb-6">
                Si alguien te compartió un link de Sorteos.club, puedes verificar que el sorteo es legítimo y transparente.
              </p>

              <div className="space-y-4">
                {[
                  { icon: Eye, text: "Ve qué números están disponibles y cuáles vendidos" },
                  { icon: Shield, text: "Verifica la información del organizador" },
                  { icon: Trophy, text: "Consulta los resultados cuando el sorteo termine" },
                  { icon: CheckCircle, text: "Confirma que el sorteo fue justo y aleatorio" },
                ].map((item, i) => (
                  <div key={i} className="flex items-center gap-3">
                    <item.icon className="w-5 h-5 text-green-500" />
                    <span className="text-slate-700 dark:text-slate-300">{item.text}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="bg-slate-100 dark:bg-slate-900 rounded-2xl p-8">
              <div className="text-center">
                <Trophy className="w-16 h-16 text-amber-500 mx-auto mb-4" />
                <h3 className="text-xl font-semibold text-slate-900 dark:text-white mb-2">
                  Página pública del sorteo
                </h3>
                <p className="text-slate-600 dark:text-slate-400 mb-6">
                  Cada sorteo tiene su propia página donde puedes ver toda la información
                </p>
                <div className="bg-white dark:bg-slate-800 rounded-lg p-4 border border-slate-200 dark:border-slate-700">
                  <p className="text-sm text-slate-500 mb-1">Ejemplo de link:</p>
                  <code className="text-blue-600 dark:text-blue-400 font-mono">
                    sorteos.club/r/iphone-navidad-2025
                  </code>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Pricing Section - DESACTIVADO TEMPORALMENTE
      <section id="pricing" className="py-20 bg-slate-50 dark:bg-slate-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-slate-900 dark:text-white mb-4">
              Comienza gratis, crece con nosotros
            </h2>
            <p className="text-slate-600 dark:text-slate-300 max-w-2xl mx-auto">
              Planes diseñados para organizadores de todos los tamaños
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            <div className="bg-white dark:bg-slate-800 rounded-2xl p-8 border-2 border-slate-200 dark:border-slate-700">
              <h3 className="text-xl font-bold text-slate-900 dark:text-white mb-2">Gratis</h3>
              <p className="text-slate-600 dark:text-slate-400 mb-4">Para empezar</p>
              <p className="text-4xl font-bold text-slate-900 dark:text-white mb-6">₡0<span className="text-lg font-normal text-slate-500">/mes</span></p>

              <ul className="space-y-3 mb-8">
                {[
                  "Hasta 3 sorteos activos",
                  "Registro ilimitado de compradores",
                  "Link público para compartir",
                  "Sorteo aleatorio verificable",
                  "Historial básico"
                ].map((feature, i) => (
                  <li key={i} className="flex items-center gap-2 text-slate-600 dark:text-slate-400">
                    <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0" />
                    {feature}
                  </li>
                ))}
              </ul>

              <Link to="/register">
                <Button variant="outline" className="w-full">
                  Comenzar Gratis
                </Button>
              </Link>
            </div>

            <div className="bg-blue-600 rounded-2xl p-8 border-2 border-blue-500 relative">
              <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-amber-400 text-amber-900 text-xs font-bold px-3 py-1 rounded-full">
                PRÓXIMAMENTE
              </div>
              <h3 className="text-xl font-bold text-white mb-2">Pro</h3>
              <p className="text-blue-200 mb-4">Para organizadores frecuentes</p>
              <p className="text-4xl font-bold text-white mb-6">₡5,000<span className="text-lg font-normal text-blue-200">/mes</span></p>

              <ul className="space-y-3 mb-8">
                {[
                  "Sorteos ilimitados",
                  "Notificaciones automáticas",
                  "Analytics y reportes",
                  "Personalización de marca",
                  "Soporte prioritario"
                ].map((feature, i) => (
                  <li key={i} className="flex items-center gap-2 text-blue-100">
                    <CheckCircle className="w-5 h-5 text-blue-300 flex-shrink-0" />
                    {feature}
                  </li>
                ))}
              </ul>

              <Button variant="outline" className="w-full bg-white text-blue-600 hover:bg-blue-50 border-white" disabled>
                Próximamente
              </Button>
            </div>

            <div className="bg-white dark:bg-slate-800 rounded-2xl p-8 border-2 border-slate-200 dark:border-slate-700 relative">
              <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-amber-400 text-amber-900 text-xs font-bold px-3 py-1 rounded-full">
                PRÓXIMAMENTE
              </div>
              <h3 className="text-xl font-bold text-slate-900 dark:text-white mb-2">Business</h3>
              <p className="text-slate-600 dark:text-slate-400 mb-4">Para empresas y grandes organizadores</p>
              <p className="text-4xl font-bold text-slate-900 dark:text-white mb-6">₡15,000<span className="text-lg font-normal text-slate-500">/mes</span></p>

              <ul className="space-y-3 mb-8">
                {[
                  "Todo lo de Pro",
                  "Múltiples usuarios",
                  "API de integración",
                  "White-label",
                  "Soporte dedicado"
                ].map((feature, i) => (
                  <li key={i} className="flex items-center gap-2 text-slate-600 dark:text-slate-400">
                    <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0" />
                    {feature}
                  </li>
                ))}
              </ul>

              <Button variant="outline" className="w-full" disabled>
                Próximamente
              </Button>
            </div>
          </div>
        </div>
      </section>
      FIN Pricing Section - DESACTIVADO TEMPORALMENTE */}

      {/* Important Note */}
      <section className="py-12 bg-slate-800 dark:bg-slate-950">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center">
            <Shield className="w-12 h-12 text-blue-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-white mb-4">
              Importante: ¿Qué NO es Sorteos.club?
            </h3>
            <div className="grid md:grid-cols-3 gap-6 text-left">
              {[
                "NO vendemos números de rifa",
                "NO cobramos comisión por ventas",
                "NO procesamos pagos de compradores"
              ].map((item, i) => (
                <div key={i} className="flex items-center gap-2 text-slate-300">
                  <div className="w-6 h-6 bg-red-500/20 rounded-full flex items-center justify-center flex-shrink-0">
                    <span className="text-red-400 text-sm">✕</span>
                  </div>
                  {item}
                </div>
              ))}
            </div>
            <p className="mt-6 text-slate-400">
              Sorteos.club es una <span className="text-white font-medium">herramienta de gestión</span>.
              El organizador cobra y entrega el premio por su cuenta.
              Nosotros solo facilitamos la administración y transparencia.
            </p>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-blue-600 to-indigo-600">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            ¿Listo para organizar tu próximo sorteo?
          </h2>
          <p className="text-blue-100 mb-8 max-w-2xl mx-auto text-lg">
            Únete a los organizadores que ya confían en Sorteos.club para gestionar sus rifas de forma profesional
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            {isAuthenticated ? (
              <Link to="/organizer/raffles/new">
                <Button size="lg" className="bg-white text-blue-600 hover:bg-slate-50 shadow-lg">
                  Crear Nuevo Sorteo
                  <ArrowRight className="w-5 h-5 ml-2" />
                </Button>
              </Link>
            ) : (
              <Link to="/register">
                <Button size="lg" className="bg-white text-blue-600 hover:bg-slate-50 shadow-lg">
                  Crear Mi Cuenta Gratis
                  <ArrowRight className="w-5 h-5 ml-2" />
                </Button>
              </Link>
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
                  <ClipboardList className="w-5 h-5 text-white" />
                </div>
                <span className="text-white font-semibold">Sorteos.club</span>
              </div>
              <p className="text-sm">
                Plataforma de gestión de sorteos para organizadores que buscan transparencia y profesionalismo.
              </p>
            </div>

            <div>
              <h4 className="text-white font-semibold mb-4">Plataforma</h4>
              <ul className="space-y-2 text-sm">
                {isAuthenticated ? (
                  <>
                    <li><Link to="/organizer" className="hover:text-white transition-colors">Mi Panel</Link></li>
                    <li><Link to="/organizer/raffles/new" className="hover:text-white transition-colors">Crear Sorteo</Link></li>
                    <li><Link to="/my-raffles" className="hover:text-white transition-colors">Mis Sorteos</Link></li>
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
                <li><a href="#" className="hover:text-white transition-colors">Términos de Servicio</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Política de Privacidad</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Responsabilidad</a></li>
              </ul>
            </div>

            <div>
              <h4 className="text-white font-semibold mb-4">Soporte</h4>
              <ul className="space-y-2 text-sm">
                <li><a href="#" className="hover:text-white transition-colors">Centro de Ayuda</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Contacto</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Reportar Problema</a></li>
              </ul>
            </div>
          </div>

          <div className="border-t border-slate-800 mt-12 pt-8 text-center text-sm">
            <p>&copy; {new Date().getFullYear()} Sorteos.club. Todos los derechos reservados. Costa Rica.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
