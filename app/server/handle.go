package server

import "comprix/app/domain/dto"

// Funcion para emitir excepcion con status code 403
func EntityNotFound[T any](response dto.RepositoryGenericResponse[T]) T {
	// Verificar si existe error
	if response.Error != nil {
		InternalErrorException("Could not found entity", nil)
	}
	// Regresar data
	return response.Data
}

func HandleEntitySave[T any](response dto.RepositoryGenericResponse[T]) T {
	// Verificar si existe error
	if response.Error != nil {
		InternalErrorException("Could not save entity", nil)
	}
	// Regresar data
	return response.Data
}
