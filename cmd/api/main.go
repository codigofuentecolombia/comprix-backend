package main

import (
	"comprix/app/config"
	"comprix/app/domain/dto"
	"comprix/app/routes"
	"comprix/app/server"
	"flag"
	"log"
)

var cnf *dto.Config
var configPath string

func init() {
	// Define flag for configuration path
	flag.StringVar(&configPath, "cnf", "settings/conf.yaml", "Path to the configuration file")
	flag.Parse()
	// Obtener configuracion
	cnf = HandleError(config.Load(configPath))
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
