export async function loadSomeUserProfile(userId) {
    try {
        const response = await fetch(`/api/user/${userId}/`, {
            method: "GET",
            credentials: "include",
        });
        if (!response.ok) throw new Error("Ошибка загрузки данных");
        const data = await response.json();

        document.getElementById("user-avatar").src = `data:image/png;base64,${data.avatar}`;
        document.getElementById("username").textContent = data.username;
        document.getElementById("email").textContent = data.email;
        document.getElementById("bio").textContent = data.bio || "Не указано";

        // Скрыть кнопку, если это не текущий пользователь
        const editProfileButton = document.getElementById("edit-profile");
        if (!data.my_account && editProfileButton) {
            editProfileButton.style.display = "none";
        }
    } catch (error) {
        console.error("Ошибка загрузки профиля:", error);
    }
}

export async function loadUserPosts(userId, limit = 10, offset = 0) {
    try {
        const response = await fetch(`/api/user/${userId}/posts/?limit=${limit}&offset=${offset}`);
        if (!response.ok) throw new Error("Ошибка загрузки постов");
        return await response.json();
    } catch (error) {
        console.error("Ошибка загрузки постов:", error);
        return [];
    }
}
