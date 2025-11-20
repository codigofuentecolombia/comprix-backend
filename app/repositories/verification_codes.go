package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type VerificationCodeRepository struct {
	db *gorm.DB
}

func InitVerificationCodeRepository(db *gorm.DB) VerificationCodeRepository {
	return VerificationCodeRepository{db: db}
}

func (repo VerificationCodeRepository) Create(vCode dao.VerificationCode) (dao.VerificationCode, error) {
	// Intentar crear el usuario en la base de datos
	if err := repo.db.Create(&vCode).Error; err != nil {
		return dao.VerificationCode{}, fails.Create("VerificationCodeRepository Create: No se pudo crear el codigo.", err)
	}
	// Regresar data
	return vCode, nil
}

func (repo VerificationCodeRepository) FindUsersCode(userID uint, code string) (*dao.VerificationCode, error) {
	var entity *dao.VerificationCode
	// Validar si existe un error al obtener
	if err := repo.db.Where("user_id = ? AND code = ? AND claimed = 0", userID, code).First(&entity).Error; err != nil {
		return nil, fails.Create("VerificationCodeRepository FindUsersCode", err)
	}
	// Regresar data
	return entity, nil
}
