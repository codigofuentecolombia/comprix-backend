package service_email

import (
	"comprix/app/domain/dto"
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

type Service struct {
	config dto.Config
}

func InitService(config dto.Config) Service {
	return Service{config}
}

func (svc *Service) GenerateMailer(email string, subject string) *gomail.Message {
	m := gomail.NewMessage()

	m.SetHeader("From", "info@comprix.com.ar")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	return m
}

func (svc *Service) GenerateDialer() *gomail.Dialer {
	return gomail.NewDialer("vps-4310054-x.dattaweb.com", 465, "info@comprix.com.ar", "Comprar25@")
}

// Función para reemplazar texto en una cadena
func replacePlaceholder(text, placeholder, value string) string {
	return strings.ReplaceAll(text, placeholder, value)
}

func (svc *Service) Send(email string, subject string, body string) {
	m := svc.GenerateMailer(email, subject)
	// Especificar cuerpo
	m.SetBody("text/html", body)
	// Enviar correo
	if err := svc.GenerateDialer().DialAndSend(m); err != nil {
		fmt.Println("Error al enviar:", err)
	}
}
