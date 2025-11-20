package order_controller

import (
	"comprix/app/server"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func (ctr OrderController) DownloadReport(ginContext *gin.Context) {
	// Obtener ordenes
	orders, err := ctr.repositories.order.FindAll()
	// Validar si se pudo obtener la informacion
	if err != nil {
		server.InternalErrorException("No se pudo obtener los datos de las ordenes", nil)
	}
	// Crear archivo excel
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rowIndex := 1
	// Header
	headers := []string{"Tienda", "Producto", "Cantidad", "Precio unitario", "Total"}
	// Crear estilo
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#1babbb"}, // fondo
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:  true,
			Color: "#ffffff",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	// Agregar header al archivo
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, rowIndex)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, style)
	}
	// Incrementar fila
	rowIndex++
	// Iterar ordenes
	for _, order := range orders {
		// Iterar items de orden
		for _, item := range order.Items {
			price := item.PageProduct.Price
			// Verificar si tiene descuento
			if item.PageProduct.DiscountPrice > 0 {
				price = item.PageProduct.DiscountPrice
			}
			// Agregar data
			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), item.PageProduct.Page.Name)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), item.PageProduct.Product.Name)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), item.Quantity)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), price)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), (price * float64(item.Quantity)))
			// Incrementar index de fila
			rowIndex++
		}
	}
	// Guardar el archivo en un buffer
	buf, err := f.WriteToBuffer()
	// Verificar si existe error
	if err != nil {
		server.InternalErrorException("No se pudo generar el archivo de Excel", nil)
	}
	// Regresar
	ginContext.Header("Content-Disposition", `attachment; filename="reporte.xlsx"`)
	ginContext.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
