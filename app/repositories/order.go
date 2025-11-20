package repositories

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func InitOrderRepository(db *gorm.DB) OrderRepository {
	return OrderRepository{db: db}
}

func (repo OrderRepository) FindByID(id interface{}) (*dao.Order, error) {
	var order *dao.Order
	// Validar si existe un error al obtener
	if err := repo.db.Preload("UserShippingAddress").Where("id = ?", id).First(&order).Error; err != nil {
		return nil, fails.Create("OrderRepository FindByID", err)
	}
	// Regresar data
	return order, nil
}

func (repo OrderRepository) FindAll() ([]dao.Order, error) {
	// Ordenes
	var orders []dao.Order
	// Query
	db := repo.db.Model(&dao.Order{}).Preload("Items", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("PageProduct", func(tx2 *gorm.DB) *gorm.DB {
			return tx2.Unscoped().Preload("Page").Preload("Product")
		})
	}).Where("status = ?", "pending")
	// Validar si existe un error al obtener
	if err := db.Find(&orders).Error; err != nil {
		return nil, fails.Create("OrderRepository FindAll", err)
	}
	// Regresar ordenes
	return orders, nil
}

func (repo OrderRepository) FindAllByID(id interface{}, user_id interface{}) ([]dao.PageProduct, error) {
	db := repo.db.
		Table("page_products AS pp ").
		Unscoped().
		Select("pp.*, op.quantity as quantity").
		Joins("JOIN order_products AS op ON op.page_product_id = pp.id").
		Joins("JOIN orders AS o ON o.id = op.order_id").
		Preload("Product").
		Preload("Page").
		Preload("Product.Category").
		Preload("Product.Brand").
		Where("o.id = ?", id)

	var products []dao.PageProduct
	// Validar si existe un error al obtener
	if err := db.Find(&products).Error; err != nil {
		return nil, fails.Create("OrderRepository FindAllByUserID", err)
	}
	// // Regresar data
	return products, nil
}

func (repo OrderRepository) FindAllByUserID(userID interface{}) ([]dao.PageProduct, error) {
	db := repo.db.
		Table("page_products AS pp ").
		Select("pp.*, op.quantity as quantity").
		Joins("JOIN order_products AS op ON op.page_product_id = pp.id").
		Joins("JOIN orders AS o ON o.id = op.order_id").
		Preload("Product").
		Preload("Page").
		Preload("Product.Category").
		Preload("Product.Brand").
		Where("o.user_id = ?", userID)

	var products []dao.PageProduct
	// Validar si existe un error al obtener
	if err := db.Find(&products).Error; err != nil {
		return nil, fails.Create("OrderRepository FindAllByUserID", err)
	}
	// // Regresar data
	return products, nil
}

func (repo OrderRepository) Create(user dao.User, products []dao.PageProduct, form dto.NewOrder) (dao.Order, error) {
	tx := repo.db.Begin()
	// Crear referencia de envio
	address := dao.UserShippingAddress{
		Time:           form.Time,
		Date:           form.ShippingAddress.Date,
		City:           "Simoca",
		State:          "Tucumán",
		Colony:         "Simoca",
		UserID:         user.ID,
		Street:         form.ShippingAddress.Street,
		Reference:      form.ShippingAddress.Reference,
		PostalCode:     "4107",
		PhoneNumber:    form.ShippingAddress.PhoneNumber,
		ExternalNumber: "",
		InternalNumber: nil,
	}
	// Guardar address
	if err := tx.Create(&address).Error; err != nil {
		tx.Rollback()
		return dao.Order{}, fails.Create("OrderRepository Create: No se pudo crear la direccion.", err)
	}
	// Inicializar orden
	order := dao.Order{
		UserID:                user.ID,
		Status:                "pending",
		ShippingCost:          form.ShippingCost,
		PaymentMethod:         form.PaymentMethod,
		UserShippingAddressID: address.ID,
	}
	// Calcular totales
	for _, product := range products {
		// Sumar total
		order.Subtotal += product.Price * float64(product.Quantity)
		// Verificar si tiene descuento y si cumple con el minimo
		if product.DiscountPrice > 0 && product.Quantity >= product.MinQuantityToApplyDiscount {
			order.TotalDiscount += product.DiscountPrice * float64(product.Quantity)
		}
	}
	// Sumar
	order.Total = order.Subtotal + order.ShippingCost
	// Intentar crear el usuario en la base de datos
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return dao.Order{}, fails.Create("OrderRepository Create: No se pudo crear la orden.", err)
	}
	// Productos de orden
	var orderProducts []dao.OrderProduct
	// Iterar productos
	for _, product := range products {
		if product.ID != 0 {
			orderProducts = append(orderProducts, dao.OrderProduct{
				OrderID:       order.ID,
				Quantity:      product.Quantity,
				PageProductID: product.ID,
			})
		}
	}
	// Guardar productos
	if err := tx.Create(&orderProducts).Error; err != nil {
		tx.Rollback() // Revertir si hay error
		return dao.Order{}, fails.Create("OrderRepository Create: No se pudo crear la lista de articulos.", err)
	}
	//
	if err := tx.Commit().Error; err != nil {
		return dao.Order{}, fails.Create("OrderRepository Create: No se pudo guardar.", err)
	}
	// Regresar data
	return order, nil
}

func (repo OrderRepository) GetPaginated(params dto.GetOrdersParams) (dto.Pagination[dao.Order], error) {
	var orders []dao.Order
	var totalItems int64
	// Crear db
	db := repo.db.Session(&gorm.Session{})
	// Validar si existe el nombre del usuario
	if params.UserID != nil {
		db = db.Where("user_id = ?", params.UserID)
	}
	// Contar el total de productos
	if err := db.Model(&dao.Order{}).Count(&totalItems).Error; err != nil {
		return dto.Pagination[dao.Order]{}, fails.Create("OrderRepository Paginated Count: No se pudieron contar las ordenes.", err)
	}
	// Obtener productos con paginación
	if err := db.
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset).
		Find(&orders).Error; err != nil {
		return dto.Pagination[dao.Order]{}, fails.Create("OrderRepository Paginated: No se pudieron obtener los productos.", err)
	}
	// Calcular el total de páginas
	totalPages := int((totalItems + int64(params.Pagination.Limit) - 1) / int64(params.Pagination.Limit)) // Redondeo hacia arriba
	// Crear y devolver la estructura de paginación
	return dto.Pagination[dao.Order]{
		Items:      orders,
		Index:      params.Pagination.Index,
		Limit:      params.Pagination.Limit,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
	}, nil
}

func (repo OrderRepository) UpdateStatus(id interface{}, status dao.OrderStatus) error {
	// Actualizar
	if err := repo.db.Model(&dao.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fails.Create("OrderRepository UpdateStatus: No se pudo actualizar.", err)
	}
	// Regresar limpio
	return nil
}

func (repo OrderRepository) GetStatistics() dto.OrderStatistics {
	var stats dto.OrderStatistics
	// Query
	query := `
		SELECT 
			(SELECT COUNT(*) FROM orders WHERE STATUS = 'pending') AS pending,
			(SELECT COUNT(*) FROM orders WHERE STATUS = 'completed') AS completed,
			(SELECT COALESCE(SUM(total), 0) FROM orders WHERE STATUS = 'completed') AS total
	`
	// Poblar estructura
	repo.db.Raw(query).Scan(&stats)
	// Regresar estructura
	return stats
}
