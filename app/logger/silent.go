package logger

import "gorm.io/gorm/logger"

// NoOpWriter es un escritor que no realiza ninguna operación, útil para silenciar los logs.
type NoOpWriter struct{}

// Write implementa la interfaz de io.Writer, pero no hace nada.
func (w *NoOpWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Printf implementa la interfaz de logger.Writer, pero no hace nada.
func (w *NoOpWriter) Printf(format string, args ...interface{}) {
	// No imprime nada
}

// Función que crea el logger de GORM en modo silencioso.
func SilentLogger() logger.Writer {
	return &NoOpWriter{}
}
