const api_endpoint = "http://localhost:8073/api/login";

const fail_p = document.getElementById("fail-message")

const form = document.forms[0];

let can_click = true;

function checkValidData(formdata) {
    return formdata.get("username").length >= 5 && formdata.get("password").length >= 5
}

async function Login(formdata) {
    if (!checkValidData(formdata)) {
        return {success: false, err: "Некорректный логин или пароль"};
    }

    const response = await fetch(api_endpoint, {
        method: 'POST',
        body: formdata
    })

    return {success: response.ok, err: await response.text()};
}

async function Submit(e) {
    e.preventDefault();

    if (!can_click) {
        return
    }

    can_click = false;

    fail_p.style = "display: none;";

    const formData = new FormData(e.target);

    let response = await Login(formData);

    if (response.success) {
        window.location.href = '/home.html';
    } else {
        fail_p.innerHTML = response.err;
        fail_p.style = "display: block;";
    }

    can_click = true;
}

form.addEventListener("submit", Submit) 