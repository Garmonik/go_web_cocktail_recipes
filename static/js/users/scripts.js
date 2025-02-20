async function loadUserProfile() {
    try {
        const response = await fetch('/api/my_user/', {
            method: 'GET',
            credentials: 'include',
        });

        if (!response.ok) throw new Error("Ошибка загрузки данных");

        const data = await response.json();

        document.getElementById("user-avatar").src = `data:image/png;base64,${data.avatar}`;
        document.getElementById("username").textContent = data.username;
        document.getElementById("email").textContent = data.email;
        document.getElementById("bio").textContent = data.bio || "Не указано";
    } catch (error) {
        console.error("Ошибка загрузки профиля:", error);
    }
}

document.addEventListener("DOMContentLoaded", () => {
    loadUserProfile();

    document.getElementById("edit-profile").addEventListener("click", () => {
        window.location.href = "/my_user/update/";
    });
});