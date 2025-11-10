ARQUITECTURA DETALLADA Y DECISIONES DE DISEÑO - PLATAFORMA DE SORTEOS
1. MANEJO DE CONCURRENCIA - EL PROBLEMA CENTRAL
El Problema Real
Imagina este escenario:
- Sorteo popular: iPhone 15 a ₡100 el número
- 100 números disponibles
- 500 personas intentando comprar simultáneamente
- Cada persona quiere 2-3 números

SIN control de concurrencia:
- Usuario A y B ven número 42 disponible
- Ambos clickean "comprar" al mismo tiempo
- Ambos pagan
- ¿Quién se queda con el número 42?
La Solución: Sistema de 3 Capas
go// CAPA 1: Lock Distribuido con Redis (Microsegundos)
func ReserveNumbers(sorteoID, userID string, numbers []int) error {
    // Redis actúa como semáforo ultrarrápido
    lockKey := fmt.Sprintf("sorteo:%s:lock", sorteoID)
    
    // Solo UN usuario puede tener el lock a la vez
    // Si 100 personas clickean al mismo tiempo, Redis las ordena en fila
    locked := redis.SetNX(lockKey, userID, 100*time.Millisecond)
    if !locked {
        return errors.New("Por favor intente de nuevo") // Usuario va al final de la fila
    }
    defer redis.Del(lockKey) // Liberar lock al terminar
    
    // CAPA 2: Verificación en Base de Datos (Milisegundos)
    tx := db.Begin() // Transacción SQL
    for _, num := range numbers {
        // PostgreSQL verifica si el número está disponible
        var count int
        tx.Raw("SELECT COUNT(*) FROM sorteo_numbers WHERE sorteo_id = ? AND number = ? AND status != 'available'", 
            sorteoID, num).Scan(&count)
        
        if count > 0 {
            tx.Rollback()
            return fmt.Errorf("Número %d ya fue tomado", num)
        }
    }
    
    // CAPA 3: Reserva Temporal (5 minutos para pagar)
    for _, num := range numbers {
        // Marcar como "reserved" no "sold"
        tx.Exec("UPDATE sorteo_numbers SET status = 'reserved', user_id = ?, reserved_until = ? 
                WHERE sorteo_id = ? AND number = ?", 
                userID, time.Now().Add(5*time.Minute), sorteoID, num)
    }
    tx.Commit()
    
    return nil
}
2. REDIS - ¿POR QUÉ ES CRÍTICO?
Redis no es solo cache, es nuestro coordinador de tráfico:
javascript// Lo que hace Redis en el sistema:

1. LOCKS DISTRIBUIDOS (más importante)
   - Previene race conditions
   - Ordena las solicitudes
   - Tiempo de lock: 50-100ms máximo

2. RESERVAS TEMPORALES
   Redis Key: "reservation:user123:sorteo456"
   Value: {
     numbers: [15, 42, 73],
     expires_at: "2025-01-15T10:30:00Z",
     amount: 300.00
   }
   TTL: 5 minutos

3. RATE LIMITING
   Key: "rate:user123"
   - Max 10 intentos por minuto
   - Previene ataques/bots

4. SESIONES ACTIVAS
   Key: "session:token123"
   - JWT tokens
   - User data cache
```

### Flujo Visual de Concurrencia
```
Tiempo | Usuario A        | Usuario B        | Usuario C
-------|-----------------|-----------------|------------------
00.000 | Click número 42 | Click número 42 | Click número 42
00.001 | Obtiene lock ✓  | Espera lock...  | Espera lock...
00.005 | Verifica DB     | Esperando...    | Esperando...
00.010 | Reserva número  | Esperando...    | Esperando...
00.015 | Lock liberado   | Obtiene lock ✓  | Espera lock...
00.016 | Paga...         | DB dice "NO"    | Esperando...
00.020 |                 | Error: ocupado  | Obtiene lock ✓
00.025 |                 |                 | DB dice "NO"
3. SEGURIDAD - MULTICAPA
go// NIVEL 1: Autenticación
type SecurityLayers struct {
    // Passwords
    PasswordHashing: "bcrypt con cost 12", // 2^12 iteraciones
    MinLength: 12,
    RequireComplexity: true, // mayúsculas, números, símbolos
    
    // JWT Tokens
    TokenExpiry: 15 * time.Minute,
    RefreshTokenExpiry: 7 * 24 * time.Hour,
    SecretRotation: "mensual",
    
    // Rate Limiting
    LoginAttempts: 5,        // por IP
    PurchaseAttempts: 10,    // por usuario
    TimeWindow: 5 * time.Minute,
}

