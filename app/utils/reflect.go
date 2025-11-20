package utils

import "reflect"

func ValueIsArrayOrSlice(value any) bool {
	kind := reflect.ValueOf(value).Kind()
	// Verificar si es array o slice
	return kind == reflect.Array || kind == reflect.Slice
}
