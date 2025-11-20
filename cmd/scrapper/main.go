package main

import (
	"flag"
	"log"
	"sync"

	"comprix/app/config"
	"comprix/app/domain/dto"
	"comprix/app/scrapper/pages"

	pages_alem "comprix/app/scrapper/pages/alem"
	pages_carrefour "comprix/app/scrapper/pages/carrefour"
	pages_hiperlibertad "comprix/app/scrapper/pages/hiperlibertad"
	pages_jumbo "comprix/app/scrapper/pages/jumbo"
	pages_masonline "comprix/app/scrapper/pages/masonline"
	pages_vea "comprix/app/scrapper/pages/vea"
	service_scrapper "comprix/app/services/scrapper"
)

var cnf *dto.Config
var svc service_scrapper.Service
var configPath string

func init() {
	// Define flag for configuration path
	flag.StringVar(&configPath, "cnf", "settings/conf.yaml", "Path to the configuration file")
	flag.Parse()
	// Obtener configuracion
	cnf = HandleError(config.Load(configPath))
}

func main() {
	// DebugProduct()
	AnalizeCategories()
}

func AnalizeExistingProducts() {
	pages.AnalizeExistingPageProducts(cnf, pages_jumbo.Initialize, cnf.Settings.Scrapping.MaxGoRutines)
	pages.AnalizeExistingPageProducts(cnf, pages_hiperlibertad.Initialize, cnf.Settings.Scrapping.MaxGoRutines)
	pages.AnalizeExistingPageProducts(cnf, pages_carrefour.Initialize, cnf.Settings.Scrapping.MaxGoRutines)
	pages.AnalizeExistingPageProducts(cnf, pages_vea.Initialize, cnf.Settings.Scrapping.MaxGoRutines)
	pages.AnalizeExistingPageProducts(cnf, pages_masonline.Initialize, cnf.Settings.Scrapping.MaxGoRutines)
}

func AnalizeCategories() {
	var wg sync.WaitGroup
	// vamos a lanzar 2 goroutines
	wg.Add(2)

	go func() {
		defer wg.Done()
		AnalizeAlemCategories()
		pages.AnalizePageProductsByCategories(cnf, pages_carrefour.Initialize, 1, true)
		pages.AnalizePageProductsByCategories(cnf, pages_jumbo.Initialize, 1, true)
	}()

	go func() {
		defer wg.Done()
		pages.AnalizePageProductsByCategories(cnf, pages_vea.Initialize, 1, true)
		pages.AnalizePageProductsByCategories(cnf, pages_hiperlibertad.Initialize, 1, true)
		pages.AnalizePageProductsByCategories(cnf, pages_masonline.Initialize, 1, true)
	}()
	// espera hasta que ambas goroutines terminen
	wg.Wait()
}

func DebugProduct() {
	scrapper, _ := pages.InitService(cnf, pages_jumbo.Initialize, 1)
	// Obtener producto
	scrapper.CreateOrUpdateProduct("https://www.jumbo.com.ar/cerveza-erdinger-negra-500-ml/p", []string{})
}

func AnalizeAlemCategories() {
	svc, err := pages_alem.Initialize(cnf)
	// Validar que no haya error
	if err == nil {
		svc.GetCategoryProducts()
	}
}

func HandleError[T any](data T, err error) T {
	// Verificar si existe error
	if err != nil {
		log.Fatal(err)
	}
	// Regresar data
	return data
}
