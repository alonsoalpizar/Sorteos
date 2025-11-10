package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config estructura principal de configuración
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Stripe   StripeConfig
	SendGrid SendGridConfig
	Twilio   TwilioConfig
	Business BusinessConfig
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port            string
	Environment     string // development, staging, production
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
}

// DatabaseConfig configuración de PostgreSQL
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig configuración de Redis
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// JWTConfig configuración de tokens JWT
type JWTConfig struct {
	Secret              string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
	Issuer              string
	RefreshTokenRotate  bool
}

// StripeConfig configuración de Stripe
type StripeConfig struct {
	SecretKey     string
	WebhookSecret string
	SuccessURL    string
	CancelURL     string
}

// SendGridConfig configuración de SendGrid (email)
type SendGridConfig struct {
	APIKey       string
	FromEmail    string
	FromName     string
	TemplatesDir string
}

// TwilioConfig configuración de Twilio (SMS)
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// BusinessConfig parámetros de negocio
type BusinessConfig struct {
	MaxActiveRafflesPerUser   int
	ReservationTTLMinutes     int
	PlatformFeePercentage     float64
	MinRaffleNumbers          int
	MaxRaffleNumbers          int
	MinPricePerNumber         float64
	RateLimitLoginPerMinute   int
	RateLimitReservePerMinute int
	RateLimitPaymentPerMinute int
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Intentar cargar archivo .env (opcional en producción)
	_ = viper.ReadInConfig()

	// Valores por defecto
	setDefaults()

	config := &Config{
		Server: ServerConfig{
			Port:            viper.GetString("CONFIG_SERVER_PORT"),
			Environment:     viper.GetString("CONFIG_ENVIRONMENT"),
			AllowedOrigins:  viper.GetStringSlice("CONFIG_ALLOWED_ORIGINS"),
			ShutdownTimeout: viper.GetDuration("CONFIG_SHUTDOWN_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("CONFIG_DB_HOST"),
			Port:            viper.GetInt("CONFIG_DB_PORT"),
			User:            viper.GetString("CONFIG_DB_USER"),
			Password:        viper.GetString("CONFIG_DB_PASSWORD"),
			DBName:          viper.GetString("CONFIG_DB_NAME"),
			SSLMode:         viper.GetString("CONFIG_DB_SSLMODE"),
			MaxOpenConns:    viper.GetInt("CONFIG_DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("CONFIG_DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("CONFIG_DB_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			Host:         viper.GetString("CONFIG_REDIS_HOST"),
			Port:         viper.GetInt("CONFIG_REDIS_PORT"),
			Password:     viper.GetString("CONFIG_REDIS_PASSWORD"),
			DB:           viper.GetInt("CONFIG_REDIS_DB"),
			PoolSize:     viper.GetInt("CONFIG_REDIS_POOL_SIZE"),
			MinIdleConns: viper.GetInt("CONFIG_REDIS_MIN_IDLE_CONNS"),
		},
		JWT: JWTConfig{
			Secret:              viper.GetString("CONFIG_JWT_SECRET"),
			AccessTokenExpiry:   viper.GetDuration("CONFIG_JWT_ACCESS_TOKEN_EXPIRY"),
			RefreshTokenExpiry:  viper.GetDuration("CONFIG_JWT_REFRESH_TOKEN_EXPIRY"),
			Issuer:              viper.GetString("CONFIG_JWT_ISSUER"),
			RefreshTokenRotate:  viper.GetBool("CONFIG_JWT_REFRESH_TOKEN_ROTATE"),
		},
		Stripe: StripeConfig{
			SecretKey:     viper.GetString("CONFIG_STRIPE_SECRET_KEY"),
			WebhookSecret: viper.GetString("CONFIG_STRIPE_WEBHOOK_SECRET"),
			SuccessURL:    viper.GetString("CONFIG_STRIPE_SUCCESS_URL"),
			CancelURL:     viper.GetString("CONFIG_STRIPE_CANCEL_URL"),
		},
		SendGrid: SendGridConfig{
			APIKey:       viper.GetString("CONFIG_SENDGRID_API_KEY"),
			FromEmail:    viper.GetString("CONFIG_SENDGRID_FROM_EMAIL"),
			FromName:     viper.GetString("CONFIG_SENDGRID_FROM_NAME"),
			TemplatesDir: viper.GetString("CONFIG_SENDGRID_TEMPLATES_DIR"),
		},
		Twilio: TwilioConfig{
			AccountSID: viper.GetString("CONFIG_TWILIO_ACCOUNT_SID"),
			AuthToken:  viper.GetString("CONFIG_TWILIO_AUTH_TOKEN"),
			FromNumber: viper.GetString("CONFIG_TWILIO_FROM_NUMBER"),
		},
		Business: BusinessConfig{
			MaxActiveRafflesPerUser:   viper.GetInt("CONFIG_RAFFLE_MAX_ACTIVE_PER_USER"),
			ReservationTTLMinutes:     viper.GetInt("CONFIG_RESERVATION_TTL_MINUTES"),
			PlatformFeePercentage:     viper.GetFloat64("CONFIG_PAYMENT_PLATFORM_FEE_PERCENTAGE"),
			MinRaffleNumbers:          viper.GetInt("CONFIG_RAFFLE_MIN_NUMBERS"),
			MaxRaffleNumbers:          viper.GetInt("CONFIG_RAFFLE_MAX_NUMBERS"),
			MinPricePerNumber:         viper.GetFloat64("CONFIG_RAFFLE_MIN_PRICE_PER_NUMBER"),
			RateLimitLoginPerMinute:   viper.GetInt("CONFIG_RATE_LIMIT_LOGIN_PER_MINUTE"),
			RateLimitReservePerMinute: viper.GetInt("CONFIG_RATE_LIMIT_RESERVE_PER_MINUTE"),
			RateLimitPaymentPerMinute: viper.GetInt("CONFIG_RATE_LIMIT_PAYMENT_PER_MINUTE"),
		},
	}

	// Validar configuración crítica
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// Server
	viper.SetDefault("CONFIG_SERVER_PORT", "8080")
	viper.SetDefault("CONFIG_ENVIRONMENT", "development")
	viper.SetDefault("CONFIG_ALLOWED_ORIGINS", []string{"http://localhost:5173"})
	viper.SetDefault("CONFIG_SHUTDOWN_TIMEOUT", "30s")

	// Database
	viper.SetDefault("CONFIG_DB_HOST", "localhost")
	viper.SetDefault("CONFIG_DB_PORT", 5432)
	viper.SetDefault("CONFIG_DB_SSLMODE", "disable")
	viper.SetDefault("CONFIG_DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("CONFIG_DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("CONFIG_DB_CONN_MAX_LIFETIME", "1h")

	// Redis
	viper.SetDefault("CONFIG_REDIS_HOST", "localhost")
	viper.SetDefault("CONFIG_REDIS_PORT", 6379)
	viper.SetDefault("CONFIG_REDIS_DB", 0)
	viper.SetDefault("CONFIG_REDIS_POOL_SIZE", 100)
	viper.SetDefault("CONFIG_REDIS_MIN_IDLE_CONNS", 10)

	// JWT
	viper.SetDefault("CONFIG_JWT_ACCESS_TOKEN_EXPIRY", "15m")
	viper.SetDefault("CONFIG_JWT_REFRESH_TOKEN_EXPIRY", "168h") // 7 días
	viper.SetDefault("CONFIG_JWT_ISSUER", "sorteos-platform")
	viper.SetDefault("CONFIG_JWT_REFRESH_TOKEN_ROTATE", true)

	// Business
	viper.SetDefault("CONFIG_RAFFLE_MAX_ACTIVE_PER_USER", 10)
	viper.SetDefault("CONFIG_RESERVATION_TTL_MINUTES", 5)
	viper.SetDefault("CONFIG_PAYMENT_PLATFORM_FEE_PERCENTAGE", 0.05)
	viper.SetDefault("CONFIG_RAFFLE_MIN_NUMBERS", 10)
	viper.SetDefault("CONFIG_RAFFLE_MAX_NUMBERS", 10000)
	viper.SetDefault("CONFIG_RAFFLE_MIN_PRICE_PER_NUMBER", 100.0)
	viper.SetDefault("CONFIG_RATE_LIMIT_LOGIN_PER_MINUTE", 5)
	viper.SetDefault("CONFIG_RATE_LIMIT_RESERVE_PER_MINUTE", 10)
	viper.SetDefault("CONFIG_RATE_LIMIT_PAYMENT_PER_MINUTE", 5)

	// SendGrid
	viper.SetDefault("CONFIG_SENDGRID_FROM_NAME", "Sorteos Platform")
	viper.SetDefault("CONFIG_SENDGRID_TEMPLATES_DIR", "./templates/email")
}

// Validate valida que la configuración sea correcta
func (c *Config) Validate() error {
	// Validar JWT secret en producción
	if c.Server.Environment == "production" && len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters in production")
	}

	// Validar database
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	// Validar Redis
	if c.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}

	return nil
}

// DSN retorna el Data Source Name para PostgreSQL
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// RedisAddr retorna la dirección de Redis
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDevelopment retorna true si está en modo development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction retorna true si está en modo production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}
