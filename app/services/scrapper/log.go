package service_scrapper

import "github.com/sirupsen/logrus"

func (service *Service) CreateDebugLog(msg string, params logrus.Fields) {
	// Verificar si existen parametros
	if params != nil {
		service.Logs.Debug.WithFields(params).Debug(msg)
	} else {
		service.Logs.Debug.Debug(msg)
	}
}

func (service *Service) CreateWarningLog(msg string, params logrus.Fields) {
	// Verificar si existen parametros
	if params != nil {
		service.Logs.Warning.WithFields(params).Error(msg)
	} else {
		service.Logs.Warning.Error(msg)
	}
}
