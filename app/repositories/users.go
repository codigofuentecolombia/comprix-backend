package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func InitUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (repo UserRepository) FindByID(id interface{}) (*dao.User, error) {
	var user *dao.User
	// Validar si existe un error al obtener
	if err := repo.db.Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fails.Create("UserRepository FindByID", err)
	}
	// Regresar data
	return user, nil
}

func (repo UserRepository) FindByUsername(username string) (*dao.User, error) {
	var user *dao.User
	// Validar si existe un error al obtener
	if err := repo.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fails.Create("UserRepository FindByUsername", err)
	}
	// Regresar data
	return user, nil
}

func (repo UserRepository) FindByUsernameNoRelation(username string) (*dao.User, error) {
	var user *dao.User
	// Validar si existe un error al obtener
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fails.Create("UserRepository FindByUsername", err)
	}
	// Regresar data
	return user, nil
}

func (repo UserRepository) FindByEmail(email string) (*dao.User, error) {
	var user *dao.User
	// Validar si existe un error al obtener
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fails.Create("UserRepository FindByEmail", err)
	}
	// Regresar data
	return user, nil
}

func (repo UserRepository) Create(user dao.User) (dao.User, error) {
	// Intentar crear el usuario en la base de datos
	if err := repo.db.Create(&user).Error; err != nil {
		return dao.User{}, fails.Create("UserRepository Create: No se pudo crear el usuario.", err)
	}
	// Regresar data
	return user, nil
}

func (repo UserRepository) Update(id interface{}, user interface{}) error {
	if err := repo.db.Model(&dao.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return fails.Create("UserRepository Update: No se pudo actualizar", err)
	}
	// Regresar  nil
	return nil
}

func (repo UserRepository) GetPaginated(params dto.GetUsersParams) (dto.Pagination[dao.User], error) {
	var items []dao.User
	var totalItems int64
	// Contar el total de productos
	if err := repo.db.Select("count(*)").Model(&dao.Order{}).Count(&totalItems).Error; err != nil {
		return dto.Pagination[dao.User]{}, fails.Create("UserRepository Paginated Count: No se pudieron contar las ordenes.", err)
	}
	// Obtener productos con paginación
	if err := repo.db.
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Find(&items).Error; err != nil {
		return dto.Pagination[dao.User]{}, fails.Create("UserRepository Paginated: No se pudieron obtener los productos.", err)
	}
	// Calcular el total de páginas
	totalPages := int((totalItems + int64(params.Pagination.Limit) - 1) / int64(params.Pagination.Limit)) // Redondeo hacia arriba
	// Crear y devolver la estructura de paginación
	return dto.Pagination[dao.User]{
		Items:      items,
		Index:      params.Pagination.Index,
		Limit:      params.Pagination.Limit,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	}, nil
}

func (repo UserRepository) VerifyCode(user dao.User, code dao.VerificationCode) error {
	// Marcar como recomendado
	if err := repo.db.Model(&dao.User{}).Where("id = ?", user.ID).Update("is_verified", 1).Error; err != nil {
		return fails.Create("UserRepository VerifyCode: No se marco.", err)
	}
	//
	if err := repo.db.Model(&dao.VerificationCode{}).Where("id = ?", code.ID).Update("claimed", 1).Error; err != nil {
		return fails.Create("UserRepository VerifyCode: No se marco.", err)
	}
	// Regresar sin error
	return nil
}
