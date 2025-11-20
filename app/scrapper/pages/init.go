package pages

import (
	"comprix/app/domain/dto"
	log "comprix/app/logger"
	"fmt"
)

func InitService(conf *dto.Config, page dto.InitScrapPageFn, maxGoRutines int) (*ScrapPage, error) {
	// Siempre limitar
	// if maxGoRutines > 3 {
	// 	maxGoRutines /= 3
	// }
	// Structura
	scrapper := &ScrapPage{
		log:            log.Create(conf.Settings.Paths.Logs, "scrapper"),
		maxGoRutines:   maxGoRutines,
		disableProduct: false,
	}
	// Generar servicio
	svc, err := page(conf)
	// Verificar si existe error
	if err != nil {
		scrapper.log.Error(fmt.Sprintf("No se pudo inicializar la pagina: error %v", err))
		// Regresar error
		return nil, err
	}
	// Asignar pagina
	scrapper.svc = svc
	// Regresar estructura
	return scrapper, nil
}
