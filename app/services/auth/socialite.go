package service_auth

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func (svc *Service) GenerateFacebookAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
		ClientID:     svc.config.Settings.Auth.Facebook.Client,
		RedirectURL:  svc.config.Settings.Auth.Facebook.Callback,
		ClientSecret: svc.config.Settings.Auth.Facebook.Secret,
	}
}

func (svc *Service) GenerateGoogleAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
		ClientID:     svc.config.Settings.Auth.Google.Client,
		RedirectURL:  svc.config.Settings.Auth.Google.Callback,
		ClientSecret: svc.config.Settings.Auth.Google.Secret,
	}
}

func (svc *Service) GenerateAuthConfigByType(sType dto.SocialiteType) *oauth2.Config {
	if sType == dto.SocialiteFacebookType {
		return svc.GenerateFacebookAuthConfig()
	} else {
		return svc.GenerateGoogleAuthConfig()
	}
}

func (svc *Service) GetGoogleDetailUri() string {
	return "https://www.googleapis.com/oauth2/v3/userinfo"
}

func (svc *Service) GetFacebookDetailUri() string {
	return "https://graph.facebook.com/me?fields=id,name,email"
}

func (svc *Service) GetDetailUriByType(sType dto.SocialiteType) string {
	if sType == dto.SocialiteFacebookType {
		return svc.GetFacebookDetailUri()
	} else {
		return svc.GetGoogleDetailUri()
	}
}

func (svc *Service) DecodeSocialiteUser(sType dto.SocialiteType, body io.ReadCloser) (*dao.User, error) {
	var socialiteUser interface{}
	// Verificar type
	if sType == dto.SocialiteFacebookType {
		socialiteUser = &dto.AuthFacebookUser{}
	} else {
		socialiteUser = &dto.AuthGoogleUser{}
	}
	// Decodificar el JSON en la estructura correcta
	if err := json.NewDecoder(body).Decode(&socialiteUser); err != nil {
		return nil, fmt.Errorf("DecodeAuthUser: Error al decodificar JSON - %v", err)
	}
	// Inicializar datos de usuario
	var email, name string
	// Convertir `socialiteUser` al tipo correcto y extraer datos
	switch v := socialiteUser.(type) {
	case *dto.AuthFacebookUser:
		email = v.Email
		name = v.Name
	case *dto.AuthGoogleUser:
		email = v.Email
		name = v.Name
	default:
		return nil, fmt.Errorf("DecodeAuthUser: No se pudo determinar el tipo de usuario")
	}
	// Crear la estructura del usuario
	user := &dao.User{
		Email:      email,
		Username:   email,
		IsVerified: true,
	}
	// Dividir el nombre por espacios
	words := strings.Fields(name)
	// Si hay más de una palabra, asignar el último como apellido
	if len(words) > 1 {
		user.LastName = words[len(words)-1]
		user.FirstName = strings.Join(words[:len(words)-1], " ")
	} else {
		user.LastName = ""
		user.FirstName = name
	}
	// Regresar usuario
	return user, nil
}
