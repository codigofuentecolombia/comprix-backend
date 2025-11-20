package config

import (
	"comprix/app/domain/dto"
	"comprix/app/fails"
	log "comprix/app/logger"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Load(filepath string) (*dto.Config, error) {
	//Iniciar configuracion
	config := dto.Config{}
	//Abrir archivo
	file, err := os.Open(filepath)
	//Verificar si existe error al abrir el archivo de configuracion
	if err != nil {
		return nil, fails.Create(fails.OpenConfigFileMsg, err)
	}
	//Indicar que cierre el archivo
	defer file.Close()
	//Parsear estructura a los ajustes
	if err := yaml.NewDecoder(file).Decode(&config.Settings); err != nil {
		return nil, fails.Create(fails.DecodeConfigMsg, err)
	}
	//Validar ajustes
	if err := ValidateSettings(&config); err != nil {
		return nil, err
	}
	// Obtener referencia de gorm
	if err := InitDatabase(&config); err != nil {
		return nil, err
	}
	// Regresar configuracion
	return &config, nil
}

func ValidateSettings(conf *dto.Config) error {
	//Validar si existe configuracion requerida
	if conf != nil {
		//Validar Base de datos
		if conf.Settings.Database.Dsn == "" {
			return fails.Create(fails.ConfigRequiredDatabaseDsnMsg, nil)
		}
		// Validar que exista la ruta del log
		if conf.Settings.Paths.Logs == "" {
			return fails.Create(fails.ConfigRequiredLogPathMsg, nil)
		}
		// Validar que exista la ruta de los assets
		if conf.Settings.Paths.Assets == "" {
			return fails.Create(fails.ConfigRequiredAssetPathMsg, nil)
		}
		// Validar puerto
		if conf.Settings.Server.Port == "" {
			conf.Settings.Server.Port = "8080"
		}
		// Validar si tiene goroutines
		if conf.Settings.Scrapping.MaxGoRutines == 0 {
			conf.Settings.Scrapping.MaxGoRutines = 5
		}
	}
	//Regresar configuracion
	return nil
}

func InitDatabase(conf *dto.Config) error {
	gormConfig := &gorm.Config{}
	// Verificar si es debug
	if conf.Settings.Database.Debug {
		gormConfig.Logger = logger.New(
			log.Create(conf.Settings.Paths.Logs, "database"), // Usar logrux como escritor
			logger.Config{
				SlowThreshold: time.Second, // Umbral para consultas lentas
				LogLevel:      logger.Info, // Nivel de logging
				Colorful:      false,       // Sin colores
			},
		)
	} else {
		// Si no es debug, no asignamos logger, evitando la impresión en consola
		gormConfig.Logger = logger.New(
			// No guardamos en archivo ni imprimimos nada en consola
			log.SilentLogger(), // No guardamos en archivo ni mostramos en consola
			logger.Config{
				SlowThreshold: time.Second,   // Umbral para consultas lentas
				LogLevel:      logger.Silent, // Silenciar los logs
				Colorful:      false,         // Sin colores
			},
		)
	}
	// Gestiona conexion
	db, err := gorm.Open(mysql.Open(conf.Settings.Database.Dsn), gormConfig)
	// Verificar si existe error
	if err != nil {
		fails.Create(fails.InitializeGormMsg, err)
	}
	// Asignar base de datos
	conf.GormDB = db
	// Regresar sin error
	return nil
}
