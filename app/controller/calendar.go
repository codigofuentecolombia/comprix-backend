package controller

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CalendarController struct {
	config     dto.Config
	repository repositories.CalendarRepository
}

func InitCalendarController(config dto.Config) CalendarController {
	return CalendarController{
		config:     config,
		repository: repositories.InitCalendarRepository(config.GormDB),
	}
}

func (ctr CalendarController) GetAll(ginContext *gin.Context) {
	times, err := ctr.repository.GetAll()
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, times)
}

func (ctr CalendarController) Create(ginContext *gin.Context) {
	form := server.GinFormBinding(ginContext, dto.Calendar{}, "Invalid form")
	calendar, err := ctr.repository.Create(server.StructCopier(dao.Calendar{}, form))
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.CreateEntityResponse(ginContext, calendar)
}

func (ctr CalendarController) Update(ginContext *gin.Context) {
	form := server.GinFormBinding(ginContext, dto.UpdateCalendar{}, "Invalid form")
	// Actualizar
	if err := ctr.repository.Update(form.Times); err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.SuccessResponse(ginContext)
}

func (ctr CalendarController) Delete(ginContext *gin.Context) {
	id := ginContext.Param("id")
	// Verificar si existe
	calendar, err := ctr.repository.FindByID(id)
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Actualizar
	if ctr.repository.Delete(*calendar) != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.SuccessResponse(ginContext)
}
