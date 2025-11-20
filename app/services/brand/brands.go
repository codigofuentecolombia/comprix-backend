package service_brand

import (
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/repositories"
	"comprix/app/utils"
	"slices"
)

type Service struct {
	config     *dto.Config
	repository *repositories.BrandhRepository
}

func InitService(cnf *dto.Config) Service {
	return Service{
		config:     cnf,
		repository: repositories.InitBrandhRepository(cnf.GormDB),
	}
}

func (svc *Service) AggroupByNames() error {
	orderFlag := true
	// Verificar si se pudo obtener
	response := svc.repository.FindAll(dto.BrandRepositoryFindParams{OrderBy: &dto.BrandRepositoryFindOrderBy{
		NameDesc:          &orderFlag,
		MaxNameLengthDesc: &orderFlag,
	}})
	// Verificar si hay error
	if response.Error != nil {
		return fails.Create("GroupByNames", response.Error)
	}
	// Limpiar nombres
	for index, brand := range response.Data {
		response.Data[index].Name = utils.SanitizeString(brand.Name)
	}
	// Variable de elementos encontrados
	namesInUse := []string{}
	groupedNames := dto.GroupedBrands{}
	uncategorized := []string{}
	// Iterar nombres
	for i, currentBrand := range response.Data {
		if !slices.Contains(namesInUse, currentBrand.Name) {
			percentaje := float64(0.80)
			brandNameSize := len(currentBrand.Name)
			// Calcular porcentaje de coincidencia
			if brandNameSize == 5 {
				percentaje = 0.75
			}
			//
			currentGroup := []int{currentBrand.ID}
			// Iterar marcas omitiendo registro actual
			for j, brand := range response.Data {
				if j != i {
					// Verificar que no se haya procesado aun
					if !slices.Contains(namesInUse, brand.Name) {
						if utils.LevenshteinSimilarity(currentBrand.Name, brand.Name) >= percentaje {
							currentGroup = append(currentGroup, brand.ID)
							namesInUse = append(namesInUse, brand.Name)
						}
					}
				}
			}
			//
			namesInUse = append(namesInUse, currentBrand.Name)
			//
			if len(currentGroup) > 1 {
				groupedNames[currentBrand.ID] = currentGroup
			} else {
				uncategorized = append(uncategorized, currentBrand.Name)
			}
		}
	}
	// Iterar grupos
	for _, group := range groupedNames {
		svc.repository.AggroupInProducts(group)
	}
	// Regresar
	return nil
}
