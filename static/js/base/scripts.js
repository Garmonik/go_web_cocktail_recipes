// Подключение шаблонов (AJAX)
fetch("/static/templates/base/header.html")
    .then(response => response.text())
    .then(data => document.getElementById("header").innerHTML = data);

fetch("/static/templates/base/sidebar.html")
    .then(response => response.text())
    .then(data => document.getElementById("sidebar").innerHTML = data);

document.addEventListener("DOMContentLoaded", function () {
    let images = [
        "/static/images/general/home/home1.png",
        "/static/images/general/home/home2.png",
        "/static/images/general/home/home3.png",
        "/static/images/general/home/home4.png",
        "/static/images/general/home/home5.png",
        "/static/images/general/home/home6.png",
    ];

    document.getElementById("randomImage1").src = images[Math.floor(Math.random() * images.length)];
    document.getElementById("randomImage2").src = images[Math.floor(Math.random() * images.length)];
});

document.addEventListener("DOMContentLoaded", function () {
    const logoutLink = document.querySelector('a[href="/login/"]');

    if (logoutLink) {
        logoutLink.addEventListener("click", async function (event) {
            event.preventDefault();

            try {
                const response = await fetch("/api/logout/", {
                    method: "POST",
                    credentials: "include",
                    headers: {
                        "X-CSRFToken": getCSRFToken()
                    }
                });

                if (response.ok) {
                    window.location.href = "/login/";
                } else {
                    console.error("Ошибка при выходе:", response.status);
                }
            } catch (error) {
                console.error("Ошибка сети:", error);
            }
        });
    }

    function getCSRFToken() {
        const match = document.cookie.match(/csrftoken=([^;]+)/);
        return match ? match[1] : "";
    }
});

async function loadUserInfo() {
    try {
        const response = await fetch('/api/user/short/', {
            method: 'GET',
            credentials: 'include',
        });

        if (!response.ok) throw new Error("Ошибка загрузки пользователя");

        const data = await response.json();
        const userContainer = document.querySelector(".user-container");

        userContainer.innerHTML = `
                <span class="username">${data.username}</span>
                <img class="avatar" src="data:image/png;base64,${data.avatar}" alt="Аватар">
            `;
    } catch (error) {
        console.error("Ошибка загрузки данных пользователя:", error);
    }
}

document.addEventListener("DOMContentLoaded", loadUserInfo);