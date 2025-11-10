package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateVerificationCode genera un código de verificación de 6 dígitos
func GenerateVerificationCode() (string, error) {
	// Generar número aleatorio entre 0 y 999999
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Formatear a 6 dígitos (con ceros a la izquierda)
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// GeneratePasswordResetToken genera un token de reset de contraseña
func GeneratePasswordResetToken() (string, error) {
	// Generar 32 bytes aleatorios
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convertir a hexadecimal
	return fmt.Sprintf("%x", bytes), nil
}
