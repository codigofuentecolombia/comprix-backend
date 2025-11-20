package pages_carrefour

func (service *Service) GetPageProductLinksScript() string {
	return `
        return Array.from(document.querySelectorAll(".valtech-carrefourar-search-result-3-x-gallery .vtex-product-summary-2-x-clearLink")).map( element => element.href );
    `
}

func (service *Service) GetTotalPagesScript() string {
	return `
        let totalPages = 1;
        let pageElements = document.querySelectorAll(".vtex-render__container-id-search-result-layout .valtech-carrefourar-search-result-3-x-paginationButtonPages");
        // Obtener el ultimo valor
        if( pageElements.length ){
            let lastPageElement = pageElements[pageElements.length - 1].querySelector("div");
            // Obtener el total de paginas
            if( lastPageElement ){
                totalPages = Number(lastPageElement.textContent)
            }
        }
        // Regresar total
        return totalPages;
    `
}

func (service *Service) GetCategoriesScript() string {
	return `
        let categoryContainer = document.querySelector(".vtex-breadcrumb-1-x-container");
        let categories = [];
        // Verificar que exista
        if( categoryContainer ){
            // Iterar elementos
            for(let el of categoryContainer.childNodes ){
                // Verificar que no sea la ultima categoria
                if( el.textContent ){
                    categories.push(el.textContent);
                }
            }
            // Verificar si existen
            categories.pop()
        }
        // Regresar categorias
        return categories;
    `
}

func (service *Service) GetImagesScript() string {
	return `
        const images = Array.from(document.querySelectorAll('.vtex-store-components-3-x-productImagesGallerySwiperContainer img'));
		return images.map(img => img.src);
    `
}

func (service *Service) GetOriginalPriceScript() string {
	return `
        let priceElement = document.querySelector(".valtech-carrefourar-product-price-0-x-listPriceValue");
        // Verificar si existe elemento
        if( priceElement ){
            let originalPrice = priceElement.textContent;
            // Verificar si es distinto a laburando
            if( originalPrice != "laburando" ){
                return originalPrice
            }
        }
        // Regresar cero
        return "0";
    `
}

func (service *Service) LoadAllPageProducts() string {
	return `
        sleep = async(seconds) => new Promise(resolve => setTimeout(resolve, seconds * 1000))
        // Esperar a que se cargen todos los productos
        waitAllProducts = async() => {
            if( document.querySelector(".valtech-carrefourar-search-result-3-x-gallery") ){
                let maxTries = 5;
                let currentTries = 0;
                let totalProducts = 16;
                let areAllProductsReady = false;
                // Verificar la cantidad de productos
                let productElements = document.querySelectorAll(".valtech-carrefourar-search-result-3-x-gallery .vtex-product-summary-2-x-clearLink")
                while( !areAllProductsReady && (maxTries > currentTries) ) {
                    let currentTotalProducts = productElements?.length ?? totalProducts; //Si es null o undefined definir el total global
                    // Verificar el total de productos
                    if( currentTotalProducts == totalProducts ){
                        areAllProductsReady = true;
                    } else {
                        // Scrollear para recargar productos
                        let reloadElement1 = document.querySelector(".render-route-store-search-category")
                        // Validar si existe
                        if(reloadElement1){
                            reloadElement1.scrollIntoView({
                                behavior: "smooth", // Desplazamiento suave
                                block: "center"     // Alineación del elemento en el centro del contenedor
                            });
                        } else {
                            areAllProductsReady = true;
                        }
                        // Esperar 3 segundos
                        await sleep(5)
                    }
                    // Aumentar intentos
                    currentTries++;
                }
                // Crear elemento para saber cuando haya finalizado
                const successElement = document.createElement("div");
                successElement.id = "se-ha-completado"
                document.body.appendChild(successElement);
            }
        }
        // Consultar funcion
        waitAllProducts()
    `
}

func (service *Service) WaitUntilProductPriceIsLoaded() string {
	return `
        let attemptCount = 0;
        const maxAttempts = 5;

        function checkPrice() {
            // Selecciona el elemento con el ID 'priceContainer'
            let priceElement = document.querySelector(".valtech-carrefourar-product-price-0-x-currencyContainer");
            // Verificar si  existe
            if (priceElement) {
                let priceText = priceElement.textContent.trim();
                // Quitar inecesarios
                priceText = priceText.replaceAll(" ", "");
                priceText = priceText.replaceAll(",", "");
                priceText = priceText.replaceAll(".","");
                // Expresión regular para verificar si el texto es una cantidad de dinero con formato "$1.350"
                let priceRegex = /^\$\d+$/;
                // Verificar si es correcto
                if (priceRegex.test(priceText)) {
                    // Si el precio es válido, crea el nuevo elemento con ID 'succes-load'
                    let successElement = document.createElement('div');
                    successElement.id = 'succes-load';
                    successElement.textContent = 'Carga exitosa';
                    // Agrega el nuevo elemento al body
                    document.body.appendChild(successElement);
                    // Sale de la función y detiene los intentos
                    return; 
                }
            }
            // Si no se encontró un precio válido, intenta nuevamente si no se alcanzó el máximo de intentos
            if (attemptCount < maxAttempts) {
                attemptCount++;
                setTimeout(checkPrice, 1500); // Vuelve a intentar después de 1.5 segundos
            }
        }
        // 
        checkPrice();
    `
}
