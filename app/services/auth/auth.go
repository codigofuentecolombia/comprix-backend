package service_auth

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/repositories"
	service_email "comprix/app/services/email"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	config  dto.Config
	service service_email.Service
}

func InitService(config dto.Config) Service {
	return Service{
		config:  config,
		service: service_email.InitService(config),
	}
}

func (a *Service) CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fails.Create("Auth CheckPasswordHash: ", err)
	}
	// Regresar nil
	return nil
}

func (a *Service) GenerateToken(user dao.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	// Token firmado
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.config.Settings.Server.Sk))
	// Verificar si hubo error
	if err != nil {
		return "", fails.Create("Auth GenerateToken: ", err)
	}
	// Regresar token firmado
	return signedToken, nil
}

func (a *Service) HandleLogin(authUser dao.User) (dto.AuthResponse, error) {
	repository := repositories.InitUserRepository(a.config.GormDB)
	// Obtener usuario
	user, err := repository.FindByUsername(authUser.Username)
	// Buscar usuario por email en la BD
	if err != nil {
		if authUser.Password == "" {
			password, err := a.GenerateRandomPassword()
			// Validar si hubo error al generar
			if err != nil {
				return dto.AuthResponse{}, fails.Create("HandleLogin: No se pudo crear el usuario.", err)
			}
			// Actualizar contraseña
			authUser.Password = password
		}
		// Encriptar contraseña
		encryptedPass, err := a.HashPassword(authUser.Password)
		// Validar si hubo error al generar
		if err != nil {
			return dto.AuthResponse{}, fails.Create("HandleLogin: No se pudo crear el usuario.", err)
		}
		//
		password := authUser.Password
		// Actualizar structura
		authUser.RoleID = 2
		authUser.Password = encryptedPass
		// Crear usuario
		_, err = repository.Create(authUser)
		// Verificar si  ocurrio un error al crear el usuario
		if err != nil {
			return dto.AuthResponse{}, fails.Create("HandleLogin: No se pudo crear el usuario.", err)
		}
		// Enviar email
		a.service.Register(authUser, password)
		// Buscar usuario actualizado
		user, err = repository.FindByUsername(authUser.Username)
		// Verificar si  ocurrio un error
		if err != nil {
			return dto.AuthResponse{}, fails.Create("HandleLogin: No se encontro usuario.", err)
		}
	}
	// Generar token JWT
	token, err := a.GenerateToken(*user)
	// Verificar si  ocurrio un error
	if err != nil {
		return dto.AuthResponse{}, fails.Create("HandleLogin: Error al generar el token.", err)
	}
	// Regresar respuesta
	return dto.AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}