// NIVEL 2: Validación de Inputs
func ValidateInput(input interface{}) error {
    // Prevenir SQL Injection (GORM ya lo hace)
    // Prevenir XSS
    // Validar rangos
    // Sanitizar HTML
}

// NIVEL 3: Autorización
type UserTrustLevels struct {
    Level1: "Email verificado",           // Puede ver sorteos
    Level2: "Teléfono verificado",       // Puede comprar hasta ₡10,000
    Level3: "Cédula verificada",         // Puede comprar hasta ₡50,000
    Level4: "Dirección verificada",      // Puede comprar hasta ₡200,000
    Level5: "PowerUser (>20 compras)",   // Sin límites
}

// NIVEL 4: Auditoría
type AuditLog struct {
    UserID    string
    Action    string // "reserve_numbers", "purchase", "create_sorteo"
    IP        string
    UserAgent string
    Timestamp time.Time
    Result    string // "success", "failed"
    Details   json.RawMessage
}
4. BACKOFFICE DEL USUARIO
typescript// Dashboard con 4 secciones principales

interface UserDashboard {
  // 1. MIS SORTEOS (Como vendedor)
  misSorteos: {
    activos: Sorteo[],      // En curso
    completados: Sorteo[],  // Finalizados
    borradores: Sorteo[],   // Sin publicar
    estadisticas: {
      ventasTotales: number,
      numerosSoldidos: number,
      tasaConversion: number,
      proximoPago: Date
    }
  },
  
  // 2. MIS PARTICIPACIONES (Como comprador)
  misParticipaciones: {
    actuales: Participacion[],    // Sorteos activos donde participo
    ganados: Premio[],            // Histórico de premios
    gastadoTotal: number,
    numerosFavoritos: number[]    // Estadística personal
  },
  
  // 3. MI PERFIL
  perfil: {
    datosPersonales: UserData,
    verificacion: {
      email: boolean,
      telefono: boolean,
      cedula: boolean,
      direccion: boolean
    },
    trustLevel: number,           // 1-5
    limiteCompra: number,         // Según trust level
    metodosPago: PaymentMethod[]
  },
  
  // 4. FINANZAS
  finanzas: {
    saldoPendiente: number,       // Por cobrar de sorteos
    historialPagos: Payment[],
    comisiones: {
      porcentaje: number,         // 5-10%
      totalPagado: number
    },
    proximoCobro: Date
  }
}
Flujo de Autogestión de Sorteos
javascript// CREAR SORTEO - Proceso paso a paso

// Paso 1: Información básica
const createSorteo = {
  title: "iPhone 15 Pro Max 256GB",
  description: "Nuevo, sellado, con factura",
  productValue: 1200000, // ₡1.2M
  
  // Paso 2: Configuración de números
  pricePerNumber: 5000,  // ₡5,000 por número
  totalNumbers: 100,     // Se necesitan vender 60% para cubrir
  
  // Paso 3: Fechas
  startDate: "2025-01-20",
  endDate: "2025-02-20",
  lotteryDate: "2025-02-21", // Lotería Nacional
  lotteryType: "JPS",
  
  // Paso 4: Verificación
  images: [img1, img2, img3], // Mínimo 3 fotos
  invoice: "factura.pdf",     // Prueba de propiedad
  
  // Paso 5: Revisión y publicación
  status: "draft" // -> "pending_review" -> "active"
}

