package repository_page_product

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

func (repo *Repository) GroupExisting(params dto.GroupPageProduct) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var mainProduct dao.PageProduct
		var secondaryProduct dao.PageProduct
		// Obtener producto principal
		if err := tx.Where("id = ?", params.ID).First(&mainProduct).Error; err != nil {
			return fails.Create("ProductRepository GroupExisting: No se pudo obtener el producto primario", err)
		}
		// Obtener producto secundario
		if err := tx.Where("id = ?", params.NewID).First(&secondaryProduct).Error; err != nil {
			return fails.Create("ProductRepository GroupExisting: No se pudo obtener el producto secundario", err)
		}
		// Verificar si el producto secundario no es de jumbo
		if secondaryProduct.PageID == 3 {
			// Actualizar datos del producto padre
			err := repo.db.Model(&dao.PageProduct{}).Where("id = ?", mainProduct.ID).Updates(map[string]interface{}{
				"product_id":      secondaryProduct.ProductID,
				"main_product_id": secondaryProduct.MainProductID,
			}).Error
			// Verificar si ocurrio un error
			if err != nil {
				return fails.Create("ProductRepository GroupExisting: No se pudo agrupar los productos", err)
			}
		} else {
			// Actualizar datos del producto padre
			err := repo.db.Model(&dao.PageProduct{}).Where("id = ?", secondaryProduct.ID).Updates(map[string]interface{}{
				"product_id":      mainProduct.ProductID,
				"main_product_id": mainProduct.MainProductID,
			}).Error
			// Verificar si ocurrio un error
			if err != nil {
				return fails.Create("ProductRepository GroupExisting: No se pudo agrupar los productos", err)
			}
		}
		// Regresar sin error
		return nil
	})
}
