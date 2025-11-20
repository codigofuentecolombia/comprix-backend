package main

import (
	"comprix/app/config"
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	service_product "comprix/app/services/product"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
)

var cnf *dto.Config
var configPath string

type Group struct {
	ID   int
	Name string
}

func init() {
	// Define flag for configuration path
	flag.StringVar(&configPath, "cnf", "settings/conf.yaml", "Path to the configuration file")
	flag.Parse()
	// Obtener configuracion
	cnf = HandleError(config.Load(configPath))
}

func main() {
	// ValidateData()
	CompareNames()
	// ReasignarProductosDuplicados()

}

func CompareNames() {
	svc := service_product.InitService(cnf)
	// Comparar nombres
	svc.AggroupByNames()
	// svc.CompareNamesWithDetails("Aceite Mezcla Siglo De Oro 900 Ml", "Aceite Mezcla Siglo De Oro 900ml")

}

func ReasignarProductosDuplicados() {
	// Listado de nombres duplicados
	var duplicatedNames []string
	// Solo nombres no eliminados y duplicados
	err := cnf.GormDB.Model(&dao.Product{}).Unscoped().Select("name").Group("name").Having("COUNT(*) > 1").Pluck("name", &duplicatedNames).Error
	// Verificar si ocurrio error
	if err != nil {
		fmt.Println(fails.Create("No se pudieron obtener nombres duplicados", err))
		return
	}
	// Concurrencia
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20)
	// Iterar nombres duplicados
	for _, name := range duplicatedNames {
		wg.Add(1)
		semaphore <- struct{}{}
		// Concurrencia
		go func(name string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			ReasignarProductosDuplicadosPorNombre(name)
		}(name)
	}
	// Esperar a que finaliice
	wg.Wait()
}

func ReasignarProductosDuplicadosPorNombre(name string) {
	// Variable para obtener el producto mas antiguo
	var oldest dao.Product
	// Consultar producto mas antiguo
	err := cnf.GormDB.Unscoped().Where("name = ?", name).Order("created_at ASC").First(&oldest).Error
	// Verificar si ocurrio un error
	if err != nil {
		fmt.Println(fails.Create(fmt.Sprintf("No se pudo obtener el producto más antiguo para '%s'", name), err))
		return
	}
	// Variable para productos mas recientes
	var otherProducts []dao.Product
	// Consultar productos mas recientes
	err = cnf.GormDB.Unscoped().Where("name = ? AND id != ?", name, oldest.ID).Find(&otherProducts).Error
	// Verificar si ocurrio un error al consultar productos
	if err != nil {
		fmt.Println(fails.Create(fmt.Sprintf("No se pudo obtener los productos mas recientes para '%s'", name), err))
		return
	}
	//
	var hasChanges bool
	// Iterar productos mas recientes
	for _, prod := range otherProducts {
		// actualizaciones
		updates := map[string]interface{}{
			"product_id":      oldest.ID,
			"main_product_id": oldest.ID,
		}
		// Actualizar
		res := cnf.GormDB.Model(&dao.PageProduct{}).Unscoped().Where("product_id = ?", prod.ID).Updates(updates)
		// Verificar si se pudo actualizar
		if res.Error != nil {
			fmt.Println(fails.Create(fmt.Sprintf("No se pudo obtener los productos mas recientes para '%s'", name), res.Error))
			return
		} else if res.RowsAffected > 0 {
			fmt.Printf("Pageproducts actualizados %d del productID %d a %d\n", res.RowsAffected, prod.ID, oldest.ID)
		} else {
			fmt.Printf("No se actualizaron productos de pagina con nombre '%s'\n", name)
		}
		// Actualizar producto dependiendo de estatus
		if prod.IsDisabled {
			hasChanges = true
			oldest.IsDisabled = prod.IsDisabled
		}
		// Actualizar producto dependiendo de estatus
		if prod.IsRecommended {
			hasChanges = true
			oldest.IsRecommended = prod.IsRecommended
		}
		// Actualizar producto dependiendo de estatus
		if prod.IsInDiscount {
			hasChanges = true
			oldest.IsInDiscount = prod.IsInDiscount
		}
		// Eliminar producto
		if err := cnf.GormDB.Unscoped().Delete(&prod).Error; err != nil {
			fmt.Println(fails.Create(fmt.Sprintf("No se pudo eliminar producto duplicado '%s'", name), err))
		} else {
			fmt.Printf("Producto eliminado (ID: %d, nombre: '%s')\n", prod.ID, name)
		}
	}
	// Actualizar viejo si tiene cambio
	if hasChanges {
		if err := cnf.GormDB.Model(&oldest).Updates(oldest).Error; err != nil {
			fmt.Println(fails.Create(fmt.Sprintf("No se pudo actualizar producto viejo '%s'", name), err))
		}
	}
}

