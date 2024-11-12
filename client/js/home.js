const images_endpoint = "http://localhost:8073/images/";
const api_endpoint = "http://localhost:8073/api/";
const category_list = ["Новинки сезона", "Кошкодевочки", "Мужские персонажи"];

const catalog = document.getElementsByClassName("catalog-section")[0];
const load_smoke = document.getElementById("load-smoke");

function formatNumber(number) {
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

async function GetCategory(category_name) {
    try {
        const response = await fetch(`${api_endpoint}get-category?category=${category_name}`, {
            method: 'GET'
        })

        if (response.ok) {
            const result = await response.json();
            console.log('Success:', result);
            return result;
        } else {
            console.error('Server error:', response.statusText);
        }
    } catch(error) {
        console.error('try error:', error);
    }
}

async function FillCatalog() {
    for (let i = 0; i < category_list.length; i++) {
        let category_name = category_list[i];
        let card = ``
        card +=`
        <div class="catalog-sub-section">
            <h3 class="regular-font sub-catalog-name">
                ${category_name}
            </h3>
            <div class="sub-catalog-selection">
        `;

        let products = await GetCategory(category_name);
        console.log("Get products: ", products);
        for (let j = 0; j < products.length; j++) {
            let product = products[j];
            card += `
                <div class="selectable sub-catalog-card">
                    <a href="product.html?id=${product.id}">
                        <img class="catalog-card-image" src="${images_endpoint}${product.image}" alt="dakimakura">
                        <p class="catalog-card-name regular-font">${product.name}</p>
                        <p class="catalog-card-price medium-font">${formatNumber(product.price)} руб.</p>
                    </a>
                </div>
            `;
        }
        
        card +=`
            </div>
        </div>
        `;

        catalog.innerHTML += card;
    }

    load_smoke.style="opacity:0;"
}

FillCatalog();