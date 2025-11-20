package service_email

import (
	"comprix/app/domain/dao"
)

func (s *Service) VerificationCode(user dao.User, vCode dao.VerificationCode) {
	// Contenido del correo en HTML
	body := `
		<html>
			<body style="font-family: Arial, sans-serif; text-align: center; padding: 20px; background-color: #f4f4f4;">
				<div style="max-width: 500px; margin: auto; background: white; padding: 20px; border-radius: 10px;">
					<img src="https://comprix.com.ar/assets/img/logo-dark.png" alt="Logo" style="max-width: 300px; filter: drop-shadow(5px 5px 10px rgba(0, 0, 0, 0.5)); rotate: -10deg;">

					<h2 style="color: #2c3e50;">🔐 ¡Código de Verificación! 🔐</h2>
					<p>Gracias por registrarte con nosotros. Para completar el proceso, por favor ingresa el siguiente código en la aplicación:</p>
					<div style="background-color: #ecf0f1; padding: 15px; border-radius: 5px; text-align: left;">
						<p><strong>🔑 Código de verificación:</strong> <span style="font-size: 24px; font-weight: bold; color: #3498db;">{code}</span></p>
					</div>
					<p style="margin-top: 20px;">Este código expirará en <strong>1 hora</strong>.</p>
					<p style="margin-top: 20px;">Si no realizaste esta solicitud, por favor ignora este mensaje.</p>
					<br>
					<small>⚠️ Este código es válido solo una vez.</small>
				</div>
			</body>
		</html>
    `

	// Reemplazar placeholders con los valores del usuario
	body = replacePlaceholder(body, "{code}", vCode.Code)

	s.Send(user.Email, "Verificación de cuenta", body)
}
