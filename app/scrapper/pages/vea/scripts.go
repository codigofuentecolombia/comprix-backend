package pages_vea

func (service *Service) GetPageProductLinksScript() string {
	return `
        return Array.from(document.querySelectorAll("#gallery-layout-container a")).map( element => element.href );
    `
}

func (service *Service) GetTotalPagesScript() string {
	return `
        let pagesElement = document.querySelector(".discoargentina-search-result-custom-1-x-span-selector-pages")
        // Verificar si existe
        if( pagesElement ){
            let pages = pagesElement.textContent.split(" de ")
            if( pages.length == 2 ){
                return Number(pages[1]);
            }
        }  
        return 1;
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
        const images = Array.from(document.querySelectorAll('.vtex-store-components-3-x-carouselGaleryCursor .vtex-store-components-3-x-productImage img'));
		return images.map(img => img.src);
    `
}

func (service *Service) GetOriginalPriceScript() string {
	return `
        let priceElement = document.querySelector("#priceContainer");
        // Verificar si existe elemento
        if( priceElement ){
            let originalPrice = priceElement.parentNode.parentNode.nextSibling.textContent;
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
            if( document.querySelector("#gallery-layout-container") ){
                let maxTries = 3;
                let currentTries = 0;
                let totalProducts = 20
                let areAllProductsReady = false;
                // Verificar la cantidad de productos
                let productElements = document.querySelectorAll("#gallery-layout-container a")
                while( !areAllProductsReady && (maxTries > currentTries) ) {
                    let currentTotalProducts = productElements?.length ?? totalProducts; //Si es null o undefined definir el total global
                    // Verificar el total de productos
                    if( currentTotalProducts == totalProducts ){
                        areAllProductsReady = true;
                    } else {
                        // Scrollear para recargar productos
                        let reloadElement1 = document.querySelector(".discoargentina-search-result-custom-1-x-pagination-container")
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
            let priceElement = document.querySelector("#priceContainer");

            if (priceElement) {
                let priceText = priceElement.textContent.trim();

                // Expresión regular para verificar si el texto es una cantidad de dinero con formato "$1.350"
                let priceRegex = /^\$\d{1,3}(?:\.\d{3})*(?:,\d+)?$/;

                if (priceRegex.test(priceText)) {
                    // Si el precio es válido, crea el nuevo elemento con ID 'succes-load'
                    let successElement = document.createElement('div');
                    successElement.id = 'succes-load';
                    successElement.textContent = 'Carga exitosa';

                    // Agrega el nuevo elemento al body
                    document.body.appendChild(successElement);
                    return; // Sale de la función y detiene los intentos
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
