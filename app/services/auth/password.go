package service_auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// Generar una contraseña aleatoria
func (svc *Service) GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 16) // 16 bytes (128 bits)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Hashear la contraseña antes de guardarla en la base de datos
func (svc *Service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
