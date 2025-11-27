package pagadito

import (
	"github.com/shopspring/decimal"
)

// Config configuración del cliente Pagadito
type Config struct {
	UID         string // Identificador del comercio
	WSK         string // Clave secreta
	SandboxMode bool   // true = sandbox, false = production
	APIURL      string // URL base de la API
	ReturnURL   string // URL de callback
}

// TransactionRequest solicitud para crear transacción
type TransactionRequest struct {
	ERN          string                 // External Reference Number (único)
	Amount       decimal.Decimal        // Monto total a cobrar
	Currency     string                 // CRC o USD
	Details      []TransactionDetail    // Detalles de la transacción
	CustomParams map[string]string      // Parámetros personalizados
	AllowPending bool                   // Permitir pagos pendientes
}

// TransactionDetail detalle de un item en la transacción
type TransactionDetail struct {
	Quantity    int             `json:"quantity"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	URLProduct  string          `json:"url_product,omitempty"`
}

// TransactionResponse respuesta al crear transacción
type TransactionResponse struct {
	Code       string // Código de respuesta (PG1002 = éxito)
	Message    string // Mensaje descriptivo
	Token      string // Token de transacción (para GetStatus)
	PaymentURL string // URL de pago donde redirigir al usuario
	DateTime   string // Fecha/hora de la respuesta
}

// StatusResponse respuesta al consultar estado de transacción
type StatusResponse struct {
	Code         string                 // Código de respuesta (PG1003 = éxito)
	Message      string                 // Mensaje descriptivo
	Status       string                 // COMPLETED, REGISTERED, VERIFYING, REVOKED, FAILED
	Reference    string                 // NAP (Número de Aprobación Pagadito)
	DateTrans    string                 // Fecha/hora de la transacción
	Amount       decimal.Decimal        // Monto de la transacción
	Currency     string                 // Moneda
	CustomParams map[string]string      // Parámetros personalizados devueltos
}

// apiResponse estructura genérica de respuesta de la API
type apiResponse struct {
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Value    interface{} `json:"value"`
	DateTime string      `json:"datetime"`
}

// Client interfaz para el cliente Pagadito
type Client interface {
	// Connect autentica con Pagadito y obtiene token de sesión
	Connect() error

	// CreateTransaction crea una transacción y retorna URL de pago
	CreateTransaction(req *TransactionRequest) (*TransactionResponse, error)

	// GetStatus consulta el estado de una transacción
	GetStatus(token string) (*StatusResponse, error)
}
