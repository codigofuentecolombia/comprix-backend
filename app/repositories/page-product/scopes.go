package repository_page_product

import (
	"comprix/app/domain/dto"
	"strings"

	"gorm.io/gorm"
)

func (repo *Repository) WithDistinctPageProducts(db *gorm.DB) *gorm.DB {
	subquery := repo.db.
		Table("page_products").
		Select("product_id").
		Group("product_id").
		Having("COUNT(DISTINCT page_id) > ?", 3)

	return db.Where("p.id IN (?)", subquery)
}

func (repo *Repository) BestPriceScope(db *gorm.DB) *gorm.DB {
	subquery := repo.db.Table("page_products as pp").
		Select([]string{
			"pp.id",
			"pp.main_product_id",
			"LEAST(pp.price, COALESCE(NULLIF(pp.discount_price, 0), pp.price)) AS min_price",
			"ROW_NUMBER() OVER (PARTITION BY pp.main_product_id ORDER BY LEAST(pp.price, COALESCE(NULLIF(pp.discount_price, 0), pp.price)) ASC) AS rn",
		}).
		Where("pp.main_product_id IS NOT NULL").
		Where("pp.deleted_at = 0")

	return db.
		Joins("JOIN (?) AS sub ON page_products.id = sub.id AND sub.rn = 1", subquery).
		Joins("JOIN products AS p on p.id = page_products.product_id")
}

func (repo *Repository) BestPricePerPageScope(db *gorm.DB) *gorm.DB {
	subquery := repo.db.Table("page_products as pp").
		Select("MIN(LEAST(pp.price, COALESCE(NULLIF(pp.discount_price, 0), pp.price))) as min_price, pp.main_product_id, pp.page_id, pp.id").
		Where("pp.main_product_id IS NOT NULL").
		Where("pp.deleted_at = 0").
		Group("pp.main_product_id, pp.page_id")

	return db.
		Joins("JOIN (?) AS sub ON page_products.id = sub.id", subquery).
		Joins("JOIN products AS p ON p.id = page_products.product_id")
}

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
		return db.Where("page_products.main_product_id IN (SELECT id FROM products WHERE category_id IN ?)", categoryIDs)
	}
}

func (repo *Repository) ShowByTypeScope(orderType dto.PageProductType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch orderType {
		case dto.DiscountsPageProductType:
			return db.Where("page_products.discount_price > 0 AND p.is_disabled = 0")
		case dto.RecommendedPageProductType:
			return db.Where("p.is_recommended = 1 AND p.is_disabled = 0")
		case dto.OffersPageProductType:
			return db.Where("p.is_in_discount = 1 AND p.is_disabled = 0")
		case dto.DisablePageProductType:
			return db.Where("p.is_disabled = 1")
		default:
			return db
		}
	}
}

func (repo *Repository) CustomOrder(orderBy string) string {
	switch strings.ToLower(orderBy) {
	case "id_asc":
		return "page_products.id ASC"
	case "id_desc":
		return "page_products.id DESC"
	case "category_asc":
		return "p.category_id ASC"
	case "category_desc":
		return "p.category_id DESC"
	case "price_asc":
		return "sub.min_price ASC"
	case "price_desc":
		return "sub.min_price DESC"
	case "admin":
		return "p.id, p.category_id ASC"
	default:
		return "page_products.id ASC"
	}
}

func (repo *Repository) SearchScope(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var categoryIDs []uint32
		// Crear buscador
		searcher := search + "%"
		// Generar consulta
		err := repo.db.Raw(`
			WITH RECURSIVE category_hierarchy AS (
				SELECT id FROM categories WHERE id = (SELECT id FROM categories WHERE name like ? ORDER BY id ASC LIMIT 1)
				UNION ALL
				SELECT c.id
				FROM categories c
				INNER JOIN category_hierarchy ch ON c.parent_id = ch.id
			)
			SELECT id FROM category_hierarchy
		`, searcher).Scan(&categoryIDs).Error
		// Verificar si hubo un error
		if err != nil {
			return db
		}
		// Verificar si se encontro algun resultado
		if len(categoryIDs) > 0 {
			return db.Where("page_products.main_product_id IN (SELECT id FROM products WHERE category_id IN ?)", categoryIDs)
		}
		// Buscar por nombre de producto
		return db.Where("p.name like ?", "%"+searcher)
	}
}
