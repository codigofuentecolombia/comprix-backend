package validators

import (
	"comprix/app/domain/dto"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func TimeFormatValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("15:04:05", fl.Field().String()) // Verifica el formato HH:MM:SS
	return err == nil
}

func RegisterValidations() {
	// Obtener la instancia del validador Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Registrar validador de formato de hora
		v.RegisterValidation("timeFormat", TimeFormatValidator)
		v.RegisterStructValidation(CalendarStructValidator, dto.Calendar{})
	}
}

func CalendarStructValidator(sl validator.StructLevel) {
	calendar := sl.Current().Interface().(dto.Calendar) // Convertir el struct

	// Validar formato de tiempo
	if _, err := time.Parse("15:04:05", calendar.StartTime); err != nil {
		sl.ReportError(calendar.StartTime, "StartTime", "StartTime", "timeFormat", "")
	}
	if _, err := time.Parse("15:04:05", calendar.EndTime); err != nil {
		sl.ReportError(calendar.EndTime, "EndTime", "EndTime", "timeFormat", "")
	}

	// Validar que StartTime < EndTime
	start, _ := time.Parse("15:04:05", calendar.StartTime)
	end, _ := time.Parse("15:04:05", calendar.EndTime)

	if !start.Before(end) {
		sl.ReportError(calendar.EndTime, "EndTime", "EndTime", "startBeforeEnd", "")
	}
}
