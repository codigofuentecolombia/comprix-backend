package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func InitCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{db}
}

func (repo CategoryRepository) GetAll() ([]dao.Category, error) {
	var categories []dao.Category
	// Obtener todas las categorías de una sola vez
	if err := repo.db.Preload("Subcategories.Subcategories").Where("parent_id IS NULL").Find(&categories).Error; err != nil {
		return nil, fails.Create("CategoryRepository GetAll: No se pudo obtener las categorías.", err)
	}
	// Regresar categorias
	return categories, nil
}

// Función auxiliar para asignar subcategorías desde el mapa
func assignSubcategories(categories []dao.Category, categoryMap map[uint32][]dao.Category) {
	for i := range categories {
		if subcats, found := categoryMap[categories[i].ID]; found {
			categories[i].Subcategories = subcats
			assignSubcategories(categories[i].Subcategories, categoryMap)
		}
	}
}

func HandleNewProductCategories(tx *gorm.DB, categories []string) (*uint32, error) {
	var parentID *uint32 = nil
	// Iterar
	for _, name := range categories {
		category := dao.Category{
			Name:     name,
			ParentID: parentID,
		}
		// Buscar o crear la categoría
		if err := tx.FirstOrCreate(&category, dao.Category{Name: name}).Error; err != nil {
			return nil, fails.Create("HandleNewProductCategories: No se pudo crear la categoria.", nil)
		}
		// Establecer el ID como parent para la siguiente iteración
		parentID = &category.ID
	}
	// Regresar parent id
	return parentID, nil
}
