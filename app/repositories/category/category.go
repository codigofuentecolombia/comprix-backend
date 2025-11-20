package repository_category

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func InitRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Find(params dto.CategoryRepositoryFindParams) dto.RepositoryGenericResponse[dao.Category] {
	return repositories.Find("CategoryRepository Find", dao.Category{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *Repository) FindAll(params dto.CategoryRepositoryFindParams) dto.RepositoryGenericResponse[[]dao.Category] {
	return repositories.Find("CategoryRepository FindAll", []dao.Category{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *Repository) GenerateQueryFromFindParams(params dto.CategoryRepositoryFindParams) *gorm.DB {
	model := dao.Category{}
	// Buscar
	db := repo.db.Model(&model)
	// Verificar si se buscara por id
	if params.ID != nil && *params.ID != "" {
		db = db.Where("id = ?", *params.ID)
	}
	// Verificar si se buscara por padre
	if params.ParentID != nil && *params.ParentID != "" {
		db = db.Where("parent_id = ?", *params.ParentID)
	}
	// Verificar si se buscara por nobmre
	if params.Name != nil && *params.Name != "" {
		db = db.Where("name = ?", *params.Name)
	}
	// Verificar si se obtendran solo principales
	if params.OnlyParents != nil && *params.OnlyParents {
		db = db.Where("parent_id IS NULL")
	}
	// Verificar si tiene preloads
	if params.Preloads != nil {
		db = db.Scopes(repositories.PreloadScope(*params.Preloads))
	}
	// Regresar instancia de base de datos
	return db
}
