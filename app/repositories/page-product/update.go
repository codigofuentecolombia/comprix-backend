package repository_page_product

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/repositories"

	"gorm.io/gorm"
)

func (repo *Repository) MarkAsUpdated(id any) error {
	return repositories.MarkAsUpdated("ProductRepository MarkAsUpdated", dao.PageProduct{}, repo.db.Where("id = ?", id))
}

func (repo *Repository) UpdateRetrieved(retrieved *dto.RetrievedProduct) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Create(&dao.PageProduct{
			Url:                        retrieved.Url,
			Price:                      retrieved.Price,
			Images:                     retrieved.Images,
			PageID:                     retrieved.PageID,
			ProductID:                  *retrieved.ProductID,
			MainProductID:              *retrieved.ProductID,
			DiscountPrice:              retrieved.DiscountPrice,
			OriginalPrice:              retrieved.OriginalPrice,
			OriginalDiscountPrice:      retrieved.OriginalDiscountPrice,
			MinQuantityToApplyDiscount: retrieved.MinQuantityToApplyDiscount,
		})
		// Tratar de crear el nuevo producto
		if err := query.Error; err != nil {
			return fails.Create("ProductRepository UpdateRetrieved: No se pudo actualizar el producto de pagina.", err)
		}
		// Finalizar funcion
		return nil
	})
}
