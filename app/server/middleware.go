package server

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateJwt(config dto.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del encabezado Authorization
		authHeader := c.GetHeader("Authorization")
		// Validar que tenga el header
		if authHeader == "" {
			UnauthorizedException("Token requerido", nil)
		}
		// Separar string por spacio
		tokenParts := strings.Split(authHeader, " ")
		// Verificar que el formato sea correcto "Bearer token"
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			UnauthorizedException("Formato de token inválido", nil)
		}
		// Obtener el valor del jwt
		tokenString := tokenParts[1]
		// Parsear y validar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inválido")
			}
			// Regresar llave
			return []byte(config.Settings.Server.Sk), nil
		})
		// Validar que el token sea correcto
		if err != nil || !token.Valid {
			UnauthorizedException("Token requerido", nil)
		}
		// Extraer los claims del token
		claims, ok := token.Claims.(jwt.MapClaims)
		// Verificar si esta correcto
		if !ok {
			UnauthorizedException("No se pudieron obtener los claims", nil)
		}
		//
		user, err := repositories.InitUserRepository(config.GormDB).FindByID(
			claims["sub"].(map[string]interface{})["id"],
		)
		// Usuario
		if err != nil {
			UnauthorizedException("Usuario no encontrado", nil)
		}
		// Agregar los claims al contexto de la solicitud
		c.Set("claims", claims)
		c.Set("authUser", user)
		c.Next()
	}
}

func Role(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("authUser").(*dao.User)
		// Verificar si tiene role
		if user.Role.Name == "admin" {
			c.Next()
		} else {
			UnauthorizedException("No permitido", nil)
		}
	}
}
