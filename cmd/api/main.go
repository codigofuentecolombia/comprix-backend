package main

import (
	"comprix/app/config"
	"comprix/app/domain/dto"
	"comprix/app/routes"
	"comprix/app/scrapper/pages"
	pages_alem "comprix/app/scrapper/pages/alem"
	pages_carrefour "comprix/app/scrapper/pages/carrefour"
	pages_hiperlibertad "comprix/app/scrapper/pages/hiperlibertad"
	pages_jumbo "comprix/app/scrapper/pages/jumbo"
	pages_masonline "comprix/app/scrapper/pages/masonline"
	pages_vea "comprix/app/scrapper/pages/vea"
	"comprix/app/server"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
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
	log.Println("🚀 Comprix Backend v1.0.1 - Iniciando servidor...")
	ginEngine := server.Initialize(*cnf)
	// Manejar rutas
	routes.RouterHandler(*cnf, ginEngine)
	
	// Iniciar scraper automáticamente en background
	log.Println("⚙️  Scraper se iniciará en background después de que el servidor esté listo")
	go startScraper(*cnf)
	
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

func startScraper(config dto.Config) {
	log.Println("🔄 Iniciando scraper automático...")
	
	wg := sync.WaitGroup{}
	
	// Primer grupo de scrapers
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Println("Scrapeando Alem...")
		analizeAlemCategories(&config)
		log.Println("Scrapeando Carrefour...")
		pages.AnalizePageProductsByCategories(&config, pages_carrefour.Initialize, 1, true)
		log.Println("Scrapeando Jumbo...")
		pages.AnalizePageProductsByCategories(&config, pages_jumbo.Initialize, 1, true)
	}()
	
	go func() {
		defer wg.Done()
		log.Println("Scrapeando Vea...")
		pages.AnalizePageProductsByCategories(&config, pages_vea.Initialize, 1, true)
		log.Println("Scrapeando Hiperlibertad...")
		pages.AnalizePageProductsByCategories(&config, pages_hiperlibertad.Initialize, 1, true)
		log.Println("Scrapeando MasOnline...")
		pages.AnalizePageProductsByCategories(&config, pages_masonline.Initialize, 1, true)
	}()
	
	wg.Wait()
	
	log.Println("✅ Scraper completado exitosamente")
}

func analizeAlemCategories(config *dto.Config) {
	svc, err := pages_alem.Initialize(config)
	if err == nil {
		svc.GetCategoryProducts()
	}
}
