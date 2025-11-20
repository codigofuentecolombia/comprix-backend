package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type CalendarRepository struct {
	db *gorm.DB
}

func InitCalendarRepository(db *gorm.DB) CalendarRepository {
	return CalendarRepository{db: db}
}

func (repo CalendarRepository) GetAll() ([]dao.Calendar, error) {
	var entities []dao.Calendar
	// Validar si existe un error al obtener
	if err := repo.db.Find(&entities).Error; err != nil {
		return nil, fails.Create("CalendarRepository GetAll", err)
	}
	// Regresar data
	return entities, nil
}

func (repo CalendarRepository) FindByID(id any) (*dao.Calendar, error) {
	var entity *dao.Calendar
	// Validar si existe un error al obtener
	if err := repo.db.Where("id = ?", id).First(&entity).Error; err != nil {
		return nil, fails.Create("CalendarRepository FindByID", err)
	}
	// Regresar data
	return entity, nil
}

func (repo CalendarRepository) Create(entity dao.Calendar) (dao.Calendar, error) {
	// Intentar crear el usuario en la base de datos
	if err := repo.db.Create(&entity).Error; err != nil {
		return dao.Calendar{}, fails.Create("CalendarRepository Create: No se pudo crear el horario.", err)
	}
	// Regresar data
	return entity, nil
}

func (repo CalendarRepository) Update(entities []dto.Calendar) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Iterar elementos
		for _, entity := range entities {
			if err := tx.Model(&dao.Calendar{}).Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
				return fails.Create("CalendarRepository Update: No se pudo actualizar", err)
			}
		}
		// Regresar sin error
		return nil
	})
}

func (repo CalendarRepository) Delete(entity dao.Calendar) error {
	if err := repo.db.Delete(&entity).Error; err != nil {
		return fails.Create("CalendarRepository Delete: No se pudo actualizar", err)
	}
	// Regresar  nil
	return nil
}
