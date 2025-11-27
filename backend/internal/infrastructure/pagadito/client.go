package pagadito

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// HTTPClient implementación del cliente HTTP para Pagadito
type HTTPClient struct {
	config     *Config
	httpClient *http.Client
	token      string    // Token de autenticación (después de Connect)
	tokenTime  time.Time // Timestamp del token (para expiración)

	// Claves de operación (constantes de la API de Pagadito)
	opConnectKey      string
	opExecTransKey    string
	opGetStatusKey    string
	opExchangeRateKey string
}

// NewHTTPClient crea un nuevo cliente HTTP para Pagadito
func NewHTTPClient(config *Config) *HTTPClient {
	// URLs por defecto según modo
	apiURL := config.APIURL
	if apiURL == "" {
		if config.SandboxMode {
			apiURL = "https://sandbox.pagadito.com/comercios/apipg/charges.php"
		} else {
			apiURL = "https://comercios.pagadito.com/apipg/charges.php"
		}
	}

	return &HTTPClient{
		config: &Config{
			UID:         config.UID,
			WSK:         config.WSK,
			SandboxMode: config.SandboxMode,
			APIURL:      apiURL,
			ReturnURL:   config.ReturnURL,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Claves de operación de la API Pagadito (constantes de la documentación)
		opConnectKey:      "f3f191ce3326905ff4403bb05b0de150",
		opExecTransKey:    "41216f8caf94aaa598db137e36d4673e",
		opGetStatusKey:    "0b50820c65b0de71ce78f6221a5cf876",
		opExchangeRateKey: "da6b597cfcd0daf129287758b3c73b76",
	}
}

// Connect autentica con Pagadito y obtiene token de sesión
func (c *HTTPClient) Connect() error {
	params := url.Values{}
	params.Set("operation", c.opConnectKey)
	params.Set("uid", c.config.UID)
	params.Set("wsk", c.config.WSK)
	params.Set("format_return", "json")

	resp, err := c.call(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error en llamada Connect: %w", err)
	}

	if resp.Code != "PG1001" {
		return MapErrorCode(resp.Code, resp.Message)
	}

	// Guardar token de sesión
	tokenValue, ok := resp.Value.(string)
	if !ok {
		return fmt.Errorf("respuesta de Pagadito inválida: value no es string")
	}

	c.token = tokenValue
	c.tokenTime = time.Now()

	return nil
}

// CreateTransaction crea una transacción en Pagadito y retorna URL de pago
func (c *HTTPClient) CreateTransaction(req *TransactionRequest) (*TransactionResponse, error) {
	// Verificar conexión
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}

	// Formatear detalles como JSON (según ejemplo oficial Pagadito.php línea 183-188)
	// El orden DEBE ser: quantity, description, price, url_product
	// price DEBE ser float64, no string
	detailsJSON := make([]map[string]interface{}, len(req.Details))
	for i, detail := range req.Details {
		priceFloat, _ := detail.Price.Float64()
		detailsJSON[i] = map[string]interface{}{
			"quantity":    detail.Quantity,
			"description": detail.Description,
			"price":       priceFloat, // DEBE ser número, no string
			"url_product": detail.URLProduct,
		}
	}
	detailsBytes, err := json.Marshal(detailsJSON)
	if err != nil {
		return nil, fmt.Errorf("error serializando detalles: %w", err)
	}

	// Formatear custom params como JSON (según ejemplo oficial)
	customParamsBytes, err := json.Marshal(req.CustomParams)
	if err != nil {
		return nil, fmt.Errorf("error serializando custom_params: %w", err)
	}

	// Preparar parámetros
	params := url.Values{}
	params.Set("operation", c.opExecTransKey)
	params.Set("token", c.token)
	params.Set("ern", req.ERN)
	params.Set("amount", req.Amount.StringFixed(2))
	params.Set("details", string(detailsBytes))
	params.Set("custom_params", string(customParamsBytes))
	params.Set("currency", req.Currency)
	params.Set("format_return", "json")

	if req.AllowPending {
		params.Set("allow_pending_payments", "true")
	} else {
		params.Set("allow_pending_payments", "false")
	}

	// LOG: Imprimir a consola (aparecerá en journalctl)
	log.Println("========================================")
	log.Println("PAGADITO CreateTransaction REQUEST")
	log.Println("========================================")
	log.Printf("ERN: %s", req.ERN)
	log.Printf("Amount: %s", req.Amount.StringFixed(2))
	log.Printf("Currency: %s", req.Currency)
	log.Printf("Details JSON: %s", string(detailsBytes))
	log.Printf("Custom Params JSON: %s", string(customParamsBytes))
	log.Printf("Allow Pending: %v", req.AllowPending)
	log.Printf("Full POST params: %s", params.Encode())
	log.Println("========================================")

	// Ejecutar llamada
	resp, err := c.call(context.Background(), params)
	if err != nil {
		log.Printf("ERROR en llamada HTTP: %v", err)
		return nil, fmt.Errorf("error en llamada CreateTransaction: %w", err)
	}

	// LOG: Imprimir respuesta
	log.Println("PAGADITO CreateTransaction RESPONSE")
	log.Printf("Code: %s", resp.Code)
	log.Printf("Message: %s", resp.Message)
	log.Printf("Value: %+v", resp.Value)
	log.Println("========================================")

	if resp.Code != "PG1002" {
		return nil, MapErrorCode(resp.Code, resp.Message)
	}

	// Extraer URL de pago del value
	paymentURL, ok := resp.Value.(string)
	if !ok {
		return nil, fmt.Errorf("respuesta de Pagadito inválida: value no es string")
	}

	return &TransactionResponse{
		Code:       resp.Code,
		Message:    resp.Message,
		Token:      req.ERN, // Usamos el ERN como token para GetStatus
		PaymentURL: paymentURL,
		DateTime:   resp.DateTime,
	}, nil
}

