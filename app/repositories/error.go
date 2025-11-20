package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type ErrorRepository struct {
	db *gorm.DB
}

func InitErrorRepository(db *gorm.DB) *ErrorRepository {
	return &ErrorRepository{db: db}
}

func (repo *ErrorRepository) Create(entity dao.Error) dto.RepositoryGenericResponse[dao.Error] {
	// Validar creacion
	if err := repo.db.Create(&entity).Error; err != nil {
		return dto.RepositoryGenericResponse[dao.Error]{
			Data:  entity,
			Error: fails.Create("ErrorRepository Create: No se pudo crear", err),
		}
	}
	// Regresar sin error
	return dto.RepositoryGenericResponse[dao.Error]{
		Data:  entity,
		Error: nil,
	}
}

func (repo *ErrorRepository) Disable(params dto.ErrorRepositoryFindParams) error {
	return Disable("ErrorRepository Disable", dao.PageProduct{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *ErrorRepository) Find(params dto.ErrorRepositoryFindParams) dto.RepositoryGenericResponse[dao.Error] {
	return Find("ErrorRepository Find", dao.Error{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *ErrorRepository) FindAll(params dto.ErrorRepositoryFindParams) dto.RepositoryGenericResponse[[]dao.Error] {
	return Find("ErrorRepository FindAll", []dao.Error{}, repo.GenerateQueryFromFindParams(params))
}

func (repo *ErrorRepository) GetPaginated(params dto.ErrorRepositoryFindParams) dto.RepositoryGenericResponse[dto.Pagination[dao.Error]] {
	return GetPaginated(params.Pagination, repo.GenerateQueryFromFindParams(params))
}

func (repo *ErrorRepository) GenerateQueryFromFindParams(params dto.ErrorRepositoryFindParams) *gorm.DB {
	db := repo.db.Model(&dao.Error{})
	// Validar si tiene nombre
	if params.PageID != nil {
		db = db.Scopes(HasPageIDScope(*params.PageID))
	}
	// Validar si tiene url
	if params.Url != nil && *params.Url != "" {
		db = db.Scopes(HasUrlScope(*params.Url))
	}
	// Verificar si tiene preloads
	if params.Preloads != nil {
		db = db.Scopes(PreloadScope(*params.Preloads))
	}
	// Verificar si se ordenara
	if params.ShouldOrderDesc != nil && *params.ShouldOrderDesc {
		db = db.Order("id DESC")
	}
	// Regresar instancia de base de datos
	return db
}
