package crypto

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/sorteos-platform/backend/pkg/errors"
)

const (
	// BcryptCost es el costo de bcrypt (2^12 = 4096 iteraciones)
	BcryptCost = 12
)

// HashPassword hashea una contraseña usando bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternalServer, err)
	}
	return string(bytes), nil
}

// ComparePassword compara una contraseña con un hash
func ComparePassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.ErrInvalidCredentials
		}
		return errors.Wrap(errors.ErrInternalServer, err)
	}
	return nil
}
