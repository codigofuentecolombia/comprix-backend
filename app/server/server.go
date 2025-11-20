package server

import (
	"comprix/app/domain/dto"
	"comprix/app/logger"
	"comprix/app/validators"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func Initialize(config dto.Config) *gin.Engine {
	// Verificar si el servidor esta en modo de produccion
	if config.Settings.Server.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}
	// Inicializar servidor de gin
	httpHandler := gin.New()
	// Agregar validadores
	validators.RegisterValidations()
	// Verificar si se deben debugear
	if config.Settings.Server.Debug {
		// Inicializar loger
		httpHandler.Use(gin.LoggerWithWriter(logger.Create(config.Settings.Paths.Logs, "server").Out))
	}
	// Hacer que no se finalice en errores
	httpHandler.Use(gin.Recovery())
	// Manejo de excepciones
	httpHandler.Use(HandleException(config))
	// Cors
	httpHandler.Use(CORSMiddleware())
	// Regresar instancia de servidor
	return httpHandler
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Start(cnf dto.Config, httpHandler *gin.Engine) {
	server := &http.Server{
		Addr:    cnf.Settings.Server.GetAddress(),
		Handler: httpHandler,
	}
	// Crear salida
	exit := make(chan os.Signal, 1)
	// Notificar
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	// Escuchar gin server
	Listen(server, exit)
}

func Listen(server *http.Server, exit chan os.Signal) {
	// Inicializar listener
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		// Verificar que no haya error y si hay error que sea distinto al cierre de servidor
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen gin server: %s\n", err)
		}
	}()
	log.Printf("Server is ready to handle requests at %s", server.Addr)
	// Esperar para salir
	<-exit
	// Apagar servidor
	Shutdown(server)
}

func Shutdown(server *http.Server) {
	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cancelar context al finalizar la funcion
	defer cancel()
	// Verificar que no haya error al apagar el server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown gin server: %s\n", err)
	}
	// Marcar contexto como finalizado
	<-ctx.Done()
}
