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
        const postsHeader = document.querySelector(".posts-header");
        let limit = 10;
        let offset = 0;
        let allPostsLoaded = false;

        // Создаем модальное окно
        const postModal = document.createElement("div");
        postModal.classList.add("post-modal");
        postModal.style.display = "none";
        document.body.appendChild(postModal);

        async function loadUserPostsPaginated() {
            if (allPostsLoaded) return;

            try {
                const posts = await loadUserPosts(userId, limit, offset);
                if (!posts || posts.length === 0) {
                    allPostsLoaded = true;
                    return;
                }

                renderUserPosts(posts);
                offset += posts.length;

                if (posts.length < limit) {
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
                postElement.dataset.postId = post.id;
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
                contentContainer.appendChild(postElement);
            });

            postsHeader.style.opacity = "1";
        }

        async function openPost(postId) {
            try {
                const response = await fetch(`/api/recipes/${postId}/`);
                if (!response.ok) throw new Error("Ошибка загрузки поста");
                const post = await response.json();

                postModal.innerHTML = `
                <div class="post-modal-overlay"></div>
                <div class="post-modal-content">
                    <button class="close-modal">&times;</button>
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
                </div>
            `;
                postModal.style.display = "flex";
                document.body.style.overflow = "hidden";

                postModal.querySelector(".close-modal").addEventListener("click", closePostModal);
                postModal.querySelector(".post-modal-overlay").addEventListener("click", closePostModal);
            } catch (error) {
                console.error("Ошибка загрузки поста:", error);
            }
        }

        function closePostModal() {
            postModal.style.display = "none";
            document.body.style.overflow = "";
        }

        document.addEventListener("click", async (event) => {
            const likeButton = event.target.closest(".like-button");
            if (likeButton) {
                const postId = likeButton.dataset.postId;
                const likeIcon = likeButton.querySelector(".like-icon");
                if (!postId || !likeIcon) return;

                try {
                    const response = await fetch(`/api/recipes/${postId}/like/`, {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        credentials: "include",
                    });

                    if (response.ok) {
                        const isLiked = likeIcon.src.includes("like_color.svg");
                        likeIcon.src = isLiked ? "/static/images/general/icon/like_black.svg" : "/static/images/general/icon/like_color.svg";
                    } else {
                        throw new Error("Ошибка при лайке поста");
                    }
                } catch (error) {
                    console.error("Ошибка лайка:", error);
                }
                return;
            }

            const postElement = event.target.closest(".post");
            if (postElement) {
                const postId = postElement.dataset.postId;
                if (postId) openPost(postId);
            }
        });

        window.addEventListener("scroll", () => {
            if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
                loadUserPostsPaginated();
            }
        });

        await loadUserPostsPaginated();

        const editProfileButton = document.getElementById("edit-profile");
        if (editProfileButton) {
            editProfileButton.addEventListener("click", () => {
                window.location.href = "/my_user/update/";
            });
        }
    } catch (error) {
        console.error("Ошибка инициализации страницы:", error);
    }
});
