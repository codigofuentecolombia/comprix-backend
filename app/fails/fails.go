package fails

import "fmt"

type Fail struct {
	Err     error
	Data    interface{}
	Message string
}

func Create(message string, err error, data ...any) error {
	fail := Fail{
		Err:     err,
		Data:    nil,
		Message: message,
	}
	//Validar si existe data
	if len(data) > 0 {
		fail.Data = data[0]
	}
	//Regresar error
	return fail
}

func (err Fail) Error() string {
	//Verificar si el error en nulo
	if err.Err == nil {
		return err.Message
	}
	//Regresar error
	return fmt.Sprintf("%s: %s", err.Message, err.Err.Error())
}
