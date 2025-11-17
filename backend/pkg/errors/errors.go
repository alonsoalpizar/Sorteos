package errors

import (
	"fmt"
	"net/http"
)

// AppError representa un error de aplicación con código HTTP
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

// Error implementa la interfaz error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap permite usar errors.Is y errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// New crea un nuevo AppError
func New(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// Errores predefinidos - Authentication
var (
	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "No autorizado",
		Status:  http.StatusUnauthorized,
	}
	ErrInvalidCredentials = &AppError{
		Code:    "INVALID_CREDENTIALS",
		Message: "Credenciales inválidas",
		Status:  http.StatusUnauthorized,
	}
	ErrTokenExpired = &AppError{
		Code:    "TOKEN_EXPIRED",
		Message: "Token expirado",
		Status:  http.StatusUnauthorized,
	}
	ErrTokenInvalid = &AppError{
		Code:    "TOKEN_INVALID",
		Message: "Token inválido",
		Status:  http.StatusUnauthorized,
	}
)

// Errores predefinidos - Authorization
var (
	ErrForbidden = &AppError{
		Code:    "FORBIDDEN",
		Message: "No tienes permisos para esta operación",
		Status:  http.StatusForbidden,
	}
	ErrInsufficientKYC = &AppError{
		Code:    "INSUFFICIENT_KYC",
		Message: "Nivel de verificación insuficiente",
		Status:  http.StatusForbidden,
	}
)

// Errores predefinidos - Validación
var (
	ErrBadRequest = &AppError{
		Code:    "BAD_REQUEST",
		Message: "Solicitud inválida",
		Status:  http.StatusBadRequest,
	}
	ErrValidationFailed = &AppError{
		Code:    "VALIDATION_FAILED",
		Message: "Validación fallida",
		Status:  http.StatusBadRequest,
	}
	ErrEmailAlreadyExists = &AppError{
		Code:    "EMAIL_ALREADY_EXISTS",
		Message: "El email ya está registrado",
		Status:  http.StatusConflict,
	}
	ErrPhoneAlreadyExists = &AppError{
		Code:    "PHONE_ALREADY_EXISTS",
		Message: "El teléfono ya está registrado",
		Status:  http.StatusConflict,
	}
)

// Errores predefinidos - Recursos
var (
	ErrNotFound = &AppError{
		Code:    "NOT_FOUND",
		Message: "Recurso no encontrado",
		Status:  http.StatusNotFound,
	}
	ErrUserNotFound = &AppError{
		Code:    "USER_NOT_FOUND",
		Message: "Usuario no encontrado",
		Status:  http.StatusNotFound,
	}
	ErrRaffleNotFound = &AppError{
		Code:    "RAFFLE_NOT_FOUND",
		Message: "Sorteo no encontrado",
		Status:  http.StatusNotFound,
	}
	ErrCategoryNotFound = &AppError{
		Code:    "CATEGORY_NOT_FOUND",
		Message: "Categoría no encontrada",
		Status:  http.StatusNotFound,
	}
	ErrReservationNotFound = &AppError{
		Code:    "RESERVATION_NOT_FOUND",
		Message: "Reserva no encontrada",
		Status:  http.StatusNotFound,
	}
)

// Errores predefinidos - Concurrencia
var (
	ErrNumberAlreadyReserved = &AppError{
		Code:    "NUMBER_ALREADY_RESERVED",
		Message: "Uno o más números ya fueron reservados",
		Status:  http.StatusConflict,
	}
	ErrReservationExpired = &AppError{
		Code:    "RESERVATION_EXPIRED",
		Message: "La reserva ha expirado",
		Status:  http.StatusGone,
	}
	ErrLockAcquisitionFailed = &AppError{
		Code:    "LOCK_ACQUISITION_FAILED",
		Message: "No se pudo obtener el lock, intente nuevamente",
		Status:  http.StatusConflict,
	}
)

// Errores predefinidos - Rate Limiting
var (
	ErrTooManyRequests = &AppError{
		Code:    "TOO_MANY_REQUESTS",
		Message: "Demasiadas solicitudes, intente más tarde",
		Status:  http.StatusTooManyRequests,
	}
)

// Errores predefinidos - Pagos
var (
	ErrPaymentFailed = &AppError{
		Code:    "PAYMENT_FAILED",
		Message: "El pago falló",
		Status:  http.StatusPaymentRequired,
	}
	ErrPaymentAlreadyProcessed = &AppError{
		Code:    "PAYMENT_ALREADY_PROCESSED",
		Message: "El pago ya fue procesado",
		Status:  http.StatusConflict,
	}
	ErrStripeError = &AppError{
		Code:    "STRIPE_ERROR",
		Message: "Error en el procesador de pagos",
		Status:  http.StatusBadGateway,
	}
)

// Errores predefinidos - Sistema
var (
	ErrInternalServer = &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Error interno del servidor",
		Status:  http.StatusInternalServerError,
	}
	ErrDatabaseError = &AppError{
		Code:    "DATABASE_ERROR",
		Message: "Error de base de datos",
		Status:  http.StatusInternalServerError,
	}
	ErrRedisError = &AppError{
		Code:    "REDIS_ERROR",
		Message: "Error de Redis",
		Status:  http.StatusInternalServerError,
	}
)

// Wrap envuelve un error con contexto adicional
func Wrap(appErr *AppError, err error) *AppError {
	return &AppError{
		Code:    appErr.Code,
		Message: appErr.Message,
		Status:  appErr.Status,
		Err:     err,
	}
}

// WrapWithMessage envuelve un error con un mensaje personalizado
func WrapWithMessage(appErr *AppError, message string, err error) *AppError {
	return &AppError{
		Code:    appErr.Code,
		Message: message,
		Status:  appErr.Status,
		Err:     err,
	}
}
