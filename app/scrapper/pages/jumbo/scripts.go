package pages_jumbo

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

func (service *Service) LoadAllPageProductsScript() string {
	return `
        let totalTries = 5;
        let currentTries = 1;

        sleep = async(ms) => {
            return new Promise(resolve => setTimeout(resolve, ms));
        }

        reloadProducts = async() => {
            // Obtener elementos de productos
            let productsElements = document.querySelectorAll("#gallery-layout-container .vtex-search-result-3-x-galleryItem")
            // Verificar si se iterara
            while( (productsElements.length < 20) && (currentTries <= totalTries) ){
                // Escrollear hasta los botones
                document.querySelector(".discoargentina-search-result-custom-1-x-pagination-container")?.scrollIntoView({
                    block: "center",     // Alineación del elemento en el centro del contenedor
                    behavior: "smooth", // Desplazamiento suave
                })
                // Esperar 1 segundo
                await sleep(1000)
                // Incrementar intentos
                currentTries++;
                // Refrescar elemento
                productsElements = document.querySelectorAll("#gallery-layout-container .vtex-search-result-3-x-galleryItem")
            }
            // Crear elemento
            const successElement = document.createElement("div");
            successElement.id = "se-ha-completado"
            document.body.appendChild(successElement);
        }

        reloadProducts()
    `
}

func (service *Service) WaitUntilProductPriceIsLoaded() string {
	return `
        let hasStock = true;
        let attemptCount = 0;
        const maxAttempts = 5;

        function checkPrice() {
            // Verificar si no hay stock
            if( document.querySelector(".vtex-flex-layout-0-x-flexCol--product-box .vtex-flex-layout-0-x-flexColChild--product-box:nth-child(3) p") ){
                hasStock = false;
                return;
            }
            // Selecciona el elemento con el ID 'priceContainer'
            let priceElement = document.querySelector(".vtex-store-components-3-x-container .vtex-flex-layout-0-x-flexCol--product-box .vtex-flex-layout-0-x-flexColChild--shelf-main-price-box div span div:nth-of-type(1) div:nth-of-type(1)");
            // Verificar si existe el elemento de precio
            if (priceElement) {
                let priceText = priceElement.textContent.trim();

                // Expresión regular para verificar si el texto es una cantidad de dinero con formato "$1.350"
                let priceRegex = /^\$\d{1,3}(?:\.\d{3})*(?:,\d+)?$/;

                if (priceRegex.test(priceText)) {
                    return; // Sale de la función y detiene los intentos
                }
            }

            // Si no se encontró un precio válido, intenta nuevamente si no se alcanzó el máximo de intentos
            if (attemptCount < maxAttempts) {
                attemptCount++;
                setTimeout(checkPrice, 1500); // Vuelve a intentar después de 1.5 segundos
            }
        }
        // Buscar si existe stock
        checkPrice();
        // Crear elemento
        let successElement = document.createElement('div');
        successElement.id = 'succes-load';
        successElement.textContent = 'Carga exitosa';
        // Agrega el nuevo elemento al body
        document.body.appendChild(successElement);
        // Regresar si tiene 
        return hasStock;
    `
}

func (service *Service) GetPriceDetailScript() string {
	return `
        const element =  document.querySelector(".vtex-store-components-3-x-container .vtex-flex-layout-0-x-flexCol--product-box .vtex-flex-layout-0-x-flexColChild--shelf-main-price-box div")
        const priceInfo = {
            price: "0",
            min_quantity: 1,
            discount_price: "0"
        };
        
        // Función para extraer el primer valor que parezca un precio
        function extractValidPrice(text) {
            const match = text.match(/\$[0-9.,]+/);
            return match ? match[0] : "0";
        }
        
        // Buscar el precio con descuento (si existe)
        const discountPriceEl = element.querySelector("div[class*='-theme-1dCOMij_MzTzZOCohX1K7w']");
        if (discountPriceEl) {
            priceInfo.discount_price = extractValidPrice(discountPriceEl.textContent.trim());
        }
        
        // Buscar el precio normal (sin descuento)
        const normalPriceEl = element.querySelector("div[class*='-theme-2t-mVsKNpKjmCAEM_AMCQH'], div[class*='-theme-3b41v8cnKrHgFgMiReJ6Lp']");
        if (normalPriceEl) {
            priceInfo.price = extractValidPrice(normalPriceEl.textContent.trim());
        }
        
        // Si hay precio con descuento pero no se detectó precio normal, tomar el siguiente div como posible precio normal
        if (priceInfo.discount_price !== "0" && priceInfo.price === "0") {
            const nextPriceEl = discountPriceEl?.parentElement?.nextElementSibling;
            if (nextPriceEl) {
                priceInfo.price = extractValidPrice(nextPriceEl.textContent.trim());
            }
        }
        
        // Si aún no hay precio normal y sí hay descuento, intercambiarlos
        if (priceInfo.price === "0" && priceInfo.discount_price !== "0") {
            priceInfo.price = priceInfo.discount_price;
            priceInfo.discount_price = "0";
        }
        
        // Buscar la cantidad mínima para descuento
        const minQuantityEl = element.querySelector("div[class*='-theme-14k7D0cUQ_45k_MeZ_yfFo']");
        if (minQuantityEl) {
            const match = minQuantityEl.textContent.match(/Llevando (\d+)/);
            if (match) {
                priceInfo.min_quantity = parseInt(match[1], 10);
            }
        }
        // Validar cantidad minima
        if(priceInfo.min_quantity <= 0){
            priceInfo.min_quantity = 4;
        }
        // 
        return priceInfo;
    `
}
