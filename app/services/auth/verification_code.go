package service_auth

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"time"

	"golang.org/x/exp/rand"
)

func (a *Service) GenerateVerificationCode(user dao.User) (dao.VerificationCode, error) {
	repository := repositories.InitVerificationCodeRepository(a.config.GormDB)
	// Crear
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const length = 6
	// Crear code
	rand.Seed(uint64(time.Now().UnixNano()))
	code := make([]byte, length)
	// Iterar
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	//
	now := time.Now()             // Hora actual
	newTime := now.Add(time.Hour) // Sumar 1 hora
	// Crear entidad
	return repository.Create(dao.VerificationCode{
		Code:       string(code),
		UserID:     user.ID,
		Expiration: newTime,
	})
}

func (a *Service) VerifyCode(params dto.AuthVerificationCode) bool {
	userRepository := repositories.InitUserRepository(a.config.GormDB)
	// Obtener usuario
	user, err := userRepository.FindByEmail(params.Email)
	// Verificar si hubo error
	if err != nil {
		return false
	}
	//
	codeRepository := repositories.InitVerificationCodeRepository(a.config.GormDB)
	// Buscar
	code, err := codeRepository.FindUsersCode(user.ID, params.Code)
	// Verificar
	if err != nil {
		return false
	}
	// Actualizar
	return userRepository.VerifyCode(*user, *code) != nil
}
