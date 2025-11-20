package repository_page_product

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

func (repo *Repository) Disable(params dto.ProductRepositoryFindParams) error {
	return repositories.Disable("ProductRepository Disable", dao.PageProduct{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *Repository) Find(params dto.ProductRepositoryFindParams) dto.RepositoryGenericResponse[dao.PageProduct] {
	return repositories.Find("ProductRepository Find", dao.PageProduct{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *Repository) FindAll(params dto.ProductRepositoryFindParams) dto.RepositoryGenericResponse[[]dao.PageProduct] {
	return repositories.Find("ProductRepository FindAll", []dao.PageProduct{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *Repository) GetPaginated(params dto.ProductRepositoryFindParams) dto.RepositoryGenericResponse[dto.Pagination[dao.PageProduct]] {
	return repositories.GetPaginated(*params.Pagination, repo.GenerateQueryFromFindParams(params))
}
