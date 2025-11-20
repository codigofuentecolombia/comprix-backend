package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"time"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func InitProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) CreateOrUpdate(scrappedProduct *dto.RetrievedProduct) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Verificar si existe un registro
		var pageProduct dao.PageProduct
		// Buscar
		err := tx.Where("product_id = (SELECT id FROM products WHERE sku = ?)", scrappedProduct.Sku).
			Where("page_id = ?", scrappedProduct.PageID).
			First(&pageProduct).Error
		// Verificar si se encontro
		if err == nil {
			return HandleExistingProduct(tx, scrappedProduct, pageProduct)
		} else {
			return HandleNewProduct(tx, scrappedProduct)
		}
	})
}

func (repo *ProductRepository) Disable(scrappedProduct *dto.RetrievedProduct) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("url = ?", scrappedProduct.Url).Delete(&dao.PageProduct{}).Error; err != nil {
			return fails.Create("ProductRepository Disable: No se pudo desactivar el producto.", err)
		}
		// Regresar sin error
		return nil
	})
}

func (repo *ProductRepository) DisableNotFoundInPage(beforeDate time.Time, pageID any) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("updated_at < ?", beforeDate).Where("page_id = ?", pageID).Delete(&dao.PageProduct{}).Error; err != nil {
			return fails.Create("ProductRepository DisableNotFoundInPage: No se pudo eliminar los productos no encontrados.", err)
		}
		// Regresar sin error
		return nil
	})
}

func (repo *ProductRepository) Update(form dto.UpdatePageProduct) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		var pageProduct dao.PageProduct
		// Verificar si existe un registro
		if err := repo.db.Where("id = ?", form.ID).Preload("Product").First(&pageProduct).Error; err != nil {
			return fails.Create("ProductRepository Update: no se pudo obtener el detalle del producto", err)
		}
		// Actualizar producto de pagina
		err := repo.db.Model(&dao.PageProduct{}).Where("id = ?", form.ID).Updates(map[string]interface{}{
			"price":                          form.Price,
			"discount_price":                 form.DiscountPrice,
			"min_quantity_to_apply_discount": form.MinQuantityToApplyDiscount,
		}).Error
		// Validar que se pueda actualizar
		if err != nil {
			return fails.Create("ProductRepository Update: no se pudo actualizar el precio del producto", err)
		}
		// Actualizar datos del producto padre
		err = repo.db.Model(&dao.Product{}).Where("id = ?", pageProduct.Product.ID).Updates(map[string]interface{}{
			"name":        form.Name,
			"description": form.Description,
		}).Error
		// Validar si se pudo actualizar
		if err != nil {
			return fails.Create("ProductRepository Update: no se pudo actualizar el detalle del producto", err)
		}
		// Regresar sin errores
		return nil
	})
}

func (repo *ProductRepository) Find(params dto.ProductRepositoryFindParams) dto.RepositoryGenericResponse[dao.PageProduct] {
	return Find("ProductRepository Find", dao.PageProduct{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *ProductRepository) FindAll(params dto.ProductRepositoryFindParams) dto.RepositoryGenericResponse[[]dao.PageProduct] {
	return Find("ProductRepository FindAll", []dao.PageProduct{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *ProductRepository) FindNotFoundInPage(beforeDate time.Time, pageID any) dto.RepositoryGenericResponse[[]dao.PageProduct] {
	// Obtener data
	query := repo.db.Model(&dao.PageProduct{}).Where("updated_at < ?", beforeDate).Where("page_id = ?", pageID)
	// Buscar productos
	return Find("ProductRepository GetNotFoundInPage", []dao.PageProduct{}, query)
}

func (repo *ProductRepository) GenerateQueryFromFindParams(params dto.ProductRepositoryFindParams) *gorm.DB {
	model := dao.PageProduct{}
	//
	db := repo.db.Model(&model)
	// Verificar si se  incluira softdeleted
	if params.Softdeleted != nil && *params.Softdeleted {
		db = db.Unscoped()
	}
	// Ver si tiene campos seleccionados
	if params.Selects != nil {
		db = db.Select(params.Selects.Query)
	}
	// Validar si sera productos que tengan mas de un dia
	if params.OlderThanOneDay != nil && *params.OlderThanOneDay {
		db = db.Scopes(OlderThanOneDay)
	}
	// Validar si tiene nombre
	if params.ID != nil {
		db = db.Scopes(HasIdScope(*params.ID, model.TableName()))
	}
	// Validar si tiene id de pagina
	if params.PageID != nil {
		db = db.Scopes(HasPageIDScope(*params.PageID, model.TableName()))
	}
	// Verificar si tiene preloads
	if params.Preloads != nil {
		db = db.Scopes(PreloadScope(*params.Preloads))
	}
	// Ordenar por
	if params.OrderBy != nil {
		db = db.Order(params.OrderBy)
	}
	// Regresar instancia de base de datos
	return db
}

func (repo *ProductRepository) GetOtherStores(id any, pageID any) dto.RepositoryGenericResponse[[]dao.PageProduct] {
	// Subconsulta para obtener el precio mínimo por grupo y página
	subquery := repo.db.Table("page_products pp1").
		Select("pgm1.group_id, pp1.page_id, MIN(pp1.price) AS min_price").
		Joins("JOIN product_group_members pgm1 ON pp1.product_id = pgm1.product_id").
		Where("pp1.deleted_at = 0 AND pp1.price > 0").
		Group("pgm1.group_id, pp1.page_id")
	// Subconsulta para obtener el group_id relacionado con el producto con id 28571
	groupSubquery := repo.db.Table("product_group_members pgm1").
		Joins("JOIN page_products pp1 ON pp1.product_id = pgm1.product_id").
		Where("pp1.id = ?", id).
		Select("pgm1.group_id")
	// Consulta principal
	return Find("ProductRepository GetOtherStores", []dao.PageProduct{}, repo.db.Table("page_products pp").
		Select("pp.*").
		Joins("JOIN product_group_members pgm ON pp.product_id = pgm.product_id").
		Joins("JOIN (?) min_prices ON min_prices.group_id = pgm.group_id AND min_prices.page_id = pp.page_id", subquery).
		Joins("JOIN (?) group_id_subquery ON group_id_subquery.group_id = pgm.group_id", groupSubquery).
		Where("pp.price = min_prices.min_price").
		Where("pgm.group_id = group_id_subquery.group_id").
		Where("pp.id != ?", id).
		Where("pp.page_id != ?", pageID).
		Where("pp.deleted_at = 0").
		Preload("Page"))
}
