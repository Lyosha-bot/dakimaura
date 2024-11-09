const api_check = "http://localhost:8073/api/fetch";
const api_logout = "http://localhost:8073/api/logout";

const header_nav = document.getElementById("header-nav");
const logout_wrapper = document.getElementById("logout-wrapper");
const logout = document.getElementById("logout")
const cancel = document.getElementById("cancel")

async function Fetch() {
    try {
        const response = await fetch(api_check, {
            method: 'GET',
        });

        console.log(response.ok)
    
        return {success: response.ok, value: await response.text()};
    } catch(error) {
        return {success: false, value: error};
    }
}

async function Logout() {
    try {
        const response = await fetch(api_logout, {
            method: 'POST',
        });

        return {success: response.ok, err: await response.text()};
    } catch(error) {
        return {success: false, value: error};
    }
}

async function LogoutFunc() {
    const result = await Logout()

    if (result.success) {
        window.location.href = '/home.html';
    }
}

Fetch().then((response) => {
    if (response.success) {
        header_nav.innerHTML += `
        <button id="profile-button" class="header-nav-tab">${response.value}</button>
        `
        document.getElementById("profile-button").onclick = () => {
            logout_wrapper.style = ""
        }
    } else {
        header_nav.innerHTML += `
        <a class="header-nav-tab" href="login.html">Вход/Регистрация</a>
        `
    }
})

logout.onclick = LogoutFunc
cancel.onclick = () => {
    logout_wrapper.style = "display: none;"
}