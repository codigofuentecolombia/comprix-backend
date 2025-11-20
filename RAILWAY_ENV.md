# Variables de Entorno para Railway

Configura estas variables en el dashboard de Railway:

## Variables Requeridas

```bash
# Base de datos
DATABASE_DSN=user:password@tcp(host:3306)/database?charset=utf8&parseTime=True&loc=Local

# Servidor
SERVER_SECRET_KEY=your-secret-key-here
PORT=5000

# Email
EMAIL_HOST=smtp.example.com
EMAIL_PORT=465
EMAIL_USERNAME=your-email@example.com
EMAIL_PASSWORD=your-email-password

# OAuth Facebook
FACEBOOK_CLIENT_ID=your-facebook-client-id
FACEBOOK_CLIENT_SECRET=your-facebook-client-secret
FACEBOOK_CALLBACK_URL=https://your-domain.com/auth/facebook

# OAuth Google
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_CALLBACK_URL=https://your-domain.com/auth/google
```

## Variables Opcionales

```bash
DATABASE_DEBUG=false
```

## Cómo configurar en Railway

1. Ve al dashboard de tu proyecto en Railway
2. Click en tu servicio
3. Ve a la pestaña "Variables"
4. Click en "New Variable"
5. Agrega cada variable una por una
6. Railway redesplegará automáticamente
