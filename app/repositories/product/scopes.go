package repository_product

import "gorm.io/gorm"

func (repo *Repository) InCategoryScope(categoryID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var categoryIDs []uint32
		// Generar consulta
		err := repo.db.Raw(`
			WITH RECURSIVE category_hierarchy AS (
				SELECT id FROM categories WHERE id = ?
				UNION ALL
				SELECT c.id
				FROM categories c
				INNER JOIN category_hierarchy ch ON c.parent_id = ch.id
			)
			SELECT id FROM category_hierarchy
		`, categoryID).Scan(&categoryIDs).Error
		// Verificar si hubo un error
		if err != nil {
			return db
		}
		// Regresar ids
		return db.Where("category_id IN ?", categoryIDs)
	}
}
