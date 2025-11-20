package utils

import (
	"fmt"
	"time"
)

func FormatDuration(start, end time.Time) string {
	// Calcular la diferencia de tiempo
	duration := end.Sub(start)
	totalSeconds := int(duration.Seconds())
	// Obtener horas, minutos y segundos
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	// Construir la cadena de tiempo
	result := ""
	// Verificar si hay horas
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	// Si hay horas, siempre mostramos minutos
	if minutes > 0 || hours > 0 {
		result += fmt.Sprintf("%dm ", minutes)
	}
	// Concatenar segundos
	result += fmt.Sprintf("%ds", seconds)
	// Regresar resultado
	return result
}
