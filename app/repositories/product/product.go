package repository_product

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func InitRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Find(params dto.ProductRepositoryFindParams) dto.RepositoryGenericResponse[dao.Product] {
	return repositories.Find("Repository Find", dao.Product{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *Repository) FindAll(params dto.ProductRepositoryFindParams) dto.RepositoryGenericResponse[[]dao.Product] {
	return repositories.Find("Repository FindAll", []dao.Product{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *Repository) GenerateQueryFromFindParams(params dto.ProductRepositoryFindParams) *gorm.DB {
	model := dao.Product{}
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
	// Por categoria
	if params.CategoryID != nil {
		db = db.Scopes(repo.InCategoryScope(*params.CategoryID))
	}
	// Filtrar por nombre
	if params.Name != nil && *params.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", *params.Name))
	}
	// Ordenar por
	if params.OrderBy != nil {
		db = db.Order(params.OrderBy)
	}
	// Regresar instancia de base de datos
	return db
}