// GESTIONAR SORTEO ACTIVO
const manageSorteo = {
  // Ver estadísticas en tiempo real
  stats: {
    numerosSold: 67,
    numerosAvailable: 33,
    montoRecaudado: 335000,
    tiempoRestante: "15 días",
    velocidadVenta: "4.5 números/día"
  },
  
  // Acciones disponibles
  actions: {
    extenderFecha: true,      // Si ventas lentas
    cancelar: true,           // Devolver dinero
    modificarPrecio: false,   // No después de primera venta
    agregarFotos: true,
    enviarRecordatorio: true  // A participantes
  }
}
5. SISTEMA DE LIMPIEZA AUTOMÁTICA
go// Worker que corre cada minuto
func CleanupWorker() {
    for {
        time.Sleep(1 * time.Minute)
        
        // 1. Liberar reservas expiradas
        db.Exec(`
            UPDATE sorteo_numbers 
            SET status = 'available', 
                user_id = NULL
            WHERE status = 'reserved' 
            AND reserved_until < NOW()
        `)
        
        // 2. Notificar usuarios con reservas por expirar
        db.Query(`
            SELECT user_id, sorteo_id, reserved_until
            FROM sorteo_numbers
            WHERE status = 'reserved'
            AND reserved_until BETWEEN NOW() AND NOW() + INTERVAL '1 minute'
        `).Scan(&expiring)
        
        for _, reservation := range expiring {
            sendNotification(reservation.UserID, "Tu reserva expira en 1 minuto!")
        }
        
        // 3. Actualizar estadísticas de sorteos
        updateSorteoStats()
    }
}
6. FLUJO COMPLETO DE COMPRA - DETALLADO
mermaidgraph TD
    A[Usuario ve sorteo] --> B[Selecciona números]
    B --> C{Números disponibles?}
    C -->|No| D[Mostrar error]
    C -->|Sí| E[Redis Lock 100ms]
    E --> F[Verificar en DB]
    F --> G{Aún disponibles?}
    G -->|No| H[Liberar lock y error]
    G -->|Sí| I[Crear reserva 5min]
    I --> J[Liberar lock]
    J --> K[Mostrar página de pago]
    K --> L[Timer 5 minutos]
    L --> M{Pago completado?}
    M -->|No| N[Liberar números]
    M -->|Sí| O[Marcar como vendidos]
    O --> P[Email confirmación]
    P --> Q[Generar PDF comprobante]
7. ESCALABILIDAD Y PERFORMANCE
go// Optimizaciones implementadas

// 1. Connection Pooling
dbConfig := &gorm.Config{
    MaxIdleConns: 10,
    MaxOpenConns: 100,
    ConnMaxLifetime: time.Hour,
}

// 2. Índices estratégicos
CREATE INDEX CONCURRENTLY idx_numbers_available 
ON sorteo_numbers(sorteo_id, status) 
WHERE status = 'available';

// 3. Paginación eficiente
func ListSorteos(page, limit int) {
    offset := (page - 1) * limit
    db.Raw(`
        SELECT * FROM sorteos 
        WHERE status = 'active'
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `, limit, offset)
}

// 4. Cache de consultas frecuentes
func GetPopularSorteos() []Sorteo {
    // Cachear en Redis por 5 minutos
    cached := redis.Get("popular_sorteos")
    if cached != nil {
        return cached
    }
    
    sorteos := db.Query("SELECT ... ORDER BY numbers_sold DESC LIMIT 10")
    redis.Set("popular_sorteos", sorteos, 5*time.Minute)
    return sorteos
}
8. MÉTRICAS Y MONITOREO
sql-- Dashboard de métricas en tiempo real
CREATE VIEW sorteo_metrics AS
SELECT 
    s.id,
    s.title,
    COUNT(sn.id) as total_numeros,
    COUNT(CASE WHEN sn.status = 'sold' THEN 1 END) as vendidos,
    COUNT(CASE WHEN sn.status = 'reserved' THEN 1 END) as reservados,
    COUNT(CASE WHEN sn.status = 'available' THEN 1 END) as disponibles,
    (COUNT(CASE WHEN sn.status = 'sold' THEN 1 END) * 100.0 / COUNT(sn.id)) as porcentaje_vendido,
    SUM(CASE WHEN sn.status = 'sold' THEN s.price_per_number ELSE 0 END) as total_recaudado,
    AVG(EXTRACT(EPOCH FROM (sn.purchased_at - sn.reserved_at))) as tiempo_promedio_compra
FROM sorteos s
LEFT JOIN sorteo_numbers sn ON s.id = sn.sorteo_id
GROUP BY s.id, s.title;
CAPACIDAD DEL SISTEMA
Con esta arquitectura, el sistema puede manejar:

1000+ usuarios concurrentes comprando números
100,000+ transacciones por día
Tiempo de respuesta < 200ms para reservas
99.9% uptime con Redis como fallback


Documento de Arquitectura - Plataforma de Sorteos v1.0
Autor: Ing. Alonso Alpízar
Fecha: Noviembre 2025