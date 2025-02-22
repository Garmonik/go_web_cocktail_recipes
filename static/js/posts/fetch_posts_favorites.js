document.addEventListener("DOMContentLoaded", async function () {
    const contentContainer = document.querySelector("#posts-container");
    const newPostsBtn = document.getElementById("new-posts-btn");
    const popularPostsBtn = document.getElementById("popular-posts-btn");
    const postsHeader = document.querySelector(".posts-header");

    let limit = 10;
    let offset = 0;
    let orderBy = "created";
    let allPostsLoaded = false;

    async function fetchPosts() {
        if (allPostsLoaded) return;

        try {
            const response = await fetch(`/api/recipes/like/list/?limit=${limit}&offset=${offset}&order_by=${orderBy}`);
            if (!response.ok) {
                throw new Error("Ошибка при загрузке постов");
            }
            const data = await response.json();

            if (data.content.length === 0) {
                allPostsLoaded = true;
                return;
            }

            renderPosts(data.content);

            if (data.content.length < limit) {
                allPostsLoaded = true;
            }
        } catch (error) {
            console.error("Ошибка:", error);
        }
    }

    function renderPosts(posts) {
        if (offset === 0) {
            contentContainer.innerHTML = "";
        }

        posts.forEach(post => {
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.innerHTML = `
                <div class="post-header">
                    <img src="data:image/png;base64,${post.author.avatar}" 
                         alt="Avatar" 
                         class="avatar user-link" 
                         data-user-id="${post.author.id}">
                    <span class="username user-link" data-user-id="${post.author.id}">
                        ${post.author.username}
                    </span>
                </div>
                <div class="post-content">
                    <h3>${post.name}</h3>
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
            contentContainer.appendChild(postElement);
        });

        postsHeader.style.opacity = "1";
    }

    contentContainer.addEventListener("click", (event) => {
        const target = event.target;
        if (target.classList.contains("user-link")) {
            const userId = target.dataset.userId;
            if (userId) {
                window.location.href = `/user/${userId}/`;
            }
        }
    });

    contentContainer.addEventListener("click", async (event) => {
        const likeButton = event.target.closest(".like-button");
        if (!likeButton) return;

        const postId = likeButton.dataset.postId;
        const likeIcon = likeButton.querySelector(".like-icon");
        if (!postId || !likeIcon) return;

        const isLiked = likeIcon.src.includes("like_color.svg");

        likeIcon.src = `/static/images/general/icon/${isLiked ? "like_black.svg" : "like_color.svg"}`;

        try {
            const response = await fetch(`/api/recipes/${postId}/like/`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
            });

            if (response.status === 201) {
                likeIcon.src = "/static/images/general/icon/like_color.svg";
            } else if (response.status === 204) {
                likeIcon.src = "/static/images/general/icon/like_black.svg";
            } else {
                throw new Error(`Неожиданный статус: ${response.status}`);
            }
        } catch (error) {
            console.error("Ошибка лайка:", error);
            likeIcon.src = `/static/images/general/icon/${isLiked ? "like_color.svg" : "like_black.svg"}`;
        }
    });

    newPostsBtn.addEventListener("click", () => {
        orderBy = "created";
        offset = 0;
        allPostsLoaded = false;
        newPostsBtn.classList.add("active");
        popularPostsBtn.classList.remove("active");
        fetchPosts();
    });

    popularPostsBtn.addEventListener("click", () => {
        orderBy = "popular";
        offset = 0;
        allPostsLoaded = false;
        popularPostsBtn.classList.add("active");
        newPostsBtn.classList.remove("active");
        fetchPosts();
    });

    window.addEventListener("scroll", () => {
        if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
            offset += limit;
            fetchPosts();
        }
    });

    await fetchPosts();
});
