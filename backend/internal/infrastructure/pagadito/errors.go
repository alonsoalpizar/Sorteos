package pagadito

import (
	"errors"
	"fmt"
)

// Errores comunes de Pagadito
var (
	ErrIncompleteData     = errors.New("PG2001: datos incompletos")
	ErrConnectionProblem  = errors.New("PG3001: problema de conexión con Pagadito")
	ErrGeneralError       = errors.New("PG3002: error general de Pagadito")
	ErrUnregistered       = errors.New("PG3003: transacción no registrada")
	ErrMatchError         = errors.New("PG3004: error de coincidencia")
	ErrDisabledConnection = errors.New("PG3005: conexión deshabilitada")
	ErrLimitExceeded      = errors.New("PG3006: límite excedido")
	ErrNotConnected       = errors.New("cliente no conectado - llame Connect() primero")
)

// MapErrorCode mapea códigos de error de Pagadito a errores Go
func MapErrorCode(code, message string) error {
	switch code {
	case "PG1001":
		// Conexión exitosa - no es error
		return nil
	case "PG1002":
		// Transacción registrada exitosamente - no es error
		return nil
	case "PG1003":
		// Estado obtenido exitosamente - no es error
		return nil
	case "PG2001":
		return fmt.Errorf("%w: %s", ErrIncompleteData, message)
	case "PG3001":
		return fmt.Errorf("%w: %s", ErrConnectionProblem, message)
	case "PG3002":
		return fmt.Errorf("%w: %s", ErrGeneralError, message)
	case "PG3003":
		return fmt.Errorf("%w: %s", ErrUnregistered, message)
	case "PG3004":
		return fmt.Errorf("%w: %s", ErrMatchError, message)
	case "PG3005":
		return fmt.Errorf("%w: %s", ErrDisabledConnection, message)
	case "PG3006":
		return fmt.Errorf("%w: %s", ErrLimitExceeded, message)
	default:
		return fmt.Errorf("código de error desconocido %s: %s", code, message)
	}
}