// GetStatus consulta el estado de una transacción
func (c *HTTPClient) GetStatus(token string) (*StatusResponse, error) {
	// Verificar conexión
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("operation", c.opGetStatusKey)
	params.Set("token", c.token)
	params.Set("token_trans", token)
	params.Set("format_return", "json")

	resp, err := c.call(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("error en llamada GetStatus: %w", err)
	}

	if resp.Code != "PG1003" {
		return nil, MapErrorCode(resp.Code, resp.Message)
	}

	// Parsear value como mapa
	valueMap, ok := resp.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("respuesta de Pagadito inválida: value no es mapa")
	}

	status, _ := valueMap["status"].(string)
	reference, _ := valueMap["reference"].(string)
	dateTrans, _ := valueMap["date_trans"].(string)

	// Parsear amount si existe
	var amount decimal.Decimal
	if amountVal, ok := valueMap["amount"]; ok {
		if amountStr, ok := amountVal.(string); ok {
			amount, _ = decimal.NewFromString(amountStr)
		} else if amountFloat, ok := amountVal.(float64); ok {
			amount = decimal.NewFromFloat(amountFloat)
		}
	}

	currency, _ := valueMap["currency"].(string)

	return &StatusResponse{
		Code:      resp.Code,
		Message:   resp.Message,
		Status:    status,
		Reference: reference,
		DateTrans: dateTrans,
		Amount:    amount,
		Currency:  currency,
	}, nil
}

// ensureConnected verifica que haya una conexión activa, si no la crea
func (c *HTTPClient) ensureConnected() error {
	// Si no hay token o expiró (más de 30 minutos), reconectar
	if c.token == "" || time.Since(c.tokenTime) > 30*time.Minute {
		return c.Connect()
	}
	return nil
}

// call ejecuta una llamada HTTP a la API de Pagadito
func (c *HTTPClient) call(ctx context.Context, params url.Values) (*apiResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.config.APIURL,
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status inválido: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("error parseando JSON: %w", err)
	}

	return &apiResp, nil
}
