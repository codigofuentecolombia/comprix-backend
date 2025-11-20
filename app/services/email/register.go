package service_email

import (
	"comprix/app/domain/dao"
)

func (s *Service) Register(user dao.User, password string) {
	// Contenido del correo en HTML
	body := `
		<html>
		<body style="font-family: Arial, sans-serif; text-align: center; padding: 20px; background-color: #f4f4f4;">
			<div style="max-width: 500px; margin: auto; background: white; padding: 20px; border-radius: 10px;">
				<!-- Logo arriba -->
				<img src="https://comprix.com.ar/assets/img/logo-dark.png" alt="Logo" style="max-width: 300px; filter: drop-shadow(5px 5px 10px rgba(0, 0, 0, 0.5)); rotate: -10deg;">

				<h2 style="color: #2c3e50;">¡Bienvenido a comprix!</h2>
				<p >Haz tus compras al mejor precio. Te dejamos tu usuario y tu contraseña:</p>
				<div style="background-color: #ecf0f1; padding: 15px; border-radius: 5px; text-align: left;">
					<p><strong>👤 Usuario:</strong> {usuario}</p>
					<p><strong>🔑 Contraseña:</strong> {contraseña}</p>
				</div>
				<p>Puedes iniciar sesión haciendo clic en el siguiente botón:</p>
				<a href="https://comprix.com.ar" style="display: inline-block; padding: 10px 20px; background-color: #3498db; color: white; text-decoration: none; border-radius: 5px;">Iniciar Sesión</a>
				<p>Si tienes algún problema, contáctanos.</p>
			</div>
		</body>
		</html>
    `

	// Reemplazar placeholders con los valores del usuario
	body = replacePlaceholder(body, "{usuario}", user.Username)
	body = replacePlaceholder(body, "{contraseña}", password)

	s.Send(user.Email, "🎉 Bienvenido a comprix", body)
}
