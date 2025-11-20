package repository_product

import (
	"comprix/app/domain/dao"
	"comprix/app/fails"

	"gorm.io/gorm"
)

func (repo *Repository) SyncRelations(mainID uint, ids []uint) error {
	// Guardar en la base de datos dentro de una transacción
	return repo.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Where("product_id in ? OR main_product_id in ?", ids, ids)
		// Actualizar relaciones
		if err := query.Updates(dao.PageProduct{ProductID: mainID, MainProductID: mainID}).Error; err != nil {
			return fails.Create("ProductRepository SyncRelations: No se pudo crear agrupacion", err)
		}
		// Regresar sin error
		return nil
	})
}
