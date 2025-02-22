import { loadSomeUserProfile, loadUserPosts } from "./user_api.js";

document.addEventListener("DOMContentLoaded", async () => {
    try {
        const userId = document.body.dataset.userId;
        if (!userId) {
            console.error("Ошибка: ID пользователя не найден!");
            return;
        }

        await loadSomeUserProfile(userId);

        const contentContainer = document.getElementById("user-posts");
        if (!contentContainer) {
            console.error("Ошибка: контейнер #user-posts не найден!");
            return;
        }

        const posts = await loadUserPosts(userId);
        renderUserPosts(posts, contentContainer);

        const editProfileButton = document.getElementById("edit-profile");
        if (editProfileButton) {
            editProfileButton.addEventListener("click", () => {
                window.location.href = "/my_user/update/";
            });
        } else {
            console.warn("Предупреждение: кнопка #edit-profile не найдена.");
        }
    } catch (error) {
        console.error("Ошибка инициализации страницы:", error);
    }
});

function renderUserPosts(posts, container) {
    if (!Array.isArray(posts)) {
        console.error("Ошибка: данные постов некорректны", posts);
        return;
    }

    container.innerHTML = "";
    posts.forEach(post => {
        const postElement = document.createElement("div");
        postElement.classList.add("post");
        postElement.innerHTML = `
            <div class="post-header">
                <h3>${post.name || "Без названия"}</h3>
            </div>
            <div class="post-content">
                <p>${post.description || "Без описания"}</p>
                ${post.image ? `<img src="data:image/png;base64,${post.image}" alt="Post Image" class="post-image">` : ""}
            </div>
            <div class="post-footer">
                <button class="like-button" data-post-id="${post.id}">
                    <img src="/static/images/general/icon/${post.like ? "like_color.svg" : "like_black.svg"}" 
                         alt="Like" class="like-icon">
                </button>
            </div>
        `;
        container.appendChild(postElement);
    });
}

