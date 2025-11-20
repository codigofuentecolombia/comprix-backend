package service_email

import (
	"comprix/app/domain/dao"
	"fmt"
)

func (s *Service) Order(user dao.User, order dao.Order, products []dao.PageProduct) {
	// Contenido del correo en HTML
	body := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Orden Canjeada con Éxito</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">

			<table align="center" width="600" style="background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0px 0px 10px #ccc;">
				<tr>
					<td align="center">
						<img src="https://comprix.com.ar/assets/img/logo-dark.png" alt="Logo" style="max-width: 300px; filter: drop-shadow(5px 5px 10px rgba(0, 0, 0, 0.5)); rotate: -10deg;">

						<h2 style="color: #4CAF50; font-family: Arial, sans-serif;">✅ ¡Tu orden ha sido canjeada con éxito!</h2>

						<table style="width: 100%; font-family: Arial, sans-serif;">
							<tr>
								<td style="font-weight: bold; padding-bottom: 5px;">Folio #{id}</td>
							</tr>
							<tr>
								<td style="padding-top: 5px;">
									{items} <!-- Aquí irían los detalles de los productos -->
								</td>
							</tr>
							<tr>
								<td style="padding-top: 10px;">
									<ul style="list-style: none; padding-left: 0; font-size: 0.875rem;">
										<li style="display: flex; justify-content: space-between; padding-bottom: 5px;">
											<span>Subtotal:</span>
											<span>{subtotal}</span>
										</li>
										<li style="display: flex; justify-content: space-between; padding-bottom: 5px; border-bottom: 1px solid #e3e9ef;">
											<span>Costo de envío:</span>
											<span>{shipping_cost}</span>
										</li>
										<li style="display: flex; justify-content: space-between; font-weight: bold;">
											<span>Total:</span>
											<span>{total}</span>
										</li>
									</ul>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>
    `
	items := ""
	// Iterar productos
	for _, product := range products {
		item := `
			<table width="100%" cellpadding="0" cellspacing="0" style="border-bottom: 1px solid #e3e9ef; padding-bottom: 0.5rem;">
				<tr>
					<td style="padding-right: 10px; padding-bottom: 12px; width: 64px;">
						<a href="#" style="display: block; width: 64px; height: 64px;">
							<img src="{image}" alt="Product" style="width: 64px; height: auto; border-radius: 4px;">
						</a>
					</td>
					<td>
						<table width="100%" cellpadding="0" cellspacing="0">
							<tr>
								<td style="font-size: 0.875rem; font-weight: 500; color: #373f50; padding-bottom: 5px;">
									<a href="#" style="color: #373f50; text-decoration: none;">{name}</a>
								</td>
							</tr>
							<tr>
								<td style="color: #7d879c; font-size: 0.75rem;">
									<div style="display: flex; align-items: center;">
										<span> Vende: </span>
										<a href="#" style="color: #7d879c; font-weight: 500; text-decoration: none;">{page}</a>
									</div>
								</td>
							</tr>
							<tr>
								<td style="padding-top: 10px;">
									<table width="100%" cellpadding="0" cellspacing="0">
										<tr>
											<td style="color: rgba(78, 84, 200, 1); font-weight: 500; padding-right: 10px; padding-bottom: 5px;">
												{price}
											</td>
											<td style="color: #7d879c; font-size: 0.75rem;">
												{quantity} unidades
											</td>
										</tr>
									</table>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		`
		// Remplazar valores
		item = replacePlaceholder(item, "{page}", product.Page.Name)
		item = replacePlaceholder(item, "{image}", product.Images[0])
		item = replacePlaceholder(item, "{name}", product.Product.Name)
		item = replacePlaceholder(item, "{price}", fmt.Sprintf("$%.2f", product.Price))
		item = replacePlaceholder(item, "{quantity}", fmt.Sprintf("%d", product.Quantity))
		item = replacePlaceholder(item, "{page_image}", product.Page.Logo)
		// Agregar al listado
		items += item
	}
	// Reemplazar placeholders con los valores del usuario
	body = replacePlaceholder(body, "{id}", fmt.Sprintf("%05d", order.ID))
	body = replacePlaceholder(body, "{items}", items)
	body = replacePlaceholder(body, "{total}", fmt.Sprintf("$%.2f", order.Total))
	body = replacePlaceholder(body, "{subtotal}", fmt.Sprintf("$%.2f", order.Subtotal))
	body = replacePlaceholder(body, "{shipping_cost}", fmt.Sprintf("$%.2f", order.ShippingCost))

	s.Send(user.Email, "✅ ¡Tu orden ha sido canjeada con éxito!", body)
}
