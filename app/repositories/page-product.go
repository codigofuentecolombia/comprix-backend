package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type PageProductRepository struct {
	db *gorm.DB
}

func InitPageProductRepository(db *gorm.DB) PageProductRepository {
	return PageProductRepository{db}
}

func (repo PageProductRepository) CreateOrUpdate(retrievedProduct *dto.RetrievedProduct) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Verificar si existe un registro
		var pageProduct dao.PageProduct
		// Buscar
		err := tx.Where("product_id = (SELECT id FROM products WHERE sku = ?)", retrievedProduct.Sku).
			Where("page_id = ?", retrievedProduct.PageID).
			First(&pageProduct).Error
		// Verificar si se encontro
		if err == nil {
			return HandleExistingProduct(tx, retrievedProduct, pageProduct)
		} else {
			return HandleNewProduct(tx, retrievedProduct)
		}
	})
}

func (repo PageProductRepository) GetAllQuery(params dto.GetProductsParams) *gorm.DB {
	query := repo.db.Table("(?) AS pp1", repo.MinPriceQuery()).
		Select("pp1.*").
		Joins("JOIN products AS p on p.id = pp1.product_id")
	// Verificar si se buscara por texto
	if params.CategoryID != "" {
		var categoryIDs []uint32
		// Obtener categorias
		categoryIDs, err := repo.GetChildCategoryIDs(params.CategoryID)
		if err != nil {
			// Manejar error
			return nil
		}
		// Adjuntar query
		query = query.Where("pp1.product_id IN (SELECT id FROM products WHERE category_id IN ?)", categoryIDs)
	}
	// Verificar si tiene sucursales
	if params.BranchIds != nil && len(*params.BranchIds) > 0 {
		query = query.Where("p.brand_id IN (?)", *params.BranchIds)
	}
	// Verificar tipo
	switch params.Type {
	case dto.DiscountsProductsType:
		query = query.Where("pp1.discount_price > 0")
		break
	case dto.RecommendedProductsType:
		query = query.Where("p.is_recommended = 1")
		break
	case dto.OffersProductsType:
		query = query.Where("p.is_in_discount = 1")
		break
	}
	// Verificar si existe filtro
	if params.Pagination.Search != "" {
		query = query.Where(
			`pp1.product_id IN (
				SELECT id FROM products 
				WHERE name LIKE ? OR sku LIKE ?
			)`,
			params.Pagination.Search+"%",
			"%"+params.Pagination.Search+"%",
		)
	}
	// Verificar si tiene orden
	if params.OrderBy != nil && *params.OrderBy != "" {
		query = query.Order("pp1.price ASC")
	} else {
		query = query.Order("pp1.price DESC")
	}
	// Regresar
	return query.
		Preload("Product").
		Preload("Page").
		Preload("Product.Category").
		Preload("Product.Brand")
}

func (repo PageProductRepository) GetChildCategoryIDs(categoryID string) ([]uint32, error) {
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
		return nil, err
	}
	// Regresar ids
	return categoryIDs, nil
}

func (repo PageProductRepository) GetAll(params dto.GetProductsParams) ([]dao.PageProduct, error) {
	var products []dao.PageProduct
	// Verificar si hay error
	if err := repo.GetAllQuery(params).Find(&products).Error; err != nil {
		return nil, fails.Create("PageProductRepository GetAll: No se pudieron obtener los productos.", err)
	}
	// Regresar data
	return products, nil
}

func (repo PageProductRepository) GetPaginated(params dto.GetProductsParams) (dto.Pagination[dao.PageProduct], error) {
	var products []dao.PageProduct
	var totalItems int64
	// Contar el total de productos
	if err := repo.GetAllQuery(params).Select("count(*)").Model(&dao.PageProduct{}).Count(&totalItems).Error; err != nil {
		return dto.Pagination[dao.PageProduct]{}, fails.Create("PageProductRepository Paginated Count: No se pudieron contar los productos.", err)
	}
	// Obtener productos con paginación
	if err := repo.GetAllQuery(params).
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Find(&products).Error; err != nil {
		return dto.Pagination[dao.PageProduct]{}, fails.Create("PageProductRepository Paginated: No se pudieron obtener los productos.", err)
	}
	// Calcular el total de páginas
	totalPages := int((totalItems + int64(params.Pagination.Limit) - 1) / int64(params.Pagination.Limit)) // Redondeo hacia arriba
	// Crear y devolver la estructura de paginación
	return dto.Pagination[dao.PageProduct]{
		Items:      products,
		Index:      params.Pagination.Index,
		Limit:      params.Pagination.Limit,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	}, nil
}

func (repo PageProductRepository) GetByID(id interface{}) (dao.PageProduct, error) {
	var product dao.PageProduct
	// Verificar si hay error
	if err := repo.GetAllQuery(dto.GetProductsParams{}).Where("pp1.id = ?", id).First(&product).Error; err != nil {
		return product, fails.Create("PageProductRepository GetByID: No se pudo obtener el producto.", err)
	}
	// Regresar data
	return product, nil
}

func (repo PageProductRepository) GetAllWithLimit(limit int) ([]dao.PageProduct, error) {
	var products []dao.PageProduct
	// Verificar si hay error
	if err := repo.GetAllQuery(dto.GetProductsParams{}).Limit(limit).Find(&products).Error; err != nil {
		return nil, fails.Create("PageProductRepository GetAllWithLimit: No se pudieron obtener los productos.", err)
	}
	// Regresar data
	return products, nil
}

func (repo PageProductRepository) GetWithDiscount() ([]dao.PageProduct, error) {
	var products []dao.PageProduct
	// Query principal
	err := repo.db.Table("page_products as pp").
		Select("pp.*").
		Joins("JOIN products as p on p.id = pp.product_id").
		Where("p.is_in_discount = 1 AND pp.discount_price > 0").
		Preload("Product").
		Preload("Page").
		Find(&products).Error
	// Verificar si hay error
	if err != nil {
		return nil, fails.Create("PageProductRepository GetRecommended: No se pudieron obtener los productos recomendados.", err)
	}
	// Regresar data
	return products, nil
}

func (repo PageProductRepository) GetRecommended() ([]dao.PageProduct, error) {
	var products []dao.PageProduct
	// Query principal
	err := repo.db.Table("page_products as pp").
		Select("pp.*").
		Joins("JOIN products as p on p.id = pp.product_id").
		Where("p.is_recommended = 1").
		Preload("Product").
		Preload("Page").
		Find(&products).Error
	// Verificar si hay error
	if err != nil {
		return nil, fails.Create("PageProductRepository GetRecommended: No se pudieron obtener los productos recomendados.", err)
	}
	// Regresar data
	return products, nil
}
