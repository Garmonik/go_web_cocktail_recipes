document.addEventListener("DOMContentLoaded", async function () {
    const contentContainer = document.querySelector("#posts-container");
    const newPostsBtn = document.getElementById("new-posts-btn");
    const popularPostsBtn = document.getElementById("popular-posts-btn");
    const postsHeader = document.querySelector(".posts-header");

    let limit = 10;
    let offset = 0;
    let orderBy = "created"; // По умолчанию "Новые"

    async function fetchPosts() {
        try {
            const response = await fetch(`/api/recipes/?limit=${limit}&offset=${offset}&order_by=${orderBy}`);
            if (!response.ok) {
                throw new Error("Ошибка при загрузке постов");
            }
            const posts = await response.json();
            renderPosts(posts);
        } catch (error) {
            console.error("Ошибка:", error);
        }
    }

    function renderPosts(posts) {
        if (offset === 0) {
            contentContainer.innerHTML = ""; // Очищаем контейнер при новой загрузке
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
                        <img src="${post.like ? "/static/images/general/icon/like_color.svg" : "/static/images/general/icon/like_black.svg"}" 
                             alt="Like" class="like-icon">
                    </button>
                </div>
            `;
            contentContainer.appendChild(postElement);
        });

        // Показываем заголовок с кнопками, если он вдруг исчезал
        postsHeader.style.opacity = "1";
    }

    // Делегирование события для клика по аватарке и username
    contentContainer.addEventListener("click", (event) => {
        const target = event.target;
        if (target.classList.contains("user-link")) {
            const userId = target.dataset.userId;
            if (userId) {
                window.location.href = `/user/${userId}/`;
            }
        }
    });

    // Обработчик клика на лайк
    contentContainer.addEventListener("click", async (event) => {
        const likeButton = event.target.closest(".like-button");
        if (!likeButton) return;

        const postId = likeButton.dataset.postId;
        const likeIcon = likeButton.querySelector(".like-icon");

        if (!postId || !likeIcon) return;

        try {
            const response = await fetch(`/api/post/${postId}/like/`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
            });

            if (!response.ok) {
                throw new Error("Ошибка при лайке поста");
            }

            // Переключаем иконку лайка
            if (likeIcon.src.includes("like_black.svg")) {
                likeIcon.src = "/static/images/general/icon/like_color.svg";
            } else {
                likeIcon.src = "/static/images/general/icon/like_black.svg";
            }
        } catch (error) {
            console.error("Ошибка лайка:", error);
        }
    });

    newPostsBtn.addEventListener("click", () => {
        orderBy = "created";
        offset = 0;
        newPostsBtn.classList.add("active");
        popularPostsBtn.classList.remove("active");
        fetchPosts();
    });

    popularPostsBtn.addEventListener("click", () => {
        orderBy = "popular";
        offset = 0;
        popularPostsBtn.classList.add("active");
        newPostsBtn.classList.remove("active");
        fetchPosts();
    });

    // Функция для подгрузки постов при прокрутке
    window.addEventListener("scroll", () => {
        if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
            offset += limit;
            fetchPosts();
        }
    });

    await fetchPosts();
});
