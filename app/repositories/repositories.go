package repositories

import (
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func Exists(context string, db *gorm.DB) dto.RepositoryGenericResponse[bool] {
	var exists bool
	// Contar registros
	err := db.Select("1").Limit(1).Find(&exists).Error
	// Verificar si hubo error
	if err != nil {
		return dto.RepositoryGenericResponse[bool]{
			Data:  false,
			Error: fails.Create(fmt.Sprintf("%s: No se pudo realizar la consulta", context), err, nil),
		}
	}
	// Regresar exito
	return dto.RepositoryGenericResponse[bool]{
		Data:  exists,
		Error: nil,
	}
}

func Disable[Type any](context string, model Type, db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Tratar de desactivar
		if err := tx.Delete(&model).Error; err != nil {
			return fails.Create("%s: No se pudo desactivar", err)
		}
		// Regresar sin error
		return nil
	})
}

func MarkAsUpdated[Type any](context string, model Type, db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Tratar de desactivar
		if err := tx.Model(&model).UpdateColumn("updated_at", time.Now()).Error; err != nil {
			return fails.Create("%s: No se pudo marcar como actualizado", err)
		}
		// Regresar sin error
		return nil
	})
}

func Find[Type any](context string, model Type, db *gorm.DB) dto.RepositoryGenericResponse[Type] {
	var err error
	var response Type
	// Verificar si es array
	if utils.ValueIsArrayOrSlice(model) {
		err = db.Find(&response).Error
	} else {
		err = db.Take(&response).Error
	}
	// Verificar si hubo error
	if err != nil {
		return dto.RepositoryGenericResponse[Type]{
			Data:  model,
			Error: fails.Create(fmt.Sprintf("%s: No se pudo realizar la consulta", context), err, nil),
		}
	}
	// Regresar exito
	return dto.RepositoryGenericResponse[Type]{
		Data:  response,
		Error: nil,
	}
}

func GetPaginated[Type any](pagination dto.Pagination[Type], db *gorm.DB) dto.RepositoryGenericResponse[dto.Pagination[Type]] {
	var items []Type
	var totalItems int64
	//
	countQuery := db.Session(&gorm.Session{})
	// Validar pagina
	pagination.Validate()
	// Contar el total de resultados
	if err := countQuery.Select("COUNT(*)").Scan(&totalItems).Error; err != nil {
		return dto.RepositoryGenericResponse[dto.Pagination[Type]]{
			Data:  pagination,
			Error: fails.Create("GetPaginated Count: No se pudieron contar los resultados.", err),
		}
	}
	// Obtener resultados con paginación
	if err := db.Limit(pagination.Limit).Offset(pagination.Offset).Find(&items).Error; err != nil {
		return dto.RepositoryGenericResponse[dto.Pagination[Type]]{
			Data:  pagination,
			Error: fails.Create("GetPaginated Find: No se pudieron obtener los resultados.", err),
		}
	}
	// Calcular el total de páginas
	totalPages := int((totalItems + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	// Actualizar pagina
	pagination.Items = items
	pagination.TotalItems = int(totalItems)
	pagination.TotalPages = totalPages
	// Crear y devolver la estructura de paginación
	return dto.RepositoryGenericResponse[dto.Pagination[Type]]{
		Data:  pagination,
		Error: nil,
	}
}

func getTableColumn(column string, table ...string) string {
	if len(table) > 0 {
		return fmt.Sprintf("%s.%s", table[0], column)
	}
	// Regresar columna
	return column
}
