package repository_page_product

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/repositories"

	"gorm.io/gorm"
)

func (repo *Repository) CreateRetrieved(retrievedProduct *dto.RetrievedProduct) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Crear variables
		var err error
		var product dao.Product
		var brandID *int
		var categoryID *uint32
		// Crear categorias
		categoryID, err = repositories.HandleNewProductCategories(tx, retrievedProduct.Categories)
		// Verificar si hay error
		if err != nil {
			return err
		}
		// Crear marca
		brandID, err = repositories.HandleNewProductBrand(tx, retrievedProduct.Brand)
		// Verificar si hay error
		if err != nil {
			return err
		}
		// Crear referencia cuando no se encuentre
		product = dao.Product{
			Sku:         retrievedProduct.Sku,
			Name:        retrievedProduct.Name,
			BrandID:     *brandID,
			CategoryID:  *categoryID,
			Description: retrievedProduct.Description,
		}
		// Crear producto
		if err := tx.Create(&product).Error; err != nil {
			return fails.Create("ProductRepository CreateRetrieved: No se pudo crear el producto nuevo.", err)
		}
		// Tratar de crear
		if err := tx.Create(&dao.PageProduct{
			Url:                        retrievedProduct.Url,
			Price:                      retrievedProduct.Price,
			Images:                     retrievedProduct.Images,
			PageID:                     retrievedProduct.PageID,
			ProductID:                  product.ID,
			MainProductID:              product.ID,
			DiscountPrice:              retrievedProduct.DiscountPrice,
			OriginalPrice:              retrievedProduct.OriginalPrice,
			OriginalDiscountPrice:      retrievedProduct.OriginalDiscountPrice,
			MinQuantityToApplyDiscount: retrievedProduct.MinQuantityToApplyDiscount,
		}).Error; err != nil {
			return fails.Create("ProductRepository CreateRetrieved: No se pudo crear el producto de la pagina nuevo.", err)
		}
		// Regresar exito
		return nil
	})
}
