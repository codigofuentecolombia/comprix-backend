package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/utils"

	"gorm.io/gorm"
)

type BrandhRepository struct {
	db *gorm.DB
}

func InitBrandhRepository(db *gorm.DB) *BrandhRepository {
	return &BrandhRepository{db: db}
}

func (repo *BrandhRepository) FindAll(params dto.BrandRepositoryFindParams) dto.RepositoryGenericResponse[[]dao.Brand] {
	return Find("BrandhRepository FindAll", []dao.Brand{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *BrandhRepository) AggroupInProducts(brands []int) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&dao.Product{}).Where("brand_id IN ?", brands).Update("brand_id", utils.GetMinNumber(brands)).Error
		// Verificar si se pudieron actualizar los productos
		if err != nil {
			return fails.Create("BrandhRepository AggroupInProducts: no se pudo actualizar el la marca del producto", err)
		}
		// Regresar sin errores
		return nil
	})
}

func (repo *BrandhRepository) GenerateQueryFromFindParams(params dto.BrandRepositoryFindParams) *gorm.DB {
	model := dao.Brand{}
	//
	db := repo.db.Model(&model)
	// Verificar si se obtendran las que tienen productos
	if params.HasProducts != nil && *params.HasProducts {
		db = db.Scopes(brandHasProducts(repo.db, model.TableName()))
	}
	// Verificar si se ordenara
	if params.OrderBy != nil {
		// Verificar por tamaño
		if params.OrderBy.MaxNameLengthDesc != nil && *params.OrderBy.MaxNameLengthDesc == true {
			db = db.Order("LENGTH(name) DESC")
		}
		// Verificar si es por nombre
		if params.OrderBy.NameDesc != nil && *params.OrderBy.NameDesc == true {
			db = db.Order("name DESC")
		}
	}
	// Regresar instancia de base de datos
	return db
}
