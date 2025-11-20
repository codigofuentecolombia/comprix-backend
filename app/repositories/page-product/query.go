package repository_page_product

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/utils"

	"gorm.io/gorm"
)

func (repo *Repository) GenerateQueryFromFindParams(params dto.ProductRepositoryFindParams) *gorm.DB {
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
		db = db.Scopes(repositories.OlderThanOneDay)
	}
	// Validar si se filtrara por los mejores precios
	if params.BestPrice != nil && *params.BestPrice {
		db = db.Scopes(repo.BestPriceScope)
	} else if params.BestPagePrice != nil && *params.BestPagePrice {
		db = db.Scopes(repo.BestPricePerPageScope)
	} else {
		db = db.Joins("JOIN products AS p on p.id = page_products.product_id")
	}
	// Validar si se filtrara por productos
	if params.WithDistinctProducts != nil && *params.WithDistinctProducts {
		db = db.Scopes(repo.WithDistinctPageProducts)
	}
	// Validar si se obtendrar recomendados
	if params.OnlyRecommended != nil && *params.OnlyRecommended {
		db = db.Where("p.is_recommended = 1")
	}
	// Validar si se obtendrar recomendados
	if params.OnlyWithDiscount != nil && *params.OnlyWithDiscount {
		db = db.Where("p.is_in_discount = 1")
	}
	// Validar si tiene nombre
	if params.ID != nil {
		db = db.Scopes(repositories.HasIdScope(*params.ID, model.TableName()))
	} else if params.ExcludeID != nil && *params.ExcludeID != 0 {
		db = db.Where("page_products.id != ?", params.ExcludeID)
	}
	// Validar si se filtra por producto
	if params.ProductID != nil && *params.ProductID != 0 {
		db = db.Where("page_products.main_product_id = ?", params.ProductID)
	}
	// Verificar url
	if params.Url != nil && *params.Url != "" {
		db = db.Where("page_products.url = ?", *params.Url)
	}
	// Validar si se filtrara por nombre
	if utils.CheckIfStringIsNotEmpty(params.Search) {
		db = db.Scopes(repo.SearchScope(*params.Search))
	}
	// Validar si tiene id de pagina
	if params.PageID != nil {
		db = db.Scopes(repositories.HasPageIDScope(*params.PageID, model.TableName()))
	}
	// Validar si hay marcas
	if params.BranchIDS != nil && len(*params.BranchIDS) > 0 {
		db = db.Where("p.brand_id IN ?", *params.BranchIDS)
	}
	// Validar si tiene categoria
	if utils.CheckIfStringIsNotEmpty(params.CategoryID) {
		db = db.Scopes(repo.InCategoryScope(*params.CategoryID))
	}
	// Validar si tiene tipo
	if params.Type != nil && *params.Type != "" {
		db = db.Scopes(repo.ShowByTypeScope(*params.Type))
	} else {
		// Validar si se obtendrar desactivados
		if params.OnlyDisabled != nil && *params.OnlyDisabled {
			db = db.Where("p.is_disabled = 1")
		} else if params.OmitDisabled != nil && *params.OmitDisabled {
			db = db.Where("p.is_disabled = 0")
		}
	}
	// Verificar si tiene preloads
	if params.Preloads != nil {
		db = db.Scopes(repositories.PreloadScope(*params.Preloads))
	}
	// Ordenar
	if utils.CheckIfStringIsNotEmpty(params.Order) {
		db = db.Order(repo.CustomOrder(*params.Order))
	}
	// Regresar instancia de base de datos
	return db
}
