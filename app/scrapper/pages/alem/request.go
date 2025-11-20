package pages_alem

import (
	"bytes"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Options    *[]Option    `json:"options"`
	OptionsMar *[]OptionMar `json:"optionsMar"`
}

type Option struct {
	Cod         int     `json:"cod"`
	Des         string  `json:"des"`
	WebDes      string  `json:"web_des"`
	WebDet      *string `json:"web_det"`
	Uni         int     `json:"uni"`
	Peso        string  `json:"peso"`
	WebLstPrs   *string `json:"web_lst_prs"`
	Mar         int     `json:"mar"`
	Dep         int     `json:"dep"`
	LstDep      *string `json:"lst_dep"`
	MoneSimbolo string  `json:"mone_simbolo"`
	StatZ       bool    `json:"stat_z"`
	StatM       *string `json:"stat_m"`
	PrecM       string  `json:"prec_m"`
	PrecbulM    string  `json:"precbul_m"`
	FileNi      int     `json:"file_ni"`
	Stkdispo    int     `json:"stkdispo"`
	OrderSearch string  `json:"order_search"`
}

type OptionMar struct {
	Mar  int    `json:"mar"`
	Des  string `json:"des"`
	Cant int    `json:"cant"`
}

func (service *Service) RequestByCategory(id string, categories []string) (*[]dto.RetrievedProduct, error) {
	number, err := strconv.Atoi(id)
	// Verificar si se pudo obtener
	if err != nil {
		return nil, fails.Create("RequestByCategory: El id no es numerico", err)
	}
	// Crear variable de parametros
	var params map[string]string
	// Verificar si es mayor a 26 que es el total de departamentos
	if number <= 26 {
		params = map[string]string{
			"grpdep": id,
		}
	} else {
		params = map[string]string{
			"dep": id,
		}
	}
	// Convierte los datos a JSON
	jsonData, err := json.Marshal(params)
	// Verificar si se pudo convertir
	if err != nil {
		return nil, fails.Create("RequestByCategory: No se pudieron crear parametros de peticion", err)
	}
	// Crea la solicitud HTTP
	req, err := http.NewRequest("POST", "https://tienda.alem500.com.ar/ccd_cliente/tienda/api/optionsPlu", bytes.NewBuffer(jsonData))
	// Verificar si hubo error
	if err != nil {
		return nil, fails.Create("RequestByCategory: No se pudo realizar la peticion", err)
	}
	// Agrega el header 'Content-Type: application/json'
	req.Header.Set("Content-Type", "application/json")
	// Crea un cliente HTTP
	client := &http.Client{}
	// ejecuta la solicitud
	resp, err := client.Do(req)
	// Verificar si ocurrio error
	if err != nil {
		return nil, fails.Create("RequestByCategory: No se pudo realizar la peticion", err)
	}
	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	// Verificar si se pudo leer
	if err != nil {
		return nil, fails.Create("RequestByCategory: No se pudo leer respuesta", err)
	}
	// Cerrar body
	defer resp.Body.Close()
	// Lee la respuesta
	var response Response
	// Verificar si se pudo leer
	err = json.Unmarshal(body, &response)
	// Aquí hacemos el Unmarshal (convertir el JSON en estructura Go)
	if err != nil {
		return nil, fails.Create("RequestByCategory: No se pudo parsear respuesta", err)
	}
	// Crear productos
	products := []dto.RetrievedProduct{}
	// Iterar opciones si es que existe
	if response.Options != nil && len(*response.Options) > 0 && response.OptionsMar != nil && len(*response.OptionsMar) > 0 {
		// Iterar opciones
		for _, option := range *response.Options {
			product := dto.RetrievedProduct{
				Url:    fmt.Sprintf("https://tienda.alem500.com.ar/%d", option.Cod),
				Sku:    fmt.Sprintf("%d", option.Cod),
				Name:   option.WebDes,
				Images: []string{fmt.Sprintf("https://tienda.alem500.com.ar/drive/file/index/0/%d", option.FileNi)},
			}
			// Verificar si tiene stock
			hasStock := option.Stkdispo > 0
			// Obtener precio original sin ultimo caracter
			if len(option.PrecM) > 0 {
				product.OriginalPrice = option.PrecM[:len(option.PrecM)-1]
			}
			// Obtener precio
			product.Price, _ = utils.CleanCurrencyFormat(product.OriginalPrice)
			product.PageID = service.svc.Page.ID
			product.HasStock = &hasStock
			product.Categories = categories
			// Iterar marcas
			for _, brand := range *response.OptionsMar {
				if brand.Mar == option.Mar {
					product.Brand = brand.Des
					break
				}
			}
			// Agregar producto
			products = append(products, product)
		}
		// Regresar productos
		return &products, nil
	}
	//
	return nil, fails.Create("RequestByCategory: No se encontraron productos", err)
}

func (service *Service) GetCategoryProducts() {
	products := []dto.RetrievedProduct{}
	// Iterar categorias
	for index, categoryLink := range service.GetCategoryLinks() {
		service.svc.CreateDebugLog(fmt.Sprintf("Obteniendo productos de la categoria #%s", categoryLink.Link), nil)
		// Mostrar en donde estamos
		fmt.Printf("Index #%d\n", index)
		// Obtener producto
		retrievedProducts, err := service.RequestByCategory(categoryLink.Link, categoryLink.Categories)
		// Verificar que no haya error
		if err != nil {
			service.svc.CreateWarningLog(err.Error(), logrus.Fields{"link": categoryLink.Link})
			// Ir al siguiente
			continue
		}
		// Verificar que haya producto
		if retrievedProducts != nil {
			products = append(products, *retrievedProducts...)
		}
	}
	// Total de productos
	totalProducts := len(products)
	// Iterar productos y guardar
	for index, retrievedProduct := range products {
		fmt.Printf("Procesando producto #%d de %d resultados\n", (index + 1), totalProducts)
		// Procesar producto
		service.svc.HandleProduct(&retrievedProduct)
	}
	// Desactivar no encontrados
	service.svc.DisableNotFound()
}
