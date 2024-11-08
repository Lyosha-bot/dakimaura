const images_endpoint = "http://localhost:8073/images/";
const api_endpoint = "http://localhost:8073/api/get-product?id=";

const imageElem = document.getElementById("image");
const nameElem = document.getElementById("name");
const IDElem = document.getElementById("id");
const priceElem = document.getElementById("price");
const characterElem = document.getElementById("character");
const materialElem = document.getElementById("material");
const brandElem = document.getElementById("brand");
const timeElem = document.getElementById("time");

const urlParams = new URLSearchParams(window.location.search);
const productId = urlParams.get('id');

function formatNumber(number) {
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

async function GetProduct(id) {
    try {
        const response = await fetch(`${api_endpoint}${id}`, {
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

async function LoadProduct(id) {
    const product = await GetProduct(id);

    if (!product) {
        window.location.href = '/home.html';
        return;
    }

    imageElem.src = `${images_endpoint}${product.image}`;
    nameElem.innerHTML = `Дакимакура ${product.name}`;
    IDElem.innerHTML = `ID товара ${product.id}`;
    priceElem.innerHTML = `${formatNumber(product.price)} руб.`;
    characterElem.innerHTML = product.name;
    materialElem.innerHTML = product.material;
    brandElem.innerHTML = product.brand;
    timeElem.innerHTML = product.produce_time;
}

if (productId) {
    LoadProduct(productId);
} else {
    window.location.href = '/home.html';
}