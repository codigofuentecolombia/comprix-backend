package service_product

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"errors"
	"fmt"
)

func (svc *Service) HandleRetrieved(retrieved *dto.RetrievedProduct) error {
	// Verificar si se va a desactivar
	if retrieved.HasStock != nil && !*retrieved.HasStock {
		// Tratar de desactivar
		if err := svc.repositories.pageProduct.Disable(dto.ProductRepositoryFindParams{Url: &retrieved.Url}); err != nil {
			return fails.Create("No se pudo desactivar el producto por url", err)
		}
		// Regresar sin error
		return nil
	}
	// Validar cantidad minima que no venga vacia
	if retrieved.MinQuantityToApplyDiscount == 0 {
		retrieved.MinQuantityToApplyDiscount = 4
	}
	// Verificar que tenga precio
	if err := svc.HandleRetrievedByPrice(retrieved); err != nil {
		return err
	}
	// Verificar si existe por url
	if err := svc.HandleRetrievedByPageUrl(retrieved); err != nil {
		return err
	}
	// Verificar si se encontro
	if retrieved.ProductID == nil {
		// Tratar de encontrar por categoria
		if err := svc.HandleRetrievedByCategories(retrieved); err != nil {
			return err
		}
	}
	// Verificar si se actualizara
	if retrieved.ProductID != nil {
		return svc.repositories.pageProduct.UpdateRetrieved(retrieved)
	}
	// Regresar sin error
	return svc.repositories.pageProduct.CreateRetrieved(retrieved)
}

func (svc *Service) HandleRetrievedByPrice(retrieved *dto.RetrievedProduct) error {
	// Verificar si no tiene precio
	if retrieved.Price == 0 {
		// Verificar que por lo menos tenga precio de descuento
		if retrieved.DiscountPrice > 0 {
			return fails.Create("Validate", errors.New("El producto no tiene precio pero si tiene descuento"))
		} else {
			return fails.Create("Validate", errors.New("El producto no tiene precio ni descuento"))
		}
	}
	// Regresar sin error
	return nil
}

func (svc *Service) HandleRetrievedByPageUrl(retrieved *dto.RetrievedProduct) error {
	order := "id_desc"
	softdeleted := true
	// Obtener producto por url
	response := svc.repositories.pageProduct.Find(dto.ProductRepositoryFindParams{
		Url:         &retrieved.Url,
		Order:       &order,
		Softdeleted: &softdeleted,
		Selects: &dto.RepositoryGormSelections{
			Query: []string{
				"page_products.id",
				"page_products.price",
				"page_products.product_id",
				"page_products.deleted_at",
				"page_products.discount_price",
			},
		},
	})
	// Verificar si hubo error si no se encontro regresar sin error para que continue el otro filtro
	if response.Error != nil {
		return nil
	}
	// Actualizar id de producto
	retrieved.ProductID = &response.Data.ProductID
	// Verificar si el precio ha cambiado
	if err := svc.HandlePriceDifferences(retrieved, response.Data); err != nil {
		return err
	}
	// Regresar sin error
	return nil
}

func (svc *Service) HandleRetrievedByCategories(retrieved *dto.RetrievedProduct) error {
	// Obtener categoria padre
	categoryName := retrieved.Categories[len(retrieved.Categories)-1]
	// Verificar el tamaño total
	if len(retrieved.Categories) >= 3 {
		categoryName = retrieved.Categories[1]
	}
	// Obtener por nombre
	categoryResponse := svc.repositories.category.Find(dto.CategoryRepositoryFindParams{Name: &categoryName})
	// Validar si hubo error
	if categoryResponse.Error != nil {
		return fails.Create("No se pudo obtener la categoria", categoryResponse.Error, true)
	}
	// Obtenerid de categoria
	categoryID := fmt.Sprintf("%d", categoryResponse.Data.ID)
	// Obtener productos
	productResponse := svc.repositories.product.FindAll(dto.ProductRepositoryFindParams{
		Selects:    &dto.RepositoryGormSelections{Query: []string{"id", "name"}},
		CategoryID: &categoryID,
	})
	// Validar si hubo error
	if productResponse.Error != nil {
		return fails.Create("No se pudieron obtener los productos relacionados", productResponse.Error, true)
	}
	// Nombre sanitizado
	productName := svc.SanitizeName(retrieved.Name)
	// Iterar productos
	for _, product := range productResponse.Data {
		if svc.CompareNames(productName, svc.SanitizeName(product.Name)) {
			// Actualizar id de producto
			retrieved.ProductID = &product.ID
			// Romper ciclo
			break
		}
	}
	// Verificar si es null setear ultima categoria para evitar creaciones
	if retrieved.ProductID == nil {
		retrieved.Categories = []string{categoryName}
	}
	// Regresar nil
	return nil
}

func (svc *Service) HandlePriceDifferences(retrieved *dto.RetrievedProduct, old dao.PageProduct) error {
	// Verificar que no este desactivado
	if old.DeletedAt == 0 {
		// Verificar si los precios no han cambiado
		if retrieved.Price == old.Price && retrieved.DiscountPrice == old.DiscountPrice {
			// Tratar de marcarlo como actualizado
			if err := svc.repositories.pageProduct.MarkAsUpdated(old.ID); err != nil {
				return err
			}
			// Regresar mensaje determinado
			return fails.Create("HandleRetrievedPriceDifferences: No existen cambios dentro del producto.", nil)
		}
		// Desactivar producto viejo
		if err := svc.repositories.pageProduct.Disable(dto.ProductRepositoryFindParams{ID: &old.ID}); err != nil {
			return err
		}
	}
	// Regresar sin error
	return nil
}
