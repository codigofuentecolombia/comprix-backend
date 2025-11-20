package service_scrapper

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"

	"github.com/sirupsen/logrus"
)

func (s *Service) HandleProduct(pageProduct *dto.RetrievedProduct) {
	// // Verificar si existe
	if pageProduct != nil {
		// Manejar producto obtenido
		if err := s.productSvc.HandleRetrieved(pageProduct); err != nil {
			// Verificar si el error
			if fail, ok := err.(fails.Fail); ok && fail.Err != nil {
				// Registrar en la base de datos si no se esta omitiendo
				if boolean, ok := fail.Data.(bool); !ok || !boolean {
					s.Repositories.Error.Create(
						dao.Error{
							Url:        pageProduct.Url,
							Error:      "",
							PageID:     s.Page.ID,
							Message:    err.Error(),
							Categories: pageProduct.Categories,
						},
					)
				}
				// Crear log de error
				s.CreateWarningLog(err.Error(), nil)
			} else {
				s.CreateWarningLog("Ocurrio un error al tratar de crear o actualizar el producto.", logrus.Fields{
					"error": err,
				})
			}
		}
	}
}

func (s *Service) GetExistingProducts() ([]dao.PageProduct, error) {
	omitRecentProducts := true
	includeSoftdeleted := true
	// Consultar productos
	response := s.Repositories.Product.FindAll(dto.ProductRepositoryFindParams{
		PageID: &s.Page.ID,
		// Selects: &dto.RepositoryGormSelections{Query: "DISTINCT url, product_id, page_id"},
		Selects:         &dto.RepositoryGormSelections{Query: "DISTINCT url"},
		OrderBy:         "price, discount_price asc",
		Softdeleted:     &includeSoftdeleted,
		OlderThanOneDay: &omitRecentProducts,
	})
	// Verificar si ocurrio un error
	if response.Error != nil {
		return []dao.PageProduct{}, response.Error
	}
	// Retornar productos
	return response.Data, nil
}

func (s *Service) DisableNotFound() {
	err := s.Repositories.Product.DisableNotFoundInPage(s.InitializedAt, s.Page.ID)
	// Validar si se encontro error
	if err != nil {
		s.CreateWarningLog("Ocurrio un error al tratar de desactivar productos no encontrados.", logrus.Fields{
			"error": err,
		})
	}
}

func (s *Service) GetNotFoundProductLinks() []string {
	response := s.Repositories.Product.FindNotFoundInPage(s.InitializedAt, s.Page.ID)
	// Validar si se encontro error
	if response.Error != nil {
		s.CreateWarningLog("Ocurrio un error al tratar de obtener los productos no encontrados.", logrus.Fields{
			"error": response.Error,
		})
	}
	// Crear links
	links := []string{}
	// Iterar productos
	for _, product := range response.Data {
		links = append(links, product.Url)
	}
	// Regresar links
	return links
}
