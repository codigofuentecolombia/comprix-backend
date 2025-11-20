package pages_masonline

func (service *Service) GetPageProductLinksScript() string {
	return `
        return Array.from(document.querySelectorAll(".valtech-gdn-search-result-0-x-gallery a")).map( element => element.href );
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
            if( document.querySelector(".valtech-gdn-search-result-0-x-gallery") ){
                let attempts = 0;
                let areAllProductsReady = false;
                // 
                const maxAttempts = 3;
                // Iterar hasta que todo este cargado
                while( !areAllProductsReady && attempts < maxAttempts ) {
                    // Validar que no esten cargados todos los elementos
                    if( document.querySelectorAll(".valtech-gdn-search-result-0-x-gallery a").length >= 24 ){
                        areAllProductsReady = true;
                    } else {
                        // Scrollear para recargar productos
                        let reloadElement = document.querySelector(".valtech-gdn-search-result-0-x-buttonShowMore a")
                        // Validar si existe
                        if(reloadElement){
                            reloadElement.scrollIntoView({
                                behavior: "smooth", // Desplazamiento suave
                                block: "center"     // Alineación del elemento en el centro del contenedor
                            });
                        }
                        // Verificar si ya estan cargados
                        if( document.querySelectorAll(".valtech-gdn-search-result-0-x-gallery a").length >= 24 ){
                            areAllProductsReady = true;
                        }
                        // Esperar 3 segundos
                        await sleep(5)
                    }
                    attempts++;
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
