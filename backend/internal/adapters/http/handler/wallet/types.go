package wallet

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
