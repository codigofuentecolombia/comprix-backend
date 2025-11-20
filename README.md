# Comprix Backend

Backend API para el comparador de precios Comprix, desarrollado en Go con Gin framework.

## 🚀 Características

- API RESTful completa
- Autenticación JWT
- OAuth con Facebook y Google
- Sistema de scraping de precios de múltiples tiendas
- Gestión de productos, categorías y marcas
- Sistema de órdenes y carrito
- Panel de administración

## 📋 Requisitos

- Go 1.21+
- MySQL 5.7+

## 🔧 Configuración Local

1. Clonar el repositorio
```bash
git clone <tu-repo>
cd comprix-backend
```

2. Copiar archivo de configuración
```bash
cp settings/conf.example.yaml settings/conf.yaml
```

3. Editar `settings/conf.yaml` con tus credenciales de base de datos

4. Ejecutar el servidor
```bash
go run cmd/api/main.go
```

El servidor estará disponible en `http://localhost:5000`

## 🗄️ Base de Datos

Ejecutar el script SQL en `database_schema.sql` para crear las tablas.

## 📡 Endpoints Principales

- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Registro
- `GET /api/v1/products` - Listar productos
- `GET /api/v1/products/:id` - Detalle de producto
- `GET /api/v1/categories` - Listar categorías
- `GET /api/v1/brands` - Listar marcas
- `POST /api/v1/orders` - Crear orden

## 🕷️ Scraper

Para ejecutar el scraper y obtener productos:

```bash
go run cmd/scrapper/main.go
```

## 🌐 Despliegue

El proyecto está configurado para desplegar en Railway.app usando Docker.

Variables de entorno necesarias:
- `DATABASE_DSN` - String de conexión MySQL

## 📝 Licencia

Privado
