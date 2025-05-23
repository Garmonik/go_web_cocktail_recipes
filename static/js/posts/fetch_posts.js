document.addEventListener("DOMContentLoaded", async function () {
    const contentContainer = document.querySelector("#posts-container");
    const newPostsBtn = document.getElementById("new-posts-btn");
    const popularPostsBtn = document.getElementById("popular-posts-btn");
    const postsHeader = document.querySelector(".posts-header");

    let limit = 10;
    let offset = 0;
    let orderBy = "created";
    let allPostsLoaded = false;

    // Создаём модальное окно для просмотра поста
    const postModal = document.createElement("div");
    postModal.classList.add("post-modal");
    postModal.style.display = "none";
    document.body.appendChild(postModal);

    async function fetchPosts() {
        if (allPostsLoaded) return;

        try {
            const response = await fetch(`/api/recipes/?limit=${limit}&offset=${offset}&order_by=${orderBy}`);
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
            postElement.dataset.postId = post.id;
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

    contentContainer.addEventListener("click", async (event) => {
        const likeButton = event.target.closest(".like-button");
        if (likeButton) {
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
            return;
        }

        const postElement = event.target.closest(".post");
        if (!postElement) return;

        const postId = postElement.dataset.postId;
        if (!postId) return;

        openPost(postId);
    });

    async function openPost(postId) {
        try {
            const response = await fetch(`/api/recipes/${postId}/`);
            if (!response.ok) {
                throw new Error("Ошибка загрузки поста");
            }
            const post = await response.json();

            postModal.innerHTML = `
            <div class="post-modal-overlay"></div>
            <div class="post-modal-content">
                <button class="close-modal">&times;</button>
                <div class="post-header">
                    <img src="data:image/png;base64,${post.author.avatar}" 
                         alt="Avatar" 
                         class="avatar">
                    <span class="username">${post.author.username}</span>
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
            </div>
        `;

            postModal.style.display = "flex";
            document.body.style.overflow = "hidden"; // Отключаем скролл

            const closeModalBtn = postModal.querySelector(".close-modal");
            closeModalBtn.addEventListener("click", closePostModal);

            postModal.querySelector(".post-modal-overlay").addEventListener("click", closePostModal);

            // Добавляем обработчик лайков внутри модального окна
            const likeButton = postModal.querySelector(".like-button");
            likeButton.addEventListener("click", async () => {
                const postId = likeButton.dataset.postId;
                const likeIcon = likeButton.querySelector(".like-icon");

                if (!postId || !likeIcon) return;

                const isLiked = likeIcon.src.includes("like_color.svg");

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
                }
            });

        } catch (error) {
            console.error("Ошибка загрузки поста:", error);
        }
    }

    function closePostModal() {
        postModal.style.display = "none";
        document.body.style.overflow = ""; // Включаем скролл обратно
    }

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

    contentContainer.addEventListener("click", (event) => {
        const userLink = event.target.closest(".user-link");
        if (!userLink) return;

        const userId = userLink.dataset.userId;
        if (userId) {
            window.location.href = `/user/${userId}/`;
        }
    });

    await fetchPosts();
});
