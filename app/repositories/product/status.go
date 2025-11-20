package repository_product

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

func (repo *Repository) UpdateStatus(id interface{}, column dto.ProductStatusColumn, status bool) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var flag int
		var product dao.Product
		// Verificar si se activara
		if status {
			flag = 1
		} else {
			flag = 0
		}
		// Validar si existe
		if err := tx.Select("id").Where("id = ?", id).First(&product).Error; err != nil {
			return fails.Create("UpdateStatus: No se encontro el producto.", err)
		}
		// Marcar como recomendado
		if err := tx.Model(&dao.Product{}).Where("id = ?", id).Update(string(column), flag).Error; err != nil {
			return fails.Create("UpdateStatus: No se marco.", err)
		}
		// Regresar sin error
		return nil
	})
}
