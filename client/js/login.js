const api_endpoint = "http://localhost:8073/api/login";

const fail_p = document.getElementById("fail-message")

const form = document.forms[0];

let can_click = true;

async function Login(formdata) {
    const response = await fetch(api_endpoint, {
        method: 'POST',
        body: formdata
    })


    if (response.ok) {
        return {success: true, err: null};
    }

    return {success: false, err: await response.text()};
}

async function Submit(e) {
    e.preventDefault();

    if (!can_click) {
        return
    }

    can_click = false;

    fail_p.style = "display: none;";

    const formData = new FormData(e.target);

    let result = await Login(formData);

    if (result.success) {
        window.location.href = '/home.html';
    } else {
        fail_p.innerHTML = result.err;
        fail_p.style = "display: block;";
    }

    can_click = true;
}

form.addEventListener("submit", Submit) 