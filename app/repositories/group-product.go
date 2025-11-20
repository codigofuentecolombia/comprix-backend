package repositories

import (
	"gorm.io/gorm"
)

func (repo *PageProductRepository) MinPriceQuery() *gorm.DB {
	subQuery := repo.db.Table("page_products pp2").
		Select("MIN(pp2.price)").
		Joins("JOIN product_group_members pgm2 ON pp2.product_id = pgm2.product_id").
		Where("pgm2.group_id = pg.id AND pp2.deleted_at = 0 AND pp2.price > 0")

	return repo.db.Table("product_groups pg").
		Select("pp.*").
		Joins("JOIN product_group_members pgm ON pg.id = pgm.group_id").
		Joins("JOIN products p ON pgm.product_id = p.id").
		Joins("JOIN page_products pp ON p.id = pp.product_id").
		Where("pp.deleted_at = 0 AND pp.price > 0").
		Where("pp.price = (?)", subQuery).
		Group("pg.id")
}
