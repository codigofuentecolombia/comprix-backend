package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func new(filename, basePath, logType string, optIsDebug ...bool) *logrus.Logger {
	// Construye la ruta según los parámetros
	var logPath string = filepath.Join(basePath, logType, fmt.Sprintf("%s.log", filename))
	// Crea los directorios si no existen
	if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
		logrus.Errorf("Error creando directorios: %v", err)
	}
	// Configura `lumberjack` para manejar la rotación
	rotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,   // Tamaño máximo del archivo en MB
		MaxAge:     0,     // No eliminar por antigüedad
		MaxBackups: 0,     // Sin límite de copias de respaldo
		Compress:   false, // No comprimir archivos antiguos
	}
	// Configura el logger de `logrus`
	logger := logrus.New()
	// Definir comportamientos
	logger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
	logger.SetLevel(logrus.DebugLevel)
	// Verificar si es debug
	if len(optIsDebug) > 0 && optIsDebug[0] {
		logger.SetOutput(io.MultiWriter(rotator, os.Stdout))
	} else {
		logger.SetOutput(rotator)
	}
	// Regresar logger
	return logger
}

// Crear logger general
func Create(basePath, logType string, optIsDebug ...bool) *logrus.Logger {
	return new("general", basePath, logType, optIsDebug...)
}

// Crear logger de debug
func CreateDebug(basePath, logType string, optIsDebug ...bool) *logrus.Logger {
	return new("debug", basePath, logType, optIsDebug...)
}

// Crear logger de alertas
func CreateWarning(basePath, logType string, optIsDebug ...bool) *logrus.Logger {
	return new("warning", basePath, logType, optIsDebug...)
}
