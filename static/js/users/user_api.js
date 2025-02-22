export async function loadSomeUserProfile(userId) {
    try {
        const response = await fetch(`/api/user/${userId}/`, {
            method: "GET",
            credentials: "include",
        });
        if (!response.ok) throw new Error("Ошибка загрузки данных");
        const data = await response.json();

        const avatarElement = document.getElementById("user-avatar");
        if (avatarElement && data.avatar) {
            avatarElement.src = `data:image/png;base64,${data.avatar}`;
        }


        document.getElementById("username").textContent = data.username;
        document.getElementById("email").textContent = data.email;
        document.getElementById("bio").textContent = data.bio || "Не указано";

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
        const response = await fetch(`/api/user/${userId}/recipes/?limit=${limit}&offset=${offset}`);
        if (!response.ok) throw new Error("Ошибка загрузки постов");
        const data = await response.json();
        return data.content || [];
    } catch (error) {
        console.error("Ошибка загрузки постов:", error);
        return [];
    }
}

document.addEventListener("DOMContentLoaded", async function () {
    const userId = document.body.dataset.userId;
    if (!userId) {
        console.error("Ошибка: ID пользователя не найден!");
        return;
    }

    await loadSomeUserProfile(userId);

    const contentContainer = document.querySelector("#user-posts");
    let limit = 10;
    let offset = 0;
    let allPostsLoaded = false;

    async function fetchPosts() {
        if (allPostsLoaded) return;

        const data = await loadUserPosts(userId, limit, offset);
        if (data.length === 0) {
            allPostsLoaded = true;
            return;
        }

        renderUserPosts(data, contentContainer);
        if (data.length < limit) {
            allPostsLoaded = true;
        }
    }

    function renderUserPosts(posts, container) {
        if (offset === 0) {
            container.innerHTML = "";
        }

        posts.forEach(post => {
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.innerHTML = `
                <div class="post-header">
                    <h3>${post.name}</h3>
                </div>
                <div class="post-content">
                    <p>${post.description}</p>
                    <img src="data:image/png;base64,${post.image}" alt="Post Image" class="post-image">
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

    document.addEventListener("click", async (event) => {
        const likeButton = event.target.closest(".like-button");
        if (!likeButton) return;

        const postId = likeButton.dataset.postId;
        const likeIcon = likeButton.querySelector(".like-icon");

        if (!postId || !likeIcon) return;

        try {
            const response = await fetch(`/api/recipes/${postId}/like/`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Ошибка при лайке поста");
            }

            // Переключение состояния лайка
            const isLiked = likeIcon.src.includes("like_color.svg");
            likeIcon.src = `/static/images/general/icon/${isLiked ? "like_black.svg" : "like_color.svg"}`;
        } catch (error) {
            console.error("Ошибка лайка:", error);
        }
    });


    window.addEventListener("scroll", () => {
        if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
            offset += limit;
            fetchPosts();
        }
    });

    await fetchPosts();
});