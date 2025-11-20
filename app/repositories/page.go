package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type PageRepository struct {
	db *gorm.DB
}

func InitPageRepository(db *gorm.DB) PageRepository {
	return PageRepository{db: db.Model(&dao.Page{})}
}

func (repo PageRepository) FindByName(name string) (*dao.Page, error) {
	var page *dao.Page
	// Validar si existe un error al obtener
	if err := repo.db.Where("name = ?", name).First(&page).Error; err != nil {
		return nil, fails.Create("PageRepository FindByName", err)
	}
	// Regresar data
	return page, nil
}

func (repo PageRepository) FindAll() ([]dao.Page, error) {
	var page []dao.Page
	// Validar si existe un error al obtener
	if err := repo.db.Find(&page).Error; err != nil {
		return nil, fails.Create("PageRepository FindAll", err)
	}
	// Regresar data
	return page, nil
}
