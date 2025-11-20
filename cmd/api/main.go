package main

import (
	"comprix/app/config"
	"comprix/app/domain/dto"
	"comprix/app/routes"
	"comprix/app/server"
	"flag"
	"fmt"
	"log"
	"os"
)

var cnf *dto.Config
var configPath string

func init() {
	// Define flag for configuration path
	flag.StringVar(&configPath, "cnf", "settings/conf.yaml", "Path to the configuration file")
	flag.Parse()
	
	// Check if config file exists, if not create from env vars
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Config file not found, creating from environment variables...")
		if err := createConfigFromEnv(configPath); err != nil {
			log.Printf("Warning: Could not create config from env: %v", err)
		}
	}
	
	// Obtener configuracion
	cnf = HandleError(config.Load(configPath))
}

func createConfigFromEnv(path string) error {
	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		return fmt.Errorf("DATABASE_DSN environment variable is required")
	}
	
	// Create settings directory if it doesn't exist
	if err := os.MkdirAll("settings", 0755); err != nil {
		return err
	}
	
	configContent := fmt.Sprintf(`paths:
  logs: "./logs"
  assets: "./assets"
database:
  dsn: "%s"
  debug: %s
scrapping:
  total_tries: 3
  max_go_rutines: 6
auth:
  facebook:
    client: "%s"
    secret: "%s"
    callback: "%s"
  google:
    client: "%s"
    secret: "%s"
    callback: "%s"
storage:
  path: "./storage"
server:
  sk: "%s"
  port: %s
  host: "0.0.0.0"
  debug: false
email:
  host: "%s"
  port: %s
  username: "%s"
  password: "%s"
`,
		databaseDSN,
		getEnvOrDefault("DATABASE_DEBUG", "true"),
		os.Getenv("FACEBOOK_CLIENT_ID"),
		os.Getenv("FACEBOOK_CLIENT_SECRET"),
		os.Getenv("FACEBOOK_CALLBACK_URL"),
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_CALLBACK_URL"),
		getEnvOrDefault("SERVER_SECRET_KEY", "default-secret"),
		getEnvOrDefault("PORT", "5000"),
		os.Getenv("EMAIL_HOST"),
		getEnvOrDefault("EMAIL_PORT", "465"),
		os.Getenv("EMAIL_USERNAME"),
		os.Getenv("EMAIL_PASSWORD"),
	)
	
	return os.WriteFile(path, []byte(configContent), 0644)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	ginEngine := server.Initialize(*cnf)
	// Manejar rutas
	routes.RouterHandler(*cnf, ginEngine)
	// Encender server
	server.Start(*cnf, ginEngine)
}

func HandleError[T any](data T, err error) T {
	// Verificar si existe error
	if err != nil {
		log.Fatal(err)
	}
	// Regresar data
	return data
}