func HandleError[T any](data T, err error) T {
	// Verificar si existe error
	if err != nil {
		fmt.Println(err)
	}
	// Regresar data
	return data
}

func ValidateData() {
	// Leer archivo
	data := ReadFile()
	// productSvc := service_product.InitService(cnf)
	// Iterar data
	for name, names := range data {
		ids := []uint{}
		_, originalIndex := getProductIndex(name)
		// Iterar nombres
		for _, subname := range names {
			_, originalSubindex := getProductIndex(subname)
			// Obtener porcentaje
			if originalSubindex == 0 {
				fmt.Printf("No se pudo obtener el subindex\n")
			} else {
				ids = append(ids, originalSubindex)
			}
		}
		// Verificar si se encontro
		if originalIndex != 0 && validateIndexes(append(ids, originalIndex)) {
			SyncProducts(append(ids, originalIndex))
		}
	}
}

func validateIndexes(ids []uint) bool {
	for i, mainID := range ids {
		for j, subID := range ids {
			if i != j && mainID > (subID-100) && mainID < (subID+100) {
				return false
			}
		}
	}
	return true
}

func getProductIndex(name string) (string, uint) {
	vals := strings.Split(name, " - ")
	// Validar tamaño
	if len(vals) == 0 {
		return "", 0
	}
	// Obtener
	index := vals[len(vals)-1]
	name = strings.Join(vals[:len(vals)-1], " - ")
	// Regresar
	id, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		return "", 0
	}
	// rEGRESAR
	return name, uint(id)
}

func ReadFile() map[string][]string {
	// Leer el archivo JSON
	data := HandleError(os.Open("datos.json"))
	defer data.Close()
	// Crear un mapa para almacenar el resultado
	var result map[string][]string
	// Decodificar el contenido del archivo JSON en el mapa
	if err := json.NewDecoder(data).Decode(&result); err != nil {
		log.Fatalf("Error al decodificar el JSON: %v", err)
	}
	// Regresar data
	return result
}

func SyncProducts(ids []uint) {
	err := cnf.GormDB.Transaction(func(tx *gorm.DB) error {
		var oldest dao.Product
		var products []dao.Product
		// Buscar producto mas viejo
		if err := tx.Unscoped().Where("id in ?", ids).Order("created_at ASC").Take(&oldest).Error; err != nil {
			return fails.Create("No se pudo obtener el producto más antiguo", err)
		}
		// Buscar productos
		if err := tx.Where("id in ? AND id != ?", ids, oldest.ID).Find(&products).Error; err != nil {
			return fails.Create("No se pudo obtener los productos mas recientes", err)
		}
		//
		var hasChanges bool
		// Iterar productos mas recientes
		for _, prod := range products {
			// actualizaciones
			updates := map[string]interface{}{
				"product_id":      oldest.ID,
				"main_product_id": oldest.ID,
			}
			// Actualizar
			res := tx.Model(&dao.PageProduct{}).Unscoped().Where("product_id = ?", prod.ID).Updates(updates)
			// Verificar si se pudo actualizar
			if res.Error != nil {
				return fails.Create("No se pudo actualizar", res.Error)
			}
			// Actualizar producto dependiendo de estatus
			if prod.IsDisabled {
				hasChanges = true
				oldest.IsDisabled = prod.IsDisabled
			}
			// Actualizar producto dependiendo de estatus
			if prod.IsRecommended {
				hasChanges = true
				oldest.IsRecommended = prod.IsRecommended
			}
			// Actualizar producto dependiendo de estatus
			if prod.IsInDiscount {
				hasChanges = true
				oldest.IsInDiscount = prod.IsInDiscount
			}
			// Eliminar producto
			if err := tx.Unscoped().Delete(&prod).Error; err != nil {
				return fails.Create("No se pudo eliminar producto duplicado", err)
			}
		}
		// Actualizar viejo si tiene cambio
		if hasChanges {
			if err := tx.Model(&oldest).Updates(oldest).Error; err != nil {
				return fails.Create("No se pudo actualizar producto viejo", err)
			}
		}
		// Regresar sin error
		return nil
	})
	// Imprimir error
	fmt.Println(err)
}
