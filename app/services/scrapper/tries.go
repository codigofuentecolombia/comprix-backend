package service_scrapper

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"fmt"

	"github.com/sirupsen/logrus"
)

func HandleTotalTries[ResponseType any](svc *Service, args dto.HandleScrapperTotalTriesParams[ResponseType]) ResponseType {
	// Declarar campos para el loger
	fields := logrus.Fields{
		"url":   args.CallbackArgs.Url,
		"error": args.Err,
		"tries": args.CallbackArgs.Tries,
	}
	// Verificar si existe url
	if args.Url != "" {
		fields["url"] = args.Url
	}
	// Verificar la cantidad de intentos
	if args.CallbackArgs.Tries < svc.Conf.Settings.Scrapping.TotalTries {
		args.CallbackArgs.Tries += 1
		// Registrar en el log
		svc.CreateWarningLog(fmt.Sprintf("%s, intento #%d - volviendo a intentar.", args.Msg, args.CallbackArgs.Tries), fields)
		// Reintentar
		return args.Callback(args.CallbackArgs)
	}
	// Mostrar log
	svc.CreateWarningLog(fmt.Sprintf("%s, se ha excedido el maximo de intentos consecutivos", args.Msg), fields)
	// Crear error
	svc.Repositories.Error.Create(
		dao.Error{
			Url:        args.CallbackArgs.Url,
			Error:      args.Err.Error(),
			PageID:     svc.Page.ID,
			Message:    args.Msg,
			Categories: args.CallbackArgs.Categories,
		},
	)
	// Finalizar funcion
	return args.DefaultResponse
}
