package pages

import (
	"comprix/app/domain/dto"
	"comprix/app/utils"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ScrapPage struct {
	svc            dto.IScrapperService
	log            *logrus.Logger
	maxGoRutines   int
	disableProduct bool
}

func AnalizePageProductsByCategories(conf *dto.Config, page dto.InitScrapPageFn, maxGoRutines int, chunked ...bool) {
	scrapper, err := InitService(conf, page, maxGoRutines)
	// Verificar si hubo error para finalizar funcion
	if err != nil {
		return
	}
	// Cerrar chrome dp
	defer scrapper.svc.CloseChromedpCtx()
	// Eliminar productos no encontrados
	defer scrapper.svc.DisableNotFound()
	// Log
	scrapper.log.Debug("Iniciando el scrapping de categorias")
	// Obtener links de las categorias
	scrapper.ProcessCategoryLinks(scrapper.svc.GetCategoryLinks(), chunked...)
	// Log
	scrapper.log.Debug("Iniciando el scrapping de productos no encontrados")
	// Procesar links no encontrados
	scrapper.processUrls(scrapper.svc.GetNotFoundProductLinks(), []string{})
	// Obtener urls con error
	scrapper.FixErrors()
}

func (scrapper *ScrapPage) processUrls(urls, categories []string) {
	for index, chunkedUrls := range utils.ChunkSlice(urls, 25) {
		start := time.Now()
		// Manejar la creacion o actualizacion de los productos
		scrapper.handleCategoryPageProducts(chunkedUrls, categories)
		// Fecha fin
		end := time.Now()
		// Mostrar mensaje
		scrapper.log.Debug(
			fmt.Sprintf("Se ha terminado de procesar el batch #%d, con una duración de %s", (index + 1), utils.FormatDuration(start, end)),
		)
		// Limpiar
		scrapper.removeChromeProcesses()
	}
}
