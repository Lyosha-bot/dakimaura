const api_endpoint = "http://localhost:8073/api/register";

const fail_p = document.getElementById("fail-message")

const form = document.forms[0];

let can_click = true;

async function Register(formdata) {
    const response = await fetch(api_endpoint, {
        method: 'POST',
        body: formdata
    })


    if (response.ok) {
        return true;
    }

    const text = await response.text();

    return false, text;
}

async function Submit(e) {
    e.preventDefault();

    if (!can_click) {
        return
    }

    can_click = false;

    fail_p.style = "display: none;";

    const formData = new FormData(e.target);

    let success, err = Register(formData);

    console.log(success)

    if (success) {
        // window.location.href = '/login.html?success=1';
    } else {
        fail_p.innerHTML = err;
        fail_p.style = "display: block;";
    }

    can_click = true;
}

form.addEventListener("submit", Submit) 