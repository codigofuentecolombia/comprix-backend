package controller

import (
	"comprix/app/domain/dto"
	"comprix/app/scrapper/pages"
	pages_alem "comprix/app/scrapper/pages/alem"
	pages_carrefour "comprix/app/scrapper/pages/carrefour"
	pages_hiperlibertad "comprix/app/scrapper/pages/hiperlibertad"
	pages_jumbo "comprix/app/scrapper/pages/jumbo"
	pages_masonline "comprix/app/scrapper/pages/masonline"
	pages_vea "comprix/app/scrapper/pages/vea"
	"comprix/app/server"
	"sync"

	"github.com/gin-gonic/gin"
)

type ScrapperController struct {
	config *dto.Config
}

func NewScrapperController(config *dto.Config) ScrapperController {
	return ScrapperController{
		config: config,
	}
}

func (ctr ScrapperController) StartScrapper(ginContext *gin.Context) {
	// Ejecutar scrapper en background
	go ctr.runScrapper()
	
	ginContext.JSON(200, server.OK("Scrapper iniciado en background", nil))
}

func (ctr ScrapperController) runScrapper() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ctr.analizeAlemCategories()
		pages.AnalizePageProductsByCategories(ctr.config, pages_carrefour.Initialize, 1, true)
		pages.AnalizePageProductsByCategories(ctr.config, pages_jumbo.Initialize, 1, true)
	}()

	go func() {
		defer wg.Done()
		pages.AnalizePageProductsByCategories(ctr.config, pages_vea.Initialize, 1, true)
		pages.AnalizePageProductsByCategories(ctr.config, pages_hiperlibertad.Initialize, 1, true)
		pages.AnalizePageProductsByCategories(ctr.config, pages_masonline.Initialize, 1, true)
	}()

	wg.Wait()
}

func (ctr ScrapperController) analizeAlemCategories() {
	svc, err := pages_alem.Initialize(ctr.config)
	if err == nil {
		svc.GetCategoryProducts()
	}
}
