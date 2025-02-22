document.addEventListener("DOMContentLoaded", async function () {
    const contentContainer = document.getElementById("user-posts");
    const postsHeader = document.querySelector(".posts-header");
    let limit = 10;
    let offset = 0;
    let allPostsLoaded = false;

    async function loadUserProfile() {
        try {
            const response = await fetch("/api/my_user/", {
                method: "GET",
                credentials: "include",
            });
            if (!response.ok) throw new Error("Ошибка загрузки данных");
            const data = await response.json();

            document.getElementById("user-avatar").src = `data:image/png;base64,${data.avatar}`;
            document.getElementById("username").textContent = data.username;
            document.getElementById("username").dataset.userId = data.id;
            document.getElementById("email").textContent = data.email;
            document.getElementById("bio").textContent = data.bio || "Не указано";

            loadUserPosts(data.id);
        } catch (error) {
            console.error("Ошибка загрузки профиля:", error);
        }
    }

    async function loadUserPosts(userId) {
        if (allPostsLoaded) return;

        try {
            const response = await fetch(`/api/user/${userId}/recipes/?limit=${limit}&offset=${offset}`);
            if (!response.ok) throw new Error("Ошибка загрузки постов");
            const posts = await response.json();

            if (posts.content.length === 0) {
                allPostsLoaded = true;
                return;
            }

            renderUserPosts(posts.content);
            offset += posts.content.length;

            if (posts.content.length < limit) {
                allPostsLoaded = true;
            }
        } catch (error) {
            console.error("Ошибка загрузки постов:", error);
        }
    }

    function renderUserPosts(posts) {
        if (offset === 0) {
            contentContainer.innerHTML = "";
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
                        <img src="${post.like ? "/static/images/general/icon/like_color.svg" : "/static/images/general/icon/like_black.svg"}"
                             alt="Like" class="like-icon">
                    </button>
                </div>
            `;
            contentContainer.appendChild(postElement);
        });

        postsHeader.style.opacity = "1";
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
            if (!response.ok) throw new Error("Ошибка при лайке поста");

            const isLiked = likeIcon.src.includes("like_color.svg");
            likeIcon.src = isLiked ? "/static/images/general/icon/like_black.svg" : "/static/images/general/icon/like_color.svg";
        } catch (error) {
            console.error("Ошибка лайка:", error);
        }
    });

    window.addEventListener("scroll", () => {
        if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
            const userId = document.getElementById("username").dataset.userId;
            if (userId) {
                loadUserPosts(userId);
            }
        }
    });

    await loadUserProfile();
    document.getElementById("edit-profile").addEventListener("click", () => {
        window.location.href = "/my_user/update/";
    });
});