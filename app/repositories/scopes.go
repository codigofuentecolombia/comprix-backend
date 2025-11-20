package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func brandHasProducts(query *gorm.DB, table ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query = query.Model(dao.Product{}).Select("DISTINCT(brand_id)")

		return db.Where(fmt.Sprintf("%s in (?)", getTableColumn("id", table...)), query)
	}
}

func hasBrandIdsScope(brandIds any, table ...string) func(db *gorm.DB) *gorm.DB {
	return isInKeyValuesScope(getTableColumn("brand_id", table...), brandIds)
}

func isInKeyValuesScope(key string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s in ?", key), value)
	}
}

func HasIdScope(id any, table ...string) func(db *gorm.DB) *gorm.DB {
	return hasKeyValueScope(getTableColumn("id", table...), id)
}

func HasPageIDScope(id any, table ...string) func(db *gorm.DB) *gorm.DB {
	return hasKeyValueScope(getTableColumn("page_id", table...), id)
}

func HasUrlScope(id any, table ...string) func(db *gorm.DB) *gorm.DB {
	return hasKeyValueScope(getTableColumn("url", table...), id)
}

func hasKeyValueScope(key string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", key), value)
	}
}

func diffNameScope(name string, table ...string) func(db *gorm.DB) *gorm.DB {
	return diffKeyValueScope(getTableColumn("name", table...), name)
}

func diffKeyValueScope(key string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s != ?", key), value)
	}
}

func PreloadScope(preload []dto.RepositoryGormParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preload {
			db = db.Preload(preload.Query, preload.Args)
		}
		// Regresar
		return db
	}
}

func OlderThanOneDay(db *gorm.DB) *gorm.DB {
	// Obtener la fecha de hace 1 día
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	// Agregar la condición a la consulta
	return db.Where("created_at < ? AND updated_at < ?", oneDayAgo, oneDayAgo)
}

func BestProductPriceScope(db *gorm.DB) *gorm.DB {
	subquery := db.Table("page_products as pp").
		Select("pp.main_product_id, MIN(pp.price) as min_price").
		Where("pp.main_product_id IS NOT NULL").
		Group("pp.main_product_id")

	return db.
		Joins("JOIN (?) AS sub ON page_products.main_product_id = sub.main_product_id AND page_products.price = sub.min_price", subquery).
		Joins("JOIN products AS p on p.id = pp1.product_id")
}
