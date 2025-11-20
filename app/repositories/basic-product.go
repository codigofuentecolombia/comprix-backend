package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

func SetAsRecommended(tx *gorm.DB, id interface{}) error {
	var product dao.Product
	// Validar si existe
	if err := tx.Where("id = ?", id).First(&product).Error; err != nil {
		return fails.Create("SetAsRecommended: No se encontro el producto.", err)
	}
	// Marcar como recomendado
	if err := tx.Model(&dao.Product{}).Where("id = ?", id).Update("is_recommended", 1).Error; err != nil {
		return fails.Create("SetAsRecommended: No se marco.", err)
	}
	// Regresar sin error
	return nil
}

func UnsetAsRecommended(tx *gorm.DB, id interface{}) error {
	var product dao.Product
	// Validar si existe
	if err := tx.Where("id = ?", id).First(&product).Error; err != nil {
		return fails.Create("SetAsRecommended: No se encontro el producto.", err)
	}
	// Marcar como recomendado
	if err := tx.Model(&dao.Product{}).Where("id = ?", id).Update("is_recommended", 0).Error; err != nil {
		return fails.Create("SetAsRecommended: No se marco.", err)
	}
	// Regresar sin error
	return nil
}

func SetAsInDiscount(tx *gorm.DB, id interface{}) error {
	var product dao.Product
	// Validar si existe
	if err := tx.Where("id = ?", id).First(&product).Error; err != nil {
		return fails.Create("SetAsInDiscount: No se encontro el producto.", err)
	}
	// Marcar como recomendado
	if err := tx.Model(&dao.Product{}).Where("id = ?", id).Update("is_in_discount", 1).Error; err != nil {
		return fails.Create("SetAsInDiscount: No se marco.", err)
	}
	// Regresar sin error
	return nil
}

func UnsetAsInDiscount(tx *gorm.DB, id interface{}) error {
	var product dao.Product
	// Validar si existe
	if err := tx.Where("id = ?", id).First(&product).Error; err != nil {
		return fails.Create("UnsetAsInDiscount: No se encontro el producto.", err)
	}
	// Marcar como recomendado
	if err := tx.Model(&dao.Product{}).Where("id = ?", id).Update("is_in_discount", 0).Error; err != nil {
		return fails.Create("UnsetAsInDiscount: No se marco.", err)
	}
	// Regresar sin error
	return nil
}

func HandleExistingProduct(tx *gorm.DB, retrievedProduct *dto.RetrievedProduct, existingProduct dao.PageProduct) error {
	// Verificar si existe algun cambio
	if existingProduct.Price != retrievedProduct.Price || existingProduct.DiscountPrice != retrievedProduct.DiscountPrice {
		// Eliminar producto por pagina
		if err := tx.Where("page_id = ?", existingProduct.PageID).Delete(&existingProduct).Error; err != nil {
			return fails.Create("HandleExistingProduct: No se pudo eliminar el producto existente.", err)
		}
		// Tratar de crear el nuevo producto
		if err := tx.Create(&dao.PageProduct{
			Url:                        retrievedProduct.Url,
			Price:                      retrievedProduct.Price,
			Images:                     retrievedProduct.Images,
			PageID:                     retrievedProduct.PageID,
			ProductID:                  existingProduct.ProductID,
			DiscountPrice:              retrievedProduct.DiscountPrice,
			OriginalPrice:              retrievedProduct.OriginalPrice,
			OriginalDiscountPrice:      retrievedProduct.OriginalDiscountPrice,
			MinQuantityToApplyDiscount: retrievedProduct.MinQuantityToApplyDiscount,
		}).Error; err != nil {
			return fails.Create("HandleExistingProduct: No se pudo crear el producto nuevo.", err)
		}
		// Finalizar funcion
		return nil
	} else {
		err := tx.Model(&dao.PageProduct{}).Where("id = ?", existingProduct.ID).Update("id", gorm.Expr("id")).Error
		// Verificar si hubo error
		if err != nil {
			return fails.Create("HandleExistingProduct: No se pudo marcar el produco como actualizado.", err)
		}
	}
	// Regresar
	return fails.Create("PageProductRepository CreateOrUpdate: No existen cambios dentro del producto.", nil)
}

func HandleNewProduct(tx *gorm.DB, retrievedProduct *dto.RetrievedProduct) error {
	// Crear variables
	var err error
	var brandID *int
	var categoryID *uint32
	var product dao.Product
	// Buscar producto
	if err = tx.Where("sku = ?", retrievedProduct.Sku).First(&product).Error; err != nil {
		// Crear categorias
		categoryID, err = HandleNewProductCategories(tx, retrievedProduct.Categories)
		// Verificar si hay error
		if err != nil {
			return err
		}
		// Crear marca
		brandID, err = HandleNewProductBrand(tx, retrievedProduct.Brand)
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
			return fails.Create("HandleNewProduct: No se pudo crear el producto nuevo.", err)
		}
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
		return fails.Create("HandleNewProduct: No se pudo crear el producto de la pagina nuevo.", err)
	}
	// Regresar exito
	return nil
}

func HandleNewProductBrand(tx *gorm.DB, brandName string) (*int, error) {
	brand := dao.Brand{Name: brandName}
	// Buscar o crear la categoría
	if err := tx.Where(brand).Attrs(brand).FirstOrCreate(&brand, brand).Error; err != nil {
		return nil, fails.Create("HandleNewProductCategories: No se pudo crear la categoria.", nil)
	}
	// Regresar id
	return &brand.ID, nil
}

// func HandleMainProduct(tx *gorm.DB, product dao.Product) (*int, error) {
// 	// Verificar si existe
// 	if err := tx.Where("sku = ?")
// }
