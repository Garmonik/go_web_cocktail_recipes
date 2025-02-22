import { loadSomeUserProfile, loadUserPosts } from "./user_api.js";

document.addEventListener("DOMContentLoaded", async function () {
    const contentContainer = document.getElementById("user-posts");
    const postsHeader = document.querySelector(".posts-header");
    let limit = 10;
    let offset = 0;

    async function loadAndRenderUserProfile() {
        const userId = document.body.dataset.userId;
        if (!userId) {
            console.error("Ошибка: ID пользователя не найден!");
            return;
        }

        await loadSomeUserProfile(userId);
        await loadUserPostsAndRender(userId);
    }

    async function loadUserPostsAndRender(userId) {
        try {
            const posts = await loadUserPosts(userId, limit, offset);
            renderUserPosts(posts);
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

        if (postsHeader) postsHeader.style.opacity = "1";
    }

    contentContainer.addEventListener("click", async (event) => {
        const likeButton = event.target.closest(".like-button");
        if (!likeButton) return;

        const postId = likeButton.dataset.postId;
        const likeIcon = likeButton.querySelector(".like-icon");
        const likeCount = likeButton.querySelector(".like-count");

        if (!postId || !likeIcon || !likeCount) return;

        try {
            const response = await fetch(`/api/post/${postId}/like/`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
            });
            if (!response.ok) throw new Error("Ошибка при лайке поста");

            const isLiked = likeIcon.src.includes("like_color.svg");
            likeIcon.src = isLiked ? "/static/images/general/icon/like_black.svg" : "/static/images/general/icon/like_color.svg";
            likeCount.textContent = isLiked ? Number(likeCount.textContent) - 1 : Number(likeCount.textContent) + 1;
        } catch (error) {
            console.error("Ошибка лайка:", error);
        }
    });

    window.addEventListener("scroll", () => {
        if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
            offset += limit;
            const userId = document.body.dataset.userId;
            if (userId) {
                loadUserPostsAndRender(userId);
            }
        }
    });

    await loadAndRenderUserProfile();

    document.getElementById("edit-profile").addEventListener("click", () => {
        window.location.href = "/my_user/update/";
    });
});
